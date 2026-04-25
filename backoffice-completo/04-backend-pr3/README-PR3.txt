BACKOFFICE BILLING MODULE — PR 3/N
===================================

Este pacote contém o terceiro Pull Request do módulo de faturamento do backoffice Azi Dourado.

CONTEÚDO
--------
- backoffice-billing-pr3.bundle  → Git bundle com 5 commits organizados
- apply-pr3.sh                   → Script de aplicação automática
- PR3_README.md                  → Overview completo do PR3
- README.txt                     → Este arquivo

O QUE ESTÁ INCLUÍDO
-------------------
✓ PDF Generation — Chromium headless para gerar faturas do cliente
✓ Worker Pool — Goroutines com FOR UPDATE SKIP LOCKED
✓ Bulk Actions — sync_all, recalculate_all, generate_pdfs, approve_all
✓ LISTEN/NOTIFY — SSE realtime (não mais polling)
✓ Documentação completa de integração

NOVOS ENDPOINTS
---------------
POST   /v1/billing/cycles/{id}/bulk           → Bulk actions
POST   /v1/billing/calculations/{id}/pdf      → Gerar PDF do cliente
GET    /v1/billing/documents/{id}             → Download PDF gerado

ATUALIZADO:
GET    /v1/billing/events/cycles/{id}         → SSE com LISTEN/NOTIFY realtime

APLICAR NO REPOSITÓRIO
----------------------

1. Extrair este arquivo no repositório:
   cd ~/repos/5g-energia-fatura
   tar -xzf backoffice-billing-pr3.tar.gz

2. Executar script de aplicação:
   bash apply-pr3.sh

   OU aplicar manualmente:
   git bundle unbundle backoffice-billing-pr3.bundle
   git merge pr3-billing-pdf-workers

3. Instalar dependências Go:
   cd services/backend-go
   go get github.com/chromedp/chromedp@v0.9.3
   go get github.com/lib/pq@v1.10.9
   go mod tidy

4. Instalar chromium (sistema):
   # Ubuntu/Debian
   apt-get install chromium-browser
   
   # Alpine
   apk add chromium
   
   # macOS
   brew install chromium

5. Criar diretórios:
   mkdir -p /var/lib/backoffice/pdfs
   mkdir -p /opt/backoffice/templates
   chown appuser:appuser /var/lib/backoffice/pdfs

6. Integrar no código:
   Siga o guia em services/backend-go/internal/app/BILLING_INTEGRATION_PR3.md

ESTRUTURA DE COMMITS
--------------------
1. feat(billing): add PDF generation with chromium
2. feat(worker): add worker pool with FOR UPDATE SKIP LOCKED
3. feat(billing): add bulk actions and LISTEN/NOTIFY for SSE
4. docs(billing): add PR3 integration guide
5. docs: add PR3 README with overview and production guide

TESTAR APÓS APLICAR
-------------------

# 1. Bulk sync de todas UCs do ciclo
curl -X POST http://localhost:8080/v1/billing/cycles/{cycle_id}/bulk \
  -H "Content-Type: application/json" \
  -d '{
    "action": "sync",
    "uc_codes": [],
    "force_all": false
  }'

Resposta:
{
  "jobs_created": 30,
  "jobs_skipped": 2,
  "skipped_reasons": ["UC 007098175908 already synced"]
}

# 2. Bulk gerar PDFs (só para cálculos aprovados)
curl -X POST http://localhost:8080/v1/billing/cycles/{cycle_id}/bulk \
  -d '{"action": "generate_pdf"}'

# 3. Bulk aprovar tudo (sync)
curl -X POST http://localhost:8080/v1/billing/cycles/{cycle_id}/bulk \
  -d '{"action": "approve"}'

# 4. Gerar PDF individual
curl -X POST http://localhost:8080/v1/billing/calculations/{calc_id}/pdf \
  -H "Content-Type: application/json" \
  -d '{"generated_by": "admin-user-id"}'

# 5. Download PDF
curl http://localhost:8080/v1/billing/documents/{doc_id} --output fatura.pdf

# 6. SSE realtime (LISTEN/NOTIFY, não polling)
curl -N http://localhost:8080/v1/billing/events/cycles/{id}

# 7. Trigger evento manualmente (para testar SSE)
psql "$BACKOFFICE_PG_URL" -c "
  SELECT pg_notify('billing_cycle_events', 
    '{\"cycle_id\":\"uuid-do-ciclo\",\"event\":\"sync_progress\",\"data\":{\"uc_code\":\"007098175908\",\"status\":\"synced\"}}'
  );
"

# 8. Monitorar fila de jobs
psql "$BACKOFFICE_PG_URL" -c "
  SELECT status, count(*) FROM billing.sync_job GROUP BY status;
"

# 9. Ver jobs rodando
psql "$BACKOFFICE_PG_URL" -c "
  SELECT id, job_type, started_at, retry_count
  FROM billing.sync_job
  WHERE status = 'running';
"

# 10. Ver jobs falhados
psql "$BACKOFFICE_PG_URL" -c "
  SELECT id, job_type, error_message, retry_count
  FROM billing.sync_job
  WHERE status = 'failed'
  ORDER BY created_at DESC
  LIMIT 10;
"

DIFERENÇAS DO PR 2 → PR 3
-------------------------
| Feature         | PR 2                  | PR 3                      |
|-----------------|-----------------------|---------------------------|
| SSE             | Polling 5s            | LISTEN/NOTIFY realtime    |
| Jobs            | Placeholders          | Worker pool funcional     |
| Bulk actions    | Não existe            | sync/recalc/pdf/approve   |
| PDF generation  | Não existe            | Chromium headless         |
| Worker pool     | Não existe            | FOR UPDATE SKIP LOCKED    |

CONFIGURAÇÃO WORKER POOL
------------------------
No server.go:

pool := worker.NewPool(db, 10, 1*time.Second)
                         ^    ^
                         |    Poll interval
                         Pool size (workers concorrentes)

Recomendações:
- Desenvolvimento: 2-5 workers, 2s poll
- Produção: 10-20 workers, 1s poll
- High load: 50+ workers, 500ms poll

REGISTRO DE HANDLERS
--------------------
IMPORTANTE: Antes de iniciar o pool, registre handlers para cada job_type.

pool.RegisterHandler("sync_uc", func(ctx context.Context, job *worker.Job) error {
    // Lógica de sync
})

pool.RegisterHandler("recalculate", func(ctx context.Context, job *worker.Job) error {
    // Lógica de recálculo
})

pool.RegisterHandler("generate_pdf", func(ctx context.Context, job *worker.Job) error {
    // Lógica de geração de PDF
})

Depois:
pool.Start(context.Background())

LIMITAÇÕES SUPERADAS
---------------------
✅ SSE agora é realtime via LISTEN/NOTIFY
✅ Worker pool funcional (não mais síncronos)
✅ Bulk actions implementadas (4 tipos)
✅ PDF generation com chromium
✅ FOR UPDATE SKIP LOCKED para lock-free job claiming

PRÓXIMO PR
----------
PR 4/N incluirá:
- Suite de testes de integração com testcontainers
- OpenAPI spec completo (swagger UI)
- Templates de PDF profissionais (html/template)
- S3 storage para PDFs (opcional filesystem)
- Métricas Prometheus
- Health checks

AÇÕES PARA PRODUÇÃO
-------------------

1. Templates HTML Profissionais
   - Criar template em /opt/backoffice/templates/customer_invoice.html
   - Usar package html/template
   - Incluir: logo Azi Dourado, QR code PIX, instruções pagamento
   - Dados do cliente formatados BR
   - Detalhamento linha a linha com cores alternadas

2. Storage S3 (recomendado)
   - Migrar de filesystem para S3
   - URLs assinadas com expiração (1 hora)
   - Lifecycle policy: mover para Glacier após 90 dias

3. Monitoring & Alerting
   - Jobs processados/segundo (Prometheus)
   - Queue depth (alerta se > 100)
   - Tempo médio de geração de PDF (alerta se > 10s)
   - SSE connections ativas
   - Worker pool utilization

4. Tuning
   - Separate pools por job type (sync pool vs PDF pool)
   - Priority queues (cálculos urgentes primeiro)
   - Dead letter queue para jobs falhados > 3x

5. Segurança
   - Rate limiting em bulk actions (max 1 req/min por ciclo)
   - Audit log de quem disparou bulk
   - Validar permissões antes de enfileirar jobs

TROUBLESHOOTING
---------------

❌ Worker pool não processa jobs
   → Verificar:
     1. pool.Start(ctx) foi chamado?
     2. Handlers registrados antes de Start?
     3. SELECT count(*) FROM billing.sync_job WHERE status='pending';
     4. Locks no Postgres? SELECT * FROM pg_locks WHERE NOT granted;

❌ PDF generation falha com "chromium not found"
   → Instalar chromium:
     apt-get install chromium-browser
   → Ou setar env var:
     export CHROME_BIN=/usr/bin/chromium-browser

❌ SSE não recebe eventos
   → Verificar:
     1. NotifyHub criado e listening?
     2. pg_notify sendo chamado? (logs do app)
     3. SELECT * FROM pg_stat_activity WHERE query LIKE '%LISTEN%';
     4. Firewall bloqueando? (SSE usa long-polling HTTP)

❌ Bulk sync pula todas UCs
   → Ver skipped_reasons na resposta
   → Comum: "already synced" → use force_all: true
   → Ou: "no active contract" → vincular contrato primeiro

❌ PDF vazio ou com layout quebrado
   → Template HTML inválido
   → Chromium não consegue carregar fonts
   → Verificar logs do chromedp

❌ Jobs ficam stuck em "running"
   → Worker crashou durante execução
   → Timeout não configurado
   → Reset manual:
     UPDATE billing.sync_job SET status='pending' WHERE status='running' AND started_at < NOW() - interval '10 minutes';

DEPENDÊNCIAS DE SISTEMA
-----------------------
- Chromium (browser headless)
- Postgres 9.4+ (LISTEN/NOTIFY support)
- Go 1.22+

Dependências Go:
- github.com/chromedp/chromedp@v0.9.3
- github.com/lib/pq@v1.10.9
- (outras já instaladas nos PRs anteriores)

PRÉ-REQUISITOS
--------------
- PR 1/N aplicado (calc engine, contracts, migrations 000001)
- PR 2/N aplicado (cycles, adjustments, migration 000002)
- BACKOFFICE_PG_URL configurado
- golang-migrate instalado

DOCUMENTAÇÃO COMPLETA
---------------------
Após aplicar, consulte:
- PR3_README.md — Overview geral
- services/backend-go/internal/app/BILLING_INTEGRATION_PR3.md — Guia técnico detalhado

EXEMPLO COMPLETO DE INTEGRAÇÃO
-------------------------------
Ver arquivo BILLING_INTEGRATION_PR3.md, seção "Complete Example".
Inclui:
- Setup de BillingDeps com PDF e Worker
- Registro de handlers
- Start do worker pool
- Graceful shutdown
- Router setup

SUPORTE
-------
Dúvidas sobre integração: consultar BILLING_INTEGRATION_PR3.md
Bugs encontrados: reportar com logs e passos para reproduzir
Performance issues: verificar tuning do worker pool primeiro

CHANGELOG COMPLETO (PR 1+2+3)
-----------------------------
PR 1: calc engine, contracts, repos, handlers básicos
PR 2: cycles, adjustments, SSE com polling, migration notification
PR 3: PDF, worker pool, bulk actions, LISTEN/NOTIFY, SSE realtime

Total: 13 módulos, 3 migrations, ~3500 linhas Go

===========================================
Gerado em abril/2026 para Azi Dourado SaaS
===========================================
