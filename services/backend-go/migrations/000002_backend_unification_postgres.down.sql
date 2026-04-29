-- =====================================================================
-- Migration 002 DOWN — Reverte unificação backend em Postgres
-- =====================================================================

DROP TABLE IF EXISTS job_items;
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS processamento_log;
DROP TABLE IF EXISTS itens_fatura;
DROP TABLE IF EXISTS contas;
DROP TABLE IF EXISTS clientes;

DROP TABLE IF EXISTS invoice_field_sources;
DROP TABLE IF EXISTS invoice_extraction_results;
DROP TABLE IF EXISTS invoice_items;
DROP TABLE IF EXISTS invoice_documents;
DROP TABLE IF EXISTS invoice_api_snapshots;
DROP TABLE IF EXISTS invoices;
DROP TABLE IF EXISTS sync_runs;
DROP TABLE IF EXISTS consumer_units;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS credentials;
