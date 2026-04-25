# Contrato de Mapeamento: `SyncUCResponse` -> `utility_invoice_item`

Baseado em chamadas reais da API Go em **18/04/2026**.

## Evidências usadas

1. UC alvo pedida (`7065654738`) — cenário de erro de permissão:
- [`syncuc-7065654738-response-redacted.json`](/home/gustavo/Projetos/5g-energia-fatura/docs/samples/syncuc-7065654738-response-redacted.json)
- Resultado: `403` nos blocos de fatura, sem `billing_record`.

2. UC autorizada do mesmo cliente (`007098175908`) — cenário de sucesso (usado para shape real de itens):
- [`syncuc-007098175908-response-redacted.json`](/home/gustavo/Projetos/5g-energia-fatura/docs/samples/syncuc-007098175908-response-redacted.json)
- [`syncuc-007098175908-summary.json`](/home/gustavo/Projetos/5g-energia-fatura/docs/samples/syncuc-007098175908-summary.json)

## Fontes de itens na resposta

Ordem de prioridade para normalização:

1. `billing_record.itens_fatura[]`
2. fallback: `document_record.ocr.itens_fatura[]`

No payload real de sucesso, os dois vieram com o mesmo shape por item:

```json
{
  "codigo": "",
  "descricao": "Consumo-TUSD",
  "quantidade": "418,00",
  "quantidade_residual": "",
  "quantidade_faturada": "",
  "tarifa": "0,56",
  "valor": "315,48",
  "base_icms": "315,48",
  "aliq_icms": "20,50",
  "icms": "64,67",
  "valor_total": "315,48"
}
```

## Regras de normalização de número

Todos os campos monetários/quantitativos em `itens_fatura` chegam como string em formato BR.

- Entrada BR: `1.234,56`
- Saída decimal canônica: `1234.56`
- String vazia (`""`) -> `null`

Campos com conversão obrigatória:
- `quantidade`
- `quantidade_residual`
- `quantidade_faturada`
- `tarifa`
- `valor`
- `base_icms`
- `icms`
- `valor_total`
- `aliq_icms` (percentual; manter número decimal sem `%`)

## Contrato sugerido de mapeamento

Mapeamento `source -> utility_invoice_item`:

- `codigo` -> `item_code` (string|null)
- `descricao` -> `description` (string, obrigatório)
- `quantidade` -> `quantity` (decimal|null)
- `quantidade_residual` -> `residual_quantity` (decimal|null)
- `quantidade_faturada` -> `billed_quantity` (decimal|null)
- `tarifa` -> `unit_rate` (decimal|null)
- `valor` -> `amount` (decimal|null)
- `base_icms` -> `icms_base` (decimal|null)
- `aliq_icms` -> `icms_rate` (decimal|null)
- `icms` -> `icms_amount` (decimal|null)
- `valor_total` -> `line_total` (decimal|null)

Campos de contexto (do header da fatura):

- `billing_record.numero_fatura` -> `invoice_number`
- `billing_record.mes_referencia` -> `reference_month`
- `billing_record.uc` -> `consumer_unit`
- `billing_record.status_fatura` -> `invoice_status`

## Regras de robustez

1. Se `billing_record` for `null` e `faturas.error.status_code == 403`:
- classificar como `forbidden_uc`
- não gerar `utility_invoice_item`
- persistir erro técnico e payload bruto para auditoria

2. Se `itens_fatura` não vier ou vier vazio:
- fallback para `document_record.ocr.itens_fatura`
- se ambos vazios: classificar como `missing_items`

3. `descricao` é a chave de negócio mínima por linha.
- Se faltar `descricao`, descartar linha e registrar warning.

## Observação importante para Fase 2

Para fechar `7065654738` com dados de item reais, precisa de credencial com permissão nessa UC.
Com a credencial usada nessa execução, a API Neoenergia retornou `403` para essa UC específica.
