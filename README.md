# Módulo Billing — PR 1/N

Este commit adiciona ao seu monorepo `5g-energia-fatura` toda a base do
módulo de **faturamento de energia compartilhada** (o "backoffice" da
Azi Dourado), sem tocar em uma linha do código que você já tem.

## O que entra

### `packages/` — Go modules irmãos, reutilizáveis

- `packages/calc-engine/` — motor determinístico de cálculo. Puro,
  sem I/O. Validado em turnos anteriores contra 5 competências reais
  da planilha do cliente. 4 testes passam (happy path, custo
  disponibilidade, bandeira com desconto, validação de contrato).

- `packages/normalizer/` — transforma `sync.BillingRecord` do
  backend-go em entrada do motor. Classifier com 3 gerações de rótulos
  (bandeiras legadas, novas, posto), extrator SCEE/MMGD dos 3 layouts
  documentais. 4 testes passam (classify real strings, scee real texts,
  normalize end-to-end Paula, ignored items).

### `services/backend-go/migrations/` — golang-migrate

- `000001_initial_schema.up.sql` — Postgres do backoffice com dois
  schemas: `core` (cadastro: customer, address, consumer_unit,
  credential_link, app_user) e `billing` (faturamento: contract,
  billing_cycle, utility_invoice_ref, utility_invoice_item,
  utility_invoice_scee, billing_calculation, manual_adjustment,
  generated_document, sync_job, audit_log). Triggers de updated_at,
  índices de produção (partial unique de contrato vigente, GIN de
  snapshot JSONB, hot path do worker).

- `000001_initial_schema.down.sql` — reverte tudo.

### `services/backend-go/internal/pgstore/`

- Pool pgx/v5 configurável via env vars (`BACKOFFICE_PG_*`). Ping na
  inicialização. Coexiste com o SQLite existente.

### `services/backend-go/internal/billing/`

- `repo/` — camada de dados (todas as queries SQL). Dois repos:
  `ContractRepo` e `CalculationRepo`. `types.go` com structs Go
  espelhando as tabelas.
- `contract/service.go` — serviço com a regra "novo contrato fecha
  anterior automaticamente", em uma única transação.

### `services/backend-go/internal/app/`

- `server_billing.go` — handlers HTTP novos + `RegisterBillingRoutes`
  que plugga no mux e catálogo do OpenAPI **já existentes**.
- `date_helpers.go` — helper mínimo de parse de data.
- `BILLING_INTEGRATION.md` — **LEIA ESTE** — instruções exatas das
  3 mudanças no `server.go` existente pra ativar o módulo.

## O que não entra (escopo do PR seguinte)

- `billing/cycle/` — orquestração de ciclo (abrir/fechar competência,
  disparar syncs em massa, agregar resultados)
- `billing/adjustment/` — aplicar ajuste manual com nova version
- `billing/pdf/` — geração do PDF do cliente (chromium headless)
- `internal/worker/` — pool de goroutines consumindo `sync_job`
- SSE `/v1/events/billing-cycles/{id}`
- Testes de integração ponta-a-ponta contra Postgres
- OpenAPI com schemas completos (hoje só método+path+summary)

Tudo isso está com interfaces e tabelas já preparadas — é continuação
natural. A próxima rodada não vai precisar refatorar nada do que está
aqui.

## Como aplicar

O bundle é um repo git separado. Você tem 2 opções:

### Opção A — Aplicar como branch nova no seu repo

```bash
# No seu repo 5g-energia-fatura:
cd /caminho/para/5g-energia-fatura
git fetch /caminho/para/backoffice-billing.bundle main:feat/backoffice-billing
git checkout feat/backoffice-billing
# Revisa, push, abre PR no GitHub
git push origin feat/backoffice-billing
```

### Opção B — Aplicar os arquivos diretamente

```bash
# Extrai o tar.gz
tar -xzf backoffice-billing.tar.gz -C /tmp/patch

# Copia os arquivos pro seu repo (reproduz a estrutura de pastas)
cp -r /tmp/patch/packages/calc-engine YOUR_REPO/packages/
cp -r /tmp/patch/packages/normalizer YOUR_REPO/packages/
cp -r /tmp/patch/services/backend-go/migrations YOUR_REPO/services/backend-go/
cp -r /tmp/patch/services/backend-go/internal/pgstore YOUR_REPO/services/backend-go/internal/
cp -r /tmp/patch/services/backend-go/internal/billing YOUR_REPO/services/backend-go/internal/
cp /tmp/patch/services/backend-go/internal/app/server_billing.go YOUR_REPO/services/backend-go/internal/app/
cp /tmp/patch/services/backend-go/internal/app/date_helpers.go YOUR_REPO/services/backend-go/internal/app/
cp /tmp/patch/services/backend-go/internal/app/BILLING_INTEGRATION.md YOUR_REPO/services/backend-go/internal/app/

# Então segue as 3 mudanças do BILLING_INTEGRATION.md
cd YOUR_REPO && cat services/backend-go/internal/app/BILLING_INTEGRATION.md
```

## Próxima rodada

Me diz "continua" e eu faço o PR 2/N com:
1. `billing/cycle/` completo
2. `billing/adjustment/`
3. `billing/pdf/` + template HTML
4. `internal/worker/` com os 4 handlers
5. Endpoints HTTP do cycle, adjustment, pdf, bulk actions
6. SSE de progresso
7. Testes de integração

Se preferir priorizar outra coisa, me diz qual.
