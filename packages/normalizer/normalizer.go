package normalizer

import (
	"fmt"

	calcengine "github.com/gustavo/5g-energia-fatura/packages/calc-engine"
	"github.com/shopspring/decimal"
)

// Normalize é o ponto de entrada público. Recebe um RawInvoice (montado
// pelo caller a partir do sync.BillingRecord do backend-go) e devolve
// um Result pronto pra alimentar o calc-engine.
//
// Itens com flag Ignore (IRRF, reativo, bandeira verde) vão pra
// IgnoredItems — preservados pra auditoria mas sem entrar no cálculo.
// Itens não reconhecidos vão pra UnclassifiedItems com warning.
func Normalize(raw RawInvoice) (*Result, error) {
	result := &Result{
		Items:             make([]calcengine.UtilityInvoiceItem, 0, len(raw.Itens)),
		IgnoredItems:      []IgnoredItem{},
		UnclassifiedItems: []RawItem{},
		Warnings:          []string{},
	}

	for idx, item := range raw.Itens {
		classif, ok := classify(item.Descricao)
		if !ok {
			result.UnclassifiedItems = append(result.UnclassifiedItems, item)
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("item %d (%q) não foi classificado — requer revisão manual",
					idx, item.Descricao))
			continue
		}

		qtd, qtdErr := parseBRDecimal(item.Quantidade)
		tarifa, tarErr := parseBRDecimal(item.Tarifa)
		valor, valErr := parseBRDecimal(item.Valor)
		valorTot, vtErr := parseBRDecimal(item.ValorTotal)

		for _, e := range []error{qtdErr, tarErr, valErr, vtErr} {
			if e != nil {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("item %d (%q): %v", idx, item.Descricao, e))
			}
		}

		// valor_total tem precedência sobre valor (na prática são iguais
		// nos payloads reais, mas em casos de desconto divergem)
		valorCanonico := valorTot
		if valorCanonico.IsZero() && !valor.IsZero() {
			valorCanonico = valor
		}

		if classif.Ignore {
			result.IgnoredItems = append(result.IgnoredItems, IgnoredItem{
				Type:        classif.Type,
				Description: item.Descricao,
				ValorTotal:  valorCanonico,
				Reason:      classif.MatchedBy,
			})
			continue
		}

		result.Items = append(result.Items, calcengine.UtilityInvoiceItem{
			Type:          classif.Type,
			Description:   item.Descricao,
			Quantidade:    qtd,
			PrecoUnitario: tarifa,
			ValorTotal:    valorCanonico,
		})
	}

	// SCEE só roda se tem texto
	if raw.InformacoesImportantes != "" {
		result.SCEE = ExtractSCEE(raw.InformacoesImportantes)
	}

	return result, nil
}

// garantir import usado
var _ = decimal.Zero
