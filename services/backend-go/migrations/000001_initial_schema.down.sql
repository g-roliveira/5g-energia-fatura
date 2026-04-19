-- =====================================================================
-- Migration 001 DOWN — reverte schema completo
-- =====================================================================

DROP TRIGGER IF EXISTS trg_invoice_ref_updated ON billing.utility_invoice_ref;
DROP TRIGGER IF EXISTS trg_contract_updated ON billing.contract;
DROP TRIGGER IF EXISTS trg_credential_link_updated ON core.credential_link;
DROP TRIGGER IF EXISTS trg_consumer_unit_updated ON core.consumer_unit;
DROP TRIGGER IF EXISTS trg_address_updated ON core.address;
DROP TRIGGER IF EXISTS trg_customer_updated ON core.customer;
DROP TRIGGER IF EXISTS trg_app_user_updated ON core.app_user;

DROP FUNCTION IF EXISTS core.set_updated_at();

DROP TABLE IF EXISTS billing.audit_log;
DROP TABLE IF EXISTS billing.sync_job;
DROP TABLE IF EXISTS billing.generated_document;
DROP TABLE IF EXISTS billing.manual_adjustment;
DROP TABLE IF EXISTS billing.billing_calculation;
DROP TABLE IF EXISTS billing.utility_invoice_scee;
DROP TABLE IF EXISTS billing.utility_invoice_item;
DROP TABLE IF EXISTS billing.utility_invoice_ref;
DROP TABLE IF EXISTS billing.billing_cycle;
DROP TABLE IF EXISTS billing.contract;

DROP TABLE IF EXISTS core.credential_link;
DROP TABLE IF EXISTS core.consumer_unit;
DROP TABLE IF EXISTS core.address;
DROP TABLE IF EXISTS core.customer;
DROP TABLE IF EXISTS core.app_user;

DROP SCHEMA IF EXISTS billing;
DROP SCHEMA IF EXISTS core;
