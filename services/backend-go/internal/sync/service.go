package sync

import (
	"context"
	"encoding/json"

	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/extractor"
	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/neoenergia"
	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/store"
)

type SyncUCRequest struct {
	BearerToken       string `json:"bearer_token"`
	CredentialID      string `json:"credential_id"`
	Documento         string `json:"documento"`
	UC                string `json:"uc"`
	IncludePDF        bool   `json:"include_pdf"`
	IncludeExtraction bool   `json:"include_extraction"`
}

type OptionalResult[T any] struct {
	Data  *T                        `json:"data,omitempty"`
	Error *neoenergia.ErrorResponse `json:"error,omitempty"`
}

type SyncUCResponse struct {
	Documento         string                                               `json:"documento"`
	UC                string                                               `json:"uc"`
	GrupoCliente      OptionalResult[neoenergia.GrupoClienteResponse]      `json:"grupo_cliente"`
	MinhaConta        OptionalResult[neoenergia.MinhaContaResponse]        `json:"minha_conta"`
	MinhaContaLegado  OptionalResult[neoenergia.MinhaContaLegadoResponse]  `json:"minha_conta_legado"`
	UCs               OptionalResult[neoenergia.UCsResponse]               `json:"ucs"`
	Imovel            OptionalResult[neoenergia.ImovelResponse]            `json:"imovel"`
	Protocolo         OptionalResult[neoenergia.ProtocoloResponse]         `json:"protocolo"`
	HistoricoConsumo  OptionalResult[neoenergia.HistoricoConsumoResponse]  `json:"historico_consumo"`
	DataCerta         OptionalResult[neoenergia.DataCertaResponse]         `json:"data_certa"`
	FaturaDigital     OptionalResult[neoenergia.FaturaDigitalResponse]     `json:"fatura_digital"`
	DebitoAutomatico  OptionalResult[neoenergia.DebitoAutomaticoResponse]  `json:"debito_automatico"`
	MotivosSegundaVia OptionalResult[neoenergia.MotivosSegundaViaResponse] `json:"motivos_segunda_via"`
	Faturas           OptionalResult[neoenergia.FaturasResponse]           `json:"faturas"`
	DadosPagamento    OptionalResult[neoenergia.DadosPagamentoResponse]    `json:"dados_pagamento_primeira_fatura"`
	PDF               OptionalResult[neoenergia.FaturaPDFResponse]         `json:"pdf_primeira_fatura"`
	Extraction        OptionalResult[extractor.Response]                   `json:"extraction"`
	BillingRecord     *BillingRecord                                       `json:"billing_record,omitempty"`
	DocumentRecord    *DocumentRecord                                      `json:"document_record,omitempty"`
	Persistence       *PersistenceResult                                   `json:"persistence,omitempty"`
}

type PersistenceResult struct {
	SyncRunID string `json:"sync_run_id"`
	InvoiceID string `json:"invoice_id,omitempty"`
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
}

type BillingRecord struct {
	UC                     string                        `json:"uc"`
	NumeroFatura           string                        `json:"numero_fatura"`
	MesReferencia          string                        `json:"mes_referencia"`
	StatusFatura           string                        `json:"status_fatura,omitempty"`
	ValorTotal             string                        `json:"valor_total"`
	CodigoBarras           *string                       `json:"codigo_barras"`
	DataEmissao            *string                       `json:"data_emissao"`
	DataVencimento         *string                       `json:"data_vencimento"`
	DataPagamento          *string                       `json:"data_pagamento"`
	DataInicioPeriodo      *string                       `json:"data_inicio_periodo,omitempty"`
	DataFimPeriodo         *string                       `json:"data_fim_periodo,omitempty"`
	HistoricoConsumo       []neoenergia.HistoricoConsumo `json:"historico_consumo,omitempty"`
	ItensFatura            any                           `json:"itens_fatura"`
	ComposicaoFornecimento any                           `json:"composicao_fornecimento"`
	NotaFiscal             any                           `json:"nota_fiscal"`
	ExtractorStatus        *string                       `json:"extractor_status,omitempty"`
	Completeness           BillingCompleteness           `json:"completeness"`
	SourceMap              map[string]string             `json:"source_map"`
	ConfidenceMap          map[string]float64            `json:"confidence_map,omitempty"`
	RawSourceAvailability  map[string]bool               `json:"raw_source_availability"`
}

type BillingCompleteness struct {
	Status        string   `json:"status"`
	MissingFields []string `json:"missing_fields,omitempty"`
}

type Service struct {
	client    *neoenergia.Client
	extractor *extractor.Client
	store     *store.SQLiteStore
}

func NewService(client *neoenergia.Client, extractorClient *extractor.Client, sqliteStore *store.SQLiteStore) *Service {
	return &Service{client: client, extractor: extractorClient, store: sqliteStore}
}

func (s *Service) SyncUC(ctx context.Context, in SyncUCRequest) SyncUCResponse {
	reqCtx := neoenergia.RequestContext{
		BearerToken: in.BearerToken,
		Documento:   in.Documento,
	}

	out := SyncUCResponse{
		Documento: in.Documento,
		UC:        in.UC,
	}

	out.GrupoCliente = callResult(func() (neoenergia.GrupoClienteResponse, error) {
		return s.client.GetGrupoCliente(ctx, reqCtx)
	})
	out.MinhaConta = callResult(func() (neoenergia.MinhaContaResponse, error) {
		return s.client.GetMinhaConta(ctx, reqCtx)
	})
	out.MinhaContaLegado = callResult(func() (neoenergia.MinhaContaLegadoResponse, error) {
		return s.client.GetMinhaContaLegado(ctx, reqCtx)
	})
	out.UCs = callResult(func() (neoenergia.UCsResponse, error) {
		return s.client.ListUCs(ctx, reqCtx)
	})
	out.Imovel = callResult(func() (neoenergia.ImovelResponse, error) {
		return s.client.GetImovel(ctx, reqCtx, in.UC)
	})

	out.Protocolo = callResult(func() (neoenergia.ProtocoloResponse, error) {
		return s.client.GetProtocolo(ctx, reqCtx, in.UC)
	})

	var protocolo string
	if out.Protocolo.Data != nil {
		protocolo = out.Protocolo.Data.ProtocoloSalesforceStr
		if protocolo == "" {
			protocolo = out.Protocolo.Data.ProtocoloLegadoStr
		}
	}

	if protocolo != "" {
		out.HistoricoConsumo = callResult(func() (neoenergia.HistoricoConsumoResponse, error) {
			return s.client.GetHistoricoConsumo(ctx, reqCtx, in.UC, protocolo)
		})
		out.Faturas = callResult(func() (neoenergia.FaturasResponse, error) {
			return s.client.ListFaturas(ctx, reqCtx, in.UC, protocolo)
		})
	}

	out.DataCerta = callResult(func() (neoenergia.DataCertaResponse, error) {
		return s.client.GetDataCerta(ctx, reqCtx, in.UC)
	})
	out.FaturaDigital = callResult(func() (neoenergia.FaturaDigitalResponse, error) {
		return s.client.GetFaturaDigital(ctx, reqCtx, in.UC)
	})
	out.MotivosSegundaVia = callResult(func() (neoenergia.MotivosSegundaViaResponse, error) {
		return s.client.GetMotivosSegundaVia(ctx, reqCtx, in.UC)
	})

	if out.Imovel.Data != nil {
		codCliente := out.Imovel.Data.Cliente.Codigo
		out.DebitoAutomatico = callResult(func() (neoenergia.DebitoAutomaticoResponse, error) {
			return s.client.GetDebitoAutomatico(ctx, reqCtx, in.UC, codCliente)
		})
	}

	var selectedFatura *neoenergia.Fatura
	if protocolo != "" && out.Faturas.Data != nil && len(out.Faturas.Data.Faturas) > 0 {
		selectedFatura = &out.Faturas.Data.Faturas[0]
		out.DadosPagamento = callResult(func() (neoenergia.DadosPagamentoResponse, error) {
			return s.client.GetDadosPagamento(ctx, reqCtx, in.UC, selectedFatura.NumeroFatura, protocolo)
		})
		if in.IncludePDF {
			motivo := "02"
			if out.MotivosSegundaVia.Data != nil && len(out.MotivosSegundaVia.Data.Motivos) > 0 {
				motivo = out.MotivosSegundaVia.Data.Motivos[0].IDMotivo
			}
			out.PDF = callResult(func() (neoenergia.FaturaPDFResponse, error) {
				return s.client.GetFaturaPDF(ctx, reqCtx, in.UC, selectedFatura.NumeroFatura, protocolo, motivo)
			})
			if in.IncludeExtraction && s.extractor != nil && out.PDF.Data != nil {
				request := extractor.Request{
					SchemaVersion: "1.0.0",
					JobID:         "sync-" + in.UC,
					UC:            in.UC,
					Documento:     reqCtx.Documento,
					NumeroFatura:  selectedFatura.NumeroFatura,
					MesReferencia: selectedFatura.MesReferencia,
					PDF: extractor.PDFPayload{
						Mode:     "base64",
						Base64:   out.PDF.Data.FileData,
						FileName: out.PDF.Data.FileName + out.PDF.Data.FileExtension,
					},
					Requested: []string{"itens_fatura", "composicao_fornecimento", "nota_fiscal"},
					APISnapshot: map[string]any{
						"fatura":          selectedFatura,
						"dados_pagamento": out.DadosPagamento.Data,
					},
				}
				out.Extraction = callResult(func() (extractor.Response, error) {
					return s.extractor.Extract(ctx, request)
				})
			}
		}
	}
	if selectedFatura != nil {
		out.BillingRecord = buildBillingRecord(out, *selectedFatura)
		out.DocumentRecord = buildDocumentRecord(out)
	}
	out.Persistence = s.persist(ctx, in, out, selectedFatura)

	return out
}

func (s *Service) persist(ctx context.Context, in SyncUCRequest, out SyncUCResponse, selectedFatura *neoenergia.Fatura) *PersistenceResult {
	if s.store == nil {
		return nil
	}
	status := "succeeded"
	errorMessage := ""
	if out.BillingRecord == nil || out.BillingRecord.Completeness.Status != "complete" {
		status = "partial"
	}
	if selectedFatura == nil {
		status = "failed"
		errorMessage = "nenhuma fatura encontrada"
	}

	persisted, err := s.store.PersistSyncResult(store.PersistSyncInput{
		CredentialID:   in.CredentialID,
		Documento:      out.Documento,
		UC:             out.UC,
		Status:         status,
		ErrorMessage:   errorMessage,
		UCRecord:       findUC(out.UCs.Data, out.UC),
		Imovel:         out.Imovel.Data,
		Fatura:         selectedFatura,
		Historico:      out.HistoricoConsumo.Data,
		DadosPagamento: out.DadosPagamento.Data,
		PDF:            out.PDF.Data,
		Extraction:     out.Extraction.Data,
		BillingRecord:  out.BillingRecord,
		DocumentRecord: out.DocumentRecord,
		RawResponse:    out,
	})
	if err != nil {
		return &PersistenceResult{Status: "error", Error: err.Error()}
	}
	_ = ctx
	return &PersistenceResult{
		SyncRunID: persisted.SyncRunID,
		InvoiceID: persisted.InvoiceID,
		Status:    status,
	}
}

func findUC(response *neoenergia.UCsResponse, uc string) *neoenergia.UC {
	if response == nil {
		return nil
	}
	for idx := range response.UCs {
		if response.UCs[idx].UC == uc {
			return &response.UCs[idx]
		}
	}
	return nil
}

func callResult[T any](fn func() (T, error)) OptionalResult[T] {
	value, err := fn()
	if err == nil {
		return OptionalResult[T]{Data: &value}
	}

	if httpErr, ok := err.(*neoenergia.ErrorResponse); ok {
		return OptionalResult[T]{Error: httpErr}
	}

	raw, _ := json.Marshal(map[string]any{"message": err.Error()})
	return OptionalResult[T]{
		Error: &neoenergia.ErrorResponse{
			StatusCode: 0,
			Method:     "",
			Path:       "",
			Body:       raw,
		},
	}
}
