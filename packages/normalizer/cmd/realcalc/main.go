package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	calcengine "github.com/gustavo/5g-energia-fatura/packages/calc-engine"
	"github.com/gustavo/5g-energia-fatura/packages/normalizer"
	"github.com/shopspring/decimal"
)

type extractedFile struct {
	FileName      string         `json:"file_name"`
	PDFPath       string         `json:"pdf_path"`
	UC            string         `json:"uc"`
	NumeroFatura  string         `json:"numero_fatura"`
	MesReferencia string         `json:"mes_referencia"`
	Fields        map[string]any `json:"fields"`
}

type runResult struct {
	FileName       string                        `json:"file_name"`
	PDFPath        string                        `json:"pdf_path"`
	UC             string                        `json:"uc"`
	NumeroFatura   string                        `json:"numero_fatura"`
	MesReferencia  string                        `json:"mes_referencia"`
	ValorFaturaPDF string                        `json:"valor_fatura_pdf"`
	TotalItens     string                        `json:"total_itens_normalizados"`
	DeltaItensPDF  string                        `json:"delta_itens_vs_pdf"`
	Calc           *calcengine.CalculationResult `json:"calc,omitempty"`
	Warnings       []string                      `json:"warnings,omitempty"`
	Error          string                        `json:"error,omitempty"`
}

type report struct {
	Contract map[string]any `json:"contract"`
	Results  []runResult    `json:"results"`
}

func main() {
	inputDir := flag.String("input-dir", "output/real-calc/extracted", "Diretório com JSON extraído por PDF")
	outputFile := flag.String("output-file", "output/real-calc/report.json", "Arquivo de saída com o relatório")
	desconto := flag.String("desconto", "0.30", "Desconto do contrato (0..1)")
	ipMode := flag.String("ip-mode", "fixed", "Modo IP usina: fixed|percent")
	ipValor := flag.String("ip-valor", "45.00", "Valor fixo da IP usina quando ip-mode=fixed")
	ipPercent := flag.String("ip-percent", "0", "Percentual IP usina quando ip-mode=percent")
	bandeiraComDesconto := flag.Bool("bandeira-com-desconto", false, "Se true, aplica desconto nas bandeiras")
	custoDispSempre := flag.Bool("custo-disp-sempre", true, "Se true, aplica custo de disponibilidade mínimo")
	consumoMinimo := flag.Int("consumo-minimo-kwh", 100, "Consumo mínimo kWh (trifásico=100)")
	flag.Parse()

	descontoDec, err := decimal.NewFromString(*desconto)
	must(err, "desconto inválido")
	ipValorDec, err := decimal.NewFromString(*ipValor)
	must(err, "ip-valor inválido")
	ipPercentDec, err := decimal.NewFromString(*ipPercent)
	must(err, "ip-percent inválido")

	entries, err := os.ReadDir(*inputDir)
	must(err, "falha ao listar input-dir")

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToLower(e.Name()), ".json") {
			files = append(files, filepath.Join(*inputDir, e.Name()))
		}
	}
	sort.Strings(files)

	contract := calcengine.Contract{
		DescontoPct:                       descontoDec,
		IPFaturamentoMode:                 calcengine.IPFaturamentoMode(*ipMode),
		IPFaturamentoValor:                ipValorDec,
		IPFaturamentoPct:                  ipPercentDec,
		BandeiraComDesconto:               *bandeiraComDesconto,
		CustoDisponibilidadeSempreCobrado: *custoDispSempre,
	}

	out := report{
		Contract: map[string]any{
			"desconto":              *desconto,
			"ip_mode":               *ipMode,
			"ip_valor":              *ipValor,
			"ip_percent":            *ipPercent,
			"bandeira_com_desconto": *bandeiraComDesconto,
			"custo_disp_sempre":     *custoDispSempre,
			"consumo_minimo_kwh":    *consumoMinimo,
		},
		Results: make([]runResult, 0, len(files)),
	}

	for _, path := range files {
		result := runResult{FileName: filepath.Base(path)}
		var raw extractedFile
		if err := readJSON(path, &raw); err != nil {
			result.Error = fmt.Sprintf("erro lendo json: %v", err)
			out.Results = append(out.Results, result)
			continue
		}
		result.FileName = raw.FileName
		result.PDFPath = raw.PDFPath

		invoice := buildRawInvoice(raw)
		result.UC = invoice.UC
		result.NumeroFatura = invoice.NumeroFatura
		result.MesReferencia = invoice.MesReferencia

		norm, err := normalizer.Normalize(invoice)
		if err != nil {
			result.Error = fmt.Sprintf("normalizer: %v", err)
			out.Results = append(out.Results, result)
			continue
		}
		result.Warnings = append(result.Warnings, norm.Warnings...)

		totalItens := decimal.Zero
		for _, item := range norm.Items {
			totalItens = totalItens.Add(item.ValorTotal)
		}
		result.TotalItens = totalItens.StringFixed(2)

		valorPDF := extractValor(raw.Fields)
		result.ValorFaturaPDF = valorPDF.StringFixed(2)
		result.DeltaItensPDF = totalItens.Sub(valorPDF).StringFixed(2)

		calc, err := calcengine.Calculate(calcengine.CalculationInput{
			Contract:         contract,
			Itens:            norm.Items,
			ConsumoMinimoKWh: *consumoMinimo,
		})
		if err != nil {
			result.Error = fmt.Sprintf("calc-engine: %v", err)
			out.Results = append(out.Results, result)
			continue
		}
		result.Calc = &calc
		out.Results = append(out.Results, result)
	}

	must(os.MkdirAll(filepath.Dir(*outputFile), 0o755), "falha criando diretório de saída")
	payload, err := json.MarshalIndent(out, "", "  ")
	must(err, "falha serializando relatório")
	must(os.WriteFile(*outputFile, payload, 0o644), "falha escrevendo relatório")

	okCount := 0
	for _, r := range out.Results {
		if r.Error == "" {
			okCount++
		}
	}
	fmt.Printf("processados=%d ok=%d erros=%d report=%s\n", len(out.Results), okCount, len(out.Results)-okCount, *outputFile)
}

func buildRawInvoice(raw extractedFile) normalizer.RawInvoice {
	fields := raw.Fields
	uc := strings.TrimSpace(raw.UC)
	if uc == "" {
		uc = asString(fields["uc"])
	}
	numFatura := strings.TrimSpace(raw.NumeroFatura)
	if numFatura == "" {
		numFatura = asString(fields["numero_fatura"])
	}
	mesRef := strings.TrimSpace(raw.MesReferencia)
	if mesRef == "" {
		mes := asInt(fields["mes"])
		ano := asInt(fields["ano"])
		if mes > 0 && ano > 0 {
			mesRef = fmt.Sprintf("%04d/%02d", ano, mes)
		}
	}
	infos := asString(fields["informacoes_gerais"])

	items := make([]normalizer.RawItem, 0)
	rawItems, _ := fields["itens_fatura"].([]any)
	for _, ri := range rawItems {
		itemMap, ok := ri.(map[string]any)
		if !ok {
			continue
		}
		items = append(items, normalizer.RawItem{
			Descricao:  asString(itemMap["descricao"]),
			Quantidade: asString(itemMap["quantidade"]),
			Tarifa:     asString(itemMap["tarifa"]),
			Valor:      asString(itemMap["valor"]),
			ValorTotal: asString(itemMap["valor_total"]),
		})
	}

	return normalizer.RawInvoice{
		UC:                     uc,
		NumeroFatura:           numFatura,
		MesReferencia:          mesRef,
		Itens:                  items,
		InformacoesImportantes: infos,
	}
}

func extractValor(fields map[string]any) decimal.Decimal {
	if v, ok := fields["normalizado_valor"]; ok {
		switch x := v.(type) {
		case float64:
			return decimal.NewFromFloat(x)
		case string:
			if d, err := decimal.NewFromString(strings.TrimSpace(x)); err == nil {
				return d
			}
		}
	}
	return decimal.Zero
}

func readJSON(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}

func asString(v any) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return decimal.NewFromFloat(x).String()
	default:
		return fmt.Sprint(v)
	}
}

func asInt(v any) int {
	switch x := v.(type) {
	case float64:
		return int(x)
	case int:
		return x
	case string:
		i, _ := decimal.NewFromString(strings.TrimSpace(x))
		return int(i.IntPart())
	default:
		return 0
	}
}

func must(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
		os.Exit(1)
	}
}
