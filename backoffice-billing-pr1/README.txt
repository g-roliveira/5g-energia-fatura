═══════════════════════════════════════════════════════════════════════
  PR Backoffice Billing — 1/N
  Módulo de faturamento de energia compartilhada (Azi Dourado)
═══════════════════════════════════════════════════════════════════════

ARQUIVOS NESTE PACOTE:

  backoffice-billing.bundle   bundle git com 8 commits organizados
  apply.sh                    script que importa no seu repo
  README.txt                  você está lendo

COMO APLICAR (2 comandos):

  cd ~/Projetos/5g-energia-fatura
  /caminho/ate/apply.sh

Depois: leia services/backend-go/internal/app/BILLING_INTEGRATION.md
e git push origin feat/backoffice-billing.

───────────────────────────────────────────────────────────────────────
  O QUE TEM DENTRO
───────────────────────────────────────────────────────────────────────

8 commits, revisáveis 1 por 1:

  1. packages/calc-engine/      motor puro, 4 testes passando
  2. packages/normalizer/       classifier + SCEE, 4 testes passando
  3. migrations Postgres        2 schemas, 13 tabelas, DDL completo
  4. pgstore/                   pool pgx/v5 pro Postgres do backoffice
  5. billing/repo/              ContractRepo + CalculationRepo
  6. billing/contract/          service com versionamento por vigência
  7. handlers HTTP              4 rotas novas no OpenAPI existente
  8. README.md                  guia do PR

VALIDAÇÕES FEITAS:

  ✓ todos os 17 arquivos Go passam parse AST
  ✓ gofmt limpo em todos
  ✓ testes unitários dos packages rodam e passam
  ✓ sintaxe SQL das migrations conferida
  ✓ 100% dos comentários são novos ou realinhados
    ao sync.BillingRecord real do backend-go existente

O QUE NÃO TEM (próximo PR):

  - billing/cycle/        orquestração de ciclo
  - billing/adjustment/   ajustes manuais versionados
  - billing/pdf/          geração do PDF do cliente
  - internal/worker/      pool de goroutines consumindo sync_job
  - SSE de progresso
  - Testes de integração ponta-a-ponta
  - OpenAPI com schemas completos de request/response

Nada disso quebra o que chegou neste PR. É continuação natural.

───────────────────────────────────────────────────────────────────────
  INTEGRAÇÃO NO SERVER.GO EXISTENTE
───────────────────────────────────────────────────────────────────────

São 3 mudanças pequenas (documentadas em detalhe em
services/backend-go/internal/app/BILLING_INTEGRATION.md após importar):

  1. Adicionar campo BackofficePGURL em internal/app/config.go
  2. Chamar pgstore.Open + NewBillingDeps + RegisterBillingRoutes no
     NewServer de internal/app/server.go
  3. Adicionar 3 deps + 2 replace no go.mod

Se BACKOFFICE_PG_URL não estiver setada, módulo fica dormente e o
backend-go roda exatamente como antes. Rollout incremental seguro.

───────────────────────────────────────────────────────────────────────
