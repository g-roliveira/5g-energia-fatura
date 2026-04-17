# Integrations

## Backend Go (concessionária integration)

**Produção**: `https://api5g.numbro.app`
**Docs (Swagger UI)**: `https://api5g.numbro.app/docs`
**OpenAPI spec**: `https://api5g.numbro.app/openapi.json`
**Base URL runtime**: `BACKEND_GO_URL` env var (server-only, nunca `NEXT_PUBLIC_`)
**Auth**: Sem auth header documentado — isolamento por URL e rede (403 para requests externos não autorizados)

### Endpoint Catalog

| Method | Path | Purpose |
|--------|------|---------|
| GET | /healthz | Health check |
| POST | /v1/credentials | Create encrypted credential |
| POST | /v1/credentials/{id}/session | Create/renew Neoenergia session |
| POST | /v1/consumer-units/{uc}/sync | Full ETL sync for a UC |
| GET | /v1/consumer-units | List UCs (params: limit, status) |
| GET | /v1/consumer-units/{uc} | UC detail with latest invoice + sync |
| GET | /v1/consumer-units/{uc}/invoices | UC invoice list (params: limit, status) |
| GET | /v1/consumer-units/{uc}/latest-invoice | Latest invoice for UC |
| GET | /v1/invoices/{id} | Invoice detail (billing_record + document_record + items) |
| GET | /v1/sync-runs/{id} | Sync audit trail |
| GET | /openapi.json | OpenAPI spec |
| GET | /docs | Swagger UI |

### Key Request Shapes

**POST /v1/credentials**
```json
{ "label": "neo-paula", "documento": "03021937586", "senha": "portal-pass", "uf": "BA", "tipo_acesso": "normal" }
```
Response: `{ id, label, documento (masked), uf, tipo_acesso, created_at }`

**POST /v1/consumer-units/{uc}/sync**
```json
{ "credential_id": "f9ba880b...", "include_pdf": true, "include_extraction": true }
```
Response: `{ uc, billing_record, document_record, persistence: { sync_run_id, invoice_id, status } }`

**GET /v1/consumer-units/{uc}**
Response: `{ uc, status, nome_cliente, latest_invoice, latest_sync_run }`

**GET /v1/invoices/{id}**
Response: `{ billing_record, document_record, items[] }`

**GET /v1/sync-runs/{id}**
Response: `{ status, error_message, raw_response, created_at }`

### Security Notes
- `senha` NEVER stored locally — forwarded once to Go backend, encrypted there
- Go backend returns masked `documento` (e.g. "*******7586")
- BFF must never expose `BACKEND_GO_URL` to client bundle
- Sync is async: POST /sync returns initial result, poll GET /sync-runs/{id} for final status

## Local PostgreSQL

Managed by Next.js backend (server-side only).
Connection via `DATABASE_URL` env var.
ORM: Prisma (see design.md).
