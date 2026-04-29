package billing

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// Calculate é o motor determinístico de cálculo de faturamento de energia
// compartilhada no modelo de negócio 5G.
//
// Dada uma entrada canônica (contrato + itens normalizados da fatura Coelba),
// produz um resultado em duas vertentes: "Sem Desconto" e "Com Desconto".
// A diferença entre elas é a economia gerada para o cliente.
//
// Algoritmo (modelo 5G):
//
//  1. Classifica os itens da fatura por Tipo.
//     TUSD Fio + TUSD Energia formam o "Custo de Disponibilidade".
//     Energia Injetada é tratada separadamente (kWh × tarifa).
//     Bandeira e IP Coelba são valores fixos em R$.
//
//  2. Custo de Disponibilidade (TUSD) — IGUAL nos dois cenários:
//     tusdFio.Qtde × tusdFio.Preco + tusdEnergia.Qtde × tusdEnergia.Preco
//     NOTA: diferente do modelo anterior, NÃO há "consumo líquido".
//     O TUSD da fatura Coelba é usado como está.
//
//  3. Energia Injetada:
//     Sem desconto: valor cheio (kWh_injetados × tarifa_injeção)
//     Com desconto: valor cheio × FatorRepasseEnergia (ex: 0,85 = 85%)
//
//  4. Bandeira tarifária:
//     Se BandeiraComDesconto=true: aplica FatorRepasseEnergia
//     Se false: valor cheio (passthrough) — igual nos dois cenários
//
//  5. Iluminação Pública (IP):
//     Sem desconto: valor cheio da fatura Coelba
//     Com desconto: valor contratual (ValorIPComDesconto)
//
//  6. Economia = TotalSemDesconto - TotalComDesconto
func Calculate(input CalcInput) (*CalcResult, error) {
	if err := validate(input); err != nil {
		return nil, err
	}

	// 1. Classify items
	tusdFio, tusdEnergia, injetada, ipCoelba, bandeiras := classifyItems(input.Itens)

	if tusdFio == nil || tusdEnergia == nil {
		return nil, fmt.Errorf("engine: TUSD fio e TUSD energia são obrigatórios")
	}

	// 2. Custo de disponibilidade (TUSD) — IGUAL nos dois cenários
	custoDisponibilidade := tusdFio.Quantidade.Mul(tusdFio.PrecoUnitario).
		Add(tusdEnergia.Quantidade.Mul(tusdEnergia.PrecoUnitario))

	// 3. Valor da energia injetada
	var valorEnergiaInjetada decimal.Decimal
	var injQty, injPreco decimal.Decimal
	if injetada != nil {
		valorEnergiaInjetada = injetada.Quantidade.Mul(injetada.PrecoUnitario)
		injQty = injetada.Quantidade
		injPreco = injetada.PrecoUnitario
	}

	// 4. Bandeira (passthrough ou com fator)
	totalBandeira := decimal.Zero
	for _, b := range bandeiras {
		totalBandeira = totalBandeira.Add(b.ValorTotal)
	}

	// 5. Iluminação pública
	var ipSemDesconto decimal.Decimal
	if ipCoelba != nil {
		ipSemDesconto = ipCoelba.ValorTotal
	}
	ipComDesconto := input.Contract.ValorIPComDesconto

	// 6. Totais — Sem Desconto
	totalSemDesconto := custoDisponibilidade.
		Add(valorEnergiaInjetada).
		Add(totalBandeira).
		Add(ipSemDesconto)

	// 7. Totais — Com Desconto
	valorComDesconto := valorEnergiaInjetada.Mul(input.Contract.FatorRepasseEnergia)

	bandeiraComDesconto := totalBandeira
	if input.Contract.BandeiraComDesconto {
		bandeiraComDesconto = totalBandeira.Mul(input.Contract.FatorRepasseEnergia)
	}

	totalComDesconto := custoDisponibilidade.
		Add(valorComDesconto).
		Add(bandeiraComDesconto).
		Add(ipComDesconto)

	// 8. Economia
	economia := totalSemDesconto.Sub(totalComDesconto)
	economiaPct := decimal.Zero
	if !totalSemDesconto.IsZero() {
		economiaPct = economia.Div(totalSemDesconto)
	}

	// 9. Line breakdown
	linhas := []LineBreakdown{
		{
			Label: "Custo de Disponibilidade (TUSD)",
			Quantidade: func() decimal.Decimal {
				if tusdFio.Quantidade.Equal(tusdEnergia.Quantidade) {
					return tusdFio.Quantidade
				}
				return decimal.Zero
			}(),
			PrecoUnitario: func() decimal.Decimal {
				if tusdFio.Quantidade.Equal(tusdEnergia.Quantidade) {
					return tusdFio.PrecoUnitario.Add(tusdEnergia.PrecoUnitario)
				}
				return decimal.Zero
			}(),
			ValorSemDesc: custoDisponibilidade,
			ValorComDesc: custoDisponibilidade,
		},
		{
			Label:         "Energia Injetada / Compensada",
			Quantidade:    injQty,
			PrecoUnitario: injPreco,
			ValorSemDesc:  valorEnergiaInjetada,
			ValorComDesc:  valorComDesconto,
		},
	}

	if totalBandeira.GreaterThan(decimal.Zero) {
		linhas = append(linhas, LineBreakdown{
			Label:        "Bandeira Tarifária",
			ValorSemDesc: totalBandeira,
			ValorComDesc: bandeiraComDesconto,
		})
	}

	linhas = append(linhas, LineBreakdown{
		Label:        "Iluminação Pública",
		ValorSemDesc: ipSemDesconto,
		ValorComDesc: ipComDesconto,
	})

	// 10. Warnings
	warnings := []string{}
	if injetada == nil || valorEnergiaInjetada.IsZero() {
		warnings = append(warnings, "Nenhuma energia injetada encontrada — verifique se a UC tem geração no período")
	}

	return &CalcResult{
		TotalSemDesconto: totalSemDesconto,
		TotalComDesconto: totalComDesconto,
		EconomiaRS:       economia,
		EconomiaPct:      economiaPct,
		Linhas:           linhas,
		Warnings:         warnings,
	}, nil
}

// classifyItems percorre os itens e os classifica por tipo.
// Retorna ponteiros para os itens de TUSD Fio, TUSD Energia, Injetada e IP Coelba,
// além de um slice com todas as bandeiras.
func classifyItems(itens []UtilityInvoiceItem) (
	tusdFio, tusdEnergia, injetada, ipCoelba *UtilityInvoiceItem,
	bandeiras []UtilityInvoiceItem,
) {
	for i := range itens {
		item := &itens[i]
		switch item.Tipo {
		case ItemTUSDFio:
			tusdFio = item
		case ItemTUSDEnergia:
			tusdEnergia = item
		case ItemEnergiaInjetada:
			injetada = item
		case ItemBandeira:
			bandeiras = append(bandeiras, itens[i])
		case ItemIPCoelba:
			ipCoelba = item
		}
	}
	return
}

func validate(input CalcInput) error {
	if input.Contract.FatorRepasseEnergia.IsZero() {
		return fmt.Errorf("engine: FatorRepasseEnergia do contrato não pode ser zero (use 1.0 para repasse integral)")
	}
	if input.Contract.FatorRepasseEnergia.GreaterThan(decimal.NewFromInt(1)) {
		return fmt.Errorf("engine: FatorRepasseEnergia deve estar entre 0 e 1 (ex: 0.85 para cobrar 85%% do valor)")
	}
	if len(input.Itens) == 0 {
		return fmt.Errorf("engine: lista de itens vazia")
	}
	return nil
}
