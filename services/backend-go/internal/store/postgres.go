package store

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStore struct {
	db *sql.DB
}

func OpenPostgres(dsn string) (*PostgresStore, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	store := &PostgresStore{db: db}
	if err := store.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}

func (s *PostgresStore) migrate() error {
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
			confidence double precision,
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

func (s *PostgresStore) InsertCredential(rec CredentialRecord) (CredentialRecord, error) {
	now := nowISO()
	if rec.ID == "" {
		rec.ID = newID()
	}
	rec.CreatedAt = now
	rec.UpdatedAt = now
	_, err := s.db.Exec(
		pgQuery(`insert into credentials (
			id, label, documento_cipher, documento_nonce, senha_cipher, senha_nonce,
			uf, tipo_acesso, key_version, created_at, updated_at
		) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`),
		rec.ID, rec.Label, rec.DocumentoCipher, rec.DocumentoNonce, rec.SenhaCipher, rec.SenhaNonce,
		rec.UF, rec.TipoAcesso, rec.KeyVersion, rec.CreatedAt, rec.UpdatedAt,
	)
	return rec, err
}

func (s *PostgresStore) GetCredentialByID(id string) (CredentialRecord, error) {
	var rec CredentialRecord
	err := s.db.QueryRow(
		pgQuery(`select id, label, documento_cipher, documento_nonce, senha_cipher, senha_nonce,
		        uf, tipo_acesso, key_version, created_at, updated_at
		   from credentials where id = ?`),
		id,
	).Scan(
		&rec.ID, &rec.Label, &rec.DocumentoCipher, &rec.DocumentoNonce, &rec.SenhaCipher, &rec.SenhaNonce,
		&rec.UF, &rec.TipoAcesso, &rec.KeyVersion, &rec.CreatedAt, &rec.UpdatedAt,
	)
	return rec, err
}

func (s *PostgresStore) InsertSession(rec SessionRecord) (SessionRecord, error) {
	now := nowISO()
	if rec.ID == "" {
		rec.ID = newID()
	}
	rec.CreatedAt = now
	rec.UpdatedAt = now
	_, err := s.db.Exec(
		pgQuery(`insert into sessions (
			id, credential_id, bearer_token_cipher, bearer_token_nonce, created_at, updated_at
		) values (?, ?, ?, ?, ?, ?)`),
		rec.ID, rec.CredentialID, rec.BearerTokenCipher, rec.BearerTokenNonce, rec.CreatedAt, rec.UpdatedAt,
	)
	return rec, err
}

func (s *PostgresStore) GetLatestSessionByCredentialID(credentialID string) (SessionRecord, error) {
	var rec SessionRecord
	err := s.db.QueryRow(
		pgQuery(`select id, credential_id, bearer_token_cipher, bearer_token_nonce, created_at, updated_at
		   from sessions where credential_id = ?
		   order by created_at desc
		   limit 1`),
		credentialID,
	).Scan(
		&rec.ID, &rec.CredentialID, &rec.BearerTokenCipher, &rec.BearerTokenNonce, &rec.CreatedAt, &rec.UpdatedAt,
	)
	return rec, err
}

func (s *PostgresStore) ListConsumerUnits(limit int, status string) ([]ConsumerUnitView, error) {
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
	rows, err := s.db.Query(pgQuery(query), args...)
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

func (s *PostgresStore) GetConsumerUnitByUC(uc string) (*ConsumerUnitDetailsView, error) {
	row := s.db.QueryRow(
		pgQuery(`select uc, coalesce(credential_id,''), coalesce(status,''), coalesce(nome_cliente,''), coalesce(instalacao,''),
		        coalesce(contrato,''), coalesce(grupo_tensao,''), coalesce(endereco_json,''), coalesce(imovel_json,''),
		        created_at, updated_at
	   from consumer_units where uc = ?`),
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

func (s *PostgresStore) ListInvoicesByUC(uc string, limit int, status string) ([]InvoiceView, error) {
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
	rows, err := s.db.Query(pgQuery(query), args...)
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

func (s *PostgresStore) GetInvoiceByID(id string) (*InvoiceView, error) {
	row := s.db.QueryRow(
		pgQuery(`select id, uc, numero_fatura, mes_referencia, coalesce(status_fatura,''), coalesce(valor_total,''),
		        coalesce(codigo_barras,''), coalesce(data_emissao,''), coalesce(data_vencimento,''),
		        coalesce(data_pagamento,''), coalesce(data_inicio_periodo,''), coalesce(data_fim_periodo,''),
		        coalesce(completeness_status,''), coalesce(completeness_missing_json,''), coalesce(billing_record_json,''),
		        coalesce(document_record_json,''), created_at, updated_at
	   from invoices where id = ?`),
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

func (s *PostgresStore) GetLatestInvoiceByUC(uc string) (*InvoiceView, error) {
	row := s.db.QueryRow(
		pgQuery(`select id, uc, numero_fatura, mes_referencia, coalesce(status_fatura,''), coalesce(valor_total,''),
		        coalesce(codigo_barras,''), coalesce(data_emissao,''), coalesce(data_vencimento,''),
		        coalesce(data_pagamento,''), coalesce(data_inicio_periodo,''), coalesce(data_fim_periodo,''),
		        coalesce(completeness_status,''), coalesce(completeness_missing_json,''), coalesce(billing_record_json,''),
		        coalesce(document_record_json,''), created_at, updated_at
	   from invoices where uc = ?
	   order by updated_at desc limit 1`),
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

func (s *PostgresStore) GetSyncRunByID(id string) (*SyncRunView, error) {
	row := s.db.QueryRow(
		pgQuery(`select id, coalesce(credential_id,''), documento, uc, status, started_at, coalesce(finished_at,''), coalesce(error_message,''), coalesce(raw_response_json,'')
		   from sync_runs where id = ?`),
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

func (s *PostgresStore) GetLatestSyncRunByUC(uc string) (*SyncRunView, error) {
	row := s.db.QueryRow(
		pgQuery(`select id, coalesce(credential_id,''), documento, uc, status, started_at, coalesce(finished_at,''), coalesce(error_message,''), coalesce(raw_response_json,'')
		   from sync_runs where uc = ? order by started_at desc limit 1`),
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

func (s *PostgresStore) listInvoiceItems(invoiceID string) ([]map[string]any, error) {
	rows, err := s.db.Query(
		pgQuery(`select coalesce(raw_json,'') from invoice_items where invoice_id = ? order by descricao asc`),
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

func (s *PostgresStore) PersistSyncResult(in PersistSyncInput) (PersistSyncResult, error) {
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
		pgQuery(`insert into sync_runs (
			id, credential_id, documento, uc, status, started_at, finished_at, error_message, raw_response_json
		) values (?, ?, ?, ?, ?, ?, ?, ?, ?)`),
		syncRunID, nullableString(in.CredentialID), in.Documento, in.UC, status, now, now, nullableString(in.ErrorMessage), nullableString(rawResponseJSON),
	)
	if err != nil {
		return PersistSyncResult{}, err
	}

	if in.UCRecord != nil {
		_, err = tx.Exec(
			pgQuery(`insert into consumer_units (
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
				updated_at=excluded.updated_at`),
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

	invoiceID, err := upsertInvoicePG(tx, now, in)
	if err != nil {
		return PersistSyncResult{}, err
	}
	if err := replaceInvoiceChildrenPG(tx, now, invoiceID, syncRunID, in); err != nil {
		return PersistSyncResult{}, err
	}
	if err := tx.Commit(); err != nil {
		return PersistSyncResult{}, err
	}
	return PersistSyncResult{SyncRunID: syncRunID, InvoiceID: invoiceID}, nil
}

func upsertInvoicePG(tx *sql.Tx, now string, in PersistSyncInput) (string, error) {
	var existingID string
	err := tx.QueryRow(
		pgQuery(`select id from invoices where uc = ? and numero_fatura = ?`),
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
		pgQuery(`insert into invoices (
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
			updated_at=excluded.updated_at`),
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

func replaceInvoiceChildrenPG(tx *sql.Tx, now string, invoiceID string, syncRunID string, in PersistSyncInput) error {
	for _, stmt := range []string{
		`delete from invoice_api_snapshots where invoice_id = ?`,
		`delete from invoice_documents where invoice_id = ?`,
		`delete from invoice_items where invoice_id = ?`,
		`delete from invoice_extraction_results where invoice_id = ?`,
		`delete from invoice_field_sources where invoice_id = ?`,
	} {
		if _, err := tx.Exec(pgQuery(stmt), invoiceID); err != nil {
			return err
		}
	}

	if _, err := tx.Exec(
		pgQuery(`insert into invoice_api_snapshots (
			id, invoice_id, sync_run_id, fatura_json, historico_json, dados_pagamento_json, created_at
		) values (?, ?, ?, ?, ?, ?, ?)`),
		newID(), invoiceID, syncRunID, nullableString(mustJSON(in.Fatura)), nullableString(mustJSON(in.Historico)), nullableString(mustJSON(in.DadosPagamento)), now,
	); err != nil {
		return err
	}

	if in.PDF != nil {
		if _, err := tx.Exec(
			pgQuery(`insert into invoice_documents (
				id, invoice_id, sync_run_id, file_name, file_extension, file_size, file_data_base64, storage_uri, created_at
			) values (?, ?, ?, ?, ?, ?, ?, ?, ?)`),
			newID(), invoiceID, syncRunID, nullableString(in.PDF.FileName), nullableString(in.PDF.FileExtension), nullableString(in.PDF.FileSize), nullableString(in.PDF.FileData), nil, now,
		); err != nil {
			return err
		}
	}

	for _, item := range extractItems(in.BillingRecord) {
		if _, err := tx.Exec(
			pgQuery(`insert into invoice_items (
				id, invoice_id, descricao, quantidade, quantidade_residual, quantidade_faturada,
				tarifa, valor, base_icms, aliq_icms, icms, valor_total, raw_json, created_at
			) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`),
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
			pgQuery(`insert into invoice_extraction_results (
				id, invoice_id, sync_run_id, status, fields_json, source_map_json,
				confidence_map_json, warnings_json, artifacts_json, created_at
			) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`),
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
			pgQuery(`insert into invoice_field_sources (
				id, invoice_id, field_path, source, confidence, created_at
			) values (?, ?, ?, ?, ?, ?)`),
			newID(), invoiceID, field, source, confidenceMap[field], now,
		); err != nil {
			return err
		}
	}
	return nil
}

func pgQuery(query string) string {
	var b strings.Builder
	b.Grow(len(query) + 16)
	idx := 0
	for _, r := range query {
		if r == '?' {
			idx++
			_, _ = b.WriteString(fmt.Sprintf("$%d", idx))
			continue
		}
		_, _ = b.WriteRune(r)
	}
	return b.String()
}
