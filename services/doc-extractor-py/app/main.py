from __future__ import annotations

import base64
import os
import tempfile
from pathlib import Path
from typing import Any, Literal

import uvicorn
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field

from fatura.config import ParserConfig
from fatura.models import conta_para_ocr_payload
from fatura.parser_pdf import CoelbaPdfParser


class PDFPayload(BaseModel):
    mode: Literal["path", "base64"]
    path: str | None = None
    base64: str | None = None
    file_name: str | None = None


class ExtractRequest(BaseModel):
    schema_version: str = "1.0.0"
    job_id: str
    uc: str
    documento: str
    numero_fatura: str
    mes_referencia: str | None = None
    pdf: PDFPayload
    requested_fields: list[str] = Field(default_factory=list)
    api_snapshot: dict[str, Any] | None = None


class ExtractResponse(BaseModel):
    schema_version: str = "1.0.0"
    job_id: str
    status: Literal["ok", "partial", "error"]
    fields: dict[str, Any]
    source_map: dict[str, str]
    confidence_map: dict[str, float]
    warnings: list[str] = Field(default_factory=list)
    artifacts: dict[str, str | None] = Field(default_factory=dict)


def _build_parser() -> CoelbaPdfParser:
    config = ParserConfig(
        engine="pymupdf",
        enable_mistral_fallback=bool(os.getenv("MISTRAL_API_KEY")),
        validate_new_pdfs_with_mistral=False,
        mistral_api_key=os.getenv("MISTRAL_API_KEY", ""),
        mistral_model=os.getenv("MISTRAL_MODEL", "mistral-ocr-latest"),
    )
    return CoelbaPdfParser(config)


def _write_base64_pdf(payload: PDFPayload) -> Path:
    suffix = payload.file_name or "document.pdf"
    tmp = tempfile.NamedTemporaryFile(prefix="extractor_", suffix="_" + suffix, delete=False)
    with tmp:
        tmp.write(base64.b64decode(payload.base64 or ""))
    return Path(tmp.name)


def _resolve_pdf_path(payload: PDFPayload) -> tuple[Path, bool]:
    if payload.mode == "path":
        if not payload.path:
            raise HTTPException(status_code=400, detail="pdf.path é obrigatório quando mode=path")
        return Path(payload.path), False
    if not payload.base64:
        raise HTTPException(status_code=400, detail="pdf.base64 é obrigatório quando mode=base64")
    return _write_base64_pdf(payload), True


def _set_source_and_confidence(value: Any, path: str, source_map: dict[str, str], confidence_map: dict[str, float]) -> None:
    if isinstance(value, dict):
        for key, inner in value.items():
            child = f"{path}.{key}" if path else key
            _set_source_and_confidence(inner, child, source_map, confidence_map)
        return
    if isinstance(value, list):
        if value:
            source_map[path] = "pymupdf"
            confidence_map[path] = 0.9
        else:
            source_map[path] = "unknown"
            confidence_map[path] = 0.0
        return
    if value in (None, "", {}):
        source_map[path] = "unknown"
        confidence_map[path] = 0.0
        return
    source_map[path] = "pymupdf"
    confidence_map[path] = 0.9


def create_app() -> FastAPI:
    app = FastAPI(
        title="doc-extractor-py",
        version="0.1.0",
        summary="Extrator documental de faturas Neoenergia",
    )
    parser = _build_parser()

    @app.get("/healthz")
    async def healthz() -> dict[str, str]:
        return {"status": "ok", "service": "doc-extractor-py"}

    @app.post("/v1/extract", response_model=ExtractResponse)
    async def extract(request: ExtractRequest) -> ExtractResponse:
        pdf_path, should_cleanup = _resolve_pdf_path(request.pdf)
        try:
            conta = parser.parse(pdf_path)
            payload = conta_para_ocr_payload(conta)
            source_map: dict[str, str] = {}
            confidence_map: dict[str, float] = {}
            _set_source_and_confidence(payload, "", source_map, confidence_map)
            status: Literal["ok", "partial", "error"] = "ok"
            warnings: list[str] = []
            if any(source == "unknown" for source in source_map.values()):
                status = "partial"
                warnings.append("Alguns campos permaneceram vazios após o parse local.")
            return ExtractResponse(
                job_id=request.job_id,
                status=status,
                fields=payload,
                source_map=source_map,
                confidence_map=confidence_map,
                warnings=warnings,
                artifacts={"pdf_path": str(pdf_path)},
            )
        except Exception as exc:
            return ExtractResponse(
                job_id=request.job_id,
                status="error",
                fields={},
                source_map={},
                confidence_map={},
                warnings=[str(exc)],
                artifacts={"pdf_path": str(pdf_path)},
            )
        finally:
            if should_cleanup:
                pdf_path.unlink(missing_ok=True)

    return app


def main() -> None:
    uvicorn.run(
        create_app(),
        host=os.getenv("DOC_EXTRACTOR_HOST", "127.0.0.1"),
        port=int(os.getenv("DOC_EXTRACTOR_PORT", "8090")),
    )


if __name__ == "__main__":
    main()
