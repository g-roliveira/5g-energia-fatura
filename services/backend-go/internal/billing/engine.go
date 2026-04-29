package billing

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

// Calculate é o motor determinístico de cálculo de faturamento de energia
// compartilhada. Dada uma entrada canônica (contrato + itens normalizados),
// produz um resultado em duas vertentes: "sem desconto" e "com desconto".
//
// Algoritmo:
//
//  1. Coleta os itens por Type. TUSD fio + TUSD energia + energia injetada
//     são tratados como preços unitários (R$/kWh). Bandeira e IP Coelba
//     são valores fixos já em R$.
//
//  2. Consumo líquido = Qtd(TUSD) - Qtd(EnergiaInjetada).
//     Se negativo ou menor que ConsumoMinimoKWh (e
//     CustoDisponibilidadeSempreCobrado=true), usa ConsumoMinimoKWh.
//
//  3. Subtotal "sem desconto" = consumo_liquido * (tusd_fio + tusd_te).
//     Subtotal "com desconto" = subtotal_sem_desconto * DescontoPct.
//
//  4. Bandeira: sempre somada em "sem desconto". Em "com desconto" só
//     sofre desconto se BandeiraComDesconto=true.
//
//  5. IP Coelba: sempre somada. Nunca tem desconto.
//
//  6. IP da Usina: não aparece na fatura Coelba. Adicionada apenas no
//     lado "com desconto" conforme IPFaturamentoMode.
//
//  7. Arredondamento: decimal.Decimal, sem perda. Arredondamento final
//     é responsabilidade do caller (geralmente 2 casas no storage).
func Calculate(in CalculationInput) (CalculationResult, error) {
	if err := validate(in); err != nil {
		return CalculationResult{}, err
	}

	tusdFio, tusdEnergia, injetada, bandeiras, ipCoelba := collectItems(in.Itens)

	if tusdFio == nil || tusdEnergia == nil {
		return CalculationResult{}, errors.New("calc: TUSD fio e TUSD energia são obrigatórios")
	}

	// 2. Consumo líquido
	consumoLiquido := tusdFio.Quantidade
	if injetada != nil {
		consumoLiquido = consumoLiquido.Sub(injetada.Quantidade)
	}

	warnings := []string{}
	minKWh := decimal.NewFromFloat(in.ConsumoMinimoKWh)
	if in.Contract.CustoDisponibilidadeSempreCobrado && consumoLiquido.LessThan(minKWh) {
		warnings = append(warnings,
			fmt.Sprintf("consumo líquido %s kWh menor que mínimo %s kWh — aplicado mínimo",
				consumoLiquido.String(), minKWh.String()))
		consumoLiquido = minKWh
	}
	if consumoLiquido.IsNegative() {
		consumoLiquido = decimal.Zero
	}

	// 3. Subtotais de energia
	precoUnitarioTotal := tusdFio.PrecoUnitario.Add(tusdEnergia.PrecoUnitario)
	subtotalEnergiaSem := consumoLiquido.Mul(precoUnitarioTotal)
	subtotalEnergiaCom := subtotalEnergiaSem.Mul(in.Contract.DescontoPct)

	linhas := []LineBreakdown{
		{
			Label:         fmt.Sprintf("Energia (TUSD+TE) — %s kWh", consumoLiquido.String()),
			Quantidade:    consumoLiquido,
			PrecoUnitario: precoUnitarioTotal,
			ValorSemDesc:  subtotalEnergiaSem,
			ValorComDesc:  subtotalEnergiaCom,
		},
	}

	totalSem := subtotalEnergiaSem
	totalCom := subtotalEnergiaCom

	// 4. Bandeiras
	for _, b := range bandeiras {
		totalSem = totalSem.Add(b.ValorTotal)
		var valorCom decimal.Decimal
		if in.Contract.BandeiraComDesconto {
			valorCom = b.ValorTotal.Mul(in.Contract.DescontoPct)
		} else {
			valorCom = b.ValorTotal
		}
		totalCom = totalCom.Add(valorCom)
		linhas = append(linhas, LineBreakdown{
			Label:        b.Description,
			ValorSemDesc: b.ValorTotal,
			ValorComDesc: valorCom,
		})
	}

	// 5. IP Coelba (quando presente)
	if !ipCoelba.IsZero() {
		totalSem = totalSem.Add(ipCoelba)
		totalCom = totalCom.Add(ipCoelba)
		linhas = append(linhas, LineBreakdown{
			Label:        "Ilum. Púb. Municipal (repasse Coelba)",
			ValorSemDesc: ipCoelba,
			ValorComDesc: ipCoelba,
		})
	}

	// 6. IP da Usina (só no lado "com desconto")
	ipUsina := computeIPUsina(in.Contract, subtotalEnergiaCom)
	if !ipUsina.IsZero() {
		totalCom = totalCom.Add(ipUsina)
		linhas = append(linhas, LineBreakdown{
			Label:        "IP Usina (faturamento da usina)",
			ValorSemDesc: decimal.Zero,
			ValorComDesc: ipUsina,
		})
	}

	// 7. Economia
	economiaRS := totalSem.Sub(totalCom)
	economiaPct := decimal.Zero
	if !totalSem.IsZero() {
		economiaPct = economiaRS.Div(totalSem)
	}

	return CalculationResult{
		TotalSemDesconto: totalSem,
		TotalComDesconto: totalCom,
		EconomiaRS:       economiaRS,
		EconomiaPct:      economiaPct,
		Linhas:           linhas,
		Warnings:         warnings,
	}, nil
}

func collectItems(itens []UtilityInvoiceItem) (
	tusdFio, tusdEnergia, injetada *UtilityInvoiceItem,
	bandeiras []UtilityInvoiceItem,
	ipCoelba decimal.Decimal,
) {
	for i := range itens {
		item := &itens[i]
		switch item.Type {
		case ItemTUSDFio:
			tusdFio = item
		case ItemTUSDEnergia:
			tusdEnergia = item
		case ItemEnergiaInjetada:
			injetada = item
		case ItemBandeira:
			bandeiras = append(bandeiras, *item)
		case ItemIPCoelba:
			ipCoelba = item.ValorTotal
		}
	}
	return
}

func computeIPUsina(c CalcContract, subtotalEnergiaCom decimal.Decimal) decimal.Decimal {
	switch c.IPFaturamentoMode {
	case IPModeFixed:
		return c.IPFaturamentoValor
	case IPModePercent:
		if c.IPFaturamentoPct.IsZero() {
			return decimal.Zero
		}
		return subtotalEnergiaCom.Mul(c.IPFaturamentoPct)
	default:
		return decimal.Zero
	}
}

func validate(in CalculationInput) error {
	if in.Contract.DescontoPct.IsZero() {
		return errors.New("calc: DescontoPct do contrato não pode ser zero (use 1.0 para 'sem desconto')")
	}
	if in.Contract.DescontoPct.GreaterThan(decimal.NewFromInt(1)) {
		return errors.New("calc: DescontoPct deve estar entre 0 e 1 (ex: 0.85 para 15% de desconto)")
	}
	if len(in.Itens) == 0 {
		return errors.New("calc: lista de itens vazia")
	}
	return nil
}
