package normalizer

import (
	"strings"

	calcengine "github.com/gustavo/5g-energia-fatura/packages/calc-engine"
)

type classification struct {
	Type       calcengine.ItemType
	Confidence float64
	MatchedBy  string
	Ignore     bool
}

// classify tenta descobrir o tipo do item a partir da descrição livre.
// Retorna ok=false quando não é capaz de classificar.
//
// EVIDÊNCIAS DOCUMENTAIS (faturas Coelba/Neoenergia reais examinadas):
//
//   (A) UC 007098175908 / 04/2026 (resposta real da API Go do backend):
//     "Consumo-TUSD", "Consumo-TE", "Ilum. Púb. Municipal"
//
//   (B) MP-BA 40101.0013.23.0000057-4 / 10/2023 (MMGD legado):
//     "Consumo-TUSD", "BANDEIRA VERDE", "Consumo-TE",
//     "Cons.Reat.Excedente", "Ilum. Púb. Municipal", "TRIBF-IRRF(1.2%)"
//
//   (C) MP-BA 40101.0013.24.0000047-3 / 07/2024 (MMGD transição):
//     "Consumo-TUSD", "BANDEIRA AMARELA", "Consumo-TE", "BANDEIRA VERDE",
//     "Consumo-TE" (segundo), "Cons.Reat.Excedente", "Ilum. Púb. Municipal"
//     Obs: duas bandeiras no mesmo ciclo quando bandeira ANEEL muda no mês.
//
//   (D) MP-BA 40101.0055.26.0000006-5 / 12/2025 (SCEE moderno):
//     "Consumo-TUSD", "Consumo-TE", "Acrés. Band. AMARELA",
//     "Acrés. Band. VERMELHA", "Ilum. Púb. Municipal"
//
//   (E) Planilha Azi Dourado out/2025: "Acrés. Band. VERMELHA- P2" (posto 2)
//
// Três gerações de rótulo pra bandeira: "BANDEIRA <COR>" (legado),
// "Acrés. Band. <COR>" (novo), "Acrés. Band. <COR>- P2" (posto).
func classify(descricao string) (classification, bool) {
	n := normalizeForMatch(descricao)

	// 1) Tributos retidos (IRRF etc) — ignore
	if containsAny(n, "tribf-irrf", "tribf irrf", "irrf", "retencao", "tributo retido") {
		return classification{
			Type: calcengine.ItemTributoRetido, Confidence: 0.95,
			MatchedBy: "observed:TRIBF-IRRF", Ignore: true,
		}, true
	}

	// 2) Reativo excedente — ignore
	if containsAny(n,
		"cons.reat", "cons reat", "reativo excedente", "reat.exc", "reat exc",
		"cons.real.excedente", "cons.real.exced", "cons real exc", "cons.real exc",
		"cons.real.exc.nponta", "cons.real exc.fponta", "cons.real.exc.fponta") {
		return classification{
			Type: calcengine.ItemReativoExcedente, Confidence: 0.95,
			MatchedBy: "observed:Cons.Reat.Excedente", Ignore: true,
		}, true
	}

	// 2.1) Demandas e impostos de grupo A — fora do cálculo atual.
	// Mantemos como ignored para auditoria enquanto o motor de grupo A
	// não estiver implementado.
	if containsAny(n, "demanda ativa", "demanda reativa", "imp.som/dim-c/impost", "imp som/dim-c/impost") {
		return classification{
			Type: calcengine.ItemTributoRetido, Confidence: 0.85,
			MatchedBy: "observed:grupoA_demanda_ou_imposto", Ignore: true,
		}, true
	}

	// 3) Energia injetada / SCEE (raro como item, normalmente no rodapé)
	if containsAny(n, "injetad", "compensad", "scee", "credito energia",
		"credito de energia", "energia ativa injetada") {
		return classification{
			Type: calcengine.ItemEnergiaInjetada, Confidence: 0.85,
			MatchedBy: "heuristic:injetada|compensada|scee",
		}, true
	}

	// 4) Bandeiras
	if containsAny(n, "bandeira verde") {
		return classification{
			Type: calcengine.ItemBandeira, Confidence: 0.95,
			MatchedBy: "observed:BANDEIRA VERDE (sem cobrança)", Ignore: true,
		}, true
	}
	if containsAny(n,
		"bandeira amarela", "bandeira vermelha",
		"acres. band", "acresc. band", "acres band",
		"acres, band", "acres, bend", "acresc, bend", "acres bend",
		"bandeira tarifaria", "adicional band") {
		return classification{
			Type: calcengine.ItemBandeira, Confidence: 0.95,
			MatchedBy: "observed:Bandeira Amarela/Vermelha",
		}, true
	}

	// 5) Iluminação Pública Municipal
	if containsAny(n,
		"ilum. pub", "ilum pub", "iluminacao publica",
		"cip municipal", "contrib. ilum", "contribuicao ilum",
		"hum. pub", "hum pub", "itum. pub", "itum pub") {
		return classification{
			Type: calcengine.ItemIPCoelba, Confidence: 0.95,
			MatchedBy: "observed:Ilum. Púb. Municipal",
		}, true
	}

	// 6) Consumo-TE (antes de TUSD pra não confundir)
	if containsAny(n,
		"consumo-te", "consumo te",
		"consumo ativo te", "tarifa de energia",
		"consumo ativo na ponta", "consumo ativo fora",
		"tarifa energia") {
		return classification{
			Type: calcengine.ItemTUSDEnergia, Confidence: 0.95,
			MatchedBy: "observed:Consumo-TE",
		}, true
	}

	// 7) Consumo-TUSD
	if containsAny(n,
		"consumo-tusd", "consumo tusd",
		"consumo ativo tusd", "consumo ativo-tusd",
		"tusd") {
		return classification{
			Type: calcengine.ItemTUSDFio, Confidence: 0.95,
			MatchedBy: "observed:Consumo-TUSD",
		}, true
	}

	return classification{}, false
}

func normalizeForMatch(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	return removePortugueseAccents(s)
}

// removePortugueseAccents é uma tabela manual porque o ambiente-alvo
// (algumas VMs corporativas) pode não ter golang.org/x/text disponível.
func removePortugueseAccents(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case 'á', 'à', 'â', 'ã', 'ä', 'Á', 'À', 'Â', 'Ã', 'Ä':
			b.WriteRune('a')
		case 'é', 'è', 'ê', 'ë', 'É', 'È', 'Ê', 'Ë':
			b.WriteRune('e')
		case 'í', 'ì', 'î', 'ï', 'Í', 'Ì', 'Î', 'Ï':
			b.WriteRune('i')
		case 'ó', 'ò', 'ô', 'õ', 'ö', 'Ó', 'Ò', 'Ô', 'Õ', 'Ö':
			b.WriteRune('o')
		case 'ú', 'ù', 'û', 'ü', 'Ú', 'Ù', 'Û', 'Ü':
			b.WriteRune('u')
		case 'ç', 'Ç':
			b.WriteRune('c')
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func containsAny(haystack string, needles ...string) bool {
	for _, n := range needles {
		if strings.Contains(haystack, n) {
			return true
		}
	}
	return false
}
