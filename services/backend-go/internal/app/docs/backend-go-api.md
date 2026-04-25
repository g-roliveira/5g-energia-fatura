# Backend Go API

API persistida para sincronização e leitura de dados da Neoenergia.

## Base URL

- Local: `http://127.0.0.1:8080`

## Autenticação

A API suporta autenticação opcional por chave:

- Variável: `BACKEND_API_KEY`
- Header: `X-API-Key: <valor-da-chave>`

Regras:

- Se `BACKEND_API_KEY` estiver vazio, endpoints `/v1/*` ficam públicos.
- Se `BACKEND_API_KEY` estiver definido, todos os endpoints `/v1/*` exigem `X-API-Key`.
- Endpoints de infra (`/healthz`, `/docs`, `/docs.md`, `/openapi.json`) permanecem públicos.

## Infra

- `GET /healthz`
- `GET /openapi.json`
- `GET /docs`
- `GET /docs.md`

## Credenciais e sessão

- `POST /v1/credentials`
- `POST /v1/credentials/{id}/session`
- `GET /v1/credentials/{id}/discover`

## Sincronização

- `POST /v1/sync/uc`
- `POST /v1/consumer-units/{uc}/sync`
- `GET /v1/sync-runs/{id}`

## Leitura persistida

- `GET /v1/consumer-units`
- `GET /v1/consumer-units/{uc}`
- `GET /v1/consumer-units/{uc}/invoices`
- `GET /v1/consumer-units/{uc}/latest-invoice`
- `GET /v1/invoices/{id}`

## Contratos do extrator

- `GET /v1/extractor/contracts`

## Fluxo recomendado

1. Criar credencial (`POST /v1/credentials`)
2. Criar sessão (`POST /v1/credentials/{id}/session`)
3. (Opcional) Descobrir perfil e UCs (`GET /v1/credentials/{id}/discover`)
4. Sincronizar UC (`POST /v1/sync/uc` ou `POST /v1/consumer-units/{uc}/sync`)
5. Ler dados persistidos pelos `GET /v1/*`

## Exemplos

### Criar credencial

```bash
curl -X POST 'http://127.0.0.1:8080/v1/credentials' \
  -H 'Content-Type: application/json' \
  -H 'X-API-Key: SUA_CHAVE' \
  -d '{
    "label": "neo-paula",
    "documento": "03021937586",
    "senha": "MinhaSenha@123",
    "uf": "BA",
    "tipo_acesso": "normal"
  }'
```

Resposta típica:

```json
{
  "id": "a0f6f5c4f2a44511b7af35d6b2e9f893",
  "label": "neo-paula",
  "documento": "*******7586",
  "uf": "BA",
  "tipo_acesso": "normal",
  "created_at": "2026-04-18T13:10:00Z"
}
```

### Criar sessão por credencial

```bash
curl -X POST 'http://127.0.0.1:8080/v1/credentials/a0f6f5c4f2a44511b7af35d6b2e9f893/session' \
  -H 'X-API-Key: SUA_CHAVE'
```

### Sincronizar UC com credencial salva

```bash
curl -X POST 'http://127.0.0.1:8080/v1/sync/uc' \
  -H 'Content-Type: application/json' \
  -H 'X-API-Key: SUA_CHAVE' \
  -d '{
    "credential_id": "a0f6f5c4f2a44511b7af35d6b2e9f893",
    "uc": "007098175908",
    "include_pdf": true,
    "include_extraction": true
  }'
```

### Sincronizar UC com token manual

```bash
curl -X POST 'http://127.0.0.1:8080/v1/sync/uc' \
  -H 'Content-Type: application/json' \
  -H 'X-API-Key: SUA_CHAVE' \
  -d '{
    "bearer_token": "eyJhbGciOi...",
    "documento": "03021937586",
    "uc": "007098175908",
    "include_pdf": true
  }'
```

### Ler faturas de uma UC

```bash
curl 'http://127.0.0.1:8080/v1/consumer-units/007098175908/invoices?limit=20&status=A%20Vencer' \
  -H 'X-API-Key: SUA_CHAVE'
```

## Query params suportados

### `GET /v1/consumer-units`

- `limit` (default: 100)
- `status`

### `GET /v1/consumer-units/{uc}/invoices`

- `limit` (default: 100)
- `status`

## Modelo de erro

Erros retornam JSON no formato:

```json
{ "error": "codigo_ou_mensagem" }
```

Códigos comuns:

- `400`: payload inválido (`invalid_json`, validação de campos)
- `401`: sem `X-API-Key` válido quando proteção está ativa
- `404`: recurso não encontrado
- `405`: método não permitido
- `500`: erro interno (`internal_error`)

## Persistência

O backend grava:

- credenciais criptografadas
- sessões/tokens criptografados
- unidades consumidoras
- sync runs
- invoices
- snapshots da API privada
- PDF/base64
- itens da fatura
- resultado da extração
- origem/confiança por campo
