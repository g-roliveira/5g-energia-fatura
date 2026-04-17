package sync

import (
	"slices"

	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/neoenergia"
)

func buildBillingRecord(out SyncUCResponse, fatura neoenergia.Fatura) *BillingRecord {
	sourceMap := map[string]string{}
	confidenceMap := map[string]float64{}
	rawAvailability := map[string]bool{
		"api.fatura":             true,
		"api.dados_pagamento":    out.DadosPagamento.Data != nil,
		"api.historico_consumo":  out.HistoricoConsumo.Data != nil,
		"api.pdf":                out.PDF.Data != nil,
		"extractor.document_pdf": out.Extraction.Data != nil,
	}

	record := &BillingRecord{
		UC:                     firstNonEmpty(fatura.UC, out.UC),
		NumeroFatura:           fatura.NumeroFatura,
		MesReferencia:          fatura.MesReferencia,
		StatusFatura:           fatura.StatusFatura,
		ValorTotal:             fatura.ValorEmissao,
		CodigoBarras:           nil,
		DataEmissao:            stringPtrIfNotEmpty(fatura.DataEmissao),
		DataVencimento:         stringPtrIfNotEmpty(fatura.DataVencimento),
		DataPagamento:          stringPtrIfNotEmpty(emptyDateAsBlank(fatura.DataPagamento)),
		DataInicioPeriodo:      stringPtrIfNotEmpty(fatura.DataInicioPeriodo),
		DataFimPeriodo:         stringPtrIfNotEmpty(fatura.DataFimPeriodo),
		ItensFatura:            nil,
		ComposicaoFornecimento: nil,
		NotaFiscal:             nil,
		SourceMap:              sourceMap,
		ConfidenceMap:          confidenceMap,
		RawSourceAvailability:  rawAvailability,
	}

	setAPISource(sourceMap, confidenceMap, "uc")
	setAPISource(sourceMap, confidenceMap, "numero_fatura")
	setAPISource(sourceMap, confidenceMap, "mes_referencia")
	setAPISource(sourceMap, confidenceMap, "status_fatura")
	setAPISource(sourceMap, confidenceMap, "valor_total")
	setAPISource(sourceMap, confidenceMap, "data_emissao")
	setAPISource(sourceMap, confidenceMap, "data_vencimento")
	setAPISource(sourceMap, confidenceMap, "data_pagamento")
	setAPISource(sourceMap, confidenceMap, "data_inicio_periodo")
	setAPISource(sourceMap, confidenceMap, "data_fim_periodo")

	if out.DadosPagamento.Data != nil && out.DadosPagamento.Data.CodBarras != "" {
		record.CodigoBarras = &out.DadosPagamento.Data.CodBarras
		setAPISource(sourceMap, confidenceMap, "codigo_barras")
	}
	if out.HistoricoConsumo.Data != nil {
		record.HistoricoConsumo = out.HistoricoConsumo.Data.HistoricoConsumo
		setAPISource(sourceMap, confidenceMap, "historico_consumo")
	}
	if out.Extraction.Data != nil {
		status := out.Extraction.Data.Status
		record.ExtractorStatus = &status
		mergeExtractorField(out.Extraction.Data.Fields, out.Extraction.Data.SourceMap, out.Extraction.Data.ConfidenceMap, "itens_fatura", &record.ItensFatura, sourceMap, confidenceMap)
		mergeExtractorField(out.Extraction.Data.Fields, out.Extraction.Data.SourceMap, out.Extraction.Data.ConfidenceMap, "composicao_fornecimento", &record.ComposicaoFornecimento, sourceMap, confidenceMap)
		mergeExtractorField(out.Extraction.Data.Fields, out.Extraction.Data.SourceMap, out.Extraction.Data.ConfidenceMap, "nota_fiscal", &record.NotaFiscal, sourceMap, confidenceMap)
	}

	record.Completeness = computeCompleteness(record)
	return record
}

func mergeExtractorField(fields map[string]any, extractorSources map[string]string, extractorConfidence map[string]float64, field string, target *any, sourceMap map[string]string, confidenceMap map[string]float64) {
	value, ok := fields[field]
	if !ok || isEmptyValue(value) {
		sourceMap[field] = "missing"
		confidenceMap[field] = 0
		return
	}
	*target = value
	source := firstNonEmpty(extractorSources[field], "extractor")
	sourceMap[field] = source
	if confidence, ok := extractorConfidence[field]; ok {
		confidenceMap[field] = confidence
		return
	}
	confidenceMap[field] = 0.8
}

func computeCompleteness(record *BillingRecord) BillingCompleteness {
	required := map[string]bool{
		"uc":              record.UC != "",
		"numero_fatura":   record.NumeroFatura != "",
		"mes_referencia":  record.MesReferencia != "",
		"valor_total":     record.ValorTotal != "",
		"data_vencimento": record.DataVencimento != nil,
		"codigo_barras":   record.CodigoBarras != nil,
		"itens_fatura":    !isEmptyValue(record.ItensFatura),
		"nota_fiscal":     !isEmptyValue(record.NotaFiscal),
	}
	missing := make([]string, 0)
	for field, ok := range required {
		if !ok {
			missing = append(missing, field)
		}
	}
	slices.Sort(missing)
	if len(missing) == 0 {
		return BillingCompleteness{Status: "complete"}
	}
	return BillingCompleteness{Status: "partial", MissingFields: missing}
}

func setAPISource(sourceMap map[string]string, confidenceMap map[string]float64, field string) {
	sourceMap[field] = "api"
	confidenceMap[field] = 1
}

func stringPtrIfNotEmpty(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func emptyDateAsBlank(value string) string {
	if value == "0000-00-00" {
		return ""
	}
	return value
}

func isEmptyValue(value any) bool {
	if value == nil {
		return true
	}
	switch typed := value.(type) {
	case string:
		return typed == ""
	case []any:
		return len(typed) == 0
	case map[string]any:
		return len(typed) == 0
	default:
		return false
	}
}
