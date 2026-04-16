"""Testes do parser de PDF da Coelba.

Usa PDFs reais em tests/fixtures/ para garantir que o parser extrai corretamente
os campos que importam para o negócio da 5G Energia.
"""

from decimal import Decimal
from pathlib import Path

import pytest

from fatura.parser_pdf import CoelbaPdfParser

FIXTURES = Path(__file__).parent / "fixtures"
PDF_DEC24_A = FIXTURES / "333775733213.pdf"  # dez/2024, 433 kWh, 31 dias
PDF_DEC24_B = FIXTURES / "334075735546.pdf"  # dez/2024, 96 kWh, 6 dias


@pytest.fixture
def parser():
    return CoelbaPdfParser()


class TestParsePdfA:
    """PDF 333775733213 — dezembro/2024, ciclo completo (31 dias)."""

    @pytest.fixture(autouse=True)
    def conta(self, parser):
        self.conta = parser.parse(PDF_DEC24_A)

    def test_uc(self):
        assert self.conta.uc == "7085489032"

    def test_mes_ano(self):
        assert self.conta.mes == 12
        assert self.conta.ano == 2024

    def test_vencimento(self):
        from datetime import date
        assert self.conta.vencimento == date(2024, 12, 20)

    def test_valor_total_itens(self):
        # Total dos itens (valor real do consumo, ignorando crédito solar)
        assert self.conta.valor == Decimal("536.56")

    def test_numero_dias(self):
        assert self.conta.numero_dias == 31

    def test_cliente_nome(self):
        assert "PAULA" in self.conta.cliente.nome
        assert "FERNANDES" in self.conta.cliente.nome

    def test_classificacao(self):
        assert self.conta.cliente.classificacao is not None
        assert "B1" in self.conta.cliente.classificacao

    def test_medidor(self):
        assert self.conta.consumo is not None
        assert self.conta.consumo.medidor == "1204300816"

    def test_leituras(self):
        assert self.conta.consumo.leitura_anterior is not None
        assert self.conta.consumo.leitura_atual is not None

    def test_historico_nao_vazio(self):
        assert len(self.conta.historico_energia) >= 4

    def test_historico_periodo_formato(self):
        # Formato esperado: "DEZ/24"
        for h in self.conta.historico_energia:
            assert "/" in h.periodo, f"Período inesperado: {h.periodo}"

    def test_historico_kwh(self):
        periodos = {h.periodo: h.kwh for h in self.conta.historico_energia}
        assert "DEZ/24" in periodos
        assert periodos["DEZ/24"] == Decimal("433")

    def test_itens_count(self):
        assert len(self.conta.itens_fatura) >= 2

    def test_item_consumo_tusd(self):
        itens = {it.descricao: it for it in self.conta.itens_fatura}
        assert "Consumo-TUSD" in itens
        tusd = itens["Consumo-TUSD"]
        assert tusd.quantidade == Decimal("433")
        assert tusd.valor == Decimal("305.27")

    def test_item_consumo_te(self):
        itens = {it.descricao: it for it in self.conta.itens_fatura}
        assert "Consumo-TE" in itens
        te = itens["Consumo-TE"]
        assert te.quantidade == Decimal("433")
        assert te.valor == Decimal("169.61")

    def test_pdf_path(self):
        assert self.conta.pdf_path is not None
        assert "333775733213" in self.conta.pdf_path

    def test_parsed_at(self):
        assert self.conta.parsed_at is not None


class TestParsePdfB:
    """PDF 334075735546 — dezembro/2024, ciclo curto (6 dias)."""

    @pytest.fixture(autouse=True)
    def conta(self, parser):
        self.conta = parser.parse(PDF_DEC24_B)

    def test_uc(self):
        assert self.conta.uc == "7085489032"

    def test_mes_ano(self):
        assert self.conta.mes == 12
        assert self.conta.ano == 2024

    def test_vencimento(self):
        from datetime import date
        assert self.conta.vencimento == date(2025, 1, 8)

    def test_valor_total_itens(self):
        assert self.conta.valor == Decimal("113.12")

    def test_numero_dias(self):
        assert self.conta.numero_dias == 6

    def test_consumo_kwh_historico(self):
        periodos = {h.periodo: h.kwh for h in self.conta.historico_energia}
        assert "DEZ/24" in periodos
        # 96 kWh no ciclo de 6 dias, mas histórico mostra acumulado DEZ24 = 529
        assert periodos["DEZ/24"] == Decimal("529")

    def test_itens_consumo_tusd_96kwh(self):
        itens = {it.descricao: it for it in self.conta.itens_fatura}
        assert "Consumo-TUSD" in itens
        assert itens["Consumo-TUSD"].quantidade == Decimal("96")

    def test_sem_banda_amarela(self):
        descricoes = [it.descricao for it in self.conta.itens_fatura]
        assert not any("Band" in d for d in descricoes)


class TestParserEdgeCases:
    def test_arquivo_inexistente(self, parser):
        from fatura.exceptions import ParserError
        with pytest.raises(ParserError, match="não encontrado"):
            parser.parse("/tmp/nao_existe.pdf")
