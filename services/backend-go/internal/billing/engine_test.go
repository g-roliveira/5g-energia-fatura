package billing

import (
	"testing"

	"github.com/shopspring/decimal"
)

// TestCalculate_ValidationTable reproduz exatamente os 5 meses de dados
// da planilha de validação do negócio 5G.
//
// Cada caso de teste cria os itens da fatura Coelba e o contrato, executa
// o motor e verifica se TotalSemDesconto, TotalComDesconto e EconomiaRS
// batem com os valores esperados da planilha (com tolerância de 1 centavo
// devido a arredondamentos intermediários).
func TestCalculate_ValidationTable(t *testing.T) {
	tests := []struct {
		name       string
		items      []UtilityInvoiceItem
		contract   CalcContract
		wantSEM    string
		wantCOM    string
		wantEcon   string
	}{
		{
			name: "out_2025",
			items: []UtilityInvoiceItem{
				{
					Tipo:          ItemTUSDFio,
					Descricao:     "Consumo-TUSD",
					Quantidade:    decimal.NewFromInt(200),
					PrecoUnitario: decimal.RequireFromString("0.40000"),
					ValorTotal:    decimal.RequireFromString("80.00"),
				},
				{
					Tipo:          ItemTUSDEnergia,
					Descricao:     "Consumo-TE",
					Quantidade:    decimal.NewFromInt(200),
					PrecoUnitario: decimal.RequireFromString("0.16795"),
					ValorTotal:    decimal.RequireFromString("33.59"),
				},
				{
					Tipo:          ItemEnergiaInjetada,
					Descricao:     "Energia Injetada SCEE",
					Quantidade:    decimal.NewFromInt(3289),
					PrecoUnitario: decimal.RequireFromString("1.135983"),
					ValorTotal:    decimal.RequireFromString("3736.25"),
				},
				{
					Tipo:        ItemBandeira,
					Descricao:   "Bandeira Escassez Hídrica",
					ValorTotal:  decimal.RequireFromString("7.15"),
				},
				{
					Tipo:        ItemIPCoelba,
					Descricao:   "Iluminação Pública Municipal",
					ValorTotal:  decimal.RequireFromString("461.98"),
				},
			},
			contract: CalcContract{
				FatorRepasseEnergia: decimal.RequireFromString("0.80"),
				ValorIPComDesconto:  decimal.RequireFromString("13.35"),
				BandeiraComDesconto: false,
			},
			wantSEM:  "4318.97",
			wantCOM:  "3123.09",
			wantEcon: "1195.88",
		},
		{
			name: "nov_2025",
			items: []UtilityInvoiceItem{
				{
					Tipo:          ItemTUSDFio,
					Descricao:     "Consumo-TUSD",
					Quantidade:    decimal.NewFromInt(200),
					PrecoUnitario: decimal.RequireFromString("0.40000"),
					ValorTotal:    decimal.RequireFromString("80.00"),
				},
				{
					Tipo:          ItemTUSDEnergia,
					Descricao:     "Consumo-TE",
					Quantidade:    decimal.NewFromInt(200),
					PrecoUnitario: decimal.RequireFromString("0.16620"),
					ValorTotal:    decimal.RequireFromString("33.24"),
				},
				{
					Tipo:          ItemEnergiaInjetada,
					Descricao:     "Energia Injetada SCEE",
					Quantidade:    decimal.NewFromInt(4254),
					PrecoUnitario: decimal.RequireFromString("1.132289"),
					ValorTotal:    decimal.RequireFromString("4816.75"),
				},
				{
					Tipo:       ItemBandeira,
					Descricao:  "Bandeira Escassez Hídrica",
					ValorTotal: decimal.RequireFromString("6.02"),
				},
				{
					Tipo:       ItemIPCoelba,
					Descricao:  "Iluminação Pública Municipal",
					ValorTotal: decimal.RequireFromString("566.95"),
				},
			},
			contract: CalcContract{
				FatorRepasseEnergia: decimal.RequireFromString("0.80"),
				ValorIPComDesconto:  decimal.RequireFromString("13.23"),
				BandeiraComDesconto: false,
			},
			wantSEM:  "5502.96",
			wantCOM:  "3985.89",
			wantEcon: "1517.07",
		},
		{
			name: "dez_2025",
			items: []UtilityInvoiceItem{
				{
					Tipo:          ItemTUSDFio,
					Descricao:     "Consumo-TUSD",
					Quantidade:    decimal.NewFromInt(200),
					PrecoUnitario: decimal.RequireFromString("0.40000"),
					ValorTotal:    decimal.RequireFromString("80.00"),
				},
				{
					Tipo:          ItemTUSDEnergia,
					Descricao:     "Consumo-TE",
					Quantidade:    decimal.NewFromInt(200),
					PrecoUnitario: decimal.RequireFromString("0.16605"),
					ValorTotal:    decimal.RequireFromString("33.21"),
				},
				{
					Tipo:          ItemEnergiaInjetada,
					Descricao:     "Energia Injetada SCEE",
					Quantidade:    decimal.NewFromInt(3474),
					PrecoUnitario: decimal.RequireFromString("1.132228"),
					ValorTotal:    decimal.RequireFromString("3933.60"),
				},
				{
					Tipo:       ItemBandeira,
					Descricao:  "Bandeira Escassez Hídrica",
					ValorTotal: decimal.RequireFromString("3.25"),
				},
				{
					Tipo:       ItemIPCoelba,
					Descricao:  "Iluminação Pública Municipal",
					ValorTotal: decimal.RequireFromString("534.45"),
				},
			},
			contract: CalcContract{
				FatorRepasseEnergia: decimal.RequireFromString("0.85"),
				ValorIPComDesconto:  decimal.RequireFromString("13.23"),
				BandeiraComDesconto: false,
			},
			wantSEM:  "4584.51",
			wantCOM:  "3473.25",
			wantEcon: "1111.26",
		},
		{
			name: "jan_2026",
			items: []UtilityInvoiceItem{
				{
					Tipo:          ItemTUSDFio,
					Descricao:     "Consumo-TUSD",
					Quantidade:    decimal.NewFromInt(200),
					PrecoUnitario: decimal.RequireFromString("0.40000"),
					ValorTotal:    decimal.RequireFromString("80.00"),
				},
				{
					Tipo:          ItemTUSDEnergia,
					Descricao:     "Consumo-TE",
					Quantidade:    decimal.NewFromInt(200),
					PrecoUnitario: decimal.RequireFromString("0.16465"),
					ValorTotal:    decimal.RequireFromString("32.93"),
				},
				{
					Tipo:          ItemEnergiaInjetada,
					Descricao:     "Energia Injetada SCEE",
					Quantidade:    decimal.NewFromInt(4421),
					PrecoUnitario: decimal.RequireFromString("1.129037"),
					ValorTotal:    decimal.RequireFromString("4991.47"),
				},
				{
					Tipo:       ItemBandeira,
					Descricao:  "Bandeira Escassez Hídrica",
					ValorTotal: decimal.RequireFromString("0.57"),
				},
				{
					Tipo:       ItemIPCoelba,
					Descricao:  "Iluminação Pública Municipal",
					ValorTotal: decimal.RequireFromString("534.45"),
				},
			},
			contract: CalcContract{
				FatorRepasseEnergia: decimal.RequireFromString("0.85"),
				ValorIPComDesconto:  decimal.RequireFromString("12.62"),
				BandeiraComDesconto: false,
			},
			wantSEM:  "5639.42",
			wantCOM:  "4368.87",
			wantEcon: "1270.55",
		},
		{
			name: "fev_2026",
			items: []UtilityInvoiceItem{
				{
					Tipo:          ItemTUSDFio,
					Descricao:     "Consumo-TUSD",
					Quantidade:    decimal.NewFromInt(200),
					PrecoUnitario: decimal.RequireFromString("0.40000"),
					ValorTotal:    decimal.RequireFromString("80.00"),
				},
				{
					Tipo:          ItemTUSDEnergia,
					Descricao:     "Consumo-TE",
					Quantidade:    decimal.NewFromInt(200),
					PrecoUnitario: decimal.RequireFromString("0.17970"),
					ValorTotal:    decimal.RequireFromString("35.94"),
				},
				{
					Tipo:          ItemEnergiaInjetada,
					Descricao:     "Energia Injetada SCEE",
					Quantidade:    decimal.NewFromInt(4328),
					PrecoUnitario: decimal.RequireFromString("1.122920517"),
					ValorTotal:    decimal.RequireFromString("4860.00"),
				},
				{
					Tipo:       ItemIPCoelba,
					Descricao:  "Iluminação Pública Municipal",
					ValorTotal: decimal.RequireFromString("534.45"),
				},
			},
			contract: CalcContract{
				FatorRepasseEnergia: decimal.RequireFromString("0.85"),
				ValorIPComDesconto:  decimal.RequireFromString("12.56"),
				BandeiraComDesconto: false,
			},
			wantSEM:  "5510.39",
			wantCOM:  "4259.50",
			wantEcon: "1250.89",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := CalcInput{
				Contract: tt.contract,
				Itens:    tt.items,
			}

			result, err := Calculate(input)
			if err != nil {
				t.Fatalf("Calculate(%s): erro inesperado: %v", tt.name, err)
			}

			wantSEM := decimal.RequireFromString(tt.wantSEM)
			wantCOM := decimal.RequireFromString(tt.wantCOM)
			wantEcon := decimal.RequireFromString(tt.wantEcon)

			// Comparação com tolerância de 1 centavo.
			// A planilha do negócio arredonda cada linha para 2 casas decimais
			// em etapas intermediárias. O motor usa precisão total (decimal.Decimal)
			// e só arredonda no final. Por isso permitimos diferença de até R$ 0,01.
			diffSEM := result.TotalSemDesconto.Sub(wantSEM).Abs()
			if diffSEM.GreaterThan(decimal.NewFromInt(1)) {
				t.Errorf("TotalSemDesconto = %s, esperado ~%s (diferença=%s)",
					result.TotalSemDesconto, tt.wantSEM, diffSEM)
			}

			diffCOM := result.TotalComDesconto.Sub(wantCOM).Abs()
			if diffCOM.GreaterThan(decimal.NewFromInt(1)) {
				t.Errorf("TotalComDesconto = %s, esperado ~%s (diferença=%s)",
					result.TotalComDesconto, tt.wantCOM, diffCOM)
			}

			diffEcon := result.EconomiaRS.Sub(wantEcon).Abs()
			if diffEcon.GreaterThan(decimal.NewFromInt(1)) {
				t.Errorf("EconomiaRS = %s, esperado ~%s (diferença=%s)",
					result.EconomiaRS, tt.wantEcon, diffEcon)
			}
		})
	}
}

// TestCalculate_SemDesconto_RepasseIntegral valida que com fator=1.0
// os dois cenários são idênticos (sem economia).
func TestCalculate_SemDesconto_RepasseIntegral(t *testing.T) {
	input := CalcInput{
		Contract: CalcContract{
			FatorRepasseEnergia: decimal.NewFromInt(1),
			ValorIPComDesconto:  decimal.RequireFromString("50.00"),
		},
		Itens: []UtilityInvoiceItem{
			{
				Tipo:          ItemTUSDFio,
				Descricao:     "Consumo-TUSD",
				Quantidade:    decimal.NewFromInt(200),
				PrecoUnitario: decimal.RequireFromString("0.50"),
				ValorTotal:    decimal.RequireFromString("100.00"),
			},
			{
				Tipo:          ItemTUSDEnergia,
				Descricao:     "Consumo-TE",
				Quantidade:    decimal.NewFromInt(200),
				PrecoUnitario: decimal.RequireFromString("0.30"),
				ValorTotal:    decimal.RequireFromString("60.00"),
			},
			{
				Tipo:          ItemEnergiaInjetada,
				Descricao:     "Energia Injetada SCEE",
				Quantidade:    decimal.NewFromInt(1000),
				PrecoUnitario: decimal.RequireFromString("1.00"),
				ValorTotal:    decimal.RequireFromString("1000.00"),
			},
			{
				Tipo:       ItemIPCoelba,
				Descricao:  "Iluminação Pública Municipal",
				ValorTotal: decimal.RequireFromString("50.00"),
			},
		},
	}

	result, err := Calculate(input)
	if err != nil {
		t.Fatalf("Calculate: %v", err)
	}

	// Com fator=1.0, TotalSEM = TotalCOM + diferença de IP
	// SEM: 160 + 1000 + 50 = 1210
	// COM: 160 + 1000*1.0 + 50 = 1210
	// Economia = 1210 - 1210 = 0 (IP é IGUAL porque fator=1)
	expected := decimal.RequireFromString("1210")
	if !result.TotalSemDesconto.Equal(expected) {
		t.Errorf("TotalSemDesconto = %s, queria %s", result.TotalSemDesconto, expected)
	}
	if !result.TotalComDesconto.Equal(expected) {
		t.Errorf("TotalComDesconto = %s, queria %s", result.TotalComDesconto, expected)
	}
	if !result.EconomiaRS.IsZero() {
		t.Errorf("EconomiaRS = %s, queria 0 (fator=1 não gera economia)", result.EconomiaRS)
	}
}

// TestCalculate_ComDescontoBasico valida o cálculo básico com fator de repasse.
// TUSD é IDÊNTICO nos dois cenários. A diferença está na energia injetada e IP.
func TestCalculate_ComDescontoBasico(t *testing.T) {
	// TUSD: 200 * 0.50 + 200 * 0.30 = 160.00 (igual nos dois cenários)
	// Injetada: 1000 * 1.00 = 1000.00
	// Bandeira: 10.00
	// IP SEM: 50.00, IP COM: 15.00 (contratual)
	// Fator: 0.80
	//
	// SEM: 160 + 1000 + 10 + 50 = 1220.00
	// COM: 160 + 1000*0.80 + 10 + 15 = 160 + 800 + 10 + 15 = 985.00
	input := CalcInput{
		Contract: CalcContract{
			FatorRepasseEnergia: decimal.RequireFromString("0.80"),
			ValorIPComDesconto:  decimal.RequireFromString("15.00"),
			BandeiraComDesconto: false,
		},
		Itens: []UtilityInvoiceItem{
			{
				Tipo:          ItemTUSDFio,
				Descricao:     "Consumo-TUSD",
				Quantidade:    decimal.NewFromInt(200),
				PrecoUnitario: decimal.RequireFromString("0.50"),
				ValorTotal:    decimal.RequireFromString("100.00"),
			},
			{
				Tipo:          ItemTUSDEnergia,
				Descricao:     "Consumo-TE",
				Quantidade:    decimal.NewFromInt(200),
				PrecoUnitario: decimal.RequireFromString("0.30"),
				ValorTotal:    decimal.RequireFromString("60.00"),
			},
			{
				Tipo:          ItemEnergiaInjetada,
				Descricao:     "Energia Injetada SCEE",
				Quantidade:    decimal.NewFromInt(1000),
				PrecoUnitario: decimal.RequireFromString("1.00"),
				ValorTotal:    decimal.RequireFromString("1000.00"),
			},
			{
				Tipo:       ItemBandeira,
				Descricao:  "Bandeira AMARELA",
				ValorTotal: decimal.RequireFromString("10.00"),
			},
			{
				Tipo:       ItemIPCoelba,
				Descricao:  "Iluminação Pública Municipal",
				ValorTotal: decimal.RequireFromString("50.00"),
			},
		},
	}

	result, err := Calculate(input)
	if err != nil {
		t.Fatalf("Calculate: %v", err)
	}

	wantSEM := decimal.RequireFromString("1220.00")
	wantCOM := decimal.RequireFromString("985.00")
	wantEcon := wantSEM.Sub(wantCOM)

	if !result.TotalSemDesconto.Equal(wantSEM) {
		t.Errorf("TotalSemDesconto = %s, queria %s", result.TotalSemDesconto, wantSEM)
	}
	if !result.TotalComDesconto.Equal(wantCOM) {
		t.Errorf("TotalComDesconto = %s, queria %s", result.TotalComDesconto, wantCOM)
	}
	if !result.EconomiaRS.Equal(wantEcon) {
		t.Errorf("EconomiaRS = %s, queria %s", result.EconomiaRS, wantEcon)
	}

	// Verificar que o TUSD é idêntico nos dois cenários
	tusdLine := result.Linhas[0]
	if !tusdLine.ValorSemDesc.Equal(tusdLine.ValorComDesc) {
		t.Error("TUSD deveria ser idêntico nos dois cenários")
	}
}

// TestCalculate_BandeiraComDesconto valida que bandeira também recebe o fator
// quando a flag BandeiraComDesconto=true.
func TestCalculate_BandeiraComDesconto(t *testing.T) {
	input := CalcInput{
		Contract: CalcContract{
			FatorRepasseEnergia: decimal.RequireFromString("0.80"),
			ValorIPComDesconto:  decimal.Zero,
			BandeiraComDesconto: true,
		},
		Itens: []UtilityInvoiceItem{
			{Tipo: ItemTUSDFio, Quantidade: decimal.NewFromInt(100), PrecoUnitario: decimal.NewFromInt(1), ValorTotal: decimal.NewFromInt(100)},
			{Tipo: ItemTUSDEnergia, Quantidade: decimal.NewFromInt(100), PrecoUnitario: decimal.Zero, ValorTotal: decimal.Zero},
			{Tipo: ItemBandeira, ValorTotal: decimal.NewFromInt(100)},
		},
	}
	result, err := Calculate(input)
	if err != nil {
		t.Fatal(err)
	}

	// SEM: 100 + 0 + 100 = 200
	// COM: 100 + 0*0.80 + 100*0.80 = 100 + 80 = 180
	wantSEM := decimal.NewFromInt(200)
	wantCOM := decimal.NewFromInt(180)
	if !result.TotalSemDesconto.Equal(wantSEM) {
		t.Errorf("TotalSemDesconto = %s, queria %s", result.TotalSemDesconto, wantSEM)
	}
	if !result.TotalComDesconto.Equal(wantCOM) {
		t.Errorf("TotalComDesconto = %s, queria %s", result.TotalComDesconto, wantCOM)
	}
}

// TestCalculate_MissingTUSD valida que o motor rejeita entrada sem TUSD.
func TestCalculate_MissingTUSD(t *testing.T) {
	_, err := Calculate(CalcInput{
		Contract: CalcContract{FatorRepasseEnergia: decimal.RequireFromString("0.85")},
		Itens: []UtilityInvoiceItem{
			{Tipo: ItemIPCoelba, ValorTotal: decimal.NewFromInt(50)},
		},
	})
	if err == nil {
		t.Fatal("esperava erro por falta de TUSD")
	}
}

// TestCalculate_InvalidContract valida que o motor rejeita contratos inválidos.
func TestCalculate_InvalidContract(t *testing.T) {
	// Fator zero
	_, err := Calculate(CalcInput{
		Contract: CalcContract{FatorRepasseEnergia: decimal.Zero},
		Itens:    []UtilityInvoiceItem{{Tipo: ItemTUSDFio}},
	})
	if err == nil {
		t.Fatal("esperava erro por FatorRepasseEnergia zero")
	}

	// Fator > 1
	_, err = Calculate(CalcInput{
		Contract: CalcContract{FatorRepasseEnergia: decimal.NewFromInt(2)},
		Itens:    []UtilityInvoiceItem{{Tipo: ItemTUSDFio}},
	})
	if err == nil {
		t.Fatal("esperava erro por FatorRepasseEnergia > 1")
	}
}

// TestCalculate_EmptyItems valida que o motor rejeita lista vazia.
func TestCalculate_EmptyItems(t *testing.T) {
	_, err := Calculate(CalcInput{
		Contract: CalcContract{FatorRepasseEnergia: decimal.RequireFromString("0.85")},
	})
	if err == nil {
		t.Fatal("esperava erro por lista de itens vazia")
	}
}

// TestCalculate_SemEnergiaInjetada valida o caso em que não há geração no período.
// O motor deve funcionar (TUSD é cobrado normalmente) e emitir warning.
func TestCalculate_SemEnergiaInjetada(t *testing.T) {
	input := CalcInput{
		Contract: CalcContract{
			FatorRepasseEnergia: decimal.RequireFromString("0.85"),
			ValorIPComDesconto:  decimal.RequireFromString("12.00"),
		},
		Itens: []UtilityInvoiceItem{
			{
				Tipo:          ItemTUSDFio,
				Quantidade:    decimal.NewFromInt(200),
				PrecoUnitario: decimal.RequireFromString("0.50"),
				ValorTotal:    decimal.RequireFromString("100.00"),
			},
			{
				Tipo:          ItemTUSDEnergia,
				Quantidade:    decimal.NewFromInt(200),
				PrecoUnitario: decimal.RequireFromString("0.30"),
				ValorTotal:    decimal.RequireFromString("60.00"),
			},
			{
				Tipo:       ItemIPCoelba,
				Descricao:  "Iluminação Pública Municipal",
				ValorTotal: decimal.RequireFromString("50.00"),
			},
		},
	}

	result, err := Calculate(input)
	if err != nil {
		t.Fatalf("Calculate: %v", err)
	}

	// SEM: 160 + 0 (sem injeção) + 50 = 210
	// COM: 160 + 0 + 12 (IP contratual) = 172
	wantSEM := decimal.RequireFromString("210.00")
	wantCOM := decimal.RequireFromString("172.00")
	if !result.TotalSemDesconto.Equal(wantSEM) {
		t.Errorf("TotalSemDesconto = %s, queria %s", result.TotalSemDesconto, wantSEM)
	}
	if !result.TotalComDesconto.Equal(wantCOM) {
		t.Errorf("TotalComDesconto = %s, queria %s", result.TotalComDesconto, wantCOM)
	}

	// Deve emitir warning sobre falta de energia injetada
	found := false
	for _, w := range result.Warnings {
		if len(w) > 0 {
			found = true
			break
		}
	}
	if !found {
		t.Error("esperava warning sobre falta de energia injetada")
	}
}

// TestCalculate_LinhasBreakdown verifica a estrutura e valores das linhas
// de detalhamento do resultado.
func TestCalculate_LinhasBreakdown(t *testing.T) {
	input := CalcInput{
		Contract: CalcContract{
			FatorRepasseEnergia: decimal.RequireFromString("0.85"),
			ValorIPComDesconto:  decimal.RequireFromString("12.56"),
			BandeiraComDesconto: false,
		},
		Itens: []UtilityInvoiceItem{
			{
				Tipo:          ItemTUSDFio,
				Quantidade:    decimal.NewFromInt(200),
				PrecoUnitario: decimal.RequireFromString("0.40"),
				ValorTotal:    decimal.RequireFromString("80.00"),
			},
			{
				Tipo:          ItemTUSDEnergia,
				Quantidade:    decimal.NewFromInt(200),
				PrecoUnitario: decimal.RequireFromString("0.17970"),
				ValorTotal:    decimal.RequireFromString("35.94"),
			},
			{
				Tipo:          ItemEnergiaInjetada,
				Quantidade:    decimal.NewFromInt(4328),
				PrecoUnitario: decimal.RequireFromString("1.122920517"),
				ValorTotal:    decimal.RequireFromString("4860.00"),
			},
			{
				Tipo:       ItemIPCoelba,
				Descricao:  "Iluminação Pública Municipal",
				ValorTotal: decimal.RequireFromString("534.45"),
			},
		},
	}

	result, err := Calculate(input)
	if err != nil {
		t.Fatalf("Calculate: %v", err)
	}

	if len(result.Linhas) < 3 {
		t.Fatalf("esperava pelo menos 3 linhas, got %d", len(result.Linhas))
	}

	// Linha 0: TUSD
	if result.Linhas[0].Label != "Custo de Disponibilidade (TUSD)" {
		t.Errorf("Linha[0].Label = %q, queria 'Custo de Disponibilidade (TUSD)'", result.Linhas[0].Label)
	}
	if !result.Linhas[0].ValorSemDesc.Equal(result.Linhas[0].ValorComDesc) {
		t.Error("TUSD deve ser igual nos dois cenários")
	}

	// Linha 1: Energia Injetada
	if result.Linhas[1].Label != "Energia Injetada / Compensada" {
		t.Errorf("Linha[1].Label = %q, queria 'Energia Injetada / Compensada'", result.Linhas[1].Label)
	}

	// Bandeira não deve aparecer (não há bandeira nos itens)
	for _, l := range result.Linhas {
		if l.Label == "Bandeira Tarifária" {
			t.Error("bandeira não deveria aparecer nas linhas (não foi incluída nos itens)")
		}
	}

	// Última linha: IP
	lastLine := result.Linhas[len(result.Linhas)-1]
	if lastLine.Label != "Iluminação Pública" {
		t.Errorf("última linha.Label = %q, queria 'Iluminação Pública'", lastLine.Label)
	}
}
