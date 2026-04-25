# AGENTS.md — 5G Energia Fatura

> Este arquivo é destinado a agentes de codificação AI. Ele descreve a arquitetura, convenções e processos deste monorepo. Leia-o antes de fazer qualquer modificação.

---

## Visão Geral do Projeto

`5g-energia-fatura` é um monorepo de automação de faturas de energia (Coelba/Neoenergia) e backoffice de faturamento de energia compartilhada (usina fotovoltaica). O projeto extrai dados de faturas em PDF, sincroniza informações via portal web e API privada da Neoenergia, normaliza os dados e executa cálculos de cobrança com desconto contratual.

O repositório contém **quatro subsistemas principais**:

1. **`src/fatura/`** — Extrator documental Python (CLI + API HTTP). Usa Playwright para navegar no portal, PyMuPDF para parse de PDF, e Mistral OCR como fallback.
2. **`services/backend-go/`** — Backend principal em Go. Orquestra credenciais criptografadas, sessões, sincronização de UCs, persistência e expõe API pública.
3. **`services/doc-extractor-py/`** — Serviço Python isolado (FastAPI) de extração de PDF, consumido pelo backend-go via HTTP.
4. **`web/`** — Frontend Next.js 16 (App Router). Interface administrativa para clientes, UCs, faturas e credenciais de integração.

---

## Stack Tecnológico

| Camada | Tecnologia |
|--------|-----------|
| Extrator / Portal | Python 3.12+, Playwright, PyMuPDF, FastAPI, Uvicorn, Pydantic v2, SQLAlchemy 2.0, structlog |
| OCR (opcional) | Mistral OCR API (`mistralai>=2.4`) |
| Backend principal | Go 1.25, `net/http`, pgx/v5, `shopspring/decimal`, modernc.org/sqlite |
| Motor de cálculo | Go puro (`packages/calc-engine/`), sem I/O |
| Frontend | Next.js 16.1.6, React 19, TypeScript 5.9, Tailwind CSS v4, Prisma 7.7, PostgreSQL 17 |
| UI | shadcn/ui (59 componentes), Radix UI, Base UI, Hugeicons |
| Bancos de dados | SQLite (dados legados Python), PostgreSQL (backoffice + dados de integração) |
| Deploy | Docker multi-stage, GitHub Actions → GHCR → Portainer (Docker Swarm) |

---

## Estrutura do Monorepo

```
5g-energia-fatura/
├── src/fatura/                  # Python principal (extrator + orquestração)
│   ├── main.py                  # CLI entry point
│   ├── api.py                   # FastAPI app factory
│   ├── jobs.py                  # BatchProcessor (orquestração de download/parse)
│   ├── runtime.py               # Semáforo de jobs concorrentes
│   ├── coelba_client.py         # Playwright client do portal Neoenergia
│   ├── neoenergia_private_api.py# Cliente HTTP da API privada + bootstrap de sessão
│   ├── parser_pdf.py            # Parser PyMuPDF com fallback Mistral OCR
│   ├── mistral_ocr.py           # Cliente Mistral OCR
│   ├── models.py                # Modelos Pydantic de domínio
│   ├── service_models.py        # Modelos da API HTTP
│   ├── repository.py            # SQLAlchemy repository (clientes, contas, jobs)
│   ├── db/schema.py             # Declarações ORM (SQLAlchemy 2.0)
│   ├── config.py                # Loader de config.yaml
│   └── exceptions.py            # Hierarquia de exceções customizadas
│
├── services/
│   ├── backend-go/              # Backend principal Go
│   │   ├── cmd/api/main.go      # Entry point
│   │   ├── internal/app/        # HTTP handlers, mux, config
│   │   ├── internal/billing/    # Módulo de faturamento (contratos, cálculos)
│   │   ├── internal/extractor/  # Cliente HTTP para doc-extractor-py
│   │   ├── internal/neoenergia/ # Cliente HTTP da API privada Neoenergia
│   │   ├── internal/session/    # Gestão de credenciais criptografadas + token bootstrap
│   │   ├── internal/security/   # Criptografia (AES-GCM)
│   │   ├── internal/store/      # Abstração de persistência (SQLite + Postgres)
│   │   ├── internal/pgstore/    # Pool pgx/v5 para Postgres
│   │   ├── internal/sync/       # Serviço de sincronização de UC
│   │   └── migrations/          # golang-migrate (Postgres backoffice)
│   │
│   └── doc-extractor-py/        # Serviço isolado de extração PDF
│       └── app/main.py          # FastAPI que recebe PDF (path/base64) e retorna JSON
│
├── packages/
│   ├── calc-engine/             # Motor determinístico de cálculo (Go puro)
│   ├── normalizer/              # Normalizador de fatura → entrada do motor (Go)
│   └── contracts/               # JSON Schema da fronteira Python/Go
│
├── web/                         # Frontend Next.js (App Router)
│   ├── app/                     # Rotas (admin + API BFF)
│   ├── components/              # Componentes React (shadcn/ui + custom)
│   ├── lib/                     # go-client.ts, db.ts (Prisma), utils.ts
│   ├── prisma/                  # schema.prisma, migrations, seed.ts
│   └── contexts/                # Providers (QueryClient, Breadcrumb)
│
├── scripts/                     # Scripts Python exploratórios (não são produção)
├── tests/                       # Testes Python (pytest)
├── docs/                        # Documentação de APIs, OpenAPI, amostras
├── deploy/                      # portainer-stack.yml
└── data/                        # SQLite local (faturas.db)
```

---

## Arquivos de Configuração Principais

| Arquivo | Propósito |
|---------|-----------|
| `pyproject.toml` | Python principal: dependências, scripts `fatura`/`fatura-api`/`fatura-openapi`, configuração do ruff e pytest |
| `services/doc-extractor-py/pyproject.toml` | Dependências isoladas do serviço de extração |
| `services/backend-go/go.mod` | Módulo Go do backend (usa `replace` local para `packages/`) |
| `packages/calc-engine/go.mod` | Motor de cálculo puro (sem dependências de I/O) |
| `packages/normalizer/go.mod` | Normalizador (depende de `calc-engine` via `replace ../calc-engine`) |
| `web/package.json` | Next.js, React, Prisma, Tailwind, Vitest |
| `web/prisma/schema.prisma` | Schema PostgreSQL (Client, ConsumerUnit, IntegrationCredential, etc.) |
| `web/prisma.config.ts` | **Prisma 7.7**: URL do banco vem de `env('DATABASE_URL')`, não do `schema.prisma` |
| `config.yaml` | Configuração runtime do Python (portal, clientes, banco, parser) |
| `config.yaml.example` | Exemplo comentado de todas as opções |
| `deploy/portainer-stack.yml` | Stack Docker Swarm para produção (backend + doc-extractor) |
| `.github/workflows/publish-backend.yml` | CI: testa Go, builda imagens Docker, deploy no Portainer |

---

## Comandos de Build e Teste

### Python (extrator principal)

```bash
# Instalar em ambiente virtual (usa src/ como pacote)
pip install -e ".[dev,ocr]"

# CLI
python -m fatura.main --help
fatura --uc 007085489032 --mes-ano 04/2026

# API local
fatura-api                    # uvicorn em host:port definido em config.yaml
# ou
python -m fatura.api

# Exportar OpenAPI estático
fatura-openapi                # gera docs/openapi.json

# Testes
pytest                        # unit tests
pytest -m real_portal         # requer RUN_REAL_PORTAL_TESTS=1 + config.yaml válido
pytest -m real_api_e2e        # requer RUN_REAL_API_E2E=1
pytest -m real_mistral_e2e    # requer RUN_REAL_MISTRAL_E2E=1 + MISTRAL_API_KEY
```

### Python (doc-extractor isolado)

```bash
PYTHONPATH=src:services/doc-extractor-py \
  DOC_EXTRACTOR_HOST=127.0.0.1 DOC_EXTRACTOR_PORT=8090 \
  ./.venv/bin/python services/doc-extractor-py/app/main.py
```

### Go (backend)

```bash
cd services/backend-go

# Testes
go test ./...

# Rodar localmente (exemplo com todas as env vars)
BACKEND_HOST=127.0.0.1 \
BACKEND_PORT=8088 \
EXTRACTOR_BASE_URL=http://127.0.0.1:8090 \
BACKEND_INTEGRATION_PG_URL='postgres://backoffice:backoffice@127.0.0.1:5432/backoffice?sslmode=disable' \
BACKEND_ENCRYPTION_KEY='troque-esta-chave' \
BOOTSTRAP_PYTHON_BIN="$PWD/../../.venv/bin/python" \
BOOTSTRAP_SCRIPT_PATH="$PWD/../../scripts/bootstrap_neoenergia_token.py" \
go run ./cmd/api
```

### Go (packages)

```bash
cd packages/calc-engine && go test ./...
cd packages/normalizer && go test ./...
```

### Next.js (frontend)

```bash
cd web
pnpm install

# Dev (sobe Docker DB na porta 5433 + Next.js com Turbopack)
pnpm dev

# Type check
pnpm typecheck

# Testes unitários (vitest)
pnpm test
pnpm test:watch

# Banco de dados
pnpm db:migrate      # prisma migrate dev
pnpm db:seed         # 30 clientes, 15 UCs, 5 credenciais
pnpm db:studio       # Prisma Studio

# Lint / Format
pnpm lint
pnpm format
```

O PostgreSQL de desenvolvimento roda via `web/docker-compose.yml` na porta **5433**.

---

## Convenções de Código

### Python

- **Ruff** é o linter/formatador. Configurado em `pyproject.toml`:
  - `target-version = "py312"`
  - `line-length = 100`
- Use `structlog` para logging estruturado em todos os módulos.
- Use Pydantic v2 para validação de config, modelos de domínio e payloads de API.
- SQLAlchemy 2.0 com `DeclarativeBase` e style novo (`mapped_column`).
- Funções principais são `async` (Playwright, FastAPI). A CLI faz `asyncio.run(...)`.
- Exceções customizadas herdam de `FaturaError` (ver `exceptions.py`).
- Todo valor monetário é `Decimal` (nunca `float`). No banco SQLite, armazena-se como `String` para preservar precisão.

### Go

- Go 1.25. Não use `float64` para dinheiro — use `shopspring/decimal`.
- Handlers HTTP usam `http.ServeMux` padrão (sem framework externo).
- O backend-go usa `routeCatalog` customizado para gerar OpenAPI/Swagger UI automaticamente a partir dos handlers registrados.
- Toda rota de billing é condicional: só registra se `BackofficePGURL` estiver configurado.
- Erros do store usam `repo.ErrNotFound` (sentinel) para 404.

### TypeScript / Next.js

- **Sempre use componentes shadcn/ui** — nunca HTML cru. O projeto tem 59 componentes em `components/ui/`.
- Ícones: `HugeiconsIcon` do pacote `@hugeicons/react` com ícones de `@hugeicons/core-free-icons`. Não use `lucide-react` para UI nova.
- Forms usam `react-hook-form` + `zodResolver(Schema as any)` (cast intencional devido a incompatibilidade entre `@hookform/resolvers` v5.2.2 e Zod v4.3.6 — **não remova o `as any`**).
- `Select` como filtro standalone → usa popover estilizado (Radix). `NativeSelect` só em forms com `register()`.
- Forms que abrem de outra página/painel devem usar `Dialog` (modal), nunca `Sheet`.
- O frontend é **dual-backend**: rotas `/api/clients/**` e `/api/ucs/**` falam com PostgreSQL local via Prisma; rotas `/api/integration/**` falam com o backend-go via `goFetch()`.
- **`BACKEND_GO_URL` é server-only** (sem prefixo `NEXT_PUBLIC_`). O browser nunca deve chamá-lo diretamente.
- Prisma 7.7 quebra a convenção antiga: a URL do banco não fica em `schema.prisma`, mas sim em `prisma.config.ts` via `env('DATABASE_URL')`.

---

## Instruções de Teste

### Python

- **Testes unitários** rodam offline com fixtures PDF reais em `tests/fixtures/`.
- **Testes de integração real** são protegidos por markers pytest e env vars:
  - `RUN_REAL_PORTAL_TESTS=1` → `test_real_portal_integration.py`
  - `RUN_REAL_API_E2E=1` → `test_real_api_e2e.py`
  - `RUN_REAL_MISTRAL_E2E=1` + `MISTRAL_API_KEY` → `test_real_mistral_e2e.py`
- O `conftest.py` expõe o fixture `fixtures_dir` apontando para `tests/fixtures/`.

### Go

- `go test ./...` em `services/backend-go/` e em cada `packages/*/`.
- O `normalizer` usa strings reais de faturas nos testes (Paula, MP-BA, Azi Dourado).
- O `calc-engine` foi validado contra 5 competências reais da planilha operacional do cliente.

### Next.js

- `pnpm test` roda vitest. Hoje só existem testes unitários para `lib/go-client.ts` e `types/clientes.ts`.
- Camadas de BFF, componentes UI e páginas são verificadas por `pnpm typecheck`.

---

## Arquitetura de Runtime

### Fluxo de Sincronização de UC (end-to-end)

```
Frontend (Next.js)
  → POST /api/integration/ucs/[uc]/sync   (BFF)
    → goFetch → backend-go: POST /v1/consumer-units/{uc}/sync
      → session manager resolve token (ou bootstrap via Playwright)
        → neoenergia private API: lista faturas, histórico, imóvel, etc.
        → doc-extractor-py: POST /v1/extract (PDF → JSON)
      → persiste em PostgreSQL (sync_run, invoice, billing_record, document_record)
    → retorna sync_run_id
  → Frontend polling: GET /api/integration/sync-runs/{id}
```

### Segurança de Credenciais

- Credenciais de portal (CPF/CNPJ + senha) são **criptografadas em repouso** com AES-GCM (`internal/security/crypto.go`).
- O bootstrap de sessão executa um script Python (`scripts/bootstrap_neoenergia_token.py`) via subprocess para fazer login no portal e extrair o Bearer token do `localStorage`. O token resultante também é criptografado antes de ser persistido.
- A API key do backend-go é opcional e verificada pelo middleware `withOptionalAPIKeyAuth`.

---

## Deploy

### CI/CD (GitHub Actions)

O workflow `.github/workflows/publish-backend.yml`:

1. **Testa** o backend-go (`go test ./...`).
2. **Builda** duas imagens Docker multi-stage (`linux/arm64`) e publica no GHCR:
   - `ghcr.io/<owner>/5g-energia-fatura-backend-go:main`
   - `ghcr.io/<owner>/5g-energia-fatura-doc-extractor-py:main`
3. **Deploya** no Portainer (Docker Swarm) via API, usando `deploy/portainer-stack.yml`.

### Imagens Docker

- **backend-go**: base `mcr.microsoft.com/playwright/python:v1.55.0-noble` (precisa do Python + Playwright para bootstrap de token). O binário Go é copiado de um stage builder `golang:1.25-bookworm`.
- **doc-extractor-py**: base `python:3.12-slim`, instala o pacote Python principal com extra `[ocr]`.

### Portainer Stack

- Rede pública (`network_public`) para o backend (exposto via Traefik).
- Rede interna (`internal`) para comunicação backend ↔ doc-extractor.
- Healthchecks HTTP em `/healthz` nas duas imagens.

---

## Variáveis de Ambiente Importantes

### Python (extrator)

| Var | Descrição |
|-----|-----------|
| `MISTRAL_API_KEY` | Chave da API Mistral OCR (opcional) |
| `MISTRAL_MODEL` | Modelo OCR (`mistral-ocr-latest`) |

### Go (backend)

| Var | Padrão | Descrição |
|-----|--------|-----------|
| `BACKEND_HOST` | `127.0.0.1` | Host do servidor HTTP |
| `BACKEND_PORT` | `8080` | Porta do servidor |
| `BACKEND_API_KEY` | `""` | API key opcional para autenticação |
| `EXTRACTOR_BASE_URL` | `http://127.0.0.1:8090` | URL do doc-extractor-py |
| `NEOENERGIA_API_BASE_URL` | `https://apineprd.neoenergia.com` | Base da API privada |
| `BACKEND_INTEGRATION_PG_URL` | `""` | Postgres de integração (obrigatório para billing) |
| `BACKOFFICE_PG_URL` | `""` | Postgres do backoffice (habilita rotas `/v1/billing/*`) |
| `BACKEND_ENCRYPTION_KEY` | `""` | Chave AES-GCM para criptografia de credenciais |
| `BOOTSTRAP_PYTHON_BIN` | `./.venv/bin/python` | Python para script de bootstrap |
| `BOOTSTRAP_SCRIPT_PATH` | `scripts/bootstrap_neoenergia_token.py` | Script de bootstrap |
| `ARTIFACTS_DIR` | `./artifacts` | Diretório de artefatos |

### Next.js (frontend)

| Var | Descrição |
|-----|-----------|
| `DATABASE_URL` | PostgreSQL local (ex: `postgresql://user:password@localhost:5433/fatura_dev`) |
| `BACKEND_GO_URL` | URL do backend-go (**server-only**, nunca exponha no browser) |

---

## Considerações de Segurança

- **Nunca comite** `config.yaml` com senhas reais (ele já está no `.gitignore`, mas verifique).
- **Não exponha** `BACKEND_GO_URL` ou `BACKEND_ENCRYPTION_KEY` no frontend.
- Tokens e senhas de portal são criptografados com AES-GCM antes de ir para o banco.
- O bootstrap de sessão executa Playwright com credenciais em plaintext apenas na memória do processo Python; o token resultante é imediatamente criptografado.
- O middleware de API key do backend-go é opcional (`X-API-Key` header). Em produção, configure `BACKEND_API_KEY`.
- Playwright roda com flags anti-detecção (`--disable-blink-features=AutomationControlled`, remoção de `navigator.webdriver`), mas o portal pode exibir CAPTCHA. O código lida com retry e evidências (screenshot + HTML) em `downloads/_errors/`.

---

## Convenções de Desenvolvimento Específicas

### Breadcrumb Automático (Next.js)

O layout `(admin)` provê `BreadcrumbProvider`. Segmentos estáticos são auto-rotulados (`clientes` → "Clientes"). Para segmentos dinâmicos, cada página registra:

```tsx
import { useSetBreadcrumbTitle } from '@/contexts/breadcrumb'
useSetBreadcrumbTitle(id, client?.nome_razao) // chamada incondicional
```

### Dual Backend no Frontend

```
Browser → Next.js App Router
  ├── /api/clients/**, /api/ucs/**         → Prisma → PostgreSQL local
  └── /api/integration/**                  → goFetch → backend-go
```

### Normalização de Decimais Brasileiros

O parser Python usa `normalizar_decimal_br()` para converter `"1.234,56"` → `Decimal('1234.56')`. O normalizer Go tem `parseBRDecimal()` com a mesma responsabilidade. **Nunca use `float`/`float64` para dinheiro.**

### Versionamento de Schema (Python ↔ Go)

A comunicação entre `doc-extractor-py` e `backend-go` usa JSON Schema versionado em `packages/contracts/`:
- `extractor-request.schema.json` (v1.0.0)
- `extractor-response.schema.json` (v1.0.0)
- `billing-record.schema.json`

---

## Documentação Complementar

- `docs/api.md` — API legada Python (FastAPI em `src/fatura/api.py`)
- `docs/backend-go-api.md` — API do backend-go (endpoints, fluxos, exemplos de env vars)
- `docs/openapi.json` / `docs/openapi-go.json` — Schemas OpenAPI exportados
- `services/backend-go/internal/app/BILLING_INTEGRATION.md` — Instruções de integração do módulo billing no server.go
- `web/CLAUDE.md` — Convenções específicas do frontend (em inglês, mas o projeto fala português)

---

## Notas para Agentes

- **Idioma principal**: português (código, comentários, documentação). Variáveis e funções usam português brasileiro (`fatura`, `cliente`, `conta`, `vencimento`, `desconto`).
- **Não altere** `config.yaml` no repositório (contém dados reais de cliente). Use `config.yaml.example` como template.
- **Não remova** o `as any` em `zodResolver(Schema as any)` no frontend — é um workaround conhecido.
- **Não use** `float64` para dinheiro em Go. Use `shopspring/decimal`.
- Ao adicionar novas rotas no backend-go, sempre registre no `routeCatalog` com `docs.add(...)` para que apareçam em `/openapi.json` e `/docs`.
