# Como integrar o módulo billing no `server.go` existente

O arquivo novo `internal/app/server_billing.go` faz todo o trabalho de
registrar rotas e handlers. Para ativá-lo, você precisa fazer 3 pequenas
mudanças no `server.go` já existente.

Essas mudanças seguem exatamente o padrão que seu backend-go já usa
(`routeCatalog.add` + `mux.HandleFunc`), então elas aparecem
automaticamente em `/docs` (Swagger UI) e `/openapi.json`.

---

## Mudança 1: adicionar `BackofficePGURL` em `internal/app/config.go`

Em `config.go`, no struct `Config`, adicione um campo e leia do env
em `LoadConfigFromEnv`:

```go
type Config struct {
    // ... campos existentes ...

    // Postgres do backoffice (domínio de faturamento + cadastro).
    // Se vazio, módulo billing fica inativo.
    BackofficePGURL string
}

func LoadConfigFromEnv() Config {
    return Config{
        // ... existentes ...
        BackofficePGURL: envOrDefault("BACKOFFICE_PG_URL", ""),
    }
}
```

## Mudança 2: wire do pool + registro de rotas em `server.go`

Em `server.go`, no `NewServer`, **depois** de todas as `mux.HandleFunc`
que já existem e **antes** do `return`, adicione o bloco:

```go
// ...todas as rotas existentes já registradas acima...

// Módulo billing — só ativa se tiver Postgres configurado.
if cfg.BackofficePGURL != "" {
    pool, err := pgstore.Open(context.Background(), pgstore.LoadConfigFromEnv())
    if err != nil {
        return nil, fmt.Errorf("pgstore.Open: %w", err)
    }
    deps := NewBillingDeps(pool)
    RegisterBillingRoutes(mux, docs, deps, logger)
    logger.Info("billing_module_enabled", "routes", "/v1/billing/*")
}

rootHandler := withRequestLogging(logger, mux)
return &Server{ ... }, nil
```

E adicione os imports:

```go
import (
    // existentes...
    "context"

    "github.com/gustavo/5g-energia-fatura/services/backend-go/internal/pgstore"
)
```

## Mudança 3: adicionar dependências em `go.mod`

No `go.mod` do `services/backend-go/`, adicione:

```
require (
    github.com/jackc/pgx/v5 v5.7.1
    github.com/google/uuid v1.6.0
    github.com/shopspring/decimal v1.4.0

    github.com/gustavo/5g-energia-fatura/packages/calc-engine v0.0.0-00010101000000-000000000000
    github.com/gustavo/5g-energia-fatura/packages/normalizer  v0.0.0-00010101000000-000000000000
)

replace (
    github.com/gustavo/5g-energia-fatura/packages/calc-engine => ../../packages/calc-engine
    github.com/gustavo/5g-energia-fatura/packages/normalizer  => ../../packages/normalizer
)
```

Depois roda:

```
cd services/backend-go && go mod tidy
```

---

## Validação final

1. `docker compose` com Postgres local (ou só `psql` contra um banco
   vazio), e roda:

    ```
    migrate -path services/backend-go/migrations \
            -database "postgres://user:pass@localhost/backoffice?sslmode=disable" \
            up
    ```

2. Sobe o backend-go com `BACKOFFICE_PG_URL` setado:

    ```
    BACKOFFICE_PG_URL="postgres://user:pass@localhost/backoffice?sslmode=disable" \
      go run ./cmd/api
    ```

3. Acessa `http://localhost:8080/docs` — você deve ver as 4 rotas novas
   sob a tag `billing`:

    - `POST /v1/billing/contracts`
    - `GET  /v1/billing/contracts/{id}`
    - `GET  /v1/billing/consumer-units/{uc_id}/active-contract`
    - `GET  /v1/billing/calculations/{id}`

4. Cria um contrato:

    ```bash
    curl -X POST http://localhost:8080/v1/billing/contracts \
      -H 'Content-Type: application/json' \
      -d '{
        "customer_id": "...uuid do customer...",
        "consumer_unit_id": "...uuid da UC...",
        "vigencia_inicio": "2025-10-01",
        "fator_repasse_energia": "0.85",
        "ip_faturamento_mode": "fixed",
        "ip_faturamento_valor": "10.00",
        "bandeira_com_desconto": false,
        "custo_disponibilidade_sempre_cobrado": true
      }'
    ```

    Resposta: 201 com o contrato criado. Se já houver contrato ativo pra
    mesma UC, ele é fechado automaticamente (vigencia_fim = 30/09/2025)
    na mesma transação.

---

## Observações

- **Segurança de write**: os handlers aceitam `created_by` no body, mas
  num próximo PR isso deve vir do middleware de auth (quando o
  `app_user` / bearer token estiver implementado). Por ora aceita o
  UUID explícito pra permitir testes.
- **Encoding de datas**: sempre `YYYY-MM-DD` em requests e responses
  (padrão ISO 8601 date-only — o que a maioria dos frameworks de form
  espera).
- **OpenAPI schemas detalhados**: o `routeCatalog` atual só registra
  método+path+summary+tags. Os schemas completos de request/response
  ainda são documentação de texto (aqui neste README + via `/docs.md`).
  Um próximo PR pode enriquecer `openapi_routes.go` pra gerar schemas
  completos — não fiz agora pra não aumentar demais o diff.
