-- Migration 004 DOWN
DROP TABLE IF EXISTS public.cycle_consumer_unit CASCADE;
DROP TABLE IF EXISTS public.notification CASCADE;
DROP INDEX IF EXISTS idx_billing_calculation_cycle;
DROP INDEX IF EXISTS idx_billing_calculation_invoice_ref;
DROP INDEX IF EXISTS idx_manual_adjustment_calculation;
DROP INDEX IF EXISTS idx_generated_document_calculation;
DROP INDEX IF EXISTS idx_sync_job_status;
DROP INDEX IF EXISTS idx_invoice_ref_uc_cycle;
