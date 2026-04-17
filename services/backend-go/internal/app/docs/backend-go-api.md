# Backend Go API

API persistida para sincronização e leitura de dados da Neoenergia.

## Health

- `GET /healthz`

## Credenciais e sessão

- `POST /v1/credentials`
- `POST /v1/credentials/{id}/session`

## Sincronização

- `POST /v1/sync/uc`
- `POST /v1/consumer-units/{uc}/sync`

## Leitura persistida

- `GET /v1/consumer-units`
- `GET /v1/consumer-units/{uc}`
- `GET /v1/consumer-units/{uc}/invoices`
- `GET /v1/consumer-units/{uc}/latest-invoice`
- `GET /v1/invoices/{id}`
- `GET /v1/sync-runs/{id}`

## Fluxo recomendado

1. Criar credencial
2. Criar sessão
3. Sincronizar UC
4. Ler do banco pelos endpoints `GET`

## Exemplo rápido

### Criar credencial

```json
{
  "label": "neo-paula",
  "documento": "03021937586",
  "senha": "senha-do-portal",
  "uf": "BA",
  "tipo_acesso": "normal"
}
```

### Sincronizar uma UC

```json
{
  "credential_id": "CREDENTIAL_ID",
  "include_pdf": true,
  "include_extraction": true
}
```

### Ler última invoice da UC

`GET /v1/consumer-units/007098175908`

Resposta típica:

```json
{
  "uc": "007098175908",
  "status": "LIGADA",
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

## Query params suportados

### `GET /v1/consumer-units`

- `limit`
- `status`

### `GET /v1/consumer-units/{uc}/invoices`

- `limit`
- `status`

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
