from __future__ import annotations

from pathlib import Path

from fatura.neoenergia_private_api import ApiExchange
from fatura.neoenergia_private_api_docs import (
    build_endpoint_summaries,
    render_markdown,
    sanitize_exchange,
)


def test_sanitize_exchange_masks_sensitive_fields() -> None:
    exchange = ApiExchange(
        method="POST",
        url="https://avapineanl.neoenergia.com/areanaologada/2.0.0/autentica",
        path="/areanaologada/2.0.0/autentica",
        query={},
        request_headers={"authorization": "Bearer abc", "content-type": "application/json"},
        request_body={
            "usuario": "03021937586",
            "senha": "segredo",
            "recaptcha": "token-grande",
        },
        status_code=200,
        response_headers={"set-cookie": "jwt=abc"},
        response_body={"token": "secret", "fileData": "abc123"},
    )

    sanitized = sanitize_exchange(exchange)

    assert sanitized["request_headers"]["authorization"] == "<redacted>"
    assert sanitized["request_body"]["usuario"] == "*******7586"
    assert sanitized["request_body"]["senha"] == "<redacted>"
    assert sanitized["request_body"]["recaptcha"] == "<redacted>"
    assert sanitized["response_headers"]["set-cookie"] == "<redacted>"
    assert sanitized["response_body"]["token"] == "<redacted>"
    assert sanitized["response_body"]["fileData"] == "<base64:6 chars>"


def test_build_endpoint_summaries_and_markdown() -> None:
    exchanges = [
        ApiExchange(
            method="GET",
            url="https://apineprd.neoenergia.com/imoveis/1.1.0/clientes/03021937586/ucs",
            path="/imoveis/1.1.0/clientes/03021937586/ucs",
            query={"documento": "03021937586"},
            request_headers={"authorization": "Bearer abc"},
            request_body=None,
            status_code=200,
            response_headers={"content-type": "application/json"},
            response_body={"ucs": [{"uc": "007098175908", "status": "LIGADA"}]},
        ),
        ApiExchange(
            method="GET",
            url="https://apineprd.neoenergia.com/imoveis/1.1.0/clientes/03021937586/ucs",
            path="/imoveis/1.1.0/clientes/03021937586/ucs",
            query={"documento": "03021937586"},
            request_headers={"authorization": "Bearer def"},
            request_body=None,
            status_code=200,
            response_headers={"content-type": "application/json"},
            response_body={"ucs": [{"uc": "007085489032", "status": "DESLIGADA"}]},
        ),
    ]

    summaries = build_endpoint_summaries(exchanges)

    assert len(summaries) == 1
    assert summaries[0].count == 2
    assert summaries[0].sample_query == {"documento": "*******7586"}

    markdown = render_markdown(
        generated_at="2026-04-16T22:00:00",
        cliente_nome="Cliente Teste",
        documento="03021937586",
        ucs=[
            {"uc": "007098175908", "status": "LIGADA", "local": {"endereco": "Rua X", "municipio": "Lapão"}}
        ],
        summaries=summaries,
        state_selection_observation="Nenhum endpoint dedicado foi observado.",
        output_dir=Path("docs/neoenergia-private-api/live_fake"),
    )

    assert "Neoenergia Private API" in markdown
    assert "Nenhum endpoint dedicado foi observado." in markdown
    assert "/imoveis/1.1.0/clientes/03021937586/ucs" in markdown
    assert "*******7586" in markdown
