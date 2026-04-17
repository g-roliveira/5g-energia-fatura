package sync

import "fmt"

type DocumentRecord struct {
	DataVencimento   *string        `json:"data_vencimento"`
	Nome             string         `json:"nome"`
	NormalizadoValor any            `json:"normalizado_valor"`
	OCR              map[string]any `json:"ocr"`
	UC               string         `json:"uc"`
	Valor            string         `json:"valor"`
	SiteReceipt      string         `json:"site_receipt"`
}

func buildDocumentRecord(out SyncUCResponse) *DocumentRecord {
	if out.BillingRecord == nil || out.Extraction.Data == nil || out.Extraction.Data.Fields == nil {
		return nil
	}

	ocr := cloneMap(out.Extraction.Data.Fields)
	if out.BillingRecord.CodigoBarras != nil && isEmptyValue(ocr["codigo_barras"]) {
		ocr["codigo_barras"] = *out.BillingRecord.CodigoBarras
	}
	if !isEmptyValue(out.BillingRecord.ItensFatura) && isEmptyValue(ocr["itens_fatura"]) {
		ocr["itens_fatura"] = out.BillingRecord.ItensFatura
	}
	if !isEmptyValue(out.BillingRecord.NotaFiscal) && isEmptyValue(ocr["nota_fiscal"]) {
		ocr["nota_fiscal"] = out.BillingRecord.NotaFiscal
	}

	normalizadoValor := ocr["normalizado_valor"]
	if isEmptyValue(normalizadoValor) {
		normalizadoValor = out.BillingRecord.ValorTotal
	}

	valor := asString(ocr["valor"])
	if valor == "" && out.BillingRecord.ValorTotal != "" {
		valor = "R$ " + out.BillingRecord.ValorTotal
	}

	return &DocumentRecord{
		DataVencimento:   out.BillingRecord.DataVencimento,
		Nome:             extractNomeFromOCR(ocr),
		NormalizadoValor: normalizadoValor,
		OCR:              ocr,
		UC:               out.BillingRecord.UC,
		Valor:            valor,
		SiteReceipt:      fmt.Sprintf("neoenergia-private-api://faturas/%s/pdf", out.BillingRecord.NumeroFatura),
	}
}

func cloneMap(input map[string]any) map[string]any {
	output := make(map[string]any, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}

func extractNomeFromOCR(ocr map[string]any) string {
	cliente, ok := ocr["cliente"].(map[string]any)
	if !ok {
		return ""
	}
	return asString(cliente["nome"])
}

func asString(value any) string {
	if value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	default:
		return fmt.Sprint(typed)
	}
}
