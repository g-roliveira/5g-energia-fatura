package calcengine

import (
	"testing"

	"github.com/shopspring/decimal"
)

// TestCalculate_HappyPath valida o caminho feliz: TUSD fio + TUSD energia +
// injeção + bandeira + IP Coelba + desconto 15% (0.85).
func TestCalculate_HappyPath(t *testing.T) {
	in := CalculationInput{
		Contract: Contract{
			DescontoPct:                       decimal.RequireFromString("0.85"),
			IPFaturamentoMode:                 IPModeFixed,
			IPFaturamentoValor:                decimal.RequireFromString("10"),
			BandeiraComDesconto:               false,
			CustoDisponibilidadeSempreCobrado: true,
		},
		ConsumoMinimoKWh: 100,
		Itens: []UtilityInvoiceItem{
			{
				Type:          ItemTUSDFio,
				Description:   "Consumo-TUSD",
				Quantidade:    decimal.NewFromInt(500),
				PrecoUnitario: decimal.RequireFromString("0.75"),
				ValorTotal:    decimal.RequireFromString("375.00"),
			},
			{
				Type:          ItemTUSDEnergia,
				Description:   "Consumo-TE",
				Quantidade:    decimal.NewFromInt(500),
				PrecoUnitario: decimal.RequireFromString("0.37"),
				ValorTotal:    decimal.RequireFromString("185.00"),
			},
			{
				Type:          ItemEnergiaInjetada,
				Description:   "Energia Injetada SCEE",
				Quantidade:    decimal.NewFromInt(200),
				PrecoUnitario: decimal.RequireFromString("1.12"),
				ValorTotal:    decimal.RequireFromString("224.00"),
			},
			{
				Type:        ItemBandeira,
				Description: "Acrés. Band. AMARELA",
				ValorTotal:  decimal.RequireFromString("8.00"),
			},
			{
				Type:        ItemIPCoelba,
				Description: "Ilum. Púb. Municipal",
				ValorTotal:  decimal.RequireFromString("50.00"),
			},
		},
	}

	result, err := Calculate(in)
	if err != nil {
		t.Fatalf("Calculate: %v", err)
	}

	// Consumo líquido = 500 - 200 = 300 kWh
	// Energia sem desconto = 300 * (0.75 + 0.37) = 336.00
	// Energia com desconto = 336 * 0.85 = 285.60
	// Total sem = 336 + 8 (bandeira) + 50 (IP) = 394.00
	// Total com = 285.60 + 8 (bandeira repassada) + 50 (IP) + 10 (IP usina) = 353.60
	expectedSem := decimal.RequireFromString("394.00")
	expectedCom := decimal.RequireFromString("353.60")
	if !result.TotalSemDesconto.Equal(expectedSem) {
		t.Errorf("TotalSemDesconto = %s, queria %s", result.TotalSemDesconto, expectedSem)
	}
	if !result.TotalComDesconto.Equal(expectedCom) {
		t.Errorf("TotalComDesconto = %s, queria %s", result.TotalComDesconto, expectedCom)
	}
}

// TestCalculate_CustoDisponibilidade valida que o mínimo é aplicado quando
// a injeção zera o consumo líquido.
func TestCalculate_CustoDisponibilidade(t *testing.T) {
	in := CalculationInput{
		Contract: Contract{
			DescontoPct:                       decimal.RequireFromString("0.85"),
			IPFaturamentoMode:                 IPModeFixed,
			IPFaturamentoValor:                decimal.Zero,
			CustoDisponibilidadeSempreCobrado: true,
		},
		ConsumoMinimoKWh: 100,
		Itens: []UtilityInvoiceItem{
			{
				Type:          ItemTUSDFio,
				Quantidade:    decimal.NewFromInt(500),
				PrecoUnitario: decimal.RequireFromString("1.00"),
			},
			{
				Type:          ItemTUSDEnergia,
				Quantidade:    decimal.NewFromInt(500),
				PrecoUnitario: decimal.Zero,
			},
			{
				Type:       ItemEnergiaInjetada,
				Quantidade: decimal.NewFromInt(600), // injetou mais do que consumiu
			},
		},
	}

	result, err := Calculate(in)
	if err != nil {
		t.Fatalf("Calculate: %v", err)
	}

	// Consumo líquido = 500 - 600 = -100 → aplica mínimo = 100 kWh
	// Total sem = 100 * 1.00 = 100.00
	expected := decimal.RequireFromString("100.00")
	if !result.TotalSemDesconto.Equal(expected) {
		t.Errorf("TotalSemDesconto = %s, queria %s (mínimo deveria ter sido aplicado)",
			result.TotalSemDesconto, expected)
	}
	if len(result.Warnings) == 0 {
		t.Error("esperava warning sobre consumo mínimo aplicado")
	}
}

// TestCalculate_BandeiraComDesconto cobre o flag do contrato.
func TestCalculate_BandeiraComDesconto(t *testing.T) {
	in := CalculationInput{
		Contract: Contract{
			DescontoPct:         decimal.RequireFromString("0.80"),
			IPFaturamentoMode:   IPModeFixed,
			BandeiraComDesconto: true,
		},
		ConsumoMinimoKWh: 100,
		Itens: []UtilityInvoiceItem{
			{Type: ItemTUSDFio, Quantidade: decimal.NewFromInt(100), PrecoUnitario: decimal.NewFromInt(1)},
			{Type: ItemTUSDEnergia, Quantidade: decimal.NewFromInt(100), PrecoUnitario: decimal.Zero},
			{Type: ItemBandeira, ValorTotal: decimal.NewFromInt(100)},
		},
	}
	result, err := Calculate(in)
	if err != nil {
		t.Fatal(err)
	}
	// Total com = (100*1.00)*0.80 + 100*0.80 = 80 + 80 = 160
	expected := decimal.RequireFromString("160")
	if !result.TotalComDesconto.Equal(expected) {
		t.Errorf("bandeira não sofreu desconto quando deveria: got %s want %s",
			result.TotalComDesconto, expected)
	}
}

func TestCalculate_MissingTUSD(t *testing.T) {
	_, err := Calculate(CalculationInput{
		Contract: Contract{DescontoPct: decimal.RequireFromString("0.85")},
		Itens: []UtilityInvoiceItem{
			{Type: ItemIPCoelba, ValorTotal: decimal.NewFromInt(50)},
		},
	})
	if err == nil {
		t.Fatal("esperava erro por falta de TUSD")
	}
}

func TestCalculate_InvalidContract(t *testing.T) {
	_, err := Calculate(CalculationInput{
		Contract: Contract{DescontoPct: decimal.Zero},
		Itens:    []UtilityInvoiceItem{{Type: ItemTUSDFio}},
	})
	if err == nil {
		t.Fatal("esperava erro por DescontoPct zero")
	}

	_, err = Calculate(CalculationInput{
		Contract: Contract{DescontoPct: decimal.NewFromInt(2)},
		Itens:    []UtilityInvoiceItem{{Type: ItemTUSDFio}},
	})
	if err == nil {
		t.Fatal("esperava erro por DescontoPct > 1")
	}
}
