package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// --- Credentials ---

func (s *pgxStore) InsertCredential(ctx context.Context, c *Credential) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	query := `
		INSERT INTO integration.credentials (id, label, documento_cipher, documento_nonce, senha_cipher, senha_nonce, uf, tipo_acesso, key_version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := s.pool.Exec(ctx, query,
		c.ID, c.Label, c.DocumentoCipher, c.DocumentoNonce, c.SenhaCipher, c.SenhaNonce,
		c.UF, c.TipoAcesso, c.KeyVersion, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (s *pgxStore) GetCredentialByID(ctx context.Context, id string) (*Credential, error) {
	query := `
		SELECT id, label, documento_cipher, documento_nonce, senha_cipher, senha_nonce, uf, tipo_acesso, key_version, created_at, updated_at
		FROM integration.credentials WHERE id = $1
	`
	row := s.pool.QueryRow(ctx, query, id)
	var c Credential
	var idStr string
	err := row.Scan(&idStr, &c.Label, &c.DocumentoCipher, &c.DocumentoNonce, &c.SenhaCipher, &c.SenhaNonce,
		&c.UF, &c.TipoAcesso, &c.KeyVersion, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	c.ID, _ = uuid.Parse(idStr)
	return &c, nil
}

// --- Sessions ---

func (s *pgxStore) InsertSession(ctx context.Context, sess *Session) error {
	if sess.ID == uuid.Nil {
		sess.ID = uuid.New()
	}
	now := time.Now()
	sess.CreatedAt = now
	sess.UpdatedAt = now

	query := `
		INSERT INTO integration.sessions (id, credential_id, bearer_token_cipher, bearer_token_nonce, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := s.pool.Exec(ctx, query,
		sess.ID, sess.CredentialID, sess.BearerTokenCipher, sess.BearerTokenNonce, sess.CreatedAt, sess.UpdatedAt,
	)
	return err
}

func (s *pgxStore) GetLatestSessionByCredentialID(ctx context.Context, credentialID string) (*Session, error) {
	query := `
		SELECT id, credential_id, bearer_token_cipher, bearer_token_nonce, created_at, updated_at
		FROM integration.sessions WHERE credential_id = $1 ORDER BY created_at DESC LIMIT 1
	`
	row := s.pool.QueryRow(ctx, query, credentialID)
	var sess Session
	var idStr string
	err := row.Scan(&idStr, &sess.CredentialID, &sess.BearerTokenCipher, &sess.BearerTokenNonce, &sess.CreatedAt, &sess.UpdatedAt)
	if err != nil {
		return nil, err
	}
	sess.ID, _ = uuid.Parse(idStr)
	return &sess, nil
}

// --- Consumer Units ---

func (s *pgxStore) UpsertConsumerUnit(ctx context.Context, u *ConsumerUnit) error {
	enderecoJSON, _ := json.Marshal(u.Endereco)
	imovelJSON, _ := json.Marshal(u.Imovel)

	query := `
		INSERT INTO integration.consumer_units (uc, credential_id, status, nome_cliente, instalacao, contrato, grupo_tensao, endereco_json, imovel_json, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		ON CONFLICT (uc) DO UPDATE SET
			credential_id = EXCLUDED.credential_id,
			status = EXCLUDED.status,
			nome_cliente = EXCLUDED.nome_cliente,
			instalacao = EXCLUDED.instalacao,
			contrato = EXCLUDED.contrato,
			grupo_tensao = EXCLUDED.grupo_tensao,
			endereco_json = EXCLUDED.endereco_json,
			imovel_json = EXCLUDED.imovel_json,
			updated_at = NOW()
	`
	_, err := s.pool.Exec(ctx, query,
		u.UC, u.CredentialID, u.Status, u.NomeCliente, u.Instalacao,
		u.Contrato, u.GrupoTensao, enderecoJSON, imovelJSON,
	)
	return err
}

func (s *pgxStore) ListConsumerUnits(ctx context.Context, limit int, status string) ([]ConsumerUnit, error) {
	if limit <= 0 {
		limit = 100
	}
	query := `SELECT uc, credential_id, status, nome_cliente, instalacao, contrato, grupo_tensao, endereco_json, imovel_json, created_at, updated_at FROM integration.consumer_units`
	args := []interface{}{}
	argNum := 1
	if status != "" {
		query += fmt.Sprintf(" WHERE status = $%d", argNum)
		args = append(args, status)
		argNum++
	}
	query += fmt.Sprintf(" ORDER BY updated_at DESC LIMIT $%d", argNum)
	args = append(args, limit)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []ConsumerUnit
	for rows.Next() {
		u, err := scanConsumerUnit(rows)
		if err != nil {
			return nil, err
		}
		units = append(units, *u)
	}
	return units, rows.Err()
}

func (s *pgxStore) GetConsumerUnitByUC(ctx context.Context, uc string) (*ConsumerUnit, error) {
	query := `SELECT uc, credential_id, status, nome_cliente, instalacao, contrato, grupo_tensao, endereco_json, imovel_json, created_at, updated_at FROM integration.consumer_units WHERE uc = $1`
	row := s.pool.QueryRow(ctx, query, uc)
	return scanConsumerUnit(row)
}

// --- Raw Invoices ---

func (s *pgxStore) UpsertRawInvoice(ctx context.Context, inv *RawInvoice) (*RawInvoice, error) {
	billingJSON, _ := json.Marshal(inv.BillingRecordJSON)
	documentJSON, _ := json.Marshal(inv.DocumentRecordJSON)

	var existingID uuid.UUID
	err := s.pool.QueryRow(ctx, `SELECT id FROM integration.raw_invoices WHERE uc = $1 AND numero_fatura = $2`, inv.UC, inv.NumeroFatura).Scan(&existingID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	if existingID != uuid.Nil {
		inv.ID = existingID
		query := `
			UPDATE integration.raw_invoices SET
				mes_referencia = $1, status_fatura = $2, valor_total = $3, codigo_barras = $4,
				data_emissao = $5, data_vencimento = $6, data_pagamento = $7,
				data_inicio_periodo = $8, data_fim_periodo = $9, completeness_status = $10,
				completeness_missing = $11, billing_record_json = $12, document_record_json = $13,
				updated_at = NOW()
			WHERE id = $14
		`
		_, err := s.pool.Exec(ctx, query,
			inv.MesReferencia, inv.StatusFatura, inv.ValorTotal, inv.CodigoBarras,
			inv.DataEmissao, inv.DataVencimento, inv.DataPagamento,
			inv.DataInicioPeriodo, inv.DataFimPeriodo, inv.CompletenessStatus,
			inv.CompletenessMissing, billingJSON, documentJSON, inv.ID,
		)
		return inv, err
	}

	if inv.ID == uuid.Nil {
		inv.ID = uuid.New()
	}
	query := `
		INSERT INTO integration.raw_invoices (id, uc, numero_fatura, mes_referencia, status_fatura, valor_total, codigo_barras, data_emissao, data_vencimento, data_pagamento, data_inicio_periodo, data_fim_periodo, completeness_status, completeness_missing, billing_record_json, document_record_json, pdf_bytes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, NOW(), NOW())
	`
	_, err = s.pool.Exec(ctx, query,
		inv.ID, inv.UC, inv.NumeroFatura, inv.MesReferencia, inv.StatusFatura, inv.ValorTotal,
		inv.CodigoBarras, inv.DataEmissao, inv.DataVencimento, inv.DataPagamento,
		inv.DataInicioPeriodo, inv.DataFimPeriodo, inv.CompletenessStatus,
		inv.CompletenessMissing, billingJSON, documentJSON, inv.PDFBytes,
	)
	return inv, err
}

func (s *pgxStore) GetRawInvoiceByID(ctx context.Context, id uuid.UUID) (*RawInvoice, error) {
	query := `SELECT id, uc, numero_fatura, mes_referencia, status_fatura, valor_total, codigo_barras, data_emissao, data_vencimento, data_pagamento, data_inicio_periodo, data_fim_periodo, completeness_status, completeness_missing, billing_record_json, document_record_json, pdf_bytes, created_at, updated_at FROM integration.raw_invoices WHERE id = $1`
	row := s.pool.QueryRow(ctx, query, id)
	return scanRawInvoice(row)
}

func (s *pgxStore) GetLatestRawInvoiceByUC(ctx context.Context, uc string) (*RawInvoice, error) {
	query := `SELECT id, uc, numero_fatura, mes_referencia, status_fatura, valor_total, codigo_barras, data_emissao, data_vencimento, data_pagamento, data_inicio_periodo, data_fim_periodo, completeness_status, completeness_missing, billing_record_json, document_record_json, pdf_bytes, created_at, updated_at FROM integration.raw_invoices WHERE uc = $1 ORDER BY updated_at DESC LIMIT 1`
	row := s.pool.QueryRow(ctx, query, uc)
	return scanRawInvoice(row)
}

func (s *pgxStore) ListRawInvoicesByUC(ctx context.Context, uc string, limit int) ([]RawInvoice, error) {
	if limit <= 0 {
		limit = 100
	}
	query := `SELECT id, uc, numero_fatura, mes_referencia, status_fatura, valor_total, codigo_barras, data_emissao, data_vencimento, data_pagamento, data_inicio_periodo, data_fim_periodo, completeness_status, completeness_missing, billing_record_json, document_record_json, pdf_bytes, created_at, updated_at FROM integration.raw_invoices WHERE uc = $1 ORDER BY updated_at DESC LIMIT $2`
	rows, err := s.pool.Query(ctx, query, uc, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []RawInvoice
	for rows.Next() {
		inv, err := scanRawInvoice(rows)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, *inv)
	}
	return invoices, rows.Err()
}

// --- Sync Runs ---

func (s *pgxStore) InsertSyncRun(ctx context.Context, sr *SyncRun) error {
	if sr.ID == uuid.Nil {
		sr.ID = uuid.New()
	}
	rawJSON, _ := json.Marshal(sr.RawResponseJSON)
	errCtx, _ := json.Marshal(sr.ErrorContext)

	query := `
		INSERT INTO integration.sync_runs (id, credential_id, documento, uc, status, step, error_message, error_context, raw_response_json, started_at, finished_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW())
	`
	_, err := s.pool.Exec(ctx, query,
		sr.ID, sr.CredentialID, sr.Documento, sr.UC, sr.Status, sr.Step,
		sr.ErrorMessage, errCtx, rawJSON, sr.StartedAt, sr.FinishedAt,
	)
	return err
}

func (s *pgxStore) GetSyncRunByID(ctx context.Context, id string) (*SyncRun, error) {
	query := `SELECT id, credential_id, documento, uc, status, step, error_message, error_context, raw_response_json, started_at, finished_at, created_at FROM integration.sync_runs WHERE id = $1`
	row := s.pool.QueryRow(ctx, query, id)
	return scanSyncRun(row)
}

func (s *pgxStore) GetLatestSyncRunByUC(ctx context.Context, uc string) (*SyncRun, error) {
	query := `SELECT id, credential_id, documento, uc, status, step, error_message, error_context, raw_response_json, started_at, finished_at, created_at FROM integration.sync_runs WHERE uc = $1 ORDER BY created_at DESC LIMIT 1`
	row := s.pool.QueryRow(ctx, query, uc)
	return scanSyncRun(row)
}

// --- Jobs (Worker Pool) ---

func (s *pgxStore) EnqueueJob(ctx context.Context, jobType string, payload map[string]any) (*Job, error) {
	payloadJSON, _ := json.Marshal(payload)
	j := &Job{
		ID:        uuid.New(),
		JobType:   jobType,
		Status:    "pending",
		Payload:   payload,
		RetryCount: 0,
		MaxRetries: 3,
		CreatedAt: time.Now(),
	}

	query := `
		INSERT INTO integration.jobs (id, job_type, status, payload, retry_count, max_retries, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := s.pool.Exec(ctx, query, j.ID, j.JobType, j.Status, payloadJSON, j.RetryCount, j.MaxRetries, j.CreatedAt)
	return j, err
}

func (s *pgxStore) ClaimNextJob(ctx context.Context, workerID string) (*Job, error) {
	// FOR UPDATE SKIP LOCKED: worker pool seguro sem concorrência
	query := `
		UPDATE integration.jobs
		SET status = 'running', claimed_by = $1, claimed_at = NOW()
		WHERE id = (
			SELECT id FROM integration.jobs
			WHERE status = 'pending'
			ORDER BY created_at ASC
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		)
		RETURNING id, job_type, status, payload, result, error_message, retry_count, max_retries, claimed_by, created_at, claimed_at, completed_at
	`
	row := s.pool.QueryRow(ctx, query, workerID)
	return scanJob(row)
}

func (s *pgxStore) CompleteJob(ctx context.Context, jobID uuid.UUID, result map[string]any) error {
	resultJSON, _ := json.Marshal(result)
	query := `
		UPDATE integration.jobs
		SET status = 'completed', result = $1, completed_at = NOW()
		WHERE id = $2
	`
	_, err := s.pool.Exec(ctx, query, resultJSON, jobID)
	return err
}

func (s *pgxStore) FailJob(ctx context.Context, jobID uuid.UUID, errMsg string) error {
	query := `
		UPDATE integration.jobs
		SET status = 'failed', error_message = $1, retry_count = retry_count + 1
		WHERE id = $2
	`
	_, err := s.pool.Exec(ctx, query, errMsg, jobID)
	return err
}

// --- Row scanners ---

func scanConsumerUnit(row pgx.Row) (*ConsumerUnit, error) {
	var u ConsumerUnit
	var enderecoJSON, imovelJSON []byte
	err := row.Scan(&u.UC, &u.CredentialID, &u.Status, &u.NomeCliente, &u.Instalacao,
		&u.Contrato, &u.GrupoTensao, &enderecoJSON, &imovelJSON, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if len(enderecoJSON) > 0 {
		json.Unmarshal(enderecoJSON, &u.Endereco)
	}
	if len(imovelJSON) > 0 {
		json.Unmarshal(imovelJSON, &u.Imovel)
	}
	return &u, nil
}

func scanRawInvoice(row pgx.Row) (*RawInvoice, error) {
	var inv RawInvoice
	var idStr string
	var billingJSON, documentJSON []byte
	var missingArr []string
	var statusFatura, valorTotal, codigoBarras, completenessStatus pgtype.Text
	var dataEmissao, dataVencimento, dataPagamento, dataInicio, dataFim pgtype.Date

	err := row.Scan(&idStr, &inv.UC, &inv.NumeroFatura, &inv.MesReferencia, &statusFatura,
		&valorTotal, &codigoBarras, &dataEmissao, &dataVencimento,
		&dataPagamento, &dataInicio, &dataFim,
		&completenessStatus, &missingArr, &billingJSON, &documentJSON,
		&inv.PDFBytes, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		return nil, err
	}
	inv.ID, _ = uuid.Parse(idStr)
	inv.StatusFatura = pgTextPtr(statusFatura)
	inv.ValorTotal = pgTextPtr(valorTotal)
	inv.CodigoBarras = pgTextPtr(codigoBarras)
	inv.DataEmissao = pgDatePtr(dataEmissao)
	inv.DataVencimento = pgDatePtr(dataVencimento)
	inv.DataPagamento = pgDatePtr(dataPagamento)
	inv.DataInicioPeriodo = pgDatePtr(dataInicio)
	inv.DataFimPeriodo = pgDatePtr(dataFim)
	inv.CompletenessStatus = pgTextPtr(completenessStatus)
	inv.CompletenessMissing = missingArr
	if len(billingJSON) > 0 {
		json.Unmarshal(billingJSON, &inv.BillingRecordJSON)
	}
	if len(documentJSON) > 0 {
		json.Unmarshal(documentJSON, &inv.DocumentRecordJSON)
	}
	return &inv, nil
}

func pgTextPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	s := t.String
	return &s
}

func pgDatePtr(d pgtype.Date) *string {
	if !d.Valid {
		return nil
	}
	s := d.Time.Format("2006-01-02")
	return &s
}

func scanSyncRun(row pgx.Row) (*SyncRun, error) {
	var sr SyncRun
	var idStr string
	var rawJSON, errCtx []byte

	err := row.Scan(&idStr, &sr.CredentialID, &sr.Documento, &sr.UC, &sr.Status, &sr.Step,
		&sr.ErrorMessage, &errCtx, &rawJSON, &sr.StartedAt, &sr.FinishedAt, &sr.CreatedAt)
	if err != nil {
		return nil, err
	}
	sr.ID, _ = uuid.Parse(idStr)
	if len(rawJSON) > 0 {
		json.Unmarshal(rawJSON, &sr.RawResponseJSON)
	}
	if len(errCtx) > 0 {
		json.Unmarshal(errCtx, &sr.ErrorContext)
	}
	return &sr, nil
}

func scanJob(row pgx.Row) (*Job, error) {
	var j Job
	var idStr string
	var payloadJSON, resultJSON []byte

	err := row.Scan(&idStr, &j.JobType, &j.Status, &payloadJSON, &resultJSON,
		&j.ErrorMessage, &j.RetryCount, &j.MaxRetries, &j.ClaimedBy,
		&j.CreatedAt, &j.ClaimedAt, &j.CompletedAt)
	if err != nil {
		return nil, err
	}
	j.ID, _ = uuid.Parse(idStr)
	if len(payloadJSON) > 0 {
		json.Unmarshal(payloadJSON, &j.Payload)
	}
	if len(resultJSON) > 0 {
		json.Unmarshal(resultJSON, &j.Result)
	}
	return &j, nil
}

// --- JSON helpers ---

func mustJSON(v any) []byte {
	if v == nil {
		return nil
	}
	b, _ := json.Marshal(v)
	return b
}

func nullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
