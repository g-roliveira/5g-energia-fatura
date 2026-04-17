# Backend Monorepo

## Objetivo

Reorganizar o projeto em um monorepo com dois serviços principais:

1. `services/backend-go`
2. `services/doc-extractor-py`

O backend em Go assume o papel de sistema principal:

- autenticação e sessão com a API privada da Neoenergia
- sincronização de clientes, UCs e faturas
- agendamento e execução de jobs
- persistência e idempotência
- API pública consumida pelo frontend
- leitura persistida via banco para consumo do frontend

O serviço Python assume o papel documental:

- parsing do PDF com `PyMuPDF`
- normalização de campos ricos da fatura
- fallback com `Mistral OCR`
- reconciliação e rastreabilidade por campo

## Decisão

Não adotar:

- `100% Go`: aumenta risco no parsing de PDF
- `100% Python`: não aproveita bem o ganho operacional da mudança para API privada
- muitos microserviços: multiplica custo operacional cedo demais

Adotar:

- monorepo
- dois deployables
- contratos explícitos entre serviços

## Estrutura

```text
/services
  /backend-go
  /doc-extractor-py
/packages
  /contracts
/docs
  /architecture
  /neoenergia-private-api
/src/fatura
  código Python atual, mantido durante a transição
```

## Responsabilidades

### `backend-go`

- obter token/session bootstrap da Neoenergia
- usar a API privada como fonte principal
- consultar:
  - UCs
  - dados do imóvel
  - faturas
  - histórico de consumo
  - dados de pagamento
  - PDF base64
- decidir quando chamar o extrator documental
- publicar API pública previsível ao frontend
- gravar origem dos dados no banco

### `doc-extractor-py`

- receber PDF já baixado ou payload equivalente
- extrair:
  - itens da fatura
  - composição do fornecimento
  - dados fiscais ricos
  - demais campos não cobertos pela API privada
- preencher `source_map`
- preencher `confidence`
- usar `Mistral` apenas quando:
  - o parse local falhar
  - o layout mudar
  - campos obrigatórios continuarem ausentes

## Fluxo principal

1. O backend em Go recebe um pedido de sincronização.
2. Ele autentica na Neoenergia e obtém token/sessão.
3. Ele consulta a API privada e salva dados operacionais.
4. Se os campos exigidos pelo produto já estiverem completos, encerra.
5. Se houver campos dependentes do documento, baixa o PDF.
6. Envia o PDF para o `doc-extractor-py`.
7. Recebe o payload enriquecido e persiste o resultado consolidado.
8. Expõe a resposta imediata do sync e também endpoints de leitura a partir do banco.

## Fonte por campo

### Preferência de fonte

1. API privada Neoenergia
2. PDF com `PyMuPDF`
3. `Mistral OCR`

### Regra

Sempre preferir a menor superfície de risco:

- não usar PDF para um campo que já exista de forma confiável na API privada
- não usar `Mistral` para um campo que o parser local já extraia bem

## Banco de dados

O banco precisa separar claramente:

- dados vindos da API privada
- dados vindos do PDF
- dados vindos do fallback OCR

Campos mínimos por registro consolidado:

- `source_map`
- `confidence_map`
- `document_version`
- `api_snapshot`
- `extractor_snapshot`
- `sync_job_id`

Leitura operacional para o frontend:

- `GET /v1/consumer-units`
- `GET /v1/consumer-units/{uc}`
- `GET /v1/consumer-units/{uc}/invoices`
- `GET /v1/invoices/{id}`
- `GET /v1/sync-runs/{id}`

## Contratos

Os contratos entre Go e Python ficam em `packages/contracts`.

Contratos iniciais:

- `extractor-request.schema.json`
- `extractor-response.schema.json`
- `billing-record.schema.json`

## Estratégia de migração

### Fase 1

- manter o código atual em `src/fatura`
- usar o mapeamento da API privada já descoberto
- estabilizar os contratos

### Fase 2

- implementar o cliente privado em Go
- usar o Python atual como base do extrator documental

### Fase 3

- parar de usar scraping de DOM como motor principal
- manter Playwright apenas para bootstrap/autenticação e observabilidade

### Fase 4

- migrar API pública e scheduler para Go
- reduzir `src/fatura` ao papel de serviço de extração

## Riscos controlados

### API privada

- menos frágil que o DOM
- ainda pode mudar, então precisa de snapshots e testes de integração reais

### PDF

- continua sendo necessário para campos detalhados
- precisa de parser local e fallback OCR

### reCAPTCHA / login

- login pode continuar exigindo browser em parte do fluxo
- isso não invalida a estratégia: o ganho vem do pós-login por HTTP puro

## Próximos passos

1. documentar formalmente os endpoints persistidos do backend Go
2. padronizar o contrato de status/erro para o frontend
3. adicionar paginação, filtros e ordenação onde necessário
4. evoluir retenção/armazenamento de PDFs e artefatos
