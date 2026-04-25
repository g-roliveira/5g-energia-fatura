# AGENTS.md — 5G Energia Fatura (services)

> Arquivo de referência para agentes de codificação. Escrito em português porque todo o código, comentários e documentação do projeto usam português como idioma principal.

---

## Visão geral do projeto

Este repositório contém dois serviços backend do sistema **5G Energia Fatura**:

- **`backend-go/`** — Serviço principal escrito em Go. Responsável por:
  - Cliente HTTP da API privada Neoenergia (Coelba);
  - Sincronização de UCs (unidades consumidoras) e faturas;
  - Gestão de credenciais e sessões criptografadas;
  - API pública para frontend e integração;
  - Orquestração do serviço `doc-extractor-py`;
  - Domínio de faturamento (billing) com contratos, ciclos e cálculos.

- **`doc-extractor-py/`** — Serviço documental em Python/FastAPI. Responsável por:
  - Parse local de PDFs de fatura com `PyMuPDF`;
  - Fallback com `Mistral OCR` (quando `MISTRAL_API_KEY` está configurada);
  - Geração de `source_map` e `confidence_map` para cada campo extraído;
  - Resposta compatível com os contratos JSON Schema do extrator.

---

## Stack tecnológico

| Camada | Tecnologia |
|--------|-----------|
| Backend principal | Go 1.25 |
| Extrator documental | Python 3.12 + FastAPI + Uvicorn |
| Banco de dados integração | SQLite (padrão) ou Postgres (via `INTEGRATION_PG_URL`) |
| Banco de dados backoffice | Postgres 15+ (schemas `core` e `billing`) |
| Driver Postgres (Go) | `pgx/v5` (pool) e `pgx/stdlib` (sql.DB) |
| Driver SQLite (Go) | `modernc.org/sqlite` (CGO-free) |
| Criptografia | AES-256-GCM com chave derivada de `BACKEND_ENCRYPTION_KEY` |
| Documentação API | OpenAPI 3.0.3 gerado manualmente em código + Swagger UI em `/docs` |
| Deploy | Docker + Docker Compose (imagens multistage) |

---

## Estrutura de diretórios

```
services/
├── backend-go/
│   ├── cmd/api/main.go              # Entry point
│   ├── internal/
│   │   ├── app/                     # HTTP server, handlers, roteamento, config
│   │   │   ├── server.go            # Registro de todas as rotas e wiring
│   │   │   ├── server_billing.go    # Rotas de faturamento (ativo só com Postgres)
│   │   │   ├── server_test.go       # Testes de integração HTTP
│   │   │   ├── config.go            # LoadConfigFromEnv
│   │   │   ├── route_catalog.go     # OpenAPI spec em código
│   │   │   ├── auth_middleware.go   # API key opcional (header X-API-Key)
│   │   │   ├── discover.go          # Handler de discovery de perfil/UCs
│   │   │   ├── docs_handler.go      # Swagger UI, /openapi.json, /docs.md
│   │   │   ├── http_helpers.go      # writeJSON, logging, sanitização de erro
│   │   │   ├── date_helpers.go      # parseISO YYYY-MM-DD
│   │   │   └── BILLING_INTEGRATION.md  # Guia de integração do módulo billing
│   │   ├── billing/
│   │   │   ├── contract/service.go  # Regras de negócio de contrato (versionado)
│   │   │   └── repo/                # Repositories Postgres (contract, calculation)
│   │   ├── extractor/client.go      # Cliente HTTP do doc-extractor-py
│   │   ├── neoenergia/client.go     # Cliente HTTP da API privada Neoenergia
│   │   ├── pgstore/                 # Conexão pgx pool para backoffice Postgres
│   │   ├── security/crypto.go       # Cipher AES-256-GCM
│   │   ├── session/                 # Gestão de credenciais e sessões
│   │   │   ├── manager.go           # CRUD criptografado de credenciais/sessões
│   │   │   └── bootstrap.go         # Runner Python para obter bearer token
│   │   ├── store/                   # Persistência do domínio de integração
│   │   │   ├── integration_store.go # Interface comum (SQLite/Postgres)
│   │   │   ├── sqlite.go            # Implementação SQLite com auto-migrate
│   │   │   └── postgres.go          # Implementação Postgres com auto-migrate
│   │   └── sync/                    # Orquestração de sync UC
│   │       ├── service.go           # Lógica de chamadas paralelas à Neoenergia
│   │       ├── billing_record.go    # Montagem do billing record
│   │       └── document_record.go   # Montagem do document record
│   ├── migrations/                  # SQL migrations para Postgres
│   │   ├── 000001_initial_schema.{up,down}.sql      # core + billing
│   │   └── 000002_backend_unification_postgres.{up,down}.sql  # integração + jobs
│   ├── go.mod / go.sum
│   └── Dockerfile                   # Multistage: Go builder + base Python/Playwright
├── doc-extractor-py/
│   ├── app/main.py                  # FastAPI app + endpoint /v1/extract
│   ├── pyproject.toml               # Dependências: fastapi, uvicorn, pydantic
│   └── Dockerfile                   # Base python:3.12-slim
```

---

## Bancos de dados

O projeto usa **dois bancos de dados** com responsabilidades bem separadas:

### 1. Banco de integração (backend-go)

- **Padrão:** SQLite local (`file:data/backend-go.db`).
- **Produção:** Postgres (mesmo host do backoffice, schema `public`).
- **Tabelas:** `credentials`, `sessions`, `consumer_units`, `sync_runs`, `invoices`, `invoice_api_snapshots`, `invoice_documents`, `invoice_items`, `invoice_extraction_results`, `invoice_field_sources`.
- **Migração:** auto-migrate no código (`sqlite.go` / `postgres.go`). Não usa ferramenta externa para esse schema.

### 2. Banco de backoffice (Postgres)

- **Schemas:** `core` (cadastro) e `billing` (faturamento).
- **Migração:** arquivos SQL versionados em `backend-go/migrations/`. Use `golang-migrate` ou execute manualmente.
- **Tabelas principais:**
  - `core.app_user`, `core.customer`, `core.address`, `core.consumer_unit`, `core.credential_link`
  - `billing.contract`, `billing.billing_cycle`, `billing.utility_invoice_ref`, `billing.utility_invoice_item`, `billing.utility_invoice_scee`, `billing.billing_calculation`, `billing.manual_adjustment`, `billing.generated_document`, `billing.sync_job`, `billing.audit_log`

**Convenção importante:** o domínio de integração (SQLite/Go) e o domínio de negócio (Postgres) não compartilham transações. Quando o billing precisa de dados da integração, lê primeiro e depois inicia a transação Postgres.

---

## Variáveis de ambiente

### backend-go

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| `BACKEND_HOST` | `127.0.0.1` | Host do servidor HTTP |
| `BACKEND_PORT` | `8080` | Porta do servidor HTTP |
| `BACKEND_API_KEY` | *(vazio)* | Se configurado, exige `X-API-Key` nos endpoints `/v1/*` |
| `BACKEND_DATABASE_URL` | `file:data/backend-go.db` | DSN do banco de integração (SQLite ou Postgres) |
| `BACKEND_INTEGRATION_PG_URL` | *(vazio)* | Alias legacy para Postgres de integração |
| `BACKOFFICE_PG_URL` | *(vazio)* | DSN do Postgres backoffice. Se vazio, módulo billing fica inativo |
| `EXTRACTOR_BASE_URL` | `http://127.0.0.1:8090` | URL base do `doc-extractor-py` |
| `NEOENERGIA_API_BASE_URL` | `https://apineprd.neoenergia.com` | URL base da API Neoenergia |
| `ARTIFACTS_DIR` | `./artifacts` | Diretório de artefatos |
| `BACKEND_ENCRYPTION_KEY` | *(vazio)* | Chave para AES-256-GCM. **Obrigatória em produção** |
| `BOOTSTRAP_PYTHON_BIN` | `./.venv/bin/python` | Python usado pelo bootstrap de token |
| `BOOTSTRAP_SCRIPT_PATH` | `scripts/bootstrap_neoenergia_token.py` | Script de login Neoenergia |

### doc-extractor-py

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| `DOC_EXTRACTOR_HOST` | `127.0.0.1` | Host do servidor FastAPI |
| `DOC_EXTRACTOR_PORT` | `8090` | Porta do servidor FastAPI |
| `MISTRAL_API_KEY` | *(vazio)* | Ativa fallback Mistral OCR quando presente |
| `MISTRAL_MODEL` | `mistral-ocr-latest` | Modelo Mistral OCR |

---

## Comandos de build e teste

### backend-go

```bash
# Build
cd backend-go
go build -o backend-go ./cmd/api

# Run (dev)
go run ./cmd/api

# Testes (incluem testes de integração com httptest)
go test ./...
go test -v ./internal/app/...

# Migrations Postgres (requer golang-migrate)
migrate -path backend-go/migrations -database "postgres://..." up
```

### doc-extractor-py

```bash
cd doc-extractor-py
pip install -e ".[ocr]"
python -m app.main
# ou
uvicorn app.main:create_app --factory --host 127.0.0.1 --port 8090
```

### Docker

```bash
# Build das imagens (contexto deve ser a raiz do monorepo, não services/)
docker build -f services/backend-go/Dockerfile -t backend-go .
docker build -f services/doc-extractor-py/Dockerfile -t doc-extractor-py .
```

**Atenção:** o `Dockerfile` de `backend-go` é multistage e espera o contexto de build na raiz do monorepo (acima de `services/`), porque copia `pyproject.toml`, `src/` e `scripts/` além do próprio serviço Go.

---

## Convenções de código

### Go

- **Idioma:** português para comentários, mensagens de erro, nomes de campos JSON e documentação.
- **Logger:** `log/slog` com handler JSON. Usar chaves snake_case: `slog.Info("evento", "chave", valor)`.
- **Erros:** mensagens de erro sanitizadas em `http_helpers.go#sanitizeError` — nunca logam senhas, tokens, documentos ou autorizações.
- **IDs:** hex strings de 32 caracteres (16 bytes random) para SQLite; `uuid.UUID` (github.com/google/uuid) para Postgres billing.
- **Datas:** ISO 8601 date-only (`YYYY-MM-DD`) em requests/responses de billing; RFC3339 para timestamps de integração.
- **Decimais:** sempre `github.com/shopspring/decimal` para campos monetários no billing. Nunca float.
- **SQL:** queries em maiúsculas/minúsculas misturadas (estilo legível), com placeholders `$N` para pgx e `?` convertido via `pgQuery()` para Postgres compatível com `database/sql`.
- **Roteamento:** `http.ServeMux` padrão, sem frameworks. Rotas documentadas manualmente em `routeCatalog` para gerar OpenAPI.

### Python

- FastAPI com Pydantic v2.
- Type hints obrigatórios (`from __future__ import annotations`).
- Nomes de variáveis e mensagens em português.

---

## Testes

### Estratégia

- **Testes de integração HTTP** em `backend-go/internal/app/server_test.go` usando `httptest.Server`.
- Mock de serviços externos via `httptest.NewServer` (Neoenergia e doc-extractor).
- Banco de testes: SQLite in-memory (`file::memory:?cache=shared`).
- Testes que precisam de Python usam `exec.LookPath("python3")` e fazem `t.Skip` se não encontram.

### Testes existentes (server_test.go)

- `TestDocsEndpoints` — verifica `/docs`, `/openapi.json`, `/docs.md`
- `TestAPIKeyProtectionOnV1Endpoints` — valida proteção por `X-API-Key`
- `TestSyncUCEndpoint` — sync completo com mock Neoenergia
- `TestCredentialAndSessionEndpointsDoNotExposeSecrets` — garante que responses não vazam documento/senha/token
- `TestCredentialActionRoutingStatusCodes` — 404/405 em rotas de credencial
- `TestDiscoverEndpointSanitizesPartialUpstreamErrors` — descoberta com falhas parciais sanitizadas
- `TestSyncUCEndpointWithExtraction` — sync + extração PDF + persistência completa

### Como rodar

```bash
cd backend-go
go test ./internal/app/ -v -count=1
```

---

## Segurança

### Criptografia em repouso

- Credenciais (`documento`, `senha`) e tokens de sessão (`bearer_token`) são criptografados com **AES-256-GCM** antes de persistir.
- Cada valor tem seu próprio nonce. O nonce é armazenado em coluna separada (`*_nonce`).
- A chave é derivada via SHA-256 de `BACKEND_ENCRYPTION_KEY`.

### API Key

- Se `BACKEND_API_KEY` estiver configurado, todos os endpoints `/v1/*` exigem header `X-API-Key`.
- Endpoints de infra (`/healthz`, `/docs`, `/openapi.json`, `/docs.md`) permanecem públicos.

### Sanitização

- `sanitizeError` remove de logs qualquer erro que contenha palavras sensíveis (`senha`, `token`, `bearer`, `authorization`, `documento`).
- Responses de credenciais e sessões retornam o documento **mascarado** (apenas últimos 4 dígitos visíveis).
- Responses de discovery sanitizam erros upstream para não vazar tokens.

---

## Deploy e arquitetura de runtime

### Serviços

```
┌─────────────────┐      HTTP      ┌─────────────────────┐
│   backend-go    │ ◄────────────► │  doc-extractor-py   │
│    :8080        │   /v1/extract  │      :8090          │
└────────┬────────┘                └─────────────────────┘
         │
         │ HTTP (privada)
         ▼
┌─────────────────┐
│  Neoenergia API │
│ (apineprd...)   │
└─────────────────┘
```

### Bancos

- **Integração:** SQLite local (dev) ou Postgres (produção).
- **Backoffice:** Postgres obrigatório para o módulo billing. Se `BACKOFFICE_PG_URL` for vazio, as rotas `/v1/billing/*` não são registradas.

### Docker

- `backend-go` usa imagem base `mcr.microsoft.com/playwright/python:v1.55.0-noble` porque precisa de ambiente Python para rodar o script de bootstrap de token e potencialmente orquestrar Playwright no futuro.
- `doc-extractor-py` usa `python:3.12-slim`.

---

## Regras de negócio principais

### Contratos (billing)

- Um contrato **nunca é atualizado in-place**. Alterar qualquer termo cria uma nova versão.
- A versão anterior tem `vigencia_fim` ajustada para o dia anterior ao início da nova.
- Invariante: **exatamente um contrato ativo por UC** (garantido por partial unique index no Postgres).
- Cálculos de faturamento são imutáveis: quando o contrato muda, cálculos antigos ficam como `superseded` e novos são inseridos com `version` incrementada.

### Sincronização de UC

1. Resolve credencial → sessão → bearer token (ou usa token manual).
2. Chama múltiplos endpoints Neoenergia em paralelo (grupo cliente, minha conta, UCs, imóvel, protocolo, etc.).
3. Se `include_pdf=true`, baixa PDF da primeira fatura e pode chamar o extrator (`include_extraction=true`).
4. Monta `BillingRecord` e `DocumentRecord`, avalia `completeness`.
5. Persiste tudo em transação atômica (sync_run + consumer_unit + invoice + filhos).

---

## Dicas para agentes

- **Nunca assuma que o billing está ativo:** sempre verifique `cfg.BackofficePGURL != ""` antes de usar pools/pgstore.
- **Nunca use float para dinheiro:** no billing, use `decimal.Decimal`.
- **Sempre adicione rotas ao `routeCatalog`:** assim elas aparecem em `/docs`, `/openapi.json` e `/docs.md`.
- **Sempre sanitize erros:** use `writeInternalError` e `sanitizeError` para evitar vazamento de credenciais em logs.
- **Testes de HTTP:** prefira `httptest.NewServer` + `http.Get`/`http.Post` ao invés de chamar handlers diretamente — o projeto testa assim.
- **Migrations SQL:** quando alterar schemas Postgres, adicione novos arquivos numerados em `migrations/` e mantenha os `.down.sql` simétricos.
