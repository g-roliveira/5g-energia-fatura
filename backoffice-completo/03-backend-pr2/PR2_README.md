# PR 2/N — Billing Cycles, Adjustments, and SSE

Este é o **segundo Pull Request** do módulo de backoffice de faturamento.

## O que está incluído

### 1. Migration 000002
- Tabela `core.notification` para sistema de notificações
- Índices de performance em `billing.*` tables

### 2. Módulo `billing/cycle/`
- **service.go** — orquestração de ciclos (create, list, close)
- **rows.go** — consulta das linhas do ciclo (tabela principal do dashboard)
- **handlers.go** — endpoints HTTP
- **sse.go** — Server-Sent Events para realtime

### 3. Módulo `billing/adjustment/`
- **service.go** — ajustes manuais em cálculos
- Cria nova versão ao ajustar (imutabilidade)
- Registra histórico de ajustes

### 4. Documentação
- `BILLING_INTEGRATION_PR2.md` — guia de integração completo

## Novos endpoints

| Method | Path | Descrição |
|--------|------|-----------|
| POST | `/v1/billing/cycles` | Cria ciclo (competência mensal) |
| GET | `/v1/billing/cycles` | Lista ciclos |
| GET | `/v1/billing/cycles/{id}` | Detalhe do ciclo |
| GET | `/v1/billing/cycles/{id}/rows` | Tabela principal (status de cada UC) |
| POST | `/v1/billing/cycles/{id}/close` | Fecha ciclo |
| **GET** | `/v1/billing/events/cycles/{id}` | **SSE realtime** |
| POST | `/v1/billing/calculations/{id}/adjust` | Aplica ajuste manual |
| GET | `/v1/billing/calculations/{id}/adjustments` | Histórico de ajustes |

## Como aplicar

### Pré-requisitos
- PR 1/N já aplicado
- `BACKOFFICE_PG_URL` configurado
- golang-migrate instalado

### Comandos

```bash
# 1. Extrair o bundle
git bundle unbundle backoffice-billing-pr2.bundle --ref-excludes=refs/heads/main

# 2. Merge na branch desejada
git merge pr2-billing-cycles

# 3. Rodar migration
cd services/backend-go
migrate -path migrations -database "$BACKOFFICE_PG_URL" up

# 4. Integrar no server.go
# Siga BILLING_INTEGRATION_PR2.md
```

## Estrutura de commits

1. `feat(billing): add migration 000002 - notification table and indices`
2. `feat(billing): add cycle orchestration module`
3. `feat(billing): add manual adjustment module`
4. `docs(billing): add PR2 integration guide`

## Testar

```bash
# Criar ciclo
curl -X POST http://localhost:8080/v1/billing/cycles \
  -H "Content-Type: application/json" \
  -d '{"year": 2026, "month": 4, "include_all_active": true, "created_by": "user-id"}'

# SSE
curl -N http://localhost:8080/v1/billing/events/cycles/{id}

# Fechar ciclo
curl -X POST http://localhost:8080/v1/billing/cycles/{id}/close \
  -d '{"closed_by": "user-id"}'
```

## Limitações do MVP

- SSE usa polling (5s) em vez de LISTEN/NOTIFY
- Ajuste manual não reroda motor de cálculo (apenas cria nova version)
- Sem PDF generation (PR 3)
- Sem bulk actions (PR 3)
- Sem worker pool (PR 3)

## Próximo PR

PR 3/N incluirá:
- PDF generation com chromium
- Worker pool para jobs assíncronos
- Bulk actions (sync_all, recalculate_all, generate_pdfs)
- LISTEN/NOTIFY para SSE realtime
- Testes de integração com testcontainers
- OpenAPI completo

## Arquivos modificados

```
services/backend-go/
├── migrations/
│   ├── 000002_notification.up.sql        [NEW]
│   └── 000002_notification.down.sql      [NEW]
└── internal/
    ├── app/
    │   └── BILLING_INTEGRATION_PR2.md    [NEW]
    └── billing/
        ├── cycle/
        │   ├── service.go                [NEW]
        │   ├── rows.go                   [NEW]
        │   ├── handlers.go               [NEW]
        │   └── sse.go                    [NEW]
        └── adjustment/
            └── service.go                [NEW]
```

## Contato

Dúvidas: ver BILLING_INTEGRATION_PR2.md ou abrir issue.
