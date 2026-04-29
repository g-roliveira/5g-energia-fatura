package cycle

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-pdf/fpdf"
)

// ---------------------------------------------------------------------------
// Data structures for PDF generation
// ---------------------------------------------------------------------------

// InvoicePDFData contem todos os dados necessarios para gerar a fatura em PDF.
type InvoicePDFData struct {
	CustomerName   string
	UCCode         string
	Address        string
	City           string
	UF             string
	ReferenceMonth string
	IssueDate      string

	Lines            []InvoicePDFLine
	TotalSemDesconto string
	TotalComDesconto string
	EconomiaRS       string
	EconomiaPct      string
}

// InvoicePDFLine representa uma linha do detalhamento da fatura.
type InvoicePDFLine struct {
	Label         string
	Quantidade    string
	PrecoUnitario string
	ValorSemDesc  string
	ValorComDesc  string
}

// PDFGenerationResult contem o resultado da geracao do PDF.
type PDFGenerationResult struct {
	Bytes    []byte
	Checksum string // SHA-256 hex string
	FileName string // e.g. fatura_007085489099_2026_04.pdf
	FilePath string // caminho absoluto completo do arquivo salvo
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

// GenerateAndSaveInvoicePDF gera o PDF da fatura, salva em disco e retorna
// metadados (checksum, caminho do arquivo). O diretorio de saida e definido
// por outputDir (tipicamente de PDF_OUTPUT_DIR env var, default "./pdfs/").
// O arquivo e salvo em um subdiretorio UUID para evitar colisoes.
func GenerateAndSaveInvoicePDF(data *InvoicePDFData, outputDir string) (*PDFGenerationResult, error) {
	raw, err := generateInvoicePDF(data)
	if err != nil {
		return nil, fmt.Errorf("generate pdf: %w", err)
	}

	// SHA-256 checksum
	h := sha256.Sum256(raw)
	checksum := fmt.Sprintf("%x", h)

	// Nome do arquivo: fatura_{uc_code}_{year}_{month}.pdf
	parts := strings.Split(data.ReferenceMonth, "/")
	monthName := strings.TrimSpace(parts[0])
	year := strings.TrimSpace(parts[1])
	monthNum := monthNumber(monthName)
	fileName := fmt.Sprintf("fatura_%s_%s_%02d.pdf", data.UCCode, year, monthNum)

	// Diretorio UUID para evitar colisoes
	uuidDir := fmt.Sprintf("%x", sha256.Sum256([]byte(fileName+data.UCCode)))[:12]
	saveDir := filepath.Join(outputDir, uuidDir)
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return nil, fmt.Errorf("mkdir %s: %w", saveDir, err)
	}

	fullPath := filepath.Join(saveDir, fileName)
	if err := os.WriteFile(fullPath, raw, 0644); err != nil {
		return nil, fmt.Errorf("write %s: %w", fullPath, err)
	}

	return &PDFGenerationResult{
		Bytes:    raw,
		Checksum: checksum,
		FileName: fileName,
		FilePath: fullPath,
	}, nil
}

// ---------------------------------------------------------------------------
// PDF rendering (private)
// ---------------------------------------------------------------------------

func generateInvoicePDF(data *InvoicePDFData) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	lm := 20.0 // left margin
	w := 170.0 // usable width

	// --- HEADER ---
	pdf.SetFont("Helvetica", "B", 20)
	pdf.SetTextColor(0, 90, 170)
	pdf.CellFormat(w, 12, "5G ENERGIA SOLAR", "", 1, "C", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(w, 6, "Fatura de Energia - Conta de Energia", "", 1, "C", false, 0, "")
	pdf.Ln(3)

	// Separator
	pdf.SetDrawColor(0, 90, 170)
	pdf.SetLineWidth(0.8)
	pdf.Line(lm, pdf.GetY(), lm+w, pdf.GetY())
	pdf.Ln(5)

	// Reference month big text
	pdf.SetFont("Helvetica", "B", 14)
	pdf.SetTextColor(40, 40, 40)
	refLabel := fmt.Sprintf("Referente a: %s", data.ReferenceMonth)
	pdf.CellFormat(w, 8, sanitize(refLabel), "", 1, "L", false, 0, "")
	pdf.Ln(3)

	// --- CLIENT INFO BOX ---
	boxY := pdf.GetY()
	pdf.SetFillColor(245, 247, 250)
	pdf.SetDrawColor(200, 210, 220)
	pdf.SetLineWidth(0.4)
	boxH := 36.0
	pdf.Rect(lm, boxY, w, boxH, "D")
	pdf.Rect(lm, boxY, w, boxH, "F")

	// Title inside box
	pdf.SetXY(lm+2, boxY+1)
	pdf.SetFont("Helvetica", "B", 7)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(w-4, 4, "DADOS DO CLIENTE", "", 1, "L", false, 0, "")

	// Client info
	pdf.SetXY(lm+2, boxY+6)
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(30, 30, 30)

	colX := lm + 2
	midX := lm + w/2 + 2
	lineH := 5.5

	// Row 1
	pdf.SetXY(colX, boxY+6)
	pdf.CellFormat(w/2-4, lineH, "Cliente: "+sanitize(data.CustomerName), "", 1, "L", false, 0, "")
	pdf.SetXY(midX, boxY+6)
	pdf.CellFormat(w/2-4, lineH, "UC: "+data.UCCode, "", 1, "L", false, 0, "")

	// Row 2
	pdf.SetXY(colX, boxY+6+lineH)
	pdf.CellFormat(w/2-4, lineH, "Endereco: "+sanitize(data.Address), "", 1, "L", false, 0, "")
	pdf.SetXY(midX, boxY+6+lineH)
	pdf.CellFormat(w/2-4, lineH, "Cidade: "+sanitize(data.City)+"/"+data.UF, "", 1, "L", false, 0, "")
	row2addr := data.City + "/" + data.UF
	_ = row2addr

	// Row 3
	pdf.SetXY(colX, boxY+6+2*lineH)
	pdf.CellFormat(w/2-4, lineH, "Data de Emissao: "+data.IssueDate, "", 1, "L", false, 0, "")
	pdf.SetXY(midX, boxY+6+2*lineH)
	pdf.CellFormat(w/2-4, lineH, "Mes de Referencia: "+data.ReferenceMonth, "", 1, "L", false, 0, "")

	pdf.SetY(boxY + boxH + 5)

	// --- TABLE: COMPOSICAO DA FATURA ---
	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetFillColor(0, 90, 170)
	pdf.SetTextColor(255, 255, 255)

	// Column widths
	colItem := 64.0
	colQtd := 22.0
	colTarifa := 26.0
	colSemDesc := 29.0
	colComDesc := 29.0

	// Table header
	headers := []string{"Item", "Qtd", "Tarifa (R$)", "Valor Ref. (R$)", "Repasse 5G (R$)"}
	colWidths := []float64{colItem, colQtd, colTarifa, colSemDesc, colComDesc}

	for i, h := range headers {
		align := "C"
		if i == 0 {
			align = "L"
		}
		pdf.CellFormat(colWidths[i], 7, h, "1", 0, align, true, 0, "")
	}
	pdf.Ln(-1)

	// Table body
	pdf.SetFont("Helvetica", "", 8)
	pdf.SetTextColor(30, 30, 30)

	for idx, line := range data.Lines {
		if idx%2 == 0 {
			pdf.SetFillColor(248, 250, 252)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		vals := []string{
			sanitize(line.Label),
			line.Quantidade,
			line.PrecoUnitario,
			line.ValorSemDesc,
			line.ValorComDesc,
		}
		for i, v := range vals {
			align := "C"
			if i == 0 {
				align = "L"
			}
			if i >= 3 {
				align = "R"
			}
			if i == 2 {
				align = "R"
			}
			pdf.CellFormat(colWidths[i], 6.5, v, "1", 0, align, true, 0, "")
		}
		pdf.Ln(-1)
	}

	// Totals row
	pdf.SetFont("Helvetica", "B", 9)
	pdf.SetFillColor(230, 238, 248)
	pdf.SetTextColor(0, 0, 0)

	totalData := []string{
		"TOTAL",
		"",
		"",
		"R$ " + data.TotalSemDesconto,
		"R$ " + data.TotalComDesconto,
	}
	for i, v := range totalData {
		align := "C"
		if i == 0 {
			align = "L"
		}
		if i >= 3 {
			align = "R"
		}
		pdf.CellFormat(colWidths[i], 7, v, "1", 0, align, true, 0, "")
	}
	pdf.Ln(-1)

	pdf.Ln(5)

	// --- ECONOMY SUMMARY ---
	economyY := pdf.GetY()

	// Economy box
	pdf.SetFillColor(235, 250, 240)
	pdf.SetDrawColor(80, 180, 120)
	pdf.SetLineWidth(0.6)
	ecoH := 28.0
	pdf.Rect(lm, economyY, w, ecoH, "D")
	pdf.Rect(lm, economyY, w, ecoH, "F")

	pdf.SetXY(lm+3, economyY+2)
	pdf.SetFont("Helvetica", "B", 11)
	pdf.SetTextColor(0, 130, 60)
	pdf.CellFormat(w-6, 7, "RESUMO FINANCEIRO", "", 1, "L", false, 0, "")

	pdf.SetXY(lm+3, economyY+10)
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(40, 40, 40)
	pdf.CellFormat(80, 6, "Valor de Referencia:         R$ "+data.TotalSemDesconto, "", 1, "L", false, 0, "")

	pdf.SetXY(lm+3, economyY+16)
	pdf.CellFormat(80, 6, "Repasse 5G:    R$ "+data.TotalComDesconto, "", 1, "L", false, 0, "")

	pdf.SetXY(lm+3, economyY+22)
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetTextColor(0, 130, 60)
	ecoLabel := fmt.Sprintf("Economia: R$ %s (%s%%)", data.EconomiaRS, data.EconomiaPct)
	pdf.CellFormat(w-6, 6, sanitize(ecoLabel), "", 1, "L", false, 0, "")

	pdf.SetY(economyY + ecoH + 8)

	// --- PIX / QR CODE PLACEHOLDER ---
	pixY := pdf.GetY()
	pixH := 26.0
	pdf.SetFillColor(255, 255, 255)
	pdf.SetDrawColor(180, 180, 180)
	pdf.SetLineWidth(0.4)
	pdf.Rect(lm, pixY, w, pixH, "D")

	// Dashed border effect for QR code area
	qrW := 40.0
	qrH := 20.0
	qrX := lm + 5
	qrY := pixY + 3
	pdf.SetDrawColor(160, 160, 160)
	pdf.SetLineWidth(0.3)
	pdf.Rect(qrX, qrY, qrW, qrH, "D")

	pdf.SetXY(qrX, qrY+7)
	pdf.SetFont("Helvetica", "", 7)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(qrW, 4, "[QR Code PIX]", "", 1, "C", false, 0, "")
	pdf.SetXY(qrX, qrY+11)
	pdf.CellFormat(qrW, 4, "Pague via PIX", "", 1, "C", false, 0, "")

	// PIX info text
	pdf.SetXY(lm+qrW+12, pixY+3)
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetTextColor(40, 40, 40)
	pdf.CellFormat(w-qrW-17, 6, "Pague com PIX", "", 1, "L", false, 0, "")

	pdf.SetXY(lm+qrW+12, pixY+10)
	pdf.SetFont("Helvetica", "", 8)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(w-qrW-17, 4, "Utilize o QR Code ao lado ou a chave abaixo:", "", 1, "L", false, 0, "")

	pdf.SetXY(lm+qrW+12, pixY+15)
	pdf.SetFont("Courier", "B", 8)
	pdf.SetTextColor(0, 90, 170)
	pdf.CellFormat(w-qrW-17, 4, "  00.000.000/0001-00", "", 1, "L", false, 0, "")

	pdf.SetY(pixY + pixH + 10)

	// --- FOOTER ---
	footerY := pdf.GetY()
	pdf.SetDrawColor(200, 200, 200)
	pdf.SetLineWidth(0.4)
	pdf.Line(lm, footerY, lm+w, footerY)
	pdf.Ln(4)

	pdf.SetFont("Helvetica", "", 8)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(w, 4, "5G Energia Solar Ltda. - CNPJ: 00.000.000/0001-00", "", 1, "C", false, 0, "")
	pdf.CellFormat(w, 4, "contato@5genergia.com.br - (71) 99999-9999", "", 1, "C", false, 0, "")
	pdf.CellFormat(w, 4, "Este documento e uma representacao da fatura com desconto de energia solar.", "", 1, "C", false, 0, "")
	pdf.CellFormat(w, 4, "Consulte a fatura original da distribuidora para pagamento.", "", 1, "C", false, 0, "")

	// Output to bytes
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("pdf.Output: %w", err)
	}

	return buf.Bytes(), nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

var accentReplacer = strings.NewReplacer(
	"Á", "A", "á", "a",
	"À", "A", "à", "a",
	"Â", "A", "â", "a",
	"Ã", "A", "ã", "a",
	"Ä", "A", "ä", "a",
	"É", "E", "é", "e",
	"È", "E", "è", "e",
	"Ê", "E", "ê", "e",
	"Ë", "E", "ë", "e",
	"Í", "I", "í", "i",
	"Ì", "I", "ì", "i",
	"Î", "I", "î", "i",
	"Ï", "I", "ï", "i",
	"Ó", "O", "ó", "o",
	"Ò", "O", "ò", "o",
	"Ô", "O", "ô", "o",
	"Õ", "O", "õ", "o",
	"Ö", "O", "ö", "o",
	"Ú", "U", "ú", "u",
	"Ù", "U", "ù", "u",
	"Û", "U", "û", "u",
	"Ü", "U", "ü", "u",
	"Ç", "C", "ç", "c",
	"Ñ", "N", "ñ", "n",
)

func sanitize(s string) string {
	return accentReplacer.Replace(s)
}

// monthName retorna o nome do mes em portugues.
func monthName(m int16) string {
	names := []string{
		"", "Janeiro", "Fevereiro", "Marco", "Abril", "Maio", "Junho",
		"Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro",
	}
	if m < 1 || m > 12 {
		return fmt.Sprintf("%d", m)
	}
	return names[m]
}

// monthNumber converte nome do mes em portugues para numero (1-12).
// Usado para extrair o numero do mes do formato "Abril/2026".
func monthNumber(name string) int {
	names := []string{
		"", "Janeiro", "Fevereiro", "Marco", "Abril", "Maio", "Junho",
		"Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro",
	}
	for i, n := range names {
		if strings.EqualFold(name, n) {
			return i
		}
	}
	return 0
}

// formatCurrency formata um decimal.Decimal como string brasileira (ex: "1.234,56").
func formatCurrency(d float64) string {
	intPart := int(d)
	cents := int((d-float64(intPart))*100 + 0.5)
	if cents < 0 {
		cents = -cents
	}

	// Format integer part with dots as thousands separators
	intStr := fmt.Sprintf("%d", intPart)
	if intPart < 0 {
		intStr = fmt.Sprintf("%d", -intPart)
	}
	var parts []string
	for i := len(intStr); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		parts = append([]string{intStr[start:i]}, parts...)
	}
	formattedInt := strings.Join(parts, ".")
	if intPart < 0 {
		formattedInt = "-" + formattedInt
	}

	return fmt.Sprintf("%s,%02d", formattedInt, cents)
}
