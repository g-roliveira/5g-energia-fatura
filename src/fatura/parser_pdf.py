"""Parser de faturas PDF da Coelba/Neoenergia usando PyMuPDF.

Estratégia:
1. PyMuPDF (`fitz`) é a engine primária e obrigatória.
2. Mistral OCR entra apenas como fallback/validação quando o parse local estiver
   incompleto ou quando a validação opcional estiver habilitada para um PDF ainda
   não visto neste processo.
"""

from __future__ import annotations

import hashlib
import re
import time
from dataclasses import dataclass
from datetime import date, datetime
from decimal import Decimal
from pathlib import Path
from typing import Any

import fitz
import structlog

from fatura.config import ParserConfig
from fatura.exceptions import FieldNotFoundError, ParserError
from fatura.mistral_ocr import MistralOCRClient, OCRClientProtocol
from fatura.models import (
    Cliente,
    ComposicaoFornecimento,
    Consumo,
    ContaDistribuidora,
    HistoricoConsumo,
    ItemFatura,
    NotaFiscal,
    normalizar_decimal_br,
)

logger = structlog.get_logger()

MIN_TEXT_LENGTH = 100
_MESES_ABREV = {"JAN", "FEV", "MAR", "ABR", "MAI", "JUN", "JUL", "AGO", "SET", "OUT", "NOV", "DEZ"}


@dataclass
class PageSnapshot:
    text: str
    words: list[dict[str, Any]]
    tables: list[list[list[Any]]]
    blocks: list[dict[str, Any]]


class CoelbaPdfParser:
    def __init__(
        self,
        config: ParserConfig | None = None,
        ocr_client: OCRClientProtocol | None = None,
    ) -> None:
        self._config = config or ParserConfig()
        self._ocr_client = ocr_client
        if self._ocr_client is None and self._config.mistral_api_key:
            self._ocr_client = MistralOCRClient(self._config)
        self._validated_fingerprints: set[str] = set()

    def parse(self, pdf_path: Path | str) -> ContaDistribuidora:
        pdf_path = Path(pdf_path)
        if not pdf_path.exists():
            raise ParserError(f"Arquivo PDF não encontrado: {pdf_path}")

        started_at = time.perf_counter()
        snapshots = self._carregar_snapshot(pdf_path)
        first_page = snapshots[0]

        # Se PDF escaneado e Mistral disponível, vai direto para OCR
        if len(first_page.text) < MIN_TEXT_LENGTH and not first_page.tables:
            if not self._config.enable_mistral_fallback:
                raise ParserError(
                    f"PDF sem texto extraível ({len(first_page.text)} chars) e sem tabelas. "
                    "Configure Mistral OCR para fallback semântico."
                )
            logger.info("pdf_parse_sem_texto_ocr_fallback", path=str(pdf_path))
            return self._parse_with_ocr_only(pdf_path)

        conta = self._parse_with_pymupdf(pdf_path, snapshots)
        missing_fields = self._missing_fields(conta, first_page)
        logger.info(
            "pdf_parse_pymupdf_concluido",
            path=str(pdf_path),
            tabelas_primeira_pagina=len(first_page.tables),
            missing_fields=missing_fields,
            duration_ms=round((time.perf_counter() - started_at) * 1000, 2),
        )

        fingerprint = self._fingerprint(pdf_path)
        should_validate = (
            self._config.validate_new_pdfs_with_mistral
            and fingerprint not in self._validated_fingerprints
        )
        should_fallback = bool(missing_fields) and self._config.enable_mistral_fallback

        if should_validate or should_fallback:
            conta = self._maybe_enrich_with_ocr(
                pdf_path=pdf_path,
                conta=conta,
                missing_fields=missing_fields,
                validate_only=should_validate,
            )
            if should_validate:
                self._validated_fingerprints.add(fingerprint)

        return conta

    def _carregar_snapshot(self, pdf_path: Path) -> list[PageSnapshot]:
        try:
            doc = fitz.open(pdf_path)
        except Exception as exc:
            raise ParserError(f"Erro ao abrir PDF com PyMuPDF: {exc}") from exc

        snapshots: list[PageSnapshot] = []
        try:
            for page in doc:
                tables = []
                try:
                    tables = [table.extract() for table in page.find_tables().tables]
                except Exception:
                    tables = []

                words = [
                    {
                        "x0": word[0],
                        "y0": word[1],
                        "x1": word[2],
                        "y1": word[3],
                        "text": word[4],
                        "block_no": word[5],
                        "line_no": word[6],
                        "word_no": word[7],
                    }
                    for word in page.get_text("words")
                ]
                page_dict = page.get_text("dict")
                snapshots.append(
                    PageSnapshot(
                        text=page.get_text("text") or "",
                        words=words,
                        tables=tables,
                        blocks=page_dict.get("blocks", []),
                    )
                )
        finally:
            doc.close()
        return snapshots

    def _parse_with_pymupdf(
        self,
        pdf_path: Path,
        snapshots: list[PageSnapshot],
    ) -> ContaDistribuidora:
        first_page = snapshots[0]
        full_text = "\n".join(snapshot.text for snapshot in snapshots if snapshot.text)

        uc = self._extrair_uc(first_page.tables, first_page.text)
        header = self._extrair_header(first_page.tables, full_text)
        total_itens = self._extrair_total_itens(first_page.tables)
        cliente = self._extrair_cliente(first_page.tables, full_text, uc)
        consumo = self._extrair_consumo(first_page.tables)
        historico = self._extrair_historico(first_page.tables)
        itens = self._extrair_itens(first_page.tables)
        nota_fiscal = self._extrair_nota_fiscal(first_page.words, full_text)
        emissao_data = self._extrair_emissao_data(full_text, first_page.words)
        aviso = self._extrair_aviso(full_text)
        informacoes_gerais = self._extrair_informacoes_gerais(full_text)
        codigo_barras = self._extrair_codigo_barras(full_text)
        composicao = self._extrair_composicao(full_text)

        valor = total_itens or header.get("valor_a_pagar") or Decimal("0")
        return ContaDistribuidora(
            uc=uc,
            mes=header["mes"],
            ano=header["ano"],
            valor=valor,
            vencimento=header["vencimento"],
            numero_dias=header.get("numero_dias"),
            codigo_barras=codigo_barras,
            nota_fiscal=nota_fiscal,
            emissao_data=emissao_data,
            controle_n=self._extrair_controle_n(full_text),
            aviso=aviso,
            informacoes_gerais=informacoes_gerais,
            cliente=cliente,
            consumo=consumo,
            historico_energia=historico,
            composicao=composicao,
            itens_fatura=itens,
            pdf_path=str(pdf_path),
            parsed_at=datetime.now(),
        )


    def _parse_with_ocr_only(self, pdf_path: Path) -> ContaDistribuidora:
        """Usa Mistral OCR diretamente em PDFs escaneados (sem texto extraível)."""
        if self._ocr_client is None:
            raise ParserError(
                "OCR fallback requisitado mas cliente OCR não configurado. "
                "Verifique MISTRAL_API_KEY."
            )

        try:
            ocr_payload = self._ocr_client.extract_with_mistral(pdf_path)
        except Exception as exc:
            raise ParserError(f"OCR fallback falhou: {exc}") from exc

        cliente_payload = ocr_payload.get("cliente") or {}
        uc = cliente_payload.get("codigo") or ""
        vencimento_raw = ocr_payload.get("vencimento")
        try:
            vencimento = _parse_date(vencimento_raw) if vencimento_raw else None
        except Exception:
            vencimento = None

        # Usa normalizado_valor (float) como valor principal, depois tenta string
        valor = ocr_payload.get("normalizado_valor") or ocr_payload.get("valor") or "0"

        conta = ContaDistribuidora(
            uc=uc,
            mes=ocr_payload.get("mes") or 0,
            ano=ocr_payload.get("ano") or 0,
            valor=valor,
            vencimento=vencimento or date.today(),
            emissao_data=ocr_payload.get("emissao_data"),
            codigo_barras=ocr_payload.get("codigo_barras"),
            cliente=cliente_payload,
            consumo=ocr_payload.get("consumo"),
            informacoes_gerais=ocr_payload.get("informacoes_gerais"),
            itens_fatura=ocr_payload.get("itens_fatura", []),
            pdf_path=str(pdf_path),
            parsed_at=datetime.now(),
        )
        logger.info("pdf_parse_ocr_only_concluido", path=str(pdf_path))
        return conta

    def _maybe_enrich_with_ocr(
        self,
        pdf_path: Path,
        conta: ContaDistribuidora,
        missing_fields: list[str],
        validate_only: bool,
    ) -> ContaDistribuidora:
        if not (validate_only or self._config.enable_mistral_fallback):
            return conta
        if self._ocr_client is None:
            logger.warning(
                "pdf_parse_ocr_nao_configurado",
                path=str(pdf_path),
                missing_fields=missing_fields,
                validate_only=validate_only,
            )
            return conta

        ocr_started_at = time.perf_counter()
        logger.info(
            "pdf_parse_ocr_iniciado",
            path=str(pdf_path),
            missing_fields=missing_fields,
            validate_only=validate_only,
        )
        try:
            ocr_payload = self._ocr_client.extract_with_mistral(pdf_path)
        except Exception as exc:
            logger.warning("pdf_parse_ocr_falhou", path=str(pdf_path), error=str(exc))
            return conta

        conta_enriquecida, preenchidos = self._reconcile_with_ocr(conta, ocr_payload)
        logger.info(
            "pdf_parse_ocr_concluido",
            path=str(pdf_path),
            preenchidos=preenchidos,
            duration_ms=round((time.perf_counter() - ocr_started_at) * 1000, 2),
        )
        return conta_enriquecida

    def _missing_fields(self, conta: ContaDistribuidora, first_page: PageSnapshot) -> list[str]:
        missing: list[str] = []
        checks = {
            "uc": conta.uc,
            "mes": conta.mes,
            "ano": conta.ano,
            "valor": conta.valor,
            "vencimento": conta.vencimento,
            "emissao_data": conta.emissao_data,
            "cliente.nome": conta.cliente.nome,
            "cliente.codigo": conta.cliente.codigo,
            "consumo.medidor": conta.consumo.medidor if conta.consumo else None,
            "historico_energia": conta.historico_energia,
            "itens_fatura": conta.itens_fatura,
            "codigo_barras": conta.codigo_barras,
            "composicao_fornecimento": conta.composicao,
        }
        for field_name, value in checks.items():
            if value in (None, "", [], {}):
                missing.append(field_name)
        if len(first_page.tables) < 2:
            missing.append("tabelas_insuficientes")
        return missing

    def _reconcile_with_ocr(
        self,
        conta: ContaDistribuidora,
        payload: dict[str, Any],
    ) -> tuple[ContaDistribuidora, list[str]]:
        filled: list[str] = []

        def fill_if_missing(path: str, has_value: bool, setter) -> None:
            if has_value:
                return
            setter()
            filled.append(path)

        cliente_payload = payload.get("cliente") or {}
        consumo_payload = payload.get("consumo") or {}
        nota_fiscal_payload = payload.get("nota_fiscal") or {}
        composicao_payload = payload.get("composicao_fornecimento") or {}
        energia_payload = payload.get("energia") or {}

        fill_if_missing("codigo_barras", bool(conta.codigo_barras), lambda: setattr(conta, "codigo_barras", payload.get("codigo_barras")))
        fill_if_missing("emissao_data", bool(conta.emissao_data), lambda: setattr(conta, "emissao_data", payload.get("emissao_data")))
        fill_if_missing("controle_n", bool(conta.controle_n), lambda: setattr(conta, "controle_n", payload.get("controle_n")))
        fill_if_missing("aviso", bool(conta.aviso), lambda: setattr(conta, "aviso", payload.get("aviso")))
        fill_if_missing("informacoes_gerais", bool(conta.informacoes_gerais), lambda: setattr(conta, "informacoes_gerais", payload.get("informacoes_gerais")))

        if not conta.nota_fiscal and nota_fiscal_payload:
            conta.nota_fiscal = NotaFiscal(
                numero_serie=nota_fiscal_payload.get("numero_serie"),
                apresentacao_data=nota_fiscal_payload.get("apresentacao_data"),
            )
            filled.append("nota_fiscal")
        elif conta.nota_fiscal:
            if not conta.nota_fiscal.numero_serie and nota_fiscal_payload.get("numero_serie"):
                conta.nota_fiscal.numero_serie = nota_fiscal_payload["numero_serie"]
                filled.append("nota_fiscal.numero_serie")
            if not conta.nota_fiscal.apresentacao_data and nota_fiscal_payload.get("apresentacao_data"):
                conta.nota_fiscal.apresentacao_data = nota_fiscal_payload["apresentacao_data"]
                filled.append("nota_fiscal.apresentacao_data")

        if cliente_payload:
            if not conta.cliente.codigo and cliente_payload.get("codigo"):
                conta.cliente.codigo = cliente_payload["codigo"]
                filled.append("cliente.codigo")
            if not conta.cliente.nome and cliente_payload.get("nome"):
                conta.cliente.nome = cliente_payload["nome"]
                filled.append("cliente.nome")
            if not conta.cliente.cpf and cliente_payload.get("cpf"):
                conta.cliente.cpf = cliente_payload["cpf"]
                filled.append("cliente.cpf")
            if not conta.cliente.cnpj and cliente_payload.get("cnpj"):
                conta.cliente.cnpj = cliente_payload["cnpj"]
                filled.append("cliente.cnpj")
            if not conta.cliente.classificacao and cliente_payload.get("classificacao"):
                conta.cliente.classificacao = cliente_payload["classificacao"]
                filled.append("cliente.classificacao")
            if not conta.cliente.tensao_nominal and cliente_payload.get("tensao_nominal"):
                conta.cliente.tensao_nominal = cliente_payload["tensao_nominal"]
                filled.append("cliente.tensao_nominal")
            if not conta.cliente.limites_tensao and cliente_payload.get("limites_tensao"):
                conta.cliente.limites_tensao = cliente_payload["limites_tensao"]
                filled.append("cliente.limites_tensao")
            if not conta.cliente.endereco and cliente_payload.get("endereco"):
                conta.cliente.endereco = cliente_payload["endereco"]
                filled.append("cliente.endereco")

        if consumo_payload:
            if conta.consumo is None:
                conta.consumo = Consumo()
                filled.append("consumo")
            if conta.consumo and not conta.consumo.medidor and consumo_payload.get("medidor"):
                conta.consumo.medidor = consumo_payload["medidor"]
                filled.append("consumo.medidor")
            if conta.consumo and not conta.consumo.constante and consumo_payload.get("constante"):
                conta.consumo.constante = consumo_payload["constante"]
                filled.append("consumo.constante")
            if conta.consumo and not conta.consumo.leitura_anterior and consumo_payload.get("leitura_anterior"):
                conta.consumo.leitura_anterior = consumo_payload["leitura_anterior"]
                filled.append("consumo.leitura_anterior")
            if conta.consumo and not conta.consumo.leitura_atual and consumo_payload.get("leitura_atual"):
                conta.consumo.leitura_atual = consumo_payload["leitura_atual"]
                filled.append("consumo.leitura_atual")
            if conta.consumo and not conta.consumo.leitura_anterior_data and consumo_payload.get("leitura_anterior_data"):
                conta.consumo.leitura_anterior_data = consumo_payload["leitura_anterior_data"]
                filled.append("consumo.leitura_anterior_data")
            if conta.consumo and not conta.consumo.leitura_data and consumo_payload.get("leitura_data"):
                conta.consumo.leitura_data = consumo_payload["leitura_data"]
                filled.append("consumo.leitura_data")
            if conta.consumo and not conta.consumo.leitura_proxima_data and consumo_payload.get("leitura_proxima_data"):
                conta.consumo.leitura_proxima_data = consumo_payload["leitura_proxima_data"]
                filled.append("consumo.leitura_proxima_data")

        if not conta.composicao and composicao_payload:
            conta.composicao = ComposicaoFornecimento(
                energia=composicao_payload.get("energia"),
                encargos=composicao_payload.get("encargos"),
                distribuicao=composicao_payload.get("distribuicao"),
                tributos=composicao_payload.get("tributos"),
                transmissao=composicao_payload.get("transmissao"),
                perdas=composicao_payload.get("perdas"),
            )
            filled.append("composicao_fornecimento")

        historico_payload = energia_payload.get("historico_consumo") or []
        if not conta.historico_energia and historico_payload:
            conta.historico_energia = [
                HistoricoConsumo(periodo=item["periodo"], kwh=item["kwh"])
                for item in historico_payload
                if item.get("periodo") and item.get("kwh")
            ]
            if conta.historico_energia:
                filled.append("energia.historico_consumo")

        if not conta.itens_fatura and payload.get("itens_fatura"):
            conta.itens_fatura = [
                ItemFatura.model_validate(item)
                for item in payload["itens_fatura"]
                if item.get("descricao")
            ]
            if conta.itens_fatura:
                filled.append("itens_fatura")

        if payload.get("normalizado_valor") is not None and (conta.valor is None or conta.valor == Decimal("0")):
            conta.valor = Decimal(str(payload["normalizado_valor"]))
            filled.append("valor")
        if payload.get("mes") and not conta.mes:
            conta.mes = int(payload["mes"])
            filled.append("mes")
        if payload.get("ano") and not conta.ano:
            conta.ano = int(payload["ano"])
            filled.append("ano")
        if payload.get("numero_dias") and not conta.numero_dias:
            conta.numero_dias = int(payload["numero_dias"])
            filled.append("numero_dias")
        if payload.get("vencimento") and not conta.vencimento:
            conta.vencimento = _parse_date(payload["vencimento"])
            filled.append("vencimento")

        return conta, filled

    def _extrair_uc(self, tables: list, text: str) -> str:
        t0 = tables[0] if tables else []
        for row in t0:
            flat = self._nonnone(row)
            for cell in flat:
                if "CÓDIGO DO CLIENTE" in cell:
                    for part in cell.split("\n"):
                        cleaned = part.strip()
                        if cleaned.isdigit() and 6 <= len(cleaned) <= 12:
                            return cleaned

        match = re.search(r"CÓDIGO DO CLIENTE\s+(\d{6,12})", text)
        if match:
            return match.group(1)

        match = re.search(r"(\d{10})\s+chave de acesso", text, re.I)
        if match:
            return match.group(1)

        raise FieldNotFoundError("uc", "UC (CÓDIGO DO CLIENTE) não encontrada no PDF")

    def _extrair_header(self, tables: list, text: str) -> dict[str, Any]:
        result: dict[str, Any] = {}
        t1 = tables[1] if len(tables) > 1 else []
        for row in t1:
            flat = self._nonnone(row)
            if not flat:
                continue
            joined = " ".join(flat)
            if "REF:MÊS/ANO" in joined or "MÊS/ANO" in joined:
                for cell in flat:
                    match = re.search(r"(\d{1,2})/(\d{4})", cell)
                    if match and ("MÊS" in cell or "REF" in cell):
                        result["mes"] = int(match.group(1))
                        result["ano"] = int(match.group(2))
                    match = re.search(r"TOTAL A PAGAR.*?([\d.,]+)$", cell, re.S)
                    if match:
                        result["valor_a_pagar"] = normalizar_decimal_br(match.group(1)) or Decimal("0")
                    match = re.search(r"VENCIMENTO\s*\n?\s*(\d{2}/\d{2}/\d{4})", cell)
                    if match:
                        result["vencimento"] = _parse_date(match.group(1))

            if "N° DE DIAS" in joined or "Nº DE DIAS" in joined:
                match = re.search(r"N[°º]\s*DE\s*DIAS\s+(\d+)", joined)
                if match:
                    result["numero_dias"] = int(match.group(1))
                for cell in flat:
                    match = re.search(r"LEITURA ANTERIOR\s+(\d{2}/\d{2}/\d{4})", cell)
                    if match:
                        result["leitura_anterior_data"] = match.group(1)
                    match = re.search(r"LEITURA ATUAL\s+(\d{2}/\d{2}/\d{4})", cell)
                    if match:
                        result["leitura_data"] = match.group(1)
                    match = re.search(r"PR[ÓO]XIMA LEITURA\s+(\d{2}/\d{2}/\d{4})", cell)
                    if match:
                        result["leitura_proxima_data"] = match.group(1)

        if "vencimento" not in result:
            match = re.search(r"VENCIMENTO\s*\n?\s*(\d{2}/\d{2}/\d{4})", text)
            if match:
                result["vencimento"] = _parse_date(match.group(1))
        if "mes" not in result:
            match = re.search(r"REF:MÊS/ANO.*?(\d{1,2})/(\d{4})", text, re.S)
            if match:
                result["mes"] = int(match.group(1))
                result["ano"] = int(match.group(2))
        if "vencimento" not in result:
            raise FieldNotFoundError("vencimento", "Data de vencimento não encontrada")
        if "mes" not in result:
            raise FieldNotFoundError("mes", "Mês de referência não encontrado")
        return result

    def _extrair_total_itens(self, tables: list) -> Decimal | None:
        for table in tables:
            for row in table:
                flat = self._nonnone(row)
                if flat and str(flat[0]).strip().upper() == "TOTAL" and len(flat) >= 2:
                    for cell in reversed(flat[1:]):
                        values = re.findall(r"\d{1,3}(?:\.\d{3})*,\d{2}|\d+\.\d{2}", cell)
                        if values:
                            return normalizar_decimal_br(values[-1])
        return None

    def _extrair_cliente(self, tables: list, text: str, uc: str) -> Cliente:
        nome = ""
        cpf = None
        endereco = None
        classificacao = None

        t2 = tables[2] if len(tables) > 2 else []
        for row in t2:
            flat = self._nonnone(row)
            for cell in flat:
                if "PAGADOR" in cell:
                    linhas = cell.split("\n")
                    if len(linhas) >= 2:
                        partes = linhas[1].rsplit(" ", 1)
                        nome = partes[0].strip() if partes else linhas[1].strip()
                        if len(partes) > 1 and re.match(r"\d{3}\.\d", partes[-1]):
                            cpf = partes[-1].strip()
                    if len(linhas) >= 3:
                        endereco = linhas[2].strip()

        t1 = tables[1] if len(tables) > 1 else []
        for row in t1:
            flat = self._nonnone(row)
            for cell in flat:
                if "CLASSIFICAÇÃO:" in cell:
                    match = re.search(r"CLASSIFICAÇÃO:\s*(.+?)(?:\s*-\s*|$)", cell)
                    if match:
                        classificacao = match.group(1).strip()

        if not nome:
            match = re.search(r"NOME DO CLIENTE:\s*\n?\s*([A-ZÁÉÍÓÚÂÊÎÔÛÃÕÀÇ ]{5,})", text)
            if match:
                nome = match.group(1).strip()
        match = re.search(
            r"ENDEREÇO:\s*\n(.+?)\n(\d{5}-\d{3}\s+[A-Z ]+\s+[A-Z]{2})",
            text,
            re.S,
        )
        if match:
            endereco = " ".join(
                part.strip()
                for part in (match.group(1) + "\n" + match.group(2)).splitlines()
            )

        return Cliente(
            codigo=uc,
            nome=nome,
            cpf=cpf,
            classificacao=classificacao,
            endereco=endereco,
        )

    def _extrair_consumo(self, tables: list) -> Consumo | None:
        leitura_anterior_data = None
        leitura_data = None
        leitura_proxima_data = None
        for table in tables:
            for row in table:
                flat = self._nonnone(row)
                joined = " ".join(flat)
                if "DATAS DE LEITURAS" in joined:
                    match = re.search(r"LEITURA ANTERIOR\s+(\d{2}/\d{2}/\d{4})", joined)
                    if match:
                        leitura_anterior_data = match.group(1)
                    match = re.search(r"LEITURA ATUAL\s+(\d{2}/\d{2}/\d{4})", joined)
                    if match:
                        leitura_data = match.group(1)
                    match = re.search(r"PR[ÓO]XIMA LEITURA\s+(\d{2}/\d{2}/\d{4})", joined)
                    if match:
                        leitura_proxima_data = match.group(1)
                if flat and re.match(r"^\d{7,12}$", str(flat[0]).strip()):
                    return Consumo(
                        medidor=flat[0].strip(),
                        leitura_anterior=flat[3].strip() if len(flat) > 3 else None,
                        leitura_atual=flat[4].strip() if len(flat) > 4 else None,
                        constante=flat[5].strip() if len(flat) > 5 else None,
                        leitura_anterior_data=leitura_anterior_data,
                        leitura_data=leitura_data,
                        leitura_proxima_data=leitura_proxima_data,
                    )
        return None

    def _extrair_historico(self, tables: list) -> list[HistoricoConsumo]:
        for table in tables:
            for row in table:
                for cell in self._nonnone(row):
                    if "CONSUMO FATURADO" in cell:
                        return self._parsear_historico_cell(cell)
        return []

    def _parsear_historico_cell(self, cell: str) -> list[HistoricoConsumo]:
        historico = []
        for match in re.finditer(r"\b([A-Z]{3})(\d{2})\s+(\d+(?:[.,]\d+)?)\b", cell):
            mes_str = match.group(1)
            if mes_str not in _MESES_ABREV:
                continue
            kwh = normalizar_decimal_br(match.group(3))
            if kwh is None or kwh <= 0:
                continue
            historico.append(HistoricoConsumo(periodo=f"{mes_str}/{match.group(2)}", kwh=kwh))
        return historico

    def _extrair_itens(self, tables: list) -> list[ItemFatura]:
        for table in tables:
            header_idx = None
            for idx, row in enumerate(table):
                flat = self._nonnone(row)
                if flat and "ITENS DA FATURA" in str(flat[0]):
                    header_idx = idx
                    break
            if header_idx is None:
                continue

            for idx in range(header_idx + 1, len(table)):
                flat = self._nonnone(table[idx])
                if len(flat) >= 4 and re.search(r"Consumo|Ilum|Acrés|IPCA|Custo|UFER", flat[0], re.I):
                    return self._parsear_itens_row(flat)
        return []

    def _parsear_itens_row(self, cells: list) -> list[ItemFatura]:
        def split_cell(index: int) -> list[str]:
            if index >= len(cells):
                return []
            return [value.strip() for value in str(cells[index]).split("\n") if value.strip()]

        descricoes = split_cell(0)
        quantidades = split_cell(2)
        valores = split_cell(4)
        base_icms_list = split_cell(6)
        aliq_icms_list = split_cell(7)
        icms_list = split_cell(8)
        tarifas = split_cell(9)

        itens: list[ItemFatura] = []
        for idx, descricao in enumerate(descricoes):
            def get(values: list[str], current_idx: int) -> str:
                return values[current_idx] if current_idx < len(values) else ""

            valor = normalizar_decimal_br(get(valores, idx)) if get(valores, idx) else None
            itens.append(
                ItemFatura(
                    descricao=descricao,
                    quantidade=normalizar_decimal_br(get(quantidades, idx)) if get(quantidades, idx) else None,
                    tarifa=normalizar_decimal_br(get(tarifas, idx)) if get(tarifas, idx) else None,
                    valor=valor,
                    base_icms=self._normalizar_decimal_coluna(get(base_icms_list, idx)),
                    aliq_icms=get(aliq_icms_list, idx) or None,
                    icms=normalizar_decimal_br(get(icms_list, idx)) if get(icms_list, idx) else None,
                    valor_total=valor,
                )
            )
        return itens

    @staticmethod
    def _normalizar_decimal_coluna(value: str) -> Decimal | None:
        if not value:
            return None
        matches = re.findall(r"\d{1,3}(?:\.\d{3})*,\d{2}|\d+\.\d{2}", value)
        if matches:
            return normalizar_decimal_br(matches[-1])
        return normalizar_decimal_br(value)

    def _extrair_emissao_data(self, text: str, words: list[dict[str, Any]]) -> str | None:
        match = re.search(r"DATA DE EMISS[ÃA]O:\s*(\d{2}/\d{2}/\d{4})", text)
        if match:
            return match.group(1)
        words_text = self._words_to_text(words)
        match = re.search(r"DATA\s+DE\s+EMISS[ÃA]O:?\s*(\d{2}/\d{2}/\d{4})", words_text)
        return match.group(1) if match else None

    def _extrair_nota_fiscal(self, words: list[dict[str, Any]], text: str) -> NotaFiscal | None:
        words_text = self._words_to_text(words)
        match = re.search(
            r"NOTA\s+FISCAL\s+N[°º]\s*([\d.]+)\s*-\s*S[ÉE]RIE\s*([A-Z0-9]+)",
            words_text,
            re.I,
        )
        if not match:
            normalized = " ".join(text.split())
            match = re.search(
                r"NOTA\s+FISCAL\s+N[°º]\s*([\d.]+)\s*-\s*S[ÉE]RIE\s*([A-Z0-9]+)",
                normalized,
                re.I,
            )
        if not match:
            return None
        return NotaFiscal(
            numero_serie=f"Nº {match.group(1)} Série {match.group(2)}",
            apresentacao_data=self._extrair_emissao_data(text, words),
        )

    def _extrair_controle_n(self, text: str) -> str | None:
        match = re.search(r"Protocolo de autorização:\s*([\d-]+.*)", text)
        if match:
            return " ".join(match.group(1).split())
        return None

    def _extrair_aviso(self, text: str) -> str | None:
        match = re.search(r"(ATENÇÃO![\s\S]+?)INFORMAÇÕES IMPORTANTES", text, re.I)
        return " ".join(match.group(1).split()) if match else None

    def _extrair_informacoes_gerais(self, text: str) -> str | None:
        blocks = re.findall(
            r"INFORMAÇÕES IMPORTANTES\s*([\s\S]+?)(?:DANFE - DOCUMENTO AUXILIAR|$)",
            text,
            re.I,
        )
        cleaned_blocks = []
        for block in blocks:
            cleaned = " ".join(block.split())
            if cleaned and cleaned not in cleaned_blocks:
                cleaned_blocks.append(cleaned)
        return "\n".join(cleaned_blocks) if cleaned_blocks else None

    def _extrair_codigo_barras(self, text: str) -> str | None:
        candidates = re.findall(r"(?:(?:\d{5}\.\d{5}\s+){3}\d{5}\.\d{6}\s+\d\s+\d{14})", text)
        if candidates:
            return re.sub(r"\s+", " ", candidates[0]).strip()
        return None

    def _extrair_composicao(self, text: str) -> ComposicaoFornecimento | None:
        patterns = {
            "energia": r"ENERGIA\s+R?\$?\s*([\d.,]+)",
            "encargos": r"ENCARGOS\s+R?\$?\s*([\d.,]+)",
            "distribuicao": r"DISTRIBUIÇÃO\s+R?\$?\s*([\d.,]+)",
            "tributos": r"TRIBUTOS\s+R?\$?\s*([\d.,]+)",
            "transmissao": r"TRANSMISS[ÃA]O\s+R?\$?\s*([\d.,]+)",
            "perdas": r"PERDAS\s+R?\$?\s*([\d.,]+)",
        }
        extracted: dict[str, Decimal] = {}
        for field_name, pattern in patterns.items():
            match = re.search(pattern, text, re.I)
            if match:
                normalized = normalizar_decimal_br(match.group(1))
                if normalized is not None:
                    extracted[field_name] = normalized
        return ComposicaoFornecimento(**extracted) if extracted else None

    @staticmethod
    def _nonnone(row: list) -> list[str]:
        return [str(cell) for cell in row if cell is not None and str(cell).strip()]

    @staticmethod
    def _words_to_text(words: list[dict[str, Any]]) -> str:
        return " ".join(word.get("text", "").strip() for word in words if word.get("text"))

    @staticmethod
    def _fingerprint(pdf_path: Path) -> str:
        return hashlib.sha256(pdf_path.read_bytes()).hexdigest()


def _parse_date(value: str) -> date:
    parts = value.strip().split("/")
    if len(parts) != 3:
        raise ValueError(f"Formato de data inválido: {value!r}")
    day, month, year = int(parts[0]), int(parts[1]), int(parts[2])
    return date(year, month, day)


def parse_fatura_coelba_ocr(pdf_path: Path) -> ContaDistribuidora:
    parser = CoelbaPdfParser(
        config=ParserConfig(enable_mistral_fallback=True, validate_new_pdfs_with_mistral=False)
    )
    return parser.parse(pdf_path)
