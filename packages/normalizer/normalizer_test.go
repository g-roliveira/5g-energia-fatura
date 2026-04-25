package normalizer

import (
	"testing"

	"github.com/shopspring/decimal"

	calcengine "github.com/gustavo/5g-energia-fatura/packages/calc-engine"
)

func TestClassify_RealPDFStrings(t *testing.T) {
	type tc struct {
		desc       string
		wantType   calcengine.ItemType
		wantOK     bool
		wantIgnore bool
		fromPDF    string
	}
	cases := []tc{
		// Paula UC 007098175908 (resposta real da API Go)
		{"Consumo-TUSD", calcengine.ItemTUSDFio, true, false, "Paula API"},
		{"Consumo-TE", calcengine.ItemTUSDEnergia, true, false, "Paula API"},
		{"Ilum. Púb. Municipal", calcengine.ItemIPCoelba, true, false, "Paula API"},

		// MP-BA 57-4 (MMGD legado)
		{"BANDEIRA VERDE", calcengine.ItemBandeira, true, true, "MP-BA 57-4"},
		{"Cons.Reat.Excedente", calcengine.ItemReativoExcedente, true, true, "MP-BA 57-4"},
		{"TRIBF-IRRF(1.2%)", calcengine.ItemTributoRetido, true, true, "MP-BA 57-4"},

		// MP-BA 47-3 (MMGD transição)
		{"BANDEIRA AMARELA", calcengine.ItemBandeira, true, false, "MP-BA 47-3"},

		// MP-BA 6-5 (SCEE moderno)
		{"Acrés. Band. AMARELA", calcengine.ItemBandeira, true, false, "MP-BA 6-5"},
		{"Acrés. Band. VERMELHA", calcengine.ItemBandeira, true, false, "MP-BA 6-5"},
		{"Acrés, Bend. AMARELA", calcengine.ItemBandeira, true, false, "OCR typo"},
		{"Acrés, Bend. VERMELHA", calcengine.ItemBandeira, true, false, "OCR typo"},
		{"Hum. Púb. Municipal", calcengine.ItemIPCoelba, true, false, "OCR typo"},
		{"Itum. Púb. Municipal", calcengine.ItemIPCoelba, true, false, "OCR typo"},
		{"Cons.Real.Excedente", calcengine.ItemReativoExcedente, true, true, "OCR variant"},
		{"Cons.Real.Exc.NPonta", calcengine.ItemReativoExcedente, true, true, "OCR variant"},
		{"Cons.Real Exc.FPonta", calcengine.ItemReativoExcedente, true, true, "OCR variant"},
		{"Demanda Ativa", calcengine.ItemTributoRetido, true, true, "Grupo A (não calculado)"},
		{"Demanda Reativa Exc.", calcengine.ItemTributoRetido, true, true, "Grupo A (não calculado)"},
		{"Imp.Som/Dim-C/Impost", calcengine.ItemTributoRetido, true, true, "Grupo A (não calculado)"},

		// Planilha Azi Dourado
		{"Acrés. Band. VERMELHA- P2", calcengine.ItemBandeira, true, false, "Azi Dourado"},

		// Negativos
		{"Taxa XYZ Inexistente", "", false, false, "fora de escopo"},
		{"", "", false, false, "vazio"},
	}

	for _, c := range cases {
		t.Run(c.desc+" ("+c.fromPDF+")", func(t *testing.T) {
			got, ok := classify(c.desc)
			if ok != c.wantOK {
				t.Fatalf("ok=%v, queria %v", ok, c.wantOK)
			}
			if !ok {
				return
			}
			if got.Type != c.wantType {
				t.Errorf("type=%s, queria %s", got.Type, c.wantType)
			}
			if got.Ignore != c.wantIgnore {
				t.Errorf("ignore=%v, queria %v", got.Ignore, c.wantIgnore)
			}
		})
	}
}

func TestExtractSCEE_RealPDFText(t *testing.T) {
	type tc struct {
		name         string
		text         string
		wantLayout   SCEELayout
		wantInjetada string
		wantExc      string
		wantCred     string
		wantSaldo    string
	}
	cases := []tc{
		{
			name: "MMGD_legado_57-4",
			text: "Unidade Microgeracao. Energia injetada no mes 8736 kWh. " +
				"Saldo total de credito para o proximo faturamento 46952 kWh.",
			wantLayout:   SCEELayoutMMGDLegado,
			wantInjetada: "8736",
			wantSaldo:    "46952",
		},
		{
			name: "MMGD_transicao_47-3",
			text: "Unidade Microgeracao. Energia injetada no mes 6248 kWh. " +
				"Saldo total de credito para o proximo faturamento 129108 kWh",
			wantLayout:   SCEELayoutMMGDLegado,
			wantInjetada: "6248",
			wantSaldo:    "129108",
		},
		{
			name: "SCEE_moderno_6-5",
			text: "SCEE, Excedente 3170 kWh e creditos utilizados 0 kWh. " +
				"Saldo para o proximo ciclo 0 kWh.",
			wantLayout: SCEELayoutSCEEModerno,
			wantExc:    "3170",
			wantCred:   "0",
			wantSaldo:  "0",
		},
		{
			name:       "sem_SCEE",
			text:       "Bandeira em vigor é Verde. Qualquer outra coisa.",
			wantLayout: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ExtractSCEE(c.text)
			if c.wantLayout == "" {
				if got != nil {
					t.Fatalf("esperava nil, recebi %+v", got)
				}
				return
			}
			if got == nil {
				t.Fatal("esperava SCEE, recebi nil")
			}
			if got.Layout != c.wantLayout {
				t.Errorf("layout=%s, queria %s", got.Layout, c.wantLayout)
			}
			check := func(name, want string, got decimal.Decimal) {
				if want == "" {
					return
				}
				if !got.Equal(decimal.RequireFromString(want)) {
					t.Errorf("%s=%s, queria %s", name, got, want)
				}
			}
			check("injetada", c.wantInjetada, got.EnergiaInjetadaKWh)
			check("exc", c.wantExc, got.ExcedenteKWh)
			check("cred", c.wantCred, got.CreditosUtilizados)
			check("saldo", c.wantSaldo, got.SaldoProximoCiclo)
		})
	}
}

// TestNormalize_EndToEnd_Paula monta um RawInvoice espelhando o que sai
// da fatura real da Paula (UC 007098175908, abril/2026) e verifica que
// o normalizer produz 3 items prontos pro motor + nenhum ignorado.
func TestNormalize_EndToEnd_Paula(t *testing.T) {
	raw := RawInvoice{
		UC:            "007098175908",
		NumeroFatura:  "339800707843",
		MesReferencia: "2026/04",
		Itens: []RawItem{
			{Descricao: "Consumo-TUSD", Quantidade: "418", Tarifa: "0,75478", Valor: "315,48", ValorTotal: "315,48"},
			{Descricao: "Consumo-TE", Quantidade: "418", Tarifa: "0,37333", Valor: "156,05", ValorTotal: "156,05"},
			{Descricao: "Ilum. Púb. Municipal", Valor: "50,00", ValorTotal: "50,00"},
		},
	}
	result, err := Normalize(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 3 {
		t.Fatalf("esperava 3 itens válidos, recebi %d", len(result.Items))
	}
	if len(result.IgnoredItems) != 0 {
		t.Errorf("não esperava ignored, recebi %d", len(result.IgnoredItems))
	}
	if len(result.UnclassifiedItems) != 0 {
		t.Errorf("não esperava unclassified, recebi %d", len(result.UnclassifiedItems))
	}

	// Soma dos valores bate com 521.53 (valor total da fatura)
	total := decimal.Zero
	for _, item := range result.Items {
		total = total.Add(item.ValorTotal)
	}
	expected := decimal.RequireFromString("521.53")
	if !total.Equal(expected) {
		t.Errorf("soma=%s, queria %s", total, expected)
	}
}

// TestNormalize_WithIgnoredItems valida que IRRF e reativo vão pra
// IgnoredItems mas não aparecem nos Items do motor.
func TestNormalize_WithIgnoredItems(t *testing.T) {
	raw := RawInvoice{
		UC: "123",
		Itens: []RawItem{
			{Descricao: "Consumo-TUSD", Quantidade: "100", Valor: "75,00", ValorTotal: "75,00"},
			{Descricao: "Consumo-TE", Quantidade: "100", Valor: "37,00", ValorTotal: "37,00"},
			{Descricao: "BANDEIRA VERDE", ValorTotal: "0,00"},
			{Descricao: "Cons.Reat.Excedente", ValorTotal: "2,50"},
			{Descricao: "TRIBF-IRRF(1.2%)", ValorTotal: "-3,00"},
		},
	}
	result, err := Normalize(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 2 {
		t.Errorf("esperava 2 itens no motor (TUSD+TE), recebi %d", len(result.Items))
	}
	if len(result.IgnoredItems) != 3 {
		t.Errorf("esperava 3 ignored (verde+reativo+irrf), recebi %d", len(result.IgnoredItems))
	}
}

func TestNormalize_WithUnclassifiedItem(t *testing.T) {
	raw := RawInvoice{
		UC: "123",
		Itens: []RawItem{
			{Descricao: "Consumo-TUSD", Quantidade: "100", Valor: "75,00", ValorTotal: "75,00"},
			{Descricao: "Taxa Esquisita Nova 2027", ValorTotal: "99,99"},
		},
	}
	result, _ := Normalize(raw)
	if len(result.UnclassifiedItems) != 1 {
		t.Fatalf("esperava 1 unclassified, recebi %d", len(result.UnclassifiedItems))
	}
	if len(result.Warnings) == 0 {
		t.Error("esperava warning sobre item não classificado")
	}
}

func TestParseBRDecimal(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"0,75184631", "0.75184631"},
		{"1.234,56", "1234.56"},
		{"R$ 315,48", "315.48"},
		{"", "0"},
		{"  ", "0"},
	}
	for _, c := range cases {
		got, err := parseBRDecimal(c.in)
		if err != nil {
			t.Errorf("parseBRDecimal(%q): %v", c.in, err)
			continue
		}
		if !got.Equal(decimal.RequireFromString(c.want)) {
			t.Errorf("parseBRDecimal(%q) = %s, queria %s", c.in, got, c.want)
		}
	}
}
