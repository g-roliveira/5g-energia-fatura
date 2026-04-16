"""Parser de faturas PDF da Coelba/Neoenergia.

Estrutura real do PDF (confirmada com pdfplumber):
  - Página 0 contém 3 tabelas:
    * Tabela 0 (2 linhas): CÓDIGO DA INSTALAÇÃO e CÓDIGO DO CLIENTE (= UC)
    * Tabela 1 (grande): cabeçalho da fatura, itens, histórico, medidor
    * Tabela 2 (boleto): PAGADOR, CÓDIGO DO CLIENTE, VENCIMENTO, VALOR DO DOCUMENTO

Layout da Tabela 1 (buscas por conteúdo, não por índice):
  - Linha com "REF:MÊS/ANO": mês/ano ref, TOTAL A PAGAR, vencimento
  - Linha com "CLASSIFICAÇÃO:": classificação tarifária
  - Linha com "DATAS DE LEITURAS": leitura anterior/atual, N° DE DIAS
  - Linha com "ITENS DA FATURA": cabeçalho de itens (próxima linha = dados)
  - Linha com "CONSUMO FATURADO": histórico de consumo kWh
  - Linha com "TOTAL" (célula [0]): total da fatura (soma de itens)
  - Linha com medidor (10 dígitos): dados do medidor

Células com múltiplos itens usam '\\n' como separador dentro da mesma célula.
"""

import re
from datetime import date, datetime
from decimal import Decimal
from pathlib import Path

import pdfplumber
import structlog

from fatura.exceptions import FieldNotFoundError, ParserError
from fatura.models import (
    Cliente,
    Consumo,
    ContaDistribuidora,
    HistoricoConsumo,
    ItemFatura,
    normalizar_decimal_br,
)

logger = structlog.get_logger()

MIN_TEXT_LENGTH = 100

# Meses em português usados no histórico de consumo ("DEZ24", "NOV24", etc.)
_MESES_ABREV = {"JAN", "FEV", "MAR", "ABR", "MAI", "JUN",
                "JUL", "AGO", "SET", "OUT", "NOV", "DEZ"}


class CoelbaPdfParser:

    def parse(self, pdf_path: Path | str) -> ContaDistribuidora:
        pdf_path = Path(pdf_path)
        if not pdf_path.exists():
            raise ParserError(f"Arquivo PDF não encontrado: {pdf_path}")

        try:
            with pdfplumber.open(pdf_path) as pdf:
                page = pdf.pages[0]
                tables = page.extract_tables()
                text = page.extract_text() or ""
        except Exception as e:
            raise ParserError(f"Erro ao abrir PDF: {e}") from e

        if len(text) < MIN_TEXT_LENGTH and not tables:
            raise ParserError(
                f"PDF sem texto extraível ({len(text)} chars) e sem tabelas. "
                "Pode ser necessário OCR. Instale com: pip install '5g-energia-fatura[ocr]'"
            )

        logger.debug(
            "pdf_carregado",
            tabelas=len(tables),
            chars_texto=len(text),
            path=str(pdf_path),
        )

        try:
            uc = self._extrair_uc(tables, text)
        except FieldNotFoundError:
            raise
        except Exception as e:
            raise FieldNotFoundError("uc", str(e)) from e

        header = self._extrair_header(tables, text)
        total_itens = self._extrair_total_itens(tables)
        cliente = self._extrair_cliente(tables, text)
        consumo = self._extrair_consumo(tables)
        historico = self._extrair_historico(tables)
        itens = self._extrair_itens(tables)

        # "valor" = total dos itens (valor real do consumo antes da compensação solar)
        # "TOTAL A PAGAR" pode ser 0 quando os créditos solares cobrem tudo
        valor = total_itens or header.get("valor_a_pagar") or Decimal("0")

        return ContaDistribuidora(
            uc=uc,
            mes=header["mes"],
            ano=header["ano"],
            valor=valor,
            vencimento=header["vencimento"],
            numero_dias=header.get("numero_dias"),
            codigo_barras=None,  # boleto gerado online, não extraível do PDF
            nota_fiscal=None,
            cliente=cliente,
            consumo=consumo,
            historico_energia=historico,
            composicao=None,
            itens_fatura=itens,
            pdf_path=str(pdf_path),
            parsed_at=datetime.now(),
        )

    # -------------------------------------------------------------------------
    # Extração de UC
    # -------------------------------------------------------------------------

    def _extrair_uc(self, tables: list, text: str) -> str:
        """UC = CÓDIGO DO CLIENTE (tabela 0, linha 1)."""
        # Busca na tabela 0 primeiro
        t0 = tables[0] if tables else []
        for row in t0:
            flat = self._nonnone(row)
            for cell in flat:
                if "CÓDIGO DO CLIENTE" in cell:
                    parts = cell.split("\n")
                    for p in parts:
                        p = p.strip()
                        if p.isdigit() and 6 <= len(p) <= 12:
                            return p

        # Fallback: busca no texto (UC aparece antes de "chave de acesso")
        m = re.search(r"(\d{10})\s+chave de acesso", text, re.I)
        if m:
            return m.group(1)

        # Fallback: número de 10 dígitos isolado próximo de "LAPAO BA" ou similar
        m = re.search(r"BA\s+(\d{10})", text)
        if m:
            return m.group(1)

        # Busca na tabela 1 (último rodapé da fatura tem UC repetida)
        t1 = tables[1] if len(tables) > 1 else []
        for row in t1:
            flat = self._nonnone(row)
            for cell in flat:
                if "CÓDIGO DO CLIENTE" in cell:
                    parts = cell.split("\n")
                    for p in parts:
                        p = p.strip()
                        if p.isdigit() and 6 <= len(p) <= 12:
                            return p

        raise FieldNotFoundError("uc", "UC (CÓDIGO DO CLIENTE) não encontrada no PDF")

    # -------------------------------------------------------------------------
    # Extração de cabeçalho (mês, ano, vencimento, valor a pagar, n° dias)
    # -------------------------------------------------------------------------

    def _extrair_header(self, tables: list, text: str) -> dict:
        result: dict = {}
        t1 = tables[1] if len(tables) > 1 else []

        for row in t1:
            flat = self._nonnone(row)
            if not flat:
                continue
            joined = " ".join(flat)

            # Linha: REF:MÊS/ANO / TOTAL A PAGAR / VENCIMENTO
            if "REF:MÊS/ANO" in joined or "MÊS/ANO" in joined:
                for cell in flat:
                    # Mês/Ano: "REF:MÊS/ANO\n12/2024"
                    m = re.search(r"(\d{1,2})/(\d{4})", cell)
                    if m and "MÊS" in cell or "REF" in cell:
                        result["mes"] = int(m.group(1))
                        result["ano"] = int(m.group(2))
                    # TOTAL A PAGAR: "TOTAL A PAGAR R$\n0,00"
                    m = re.search(r"TOTAL A PAGAR.*?([\d.,]+)$", cell, re.S)
                    if m:
                        result["valor_a_pagar"] = normalizar_decimal_br(m.group(1)) or Decimal("0")
                    # Vencimento: "VENCIMENTO\n20/12/2024"
                    m = re.search(r"VENCIMENTO\s*\n?\s*(\d{2}/\d{2}/\d{4})", cell)
                    if m:
                        result["vencimento"] = _parse_date(m.group(1))

                # Fallback para mês/ano: se alguma célula tem "MM/AAAA" sem rótulo
                if "mes" not in result:
                    for cell in flat:
                        m = re.search(r"\b(\d{1,2})/(\d{4})\b", cell)
                        if m:
                            result["mes"] = int(m.group(1))
                            result["ano"] = int(m.group(2))
                            break

            # Linha: N° DE DIAS
            if "N° DE DIAS" in joined or "N° DE DIAS" in joined:
                m = re.search(r"N[°º]\s*DE\s*DIAS\s+(\d+)", joined)
                if m:
                    result["numero_dias"] = int(m.group(1))

        # Fallback: extrair vencimento do texto
        if "vencimento" not in result:
            m = re.search(r"VENCIMENTO\s*\n?\s*(\d{2}/\d{2}/\d{4})", text)
            if m:
                result["vencimento"] = _parse_date(m.group(1))

        # Fallback: extrair mês/ano do texto
        if "mes" not in result:
            m = re.search(r"REF:MÊS/ANO.*?(\d{1,2})/(\d{4})", text, re.S)
            if m:
                result["mes"] = int(m.group(1))
                result["ano"] = int(m.group(2))

        if "vencimento" not in result:
            raise FieldNotFoundError("vencimento", "Data de vencimento não encontrada")
        if "mes" not in result:
            raise FieldNotFoundError("mes", "Mês de referência não encontrado")

        return result

    # -------------------------------------------------------------------------
    # Extração do total de itens
    # -------------------------------------------------------------------------

    def _extrair_total_itens(self, tables: list) -> Decimal | None:
        """Linha 'TOTAL' na tabela 1 contém a soma de todos os itens."""
        t1 = tables[1] if len(tables) > 1 else []
        for row in t1:
            flat = self._nonnone(row)
            if not flat:
                continue
            if str(flat[0]).strip().upper() == "TOTAL" and len(flat) >= 2:
                # O valor é o segundo elemento não-None
                return normalizar_decimal_br(flat[-1])
        return None

    # -------------------------------------------------------------------------
    # Extração de cliente
    # -------------------------------------------------------------------------

    def _extrair_cliente(self, tables: list, text: str) -> Cliente:
        nome = ""
        cpf = None
        cnpj = None
        classificacao = None
        endereco = None

        # Nome e CPF da tabela 2 (boleto) — linha PAGADOR
        t2 = tables[2] if len(tables) > 2 else []
        for row in t2:
            flat = self._nonnone(row)
            for cell in flat:
                if "PAGADOR" in cell:
                    linhas = cell.split("\n")
                    if len(linhas) >= 2:
                        # Segunda linha: "NOME CPF_MASCARADO"
                        partes = linhas[1].rsplit(" ", 1)
                        nome = partes[0].strip() if partes else linhas[1].strip()
                        if len(partes) > 1 and re.match(r"\d{3}\.\d", partes[-1]):
                            cpf = partes[-1].strip()
                    if len(linhas) >= 3:
                        endereco = linhas[2].strip()
                    break

        # Classificação da tabela 1
        t1 = tables[1] if len(tables) > 1 else []
        for row in t1:
            flat = self._nonnone(row)
            for cell in flat:
                if "CLASSIFICAÇÃO:" in cell:
                    m = re.search(r"CLASSIFICAÇÃO:\s*(.+?)(?:\s*-\s*|$)", cell)
                    if m:
                        classificacao = m.group(1).strip()
                    break

        # Fallback nome do texto
        if not nome:
            m = re.search(r"NOME DO CLIENTE:\s*\n?\s*([A-ZÁÉÍÓÚÂÊÎÔÛÃÕÀÇ ]{5,})", text)
            if m:
                nome = m.group(1).strip()

        return Cliente(
            nome=nome,
            cpf=cpf,
            cnpj=cnpj,
            classificacao=classificacao,
            endereco=endereco,
        )

    # -------------------------------------------------------------------------
    # Extração de dados do medidor
    # -------------------------------------------------------------------------

    def _extrair_consumo(self, tables: list) -> Consumo | None:
        """Linha com número do medidor (10 dígitos) na tabela 1."""
        t1 = tables[1] if len(tables) > 1 else []
        for row in t1:
            flat = self._nonnone(row)
            if not flat:
                continue
            # Linha do medidor: [numero_medidor, grandeza, postos, leit_ant, leit_atu, constante, consumo]
            if re.match(r"^\d{7,12}$", str(flat[0]).strip()):
                cells = flat
                constante = cells[5] if len(cells) > 5 else None
                return Consumo(
                    medidor=cells[0].strip(),
                    leitura_anterior=cells[3].strip() if len(cells) > 3 else None,
                    leitura_atual=cells[4].strip() if len(cells) > 4 else None,
                    constante=constante.strip() if constante else None,
                )
        return None

    # -------------------------------------------------------------------------
    # Extração do histórico de consumo
    # -------------------------------------------------------------------------

    def _extrair_historico(self, tables: list) -> list[HistoricoConsumo]:
        """Célula com 'CONSUMO FATURADO N°DIAS FAT' na tabela 1."""
        t1 = tables[1] if len(tables) > 1 else []
        for row in t1:
            flat = self._nonnone(row)
            for cell in flat:
                if "CONSUMO FATURADO" in cell:
                    return self._parsear_historico_cell(cell)
        return []

    def _parsear_historico_cell(self, cell: str) -> list[HistoricoConsumo]:
        historico = []
        # Padrão: "DEZ24 433 31" ou "DEZ24 433" (sem dias)
        for match in re.finditer(
            r"\b([A-Z]{3})(\d{2})\s+(\d+(?:[.,]\d+)?)\b",
            cell
        ):
            mes_str = match.group(1)
            if mes_str not in _MESES_ABREV:
                continue
            ano_str = match.group(2)
            kwh_str = match.group(3)
            kwh = normalizar_decimal_br(kwh_str)
            if kwh is None or kwh <= 0:
                continue
            historico.append(HistoricoConsumo(
                periodo=f"{mes_str}/{ano_str}",
                kwh=kwh,
            ))
        return historico

    # -------------------------------------------------------------------------
    # Extração dos itens da fatura
    # -------------------------------------------------------------------------

    def _extrair_itens(self, tables: list) -> list[ItemFatura]:
        """Linha de dados logo após o cabeçalho 'ITENS DA FATURA' na tabela 1."""
        t1 = tables[1] if len(tables) > 1 else []
        header_idx = None
        for i, row in enumerate(t1):
            flat = self._nonnone(row)
            if flat and "ITENS DA FATURA" in str(flat[0]):
                header_idx = i
                break

        if header_idx is None:
            return []

        # Próxima linha de dados (ignorar linhas vazias ou com apenas uma célula)
        data_row = None
        for i in range(header_idx + 1, len(t1)):
            flat = self._nonnone(t1[i])
            if len(flat) >= 4 and re.search(r"Consumo|Ilum|Acrés|IPCA|Custo|UFER", flat[0], re.I):
                data_row = flat
                break

        if data_row is None:
            return []

        return self._parsear_itens_row(data_row)

    def _parsear_itens_row(self, cells: list) -> list[ItemFatura]:
        """
        Linha de dados dos itens tem células com valores separados por '\\n'.
        Estrutura (14 colunas não-None):
          [0] Descrições     [1] Unidades     [2] Qtd
          [3] Preço unit c/trib  [4] Valor R$  [5] PIS/COFINS
          [6] Base ICMS      [7] Alíq ICMS%   [8] ICMS R$
          [9] Tarifa unit    [10] Tributo      [11] Base tributo
          [12] Alíq trib%    [13] Val tributo
        """
        def split_cell(idx: int) -> list[str]:
            if idx >= len(cells):
                return []
            return [v.strip() for v in str(cells[idx]).split("\n") if v.strip()]

        descricoes = split_cell(0)
        quantidades = split_cell(2)
        valores = split_cell(4)       # VALOR (R$) = qtd × tarifa
        base_icms_list = split_cell(6)
        aliq_icms_list = split_cell(7)
        icms_list = split_cell(8)
        tarifas = split_cell(9)       # Tarifa limpa (sem tributos)

        itens = []
        for i, descricao in enumerate(descricoes):
            def get(lst: list, idx: int, default: str = "") -> str:
                return lst[idx] if idx < len(lst) else default

            item = ItemFatura(
                descricao=descricao,
                quantidade=normalizar_decimal_br(get(quantidades, i)) if get(quantidades, i) else None,
                tarifa=normalizar_decimal_br(get(tarifas, i)) if get(tarifas, i) else None,
                valor=normalizar_decimal_br(get(valores, i)) if get(valores, i) else None,
                base_icms=normalizar_decimal_br(get(base_icms_list, i)) if get(base_icms_list, i) else None,
                aliq_icms=get(aliq_icms_list, i) or None,
                icms=normalizar_decimal_br(get(icms_list, i)) if get(icms_list, i) else None,
            )
            itens.append(item)
        return itens

    # -------------------------------------------------------------------------
    # Utilitários
    # -------------------------------------------------------------------------

    @staticmethod
    def _nonnone(row: list) -> list:
        """Remove None e células em branco de uma linha da tabela."""
        return [c for c in row if c is not None and str(c).strip()]


def _parse_date(s: str) -> date:
    """Converte 'DD/MM/AAAA' para date."""
    parts = s.strip().split("/")
    if len(parts) == 3:
        dia, mes, ano = int(parts[0]), int(parts[1]), int(parts[2])
        return date(ano, mes, dia)
    raise ValueError(f"Formato de data inválido: {s!r}")


def parse_fatura_coelba_ocr(pdf_path: Path) -> ContaDistribuidora:
    """Ponto de extensão OCR. Implementar quando necessário."""
    raise NotImplementedError(
        "OCR não implementado. Instale pytesseract e implemente esta função."
    )
