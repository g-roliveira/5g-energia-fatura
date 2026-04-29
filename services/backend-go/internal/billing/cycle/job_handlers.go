package cycle

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/billing"
	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/billing/repo"
)

// JobDeps holds dependencies for job handlers.
type JobDeps struct {
	Pool        *pgxpool.Pool
	ContractRepo *repo.ContractRepo
	CalcRepo    *repo.CalculationRepo
}

// NewJobDeps creates job dependencies.
func NewJobDeps(pool *pgxpool.Pool) *JobDeps {
	return &JobDeps{
		Pool:         pool,
		ContractRepo: repo.NewContractRepo(pool),
		CalcRepo:     repo.NewCalculationRepo(pool),
	}
}

// BuildHandlers registra todos os handlers no worker pool.
func BuildHandlers(pool *WorkerPool, deps *JobDeps) {
	pool.RegisterHandler("calculate", handleCalculate(deps))
	pool.RegisterHandler("approve", handleApprove(deps))
	pool.RegisterHandler("generate_pdf", handleGeneratePDF(deps))
	pool.RegisterHandler("sync_uc", handleSyncUC(deps))
}

// handleCalculate processa um job de cálculo de fatura.
func handleCalculate(deps *JobDeps) SyncJobHandler {
	return func(ctx context.Context, job *SyncJob) error {
		cycleID, err := uuid.Parse(getString(job.Payload, "cycle_id"))
		if err != nil {
			return fmt.Errorf("invalid cycle_id: %w", err)
		}
		ucID, err := uuid.Parse(getString(job.Payload, "uc_id"))
		if err != nil {
			return fmt.Errorf("invalid uc_id: %w", err)
		}

		// 1. Buscar utility_invoice_ref para esta UC no ciclo
		var refID uuid.UUID
		var snapshot []byte
		err = deps.Pool.QueryRow(ctx, `
			SELECT id, billing_record_snapshot
			FROM public.utility_invoice_ref
			WHERE billing_cycle_id = $1 AND consumer_unit_id = $2
			ORDER BY synced_at DESC
			LIMIT 1
		`, cycleID, ucID).Scan(&refID, &snapshot)
		if err != nil {
			return fmt.Errorf("no invoice ref for uc %s in cycle %s: %w", ucID, cycleID, err)
		}

		// 2. Buscar contrato ativo
		contract, err := deps.ContractRepo.GetActiveByConsumerUnit(ctx, ucID)
		if err != nil {
			return fmt.Errorf("no active contract for uc %s: %w", ucID, err)
		}

		// 3. Extrair itens do snapshot
		items, err := extractItemsFromSnapshot(snapshot)
		if err != nil {
			return fmt.Errorf("extract items: %w", err)
		}
		if len(items) == 0 {
			return fmt.Errorf("no items extracted from invoice ref %s", refID)
		}

		// 4. Montar input do motor
		input := billing.CalculationInput{
			Contract: billing.CalcContract{
				DescontoPct:                       contract.DescontoPercentual,
				IPFaturamentoMode:                 billing.IPFaturamentoMode(contract.IPFaturamentoMode),
				IPFaturamentoValor:                contract.IPFaturamentoValor,
				IPFaturamentoPct:                  contract.IPFaturamentoPercent,
				BandeiraComDesconto:               contract.BandeiraComDesconto,
				CustoDisponibilidadeSempreCobrado: contract.CustoDisponibilidadeSempreCobrado,
			},
			Itens:            items,
			ConsumoMinimoKWh: 30, // default, pode vir do contrato no futuro
		}

		// 5. Calcular
		result, err := billing.Calculate(input)
		if err != nil {
			return fmt.Errorf("calculation failed: %w", err)
		}

		// 6. Snapshots JSON
		contractSnap, _ := json.Marshal(contract)
		inputsSnap, _ := json.Marshal(input)
		resultSnap, _ := json.Marshal(result)

		// 7. Persistir em transação
		tx, err := deps.Pool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			return fmt.Errorf("begin tx: %w", err)
		}
		defer tx.Rollback(ctx)

		version, err := deps.CalcRepo.NextVersionFor(ctx, tx, refID)
		if err != nil {
			return fmt.Errorf("next version: %w", err)
		}

		// Marcar versões anteriores como superseded
		if err := deps.CalcRepo.MarkSuperseded(ctx, tx, refID); err != nil {
			return fmt.Errorf("mark superseded: %w", err)
		}

		status := repo.CalcStatusDraft
		needsReview := []string{}
		if len(result.Warnings) > 0 {
			status = repo.CalcStatusNeedsReview
			needsReview = result.Warnings
		}

		calc := &repo.BillingCalculation{
			UtilityInvoiceRefID:  refID,
			BillingCycleID:       cycleID,
			ConsumerUnitID:       ucID,
			ContractID:           contract.ID,
			ContractSnapshotJSON: contractSnap,
			InputsSnapshotJSON:   inputsSnap,
			ResultSnapshotJSON:   resultSnap,
			TotalSemDesconto:     result.TotalSemDesconto,
			TotalComDesconto:     result.TotalComDesconto,
			EconomiaRS:           result.EconomiaRS,
			EconomiaPct:          result.EconomiaPct,
			Status:               status,
			NeedsReviewReasons:   needsReview,
			Version:              version,
			CalculatedAt:         time.Now(),
		}

		if err := deps.CalcRepo.Insert(ctx, tx, calc); err != nil {
			return fmt.Errorf("insert calculation: %w", err)
		}

		// Atualizar status do cycle_consumer_unit
		_, err = tx.Exec(ctx, `
			UPDATE public.cycle_consumer_unit
			SET status = 'calculated', calculation_id = $1, updated_at = NOW()
			WHERE billing_cycle_id = $2 AND consumer_unit_id = $3
		`, calc.ID, cycleID, ucID)
		if err != nil {
			return fmt.Errorf("update cycle_consumer_unit: %w", err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit tx: %w", err)
		}

		return nil
	}
}

// handleApprove aprova o cálculo mais recente de uma UC no ciclo.
func handleApprove(deps *JobDeps) SyncJobHandler {
	return func(ctx context.Context, job *SyncJob) error {
		cycleID, err := uuid.Parse(getString(job.Payload, "cycle_id"))
		if err != nil {
			return fmt.Errorf("invalid cycle_id: %w", err)
		}
		ucID, err := uuid.Parse(getString(job.Payload, "uc_id"))
		if err != nil {
			return fmt.Errorf("invalid uc_id: %w", err)
		}

		// Buscar cálculo atual
		var calcID uuid.UUID
		err = deps.Pool.QueryRow(ctx, `
			SELECT id FROM public.billing_calculation
			WHERE billing_cycle_id = $1 AND consumer_unit_id = $2 AND status != 'superseded'
			ORDER BY version DESC
			LIMIT 1
		`, cycleID, ucID).Scan(&calcID)
		if err != nil {
			return fmt.Errorf("no calculation found for uc %s in cycle %s: %w", ucID, cycleID, err)
		}

		approverID := uuid.Nil
		if s := getString(job.Payload, "approver_id"); s != "" {
			approverID, _ = uuid.Parse(s)
		}

		if err := deps.CalcRepo.Approve(ctx, calcID, approverID); err != nil {
			return fmt.Errorf("approve calculation: %w", err)
		}

		// Atualizar status do cycle_consumer_unit
		_, err = deps.Pool.Exec(ctx, `
			UPDATE public.cycle_consumer_unit
			SET status = 'approved', updated_at = NOW()
			WHERE billing_cycle_id = $1 AND consumer_unit_id = $2
		`, cycleID, ucID)
		if err != nil {
			return fmt.Errorf("update cycle_consumer_unit: %w", err)
		}

		return nil
	}
}

// handleGeneratePDF gera PDF para um cálculo aprovado (stub).
func handleGeneratePDF(deps *JobDeps) SyncJobHandler {
	return func(ctx context.Context, job *SyncJob) error {
		cycleID, _ := uuid.Parse(getString(job.Payload, "cycle_id"))
		ucCode := getString(job.Payload, "uc_code")

		// TODO: implementar geração real de PDF usando chromedp ou similar
		// Por enquanto, apenas marca como gerado no cycle_consumer_unit
		_, err := deps.Pool.Exec(ctx, `
			UPDATE public.cycle_consumer_unit
			SET pdf_generated = true, updated_at = NOW()
			WHERE billing_cycle_id = $1 AND uc_code = $2
		`, cycleID, ucCode)
		if err != nil {
			return fmt.Errorf("mark pdf generated: %w", err)
		}

		return nil
	}
}

// handleSyncUC executa sync para uma UC (stub — delega para integration worker).
func handleSyncUC(deps *JobDeps) SyncJobHandler {
	return func(ctx context.Context, job *SyncJob) error {
		cycleID, _ := uuid.Parse(getString(job.Payload, "cycle_id"))
		ucID, _ := uuid.Parse(getString(job.Payload, "uc_id"))
		ucCode := getString(job.Payload, "uc_code")

		// TODO: implementar sync real usando credenciais + Playwright
		// Por enquanto, apenas marca como synced
		_, err := deps.Pool.Exec(ctx, `
			UPDATE public.cycle_consumer_unit
			SET status = 'synced', synced_at = NOW(), updated_at = NOW()
			WHERE billing_cycle_id = $1 AND consumer_unit_id = $2
		`, cycleID, ucID)
		if err != nil {
			return fmt.Errorf("mark synced: %w", err)
		}

		// Criar um utility_invoice_ref stub para permitir cálculo
		_, err = deps.Pool.Exec(ctx, `
			INSERT INTO public.utility_invoice_ref (
				id, consumer_unit_id, billing_cycle_id, numero_fatura,
				mes_referencia, billing_record_snapshot, synced_at, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW(), NOW())
			ON CONFLICT DO NOTHING
		`, uuid.New(), ucID, cycleID, "STUB-"+ucCode, "2024/01",
			[]byte(`{"itens_fatura":[]}`),
		)
		if err != nil {
			return fmt.Errorf("insert stub invoice ref: %w", err)
		}

		return nil
	}
}

// -------------------------------------------------------------------
// Helpers
// -------------------------------------------------------------------

func getString(m map[string]any, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case fmt.Stringer:
		return t.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

// extractItemsFromSnapshot extrai UtilityInvoiceItem[] do billing_record_snapshot JSONB.
func extractItemsFromSnapshot(snapshot []byte) ([]billing.UtilityInvoiceItem, error) {
	if len(snapshot) == 0 {
		return nil, fmt.Errorf("empty snapshot")
	}

	var record map[string]any
	if err := json.Unmarshal(snapshot, &record); err != nil {
		return nil, fmt.Errorf("unmarshal snapshot: %w", err)
	}

	// O snapshot pode estar embrulhado em "fields" (formato do extractor)
	fields, ok := record["fields"].(map[string]any)
	if ok {
		record = fields
	}

	rawItems, ok := record["itens_fatura"].([]any)
	if !ok {
		return nil, fmt.Errorf("itens_fatura not found or not array")
	}

	items := make([]billing.UtilityInvoiceItem, 0, len(rawItems))
	for _, raw := range rawItems {
		m, ok := raw.(map[string]any)
		if !ok {
			continue
		}

		desc := getString(m, "descricao")
		itemType := classifyItem(desc)
		if itemType == "" {
			continue // item não reconhecido, ignora
		}

		qtd := parseDecimal(getString(m, "quantidade"))
		tarifa := parseDecimal(getString(m, "tarifa"))
		valor := parseDecimal(getString(m, "valor"))
		valorTotal := parseDecimal(getString(m, "valor_total"))
		if valorTotal.IsZero() && !valor.IsZero() {
			valorTotal = valor
		}

		// Para itens de valor fixo (bandeira, IP), preço unitário é zero
		pu := tarifa
		if itemType == billing.ItemBandeira || itemType == billing.ItemIPCoelba {
			pu = decimal.Zero
			qtd = decimal.NewFromInt(1)
		}

		items = append(items, billing.UtilityInvoiceItem{
			Type:          itemType,
			Description:   desc,
			Quantidade:    qtd,
			PrecoUnitario: pu,
			ValorTotal:    valorTotal,
		})
	}

	return items, nil
}

// classifyItem mapeia descrição do item da fatura para ItemType.
func classifyItem(desc string) billing.ItemType {
	d := strings.ToUpper(strings.TrimSpace(desc))
	switch {
	case strings.Contains(d, "TUSD") && !strings.Contains(d, "TE"):
		return billing.ItemTUSDFio
	case strings.Contains(d, "TE") && !strings.Contains(d, "TUSD"):
		return billing.ItemTUSDEnergia
	case strings.Contains(d, "BANDEIRA"):
		return billing.ItemBandeira
	case strings.Contains(d, "ILUM") || strings.Contains(d, "IP") || strings.Contains(d, "PÚBLICA"):
		return billing.ItemIPCoelba
	case strings.Contains(d, "INJETADA") || strings.Contains(d, "INJEÇÃO") || strings.Contains(d, "GERAÇÃO"):
		return billing.ItemEnergiaInjetada
	case strings.Contains(d, "REATIVO"):
		return billing.ItemReativoExcedente
	case strings.Contains(d, "TRIBF") || strings.Contains(d, "IRRF") || strings.Contains(d, "TRIBUT"):
		return billing.ItemTributoRetido
	default:
		return "" // ignora itens não classificados
	}
}

// parseDecimal converte string brasileira (vírgula como decimal) para decimal.Decimal.
func parseDecimal(s string) decimal.Decimal {
	s = strings.TrimSpace(s)
	if s == "" {
		return decimal.Zero
	}
	// Remove pontos de milhar, substitui vírgula por ponto
	s = strings.ReplaceAll(s, ".", "")
	s = strings.Replace(s, ",", ".", 1)
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return d
}
