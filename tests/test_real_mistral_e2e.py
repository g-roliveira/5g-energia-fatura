from __future__ import annotations

import os
from pathlib import Path

import pytest

from fatura.config import ParserConfig
from fatura.mistral_ocr import MistralOCRClient
from fatura.parser_pdf import CoelbaPdfParser


def _should_run_real_mistral_e2e() -> bool:
    return os.getenv("RUN_REAL_MISTRAL_E2E") == "1"


def _require_mistral_api_key() -> str:
    api_key = os.getenv("MISTRAL_API_KEY", "").strip()
    if not api_key:
        pytest.skip("MISTRAL_API_KEY não definida para teste real de Mistral OCR")
    return api_key


pytestmark = [
    pytest.mark.real_mistral_e2e,
    pytest.mark.skipif(
        not _should_run_real_mistral_e2e(),
        reason="RUN_REAL_MISTRAL_E2E=1 não definido",
    ),
]


class CountingRealOCRClient:
    def __init__(self, config: ParserConfig) -> None:
        self._inner = MistralOCRClient(config)
        self.calls = 0
        self.last_payload: dict | None = None

    def extract_with_mistral(self, pdf_path: Path) -> dict:
        self.calls += 1
        payload = self._inner.extract_with_mistral(pdf_path)
        self.last_payload = payload
        return payload


def test_real_mistral_ocr_enrichment_on_real_pdf(fixtures_dir: Path):
    # Import guard: teste real precisa do SDK instalado.
    try:
        import mistralai  # noqa: F401
    except ImportError:
        pytest.skip("Pacote 'mistralai' não instalado no ambiente atual")

    api_key = _require_mistral_api_key()
    pdf_path = fixtures_dir / "334075735546.pdf"
    assert pdf_path.exists(), f"Fixture não encontrada: {pdf_path}"

    config = ParserConfig(
        enable_mistral_fallback=False,
        validate_new_pdfs_with_mistral=True,
        mistral_api_key=api_key,
        mistral_model=os.getenv("MISTRAL_MODEL", "mistral-ocr-latest"),
        mistral_timeout_ms=int(os.getenv("MISTRAL_TIMEOUT_MS", "120000")),
    )
    ocr_client = CountingRealOCRClient(config)
    parser = CoelbaPdfParser(config=config, ocr_client=ocr_client)

    conta = parser.parse(pdf_path)

    assert ocr_client.calls == 1, "Mistral OCR não foi acionado"
    assert ocr_client.last_payload is not None, "Mistral OCR foi acionado mas não retornou payload"
    assert isinstance(ocr_client.last_payload, dict)
    assert any(
        key in ocr_client.last_payload
        for key in ("cliente", "consumo", "itens_fatura", "mes", "ano", "valor")
    ), f"Payload Mistral inesperado: {ocr_client.last_payload}"

    # Sanidade mínima do parse completo (PyMuPDF + ciclo de validação Mistral).
    assert conta.uc
    assert conta.mes > 0
    assert conta.ano >= 2024
    assert conta.cliente.nome
    assert conta.consumo is not None

