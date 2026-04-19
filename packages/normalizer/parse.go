package normalizer

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

// parseBRDecimal converte um número em formato brasileiro para decimal.Decimal.
// Aceita:
//   "0,75184631" → 0.75184631
//   "1.234,56"   → 1234.56
//   "R$ 315,48"  → 315.48
//   ""           → 0 (não é erro)
//   "2 kWh"      → 2 (unidades são ignoradas)
//   nil / em branco / "0,00" → 0
func parseBRDecimal(raw string) (decimal.Decimal, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return decimal.Zero, nil
	}

	// Remove símbolos e unidades comuns
	s = strings.ReplaceAll(s, "R$", "")
	s = strings.ReplaceAll(s, "kWh", "")
	s = strings.ReplaceAll(s, "kwh", "")
	s = strings.TrimSpace(s)

	// Formato BR: separador de milhar é '.' e decimal é ','
	// Remove pontos (milhar), troca vírgula por ponto (decimal)
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", ".")

	if s == "" {
		return decimal.Zero, nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero, fmt.Errorf("parseBRDecimal: %q não é um número válido: %w", raw, err)
	}
	return d, nil
}
