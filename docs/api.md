# API HTTP

Esta API expõe o MVP como serviço assíncrono interno.

Documentação automática:
- Swagger UI: `GET /docs`
- ReDoc: `GET /redoc`
- OpenAPI JSON: `GET /openapi.json`

Autenticação:
- Envie `X-API-Key` quando `service.api_key` estiver configurado.

Fluxo recomendado:
1. `POST /jobs/faturas`
2. `GET /jobs/{job_id}`
3. `GET /jobs/{job_id}/result`

## Subir a API

```bash
./.venv/bin/fatura-api
```

## Exportar OpenAPI

```bash
./.venv/bin/fatura-openapi
```

Arquivo gerado:
- [openapi.json](/home/gustavo/Projetos/5g-energia-fatura/docs/openapi.json)

## Health check

```bash
curl http://127.0.0.1:8000/health
```

Resposta típica:

```json
{
  "status": "ok",
  "reset_jobs": 0
}
```

## Criar job

```bash
curl -X POST http://127.0.0.1:8000/jobs/faturas \
  -H 'Content-Type: application/json' \
  -H 'X-API-Key: troque-este-token' \
  -d '{
    "cpf_cnpj": "12345678901",
    "senha_portal": "minha_senha_portal",
    "uf": "BA",
    "tipo_acesso": "normal",
    "mes_ano": "122024",
    "force": false,
    "ucs": [
      {"uc": "007085489032", "nome": "Paula Fernandes"}
    ]
  }'
```

Resposta típica:

```json
{
  "job_id": "8b4f0ab3-9e75-4f3b-9326-88f07c9d4d6d",
  "status": "queued",
  "created_at": "2026-04-16T12:00:00.000000",
  "started_at": null,
  "finished_at": null,
  "progress_total": 1,
  "progress_done": 0,
  "summary": {
    "total": 1,
    "completed": 0,
    "success": 0,
    "error": 0
  }
}
```

## Consultar status

```bash
curl -H 'X-API-Key: troque-este-token' \
  http://127.0.0.1:8000/jobs/JOB_ID
```

Estados esperados:
- `queued`
- `running`
- `succeeded`
- `partial_failure`
- `failed`

## Consultar resultado

```bash
curl -H 'X-API-Key: troque-este-token' \
  http://127.0.0.1:8000/jobs/JOB_ID/result
```

Campos relevantes por item:
- `status`
- `mensagem`
- `erro_tipo`
- `pdf_path`
- `screenshot_path`
- `html_path`
- `mes`
- `ano`
- `valor`
- `attempts`

## Testes automatizados da API

Testes HTTP sem portal real:

```bash
./.venv/bin/pytest -q tests/test_api.py
```

Teste ponta a ponta da API com portal real:

```bash
RUN_REAL_API_E2E=1 ./.venv/bin/pytest -q tests/test_real_api_e2e.py -s
```

Esse teste:
- sobe a API em subprocesso com config temporário
- chama `/docs` e `/openapi.json`
- cria um job real via `POST /jobs/faturas`
- faz polling de `/jobs/{job_id}`
- valida `/jobs/{job_id}/result`
- confirma que o PDF retornado existe e é um arquivo PDF válido
