from __future__ import annotations

import json
from pathlib import Path
from typing import Any, Protocol

import structlog
from pydantic import BaseModel, Field

from fatura.config import ParserConfig
from fatura.exceptions import ParserError

logger = structlog.get_logger()


class OCRClientePayload(BaseModel):
    codigo: str | None = Field(default=None, description="Código do cliente / UC.")
    cpf: str | None = None
    cnpj: str | None = None
    nome: str | None = None
    classificacao: str | None = None
    tensao_nominal: str | None = None
    limites_tensao: str | None = None
    endereco: str | None = None


class OCRConsumoPayload(BaseModel):
    medidor: str | None = None
    constante: str | None = None
    leitura_anterior: str | None = None
    leitura_atual: str | None = None
    leitura_anterior_data: str | None = None
    leitura_data: str | None = None
    leitura_proxima_data: str | None = None


class OCRNotaFiscalPayload(BaseModel):
    numero_serie: str | None = None
    apresentacao_data: str | None = None


class OCRHistoricoPayload(BaseModel):
    periodo: str | None = None
    kwh: str | None = None


class OCRComposicaoPayload(BaseModel):
    energia: str | None = None
    encargos: str | None = None
    distribuicao: str | None = None
    tributos: str | None = None
    transmissao: str | None = None
    perdas: str | None = None


class OCRItemPayload(BaseModel):
    codigo: str | None = None
    descricao: str | None = None
    quantidade: str | None = None
    quantidade_residual: str | None = None
    quantidade_faturada: str | None = None
    tarifa: str | None = None
    valor: str | None = None
    base_icms: str | None = None
    aliq_icms: str | None = None
    icms: str | None = None
    valor_total: str | None = None


class OCRExtractionPayload(BaseModel):
    mes: int | None = None
    ano: int | None = None
    valor: str | None = None
    normalizado_valor: float | None = None
    vencimento: str | None = None
    leitura_anterior_data: str | None = None
    leitura_data: str | None = None
    leitura_proxima_data: str | None = None
    emissao_data: str | None = None
    controle_n: str | None = None
    numero_dias: int | None = None
    codigo_barras: str | None = None
    aviso: str | None = None
    nota_fiscal: OCRNotaFiscalPayload | None = None
    cliente: OCRClientePayload | None = None
    consumo: OCRConsumoPayload | None = None
    energia: dict[str, list[OCRHistoricoPayload]] | None = None
    composicao_fornecimento: OCRComposicaoPayload | None = None
    informacoes_gerais: str | None = None
    itens_fatura: list[OCRItemPayload] = []


class OCRClientProtocol(Protocol):
    def extract_with_mistral(self, pdf_path: Path) -> dict[str, Any]: ...


class MistralOCRClient:
    def __init__(self, config: ParserConfig) -> None:
        self._config = config

    @property
    def enabled(self) -> bool:
        return bool(self._config.mistral_api_key)

    def extract_with_mistral(self, pdf_path: Path) -> dict[str, Any]:
        if not self.enabled:
            raise ParserError(
                "Fallback Mistral OCR solicitado, mas MISTRAL_API_KEY não está configurada."
            )

        try:
            from mistralai.client import Mistral
            from mistralai.client.models import FileChunk
            from mistralai.extra import response_format_from_pydantic_model
        except ImportError as exc:  # pragma: no cover - depende do ambiente
            raise ParserError(
                "Biblioteca 'mistralai' não instalada. Instale com: pip install '5g-energia-fatura[ocr]'"
            ) from exc

        client = Mistral(api_key=self._config.mistral_api_key)
        logger.info("mistral_ocr_upload_iniciado", path=str(pdf_path))
        with pdf_path.open("rb") as handle:
            uploaded = client.files.upload(
                file={
                    "file_name": pdf_path.name,
                    "content": handle,
                    "content_type": "application/pdf",
                },
                purpose="ocr",
            )

        logger.info("mistral_ocr_process_iniciado", path=str(pdf_path), file_id=uploaded.id)
        response = client.ocr.process(
            model=self._config.mistral_model,
            document=FileChunk(file_id=uploaded.id),
            document_annotation_format=response_format_from_pydantic_model(OCRExtractionPayload),
            document_annotation_prompt=(
                "Extraia os campos da fatura de energia em JSON. "
                "Preserve datas no formato DD/MM/AAAA, valores monetários com vírgula, "
                "e preencha apenas o que estiver claramente presente no documento. "
                "Não invente código de barras ou composição do fornecimento quando ausentes."
            ),
            timeout_ms=self._config.mistral_timeout_ms,
        )

        annotation = response.document_annotation
        if annotation is None:
            raise ParserError("Mistral OCR respondeu sem document_annotation.")

        if hasattr(annotation, "model_dump"):
            payload = annotation.model_dump()
        elif isinstance(annotation, str):
            payload = json.loads(annotation)
        else:
            payload = dict(annotation)

        validated = OCRExtractionPayload.model_validate(payload)
        logger.info(
            "mistral_ocr_process_concluido",
            path=str(pdf_path),
            model=response.model,
            has_codigo_barras=bool(validated.codigo_barras),
            has_composicao=validated.composicao_fornecimento is not None,
        )
        return validated.model_dump(exclude_none=True)
