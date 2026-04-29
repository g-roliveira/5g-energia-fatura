# Backend Go API

Esta é a API nova do backend persistido.

Ela é diferente da API legada em Python documentada em [api.md](/home/gustavo/Projetos/5g-energia-fatura/docs/api.md).

Regra operacional:
- a Neoenergia é usada no momento da sincronização
- o frontend deve ler os dados persistidos do banco pelos endpoints `GET`
- o frontend não deve depender da resposta imediata do sync como fonte principal

## Estado atual

O backend já persiste:
- credenciais criptografadas
- sessões/tokens criptografados
- unidades consumidoras
- execuções de sincronização
- invoices
- snapshots da API privada
- PDF/base64 da fatura
- itens da fatura
- resultado da extração documental
- origem e confiança por campo

## Fluxo recomendado

1. `POST /v1/credentials`
2. `POST /v1/credentials/{id}/session`
3. `POST /v1/consumer-units/{uc}/sync`
4. `GET /v1/consumer-units/{uc}`
5. `GET /v1/consumer-units/{uc}/invoices`
6. `GET /v1/invoices/{id}`
7. `GET /v1/sync-runs/{id}`

## Subir localmente

```bash
cd services/backend-go
BACKEND_HOST=127.0.0.1 \
BACKEND_PORT=8088 \
EXTRACTOR_BASE_URL=http://127.0.0.1:8090 \
BACKEND_INTEGRATION_PG_URL='postgres://backoffice:backoffice@127.0.0.1:5432/backoffice?sslmode=disable' \
BACKEND_ENCRYPTION_KEY='troque-esta-chave' \
BOOTSTRAP_PYTHON_BIN="$PWD/../../.venv/bin/python" \
BOOTSTRAP_SCRIPT_PATH="$PWD/../../scripts/bootstrap_neoenergia_token.py" \
go run ./cmd/api
```

O extrator Python precisa estar no ar:

```bash
PYTHONPATH=src:services/doc-extractor-py \
DOC_EXTRACTOR_HOST=127.0.0.1 \
DOC_EXTRACTOR_PORT=8090 \
./.venv/bin/python services/doc-extractor-py/app/main.py
```

## Endpoints

### `GET /healthz`

Resposta:

```json
{
  "status": "ok",
  "service": "backend-go"
}
```

### `POST /v1/credentials`

Cria uma credencial criptografada para login na Neoenergia.

Request:

```json
{
  "label": "neo-paula",
  "documento": "03021937586",
  "senha": "senha-do-portal",
  "uf": "BA",
  "tipo_acesso": "normal"
}
```

Response:

```json
{
  "id": "f9ba880b894f63952add8c4feeb70d0c",
  "label": "neo-paula",
  "documento": "*******7586",
  "uf": "BA",
  "tipo_acesso": "normal",
  "created_at": "2026-04-17T00:00:00Z"
}
```

### `POST /v1/credentials/{id}/session`

Cria ou renova a sessão Neoenergia via bootstrap Playwright.

Response:

```json
{
  "id": "8b9ceb33f6877099fef8fb387bf3730e",
  "credential_id": "f9ba880b894f63952add8c4feeb70d0c",
  "created_at": "2026-04-17T00:00:30Z"
}
```

### `POST /v1/consumer-units/{uc}/sync`

Sincroniza uma UC, consulta a API privada, baixa PDF, extrai dados e persiste tudo.

Request:

```json
{
  "credential_id": "f9ba880b894f63952add8c4feeb70d0c",
  "include_pdf": true,
  "include_extraction": true
}
```

Campos importantes da resposta:
- `billing_record`
- `document_record`
- `persistence.sync_run_id`
- `persistence.invoice_id`
- blocos brutos da API privada para auditoria imediata

Trecho típico:

```json
{
  "uc": "007098175908",
  "billing_record": {
    "numero_fatura": "339800707843",
    "mes_referencia": "2026/04",
    "valor_total": "521.53",
    "completeness": {
      "status": "complete"
    }
  },
  "persistence": {
    "sync_run_id": "9b4c8b8d6a34e359d5add0397388fdce",
    "invoice_id": "bc874a8ce1af6601060bd12a91c1603f",
    "status": "succeeded"
  }
}
```

### `GET /v1/consumer-units`

Lista UCs persistidas no banco.

Query params:
- `limit`
- `status`

Exemplo:

```bash
curl 'http://127.0.0.1:8088/v1/consumer-units?limit=20&status=LIGADA'
```

### `GET /v1/consumer-units/{uc}`

Retorna a UC com status consolidado para o frontend.

Inclui:
- dados da UC
- `latest_invoice`
- `latest_sync_run`

Exemplo de resposta:

```json
{
  "uc": "007098175908",
  "status": "LIGADA",
  "nome_cliente": "PAULA",
  "latest_invoice": {
    "numero_fatura": "339800707843",
    "mes_referencia": "2026/04",
    "completeness_status": "complete"
  },
  "latest_sync_run": {
    "status": "succeeded"
  }
}
```

### `GET /v1/consumer-units/{uc}/invoices`

Lista invoices persistidas para a UC.

Query params:
- `limit`
- `status`

Exemplo:

```bash
curl 'http://127.0.0.1:8088/v1/consumer-units/007098175908/invoices?status=A%20Vencer&limit=1'
```

### `GET /v1/consumer-units/{uc}/latest-invoice`

Retorna a última invoice persistida da UC.

### `GET /v1/invoices/{id}`

Retorna a invoice persistida com:
- dados principais
- `billing_record`
- `document_record`
- `items`

### `GET /v1/sync-runs/{id}`

Retorna a auditoria de uma sincronização persistida.

Inclui:
- `status`
- `error_message`
- `raw_response`

Esse endpoint é o principal para a tela de acompanhamento de sincronização.

## Modelo para frontend

Tela de listagem de UCs:
- usar `GET /v1/consumer-units`

Tela de detalhe da UC:
- usar `GET /v1/consumer-units/{uc}`
- usar `GET /v1/consumer-units/{uc}/invoices`

Tela de detalhe da fatura:
- usar `GET /v1/invoices/{id}`

Ação de sincronizar:
- usar `POST /v1/consumer-units/{uc}/sync`
- depois recarregar `GET /v1/consumer-units/{uc}`

Tela de auditoria:
- usar `GET /v1/sync-runs/{id}`

## Status e erro

Hoje os estados principais são:
- `succeeded`
- `partial`
- `failed`

Leitura recomendada no frontend:
- mostrar `latest_sync_run.status`
- mostrar `latest_invoice.completeness_status`
- mostrar `latest_sync_run.error_message` quando existir

## Testes

Unit/integration do backend Go:

```bash
cd services/backend-go
go test ./...
```

E2E real validado:
- criação de credencial
- criação de sessão
- sync real da UC `007098175908`
- leitura posterior por:
  - `GET /v1/consumer-units/{uc}`
  - `GET /v1/consumer-units/{uc}/invoices`
  - `GET /v1/invoices/{id}`
  - `GET /v1/sync-runs/{id}`
