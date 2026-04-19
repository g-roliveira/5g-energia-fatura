// Package normalizer converte o BillingRecord produzido pelo pacote sync
// do backend-go em entradas canônicas aceitas pelo calc-engine.
//
// A responsabilidade principal é:
//
//   (a) classificar cada item de BillingRecord.ItensFatura (que vem como
//       []map[string]any com descrições livres em português) no ItemType
//       canônico que o motor entende (tusd_fio, tusd_energia, bandeira,
//       ip_coelba, ou ignorados: tributo_retido, reativo_excedente,
//       bandeira_verde);
//
//   (b) parsear valores em formato brasileiro ("0,75184631", "315,48")
//       para decimal.Decimal sem perda;
//
//   (c) extrair o bloco SCEE/MMGD do texto livre do rodapé "Informações
//       Importantes" — SCEE não aparece em itens_fatura em nenhum layout
//       real observado, só no rodapé.
//
// O normalizer NÃO depende do backend-go por design. O caller (código no
// internal/billing/ do backend-go) é quem monta o RawInvoice a partir do
// sync.BillingRecord.
package normalizer

import (
	calcengine "github.com/gustavo/5g-energia-fatura/packages/calc-engine"
	"github.com/shopspring/decimal"
)

// RawInvoice é a entrada minimamente estruturada que o normalizer aceita.
// O caller monta esse shape a partir do sync.BillingRecord do backend-go —
// esse contrato é propositalmente pequeno e desacoplado.
type RawInvoice struct {
	UC            string   `json:"uc"`
	NumeroFatura  string   `json:"numero_fatura"`
	MesReferencia string   `json:"mes_referencia"`
	Itens         []RawItem `json:"itens"`

	// InformacoesImportantes é o texto livre do rodapé do PDF onde o SCEE
	// aparece. Vem de document_record.ocr.informacoes_gerais do backend-go
	// (ou do equivalente). Opcional — se vazio, normalizer devolve SCEE=nil.
	InformacoesImportantes string `json:"informacoes_importantes,omitempty"`
}

// RawItem é um item da fatura como sai do sync.BillingRecord.ItensFatura
// (que é any mas na prática é []map[string]any com valores string em BR).
// Aqui já aceita os valores como string pra o caller simplesmente fazer
// cast dos maps.
type RawItem struct {
	Descricao    string `json:"descricao"`
	Quantidade   string `json:"quantidade"`
	Tarifa       string `json:"tarifa"`
	Valor        string `json:"valor"`
	ValorTotal   string `json:"valor_total"`
}

// Result é a saída do normalizer.
type Result struct {
	// Items são os itens aceitos pelo motor de cálculo.
	Items []calcengine.UtilityInvoiceItem `json:"items"`

	// IgnoredItems são itens classificados mas que não alimentam o cálculo
	// (TRIBF-IRRF, reativo excedente, bandeira verde). Preservados aqui
	// para auditoria e exibição na UI do backoffice.
	IgnoredItems []IgnoredItem `json:"ignored_items"`

	// UnclassifiedItems são linhas cuja descrição o classifier não
	// reconheceu. Requerem revisão manual — o cálculo não deve ser
	// aprovado com esses presentes.
	UnclassifiedItems []RawItem `json:"unclassified_items"`

	// SCEE é preenchido quando o texto do rodapé contém um dos três
	// layouts conhecidos (MMGD legado, MMGD transição, SCEE moderno).
	SCEE *SCEESummary `json:"scee,omitempty"`

	Warnings []string `json:"warnings,omitempty"`
}

type IgnoredItem struct {
	Type        calcengine.ItemType `json:"type"`
	Description string              `json:"description"`
	ValorTotal  decimal.Decimal     `json:"valor_total"`
	Reason      string              `json:"reason"`
}

// SCEELayout identifica qual geração de nomenclatura foi detectada.
type SCEELayout string

const (
	SCEELayoutMMGDLegado    SCEELayout = "mmgd_legado"
	SCEELayoutMMGDTransicao SCEELayout = "mmgd_transicao"
	SCEELayoutSCEEModerno   SCEELayout = "scee_moderno"
)

// SCEESummary é o resumo extraído do rodapé textual.
type SCEESummary struct {
	Layout             SCEELayout      `json:"layout"`
	EnergiaInjetadaKWh decimal.Decimal `json:"energia_injetada_kwh"`
	ExcedenteKWh       decimal.Decimal `json:"excedente_kwh"`
	CreditosUtilizados decimal.Decimal `json:"creditos_utilizados"`
	SaldoProximoCiclo  decimal.Decimal `json:"saldo_proximo_ciclo"`
	RawText            string          `json:"raw_text,omitempty"`
	Confidence         float64         `json:"confidence"`
}
