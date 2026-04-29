package billing

import "github.com/shopspring/decimal"

// ItemType classifica um item da fatura para que o motor saiba como tratá-lo.
type ItemType string

const (
	ItemTUSDFio          ItemType = "tusd_fio"           // Consumo-TUSD (parte fio, R$/kWh)
	ItemTUSDEnergia      ItemType = "tusd_energia"       // Consumo-TE (parte energia, R$/kWh)
	ItemEnergiaInjetada  ItemType = "energia_injetada"   // kWh injetados (compensação SCEE/MMGD)
	ItemBandeira         ItemType = "bandeira"           // acréscimo de bandeira tarifária (R$)
	ItemIPCoelba         ItemType = "iluminacao_publica" // Iluminação Pública Municipal (R$)
	ItemReativoExcedente ItemType = "reativo_excedente"  // não entra no cálculo
	ItemTributoRetido    ItemType = "tributo_retido"     // IRRF etc — não entra no cálculo
)

// UtilityInvoiceItem é um item normalizado da fatura da distribuidora.
type UtilityInvoiceItem struct {
	Tipo          ItemType        `json:"tipo"`
	Descricao     string          `json:"descricao"`
	Quantidade    decimal.Decimal `json:"quantidade"`
	PrecoUnitario decimal.Decimal `json:"preco_unitario"`
	ValorTotal    decimal.Decimal `json:"valor_total"`
}

// CalcContract é o snapshot imutável do contrato usado no cálculo.
type CalcContract struct {
	FatorRepasseEnergia               decimal.Decimal `json:"fator_repasse_energia"`
	ValorIPComDesconto                decimal.Decimal `json:"valor_ip_com_desconto"`
	BandeiraComDesconto               bool            `json:"bandeira_com_desconto"`
	CustoDisponibilidadeSempreCobrado bool            `json:"custo_disponibilidade_sempre_cobrado"`
	ConsumoMinimoKWh                  decimal.Decimal `json:"consumo_minimo_kwh"`
}

// CalcInput é a entrada completa do motor.
type CalcInput struct {
	Contract CalcContract         `json:"contract"`
	Itens    []UtilityInvoiceItem `json:"itens"`
}

// LineBreakdown é uma linha do detalhamento do resultado.
type LineBreakdown struct {
	Label         string          `json:"label"`
	Quantidade    decimal.Decimal `json:"quantidade,omitempty"`
	PrecoUnitario decimal.Decimal `json:"preco_unitario,omitempty"`
	ValorSemDesc  decimal.Decimal `json:"valor_sem_desconto"`
	ValorComDesc  decimal.Decimal `json:"valor_com_desconto"`
}

// CalcResult é a saída determinística do motor.
type CalcResult struct {
	TotalSemDesconto decimal.Decimal `json:"total_sem_desconto"`
	TotalComDesconto decimal.Decimal `json:"total_com_desconto"`
	EconomiaRS       decimal.Decimal `json:"economia_rs"`
	EconomiaPct      decimal.Decimal `json:"economia_pct"`
	Linhas           []LineBreakdown `json:"linhas"`
	Warnings         []string        `json:"warnings,omitempty"`
}
