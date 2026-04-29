# PR 3/N — PDF Generation, Worker Pool, Bulk Actions, LISTEN/NOTIFY

Este é o **terceiro Pull Request** do módulo de backoffice de faturamento.

## O que está incluído

### 1. Módulo `billing/pdf/`
- **service.go** — Geração de PDF com chromium headless
- Converte HTML para PDF via chromedp
- Salva em disco e registra em `billing.generated_document`
- Só gera para cálculos aprovados
- Template básico incluído (produção deve usar `html/template`)

### 2. Módulo `internal/worker/`
- **pool.go** — Pool de goroutines com FOR UPDATE SKIP LOCKED
- Consumo lock-free da fila `billing.sync_job`
- Handlers registráveis por job_type
- Retry automático com backoff
- Graceful shutdown

### 3. Módulo `billing/cycle/` — Bulk Actions
- **bulk.go** — Ações em massa (sync, recalculate, generate_pdf, approve)
- Enfileira jobs em paralelo
- Smart skipping (já sincronizado, PDF existe, etc)
- Modo force para ignorar checks
- Retorna `jobs_created` + `skipped_reasons`

### 4. Módulo `billing/cycle/` — LISTEN/NOTIFY
- **notify.go** — NotifyHub para pub/sub via Postgres
- SSE em tempo real (não mais polling)
- `pg_notify('billing_cycle_events', payload)`
- Subscribe/Unsubscribe por cycle_id
- Connection keep-alive com ping

## Novos endpoints

| Method | Path | Descrição |
|--------|------|-----------|
| POST | `/v1/billing/cycles/{id}/bulk` | Bulk actions |
| POST | `/v1/billing/calculations/{id}/pdf` | Gerar PDF do cliente |
| GET | `/v1/billing/documents/{id}` | Download do PDF |

**Atualizado:**
| GET | `/v1/billing/events/cycles/{id}` | SSE com LISTEN/NOTIFY real |

## Como aplicar

### Pré-requisitos
- PR 1/N e PR 2/N já aplicados
- Chromium instalado
- `BACKOFFICE_PG_URL` configurado

### Dependências

Adicionar ao `go.mod`:
```go
require (
    github.com/chromedp/chromedp v0.9.3
    github.com/lib/pq v1.10.9
)
```

Rodar:
```bash
go mod tidy
```

### Comandos

```bash
# 1. Extrair bundle
git bundle unbundle backoffice-billing-pr3.bundle

# 2. Merge
git merge pr3-billing-pdf-workers

# 3. Integrar no server.go
# Siga BILLING_INTEGRATION_PR3.md (essencial!)

# 4. Instalar chromium
apt-get install chromium-browser

# 5. Criar diretórios
mkdir -p /var/lib/backoffice/pdfs
mkdir -p /opt/backoffice/templates
```

## Estrutura de commits

1. `feat(billing): add PDF generation with chromium`
2. `feat(worker): add worker pool with FOR UPDATE SKIP LOCKED`
3. `feat(billing): add bulk actions and LISTEN/NOTIFY for SSE`
4. `docs(billing): add PR3 integration guide`

## Testar

### Bulk sync de todas UCs

```bash
curl -X POST http://localhost:8080/v1/billing/cycles/{cycle_id}/bulk \
  -H "Content-Type: application/json" \
  -d '{"action": "sync", "uc_codes": [], "force_all": false}'
```

Resposta:
```json
{
  "jobs_created": 30,
  "jobs_skipped": 2,
  "skipped_reasons": ["UC 007098175908 already synced"]
}
```

### Bulk gerar PDFs

```bash
curl -X POST http://localhost:8080/v1/billing/cycles/{cycle_id}/bulk \
  -d '{"action": "generate_pdf"}'
```

### Gerar PDF individual

```bash
curl -X POST http://localhost:8080/v1/billing/calculations/{calc_id}/pdf \
  -d '{"generated_by": "admin-user-id"}'
```

### Download PDF

```bash
curl http://localhost:8080/v1/billing/documents/{doc_id} --output fatura.pdf
```

### SSE com LISTEN/NOTIFY

```bash
curl -N http://localhost:8080/v1/billing/events/cycles/{id}
```

Agora envia eventos **imediatamente** via `pg_notify`, não a cada 5s.

Testar manualmente:
```sql
SELECT pg_notify('billing_cycle_events', 
  '{"cycle_id":"uuid","event":"sync_progress","data":{"uc_code":"007098175908","status":"synced"}}'
);
```

## Diferenças do PR 2

| Feature | PR 2 | PR 3 |
|---------|------|------|
| SSE | Polling 5s | LISTEN/NOTIFY realtime |
| Jobs | Placeholders | Worker pool real |
| Bulk | Não existe | sync/recalc/pdf/approve |
| PDF | Não existe | Chromium headless |

## Configuração Worker Pool

```go
// Tamanho = workers concorrentes
// PollInterval = frequência de check
pool := worker.NewPool(db, 10, 1*time.Second)
```

Recomendado:
- **Dev:** 2-5 workers, 2s
- **Prod:** 10-20 workers, 1s
- **High load:** 50+ workers, 500ms

## Arquivos modificados/criados

```
services/backend-go/
└── internal/
    ├── app/
    │   └── BILLING_INTEGRATION_PR3.md    [NEW]
    ├── worker/
    │   └── pool.go                       [NEW]
    └── billing/
        ├── pdf/
        │   └── service.go                [NEW]
        └── cycle/
            ├── bulk.go                   [NEW]
            ├── notify.go                 [NEW]
            ├── sse.go                    [MODIFIED - LISTEN/NOTIFY]
            └── handlers.go               [MODIFIED - bulk endpoint]
```

## Limitações superadas

✅ SSE agora é realtime (não mais polling)
✅ Worker pool funcional (não mais síncronos)
✅ Bulk actions implementadas
✅ PDF generation funcional

## Próximo PR

PR 4/N incluirá:
- Suite de testes de integração com testcontainers
- OpenAPI spec completo (swagger)
- Templates de PDF profissionais
- Armazenamento S3 para PDFs
- Métricas e monitoring

## Produção — Ações necessárias

1. **Templates HTML profissionais**
   - Criar `/opt/backoffice/templates/customer_invoice.html`
   - Incluir logo, QR code, instruções de pagamento

2. **Storage S3** (opcional mas recomendado)
   - Migrar de filesystem para S3
   - URLs assinadas com expiração

3. **Monitoring**
   - Jobs processados/segundo
   - Queue depth
   - Tempo de geração de PDF
   - SSE connections ativas

4. **Tuning**
   - Ajustar pool size baseado em carga
   - Separate pools por job type
   - Dead letter queue para jobs falhados

## Troubleshooting

### Worker pool não processa jobs

Verificar:
1. Pool iniciado? `pool.Start(ctx)`
2. Handlers registrados? `pool.RegisterHandler(...)`
3. Jobs na fila? `SELECT count(*) FROM billing.sync_job WHERE status='pending'`

### PDF generation falha

Instalar chromium:
```bash
apt-get install chromium-browser
# ou
apk add chromium
```

### SSE não recebe eventos

1. NotifyHub criado?
2. `pg_notify` sendo chamado?
3. Verificar: `SELECT * FROM pg_stat_activity WHERE query LIKE '%LISTEN%'`

## Dependências de sistema

- **Chromium:** `apt-get install chromium-browser` ou `apk add chromium`
- **Postgres 9.4+** (LISTEN/NOTIFY)
- **Go 1.22+**

## Contato

Dúvidas: ver BILLING_INTEGRATION_PR3.md ou abrir issue.
