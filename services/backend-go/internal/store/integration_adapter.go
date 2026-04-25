package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/integration"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// IntegrationPostgresAdapter implementa IntegrationStore usando o novo domínio
// integration.* (schema integration, pgxpool).
type IntegrationPostgresAdapter struct {
	pool  *pgxpool.Pool
	store integration.Store
}

// OpenIntegrationPostgres abre uma conexão ao Postgres e retorna um IntegrationStore
// que usa o schema integration.*.
func OpenIntegrationPostgres(dsn string) (IntegrationStore, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	return &IntegrationPostgresAdapter{
		pool:  pool,
		store: integration.NewStore(pool),
	}, nil
}

func (a *IntegrationPostgresAdapter) Close() error {
	a.pool.Close()
	return nil
}

// --- Credentials ---

func (a *IntegrationPostgresAdapter) InsertCredential(rec CredentialRecord) (CredentialRecord, error) {
	id := rec.ID
	if id == "" {
		id = uuid.New().String()
	}
	c := &integration.Credential{
		ID:              uuid.MustParse(id),
		Label:           rec.Label,
		DocumentoCipher: rec.DocumentoCipher,
		DocumentoNonce:  rec.DocumentoNonce,
		SenhaCipher:     rec.SenhaCipher,
		SenhaNonce:      rec.SenhaNonce,
		UF:              rec.UF,
		TipoAcesso:      rec.TipoAcesso,
		KeyVersion:      rec.KeyVersion,
	}
	if err := a.store.InsertCredential(context.Background(), c); err != nil {
		return rec, err
	}
	rec.ID = c.ID.String()
	rec.CreatedAt = c.CreatedAt.Format(time.RFC3339)
	rec.UpdatedAt = c.UpdatedAt.Format(time.RFC3339)
	return rec, nil
}

func (a *IntegrationPostgresAdapter) GetCredentialByID(id string) (CredentialRecord, error) {
	c, err := a.store.GetCredentialByID(context.Background(), id)
	if err != nil {
		return CredentialRecord{}, err
	}
	return CredentialRecord{
		ID:              c.ID.String(),
		Label:           c.Label,
		DocumentoCipher: c.DocumentoCipher,
		DocumentoNonce:  c.DocumentoNonce,
		SenhaCipher:     c.SenhaCipher,
		SenhaNonce:      c.SenhaNonce,
		UF:              c.UF,
		TipoAcesso:      c.TipoAcesso,
		KeyVersion:      c.KeyVersion,
		CreatedAt:       c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       c.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// --- Sessions ---

func (a *IntegrationPostgresAdapter) InsertSession(rec SessionRecord) (SessionRecord, error) {
	id := rec.ID
	if id == "" {
		id = uuid.New().String()
	}
	sess := &integration.Session{
		ID:                uuid.MustParse(id),
		CredentialID:      uuid.MustParse(rec.CredentialID),
		BearerTokenCipher: rec.BearerTokenCipher,
		BearerTokenNonce:  rec.BearerTokenNonce,
	}
	if err := a.store.InsertSession(context.Background(), sess); err != nil {
		return rec, err
	}
	rec.ID = sess.ID.String()
	rec.CreatedAt = sess.CreatedAt.Format(time.RFC3339)
	rec.UpdatedAt = sess.UpdatedAt.Format(time.RFC3339)
	return rec, nil
}

func (a *IntegrationPostgresAdapter) GetLatestSessionByCredentialID(credentialID string) (SessionRecord, error) {
	sess, err := a.store.GetLatestSessionByCredentialID(context.Background(), credentialID)
	if err != nil {
		return SessionRecord{}, err
	}
	return SessionRecord{
		ID:                sess.ID.String(),
		CredentialID:      sess.CredentialID.String(),
		BearerTokenCipher: sess.BearerTokenCipher,
		BearerTokenNonce:  sess.BearerTokenNonce,
		CreatedAt:         sess.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         sess.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// --- Consumer Units ---

func (a *IntegrationPostgresAdapter) ListConsumerUnits(limit int, status string) ([]ConsumerUnitView, error) {
	units, err := a.store.ListConsumerUnits(context.Background(), limit, status)
	if err != nil {
		return nil, err
	}
	views := make([]ConsumerUnitView, len(units))
	for i, u := range units {
		views[i] = toConsumerUnitView(u)
	}
	return views, nil
}

func (a *IntegrationPostgresAdapter) GetConsumerUnitByUC(uc string) (*ConsumerUnitDetailsView, error) {
	u, err := a.store.GetConsumerUnitByUC(context.Background(), uc)
	if err != nil {
		return nil, err
	}
	view := &ConsumerUnitDetailsView{
		ConsumerUnitView: toConsumerUnitView(*u),
	}
	latestInvoice, err := a.store.GetLatestRawInvoiceByUC(context.Background(), uc)
	if err == nil && latestInvoice != nil {
		view.LatestInvoice = toInvoiceView(*latestInvoice)
	}
	latestSyncRun, err := a.store.GetLatestSyncRunByUC(context.Background(), uc)
	if err == nil && latestSyncRun != nil {
		view.LatestSyncRun = toSyncRunView(*latestSyncRun)
	}
	return view, nil
}

// --- Invoices ---

func (a *IntegrationPostgresAdapter) ListInvoicesByUC(uc string, limit int, status string) ([]InvoiceView, error) {
	invoices, err := a.store.ListRawInvoicesByUC(context.Background(), uc, limit)
	if err != nil {
		return nil, err
	}
	views := make([]InvoiceView, len(invoices))
	for i, inv := range invoices {
		views[i] = *toInvoiceView(inv)
	}
	return views, nil
}

func (a *IntegrationPostgresAdapter) GetInvoiceByID(id string) (*InvoiceView, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	inv, err := a.store.GetRawInvoiceByID(context.Background(), uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	view := toInvoiceView(*inv)
	// Items: buscar de integration.raw_invoice_items
	items, err := a.listInvoiceItems(id)
	if err != nil {
		return nil, err
	}
	view.Items = items
	return view, nil
}

func (a *IntegrationPostgresAdapter) GetLatestInvoiceByUC(uc string) (*InvoiceView, error) {
	inv, err := a.store.GetLatestRawInvoiceByUC(context.Background(), uc)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toInvoiceView(*inv), nil
}

func (a *IntegrationPostgresAdapter) listInvoiceItems(invoiceID string) ([]map[string]any, error) {
	query := `SELECT raw_json FROM integration.raw_invoice_items WHERE raw_invoice_id = $1 ORDER BY order_index`
	rows, err := a.pool.Query(context.Background(), query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []map[string]any
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}
		var item map[string]any
		if len(raw) > 0 {
			json.Unmarshal(raw, &item)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// --- Sync Runs ---

func (a *IntegrationPostgresAdapter) GetSyncRunByID(id string) (*SyncRunView, error) {
	sr, err := a.store.GetSyncRunByID(context.Background(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toSyncRunView(*sr), nil
}

func (a *IntegrationPostgresAdapter) GetLatestSyncRunByUC(uc string) (*SyncRunView, error) {
	sr, err := a.store.GetLatestSyncRunByUC(context.Background(), uc)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toSyncRunView(*sr), nil
}

// --- PersistSyncResult ---

func (a *IntegrationPostgresAdapter) PersistSyncResult(in PersistSyncInput) (PersistSyncResult, error) {
	ctx := context.Background()
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return PersistSyncResult{}, err
	}
	defer tx.Rollback(ctx)

	now := time.Now()
	syncRunID := in.SyncRunID
	if syncRunID == "" {
		syncRunID = uuid.New().String()
	}
	status := in.Status
	if status == "" {
		status = "succeeded"
	}

	var credentialID *uuid.UUID
	if in.CredentialID != "" {
		uid, _ := uuid.Parse(in.CredentialID)
		credentialID = &uid
	}

	// Insert sync_run
	rawJSON, _ := json.Marshal(in.RawResponse)
	errCtxJSON, _ := json.Marshal(map[string]any{"message": in.ErrorMessage})
	_, err = tx.Exec(ctx, `
		INSERT INTO integration.sync_runs (id, credential_id, documento, uc, status, step, error_message, error_context, raw_response_json, started_at, finished_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, syncRunID, credentialID, in.Documento, in.UC, status, "persist", nullableString(in.ErrorMessage), errCtxJSON, rawJSON, now, now, now)
	if err != nil {
		return PersistSyncResult{}, err
	}

	// Upsert consumer_unit
	if in.UCRecord != nil {
		enderecoJSON, _ := json.Marshal(in.UCRecord.Local)
		imovelJSON, _ := json.Marshal(in.Imovel)
		_, err = tx.Exec(ctx, `
			INSERT INTO integration.consumer_units (uc, credential_id, status, nome_cliente, instalacao, contrato, grupo_tensao, endereco_json, imovel_json, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $10)
			ON CONFLICT (uc) DO UPDATE SET
				credential_id = EXCLUDED.credential_id,
				status = EXCLUDED.status,
				nome_cliente = EXCLUDED.nome_cliente,
				instalacao = EXCLUDED.instalacao,
				contrato = EXCLUDED.contrato,
				grupo_tensao = EXCLUDED.grupo_tensao,
				endereco_json = EXCLUDED.endereco_json,
				imovel_json = EXCLUDED.imovel_json,
				updated_at = EXCLUDED.updated_at
		`, in.UCRecord.UC, credentialID, in.UCRecord.Status, in.UCRecord.NomeCliente, in.UCRecord.Instalacao, in.UCRecord.Contrato, in.UCRecord.GrupoTensao, enderecoJSON, imovelJSON, now)
		if err != nil {
			return PersistSyncResult{}, err
		}
	}

	if in.Fatura == nil {
		if err := tx.Commit(ctx); err != nil {
			return PersistSyncResult{}, err
		}
		return PersistSyncResult{SyncRunID: syncRunID}, nil
	}

	// Upsert raw_invoice
	billing := asMap(in.BillingRecord)
	completeness := asMap(billing["completeness"])
	missingArr := []string{}
	if m, ok := completeness["missing_fields"].([]any); ok {
		for _, v := range m {
			missingArr = append(missingArr, fmt.Sprint(v))
		}
	}
	codigoBarras := ""
	if in.DadosPagamento != nil {
		codigoBarras = in.DadosPagamento.CodBarras
	}

	var existingID string
	err = tx.QueryRow(ctx, `SELECT id::text FROM integration.raw_invoices WHERE uc = $1 AND numero_fatura = $2`, in.Fatura.UC, in.Fatura.NumeroFatura).Scan(&existingID)
	if err != nil && err != pgx.ErrNoRows {
		return PersistSyncResult{}, err
	}

	invoiceID := existingID
	if invoiceID == "" {
		invoiceID = uuid.New().String()
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO integration.raw_invoices (id, uc, numero_fatura, mes_referencia, status_fatura, valor_total, codigo_barras, data_emissao, data_vencimento, data_pagamento, data_inicio_periodo, data_fim_periodo, completeness_status, completeness_missing, billing_record_json, document_record_json, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $17)
		ON CONFLICT (uc, numero_fatura) DO UPDATE SET
			mes_referencia = EXCLUDED.mes_referencia,
			status_fatura = EXCLUDED.status_fatura,
			valor_total = EXCLUDED.valor_total,
			codigo_barras = EXCLUDED.codigo_barras,
			data_emissao = EXCLUDED.data_emissao,
			data_vencimento = EXCLUDED.data_vencimento,
			data_pagamento = EXCLUDED.data_pagamento,
			data_inicio_periodo = EXCLUDED.data_inicio_periodo,
			data_fim_periodo = EXCLUDED.data_fim_periodo,
			completeness_status = EXCLUDED.completeness_status,
			completeness_missing = EXCLUDED.completeness_missing,
			billing_record_json = EXCLUDED.billing_record_json,
			document_record_json = EXCLUDED.document_record_json,
			updated_at = EXCLUDED.updated_at
	`, invoiceID, firstNonEmpty(in.Fatura.UC, in.UC), in.Fatura.NumeroFatura, in.Fatura.MesReferencia,
		in.Fatura.StatusFatura, in.Fatura.ValorEmissao, codigoBarras,
		in.Fatura.DataEmissao, in.Fatura.DataVencimento, emptyDateAsBlank(in.Fatura.DataPagamento),
		in.Fatura.DataInicioPeriodo, in.Fatura.DataFimPeriodo,
		asString(completeness["status"]), missingArr,
		billing, in.DocumentRecord, now)
	if err != nil {
		return PersistSyncResult{}, err
	}

	// Insert raw_invoice_items
	if err := a.replaceInvoiceItems(ctx, tx, invoiceID, in); err != nil {
		return PersistSyncResult{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return PersistSyncResult{}, err
	}
	return PersistSyncResult{SyncRunID: syncRunID, InvoiceID: invoiceID}, nil
}

func (a *IntegrationPostgresAdapter) replaceInvoiceItems(ctx context.Context, tx pgx.Tx, invoiceID string, in PersistSyncInput) error {
	_, err := tx.Exec(ctx, `DELETE FROM integration.raw_invoice_items WHERE raw_invoice_id = $1`, invoiceID)
	if err != nil {
		return err
	}
	for i, item := range extractItems(in.BillingRecord) {
		descr := asString(item["descricao"])
		if descr == "" {
			descr = asString(item["description"])
		}
		var qtd, pu, val float64
		fmt.Sscan(asString(item["quantidade"]), &qtd)
		fmt.Sscan(asString(item["tarifa"]), &pu)
		fmt.Sscan(asString(item["valor_total"]), &val)
		rawJSON, _ := json.Marshal(item)
		_, err = tx.Exec(ctx, `
			INSERT INTO integration.raw_invoice_items (id, raw_invoice_id, type, description, quantidade, preco_unitario, valor_total, order_index, raw_json, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, uuid.New().String(), invoiceID, inferItemType(descr), descr, qtd, pu, val, i, rawJSON, time.Now())
		if err != nil {
			return err
		}
	}
	return nil
}

// --- Helpers ---

func toConsumerUnitView(u integration.ConsumerUnit) ConsumerUnitView {
	credID := ""
	if u.CredentialID != nil {
		credID = u.CredentialID.String()
	}
	status := ""
	if u.Status != nil {
		status = *u.Status
	}
	nome := ""
	if u.NomeCliente != nil {
		nome = *u.NomeCliente
	}
	inst := ""
	if u.Instalacao != nil {
		inst = *u.Instalacao
	}
	contrato := ""
	if u.Contrato != nil {
		contrato = *u.Contrato
	}
	grupo := ""
	if u.GrupoTensao != nil {
		grupo = *u.GrupoTensao
	}
	return ConsumerUnitView{
		UC:           u.UC,
		CredentialID: credID,
		Status:       status,
		NomeCliente:  nome,
		Instalacao:   inst,
		Contrato:     contrato,
		GrupoTensao:  grupo,
		Endereco:     u.Endereco,
		Imovel:       u.Imovel,
		CreatedAt:    u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    u.UpdatedAt.Format(time.RFC3339),
	}
}

func toInvoiceView(inv integration.RawInvoice) *InvoiceView {
	status := ""
	if inv.StatusFatura != nil {
		status = *inv.StatusFatura
	}
	valor := ""
	if inv.ValorTotal != nil {
		valor = *inv.ValorTotal
	}
	cod := ""
	if inv.CodigoBarras != nil {
		cod = *inv.CodigoBarras
	}
	emissao := ""
	if inv.DataEmissao != nil {
		emissao = *inv.DataEmissao
	}
	venc := ""
	if inv.DataVencimento != nil {
		venc = *inv.DataVencimento
	}
	pag := ""
	if inv.DataPagamento != nil {
		pag = *inv.DataPagamento
	}
	ini := ""
	if inv.DataInicioPeriodo != nil {
		ini = *inv.DataInicioPeriodo
	}
	fim := ""
	if inv.DataFimPeriodo != nil {
		fim = *inv.DataFimPeriodo
	}
	compStatus := ""
	if inv.CompletenessStatus != nil {
		compStatus = *inv.CompletenessStatus
	}
	var missing []any
	for _, m := range inv.CompletenessMissing {
		missing = append(missing, m)
	}
	return &InvoiceView{
		ID:                  inv.ID.String(),
		UC:                  inv.UC,
		NumeroFatura:        inv.NumeroFatura,
		MesReferencia:       inv.MesReferencia,
		StatusFatura:        status,
		ValorTotal:          valor,
		CodigoBarras:        cod,
		DataEmissao:         emissao,
		DataVencimento:      venc,
		DataPagamento:       pag,
		DataInicioPeriodo:   ini,
		DataFimPeriodo:      fim,
		CompletenessStatus:  compStatus,
		CompletenessMissing: missing,
		BillingRecord:       inv.BillingRecordJSON,
		DocumentRecord:      inv.DocumentRecordJSON,
		CreatedAt:           inv.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           inv.UpdatedAt.Format(time.RFC3339),
	}
}

func toSyncRunView(sr integration.SyncRun) *SyncRunView {
	credID := ""
	if sr.CredentialID != nil {
		credID = sr.CredentialID.String()
	}
	started := ""
	if sr.StartedAt != nil {
		started = sr.StartedAt.Format(time.RFC3339)
	}
	finished := ""
	if sr.FinishedAt != nil {
		finished = sr.FinishedAt.Format(time.RFC3339)
	}
	errMsg := ""
	if sr.ErrorMessage != nil {
		errMsg = *sr.ErrorMessage
	}
	return &SyncRunView{
		ID:           sr.ID.String(),
		CredentialID: credID,
		Documento:    sr.Documento,
		UC:           sr.UC,
		Status:       sr.Status,
		StartedAt:    started,
		FinishedAt:   finished,
		ErrorMessage: errMsg,
		RawResponse:  sr.RawResponseJSON,
	}
}

func inferItemType(desc string) string {
	d := lowerNoSpace(desc)
	switch {
	case contains(d, "tusd") && contains(d, "fio"):
		return "tusd_fio"
	case contains(d, "tusd") && contains(d, "energia"):
		return "tusd_energia"
	case contains(d, "injetada") || contains(d, "scee"):
		return "energia_injetada"
	case contains(d, "bandeira"):
		return "bandeira"
	case contains(d, "ipcoelba") || contains(d, "ip") && contains(d, "coelba"):
		return "ip_coelba"
	case contains(d, "reativo"):
		return "reativo_excedente"
	case contains(d, "tributo"):
		return "tributo_retido"
	default:
		return "tusd_energia"
	}
}

func lowerNoSpace(s string) string {
	var out []rune
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			out = append(out, r+('a'-'A'))
		} else if r != ' ' && r != '\t' && r != '\n' {
			out = append(out, r)
		}
	}
	return string(out)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) > 0 && findSubstr(s, substr))
}

func findSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
