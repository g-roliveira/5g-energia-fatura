from __future__ import annotations

import json
from dataclasses import asdict, dataclass
from pathlib import Path
from typing import Any

from fatura.neoenergia_private_api import ApiExchange


SENSITIVE_KEYS = {
    "authorization",
    "cookie",
    "set-cookie",
    "senha",
    "senha_portal",
    "recaptcha",
    "recaptchaanl",
    "token",
    "tokennese",
}

DOCUMENT_KEYS = {
    "documento",
    "documentosolicitante",
    "documentocliente",
    "usuario",
    "cpf",
    "cpf_cnpj",
}


@dataclass(slots=True)
class EndpointSummary:
    method: str
    path: str
    count: int
    statuses: list[int]
    sample_query: dict[str, Any] | None
    sample_request_body: Any
    sample_response_body: Any
    sample_request_headers: dict[str, Any] | None
    sample_response_headers: dict[str, Any] | None


def mask_document(value: str) -> str:
    digits = "".join(ch for ch in value if ch.isdigit())
    if len(digits) < 4:
        return "***"
    return f"{'*' * (len(digits) - 4)}{digits[-4:]}"


def sanitize_data(value: Any, *, key: str | None = None) -> Any:
    normalized_key = (key or "").lower()
    if normalized_key in SENSITIVE_KEYS:
        return "<redacted>"
    if normalized_key in DOCUMENT_KEYS and isinstance(value, str):
        return mask_document(value)
    if normalized_key == "filedata" and isinstance(value, str):
        return f"<base64:{len(value)} chars>"

    if isinstance(value, dict):
        return {k: sanitize_data(v, key=k) for k, v in value.items()}
    if isinstance(value, list):
        return [sanitize_data(item, key=key) for item in value]
    if isinstance(value, str):
        if normalized_key.endswith("token") or normalized_key.startswith("token"):
            return "<redacted>"
    return value


def sanitize_exchange(exchange: ApiExchange) -> dict[str, Any]:
    payload = asdict(exchange)
    return sanitize_data(payload)


def build_endpoint_summaries(exchanges: list[ApiExchange]) -> list[EndpointSummary]:
    grouped: dict[tuple[str, str], list[ApiExchange]] = {}
    for exchange in exchanges:
        grouped.setdefault((exchange.method, exchange.path), []).append(exchange)

    summaries: list[EndpointSummary] = []
    for method, path in sorted(grouped):
        samples = grouped[(method, path)]
        first = samples[0]
        summaries.append(
            EndpointSummary(
                method=method,
                path=path,
                count=len(samples),
                statuses=sorted({sample.status_code for sample in samples}),
                sample_query=sanitize_data(first.query),
                sample_request_body=sanitize_data(first.request_body),
                sample_response_body=sanitize_data(first.response_body),
                sample_request_headers=sanitize_data(first.request_headers),
                sample_response_headers=sanitize_data(first.response_headers),
            )
        )
    return summaries


def save_exchanges(out_dir: Path, exchanges: list[ApiExchange]) -> None:
    out_dir.mkdir(parents=True, exist_ok=True)
    for index, exchange in enumerate(exchanges, start=1):
        path_slug = exchange.path.strip("/").replace("/", "__") or "root"
        filename = f"{index:03d}_{exchange.method.lower()}_{path_slug}.json"
        (out_dir / filename).write_text(
            json.dumps(sanitize_exchange(exchange), ensure_ascii=False, indent=2),
            encoding="utf-8",
        )


def render_markdown(
    *,
    generated_at: str,
    cliente_nome: str,
    documento: str,
    ucs: list[dict[str, Any]],
    summaries: list[EndpointSummary],
    state_selection_observation: str,
    output_dir: Path,
) -> str:
    lines = [
        "# Neoenergia Private API",
        "",
        f"- gerado em: `{generated_at}`",
        f"- cliente: `{cliente_nome}`",
        f"- documento: `{mask_document(documento)}`",
        f"- artefatos: `{output_dir}`",
        "",
        "## Fluxo observado",
        "",
        "1. Login no frontend da Agência Virtual.",
        "2. Captura do Bearer token no `localStorage` após autenticação.",
        "3. Consumo dos endpoints privados com `Authorization: Bearer ...`.",
        "4. Obtenção de protocolo e chamadas por UC para conta, faturas, histórico e PDF.",
        "",
        "## Seleção de estado",
        "",
        state_selection_observation,
        "",
        "## Unidades consumidoras observadas",
        "",
    ]

    for uc in ucs:
        endereco = uc.get("local", {}).get("endereco")
        municipio = uc.get("local", {}).get("municipio")
        status = uc.get("status")
        codigo = uc.get("uc")
        lines.append(
            f"- `{codigo}` | status `{status}` | endereço `{endereco}` | município `{municipio}`"
        )

    lines.extend(["", "## Endpoints documentados", ""])
    for summary in summaries:
        lines.append(f"### `{summary.method} {summary.path}`")
        lines.append(f"- chamadas observadas: `{summary.count}`")
        lines.append(f"- status HTTP: `{summary.statuses}`")
        if summary.sample_query:
            lines.append("```json")
            lines.append(json.dumps(summary.sample_query, ensure_ascii=False, indent=2))
            lines.append("```")
        if summary.sample_request_body is not None:
            lines.append("request body:")
            lines.append("```json")
            lines.append(json.dumps(summary.sample_request_body, ensure_ascii=False, indent=2))
            lines.append("```")
        if summary.sample_response_body is not None:
            lines.append("response body:")
            lines.append("```json")
            lines.append(json.dumps(summary.sample_response_body, ensure_ascii=False, indent=2)[:5000])
            lines.append("```")
        lines.append("")

    return "\n".join(lines)
