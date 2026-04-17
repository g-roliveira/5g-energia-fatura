"""Mapeia a cadeia real de APIs usada pela Agência Virtual Neoenergia.

Uso:
    ./.venv/bin/python scripts/mapear_api_neoenergia.py --uc 007098175908
"""

from __future__ import annotations

import argparse
import asyncio
import json
from collections import defaultdict
from datetime import datetime
from pathlib import Path
from typing import Any
from urllib.parse import parse_qsl, urlparse

from playwright.async_api import Request, Response

from fatura.coelba_client import CoelbaClient
from fatura.config import load_config

OUTPUT_ROOT = Path("output/playwright")


def _now_slug() -> str:
    return datetime.now().strftime("%Y%m%d_%H%M%S")


def _safe_json_loads(raw: str) -> Any:
    try:
        return json.loads(raw)
    except Exception:
        return None


class NetworkRecorder:
    def __init__(self, out_dir: Path) -> None:
        self.out_dir = out_dir
        self.events_path = out_dir / "network_events.jsonl"
        self.requests: dict[str, dict[str, Any]] = {}
        self.responses: list[dict[str, Any]] = []
        self.summary: dict[tuple[str, str], dict[str, Any]] = defaultdict(
            lambda: {
                "count": 0,
                "statuses": set(),
                "sample_query": None,
                "sample_request_headers": None,
                "sample_request_body": None,
                "sample_response_headers": None,
                "sample_response_json_keys": None,
                "sample_response_preview": None,
            }
        )

    @staticmethod
    def _interesting(url: str) -> bool:
        host = urlparse(url).netloc.lower()
        return "neoenergia.com" in host

    async def on_request(self, request: Request) -> None:
        if not self._interesting(request.url):
            return
        parsed = urlparse(request.url)
        headers = dict(request.headers)
        body = request.post_data
        event = {
            "kind": "request",
            "timestamp": datetime.now().isoformat(),
            "method": request.method,
            "url": request.url,
            "path": parsed.path,
            "query": dict(parse_qsl(parsed.query, keep_blank_values=True)),
            "headers": headers,
            "post_data": _safe_json_loads(body) if body else body,
            "resource_type": request.resource_type,
        }
        self.requests[id(request)] = event
        with self.events_path.open("a", encoding="utf-8") as handle:
            handle.write(json.dumps(event, ensure_ascii=False) + "\n")

    async def on_response(self, response: Response) -> None:
        request = response.request
        if not self._interesting(request.url):
            return

        parsed = urlparse(request.url)
        req_entry = self.requests.get(id(request), {})
        content_type = (response.headers.get("content-type") or "").lower()
        body_preview: str | None = None
        json_keys: list[str] | None = None

        try:
            if "application/json" in content_type:
                payload = await response.json()
                if isinstance(payload, dict):
                    json_keys = sorted(payload.keys())
                body_preview = json.dumps(payload, ensure_ascii=False)[:4000]
            elif "application/pdf" in content_type:
                body_preview = "<binary pdf>"
            else:
                text = await response.text()
                body_preview = text[:2000]
        except Exception as exc:
            body_preview = f"<unreadable: {exc}>"

        event = {
            "kind": "response",
            "timestamp": datetime.now().isoformat(),
            "method": request.method,
            "url": request.url,
            "path": parsed.path,
            "query": dict(parse_qsl(parsed.query, keep_blank_values=True)),
            "status": response.status,
            "request_headers": req_entry.get("headers"),
            "request_body": req_entry.get("post_data"),
            "response_headers": dict(response.headers),
            "content_type": content_type,
            "response_json_keys": json_keys,
            "response_preview": body_preview,
        }
        self.responses.append(event)
        with self.events_path.open("a", encoding="utf-8") as handle:
            handle.write(json.dumps(event, ensure_ascii=False) + "\n")

        key = (request.method, parsed.path)
        summary = self.summary[key]
        summary["count"] += 1
        summary["statuses"].add(response.status)
        summary["sample_query"] = summary["sample_query"] or event["query"]
        summary["sample_request_headers"] = summary["sample_request_headers"] or event["request_headers"]
        summary["sample_request_body"] = summary["sample_request_body"] or event["request_body"]
        summary["sample_response_headers"] = summary["sample_response_headers"] or event["response_headers"]
        summary["sample_response_json_keys"] = summary["sample_response_json_keys"] or event["response_json_keys"]
        summary["sample_response_preview"] = summary["sample_response_preview"] or event["response_preview"]

    def write_summary(self) -> None:
        serializable = []
        for (method, path), data in sorted(self.summary.items(), key=lambda item: item[0][1]):
            serializable.append(
                {
                    "method": method,
                    "path": path,
                    "count": data["count"],
                    "statuses": sorted(data["statuses"]),
                    "sample_query": data["sample_query"],
                    "sample_request_headers": data["sample_request_headers"],
                    "sample_request_body": data["sample_request_body"],
                    "sample_response_headers": data["sample_response_headers"],
                    "sample_response_json_keys": data["sample_response_json_keys"],
                    "sample_response_preview": data["sample_response_preview"],
                }
            )

        (self.out_dir / "endpoint_summary.json").write_text(
            json.dumps(serializable, ensure_ascii=False, indent=2),
            encoding="utf-8",
        )

        lines = ["# Neoenergia Endpoint Map", ""]
        for entry in serializable:
            lines.append(f"## `{entry['method']} {entry['path']}`")
            lines.append(f"- chamadas: `{entry['count']}`")
            lines.append(f"- status: `{entry['statuses']}`")
            if entry["sample_query"]:
                lines.append(f"- query: `{json.dumps(entry['sample_query'], ensure_ascii=False)}`")
            auth_header = None
            if entry["sample_request_headers"]:
                auth_header = entry["sample_request_headers"].get("authorization")
            if auth_header:
                lines.append("- authorization header: presente")
            if entry["sample_request_body"] is not None:
                lines.append("```json")
                lines.append(json.dumps(entry["sample_request_body"], ensure_ascii=False, indent=2)[:3000])
                lines.append("```")
            if entry["sample_response_json_keys"]:
                lines.append(f"- response keys: `{entry['sample_response_json_keys']}`")
            if entry["sample_response_preview"]:
                lines.append("```text")
                lines.append(str(entry["sample_response_preview"])[:3000])
                lines.append("```")
            lines.append("")

        (self.out_dir / "endpoint_summary.md").write_text("\n".join(lines), encoding="utf-8")


async def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--uc", required=True, help="Unidade consumidora a explorar.")
    parser.add_argument("--mes-ano", default=None, help="Competência MMAAAA para tentar baixar.")
    args = parser.parse_args()

    config = load_config("config.yaml")
    ativos = [cliente for cliente in config.clientes if cliente.ativo]
    if not ativos:
        raise SystemExit("Nenhum cliente ativo em config.yaml.")

    cred = ativos[0]
    run_dir = OUTPUT_ROOT / f"neoenergia_api_map_{args.uc}_{_now_slug()}"
    run_dir.mkdir(parents=True, exist_ok=True)
    recorder = NetworkRecorder(run_dir)

    async with CoelbaClient(config.portal) as client:
        client.page.on("request", lambda req: asyncio.create_task(recorder.on_request(req)))
        client.page.on("response", lambda resp: asyncio.create_task(recorder.on_response(resp)))

        await client.login(
            cpf_cnpj=cred.cpf_cnpj,
            senha=cred.senha_portal,
            uf=cred.uf,
            tipo_acesso=cred.tipo_acesso,
        )

        await client.context.storage_state(path=run_dir / "storage_state.json")
        cookies_after_login = await client.context.cookies()
        (run_dir / "cookies_after_login.json").write_text(
            json.dumps(cookies_after_login, ensure_ascii=False, indent=2),
            encoding="utf-8",
        )

        faturas = await client.listar_faturas(args.uc)
        (run_dir / "faturas_listadas.json").write_text(
            json.dumps([f.__dict__ for f in faturas], ensure_ascii=False, indent=2),
            encoding="utf-8",
        )

        pdf_path = None
        if faturas:
            target_mes_ano = args.mes_ano
            if target_mes_ano is None:
                ref = faturas[0].referencia
                meses = {
                    "JANEIRO": "01", "FEVEREIRO": "02", "MARÇO": "03", "ABRIL": "04",
                    "MAIO": "05", "JUNHO": "06", "JULHO": "07", "AGOSTO": "08",
                    "SETEMBRO": "09", "OUTUBRO": "10", "NOVEMBRO": "11", "DEZEMBRO": "12",
                }
                mes_nome, ano = ref.split("/", maxsplit=1)
                target_mes_ano = f"{meses[mes_nome.upper()]}{ano}"
            pdf_path = await client.baixar_fatura(args.uc, mes_ano=target_mes_ano, destino_dir=str(run_dir / "pdfs"))

        await client.page.screenshot(path=run_dir / "final_page.png", full_page=True)
        (run_dir / "final_page.html").write_text(await client.page.content(), encoding="utf-8")

        cookies_final = await client.context.cookies()
        (run_dir / "cookies_final.json").write_text(
            json.dumps(cookies_final, ensure_ascii=False, indent=2),
            encoding="utf-8",
        )

        await asyncio.sleep(2)
        recorder.write_summary()

        meta = {
            "uc": args.uc,
            "mes_ano_download": args.mes_ano,
            "final_url": client.page.url,
            "faturas_encontradas": len(faturas),
            "pdf_path": str(pdf_path) if pdf_path else None,
            "output_dir": str(run_dir),
        }
        (run_dir / "meta.json").write_text(json.dumps(meta, ensure_ascii=False, indent=2), encoding="utf-8")
        print(json.dumps(meta, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    asyncio.run(main())
