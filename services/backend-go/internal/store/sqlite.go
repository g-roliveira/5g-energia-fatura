package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/extractor"
	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/neoenergia"

	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

type CredentialRecord struct {
	ID              string
	Label           string
	DocumentoCipher string
	DocumentoNonce  string
	SenhaCipher     string
	SenhaNonce      string
	UF              string
	TipoAcesso      string
	KeyVersion      string
	CreatedAt       string
	UpdatedAt       string
}

type SessionRecord struct {
	ID                string
	CredentialID      string
	BearerTokenCipher string
	BearerTokenNonce  string
	CreatedAt         string
	UpdatedAt         string
}

type SyncRunRecord struct {
	ID           string
	CredentialID string
	Documento    string
	UC           string
	Status       string
	StartedAt    string
	FinishedAt   string
	ErrorMessage string
}

type PersistSyncInput struct {
	SyncRunID      string
	CredentialID   string
	Documento      string
	UC             string
	Status         string
	ErrorMessage   string
	UCRecord       *neoenergia.UC
	Imovel         *neoenergia.ImovelResponse
	Fatura         *neoenergia.Fatura
	Historico      *neoenergia.HistoricoConsumoResponse
	DadosPagamento *neoenergia.DadosPagamentoResponse
	PDF            *neoenergia.FaturaPDFResponse
	Extraction     *extractor.Response
	BillingRecord  any
	DocumentRecord any
	RawResponse    any
}

type PersistSyncResult struct {
	SyncRunID string
	InvoiceID string
}

type ConsumerUnitView struct {
	UC           string         `json:"uc"`
	CredentialID string         `json:"credential_id,omitempty"`
	Status       string         `json:"status,omitempty"`
	NomeCliente  string         `json:"nome_cliente,omitempty"`
	Instalacao   string         `json:"instalacao,omitempty"`
	Contrato     string         `json:"contrato,omitempty"`
	GrupoTensao  string         `json:"grupo_tensao,omitempty"`
	Endereco     map[string]any `json:"endereco,omitempty"`
	Imovel       map[string]any `json:"imovel,omitempty"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    string         `json:"updated_at"`
}

type ConsumerUnitDetailsView struct {
	ConsumerUnitView
	LatestInvoice *InvoiceView `json:"latest_invoice,omitempty"`
	LatestSyncRun *SyncRunView `json:"latest_sync_run,omitempty"`
}

type InvoiceView struct {
	ID                  string           `json:"id"`
	UC                  string           `json:"uc"`
	NumeroFatura        string           `json:"numero_fatura"`
	MesReferencia       string           `json:"mes_referencia"`
	StatusFatura        string           `json:"status_fatura,omitempty"`
	ValorTotal          string           `json:"valor_total,omitempty"`
	CodigoBarras        string           `json:"codigo_barras,omitempty"`
	DataEmissao         string           `json:"data_emissao,omitempty"`
	DataVencimento      string           `json:"data_vencimento,omitempty"`
	DataPagamento       string           `json:"data_pagamento,omitempty"`
	DataInicioPeriodo   string           `json:"data_inicio_periodo,omitempty"`
	DataFimPeriodo      string           `json:"data_fim_periodo,omitempty"`
	CompletenessStatus  string           `json:"completeness_status,omitempty"`
	CompletenessMissing []any            `json:"completeness_missing,omitempty"`
	BillingRecord       map[string]any   `json:"billing_record,omitempty"`
	DocumentRecord      map[string]any   `json:"document_record,omitempty"`
	Items               []map[string]any `json:"items,omitempty"`
	CreatedAt           string           `json:"created_at"`
	UpdatedAt           string           `json:"updated_at"`
}

type SyncRunView struct {
	ID           string         `json:"id"`
	CredentialID string         `json:"credential_id,omitempty"`
	Documento    string         `json:"documento"`
	UC           string         `json:"uc"`
	Status       string         `json:"status"`
	StartedAt    string         `json:"started_at"`
	FinishedAt   string         `json:"finished_at,omitempty"`
	ErrorMessage string         `json:"error_message,omitempty"`
	RawResponse  map[string]any `json:"raw_response,omitempty"`
}

func OpenSQLite(dsn string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	store := &SQLiteStore{db: db}
	if err := store.migrate(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) migrate() error {
	stmts := []string{
		`create table if not exists credentials (
			id text primary key,
			label text not null,
			documento_cipher text not null,
			documento_nonce text not null,
			senha_cipher text not null,
			senha_nonce text not null,
			uf text not null,
			tipo_acesso text not null,
			key_version text not null,
			created_at text not null,
			updated_at text not null
		)`,
		`create table if not exists sessions (
			id text primary key,
			credential_id text not null,
			bearer_token_cipher text not null,
			bearer_token_nonce text not null,
			created_at text not null,
			updated_at text not null,
			foreign key (credential_id) references credentials(id)
		)`,
		`create index if not exists idx_sessions_credential_created_at on sessions(credential_id, created_at desc)`,
		`create table if not exists consumer_units (
			uc text primary key,
			credential_id text,
			status text,
			nome_cliente text,
			instalacao text,
			contrato text,
			grupo_tensao text,
			endereco_json text,
			imovel_json text,
			created_at text not null,
			updated_at text not null,
			foreign key (credential_id) references credentials(id)
		)`,
		`create table if not exists sync_runs (
			id text primary key,
			credential_id text,
			documento text not null,
			uc text not null,
			status text not null,
			started_at text not null,
			finished_at text,
			error_message text,
			raw_response_json text,
			foreign key (credential_id) references credentials(id)
		)`,
		`create index if not exists idx_sync_runs_uc_started_at on sync_runs(uc, started_at desc)`,
		`create table if not exists invoices (
			id text primary key,
			uc text not null,
			numero_fatura text not null,
			mes_referencia text not null,
			status_fatura text,
			valor_total text,
			codigo_barras text,
			data_emissao text,
			data_vencimento text,
			data_pagamento text,
			data_inicio_periodo text,
			data_fim_periodo text,
			completeness_status text,
			completeness_missing_json text,
			billing_record_json text,
			document_record_json text,
			created_at text not null,
			updated_at text not null,
			unique (uc, numero_fatura)
		)`,
		`create index if not exists idx_invoices_uc_mes on invoices(uc, mes_referencia)`,
		`create table if not exists invoice_api_snapshots (
			id text primary key,
			invoice_id text not null,
			sync_run_id text not null,
			fatura_json text,
			historico_json text,
			dados_pagamento_json text,
			created_at text not null,
			foreign key (invoice_id) references invoices(id),
			foreign key (sync_run_id) references sync_runs(id)
		)`,
		`create table if not exists invoice_documents (
			id text primary key,
			invoice_id text not null,
			sync_run_id text not null,
			file_name text,
			file_extension text,
			file_size text,
			file_data_base64 text,
			storage_uri text,
			created_at text not null,
			foreign key (invoice_id) references invoices(id),
			foreign key (sync_run_id) references sync_runs(id)
		)`,
		`create table if not exists invoice_items (
			id text primary key,
			invoice_id text not null,
			descricao text,
			quantidade text,
			quantidade_residual text,
			quantidade_faturada text,
			tarifa text,
			valor text,
			base_icms text,
			aliq_icms text,
			icms text,
			valor_total text,
			raw_json text,
			created_at text not null,
			foreign key (invoice_id) references invoices(id)
		)`,
		`create table if not exists invoice_extraction_results (
			id text primary key,
			invoice_id text not null,
			sync_run_id text not null,
			status text,
			fields_json text,
			source_map_json text,
			confidence_map_json text,
			warnings_json text,
			artifacts_json text,
			created_at text not null,
			foreign key (invoice_id) references invoices(id),
			foreign key (sync_run_id) references sync_runs(id)
		)`,
		`create table if not exists invoice_field_sources (
			id text primary key,
			invoice_id text not null,
			field_path text not null,
			source text,
			confidence real,
			created_at text not null,
			foreign key (invoice_id) references invoices(id)
		)`,
		`create index if not exists idx_invoice_field_sources_invoice on invoice_field_sources(invoice_id)`,
	}
	for _, stmt := range stmts {
		if _, err := s.db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}

func nowISO() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func newID() string {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}

func (s *SQLiteStore) InsertCredential(rec CredentialRecord) (CredentialRecord, error) {
	now := nowISO()
	if rec.ID == "" {
		rec.ID = newID()
	}
	rec.CreatedAt = now
	rec.UpdatedAt = now
	_, err := s.db.Exec(
		`insert into credentials (
			id, label, documento_cipher, documento_nonce, senha_cipher, senha_nonce,
			uf, tipo_acesso, key_version, created_at, updated_at
		) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		rec.ID, rec.Label, rec.DocumentoCipher, rec.DocumentoNonce, rec.SenhaCipher, rec.SenhaNonce,
		rec.UF, rec.TipoAcesso, rec.KeyVersion, rec.CreatedAt, rec.UpdatedAt,
	)
	return rec, err
}

func (s *SQLiteStore) GetCredentialByID(id string) (CredentialRecord, error) {
	var rec CredentialRecord
	err := s.db.QueryRow(
		`select id, label, documento_cipher, documento_nonce, senha_cipher, senha_nonce,
		        uf, tipo_acesso, key_version, created_at, updated_at
		   from credentials where id = ?`,
		id,
	).Scan(
		&rec.ID, &rec.Label, &rec.DocumentoCipher, &rec.DocumentoNonce, &rec.SenhaCipher, &rec.SenhaNonce,
		&rec.UF, &rec.TipoAcesso, &rec.KeyVersion, &rec.CreatedAt, &rec.UpdatedAt,
	)
	return rec, err
}

func (s *SQLiteStore) InsertSession(rec SessionRecord) (SessionRecord, error) {
	now := nowISO()
	if rec.ID == "" {
		rec.ID = newID()
	}
	rec.CreatedAt = now
	rec.UpdatedAt = now
	_, err := s.db.Exec(
		`insert into sessions (
			id, credential_id, bearer_token_cipher, bearer_token_nonce, created_at, updated_at
		) values (?, ?, ?, ?, ?, ?)`,
		rec.ID, rec.CredentialID, rec.BearerTokenCipher, rec.BearerTokenNonce, rec.CreatedAt, rec.UpdatedAt,
	)
	return rec, err
}

func (s *SQLiteStore) GetLatestSessionByCredentialID(credentialID string) (SessionRecord, error) {
	var rec SessionRecord
	err := s.db.QueryRow(
		`select id, credential_id, bearer_token_cipher, bearer_token_nonce, created_at, updated_at
		   from sessions where credential_id = ?
		   order by created_at desc
		   limit 1`,
		credentialID,
	).Scan(
		&rec.ID, &rec.CredentialID, &rec.BearerTokenCipher, &rec.BearerTokenNonce, &rec.CreatedAt, &rec.UpdatedAt,
	)
	return rec, err
}

func (s *SQLiteStore) CountRows(table string) (int, error) {
	var count int
	err := s.db.QueryRow(`select count(*) from ` + table).Scan(&count)
	return count, err
}

func (s *SQLiteStore) ListConsumerUnits(limit int, status string) ([]ConsumerUnitView, error) {
	if limit <= 0 {
		limit = 100
	}
	query := `select uc, coalesce(credential_id,''), coalesce(status,''), coalesce(nome_cliente,''), coalesce(instalacao,''),
		        coalesce(contrato,''), coalesce(grupo_tensao,''), coalesce(endereco_json,''), coalesce(imovel_json,''),
		        created_at, updated_at
		   from consumer_units`
	args := []any{}
	if status != "" {
		query += ` where status = ?`
		args = append(args, status)
	}
	query += ` order by updated_at desc limit ?`
	args = append(args, limit)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	views := []ConsumerUnitView{}
	for rows.Next() {
		var view ConsumerUnitView
		var enderecoJSON, imovelJSON string
		if err := rows.Scan(
			&view.UC, &view.CredentialID, &view.Status, &view.NomeCliente, &view.Instalacao,
			&view.Contrato, &view.GrupoTensao, &enderecoJSON, &imovelJSON, &view.CreatedAt, &view.UpdatedAt,
		); err != nil {
			return nil, err
		}
		view.Endereco = decodeJSONMap(enderecoJSON)
		view.Imovel = decodeJSONMap(imovelJSON)
		views = append(views, view)
	}
	return views, rows.Err()
}

func (s *SQLiteStore) GetConsumerUnitByUC(uc string) (*ConsumerUnitDetailsView, error) {
	row := s.db.QueryRow(
		`select uc, coalesce(credential_id,''), coalesce(status,''), coalesce(nome_cliente,''), coalesce(instalacao,''),
		        coalesce(contrato,''), coalesce(grupo_tensao,''), coalesce(endereco_json,''), coalesce(imovel_json,''),
		        created_at, updated_at
		   from consumer_units where uc = ?`,
		uc,
	)
	var view ConsumerUnitDetailsView
	var enderecoJSON, imovelJSON string
	if err := row.Scan(
		&view.UC, &view.CredentialID, &view.Status, &view.NomeCliente, &view.Instalacao,
		&view.Contrato, &view.GrupoTensao, &enderecoJSON, &imovelJSON, &view.CreatedAt, &view.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	view.Endereco = decodeJSONMap(enderecoJSON)
	view.Imovel = decodeJSONMap(imovelJSON)
	latestInvoice, err := s.GetLatestInvoiceByUC(uc)
	if err != nil {
		return nil, err
	}
	view.LatestInvoice = latestInvoice
	latestSyncRun, err := s.GetLatestSyncRunByUC(uc)
	if err != nil {
		return nil, err
	}
	view.LatestSyncRun = latestSyncRun
	return &view, nil
}

func (s *SQLiteStore) ListInvoicesByUC(uc string, limit int, status string) ([]InvoiceView, error) {
	if limit <= 0 {
		limit = 100
	}
	query := `select id, uc, numero_fatura, mes_referencia, coalesce(status_fatura,''), coalesce(valor_total,''),
		        coalesce(codigo_barras,''), coalesce(data_emissao,''), coalesce(data_vencimento,''),
		        coalesce(data_pagamento,''), coalesce(data_inicio_periodo,''), coalesce(data_fim_periodo,''),
		        coalesce(completeness_status,''), coalesce(completeness_missing_json,''), coalesce(billing_record_json,''),
		        coalesce(document_record_json,''), created_at, updated_at
		   from invoices where uc = ?`
	args := []any{uc}
	if status != "" {
		query += ` and status_fatura = ?`
		args = append(args, status)
	}
	query += ` order by updated_at desc limit ?`
	args = append(args, limit)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	views := []InvoiceView{}
	for rows.Next() {
		view, err := scanInvoice(rows)
		if err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	return views, rows.Err()
}

func (s *SQLiteStore) GetInvoiceByID(id string) (*InvoiceView, error) {
	row := s.db.QueryRow(
		`select id, uc, numero_fatura, mes_referencia, coalesce(status_fatura,''), coalesce(valor_total,''),
		        coalesce(codigo_barras,''), coalesce(data_emissao,''), coalesce(data_vencimento,''),
		        coalesce(data_pagamento,''), coalesce(data_inicio_periodo,''), coalesce(data_fim_periodo,''),
		        coalesce(completeness_status,''), coalesce(completeness_missing_json,''), coalesce(billing_record_json,''),
		        coalesce(document_record_json,''), created_at, updated_at
		   from invoices where id = ?`,
		id,
	)
	view, err := scanInvoice(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	items, err := s.listInvoiceItems(id)
	if err != nil {
		return nil, err
	}
	view.Items = items
	return &view, nil
}

func (s *SQLiteStore) GetLatestInvoiceByUC(uc string) (*InvoiceView, error) {
	row := s.db.QueryRow(
		`select id, uc, numero_fatura, mes_referencia, coalesce(status_fatura,''), coalesce(valor_total,''),
		        coalesce(codigo_barras,''), coalesce(data_emissao,''), coalesce(data_vencimento,''),
		        coalesce(data_pagamento,''), coalesce(data_inicio_periodo,''), coalesce(data_fim_periodo,''),
		        coalesce(completeness_status,''), coalesce(completeness_missing_json,''), coalesce(billing_record_json,''),
		        coalesce(document_record_json,''), created_at, updated_at
		   from invoices where uc = ?
		   order by updated_at desc limit 1`,
		uc,
	)
	view, err := scanInvoice(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	items, err := s.listInvoiceItems(view.ID)
	if err != nil {
		return nil, err
	}
	view.Items = items
	return &view, nil
}

func (s *SQLiteStore) GetSyncRunByID(id string) (*SyncRunView, error) {
	row := s.db.QueryRow(
		`select id, coalesce(credential_id,''), documento, uc, status, started_at, coalesce(finished_at,''), coalesce(error_message,''), coalesce(raw_response_json,'')
		   from sync_runs where id = ?`,
		id,
	)
	var view SyncRunView
	var rawJSON string
	if err := row.Scan(&view.ID, &view.CredentialID, &view.Documento, &view.UC, &view.Status, &view.StartedAt, &view.FinishedAt, &view.ErrorMessage, &rawJSON); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	view.RawResponse = decodeJSONMap(rawJSON)
	return &view, nil
}

func (s *SQLiteStore) GetLatestSyncRunByUC(uc string) (*SyncRunView, error) {
	row := s.db.QueryRow(
		`select id, coalesce(credential_id,''), documento, uc, status, started_at, coalesce(finished_at,''), coalesce(error_message,''), coalesce(raw_response_json,'')
		   from sync_runs where uc = ? order by started_at desc limit 1`,
		uc,
	)
	var view SyncRunView
	var rawJSON string
	if err := row.Scan(&view.ID, &view.CredentialID, &view.Documento, &view.UC, &view.Status, &view.StartedAt, &view.FinishedAt, &view.ErrorMessage, &rawJSON); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	view.RawResponse = decodeJSONMap(rawJSON)
	return &view, nil
}

func (s *SQLiteStore) listInvoiceItems(invoiceID string) ([]map[string]any, error) {
	rows, err := s.db.Query(
		`select coalesce(raw_json,'') from invoice_items where invoice_id = ? order by descricao asc`,
		invoiceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}
		items = append(items, decodeJSONMap(raw))
	}
	return items, rows.Err()
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanInvoice(row rowScanner) (InvoiceView, error) {
	var view InvoiceView
	var missingJSON, billingJSON, documentJSON string
	err := row.Scan(
		&view.ID, &view.UC, &view.NumeroFatura, &view.MesReferencia, &view.StatusFatura, &view.ValorTotal,
		&view.CodigoBarras, &view.DataEmissao, &view.DataVencimento, &view.DataPagamento,
		&view.DataInicioPeriodo, &view.DataFimPeriodo, &view.CompletenessStatus, &missingJSON,
		&billingJSON, &documentJSON, &view.CreatedAt, &view.UpdatedAt,
	)
	if err != nil {
		return InvoiceView{}, err
	}
	view.CompletenessMissing = decodeJSONArray(missingJSON)
	view.BillingRecord = decodeJSONMap(billingJSON)
	view.DocumentRecord = decodeJSONMap(documentJSON)
	return view, nil
}

func decodeJSONMap(raw string) map[string]any {
	if raw == "" {
		return nil
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return out
}

func decodeJSONArray(raw string) []any {
	if raw == "" {
		return nil
	}
	var out []any
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return out
}

func (s *SQLiteStore) PersistSyncResult(in PersistSyncInput) (PersistSyncResult, error) {
	now := nowISO()
	syncRunID := in.SyncRunID
	if syncRunID == "" {
		syncRunID = newID()
	}
	status := in.Status
	if status == "" {
		status = "succeeded"
	}

	tx, err := s.db.Begin()
	if err != nil {
		return PersistSyncResult{}, err
	}
	defer tx.Rollback()

	rawResponseJSON := mustJSON(in.RawResponse)
	_, err = tx.Exec(
		`insert into sync_runs (
			id, credential_id, documento, uc, status, started_at, finished_at, error_message, raw_response_json
		) values (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		syncRunID, nullableString(in.CredentialID), in.Documento, in.UC, status, now, now, nullableString(in.ErrorMessage), nullableString(rawResponseJSON),
	)
	if err != nil {
		return PersistSyncResult{}, err
	}

	if in.UCRecord != nil {
		_, err = tx.Exec(
			`insert into consumer_units (
				uc, credential_id, status, nome_cliente, instalacao, contrato, grupo_tensao,
				endereco_json, imovel_json, created_at, updated_at
			) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			on conflict(uc) do update set
				credential_id=excluded.credential_id,
				status=excluded.status,
				nome_cliente=excluded.nome_cliente,
				instalacao=excluded.instalacao,
				contrato=excluded.contrato,
				grupo_tensao=excluded.grupo_tensao,
				endereco_json=excluded.endereco_json,
				imovel_json=excluded.imovel_json,
				updated_at=excluded.updated_at`,
			in.UCRecord.UC,
			nullableString(in.CredentialID),
			in.UCRecord.Status,
			in.UCRecord.NomeCliente,
			in.UCRecord.Instalacao,
			in.UCRecord.Contrato,
			in.UCRecord.GrupoTensao,
			nullableString(mustJSON(in.UCRecord.Local)),
			nullableString(mustJSON(in.Imovel)),
			now,
			now,
		)
		if err != nil {
			return PersistSyncResult{}, err
		}
	}

	if in.Fatura == nil {
		if err := tx.Commit(); err != nil {
			return PersistSyncResult{}, err
		}
		return PersistSyncResult{SyncRunID: syncRunID}, nil
	}

	invoiceID, err := upsertInvoice(tx, now, in)
	if err != nil {
		return PersistSyncResult{}, err
	}
	if err := replaceInvoiceChildren(tx, now, invoiceID, syncRunID, in); err != nil {
		return PersistSyncResult{}, err
	}
	if err := tx.Commit(); err != nil {
		return PersistSyncResult{}, err
	}
	return PersistSyncResult{SyncRunID: syncRunID, InvoiceID: invoiceID}, nil
}

func upsertInvoice(tx *sql.Tx, now string, in PersistSyncInput) (string, error) {
	var existingID string
	err := tx.QueryRow(
		`select id from invoices where uc = ? and numero_fatura = ?`,
		in.Fatura.UC,
		in.Fatura.NumeroFatura,
	).Scan(&existingID)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	invoiceID := existingID
	if invoiceID == "" {
		invoiceID = newID()
	}

	billing := asMap(in.BillingRecord)
	completeness := asMap(billing["completeness"])
	missingJSON := mustJSON(completeness["missing_fields"])
	codigoBarras := ""
	if in.DadosPagamento != nil {
		codigoBarras = in.DadosPagamento.CodBarras
	}

	_, err = tx.Exec(
		`insert into invoices (
			id, uc, numero_fatura, mes_referencia, status_fatura, valor_total, codigo_barras,
			data_emissao, data_vencimento, data_pagamento, data_inicio_periodo, data_fim_periodo,
			completeness_status, completeness_missing_json, billing_record_json, document_record_json,
			created_at, updated_at
		) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		on conflict(uc, numero_fatura) do update set
			mes_referencia=excluded.mes_referencia,
			status_fatura=excluded.status_fatura,
			valor_total=excluded.valor_total,
			codigo_barras=excluded.codigo_barras,
			data_emissao=excluded.data_emissao,
			data_vencimento=excluded.data_vencimento,
			data_pagamento=excluded.data_pagamento,
			data_inicio_periodo=excluded.data_inicio_periodo,
			data_fim_periodo=excluded.data_fim_periodo,
			completeness_status=excluded.completeness_status,
			completeness_missing_json=excluded.completeness_missing_json,
			billing_record_json=excluded.billing_record_json,
			document_record_json=excluded.document_record_json,
			updated_at=excluded.updated_at`,
		invoiceID,
		firstNonEmpty(in.Fatura.UC, in.UC),
		in.Fatura.NumeroFatura,
		in.Fatura.MesReferencia,
		nullableString(in.Fatura.StatusFatura),
		nullableString(in.Fatura.ValorEmissao),
		nullableString(codigoBarras),
		nullableString(in.Fatura.DataEmissao),
		nullableString(in.Fatura.DataVencimento),
		nullableString(emptyDateAsBlank(in.Fatura.DataPagamento)),
		nullableString(in.Fatura.DataInicioPeriodo),
		nullableString(in.Fatura.DataFimPeriodo),
		nullableString(asString(completeness["status"])),
		nullableString(missingJSON),
		nullableString(mustJSON(in.BillingRecord)),
		nullableString(mustJSON(in.DocumentRecord)),
		now,
		now,
	)
	return invoiceID, err
}

func replaceInvoiceChildren(tx *sql.Tx, now string, invoiceID string, syncRunID string, in PersistSyncInput) error {
	for _, stmt := range []string{
		`delete from invoice_api_snapshots where invoice_id = ?`,
		`delete from invoice_documents where invoice_id = ?`,
		`delete from invoice_items where invoice_id = ?`,
		`delete from invoice_extraction_results where invoice_id = ?`,
		`delete from invoice_field_sources where invoice_id = ?`,
	} {
		if _, err := tx.Exec(stmt, invoiceID); err != nil {
			return err
		}
	}

	if _, err := tx.Exec(
		`insert into invoice_api_snapshots (
			id, invoice_id, sync_run_id, fatura_json, historico_json, dados_pagamento_json, created_at
		) values (?, ?, ?, ?, ?, ?, ?)`,
		newID(), invoiceID, syncRunID, nullableString(mustJSON(in.Fatura)), nullableString(mustJSON(in.Historico)), nullableString(mustJSON(in.DadosPagamento)), now,
	); err != nil {
		return err
	}

	if in.PDF != nil {
		if _, err := tx.Exec(
			`insert into invoice_documents (
				id, invoice_id, sync_run_id, file_name, file_extension, file_size, file_data_base64, storage_uri, created_at
			) values (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			newID(), invoiceID, syncRunID, nullableString(in.PDF.FileName), nullableString(in.PDF.FileExtension), nullableString(in.PDF.FileSize), nullableString(in.PDF.FileData), nil, now,
		); err != nil {
			return err
		}
	}

	for _, item := range extractItems(in.BillingRecord) {
		if _, err := tx.Exec(
			`insert into invoice_items (
				id, invoice_id, descricao, quantidade, quantidade_residual, quantidade_faturada,
				tarifa, valor, base_icms, aliq_icms, icms, valor_total, raw_json, created_at
			) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			newID(), invoiceID,
			nullableString(asString(item["descricao"])),
			nullableString(asString(item["quantidade"])),
			nullableString(asString(item["quantidade_residual"])),
			nullableString(asString(item["quantidade_faturada"])),
			nullableString(asString(item["tarifa"])),
			nullableString(asString(item["valor"])),
			nullableString(asString(item["base_icms"])),
			nullableString(asString(item["aliq_icms"])),
			nullableString(asString(item["icms"])),
			nullableString(asString(item["valor_total"])),
			nullableString(mustJSON(item)),
			now,
		); err != nil {
			return err
		}
	}

	if in.Extraction != nil {
		if _, err := tx.Exec(
			`insert into invoice_extraction_results (
				id, invoice_id, sync_run_id, status, fields_json, source_map_json,
				confidence_map_json, warnings_json, artifacts_json, created_at
			) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			newID(), invoiceID, syncRunID, nullableString(in.Extraction.Status), nullableString(mustJSON(in.Extraction.Fields)),
			nullableString(mustJSON(in.Extraction.SourceMap)), nullableString(mustJSON(in.Extraction.ConfidenceMap)),
			nullableString(mustJSON(in.Extraction.Warnings)), nullableString(mustJSON(in.Extraction.Artifacts)), now,
		); err != nil {
			return err
		}
	}

	billing := asMap(in.BillingRecord)
	sourceMap := asStringMap(billing["source_map"])
	confidenceMap := asFloatMap(billing["confidence_map"])
	for field, source := range sourceMap {
		if _, err := tx.Exec(
			`insert into invoice_field_sources (
				id, invoice_id, field_path, source, confidence, created_at
			) values (?, ?, ?, ?, ?, ?)`,
			newID(), invoiceID, field, source, confidenceMap[field], now,
		); err != nil {
			return err
		}
	}
	return nil
}

func mustJSON(value any) string {
	if value == nil {
		return ""
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(raw)
}

func nullableString(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func emptyDateAsBlank(value string) string {
	if value == "0000-00-00" {
		return ""
	}
	return value
}

func asMap(value any) map[string]any {
	if typed, ok := value.(map[string]any); ok {
		return typed
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return map[string]any{}
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		return map[string]any{}
	}
	return out
}

func asString(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case nil:
		return ""
	default:
		return ""
	}
}

func asStringMap(value any) map[string]string {
	rawMap := asMap(value)
	out := make(map[string]string, len(rawMap))
	for key, value := range rawMap {
		out[key] = asString(value)
	}
	return out
}

func asFloatMap(value any) map[string]float64 {
	rawMap := asMap(value)
	out := make(map[string]float64, len(rawMap))
	for key, value := range rawMap {
		switch typed := value.(type) {
		case float64:
			out[key] = typed
		case int:
			out[key] = float64(typed)
		}
	}
	return out
}

func extractItems(value any) []map[string]any {
	billing := asMap(value)
	rawItems, ok := billing["itens_fatura"].([]any)
	if !ok {
		return nil
	}
	items := make([]map[string]any, 0, len(rawItems))
	for _, raw := range rawItems {
		item := asMap(raw)
		if len(item) > 0 {
			items = append(items, item)
		}
	}
	return items
}
