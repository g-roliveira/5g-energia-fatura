BACKOFFICE BILLING MODULE — PR 2/N
===================================

Este pacote contém o segundo Pull Request do módulo de faturamento do backoffice Azi Dourado.

CONTEÚDO
--------
- backoffice-billing-pr2.bundle  → Git bundle com 5 commits organizados
- apply-pr2.sh                   → Script de aplicação automática
- PR2_README.md                  → Overview completo do PR2
- README.txt                     → Este arquivo

O QUE ESTÁ INCLUÍDO
-------------------
✓ Migration 000002 — Tabela core.notification + índices de performance
✓ Módulo billing/cycle/ — Orquestração de ciclos (criar, listar, fechar)
✓ Módulo billing/adjustment/ — Ajustes manuais com versionamento
✓ Server-Sent Events (SSE) — Endpoint /v1/billing/events/cycles/{id}
✓ Documentação completa de integração

NOVOS ENDPOINTS
---------------
POST   /v1/billing/cycles                     → Criar competência
GET    /v1/billing/cycles                     → Listar competências
GET    /v1/billing/cycles/{id}                → Detalhe da competência
GET    /v1/billing/cycles/{id}/rows           → Tabela principal (dashboard)
POST   /v1/billing/cycles/{id}/close          → Fechar competência
GET    /v1/billing/events/cycles/{id}         → SSE realtime
POST   /v1/billing/calculations/{id}/adjust   → Aplicar ajuste manual
GET    /v1/billing/calculations/{id}/adjustments → Histórico de ajustes

APLICAR NO REPOSITÓRIO
----------------------

1. Extrair este arquivo no repositório:
   cd ~/repos/5g-energia-fatura
   tar -xzf backoffice-billing-pr2.tar.gz

2. Executar script de aplicação:
   bash apply-pr2.sh

   OU aplicar manualmente:
   git bundle unbundle backoffice-billing-pr2.bundle
   git merge pr2-billing-cycles

3. Rodar migration:
   cd services/backend-go
   migrate -path migrations -database "$BACKOFFICE_PG_URL" up

4. Integrar no código:
   Siga o guia em services/backend-go/internal/app/BILLING_INTEGRATION_PR2.md

ESTRUTURA DE COMMITS
--------------------
1. feat(billing): add migration 000002 - notification table and indices
2. feat(billing): add cycle orchestration module
3. feat(billing): add manual adjustment module
4. docs(billing): add PR2 integration guide
5. docs: add PR2 README with overview and instructions

TESTAR APÓS APLICAR
-------------------

# 1. Criar ciclo
curl -X POST http://localhost:8080/v1/billing/cycles \
  -H "Content-Type: application/json" \
  -d '{"year": 2026, "month": 4, "include_all_active": true, "created_by": "user-id"}'

# 2. Listar ciclos
curl http://localhost:8080/v1/billing/cycles

# 3. Ver linhas do ciclo (dashboard)
curl http://localhost:8080/v1/billing/cycles/{id}/rows?limit=50

# 4. SSE realtime
curl -N http://localhost:8080/v1/billing/events/cycles/{id}

# 5. Aplicar ajuste manual
curl -X POST http://localhost:8080/v1/billing/calculations/{calc_id}/adjust \
  -H "Content-Type: application/json" \
  -d '{
    "field": "ip_coelba",
    "old_value": 10.0,
    "new_value": 12.5,
    "reason": "Correção manual verificada",
    "adjusted_by": "admin-user-id"
  }'

# 6. Fechar ciclo (requer todos cálculos aprovados)
curl -X POST http://localhost:8080/v1/billing/cycles/{id}/close \
  -d '{"closed_by": "admin-user-id"}'

LIMITAÇÕES DO MVP
-----------------
- SSE usa polling (5s) em vez de LISTEN/NOTIFY (prod)
- Ajuste manual não reroda motor (apenas cria nova version)
- Sem PDF generation ainda (PR 3)
- Sem bulk actions ainda (PR 3)
- Sem worker pool ainda (PR 3)

PRÓXIMO PR
----------
PR 3/N incluirá:
- PDF generation com chromium
- Worker pool para jobs assíncronos
- Bulk actions (sync_all, recalculate_all)
- LISTEN/NOTIFY para SSE realtime
- Testes de integração
- OpenAPI completo

PRÉ-REQUISITOS
--------------
- PR 1/N já aplicado
- Go 1.22+
- Postgres com schema billing.* criado
- BACKOFFICE_PG_URL configurado
- golang-migrate instalado

DOCUMENTAÇÃO COMPLETA
---------------------
Após aplicar, consulte:
- PR2_README.md — Overview geral
- services/backend-go/internal/app/BILLING_INTEGRATION_PR2.md — Guia técnico detalhado

SUPORTE
-------
Dúvidas sobre integração: abrir issue ou consultar BILLING_INTEGRATION_PR2.md
Bugs encontrados: reportar com logs e passos para reproduzir

===========================================
Gerado em abril/2026 para Azi Dourado SaaS
===========================================
