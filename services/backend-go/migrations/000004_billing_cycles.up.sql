-- =====================================================================
-- Migration 002 UP — Billing cycles, notifications, performance indices
-- =====================================================================

-- ---------------------------------------------------------------------
-- CYCLE-CONSUMER-UNIT JOIN TABLE
-- Necessário para associar UCs a um ciclo antes do sync ocorrer.
-- Também guarda o status de cada UC no ciclo (pending/synced/calculated).
-- ---------------------------------------------------------------------
CREATE TABLE public.cycle_consumer_unit (
    billing_cycle_id   UUID NOT NULL REFERENCES public.billing_cycle(id) ON DELETE CASCADE,
    consumer_unit_id   UUID NOT NULL REFERENCES public.consumer_unit(id) ON DELETE CASCADE,
    synced_at          TIMESTAMPTZ,
    calculation_id     UUID REFERENCES public.billing_calculation(id) ON DELETE SET NULL,
    status             TEXT NOT NULL DEFAULT 'pending'
                           CHECK (status IN ('pending','synced','calculated','approved','error')),
    error_message      TEXT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (billing_cycle_id, consumer_unit_id)
);

CREATE INDEX idx_ccu_cycle_status ON public.cycle_consumer_unit(billing_cycle_id, status);
CREATE INDEX idx_ccu_consumer_unit ON public.cycle_consumer_unit(consumer_unit_id);

CREATE TRIGGER trg_cycle_consumer_unit_updated
    BEFORE UPDATE ON public.cycle_consumer_unit
    FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

-- ---------------------------------------------------------------------
-- NOTIFICATIONS
-- Sistema de notificações interno para o backoffice.
-- ---------------------------------------------------------------------
CREATE TABLE public.notification (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID REFERENCES public.app_user(id) ON DELETE CASCADE,
    type          TEXT NOT NULL,
    title         TEXT NOT NULL,
    description   TEXT,
    entity_type   TEXT,
    entity_id     UUID,
    link          TEXT,
    severity      TEXT NOT NULL DEFAULT 'info'
                      CHECK (severity IN ('info','warning','error','success')),
    read_at       TIMESTAMPTZ,
    archived_at   TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notification_user_unread
    ON public.notification(user_id, created_at DESC)
    WHERE read_at IS NULL AND archived_at IS NULL;

CREATE INDEX idx_notification_entity
    ON public.notification(entity_type, entity_id, created_at DESC);

CREATE INDEX idx_notification_broadcast
    ON public.notification(created_at DESC)
    WHERE user_id IS NULL AND archived_at IS NULL;

-- ---------------------------------------------------------------------
-- PERFORMANCE INDICES on existing billing tables
-- ---------------------------------------------------------------------

-- billing_calculation
CREATE INDEX IF NOT EXISTS idx_billing_calculation_cycle
    ON public.billing_calculation(billing_cycle_id, status);
CREATE INDEX IF NOT EXISTS idx_billing_calculation_invoice_ref
    ON public.billing_calculation(utility_invoice_ref_id);

-- manual_adjustment
CREATE INDEX IF NOT EXISTS idx_manual_adjustment_calculation
    ON public.manual_adjustment(billing_calculation_id, created_at DESC);

-- generated_document
CREATE INDEX IF NOT EXISTS idx_generated_document_calculation
    ON public.generated_document(billing_calculation_id, type);

-- sync_job
CREATE INDEX IF NOT EXISTS idx_sync_job_status
    ON public.sync_job(status, created_at DESC)
    WHERE status IN ('pending','running','retrying');

-- utility_invoice_ref
CREATE INDEX IF NOT EXISTS idx_invoice_ref_uc_cycle
    ON public.utility_invoice_ref(consumer_unit_id, billing_cycle_id);
