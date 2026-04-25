-- =====================================================================
-- Migration 002 UP — Unificação de backend em Postgres (integração + jobs)
-- =====================================================================
-- Este migration cria as tabelas usadas por:
-- - services/backend-go (domínio de integração Neoenergia)
-- - src/fatura (jobs batch e parser OCR)
--
-- Convenção: tabelas em public para compatibilidade com SQLAlchemy e
-- com o store Go atual sem alteração de search_path.

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ---------------------------------------------------------------------
-- backend-go integration domain
-- ---------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS credentials (
    id               TEXT PRIMARY KEY,
    label            TEXT NOT NULL,
    documento_cipher TEXT NOT NULL,
    documento_nonce  TEXT NOT NULL,
    senha_cipher     TEXT NOT NULL,
    senha_nonce      TEXT NOT NULL,
    uf               TEXT NOT NULL,
    tipo_acesso      TEXT NOT NULL,
    key_version      TEXT NOT NULL,
    created_at       TEXT NOT NULL,
    updated_at       TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
    id                 TEXT PRIMARY KEY,
    credential_id      TEXT NOT NULL REFERENCES credentials(id),
    bearer_token_cipher TEXT NOT NULL,
    bearer_token_nonce  TEXT NOT NULL,
    created_at          TEXT NOT NULL,
    updated_at          TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_sessions_credential_created_at ON sessions(credential_id, created_at DESC);

CREATE TABLE IF NOT EXISTS consumer_units (
    uc            TEXT PRIMARY KEY,
    credential_id TEXT REFERENCES credentials(id),
    status        TEXT,
    nome_cliente  TEXT,
    instalacao    TEXT,
    contrato      TEXT,
    grupo_tensao  TEXT,
    endereco_json TEXT,
    imovel_json   TEXT,
    created_at    TEXT NOT NULL,
    updated_at    TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sync_runs (
    id                TEXT PRIMARY KEY,
    credential_id     TEXT REFERENCES credentials(id),
    documento         TEXT NOT NULL,
    uc                TEXT NOT NULL,
    status            TEXT NOT NULL,
    started_at        TEXT NOT NULL,
    finished_at       TEXT,
    error_message     TEXT,
    raw_response_json TEXT
);
CREATE INDEX IF NOT EXISTS idx_sync_runs_uc_started_at ON sync_runs(uc, started_at DESC);

CREATE TABLE IF NOT EXISTS invoices (
    id                         TEXT PRIMARY KEY,
    uc                         TEXT NOT NULL,
    numero_fatura              TEXT NOT NULL,
    mes_referencia             TEXT NOT NULL,
    status_fatura              TEXT,
    valor_total                TEXT,
    codigo_barras              TEXT,
    data_emissao               TEXT,
    data_vencimento            TEXT,
    data_pagamento             TEXT,
    data_inicio_periodo        TEXT,
    data_fim_periodo           TEXT,
    completeness_status        TEXT,
    completeness_missing_json  TEXT,
    billing_record_json        TEXT,
    document_record_json       TEXT,
    created_at                 TEXT NOT NULL,
    updated_at                 TEXT NOT NULL,
    UNIQUE (uc, numero_fatura)
);
CREATE INDEX IF NOT EXISTS idx_invoices_uc_mes ON invoices(uc, mes_referencia);

CREATE TABLE IF NOT EXISTS invoice_api_snapshots (
    id                  TEXT PRIMARY KEY,
    invoice_id          TEXT NOT NULL REFERENCES invoices(id),
    sync_run_id         TEXT NOT NULL REFERENCES sync_runs(id),
    fatura_json         TEXT,
    historico_json      TEXT,
    dados_pagamento_json TEXT,
    created_at          TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS invoice_documents (
    id               TEXT PRIMARY KEY,
    invoice_id       TEXT NOT NULL REFERENCES invoices(id),
    sync_run_id      TEXT NOT NULL REFERENCES sync_runs(id),
    file_name        TEXT,
    file_extension   TEXT,
    file_size        TEXT,
    file_data_base64 TEXT,
    storage_uri      TEXT,
    created_at       TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS invoice_items (
    id                  TEXT PRIMARY KEY,
    invoice_id          TEXT NOT NULL REFERENCES invoices(id),
    descricao           TEXT,
    quantidade          TEXT,
    quantidade_residual TEXT,
    quantidade_faturada TEXT,
    tarifa              TEXT,
    valor               TEXT,
    base_icms           TEXT,
    aliq_icms           TEXT,
    icms                TEXT,
    valor_total         TEXT,
    raw_json            TEXT,
    created_at          TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS invoice_extraction_results (
    id                  TEXT PRIMARY KEY,
    invoice_id          TEXT NOT NULL REFERENCES invoices(id),
    sync_run_id         TEXT NOT NULL REFERENCES sync_runs(id),
    status              TEXT,
    fields_json         TEXT,
    source_map_json     TEXT,
    confidence_map_json TEXT,
    warnings_json       TEXT,
    artifacts_json      TEXT,
    created_at          TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS invoice_field_sources (
    id         TEXT PRIMARY KEY,
    invoice_id TEXT NOT NULL REFERENCES invoices(id),
    field_path TEXT NOT NULL,
    source     TEXT,
    confidence DOUBLE PRECISION,
    created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_invoice_field_sources_invoice ON invoice_field_sources(invoice_id);

-- ---------------------------------------------------------------------
-- src/fatura jobs domain
-- ---------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS clientes (
    id              BIGSERIAL PRIMARY KEY,
    codigo          VARCHAR(50) NOT NULL DEFAULT '',
    cpf             VARCHAR(20),
    cnpj            VARCHAR(25),
    nome            VARCHAR(200) NOT NULL,
    classificacao   VARCHAR(100),
    tensao_nominal  VARCHAR(50),
    endereco        VARCHAR(500),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS contas (
    id               BIGSERIAL PRIMARY KEY,
    uc               VARCHAR(20) NOT NULL,
    mes              INTEGER NOT NULL,
    ano              INTEGER NOT NULL,
    valor            VARCHAR(20) NOT NULL,
    vencimento       DATE NOT NULL,
    numero_dias      INTEGER,
    codigo_barras    VARCHAR(60),
    pdf_path         VARCHAR(500),
    parsed_at        TIMESTAMPTZ,
    cliente_id       BIGINT REFERENCES clientes(id),
    composicao_json  TEXT,
    consumo_json     TEXT,
    energia_json     TEXT,
    nota_fiscal_json TEXT,
    ocr_json         TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_conta_uc_mes_ano UNIQUE (uc, mes, ano)
);
CREATE INDEX IF NOT EXISTS idx_contas_uc ON contas(uc);

CREATE TABLE IF NOT EXISTS itens_fatura (
    id          BIGSERIAL PRIMARY KEY,
    conta_id    BIGINT NOT NULL REFERENCES contas(id) ON DELETE CASCADE,
    codigo      VARCHAR(20) NOT NULL DEFAULT '',
    descricao   VARCHAR(300) NOT NULL DEFAULT '',
    quantidade  VARCHAR(20),
    tarifa      VARCHAR(20),
    valor       VARCHAR(20),
    base_icms   VARCHAR(20),
    aliq_icms   VARCHAR(20),
    icms        VARCHAR(20),
    valor_total VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS processamento_log (
    id         BIGSERIAL PRIMARY KEY,
    uc         VARCHAR(20) NOT NULL,
    mes        INTEGER NOT NULL,
    ano        INTEGER NOT NULL,
    status     VARCHAR(30) NOT NULL,
    mensagem   VARCHAR(1000),
    tentativa  INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_processamento_log_uc ON processamento_log(uc);

CREATE TABLE IF NOT EXISTS jobs (
    id              VARCHAR(36) PRIMARY KEY,
    kind            VARCHAR(50) NOT NULL DEFAULT 'neoenergia_fatura',
    status          VARCHAR(30) NOT NULL,
    request_json    TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at      TIMESTAMPTZ,
    finished_at     TIMESTAMPTZ,
    total_items     INTEGER NOT NULL DEFAULT 0,
    completed_items INTEGER NOT NULL DEFAULT 0,
    success_items   INTEGER NOT NULL DEFAULT 0,
    error_items     INTEGER NOT NULL DEFAULT 0,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status);

CREATE TABLE IF NOT EXISTS job_items (
    id              BIGSERIAL PRIMARY KEY,
    job_id          VARCHAR(36) NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    uc              VARCHAR(20) NOT NULL,
    nome            VARCHAR(200) NOT NULL DEFAULT '',
    status          VARCHAR(30) NOT NULL DEFAULT 'queued',
    mensagem        VARCHAR(1000),
    error_type      VARCHAR(100),
    pdf_path        VARCHAR(500),
    screenshot_path VARCHAR(500),
    html_path       VARCHAR(500),
    step_name       VARCHAR(100),
    mes             INTEGER,
    ano             INTEGER,
    valor           VARCHAR(50),
    conta_id        BIGINT REFERENCES contas(id),
    attempts        INTEGER NOT NULL DEFAULT 0,
    started_at      TIMESTAMPTZ,
    finished_at     TIMESTAMPTZ,
    result_json     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_job_item_job_uc UNIQUE (job_id, uc)
);
CREATE INDEX IF NOT EXISTS idx_job_items_job_id ON job_items(job_id);
CREATE INDEX IF NOT EXISTS idx_job_items_uc ON job_items(uc);
CREATE INDEX IF NOT EXISTS idx_job_items_status ON job_items(status);
