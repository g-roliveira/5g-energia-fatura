-- =====================================================================
-- Migration 003 UP — Schema integration (domínio de scraping/sync)
-- =====================================================================
-- Separa as tabelas de integração com Neoenergia/Coelba em schema
-- próprio, isolado do cadastro (core) e faturamento (billing).

CREATE SCHEMA IF NOT EXISTS integration;

-- ---------------------------------------------------------------------
-- CREDENTIALS (credenciais de acesso à Neoenergia)
-- ---------------------------------------------------------------------
CREATE TABLE integration.credentials (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    label            TEXT NOT NULL,
    documento_cipher TEXT NOT NULL,
    documento_nonce  TEXT NOT NULL,
    senha_cipher     TEXT NOT NULL,
    senha_nonce      TEXT NOT NULL,
    uf               TEXT NOT NULL,
    tipo_acesso      TEXT NOT NULL,
    key_version      TEXT NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_credentials_label ON integration.credentials(label);

-- ---------------------------------------------------------------------
-- SESSIONS (sessões ativas do Playwright)
-- ---------------------------------------------------------------------
CREATE TABLE integration.sessions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    credential_id       UUID NOT NULL REFERENCES integration.credentials(id) ON DELETE CASCADE,
    bearer_token_cipher TEXT NOT NULL,
    bearer_token_nonce  TEXT NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_credential ON integration.sessions(credential_id, created_at DESC);

-- ---------------------------------------------------------------------
-- CONSUMER UNITS (UCs descobertas via scraping — mirror do core.consumer_unit)
-- ---------------------------------------------------------------------
CREATE TABLE integration.consumer_units (
    uc             TEXT PRIMARY KEY,
    credential_id  UUID REFERENCES integration.credentials(id) ON DELETE SET NULL,
    status         TEXT,
    nome_cliente   TEXT,
    instalacao     TEXT,
    contrato       TEXT,
    grupo_tensao   TEXT,
    endereco_json  JSONB,
    imovel_json    JSONB,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_integration_cus_credential ON integration.consumer_units(credential_id);

-- ---------------------------------------------------------------------
-- SYNC RUNS (registro de execuções de sync)
-- ---------------------------------------------------------------------
CREATE TABLE integration.sync_runs (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    credential_id     UUID REFERENCES integration.credentials(id) ON DELETE SET NULL,
    documento         TEXT NOT NULL,
    uc                TEXT NOT NULL,
    status            TEXT NOT NULL DEFAULT 'pending',
    step              TEXT,                      -- login, navigate, download, extract
    error_message     TEXT,
    error_context     JSONB,                     -- screenshot, HTML, etc
    raw_response_json JSONB,
    started_at        TIMESTAMPTZ,
    finished_at       TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sync_runs_uc ON integration.sync_runs(uc, created_at DESC);
CREATE INDEX idx_sync_runs_status ON integration.sync_runs(status, created_at DESC);

-- ---------------------------------------------------------------------
-- RAW INVOICES (faturas brutas baixadas — substitui invoices do SQLite)
-- ---------------------------------------------------------------------
CREATE TABLE integration.raw_invoices (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uc                      TEXT NOT NULL,
    numero_fatura           TEXT NOT NULL,
    mes_referencia          TEXT NOT NULL,
    status_fatura           TEXT,
    valor_total             NUMERIC(12,2),
    codigo_barras           TEXT,
    data_emissao            DATE,
    data_vencimento         DATE,
    data_pagamento          DATE,
    data_inicio_periodo     DATE,
    data_fim_periodo        DATE,
    completeness_status     TEXT,
    completeness_missing    TEXT[],
    billing_record_json     JSONB,
    document_record_json    JSONB,
    pdf_bytes               BYTEA,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (uc, numero_fatura)
);

CREATE INDEX idx_raw_invoices_uc_mes ON integration.raw_invoices(uc, mes_referencia);
CREATE INDEX idx_raw_invoices_billing ON integration.raw_invoices USING GIN (billing_record_json);

-- ---------------------------------------------------------------------
-- RAW INVOICE ITEMS (itens normalizados da fatura)
-- ---------------------------------------------------------------------
CREATE TABLE integration.raw_invoice_items (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    raw_invoice_id      UUID NOT NULL REFERENCES integration.raw_invoices(id) ON DELETE CASCADE,
    type                TEXT NOT NULL CHECK (type IN (
                            'tusd_fio','tusd_energia','energia_injetada',
                            'bandeira','ip_coelba',
                            'reativo_excedente','tributo_retido'
                        )),
    description         TEXT NOT NULL,
    quantidade          NUMERIC(14,4),
    preco_unitario      NUMERIC(14,8),
    valor_total         NUMERIC(12,4) NOT NULL,
    ignored_in_calc     BOOLEAN NOT NULL DEFAULT FALSE,
    source              TEXT CHECK (source IN ('api','pymupdf','mistral','manual')),
    confidence          NUMERIC(3,2),
    order_index         INTEGER NOT NULL DEFAULT 0,
    raw_json            JSONB,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_raw_invoice_items_invoice ON integration.raw_invoice_items(raw_invoice_id, type);

-- ---------------------------------------------------------------------
-- JOB QUEUE (fila de jobs genérica — worker pool)
-- ---------------------------------------------------------------------
CREATE TABLE integration.jobs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_type        TEXT NOT NULL,   -- sync_uc, extract, calculate, generate_pdf
    status          TEXT NOT NULL DEFAULT 'pending'
                        CHECK (status IN ('pending','claimed','running','completed','failed')),
    payload         JSONB NOT NULL,
    result          JSONB,
    error_message   TEXT,
    retry_count     INT NOT NULL DEFAULT 0,
    max_retries     INT NOT NULL DEFAULT 3,
    claimed_by      TEXT,            -- worker ID
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    claimed_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ
);

-- Índice crítico para worker pool (FOR UPDATE SKIP LOCKED)
CREATE INDEX idx_jobs_pending ON integration.jobs(created_at ASC)
    WHERE status = 'pending';

CREATE INDEX idx_jobs_running ON integration.jobs(claimed_at)
    WHERE status = 'running';

-- ---------------------------------------------------------------------
-- Triggers updated_at
-- ---------------------------------------------------------------------
CREATE OR REPLACE FUNCTION integration.set_updated_at() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_credentials_updated BEFORE UPDATE ON integration.credentials
    FOR EACH ROW EXECUTE FUNCTION integration.set_updated_at();
CREATE TRIGGER trg_sessions_updated BEFORE UPDATE ON integration.sessions
    FOR EACH ROW EXECUTE FUNCTION integration.set_updated_at();
CREATE TRIGGER trg_consumer_units_updated BEFORE UPDATE ON integration.consumer_units
    FOR EACH ROW EXECUTE FUNCTION integration.set_updated_at();
CREATE TRIGGER trg_raw_invoices_updated BEFORE UPDATE ON integration.raw_invoices
    FOR EACH ROW EXECUTE FUNCTION integration.set_updated_at();
