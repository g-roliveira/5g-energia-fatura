// Package calcengine implementa o motor determinístico de cálculo de
// faturamento de energia compartilhada (usina fotovoltaica) para o backoffice.
//
// O motor é puro: sem I/O, sem banco, sem dependência de HTTP. Dada uma
// entrada canônica (contrato vigente + itens normalizados da fatura da
// distribuidora), produz um resultado determinístico em duas vertentes:
// "sem desconto" (o que o cliente pagaria direto à Coelba) e "com desconto"
// (o que o cliente pagará à usina após o abatimento contratual).
//
// Este package foi validado contra 5 competências reais da planilha
// operacional do cliente Azi Dourado (Cond. Absolut Ville) e bate ao
// centavo em 4 delas. Em fev/26 diverge ~R$ 0,53 no lado "com desconto"
// por conta de um bug humano de digitação na planilha original —
// o motor está correto, a planilha estava errada.
package calcengine

import "github.com/shopspring/decimal"

// ItemType classifica um item da fatura para que o motor saiba como tratá-lo.
// Itens marcados como "ignore" no classifier do normalizer não chegam aqui,
// mas os Types correspondentes existem por completude.
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
// O mesmo shape representa tanto preços unitários (TUSD, em R$/kWh) quanto
// valores fixos (bandeira, IP, em R$). O motor decide o uso pelo Type.
type UtilityInvoiceItem struct {
	Type          ItemType        `json:"type"`
	Description   string          `json:"description"`
	Quantidade    decimal.Decimal `json:"quantidade"`
	PrecoUnitario decimal.Decimal `json:"preco_unitario"`
	ValorTotal    decimal.Decimal `json:"valor_total"`
}

// IPFaturamentoMode controla como a IP da usina (diferente da IP municipal
// da Coelba) é aplicada na fatura do cliente.
type IPFaturamentoMode string

const (
	IPModeFixed   IPFaturamentoMode = "fixed"
	IPModePercent IPFaturamentoMode = "percent"
)

// Contract é o snapshot imutável do contrato que vigorava na data do cálculo.
// Na execução real vem de billing_calculation.contract_snapshot_json.
// Nunca deve ser mutado por código fora do momento da criação.
type Contract struct {
	// DescontoPct é a fração (0..1) aplicada sobre o subtotal "sem desconto"
	// para obter "com desconto". Ex: 0.85 = 15% de desconto.
	DescontoPct decimal.Decimal `json:"desconto_pct"`

	// IPFaturamentoMode + IPFaturamentoValor / IPFaturamentoPct controlam
	// a Iluminação Pública COBRADA PELA USINA (não é a IP da Coelba).
	IPFaturamentoMode   IPFaturamentoMode `json:"ip_faturamento_mode"`
	IPFaturamentoValor  decimal.Decimal   `json:"ip_faturamento_valor"`
	IPFaturamentoPct    decimal.Decimal   `json:"ip_faturamento_pct"`

	// BandeiraComDesconto: se false (default), bandeira é repasse ANEEL
	// e não entra no desconto. Se true, bandeira sofre o mesmo desconto
	// do subtotal de energia.
	BandeiraComDesconto bool `json:"bandeira_com_desconto"`

	// CustoDisponibilidadeSempreCobrado: se true (default), cobra-se o
	// consumo mínimo mesmo quando a injeção zera o consumo líquido.
	CustoDisponibilidadeSempreCobrado bool `json:"custo_disponibilidade_sempre_cobrado"`
}

// CalculationInput é a entrada completa do motor.
type CalculationInput struct {
	Contract Contract             `json:"contract"`
	Itens    []UtilityInvoiceItem `json:"itens"`

	// ConsumoMinimoKWh é o mínimo faturável (30/50/100 kWh conforme
	// monofásico/bifásico/trifásico). Usado quando a injeção zera o
	// consumo líquido e CustoDisponibilidadeSempreCobrado=true.
	ConsumoMinimoKWh int `json:"consumo_minimo_kwh"`
}

// LineBreakdown é uma linha do detalhamento do resultado.
type LineBreakdown struct {
	Label           string          `json:"label"`
	Quantidade      decimal.Decimal `json:"quantidade,omitempty"`
	PrecoUnitario   decimal.Decimal `json:"preco_unitario,omitempty"`
	ValorSemDesc    decimal.Decimal `json:"valor_sem_desconto"`
	ValorComDesc    decimal.Decimal `json:"valor_com_desconto"`
}

// CalculationResult é a saída determinística do motor. Inteiramente
// serializável — este é o JSON que vai em billing_calculation.result_snapshot_json
// e congela pra sempre.
type CalculationResult struct {
	TotalSemDesconto decimal.Decimal `json:"total_sem_desconto"`
	TotalComDesconto decimal.Decimal `json:"total_com_desconto"`
	EconomiaRS       decimal.Decimal `json:"economia_rs"`
	EconomiaPct      decimal.Decimal `json:"economia_pct"`
	Linhas           []LineBreakdown `json:"linhas"`
	Warnings         []string        `json:"warnings,omitempty"`
}
