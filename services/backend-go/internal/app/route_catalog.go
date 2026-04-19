package app

import (
	"net/http"
	"strings"
)

type routeCatalog struct {
	routes []routeOperation
}

type routeOperation struct {
	Method        string
	Path          string
	Summary       string
	Tags          []string
	SuccessStatus int
}

func newRouteCatalog() *routeCatalog {
	return &routeCatalog{
		routes: make([]routeOperation, 0, 16),
	}
}

func (c *routeCatalog) add(method, path, summary string, tags []string, successStatus int) {
	c.routes = append(c.routes, routeOperation{
		Method:        method,
		Path:          path,
		Summary:       summary,
		Tags:          tags,
		SuccessStatus: successStatus,
	})
}

func (c *routeCatalog) spec() map[string]any {
	paths := map[string]any{}
	for _, route := range c.routes {
		operation := c.buildOperation(route)
		if operation == nil {
			continue
		}
		pathItem, ok := paths[route.Path].(map[string]any)
		if !ok {
			pathItem = map[string]any{}
			paths[route.Path] = pathItem
		}
		pathItem[strings.ToLower(route.Method)] = operation
	}

	return map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title":   "5G Energia Fatura Backend Go API",
			"version": "1.0.0",
			"summary": "API persistida para sincronização e leitura de dados da Neoenergia",
			"description": strings.TrimSpace(`
API para:
- criar credenciais com criptografia em repouso;
- criar/reusar sessão Neoenergia;
- sincronizar UC e persistir dados normalizados;
- consultar UCs/faturas/sync-runs persistidos.

Autenticação da API:
- se ` + "`BACKEND_API_KEY`" + ` estiver configurado, enviar header ` + "`X-API-Key`" + ` em todos os endpoints ` + "`/v1/*`" + `;
- endpoints de infra (` + "`/healthz`" + `, ` + "`/docs`" + `, ` + "`/openapi.json`" + `, ` + "`/docs.md`" + `) permanecem públicos.

Fluxo recomendado:
1. ` + "`POST /v1/credentials`" + `
2. ` + "`POST /v1/credentials/{id}/session`" + ` (opcional: ` + "`GET /v1/credentials/{id}/discover`" + `)
3. ` + "`POST /v1/sync/uc`" + ` ou ` + "`POST /v1/consumer-units/{uc}/sync`" + `
4. consultas via endpoints ` + "`GET /v1/*`" + `.
`),
		},
		"tags": []any{
			map[string]any{"name": "infra", "description": "Healthcheck, OpenAPI e documentação."},
			map[string]any{"name": "credentials", "description": "Gestão de credenciais e sessões Neoenergia."},
			map[string]any{"name": "sync", "description": "Sincronização com Neoenergia e estado do processamento."},
			map[string]any{"name": "consumer-units", "description": "Consulta de unidades consumidoras persistidas."},
			map[string]any{"name": "invoices", "description": "Consulta de faturas persistidas."},
			map[string]any{"name": "extractor", "description": "Descoberta dos contratos JSON Schema do extrator."},
			map[string]any{"name": "billing", "description": "Domínio de faturamento no Postgres do backoffice."},
			map[string]any{"name": "contracts", "description": "Versionamento e leitura de contratos comerciais."},
			map[string]any{"name": "calculations", "description": "Consulta de cálculos de faturamento."},
		},
		"paths":      paths,
		"components": componentsSpec(),
	}
}

func (c *routeCatalog) buildOperation(route routeOperation) map[string]any {
	key := route.Method + " " + route.Path
	op := map[string]any{
		"summary":     route.Summary,
		"operationId": operationID(route.Method, route.Path),
		"tags":        route.Tags,
	}
	if requiresAPIKey(route.Path) {
		op["security"] = []any{map[string]any{"ApiKeyAuth": []any{}}}
	}

	switch key {
	case http.MethodGet + " /healthz":
		op["description"] = "Verifica disponibilidade do serviço."
		op["responses"] = map[string]any{
			"200": responseJSON("Status da API.", "#/components/schemas/HealthResponse", map[string]any{
				"status":  "ok",
				"service": "backend-go",
			}),
		}
	case http.MethodGet + " /openapi.json":
		op["description"] = "Schema OpenAPI desta API."
		op["responses"] = map[string]any{
			"200": responseRaw("OpenAPI JSON.", "application/json"),
		}
	case http.MethodGet + " /docs":
		op["description"] = "Swagger UI gerado a partir de /openapi.json."
		op["responses"] = map[string]any{
			"200": responseRaw("Página HTML do Swagger UI.", "text/html"),
		}
	case http.MethodGet + " /docs.md":
		op["description"] = "Guia operacional em Markdown."
		op["responses"] = map[string]any{
			"200": responseRaw("Documentação markdown.", "text/markdown"),
		}
	case http.MethodPost + " /v1/credentials":
		op["description"] = "Cria uma credencial criptografada para uso em sessões/sincronização."
		op["requestBody"] = map[string]any{
			"required": true,
			"content": map[string]any{
				"application/json": map[string]any{
					"schema": refSchema("#/components/schemas/CredentialCreateRequest"),
					"examples": map[string]any{
						"cpf": map[string]any{
							"summary": "Credencial CPF",
							"value": map[string]any{
								"label":       "neo-paula",
								"documento":   "03021937586",
								"senha":       "MinhaSenha@123",
								"uf":          "BA",
								"tipo_acesso": "normal",
							},
						},
					},
				},
			},
		}
		op["responses"] = withDefaultErrors(map[string]any{
			"201": responseJSON("Credencial criada.", "#/components/schemas/CredentialView", map[string]any{
				"id":          "a0f6f5c4f2a44511b7af35d6b2e9f893",
				"label":       "neo-paula",
				"documento":   "*******7586",
				"uf":          "BA",
				"tipo_acesso": "normal",
				"created_at":  "2026-04-18T13:10:00Z",
			}),
		}, false, true)
	case http.MethodPost + " /v1/credentials/{id}/session":
		op["description"] = "Cria uma sessão Neoenergia para a credencial informada. O token é armazenado criptografado."
		op["parameters"] = []any{pathParam("id", "ID da credencial.")}
		op["responses"] = withDefaultErrors(map[string]any{
			"200": responseJSON("Sessão criada.", "#/components/schemas/SessionView", map[string]any{
				"id":            "d8c0c6624877463bb3027cbcb0ea2197",
				"credential_id": "a0f6f5c4f2a44511b7af35d6b2e9f893",
				"created_at":    "2026-04-18T13:11:00Z",
			}),
		}, true, true)
	case http.MethodGet + " /v1/credentials/{id}/discover":
		op["description"] = "Resolve sessão e retorna dados de perfil + UCs disponíveis no portal."
		op["parameters"] = []any{pathParam("id", "ID da credencial.")}
		op["responses"] = withDefaultErrors(map[string]any{
			"200": responseJSON("Perfil e UCs descobertos.", "#/components/schemas/DiscoveryResult", map[string]any{
				"minha_conta": map[string]any{
					"nome":                "PAULA",
					"usuarioAcesso":       "PAULA",
					"email":               "a@b.com",
					"celular":             "(74) 99999-9999",
					"dtUltimaAtualizacao": "2026-04-16T22:00:00-03:00",
				},
				"ucs": []any{
					map[string]any{
						"status":      "LIGADA",
						"uc":          "007098175908",
						"nomeCliente": "PAULA",
						"instalacao":  "0001",
						"grupoTensao": "B",
						"contrato":    "0023",
						"dt_inicio":   "2026-02-20",
						"dt_fim":      "9999-12-31",
						"local": map[string]any{
							"endereco":  "Rua X",
							"bairro":    "Centro",
							"municipio": "Lapão",
							"cep":       "44905-000",
							"uf":        "BA",
						},
					},
				},
			}),
		}, true, true)
	case http.MethodPost + " /v1/sync/uc":
		op["description"] = "Sincroniza uma UC e persiste resultados. Pode usar `credential_id` ou `bearer_token`+`documento`."
		op["requestBody"] = map[string]any{
			"required": true,
			"content": map[string]any{
				"application/json": map[string]any{
					"schema": refSchema("#/components/schemas/SyncUCRequest"),
					"examples": map[string]any{
						"credential": map[string]any{
							"summary": "Sincronização com credencial salva",
							"value": map[string]any{
								"credential_id":      "a0f6f5c4f2a44511b7af35d6b2e9f893",
								"uc":                 "007098175908",
								"include_pdf":        true,
								"include_extraction": true,
							},
						},
						"manual_token": map[string]any{
							"summary": "Sincronização com token manual",
							"value": map[string]any{
								"bearer_token": "eyJhbGciOi...",
								"documento":    "03021937586",
								"uc":           "007098175908",
								"include_pdf":  true,
							},
						},
					},
				},
			},
		}
		op["responses"] = withDefaultErrors(map[string]any{
			"200": responseJSON("Sincronização executada.", "#/components/schemas/SyncUCResponse", syncResponseExample()),
		}, false, true)
	case http.MethodGet + " /v1/consumer-units":
		op["description"] = "Lista UCs persistidas."
		op["parameters"] = []any{
			queryParamInt("limit", "Limite de itens.", 100),
			queryParamString("status", "Filtro por status da UC (ex.: LIGADA)."),
		}
		op["responses"] = withDefaultErrors(map[string]any{
			"200": responseJSON("Lista de UCs.", "#/components/schemas/ConsumerUnitListResponse", map[string]any{
				"items": []any{
					map[string]any{
						"uc":            "007098175908",
						"status":        "LIGADA",
						"nome_cliente":  "PAULA",
						"instalacao":    "0001",
						"contrato":      "0023",
						"grupo_tensao":  "B",
						"credential_id": "a0f6f5c4f2a44511b7af35d6b2e9f893",
						"created_at":    "2026-04-18T13:20:00Z",
						"updated_at":    "2026-04-18T13:20:00Z",
					},
				},
				"limit":  100,
				"status": "LIGADA",
			}),
		}, false, true)
	case http.MethodGet + " /v1/consumer-units/{uc}":
		op["description"] = "Consulta uma UC específica com ponteiros para última fatura e último sync."
		op["parameters"] = []any{pathParam("uc", "Código da unidade consumidora.")}
		op["responses"] = withDefaultErrors(map[string]any{
			"200": responseJSON("Detalhes da UC.", "#/components/schemas/ConsumerUnitDetailsView", map[string]any{
				"uc":           "007098175908",
				"status":       "LIGADA",
				"nome_cliente": "PAULA",
				"latest_invoice": map[string]any{
					"id":             "4c1a4c03a38b404987343123bce75ea8",
					"uc":             "007098175908",
					"numero_fatura":  "339800707843",
					"mes_referencia": "2026/04",
				},
				"latest_sync_run": map[string]any{
					"id":         "f4f896f8f4f64f7a935b8762ff9df3a8",
					"documento":  "03021937586",
					"uc":         "007098175908",
					"status":     "succeeded",
					"started_at": "2026-04-18T13:20:00Z",
				},
			}),
		}, true, true)
	case http.MethodGet + " /v1/consumer-units/{uc}/invoices":
		op["description"] = "Lista faturas de uma UC com filtros."
		op["parameters"] = []any{
			pathParam("uc", "Código da unidade consumidora."),
			queryParamInt("limit", "Limite de itens.", 100),
			queryParamString("status", "Filtro por status da fatura (ex.: A Vencer)."),
		}
		op["responses"] = withDefaultErrors(map[string]any{
			"200": responseJSON("Lista de faturas da UC.", "#/components/schemas/InvoiceListResponse", map[string]any{
				"uc": "007098175908",
				"items": []any{
					map[string]any{
						"id":                  "4c1a4c03a38b404987343123bce75ea8",
						"uc":                  "007098175908",
						"numero_fatura":       "339800707843",
						"mes_referencia":      "2026/04",
						"status_fatura":       "A Vencer",
						"valor_total":         "521.53",
						"completeness_status": "complete",
						"created_at":          "2026-04-18T13:20:00Z",
						"updated_at":          "2026-04-18T13:20:00Z",
					},
				},
				"limit":  100,
				"status": "A Vencer",
			}),
		}, false, true)
	case http.MethodGet + " /v1/consumer-units/{uc}/latest-invoice":
		op["description"] = "Retorna a fatura mais recente de uma UC."
		op["parameters"] = []any{pathParam("uc", "Código da unidade consumidora.")}
		op["responses"] = withDefaultErrors(map[string]any{
			"200": responseJSON("Última fatura da UC.", "#/components/schemas/InvoiceView", map[string]any{
				"id":                  "4c1a4c03a38b404987343123bce75ea8",
				"uc":                  "007098175908",
				"numero_fatura":       "339800707843",
				"mes_referencia":      "2026/04",
				"status_fatura":       "A Vencer",
				"valor_total":         "521.53",
				"completeness_status": "complete",
				"created_at":          "2026-04-18T13:20:00Z",
				"updated_at":          "2026-04-18T13:20:00Z",
			}),
		}, true, true)
	case http.MethodPost + " /v1/consumer-units/{uc}/sync":
		op["description"] = "Sincroniza uma UC específica definida no path."
		op["parameters"] = []any{pathParam("uc", "Código da unidade consumidora.")}
		op["requestBody"] = map[string]any{
			"required":    false,
			"description": "Mesmo payload de /v1/sync/uc sem o campo `uc` (ele vem do path).",
			"content": map[string]any{
				"application/json": map[string]any{
					"schema": refSchema("#/components/schemas/SyncUCPathRequest"),
				},
			},
		}
		op["responses"] = withDefaultErrors(map[string]any{
			"200": responseJSON("Sincronização executada.", "#/components/schemas/SyncUCResponse", syncResponseExample()),
		}, false, true)
	case http.MethodGet + " /v1/invoices/{id}":
		op["description"] = "Consulta uma fatura persistida por ID."
		op["parameters"] = []any{pathParam("id", "ID da fatura.")}
		op["responses"] = withDefaultErrors(map[string]any{
			"200": responseJSON("Fatura persistida.", "#/components/schemas/InvoiceView", map[string]any{
				"id":                  "4c1a4c03a38b404987343123bce75ea8",
				"uc":                  "007098175908",
				"numero_fatura":       "339800707843",
				"mes_referencia":      "2026/04",
				"status_fatura":       "A Vencer",
				"valor_total":         "521.53",
				"completeness_status": "complete",
				"created_at":          "2026-04-18T13:20:00Z",
				"updated_at":          "2026-04-18T13:20:00Z",
			}),
		}, true, true)
	case http.MethodGet + " /v1/sync-runs/{id}":
		op["description"] = "Consulta execução de sincronização por ID."
		op["parameters"] = []any{pathParam("id", "ID da execução de sync.")}
		op["responses"] = withDefaultErrors(map[string]any{
			"200": responseJSON("Sync run persistido.", "#/components/schemas/SyncRunView", map[string]any{
				"id":            "f4f896f8f4f64f7a935b8762ff9df3a8",
				"credential_id": "a0f6f5c4f2a44511b7af35d6b2e9f893",
				"documento":     "03021937586",
				"uc":            "007098175908",
				"status":        "succeeded",
				"started_at":    "2026-04-18T13:20:00Z",
				"finished_at":   "2026-04-18T13:21:00Z",
			}),
		}, true, true)
	case http.MethodPost + " /v1/billing/contracts":
		op["description"] = "Cria um novo contrato da UC e encerra automaticamente o ativo anterior."
		op["requestBody"] = map[string]any{
			"required": true,
			"content": map[string]any{
				"application/json": map[string]any{
					"schema": refSchema("#/components/schemas/BillingContractCreateRequest"),
					"examples": map[string]any{
						"default": map[string]any{
							"value": map[string]any{
								"customer_id":                          "d2bc53b4-034c-4e76-a8d8-e2b4f438b9cf",
								"consumer_unit_id":                     "0e446b0c-bfa1-4f42-8a9c-f92eac68538d",
								"vigencia_inicio":                      "2026-03-01",
								"desconto_percentual":                  "0.30",
								"ip_faturamento_mode":                  "fixed",
								"ip_faturamento_valor":                 "45.00",
								"bandeira_com_desconto":                false,
								"custo_disponibilidade_sempre_cobrado": true,
								"notes":                                "Contrato inicial",
								"created_by":                           "9ebd13f5-b7e0-4235-9f47-4b1df50e80cb",
							},
						},
					},
				},
			},
		}
		op["responses"] = withDefaultErrors(map[string]any{
			"201": responseJSON("Contrato criado.", "#/components/schemas/BillingContractView", map[string]any{
				"id":                                   "be7d933f-5607-4f70-90f3-31729d3ecb2f",
				"customer_id":                          "d2bc53b4-034c-4e76-a8d8-e2b4f438b9cf",
				"consumer_unit_id":                     "0e446b0c-bfa1-4f42-8a9c-f92eac68538d",
				"vigencia_inicio":                      "2026-03-01",
				"desconto_percentual":                  "0.30",
				"ip_faturamento_mode":                  "fixed",
				"ip_faturamento_valor":                 "45",
				"ip_faturamento_percent":               "0",
				"bandeira_com_desconto":                false,
				"custo_disponibilidade_sempre_cobrado": true,
				"status":                               "active",
				"created_at":                           "2026-04-19T13:40:00Z",
				"updated_at":                           "2026-04-19T13:40:00Z",
			}),
		}, false, true)
	case http.MethodGet + " /v1/billing/contracts/{id}":
		op["description"] = "Busca contrato por ID."
		op["parameters"] = []any{pathParam("id", "ID do contrato.")}
		op["responses"] = withDefaultErrors(map[string]any{
			"200": responseJSON("Contrato encontrado.", "#/components/schemas/BillingContractView", map[string]any{
				"id":                                   "be7d933f-5607-4f70-90f3-31729d3ecb2f",
				"customer_id":                          "d2bc53b4-034c-4e76-a8d8-e2b4f438b9cf",
				"consumer_unit_id":                     "0e446b0c-bfa1-4f42-8a9c-f92eac68538d",
				"vigencia_inicio":                      "2026-03-01",
				"desconto_percentual":                  "0.30",
				"ip_faturamento_mode":                  "fixed",
				"ip_faturamento_valor":                 "45",
				"ip_faturamento_percent":               "0",
				"bandeira_com_desconto":                false,
				"custo_disponibilidade_sempre_cobrado": true,
				"status":                               "active",
				"created_at":                           "2026-04-19T13:40:00Z",
				"updated_at":                           "2026-04-19T13:40:00Z",
			}),
		}, true, true)
	case http.MethodGet + " /v1/billing/consumer-units/{uc_id}/active-contract":
		op["description"] = "Retorna o contrato ativo da UC no momento da consulta."
		op["parameters"] = []any{pathParam("uc_id", "ID da unidade consumidora (UUID).")}
		op["responses"] = withDefaultErrors(map[string]any{
			"200": responseJSON("Contrato ativo da UC.", "#/components/schemas/BillingContractView", map[string]any{
				"id":                                   "be7d933f-5607-4f70-90f3-31729d3ecb2f",
				"customer_id":                          "d2bc53b4-034c-4e76-a8d8-e2b4f438b9cf",
				"consumer_unit_id":                     "0e446b0c-bfa1-4f42-8a9c-f92eac68538d",
				"vigencia_inicio":                      "2026-03-01",
				"desconto_percentual":                  "0.30",
				"ip_faturamento_mode":                  "fixed",
				"ip_faturamento_valor":                 "45",
				"ip_faturamento_percent":               "0",
				"bandeira_com_desconto":                false,
				"custo_disponibilidade_sempre_cobrado": true,
				"status":                               "active",
				"created_at":                           "2026-04-19T13:40:00Z",
				"updated_at":                           "2026-04-19T13:40:00Z",
			}),
		}, true, true)
	case http.MethodGet + " /v1/billing/calculations/{id}":
		op["description"] = "Retorna um cálculo de faturamento persistido."
		op["parameters"] = []any{pathParam("id", "ID do cálculo.")}
		op["responses"] = withDefaultErrors(map[string]any{
			"200": responseJSON("Cálculo encontrado.", "#/components/schemas/BillingCalculationView", map[string]any{
				"id":                     "a7f810e8-b4dc-425e-a64d-087f5e6a9273",
				"utility_invoice_ref_id": "94948c16-9f47-4de8-86dd-3f1d4de59176",
				"billing_cycle_id":       "62f931a9-e8fa-4a28-80fc-b4e48125768a",
				"consumer_unit_id":       "0e446b0c-bfa1-4f42-8a9c-f92eac68538d",
				"contract_id":            "be7d933f-5607-4f70-90f3-31729d3ecb2f",
				"total_sem_desconto":     "521.53",
				"total_com_desconto":     "389.11",
				"economia_rs":            "132.42",
				"economia_pct":           "0.2539",
				"status":                 "draft",
				"version":                1,
				"calculated_at":          "2026-04-19T13:40:00Z",
			}),
		}, true, true)
	case http.MethodGet + " /v1/extractor/contracts":
		op["description"] = "Informa os caminhos dos contratos JSON Schema usados na integração com o extrator."
		op["responses"] = withDefaultErrors(map[string]any{
			"200": responseJSON("Contratos do extrator.", "#/components/schemas/ExtractorContractsResponse", map[string]any{
				"extractor_request":  "packages/contracts/extractor-request.schema.json",
				"extractor_response": "packages/contracts/extractor-response.schema.json",
			}),
		}, false, true)
	default:
		return nil
	}

	return op
}

func componentsSpec() map[string]any {
	return map[string]any{
		"securitySchemes": map[string]any{
			"ApiKeyAuth": map[string]any{
				"type":        "apiKey",
				"in":          "header",
				"name":        "X-API-Key",
				"description": "Obrigatório quando BACKEND_API_KEY estiver configurado.",
			},
		},
		"responses": map[string]any{
			"BadRequest":          errorResponseRef("Requisição inválida.", "invalid_json"),
			"Unauthorized":        errorResponseRef("Não autenticado.", "unauthorized"),
			"NotFound":            errorResponseRef("Recurso não encontrado.", "not_found"),
			"MethodNotAllowed":    errorResponseRef("Método não suportado para o endpoint.", "method_not_allowed"),
			"InternalServerError": errorResponseRef("Erro interno.", "internal_error"),
		},
		"schemas": map[string]any{
			"ErrorResponse": map[string]any{
				"type":       "object",
				"required":   []string{"error"},
				"properties": map[string]any{"error": map[string]any{"type": "string"}},
			},
			"HealthResponse": map[string]any{
				"type":     "object",
				"required": []string{"status", "service"},
				"properties": map[string]any{
					"status":  map[string]any{"type": "string", "example": "ok"},
					"service": map[string]any{"type": "string", "example": "backend-go"},
				},
			},
			"CredentialCreateRequest": map[string]any{
				"type":     "object",
				"required": []string{"documento", "senha"},
				"properties": map[string]any{
					"label":       map[string]any{"type": "string", "description": "Nome amigável da credencial."},
					"documento":   map[string]any{"type": "string", "description": "CPF/CNPJ usado no login Neoenergia."},
					"senha":       map[string]any{"type": "string", "format": "password", "description": "Senha do portal Neoenergia."},
					"uf":          map[string]any{"type": "string", "default": "BA"},
					"tipo_acesso": map[string]any{"type": "string", "default": "normal"},
				},
			},
			"CredentialView": map[string]any{
				"type":     "object",
				"required": []string{"id", "label", "documento", "uf", "tipo_acesso", "created_at"},
				"properties": map[string]any{
					"id":          map[string]any{"type": "string"},
					"label":       map[string]any{"type": "string"},
					"documento":   map[string]any{"type": "string", "description": "Documento mascarado."},
					"uf":          map[string]any{"type": "string"},
					"tipo_acesso": map[string]any{"type": "string"},
					"created_at":  map[string]any{"type": "string", "format": "date-time"},
				},
			},
			"SessionView": map[string]any{
				"type":     "object",
				"required": []string{"id", "credential_id", "created_at"},
				"properties": map[string]any{
					"id":            map[string]any{"type": "string"},
					"credential_id": map[string]any{"type": "string"},
					"created_at":    map[string]any{"type": "string", "format": "date-time"},
				},
			},
			"DiscoveryResult": map[string]any{
				"type":     "object",
				"required": []string{"ucs"},
				"properties": map[string]any{
					"minha_conta":        map[string]any{"$ref": "#/components/schemas/MinhaContaResponse"},
					"minha_conta_legado": map[string]any{"$ref": "#/components/schemas/MinhaContaLegadoResponse"},
					"ucs": map[string]any{
						"type":  "array",
						"items": refSchema("#/components/schemas/UC"),
					},
					"errors": map[string]any{
						"type":                 "object",
						"additionalProperties": map[string]any{"type": "string"},
					},
				},
			},
			"SyncUCRequest": map[string]any{
				"type":     "object",
				"required": []string{"uc"},
				"properties": map[string]any{
					"bearer_token":       map[string]any{"type": "string", "description": "Token Bearer Neoenergia (opcional quando credential_id for enviado)."},
					"credential_id":      map[string]any{"type": "string", "description": "Credencial salva para resolver sessão automaticamente."},
					"documento":          map[string]any{"type": "string", "description": "Obrigatório quando bearer_token é enviado manualmente."},
					"uc":                 map[string]any{"type": "string", "description": "Unidade consumidora alvo."},
					"include_pdf":        map[string]any{"type": "boolean", "default": false},
					"include_extraction": map[string]any{"type": "boolean", "default": false},
				},
			},
			"SyncUCPathRequest": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"bearer_token":       map[string]any{"type": "string"},
					"credential_id":      map[string]any{"type": "string"},
					"documento":          map[string]any{"type": "string"},
					"include_pdf":        map[string]any{"type": "boolean", "default": false},
					"include_extraction": map[string]any{"type": "boolean", "default": false},
				},
			},
			"BillingContractCreateRequest": map[string]any{
				"type":     "object",
				"required": []string{"customer_id", "consumer_unit_id", "vigencia_inicio", "desconto_percentual"},
				"properties": map[string]any{
					"customer_id":                          map[string]any{"type": "string", "format": "uuid"},
					"consumer_unit_id":                     map[string]any{"type": "string", "format": "uuid"},
					"vigencia_inicio":                      map[string]any{"type": "string", "format": "date"},
					"desconto_percentual":                  map[string]any{"type": "string", "description": "Decimal no intervalo (0,1], ex.: 0.30"},
					"ip_faturamento_mode":                  map[string]any{"type": "string", "enum": []any{"fixed", "percent"}, "default": "fixed"},
					"ip_faturamento_valor":                 map[string]any{"type": "string"},
					"ip_faturamento_percent":               map[string]any{"type": "string"},
					"bandeira_com_desconto":                map[string]any{"type": "boolean", "default": false},
					"custo_disponibilidade_sempre_cobrado": map[string]any{"type": "boolean", "default": true},
					"notes":                                map[string]any{"type": "string"},
					"created_by":                           map[string]any{"type": "string", "format": "uuid"},
				},
			},
			"BillingContractView": map[string]any{
				"type":     "object",
				"required": []string{"id", "customer_id", "consumer_unit_id", "vigencia_inicio", "desconto_percentual", "ip_faturamento_mode", "status", "created_at", "updated_at"},
				"properties": map[string]any{
					"id":                                   map[string]any{"type": "string", "format": "uuid"},
					"customer_id":                          map[string]any{"type": "string", "format": "uuid"},
					"consumer_unit_id":                     map[string]any{"type": "string", "format": "uuid"},
					"vigencia_inicio":                      map[string]any{"type": "string", "format": "date"},
					"vigencia_fim":                         map[string]any{"type": "string", "format": "date"},
					"desconto_percentual":                  map[string]any{"type": "string"},
					"ip_faturamento_mode":                  map[string]any{"type": "string"},
					"ip_faturamento_valor":                 map[string]any{"type": "string"},
					"ip_faturamento_percent":               map[string]any{"type": "string"},
					"bandeira_com_desconto":                map[string]any{"type": "boolean"},
					"custo_disponibilidade_sempre_cobrado": map[string]any{"type": "boolean"},
					"notes":                                map[string]any{"type": "string"},
					"status":                               map[string]any{"type": "string", "enum": []any{"draft", "active", "ended"}},
					"created_by":                           map[string]any{"type": "string", "format": "uuid"},
					"created_at":                           map[string]any{"type": "string", "format": "date-time"},
					"updated_at":                           map[string]any{"type": "string", "format": "date-time"},
				},
			},
			"BillingCalculationView": map[string]any{
				"type":     "object",
				"required": []string{"id", "utility_invoice_ref_id", "billing_cycle_id", "consumer_unit_id", "contract_id", "total_sem_desconto", "total_com_desconto", "economia_rs", "economia_pct", "status", "version", "calculated_at"},
				"properties": map[string]any{
					"id":                     map[string]any{"type": "string", "format": "uuid"},
					"utility_invoice_ref_id": map[string]any{"type": "string", "format": "uuid"},
					"billing_cycle_id":       map[string]any{"type": "string", "format": "uuid"},
					"consumer_unit_id":       map[string]any{"type": "string", "format": "uuid"},
					"contract_id":            map[string]any{"type": "string", "format": "uuid"},
					"total_sem_desconto":     map[string]any{"type": "string"},
					"total_com_desconto":     map[string]any{"type": "string"},
					"economia_rs":            map[string]any{"type": "string"},
					"economia_pct":           map[string]any{"type": "string"},
					"status":                 map[string]any{"type": "string", "enum": []any{"draft", "needs_review", "approved", "superseded"}},
					"version":                map[string]any{"type": "integer"},
					"calculated_at":          map[string]any{"type": "string", "format": "date-time"},
					"needs_review_reasons": map[string]any{
						"type":  "array",
						"items": map[string]any{"type": "string"},
					},
					"approved_at":       map[string]any{"type": "string", "format": "date-time"},
					"approved_by":       map[string]any{"type": "string", "format": "uuid"},
					"contract_snapshot": map[string]any{"type": "object", "additionalProperties": true},
					"inputs_snapshot":   map[string]any{"type": "object", "additionalProperties": true},
					"result_snapshot":   map[string]any{"type": "object", "additionalProperties": true},
				},
			},
			"SyncUCResponse": map[string]any{
				"type":     "object",
				"required": []string{"documento", "uc"},
				"properties": map[string]any{
					"documento":                       map[string]any{"type": "string"},
					"uc":                              map[string]any{"type": "string"},
					"grupo_cliente":                   refSchema("#/components/schemas/OptionalGrupoClienteResult"),
					"minha_conta":                     refSchema("#/components/schemas/OptionalMinhaContaResult"),
					"minha_conta_legado":              refSchema("#/components/schemas/OptionalMinhaContaLegadoResult"),
					"ucs":                             refSchema("#/components/schemas/OptionalUCsResult"),
					"imovel":                          refSchema("#/components/schemas/OptionalImovelResult"),
					"protocolo":                       refSchema("#/components/schemas/OptionalProtocoloResult"),
					"historico_consumo":               refSchema("#/components/schemas/OptionalHistoricoConsumoResult"),
					"data_certa":                      refSchema("#/components/schemas/OptionalDataCertaResult"),
					"fatura_digital":                  refSchema("#/components/schemas/OptionalFaturaDigitalResult"),
					"debito_automatico":               refSchema("#/components/schemas/OptionalDebitoAutomaticoResult"),
					"motivos_segunda_via":             refSchema("#/components/schemas/OptionalMotivosSegundaViaResult"),
					"faturas":                         refSchema("#/components/schemas/OptionalFaturasResult"),
					"dados_pagamento_primeira_fatura": refSchema("#/components/schemas/OptionalDadosPagamentoResult"),
					"pdf_primeira_fatura":             refSchema("#/components/schemas/OptionalFaturaPDFResult"),
					"extraction":                      refSchema("#/components/schemas/OptionalExtractorResponseResult"),
					"billing_record":                  refSchema("#/components/schemas/BillingRecord"),
					"document_record":                 refSchema("#/components/schemas/DocumentRecord"),
					"persistence":                     refSchema("#/components/schemas/PersistenceResult"),
				},
			},
			"OptionalGrupoClienteResult":      optionalResultSchema("#/components/schemas/GrupoClienteResponse"),
			"OptionalMinhaContaResult":        optionalResultSchema("#/components/schemas/MinhaContaResponse"),
			"OptionalMinhaContaLegadoResult":  optionalResultSchema("#/components/schemas/MinhaContaLegadoResponse"),
			"OptionalUCsResult":               optionalResultSchema("#/components/schemas/UCsResponse"),
			"OptionalImovelResult":            optionalResultSchema("#/components/schemas/ImovelResponse"),
			"OptionalProtocoloResult":         optionalResultSchema("#/components/schemas/ProtocoloResponse"),
			"OptionalHistoricoConsumoResult":  optionalResultSchema("#/components/schemas/HistoricoConsumoResponse"),
			"OptionalDataCertaResult":         optionalResultSchema("#/components/schemas/DataCertaResponse"),
			"OptionalFaturaDigitalResult":     optionalResultSchema("#/components/schemas/FaturaDigitalResponse"),
			"OptionalDebitoAutomaticoResult":  optionalResultSchema("#/components/schemas/DebitoAutomaticoResponse"),
			"OptionalMotivosSegundaViaResult": optionalResultSchema("#/components/schemas/MotivosSegundaViaResponse"),
			"OptionalFaturasResult":           optionalResultSchema("#/components/schemas/FaturasResponse"),
			"OptionalDadosPagamentoResult":    optionalResultSchema("#/components/schemas/DadosPagamentoResponse"),
			"OptionalFaturaPDFResult":         optionalResultSchema("#/components/schemas/FaturaPDFResponse"),
			"OptionalExtractorResponseResult": optionalResultSchema("#/components/schemas/ExtractorResponse"),
			"NeoenergiaErrorResponse": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"status_code": map[string]any{"type": "integer"},
					"method":      map[string]any{"type": "string"},
					"path":        map[string]any{"type": "string"},
					"body":        map[string]any{"type": "object", "additionalProperties": true},
				},
			},
			"MinhaContaResponse": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"nome":                map[string]any{"type": "string"},
					"usuarioAcesso":       map[string]any{"type": "string"},
					"email":               map[string]any{"type": "string"},
					"celular":             map[string]any{"type": "string"},
					"dtUltimaAtualizacao": map[string]any{"type": "string"},
				},
			},
			"GrupoClienteResponse": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"dados": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"codigo": map[string]any{"type": "integer"},
								"nome":   map[string]any{"type": "string"},
								"ativo":  map[string]any{"type": "boolean"},
								"funcionalidades": map[string]any{
									"type": "array",
									"items": map[string]any{
										"type": "object",
										"properties": map[string]any{
											"codigo": map[string]any{"type": "string"},
										},
									},
								},
							},
						},
					},
				},
			},
			"MinhaContaLegadoResponse": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"nomeTitular":     map[string]any{"type": "string"},
					"dtNascimento":    map[string]any{"type": "string"},
					"emailCadastro":   map[string]any{"type": "string"},
					"telefoneContato": map[string]any{"type": "string"},
				},
			},
			"UC": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"status":      map[string]any{"type": "string"},
					"uc":          map[string]any{"type": "string"},
					"nomeCliente": map[string]any{"type": "string"},
					"instalacao":  map[string]any{"type": "string"},
					"grupoTensao": map[string]any{"type": "string"},
					"contrato":    map[string]any{"type": "string"},
					"dt_inicio":   map[string]any{"type": "string"},
					"dt_fim":      map[string]any{"type": "string"},
					"local": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"endereco":  map[string]any{"type": "string"},
							"bairro":    map[string]any{"type": "string"},
							"municipio": map[string]any{"type": "string"},
							"cep":       map[string]any{"type": "string"},
							"uf":        map[string]any{"type": "string"},
						},
					},
				},
			},
			"UCsResponse": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"ucs": map[string]any{
						"type":  "array",
						"items": refSchema("#/components/schemas/UC"),
					},
				},
			},
			"ImovelResponse": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"codigo":      map[string]any{"type": "string"},
					"instalacao":  map[string]any{"type": "string"},
					"medidor":     map[string]any{"type": "string"},
					"fase":        map[string]any{"type": "string"},
					"dataLigacao": map[string]any{"type": "string"},
					"situacao": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"codigo":         map[string]any{"type": "string"},
							"descricao":      map[string]any{"type": "string"},
							"dataSituacaoUC": map[string]any{"type": "string"},
						},
					},
					"cliente": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"codigo": map[string]any{"type": "string"},
							"nome":   map[string]any{"type": "string"},
						},
					},
				},
			},
			"ProtocoloResponse": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"protocoloSalesforceStr": map[string]any{"type": "string"},
					"protocoloLegadoStr":     map[string]any{"type": "string"},
				},
			},
			"HistoricoConsumoItem": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"dataPagamento":         map[string]any{"type": "string"},
					"dataVencimento":        map[string]any{"type": "string"},
					"dataLeitura":           map[string]any{"type": "string"},
					"consumoKw":             map[string]any{"type": "string"},
					"mesReferencia":         map[string]any{"type": "string"},
					"numeroLeitura":         map[string]any{"type": "string"},
					"tipoLeitura":           map[string]any{"type": "string"},
					"dataInicioPeriodoCalc": map[string]any{"type": "string"},
					"dataFimPeriodoCalc":    map[string]any{"type": "string"},
					"dataProxLeitura":       map[string]any{"type": "string"},
					"valorFatura":           map[string]any{"type": "string"},
					"statusFatura":          map[string]any{"type": "string"},
					"numeroFatura":          map[string]any{"type": "string"},
				},
			},
			"HistoricoConsumoResponse": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"historicoConsumo": map[string]any{
						"type":  "array",
						"items": refSchema("#/components/schemas/HistoricoConsumoItem"),
					},
					"mediamensal": map[string]any{"type": "string"},
				},
			},
			"DataCertaResponse": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"possuiDataBoa": map[string]any{"type": "string"},
					"dataAtual":     map[string]any{"type": "string"},
				},
			},
			"FaturaDigitalResponse": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"PossuiFaturaDigital": map[string]any{"type": "string"},
					"emailFatura":         map[string]any{"type": "string"},
					"emailCadastro":       map[string]any{"type": "string"},
				},
			},
			"DebitoAutomaticoResponse": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"retorno": map[string]any{
						"type":                 "object",
						"additionalProperties": true,
					},
				},
			},
			"MotivosSegundaViaResponse": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"motivos": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"idMotivo":  map[string]any{"type": "string"},
								"descricao": map[string]any{"type": "string"},
							},
						},
					},
				},
			},
			"Fatura": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"numeroFatura":      map[string]any{"type": "string"},
					"mesReferencia":     map[string]any{"type": "string"},
					"statusFatura":      map[string]any{"type": "string"},
					"dataCompetencia":   map[string]any{"type": "string"},
					"dataEmissao":       map[string]any{"type": "string"},
					"dataPagamento":     map[string]any{"type": "string"},
					"dataVencimento":    map[string]any{"type": "string"},
					"dataInicioPeriodo": map[string]any{"type": "string"},
					"dataFimPeriodo":    map[string]any{"type": "string"},
					"valorEmissao":      map[string]any{"type": "string"},
					"uc":                map[string]any{"type": "string"},
					"tipoDoc":           map[string]any{"type": "string"},
					"origemFatura":      map[string]any{"type": "string"},
					"tipoEntrega":       map[string]any{"type": "string"},
					"tipoLeitura":       map[string]any{"type": "string"},
					"aceitaPix":         map[string]any{"type": "string"},
				},
			},
			"FaturasResponse": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"entregaFaturas": map[string]any{"type": "object", "additionalProperties": true},
					"faturas": map[string]any{
						"type":  "array",
						"items": refSchema("#/components/schemas/Fatura"),
					},
				},
			},
			"DadosPagamentoResponse": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"codBarras": map[string]any{"type": "string"},
				},
			},
			"FaturaPDFResponse": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"fileName":      map[string]any{"type": "string"},
					"fileSize":      map[string]any{"type": "string"},
					"fileData":      map[string]any{"type": "string"},
					"fileExtension": map[string]any{"type": "string"},
				},
			},
			"ExtractorResponse": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"schema_version": map[string]any{"type": "string"},
					"job_id":         map[string]any{"type": "string"},
					"status":         map[string]any{"type": "string"},
					"fields":         map[string]any{"type": "object", "additionalProperties": true},
					"source_map":     map[string]any{"type": "object", "additionalProperties": map[string]any{"type": "string"}},
					"confidence_map": map[string]any{"type": "object", "additionalProperties": map[string]any{"type": "number"}},
					"warnings":       map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
					"artifacts":      map[string]any{"type": "object", "additionalProperties": true},
				},
			},
			"BillingRecord": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"uc":                      map[string]any{"type": "string"},
					"numero_fatura":           map[string]any{"type": "string"},
					"mes_referencia":          map[string]any{"type": "string"},
					"status_fatura":           map[string]any{"type": "string"},
					"valor_total":             map[string]any{"type": "string"},
					"codigo_barras":           map[string]any{"type": "string"},
					"data_emissao":            map[string]any{"type": "string"},
					"data_vencimento":         map[string]any{"type": "string"},
					"data_pagamento":          map[string]any{"type": "string"},
					"data_inicio_periodo":     map[string]any{"type": "string"},
					"data_fim_periodo":        map[string]any{"type": "string"},
					"historico_consumo":       map[string]any{"type": "array", "items": map[string]any{"type": "object", "additionalProperties": true}},
					"itens_fatura":            map[string]any{"type": "array", "items": map[string]any{"type": "object", "additionalProperties": true}},
					"composicao_fornecimento": map[string]any{"type": "array", "items": map[string]any{"type": "object", "additionalProperties": true}},
					"nota_fiscal":             map[string]any{"type": "object", "additionalProperties": true},
					"extractor_status":        map[string]any{"type": "string"},
					"completeness":            refSchema("#/components/schemas/BillingCompleteness"),
					"source_map":              map[string]any{"type": "object", "additionalProperties": map[string]any{"type": "string"}},
					"confidence_map":          map[string]any{"type": "object", "additionalProperties": map[string]any{"type": "number"}},
					"raw_source_availability": map[string]any{"type": "object", "additionalProperties": map[string]any{"type": "boolean"}},
				},
			},
			"BillingCompleteness": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"status":         map[string]any{"type": "string", "enum": []string{"complete", "partial"}},
					"missing_fields": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
				},
			},
			"DocumentRecord": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"data_vencimento":   map[string]any{"type": "string"},
					"nome":              map[string]any{"type": "string"},
					"normalizado_valor": map[string]any{"oneOf": []any{map[string]any{"type": "string"}, map[string]any{"type": "number"}}},
					"ocr":               map[string]any{"type": "object", "additionalProperties": true},
					"uc":                map[string]any{"type": "string"},
					"valor":             map[string]any{"type": "string"},
					"site_receipt":      map[string]any{"type": "string"},
				},
			},
			"PersistenceResult": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"sync_run_id": map[string]any{"type": "string"},
					"invoice_id":  map[string]any{"type": "string"},
					"status":      map[string]any{"type": "string"},
					"error":       map[string]any{"type": "string"},
				},
			},
			"ConsumerUnitView": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"uc":            map[string]any{"type": "string"},
					"credential_id": map[string]any{"type": "string"},
					"status":        map[string]any{"type": "string"},
					"nome_cliente":  map[string]any{"type": "string"},
					"instalacao":    map[string]any{"type": "string"},
					"contrato":      map[string]any{"type": "string"},
					"grupo_tensao":  map[string]any{"type": "string"},
					"endereco":      map[string]any{"type": "object", "additionalProperties": true},
					"imovel":        map[string]any{"type": "object", "additionalProperties": true},
					"created_at":    map[string]any{"type": "string", "format": "date-time"},
					"updated_at":    map[string]any{"type": "string", "format": "date-time"},
				},
			},
			"ConsumerUnitDetailsView": map[string]any{
				"allOf": []any{
					refSchema("#/components/schemas/ConsumerUnitView"),
					map[string]any{
						"type": "object",
						"properties": map[string]any{
							"latest_invoice":  refSchema("#/components/schemas/InvoiceView"),
							"latest_sync_run": refSchema("#/components/schemas/SyncRunView"),
						},
					},
				},
			},
			"ConsumerUnitListResponse": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"items": map[string]any{
						"type":  "array",
						"items": refSchema("#/components/schemas/ConsumerUnitView"),
					},
					"limit":  map[string]any{"type": "integer"},
					"status": map[string]any{"type": "string"},
				},
			},
			"InvoiceView": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id":                   map[string]any{"type": "string"},
					"uc":                   map[string]any{"type": "string"},
					"numero_fatura":        map[string]any{"type": "string"},
					"mes_referencia":       map[string]any{"type": "string"},
					"status_fatura":        map[string]any{"type": "string"},
					"valor_total":          map[string]any{"type": "string"},
					"codigo_barras":        map[string]any{"type": "string"},
					"data_emissao":         map[string]any{"type": "string"},
					"data_vencimento":      map[string]any{"type": "string"},
					"data_pagamento":       map[string]any{"type": "string"},
					"data_inicio_periodo":  map[string]any{"type": "string"},
					"data_fim_periodo":     map[string]any{"type": "string"},
					"completeness_status":  map[string]any{"type": "string"},
					"completeness_missing": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
					"billing_record":       map[string]any{"type": "object", "additionalProperties": true},
					"document_record":      map[string]any{"type": "object", "additionalProperties": true},
					"items":                map[string]any{"type": "array", "items": map[string]any{"type": "object", "additionalProperties": true}},
					"created_at":           map[string]any{"type": "string", "format": "date-time"},
					"updated_at":           map[string]any{"type": "string", "format": "date-time"},
				},
			},
			"InvoiceListResponse": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"uc": map[string]any{"type": "string"},
					"items": map[string]any{
						"type":  "array",
						"items": refSchema("#/components/schemas/InvoiceView"),
					},
					"limit":  map[string]any{"type": "integer"},
					"status": map[string]any{"type": "string"},
				},
			},
			"SyncRunView": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id":            map[string]any{"type": "string"},
					"credential_id": map[string]any{"type": "string"},
					"documento":     map[string]any{"type": "string"},
					"uc":            map[string]any{"type": "string"},
					"status":        map[string]any{"type": "string"},
					"started_at":    map[string]any{"type": "string", "format": "date-time"},
					"finished_at":   map[string]any{"type": "string", "format": "date-time"},
					"error_message": map[string]any{"type": "string"},
					"raw_response":  map[string]any{"type": "object", "additionalProperties": true},
				},
			},
			"ExtractorContractsResponse": map[string]any{
				"type":     "object",
				"required": []string{"extractor_request", "extractor_response"},
				"properties": map[string]any{
					"extractor_request":  map[string]any{"type": "string"},
					"extractor_response": map[string]any{"type": "string"},
				},
			},
		},
	}
}

func withDefaultErrors(responses map[string]any, includeNotFound bool, includeMethodNotAllowed bool) map[string]any {
	out := map[string]any{}
	for code, value := range responses {
		out[code] = value
	}
	out["400"] = map[string]any{"$ref": "#/components/responses/BadRequest"}
	out["401"] = map[string]any{"$ref": "#/components/responses/Unauthorized"}
	out["500"] = map[string]any{"$ref": "#/components/responses/InternalServerError"}
	if includeNotFound {
		out["404"] = map[string]any{"$ref": "#/components/responses/NotFound"}
	}
	if includeMethodNotAllowed {
		out["405"] = map[string]any{"$ref": "#/components/responses/MethodNotAllowed"}
	}
	return out
}

func responseJSON(description, schemaRef string, example any) map[string]any {
	content := map[string]any{
		"application/json": map[string]any{
			"schema": refSchema(schemaRef),
		},
	}
	if example != nil {
		content["application/json"].(map[string]any)["example"] = example
	}
	return map[string]any{
		"description": description,
		"content":     content,
	}
}

func responseRaw(description, mediaType string) map[string]any {
	return map[string]any{
		"description": description,
		"content": map[string]any{
			mediaType: map[string]any{
				"schema": map[string]any{"type": "string"},
			},
		},
	}
}

func errorResponseRef(description, exampleError string) map[string]any {
	return map[string]any{
		"description": description,
		"content": map[string]any{
			"application/json": map[string]any{
				"schema": refSchema("#/components/schemas/ErrorResponse"),
				"example": map[string]any{
					"error": exampleError,
				},
			},
		},
	}
}

func pathParam(name, description string) map[string]any {
	return map[string]any{
		"name":        name,
		"in":          "path",
		"required":    true,
		"description": description,
		"schema":      map[string]any{"type": "string"},
	}
}

func queryParamString(name, description string) map[string]any {
	return map[string]any{
		"name":        name,
		"in":          "query",
		"required":    false,
		"description": description,
		"schema":      map[string]any{"type": "string"},
	}
}

func queryParamInt(name, description string, defaultValue int) map[string]any {
	return map[string]any{
		"name":        name,
		"in":          "query",
		"required":    false,
		"description": description,
		"schema": map[string]any{
			"type":    "integer",
			"minimum": 1,
			"default": defaultValue,
		},
	}
}

func refSchema(path string) map[string]any {
	return map[string]any{"$ref": path}
}

func optionalResultSchema(dataSchemaRef string) map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"data":  refSchema(dataSchemaRef),
			"error": refSchema("#/components/schemas/NeoenergiaErrorResponse"),
		},
	}
}

func operationID(method, path string) string {
	clean := strings.Trim(path, "/")
	if clean == "" {
		clean = "root"
	}
	clean = strings.ReplaceAll(clean, "/", "_")
	clean = strings.ReplaceAll(clean, "{", "")
	clean = strings.ReplaceAll(clean, "}", "")
	clean = strings.ReplaceAll(clean, "-", "_")
	return strings.ToLower(method) + "_" + clean
}

func requiresAPIKey(path string) bool {
	return strings.HasPrefix(path, "/v1/")
}

func syncResponseExample() map[string]any {
	return map[string]any{
		"documento": "03021937586",
		"uc":        "007098175908",
		"faturas": map[string]any{
			"data": map[string]any{
				"faturas": []any{
					map[string]any{
						"numeroFatura":   "339800707843",
						"mesReferencia":  "2026/04",
						"statusFatura":   "A Vencer",
						"dataVencimento": "2026-05-06",
						"valorEmissao":   "521.53",
						"uc":             "007098175908",
					},
				},
			},
		},
		"dados_pagamento_primeira_fatura": map[string]any{
			"data": map[string]any{
				"codBarras": "838700000052215300300078098175908215056066285939",
			},
		},
		"pdf_primeira_fatura": map[string]any{
			"data": map[string]any{
				"fileName":      "007098175908",
				"fileExtension": ".pdf",
				"fileData":      "JVBERi0xLjQ=",
			},
		},
		"billing_record": map[string]any{
			"uc":                      "007098175908",
			"numero_fatura":           "339800707843",
			"mes_referencia":          "2026/04",
			"valor_total":             "521.53",
			"codigo_barras":           "838700000052215300300078098175908215056066285939",
			"completeness":            map[string]any{"status": "complete"},
			"source_map":              map[string]any{"valor_total": "api", "codigo_barras": "api"},
			"raw_source_availability": map[string]any{"api.fatura": true, "api.pdf": true},
		},
		"persistence": map[string]any{
			"sync_run_id": "f4f896f8f4f64f7a935b8762ff9df3a8",
			"invoice_id":  "4c1a4c03a38b404987343123bce75ea8",
			"status":      "succeeded",
		},
	}
}
