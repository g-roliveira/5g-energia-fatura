"""Documenta a API privada da Neoenergia a partir de chamadas reais.

Uso:
    ./.venv/bin/python scripts/documentar_api_neoenergia.py
    ./.venv/bin/python scripts/documentar_api_neoenergia.py --uc 007098175908
"""

from __future__ import annotations

import argparse
import json
from dataclasses import asdict
from datetime import datetime
from pathlib import Path

import httpx

from fatura.config import ClienteConfig, load_config
from fatura.neoenergia_private_api import (
    ApiExchange,
    NeoenergiaPrivateApiClient,
    bootstrap_session_and_capture,
)
from fatura.neoenergia_private_api_docs import (
    build_endpoint_summaries,
    render_markdown,
    sanitize_data,
    save_exchanges,
)


DOCS_ROOT = Path("docs/neoenergia-private-api")
API_PATH_PREFIXES = (
    "/areanaologada/",
    "/multilogin/",
    "/imoveis/",
    "/protocolo/",
)


def _now_slug() -> str:
    return datetime.now().strftime("%Y%m%d_%H%M%S")


def _digits(value: str) -> str:
    return "".join(ch for ch in value if ch.isdigit())


def _pick_cliente(clientes: list[ClienteConfig], nome: str | None) -> ClienteConfig:
    ativos = [cliente for cliente in clientes if cliente.ativo]
    if not ativos:
        raise SystemExit("Nenhum cliente ativo em config.yaml.")
    if not nome:
        return ativos[0]
    for cliente in ativos:
        if cliente.nome == nome:
            return cliente
    raise SystemExit(f"Cliente ativo não encontrado: {nome}")


def _safe_call(label: str, fn, *args, **kwargs):
    try:
        return sanitize_data(fn(*args, **kwargs))
    except httpx.HTTPStatusError as exc:
        response = exc.response
        try:
            body = response.json()
        except Exception:
            body = response.text
        return {
            "_error": {
                "label": label,
                "status_code": response.status_code,
                "path": response.request.url.path,
                "message": str(exc),
                "response": sanitize_data(body),
            }
        }


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--config", default="config.yaml")
    parser.add_argument("--cliente", default=None, help="Nome do cliente em config.yaml")
    parser.add_argument("--uc", action="append", dest="ucs", default=None)
    args = parser.parse_args()

    config = load_config(args.config)
    cliente = _pick_cliente(config.clientes, args.cliente)
    documento = _digits(cliente.cpf_cnpj)

    run_dir = DOCS_ROOT / f"live_{_now_slug()}"
    run_dir.mkdir(parents=True, exist_ok=True)

    session, login_exchanges = bootstrap_session_and_capture(config, cliente)
    api = NeoenergiaPrivateApiClient(config, session)
    all_exchanges: list[ApiExchange] = list(login_exchanges)

    try:
        grupo = _safe_call("grupo_cliente", api.get_grupo_cliente, documento)
        minha_conta = _safe_call("minha_conta", api.get_minha_conta, documento)
        minha_conta_legado = _safe_call("minha_conta_legado", api.get_minha_conta_legado, documento)
        ucs_response_raw = api.list_ucs(documento)
        ucs_response = sanitize_data(ucs_response_raw)

        ucs = ucs_response_raw.get("ucs", [])
        requested_ucs = set(args.ucs or [])
        target_ucs = [uc for uc in ucs if not requested_ucs or uc.get("uc") in requested_ucs]

        per_uc: dict[str, dict[str, object]] = {}
        for uc_info in target_ucs:
            uc = str(uc_info["uc"])
            protocolo = api.get_protocolo(documento=documento, cod_cliente=uc)
            protocolo_str = (
                protocolo.get("protocoloSalesforceStr")
                or protocolo.get("protocoloLegadoStr")
                or ""
            )
            imovel_raw = api.get_imovel(uc)
            historico = _safe_call(
                "historico_consumo",
                api.get_historico_consumo,
                uc,
                documento=documento,
                protocolo=protocolo_str,
            )
            data_certa = _safe_call("data_certa", api.get_data_certa, uc)
            fatura_digital = _safe_call("fatura_digital", api.get_fatura_digital, uc)
            cod_cliente = (
                imovel_raw.get("cliente", {}).get("codigo")
                or uc_info.get("contrato")
                or ""
            )
            debito_automatico = _safe_call(
                "debito_automatico",
                api.get_debito_automatico,
                uc,
                cod_cliente=str(cod_cliente),
            )
            faturas_raw = api.list_faturas(uc, documento=documento, protocolo=protocolo_str)
            faturas = sanitize_data(faturas_raw)
            motivos_raw = api.get_motivos_segunda_via(uc, documento=documento)
            motivos = sanitize_data(motivos_raw)

            uc_artifacts: dict[str, object] = {
                "protocolo": sanitize_data(protocolo),
                "imovel": sanitize_data(imovel_raw),
                "historico_consumo": historico,
                "data_certa": data_certa,
                "fatura_digital": fatura_digital,
                "debito_automatico": debito_automatico,
                "faturas": faturas,
                "motivos_segunda_via": motivos,
            }

            lista_faturas = faturas_raw.get("faturas", []) or []
            if lista_faturas:
                numero_fatura = str(lista_faturas[0]["numeroFatura"])
                dados_pagamento = _safe_call(
                    "dados_pagamento",
                    api.get_dados_pagamento,
                    uc,
                    numero_fatura=numero_fatura,
                    documento=documento,
                    protocolo=protocolo_str,
                )
                motivos_lista = motivos_raw.get("motivos", []) or []
                motivo = str(motivos_lista[0]["idMotivo"]) if motivos_lista else "02"
                pdf = _safe_call(
                    "pdf_primeira_fatura",
                    api.get_fatura_pdf,
                    uc,
                    numero_fatura=numero_fatura,
                    documento=documento,
                    protocolo=protocolo_str,
                    motivo=motivo,
                )
                uc_artifacts["dados_pagamento_primeira_fatura"] = dados_pagamento
                uc_artifacts["pdf_primeira_fatura"] = pdf

            per_uc[uc] = uc_artifacts

        all_exchanges.extend(api.exchanges)
    finally:
        api.close()

    save_exchanges(run_dir / "exchanges", all_exchanges)
    api_only_exchanges = [exchange for exchange in all_exchanges if exchange.path.startswith(API_PATH_PREFIXES)]
    save_exchanges(run_dir / "api_exchanges", api_only_exchanges)

    (run_dir / "session.json").write_text(
        json.dumps(
            sanitize_data(
                {
                    "token_ne_se": session.token_ne_se,
                    "local_storage": session.local_storage,
                    "cookies": session.cookies,
                }
            ),
            ensure_ascii=False,
            indent=2,
        ),
        encoding="utf-8",
    )
    (run_dir / "grupo_cliente.json").write_text(
        json.dumps(sanitize_data(grupo), ensure_ascii=False, indent=2),
        encoding="utf-8",
    )
    (run_dir / "minha_conta.json").write_text(
        json.dumps(sanitize_data(minha_conta), ensure_ascii=False, indent=2),
        encoding="utf-8",
    )
    (run_dir / "minha_conta_legado.json").write_text(
        json.dumps(sanitize_data(minha_conta_legado), ensure_ascii=False, indent=2),
        encoding="utf-8",
    )
    (run_dir / "ucs.json").write_text(
        json.dumps(sanitize_data(ucs_response), ensure_ascii=False, indent=2),
        encoding="utf-8",
    )
    (run_dir / "per_uc.json").write_text(
        json.dumps(per_uc, ensure_ascii=False, indent=2),
        encoding="utf-8",
    )

    summaries = build_endpoint_summaries(api_only_exchanges)
    summary_payload = {
        "generated_at": datetime.now().isoformat(),
        "cliente": cliente.nome,
        "documento": sanitize_data(documento, key="documento"),
        "state_selection_observation": (
            "Nenhum endpoint dedicado de seleção de estado foi observado no fluxo capturado. "
            "Após o login, as chamadas já seguem com `distribuidora=COELBA`, o que indica "
            "que a escolha de estado é resolvida no frontend ou por configuração local de sessão."
        ),
        "ucs": sanitize_data(ucs),
        "endpoints": [asdict(summary) for summary in summaries],
    }
    (run_dir / "summary.json").write_text(
        json.dumps(summary_payload, ensure_ascii=False, indent=2),
        encoding="utf-8",
    )

    markdown = render_markdown(
        generated_at=summary_payload["generated_at"],
        cliente_nome=cliente.nome,
        documento=documento,
        ucs=ucs,
        summaries=summaries,
        state_selection_observation=summary_payload["state_selection_observation"],
        output_dir=run_dir,
    )
    (run_dir / "README.md").write_text(markdown, encoding="utf-8")
    (DOCS_ROOT / "README.md").write_text(markdown, encoding="utf-8")
    (DOCS_ROOT / "latest_run.txt").write_text(str(run_dir), encoding="utf-8")

    print(run_dir)


if __name__ == "__main__":
    main()
