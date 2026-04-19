package normalizer

import (
	"regexp"
	"strings"
)

// ExtractSCEE tenta achar e parsear o bloco de compensação no texto livre
// do rodapé "Informações Importantes" da fatura. Retorna nil se não
// encontrar nenhum dos 3 layouts conhecidos.
//
// Padrões reais observados:
//
// LAYOUT 1 (MMGD legado — 10/2023):
//   "Unidade Microgeracao. Energia injetada no mes 8736 kWh.
//    Saldo total de credito para o proximo faturamento 46952 kWh"
//
// LAYOUT 2 (MMGD transição — 07/2024): mesma nomenclatura do legado.
//
// LAYOUT 3 (SCEE moderno — 12/2025+):
//   "SCEE, Excedente 3170 kWh e creditos utilizados 0 kWh.
//    Saldo para o proximo ciclo 0 kWh"
//
// É tolerante a variações de OCR (acentos faltando, "creditos" sem acento,
// espaçamento irregular).
func ExtractSCEE(rawText string) *SCEESummary {
	t := normalizeForMatch(rawText)
	t = collapseWhitespace(t)

	if s := extractSCEEModerno(t, rawText); s != nil {
		return s
	}
	if s := extractMMGDLegado(t, rawText); s != nil {
		return s
	}
	return nil
}

func extractSCEEModerno(normalized, raw string) *SCEESummary {
	if !strings.Contains(normalized, "scee") {
		return nil
	}
	reExc := regexp.MustCompile(`excedente\s+([\d\.,]+)\s*kwh`)
	reCred := regexp.MustCompile(`creditos\s+utilizados\s+([\d\.,]+)\s*kwh`)
	reSaldo := regexp.MustCompile(`saldo\s+(?:para\s+o\s+)?proximo\s+ciclo\s+([\d\.,]+)\s*kwh`)

	s := &SCEESummary{Layout: SCEELayoutSCEEModerno, RawText: raw, Confidence: 0.9}
	if m := reExc.FindStringSubmatch(normalized); len(m) == 2 {
		if v, err := parseBRDecimal(m[1]); err == nil {
			s.ExcedenteKWh = v
		}
	}
	if m := reCred.FindStringSubmatch(normalized); len(m) == 2 {
		if v, err := parseBRDecimal(m[1]); err == nil {
			s.CreditosUtilizados = v
		}
	}
	if m := reSaldo.FindStringSubmatch(normalized); len(m) == 2 {
		if v, err := parseBRDecimal(m[1]); err == nil {
			s.SaldoProximoCiclo = v
		}
	}
	if s.ExcedenteKWh.IsZero() && s.CreditosUtilizados.IsZero() && s.SaldoProximoCiclo.IsZero() {
		return nil
	}
	return s
}

func extractMMGDLegado(normalized, raw string) *SCEESummary {
	if !strings.Contains(normalized, "microgeracao") {
		return nil
	}
	reInj := regexp.MustCompile(`energia\s+injetada\s+no\s+mes\s+([\d\.,]+)\s*kwh`)
	reSaldo := regexp.MustCompile(`saldo\s+total\s+de\s+credito\s+(?:para\s+o\s+proximo\s+faturamento\s+)?([\d\.,]+)\s*kwh`)

	s := &SCEESummary{Layout: SCEELayoutMMGDLegado, RawText: raw, Confidence: 0.9}
	if m := reInj.FindStringSubmatch(normalized); len(m) == 2 {
		if v, err := parseBRDecimal(m[1]); err == nil {
			s.EnergiaInjetadaKWh = v
		}
	}
	if m := reSaldo.FindStringSubmatch(normalized); len(m) == 2 {
		if v, err := parseBRDecimal(m[1]); err == nil {
			s.SaldoProximoCiclo = v
		}
	}
	if s.EnergiaInjetadaKWh.IsZero() && s.SaldoProximoCiclo.IsZero() {
		return nil
	}
	return s
}

func collapseWhitespace(s string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(s, " ")
}
