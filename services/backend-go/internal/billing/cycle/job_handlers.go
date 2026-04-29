package cycle

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/billing"
	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/billing/repo"
	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/session"
	syncsvc "github.com/gustavo/5g-energia-fatura/services/backend-go/internal/sync"
)

// JobDeps holds dependencies for job handlers.
type JobDeps struct {
	Pool           *pgxpool.Pool
	ContractRepo   *repo.ContractRepo
	CalcRepo       *repo.CalculationRepo
	SyncService    *syncsvc.Service
	SessionManager *session.Manager
}

// NewJobDeps creates job dependencies.
func NewJobDeps(pool *pgxpool.Pool, syncSvc *syncsvc.Service, sessMgr *session.Manager) *JobDeps {
	return &JobDeps{
		Pool:           pool,
		ContractRepo:   repo.NewContractRepo(pool),
		CalcRepo:       repo.NewCalculationRepo(pool),
		SyncService:    syncSvc,
		SessionManager: sessMgr,
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
		input := billing.CalcInput{
			Contract: billing.CalcContract{
				FatorRepasseEnergia:               contract.FatorRepasseEnergia,
				ValorIPComDesconto:                contract.ValorIPComDesconto,
				BandeiraComDesconto:               contract.BandeiraComDesconto,
				CustoDisponibilidadeSempreCobrado: contract.CustoDisponibilidadeSempreCobrado,
				ConsumoMinimoKWh:                  contract.ConsumoMinimoKWh,
			},
			Itens: items,
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

// handleGeneratePDF gera PDF para um cálculo aprovado.
func handleGeneratePDF(deps *JobDeps) SyncJobHandler {
	return func(ctx context.Context, job *SyncJob) error {
		calcID, err := uuid.Parse(getString(job.Payload, "calculation_id"))
		if err != nil {
			return fmt.Errorf("invalid calculation_id in payload: %w", err)
		}
		cycleID, err := uuid.Parse(getString(job.Payload, "cycle_id"))
		if err != nil {
			return fmt.Errorf("invalid cycle_id: %w", err)
		}

		// 1. Buscar cálculo com resultado
		calc, err := deps.CalcRepo.GetByID(ctx, calcID)
		if err != nil {
			return fmt.Errorf("fetch calculation %s: %w", calcID, err)
		}

		// 2. Parse do result_snapshot_json → CalculationResult
		var result billing.CalcResult
		if err := json.Unmarshal(calc.ResultSnapshotJSON, &result); err != nil {
			return fmt.Errorf("unmarshal result_snapshot_json: %w", err)
		}

		// 3. Buscar UC + Customer + Cycle info
		var ucCode, address, city, uf, customerName string
		err = deps.Pool.QueryRow(ctx, `
			SELECT cu.uc_code, COALESCE(cu.endereco_unidade, ''), COALESCE(cu.cidade, ''),
			       COALESCE(cu.uf, ''), COALESCE(c.nome_razao, '')
			FROM public.consumer_unit cu
			LEFT JOIN public.customer c ON c.id = cu.customer_id
			WHERE cu.id = $1
		`, calc.ConsumerUnitID).Scan(&ucCode, &address, &city, &uf, &customerName)
		if err != nil {
			return fmt.Errorf("fetch consumer_unit %s: %w", calc.ConsumerUnitID, err)
		}

		var year, month int16
		err = deps.Pool.QueryRow(ctx, `
			SELECT year, month FROM public.billing_cycle WHERE id = $1
		`, cycleID).Scan(&year, &month)
		if err != nil {
			return fmt.Errorf("fetch billing_cycle %s: %w", cycleID, err)
		}

		// 4. Montar InvoicePDFData
		pdfData := &InvoicePDFData{
			CustomerName:     customerName,
			UCCode:           ucCode,
			Address:          address,
			City:             city,
			UF:               uf,
			ReferenceMonth:   fmt.Sprintf("%s/%d", monthName(month), year),
			IssueDate:        time.Now().Format("02/01/2006"),
			TotalSemDesconto: formatDecimal(calc.TotalSemDesconto),
			TotalComDesconto: formatDecimal(calc.TotalComDesconto),
			EconomiaRS:       formatDecimal(calc.EconomiaRS),
			EconomiaPct:      formatDecimal(calc.EconomiaPct),
		}

		for _, line := range result.Linhas {
			pdfData.Lines = append(pdfData.Lines, InvoicePDFLine{
				Label:         line.Label,
				Quantidade:    formatDecimal(line.Quantidade),
				PrecoUnitario: formatDecimal(line.PrecoUnitario),
				ValorSemDesc:  formatDecimal(line.ValorSemDesc),
				ValorComDesc:  formatDecimal(line.ValorComDesc),
			})
		}

		// 5. Determinar diretório de saída
		outputDir := os.Getenv("PDF_OUTPUT_DIR")
		if outputDir == "" {
			outputDir = "./pdfs"
		}

		// 6. Gerar PDF (em buffer) + salvar em disco
		pdfResult, err := GenerateAndSaveInvoicePDF(pdfData, outputDir)
		if err != nil {
			return fmt.Errorf("generate pdf: %w", err)
		}

		// 7. Inserir registro em generated_document
		_, err = deps.Pool.Exec(ctx, `
			INSERT INTO public.generated_document
				(billing_calculation_id, type, file_path, checksum_sha256, version)
			VALUES ($1, 'customer_invoice_pdf', $2, $3, 1)
		`, calcID, pdfResult.FilePath, pdfResult.Checksum)
		if err != nil {
			return fmt.Errorf("insert generated_document: %w", err)
		}

		// 8. Atualizar cycle_consumer_unit
		_, err = deps.Pool.Exec(ctx, `
			UPDATE public.cycle_consumer_unit
			SET pdf_generated = true, updated_at = NOW()
			WHERE billing_cycle_id = $1 AND consumer_unit_id = $2
		`, cycleID, calc.ConsumerUnitID)
		if err != nil {
			return fmt.Errorf("update cycle_consumer_unit pdf_generated: %w", err)
		}

		return nil
	}
}

// handleSyncUC executa sync real para uma UC usando credenciais + API Neoenergia.
func handleSyncUC(deps *JobDeps) SyncJobHandler {
	return func(ctx context.Context, job *SyncJob) error {
		cycleID, err := uuid.Parse(getString(job.Payload, "cycle_id"))
		if err != nil {
			return fmt.Errorf("cycle_id inválido: %w", err)
		}
		ucID, err := uuid.Parse(getString(job.Payload, "uc_id"))
		if err != nil {
			return fmt.Errorf("uc_id inválido: %w", err)
		}
		ucCode := getString(job.Payload, "uc_code")
		if ucCode == "" {
			return fmt.Errorf("uc_code é obrigatório no payload")
		}

		// 1. Buscar sync_credential_id da consumer_unit
		var syncCredentialID *string
		err = deps.Pool.QueryRow(ctx, `
			SELECT sync_credential_id FROM public.consumer_unit WHERE id = $1
		`, ucID).Scan(&syncCredentialID)
		if err != nil {
			return fmt.Errorf("consumer_unit %s não encontrada: %w", ucID, err)
		}

		// 2. Determinar o credential_id a usar
		credentialID, err := resolveCredentialID(ctx, deps.Pool, ucID, ucCode, syncCredentialID)
		if err != nil {
			return fmt.Errorf("resolver credencial para UC %s: %w", ucCode, err)
		}

		// 3. Resolver sessão (token bearer + documento)
		resolved, err := deps.SessionManager.ResolveToken(ctx, credentialID)
		if err != nil {
			return fmt.Errorf("falha ao obter token de acesso para UC %s: %w", ucCode, err)
		}

		// 4. Executar sync completo via sync service
		syncResult := deps.SyncService.SyncUC(ctx, syncsvc.SyncUCRequest{
			BearerToken:       resolved.Token,
			Documento:         resolved.Documento,
			CredentialID:      credentialID,
			UC:                ucCode,
			IncludePDF:        true,
			IncludeExtraction: false,
		})

		// 5. Verificar resultado do sync
		if syncResult.BillingRecord == nil {
			errMsg := "Nenhuma fatura encontrada para UC " + ucCode
			if syncResult.Persistence != nil && syncResult.Persistence.Error != "" {
				errMsg = syncResult.Persistence.Error
			}
			// Atualizar status para error
			_, updateErr := deps.Pool.Exec(ctx, `
				UPDATE public.cycle_consumer_unit
				SET status = 'error', error_message = $1, updated_at = NOW()
				WHERE billing_cycle_id = $2 AND consumer_unit_id = $3
			`, errMsg, cycleID, ucID)
			if updateErr != nil {
				return fmt.Errorf("sync falhou para UC %s: %s (erro ao atualizar status: %s)", ucCode, errMsg, updateErr)
			}
			return fmt.Errorf("sync falhou para UC %s: %s", ucCode, errMsg)
		}

		// 6. Extrair dados do resultado sync
		billingRecordJSON, err := json.Marshal(syncResult.BillingRecord)
		if err != nil {
			return fmt.Errorf("erro ao serializar billing record da UC %s: %w", ucCode, err)
		}

		invoiceID := ""
		syncRunID := ""
		if syncResult.Persistence != nil {
			invoiceID = syncResult.Persistence.InvoiceID
			syncRunID = syncResult.Persistence.SyncRunID
		}

		// 7. Criar/atualizar utility_invoice_ref
		var nfArg, mrArg any
		if syncResult.BillingRecord.NumeroFatura != "" {
			nfArg = syncResult.BillingRecord.NumeroFatura
		}
		if syncResult.BillingRecord.MesReferencia != "" {
			mrArg = syncResult.BillingRecord.MesReferencia
		}
		var srArg any
		if syncRunID != "" {
			srArg = syncRunID
		}

		_, err = deps.Pool.Exec(ctx, `
			INSERT INTO public.utility_invoice_ref (
				id, consumer_unit_id, billing_cycle_id, sync_invoice_id, sync_run_id,
				numero_fatura, mes_referencia, billing_record_snapshot, synced_at,
				created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW(), NOW())
			ON CONFLICT (consumer_unit_id, billing_cycle_id) DO UPDATE SET
				sync_invoice_id = EXCLUDED.sync_invoice_id,
				sync_run_id = EXCLUDED.sync_run_id,
				numero_fatura = EXCLUDED.numero_fatura,
				mes_referencia = EXCLUDED.mes_referencia,
				billing_record_snapshot = EXCLUDED.billing_record_snapshot,
				synced_at = EXCLUDED.synced_at,
				updated_at = NOW()
		`, uuid.New(), ucID, cycleID, invoiceID, srArg,
			nfArg, mrArg, billingRecordJSON)
		if err != nil {
			return fmt.Errorf("erro ao salvar utility_invoice_ref para UC %s: %w", ucCode, err)
		}

		// 8. Marcar UC como synced no ciclo
		_, err = deps.Pool.Exec(ctx, `
			UPDATE public.cycle_consumer_unit
			SET status = 'synced', synced_at = NOW(), updated_at = NOW()
			WHERE billing_cycle_id = $1 AND consumer_unit_id = $2
		`, cycleID, ucID)
		if err != nil {
			return fmt.Errorf("erro ao atualizar status da UC %s no ciclo: %w", ucCode, err)
		}

		return nil
	}
}

// resolveCredentialID busca o credential_id para uma UC, tentando fontes em ordem.
func resolveCredentialID(ctx context.Context, pool *pgxpool.Pool, ucID uuid.UUID, ucCode string, syncCredentialID *string) (string, error) {
	// Fonte 1: sync_credential_id direto na consumer_unit
	if syncCredentialID != nil && *syncCredentialID != "" {
		return *syncCredentialID, nil
	}

	// Fonte 2: integration_consumer_units (mapeamento UC -> credential)
	var credentialID string
	err := pool.QueryRow(ctx, `
		SELECT credential_id::text FROM public.integration_consumer_units
		WHERE uc = $1 AND credential_id IS NOT NULL
	`, ucCode).Scan(&credentialID)
	if err == nil {
		return credentialID, nil
	}

	// Fonte 3: credential_link via customer da UC
	err = pool.QueryRow(ctx, `
		SELECT cl.go_credential_id
		FROM public.credential_link cl
		JOIN public.consumer_unit cu ON cu.customer_id = cl.customer_id
		WHERE cu.id = $1
		LIMIT 1
	`, ucID).Scan(&credentialID)
	if err == nil {
		return credentialID, nil
	}

	return "", fmt.Errorf("nenhuma credencial de integração encontrada")
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
			Tipo:          itemType,
			Descricao:     desc,
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
	case strings.Contains(d, "INJETADA") || strings.Contains(d, "INJEÇÃO") || strings.Contains(d, "GERAÇÃO") || strings.Contains(d, "SCEE") || strings.Contains(d, "COMPENSA"):
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

// formatDecimal converte decimal.Decimal para string no formato brasileiro
// (ex: "1.234,56"). Usado para exibição em PDFs.
func formatDecimal(d decimal.Decimal) string {
	if d.IsZero() {
		return "0,00"
	}

	// Separar parte inteira e centavos
	// Multiplica por 100 e trunca para obter o valor em centavos
	cents := d.Mul(decimal.NewFromInt(100)).IntPart()
	signal := ""
	if cents < 0 {
		signal = "-"
		cents = -cents
	}

	intPart := cents / 100
	fracPart := cents % 100

	// Formatar parte inteira com separador de milhar (ponto)
	intStr := fmt.Sprintf("%d", intPart)
	var parts []string
	for i := len(intStr); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		parts = append([]string{intStr[start:i]}, parts...)
	}
	formattedInt := strings.Join(parts, ".")

	return fmt.Sprintf("%s%s,%02d", signal, formattedInt, fracPart)
}
