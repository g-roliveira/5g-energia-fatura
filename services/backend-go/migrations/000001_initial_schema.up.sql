-- =====================================================================
-- Migration 001 UP — schema core (cadastro) + schema billing
-- =====================================================================
-- Este é o Postgres do BACKOFFICE. Complementa (não substitui) o SQLite
-- do backend-go, que permanece responsável pelo domínio de INTEGRAÇÃO
-- com a concessionária (credentials, sessions, sync_runs, invoices).
--
-- Este banco é o que o front (Next.js) vai ler diretamente para
-- cadastro e, via BFF do backend-go, para faturamento.

CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Schema core removed: tables now in public
-- Schema billing removed: tables now in public

-- =====================================================================
-- SCHEMA core — cadastro de clientes, UCs, usuários do backoffice
-- =====================================================================

CREATE TABLE public.app_user (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           TEXT NOT NULL UNIQUE,
    name            TEXT NOT NULL,
    password_hash   TEXT NOT NULL,
    role            TEXT NOT NULL CHECK (role IN ('admin', 'operator', 'reviewer')),
    active          BOOLEAN NOT NULL DEFAULT TRUE,
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_app_user_active ON public.app_user(email) WHERE active = TRUE;

-- Cliente final (condomínio, PF, PJ)
CREATE TABLE public.customer (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tipo_pessoa     TEXT NOT NULL CHECK (tipo_pessoa IN ('PF', 'PJ')),
    nome_razao      TEXT NOT NULL,
    nome_fantasia   TEXT,
    cpf_cnpj        TEXT UNIQUE,
    email           TEXT,
    phone           TEXT,
    tipo_cliente    TEXT CHECK (tipo_cliente IN (
        'residencial','condominio','empresa','imobiliaria','outro'
    )),
    notes           TEXT,
    status          TEXT NOT NULL DEFAULT 'active'
                      CHECK (status IN ('active','inactive','prospect','archived')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    archived_at     TIMESTAMPTZ
);

-- Endereço do cliente
CREATE TABLE public.address (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id     UUID NOT NULL REFERENCES public.customer(id) ON DELETE CASCADE,
    cep             TEXT,
    logradouro      TEXT,
    numero          TEXT,
    complemento     TEXT,
    bairro          TEXT,
    cidade          TEXT,
    uf              TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_address_customer ON public.address(customer_id);

-- Unidade consumidora (cadastral — distinta da consumer_units da API no
-- SQLite do backend-go). Aqui é o vínculo comercial cliente ↔ UC.
CREATE TABLE public.consumer_unit (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id              UUID REFERENCES public.customer(id) ON DELETE SET NULL,
    uc_code                  TEXT NOT NULL UNIQUE,
    distribuidora            TEXT NOT NULL DEFAULT 'neoenergia_ba',
    apelido                  TEXT,
    classe_consumo           TEXT,
    endereco_unidade         TEXT,
    cidade                   TEXT,
    uf                       TEXT,
    grupo_tensao             TEXT CHECK (grupo_tensao IN ('monofasico','bifasico','trifasico')),
    consumo_minimo_kwh       INTEGER CHECK (consumo_minimo_kwh >= 0),
    ativa                    BOOLEAN NOT NULL DEFAULT TRUE,
    -- Ponte lógica (não-FK) pra credentials.id no SQLite do backend-go
    sync_credential_id       TEXT,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_consumer_unit_customer ON public.consumer_unit(customer_id);
CREATE INDEX idx_consumer_unit_active   ON public.consumer_unit(uc_code) WHERE ativa;

-- Credencial de integração (metadados; senha mesmo fica criptografada no SQLite)
CREATE TABLE public.credential_link (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id         UUID NOT NULL REFERENCES public.customer(id) ON DELETE CASCADE,
    label               TEXT NOT NULL,
    documento_masked    TEXT NOT NULL,
    uf                  TEXT NOT NULL,
    tipo_acesso         TEXT NOT NULL,
    -- UUID retornado por POST /v1/credentials do backend-go
    go_credential_id    TEXT NOT NULL UNIQUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_credential_link_customer ON public.credential_link(customer_id);
-- =====================================================================
-- Migration 001 UP — parte 2: schema billing (faturamento)
-- =====================================================================

-- ---------------------------------------------------------------------
-- CONTRATO (VERSIONADO POR VIGÊNCIA)
-- ---------------------------------------------------------------------
-- Um contrato NUNCA é UPDATE — mudança = INSERT novo registro com
-- vigencia_inicio nova, UPDATE de vigencia_fim no anterior.
-- Uma UC só pode ter 1 contrato ativo (vigencia_fim IS NULL).

CREATE TABLE public.contract (
    id                                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id                          UUID NOT NULL REFERENCES public.customer(id),
    consumer_unit_id                     UUID NOT NULL REFERENCES public.consumer_unit(id),
    vigencia_inicio                      DATE NOT NULL,
    vigencia_fim                         DATE,
    desconto_percentual                  NUMERIC(5,4) NOT NULL
                                           CHECK (desconto_percentual > 0
                                                  AND desconto_percentual <= 1),
    ip_faturamento_mode                  TEXT NOT NULL DEFAULT 'fixed'
                                           CHECK (ip_faturamento_mode IN ('fixed','percent')),
    ip_faturamento_valor                 NUMERIC(12,4) NOT NULL DEFAULT 0,
    ip_faturamento_percent               NUMERIC(5,4) NOT NULL DEFAULT 0,
    bandeira_com_desconto                BOOLEAN NOT NULL DEFAULT FALSE,
    custo_disponibilidade_sempre_cobrado BOOLEAN NOT NULL DEFAULT TRUE,
    notes                                TEXT,
    status                               TEXT NOT NULL DEFAULT 'active'
                                           CHECK (status IN ('draft','active','ended')),
    created_at                           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by                           UUID REFERENCES public.app_user(id),
    updated_at                           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT contract_vigencia_coerente
      CHECK (vigencia_fim IS NULL OR vigencia_fim >= vigencia_inicio)
);

-- Exatamente 1 contrato ativo por UC (partial unique index)
CREATE UNIQUE INDEX idx_contract_one_active_per_uc
    ON public.contract(consumer_unit_id)
    WHERE vigencia_fim IS NULL AND status = 'active';

CREATE INDEX idx_contract_uc_vigencia
    ON public.contract(consumer_unit_id, vigencia_inicio DESC);

-- ---------------------------------------------------------------------
-- COMPETÊNCIA (CICLO MENSAL)
-- ---------------------------------------------------------------------

CREATE TABLE public.billing_cycle (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    year            SMALLINT NOT NULL CHECK (year BETWEEN 2020 AND 2100),
    month           SMALLINT NOT NULL CHECK (month BETWEEN 1 AND 12),
    reference_date  DATE NOT NULL,
    status          TEXT NOT NULL DEFAULT 'open'
                      CHECK (status IN ('open','syncing','processing','review','approved','closed')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES public.app_user(id),
    closed_at       TIMESTAMPTZ,
    closed_by       UUID REFERENCES public.app_user(id),
    UNIQUE (year, month)
);

CREATE INDEX idx_billing_cycle_status ON public.billing_cycle(status);

-- ---------------------------------------------------------------------
-- REFERÊNCIA PARA FATURA ORIGINAL (que mora no SQLite do backend-go)
-- ---------------------------------------------------------------------
-- O utility_invoice_ref existe porque a fatura original da Coelba vive no
-- SQLite do backend-go (domínio de integração). Aqui só guardamos o
-- ponteiro lógico + os dados que o faturamento precisa ter "em mãos"
-- no Postgres pra join/query sem cross-database.

CREATE TABLE public.utility_invoice_ref (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    consumer_unit_id        UUID NOT NULL REFERENCES public.consumer_unit(id),
    billing_cycle_id        UUID NOT NULL REFERENCES public.billing_cycle(id),
    -- ID no SQLite do backend-go (invoices.id)
    sync_invoice_id         TEXT NOT NULL,
    -- ID no SQLite do backend-go (sync_runs.id) do sync que trouxe isso
    sync_run_id             TEXT,
    numero_fatura           TEXT,
    mes_referencia          TEXT,
    valor_total_coelba      NUMERIC(12,2),
    status_fatura           TEXT,
    data_emissao            DATE,
    data_vencimento         DATE,
    data_inicio_periodo     DATE,
    data_fim_periodo        DATE,
    completeness_status     TEXT,
    completeness_missing    TEXT[],
    extractor_status        TEXT,
    extractor_confidence    NUMERIC(3,2),
    -- Snapshot do BillingRecord (sync.BillingRecord) no momento do sync
    -- Guardamos aqui também (redundante com SQLite) para permitir
    -- auditoria do que o motor viu sem depender de cross-database join
    billing_record_snapshot JSONB,
    synced_at               TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (consumer_unit_id, billing_cycle_id)
);

CREATE INDEX idx_invoice_ref_cycle    ON public.utility_invoice_ref(billing_cycle_id);
CREATE INDEX idx_invoice_ref_snapshot ON public.utility_invoice_ref USING GIN (billing_record_snapshot);

-- ---------------------------------------------------------------------
-- ITEM NORMALIZADO DA FATURA
-- ---------------------------------------------------------------------
-- Produzido pelo packages/normalizer a partir do billing_record. É o que
-- alimenta o motor de cálculo. ignored_in_calc=TRUE para itens preservados
-- por auditoria mas que não entram no cálculo (IRRF, reativo, bandeira verde).

CREATE TABLE public.utility_invoice_item (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    utility_invoice_ref_id  UUID NOT NULL REFERENCES public.utility_invoice_ref(id) ON DELETE CASCADE,
    type                    TEXT NOT NULL CHECK (type IN (
                                'tusd_fio','tusd_energia','energia_injetada',
                                'bandeira','ip_coelba',
                                'reativo_excedente','tributo_retido'
                            )),
    description             TEXT NOT NULL,
    quantidade              NUMERIC(14,4),
    unidade                 TEXT,
    preco_unitario          NUMERIC(14,8),
    valor_total             NUMERIC(12,4) NOT NULL,
    ignored_in_calc         BOOLEAN NOT NULL DEFAULT FALSE,
    source                  TEXT CHECK (source IN ('api','pymupdf','mistral','manual')),
    confidence              NUMERIC(3,2),
    order_index             INTEGER NOT NULL DEFAULT 0,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_invoice_item_ref ON public.utility_invoice_item(utility_invoice_ref_id, type);

-- ---------------------------------------------------------------------
-- SCEE EXTRAÍDO DO RODAPÉ
-- ---------------------------------------------------------------------

CREATE TABLE public.utility_invoice_scee (
    utility_invoice_ref_id  UUID PRIMARY KEY REFERENCES public.utility_invoice_ref(id) ON DELETE CASCADE,
    layout                  TEXT NOT NULL CHECK (layout IN ('mmgd_legado','mmgd_transicao','scee_moderno')),
    energia_injetada_kwh    NUMERIC(14,4),
    excedente_kwh           NUMERIC(14,4),
    creditos_utilizados     NUMERIC(14,4),
    saldo_proximo_ciclo     NUMERIC(14,4),
    raw_text                TEXT,
    confidence              NUMERIC(3,2),
    extracted_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ---------------------------------------------------------------------
-- CÁLCULO (COM SNAPSHOTS IMUTÁVEIS)
-- ---------------------------------------------------------------------
-- Três snapshots JSONB congelam:
--   contract_snapshot_json: a regra que vigorava no momento do cálculo
--   inputs_snapshot_json:   os itens que entraram no motor
--   result_snapshot_json:   o output completo com breakdown por linha
-- Mudou contrato? Nova version. Mudou ajuste? Nova version. Original fica
-- marcado como superseded, nunca é sobrescrito.

CREATE TABLE public.billing_calculation (
    id                        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    utility_invoice_ref_id    UUID NOT NULL REFERENCES public.utility_invoice_ref(id),
    billing_cycle_id          UUID NOT NULL REFERENCES public.billing_cycle(id),
    consumer_unit_id          UUID NOT NULL REFERENCES public.consumer_unit(id),
    contract_id               UUID NOT NULL REFERENCES public.contract(id),
    contract_snapshot_json    JSONB NOT NULL,
    inputs_snapshot_json      JSONB NOT NULL,
    result_snapshot_json      JSONB NOT NULL,
    total_sem_desconto        NUMERIC(12,4) NOT NULL,
    total_com_desconto        NUMERIC(12,4) NOT NULL,
    economia_rs               NUMERIC(12,4) NOT NULL,
    economia_pct              NUMERIC(5,4) NOT NULL,
    status                    TEXT NOT NULL DEFAULT 'draft'
                                CHECK (status IN ('draft','needs_review','approved','superseded')),
    needs_review_reasons      TEXT[],
    version                   INTEGER NOT NULL DEFAULT 1 CHECK (version >= 1),
    calculated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    approved_at               TIMESTAMPTZ,
    approved_by               UUID REFERENCES public.app_user(id),
    UNIQUE (utility_invoice_ref_id, version)
);

CREATE INDEX idx_billing_calc_cycle_status
    ON public.billing_calculation(billing_cycle_id, status);
CREATE INDEX idx_billing_calc_current
    ON public.billing_calculation(utility_invoice_ref_id)
    WHERE status != 'superseded';

-- ---------------------------------------------------------------------
-- AJUSTE MANUAL (IMUTÁVEL — UM REGISTRO POR EDIÇÃO)
-- ---------------------------------------------------------------------

CREATE TABLE public.manual_adjustment (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    billing_calculation_id   UUID NOT NULL REFERENCES public.billing_calculation(id) ON DELETE CASCADE,
    field_path               TEXT NOT NULL,
    old_value                JSONB,
    new_value                JSONB NOT NULL,
    reason                   TEXT NOT NULL,
    created_by               UUID NOT NULL REFERENCES public.app_user(id),
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_manual_adjustment_calc
    ON public.manual_adjustment(billing_calculation_id);

-- ---------------------------------------------------------------------
-- DOCUMENTO GERADO (PDF DO CLIENTE)
-- ---------------------------------------------------------------------

CREATE TABLE public.generated_document (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    billing_calculation_id   UUID NOT NULL REFERENCES public.billing_calculation(id) ON DELETE CASCADE,
    type                     TEXT NOT NULL CHECK (type IN ('customer_invoice_pdf','preview_pdf')),
    file_path                TEXT NOT NULL,
    checksum_sha256          TEXT NOT NULL,
    version                  INTEGER NOT NULL DEFAULT 1,
    generated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    generated_by             UUID REFERENCES public.app_user(id),
    UNIQUE (billing_calculation_id, type, version)
);

-- ---------------------------------------------------------------------
-- FILA DE JOBS (VIA POSTGRES — FOR UPDATE SKIP LOCKED)
-- ---------------------------------------------------------------------
-- Workers (internal/worker/ do backend-go) fazem:
--   SELECT ... FROM public.sync_job
--     WHERE status='pending' AND scheduled_for <= NOW()
--     ORDER BY scheduled_for
--     LIMIT 1 FOR UPDATE SKIP LOCKED;
-- Zero Kafka, zero Redis, zero Trigger.dev. Postgres resolve.

CREATE TABLE public.sync_job (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type              TEXT NOT NULL CHECK (type IN (
                        'sync_uc','calculate','generate_pdf','recalculate_cycle'
                      )),
    payload_json      JSONB NOT NULL,
    status            TEXT NOT NULL DEFAULT 'pending'
                        CHECK (status IN ('pending','running','success','failed','retrying')),
    retry_count       INTEGER NOT NULL DEFAULT 0,
    max_retries       INTEGER NOT NULL DEFAULT 3,
    error_message     TEXT,
    idempotency_key   TEXT NOT NULL,
    scheduled_for     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at        TIMESTAMPTZ,
    finished_at       TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Idempotência: não pode haver 2 jobs "vivos" com mesma chave
CREATE UNIQUE INDEX idx_sync_job_idempotency_live
    ON public.sync_job(type, idempotency_key)
    WHERE status IN ('pending','running','retrying','success');

-- Worker hot path: "me dá o próximo job"
CREATE INDEX idx_sync_job_ready
    ON public.sync_job(scheduled_for)
    WHERE status = 'pending';

-- ---------------------------------------------------------------------
-- AUDIT LOG (TRILHA DE AUDITORIA)
-- ---------------------------------------------------------------------

CREATE TABLE public.audit_log (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_type     TEXT NOT NULL CHECK (actor_type IN ('user','system','job')),
    actor_id       TEXT,
    entity_type    TEXT NOT NULL,
    entity_id      UUID NOT NULL,
    action         TEXT NOT NULL,
    before_json    JSONB,
    after_json     JSONB,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_entity ON public.audit_log(entity_type, entity_id, created_at DESC);
CREATE INDEX idx_audit_time   ON public.audit_log(created_at DESC);

-- ---------------------------------------------------------------------
-- TRIGGERS updated_at
-- ---------------------------------------------------------------------

CREATE OR REPLACE FUNCTION public.set_updated_at() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_app_user_updated BEFORE UPDATE ON public.app_user
    FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();
CREATE TRIGGER trg_customer_updated BEFORE UPDATE ON public.customer
    FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();
CREATE TRIGGER trg_address_updated BEFORE UPDATE ON public.address
    FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();
CREATE TRIGGER trg_consumer_unit_updated BEFORE UPDATE ON public.consumer_unit
    FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();
CREATE TRIGGER trg_credential_link_updated BEFORE UPDATE ON public.credential_link
    FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();
CREATE TRIGGER trg_contract_updated BEFORE UPDATE ON public.contract
    FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();
CREATE TRIGGER trg_invoice_ref_updated BEFORE UPDATE ON public.utility_invoice_ref
    FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();
