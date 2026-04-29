package billing

import "github.com/shopspring/decimal"

// ItemType classifica um item da fatura para que o motor saiba como tratá-lo.
type ItemType string

const (
	ItemTUSDFio          ItemType = "tusd_fio"          // Consumo-TUSD (parte fio, R$/kWh)
	ItemTUSDEnergia      ItemType = "tusd_energia"      // Consumo-TE (parte energia, R$/kWh)
	ItemEnergiaInjetada  ItemType = "energia_injetada"  // kWh injetados (compensação SCEE/MMGD)
	ItemBandeira         ItemType = "bandeira"          // acréscimo de bandeira tarifária (R$)
	ItemIPCoelba         ItemType = "ip_coelba"         // Iluminação Pública Municipal (R$)
	ItemReativoExcedente ItemType = "reativo_excedente" // não entra no cálculo
	ItemTributoRetido    ItemType = "tributo_retido"    // IRRF etc — não entra no cálculo
)

// UtilityInvoiceItem é um item normalizado da fatura da distribuidora.
type UtilityInvoiceItem struct {
	Type          ItemType        `json:"type"`
	Description   string          `json:"description"`
	Quantidade    decimal.Decimal `json:"quantidade"`
	PrecoUnitario decimal.Decimal `json:"preco_unitario"`
	ValorTotal    decimal.Decimal `json:"valor_total"`
}

// IPFaturamentoMode controla como a IP da usina é aplicada.
type IPFaturamentoMode string

const (
	IPModeFixed   IPFaturamentoMode = "fixed"
	IPModePercent IPFaturamentoMode = "percent"
)

// CalcContract é o snapshot imutável do contrato usado no cálculo.
type CalcContract struct {
	DescontoPct                       decimal.Decimal   `json:"desconto_pct"`
	IPFaturamentoMode                 IPFaturamentoMode `json:"ip_faturamento_mode"`
	IPFaturamentoValor                decimal.Decimal   `json:"ip_faturamento_valor"`
	IPFaturamentoPct                  decimal.Decimal   `json:"ip_faturamento_pct"`
	BandeiraComDesconto               bool              `json:"bandeira_com_desconto"`
	CustoDisponibilidadeSempreCobrado bool              `json:"custo_disponibilidade_sempre_cobrado"`
}

// CalculationInput é a entrada completa do motor.
type CalculationInput struct {
	Contract         CalcContract         `json:"contract"`
	Itens            []UtilityInvoiceItem `json:"itens"`
	ConsumoMinimoKWh float64              `json:"consumo_minimo_kwh"`
}

// LineBreakdown é uma linha do detalhamento do resultado.
type LineBreakdown struct {
	Label         string          `json:"label"`
	Quantidade    decimal.Decimal `json:"quantidade,omitempty"`
	PrecoUnitario decimal.Decimal `json:"preco_unitario,omitempty"`
	ValorSemDesc  decimal.Decimal `json:"valor_sem_desconto"`
	ValorComDesc  decimal.Decimal `json:"valor_com_desconto"`
}

// CalculationResult é a saída determinística do motor.
type CalculationResult struct {
	TotalSemDesconto decimal.Decimal `json:"total_sem_desconto"`
	TotalComDesconto decimal.Decimal `json:"total_com_desconto"`
	EconomiaRS       decimal.Decimal `json:"economia_rs"`
	EconomiaPct      decimal.Decimal `json:"economia_pct"`
	Linhas           []LineBreakdown `json:"linhas"`
	Warnings         []string        `json:"warnings,omitempty"`
}
