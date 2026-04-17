from __future__ import annotations

import asyncio
import json
from dataclasses import dataclass
from typing import Any
from urllib.parse import parse_qsl, urlparse

import httpx
from playwright.async_api import Request, Response

from fatura.coelba_client import CoelbaClient
from fatura.config import AppConfig, ClienteConfig
from fatura.exceptions import LoginError


DEFAULT_HEADERS = {
    "accept": "application/json",
    "accept-language": "pt-BR",
    "origin": "https://agenciavirtual.neoenergia.com",
    "referer": "https://agenciavirtual.neoenergia.com/",
    "user-agent": (
        "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 "
        "(KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
    ),
}


@dataclass(slots=True)
class NeoenergiaSession:
    token: str
    token_ne_se: dict[str, str]
    local_storage: dict[str, str]
    cookies: list[dict[str, Any]]


@dataclass(slots=True)
class ApiExchange:
    method: str
    url: str
    path: str
    query: dict[str, Any]
    request_headers: dict[str, Any]
    request_body: Any
    status_code: int
    response_headers: dict[str, Any]
    response_body: Any


def _try_json_loads(value: str | None) -> Any:
    if value is None:
        return None
    try:
        return json.loads(value)
    except Exception:
        return value


def _is_interesting_url(url: str) -> bool:
    host = urlparse(url).netloc.lower()
    return "neoenergia.com" in host


class NeoenergiaPrivateApiClient:
    def __init__(self, config: AppConfig, session: NeoenergiaSession) -> None:
        self._config = config
        self._session = session
        self.exchanges: list[ApiExchange] = []
        self._client = httpx.Client(
            base_url="https://apineprd.neoenergia.com",
            headers={**DEFAULT_HEADERS, "authorization": f"Bearer {session.token}"},
            timeout=45.0,
            follow_redirects=True,
        )

    def close(self) -> None:
        self._client.close()

    def _record_exchange(self, response: httpx.Response, request_body: Any) -> None:
        body: Any
        content_type = (response.headers.get("content-type") or "").lower()
        if "application/json" in content_type:
            try:
                body = response.json()
            except Exception:
                body = response.text
        else:
            body = response.text

        request = response.request
        parsed = urlparse(str(request.url))
        self.exchanges.append(
            ApiExchange(
                method=request.method,
                url=str(request.url),
                path=parsed.path,
                query=dict(parse_qsl(parsed.query, keep_blank_values=True)),
                request_headers=dict(request.headers),
                request_body=request_body,
                status_code=response.status_code,
                response_headers=dict(response.headers),
                response_body=body,
            )
        )

    def _get(self, path: str, params: dict[str, Any]) -> dict[str, Any]:
        response = self._client.get(path, params=params)
        self._record_exchange(response, None)
        response.raise_for_status()
        return response.json()

    def _post(self, path: str, payload: dict[str, Any]) -> dict[str, Any]:
        response = self._client.post(
            path,
            json=payload,
            headers={**self._client.headers, "content-type": "application/json"},
        )
        self._record_exchange(response, payload)
        response.raise_for_status()
        return response.json()

    def get_grupo_cliente(self, documento: str, distribuidora: str = "COELBA") -> dict[str, Any]:
        return self._get(
            f"/multilogin/2.0.0/agv/cliente/{documento}/{distribuidora}/grupo-de-cliente",
            {"tipoPerfil": "0"},
        )

    def list_ucs(self, documento: str, distribuidora: str = "COELBA", tipo_perfil: int = 1) -> dict[str, Any]:
        return self._get(
            f"/imoveis/1.1.0/clientes/{documento}/ucs",
            {
                "documento": documento,
                "canalSolicitante": "AGC",
                "distribuidora": distribuidora,
                "usuario": "WSO2_CONEXAO",
                "indMaisUcs": "X",
                "protocolo": "123",
                "opcaoSSOS": "S",
                "tipoPerfil": str(tipo_perfil),
            },
        )

    def get_minha_conta(self, documento: str, distribuidora: str = "COELBA", tipo_perfil: int = 1) -> dict[str, Any]:
        return self._get(
            "/multilogin/2.0.0/servicos/minha-conta",
            {
                "canalSolicitante": "AGC",
                "distribuidora": distribuidora,
                "usuario": documento,
                "tipoPerfil": str(tipo_perfil),
                "documentoSolicitante": documento,
            },
        )

    def get_minha_conta_legado(self, documento: str, distribuidora: str = "COELBA", tipo_perfil: int = 1) -> dict[str, Any]:
        return self._get(
            "/multilogin/2.0.0/servicos/minha-conta/minha-conta-legado",
            {
                "canalSolicitante": "AGC",
                "usuario": documento,
                "usuarioSap": "WSO2_CONEXAO",
                "usuarioSonda": "WSO2_CONEXAO",
                "distribuidora": distribuidora,
                "tipoPerfil": str(tipo_perfil),
                "documentoSolicitante": documento,
            },
        )

    def get_imovel(self, uc: str, distribuidora: str = "COELBA", tipo_perfil: int = 1) -> dict[str, Any]:
        return self._get(
            f"/multilogin/2.0.0/servicos/imoveis/ucs/{uc}",
            {
                "usuario": "WSO2_CONEXAO",
                "canalSolicitante": "AGC",
                "distribuidora": distribuidora,
                "protocolo": "123",
                "tipoPerfil": str(tipo_perfil),
                "opcaoSSOS": "N",
            },
        )

    def get_protocolo(
        self,
        documento: str,
        cod_cliente: str,
        distribuidora: str = "COEL",
        canal_solicitante: str = "AGC",
        regiao: str = "NE",
    ) -> dict[str, Any]:
        return self._get(
            "/protocolo/1.1.0/obterProtocolo",
            {
                "distribuidora": distribuidora,
                "canalSolicitante": canal_solicitante,
                "documento": documento,
                "codCliente": cod_cliente,
                "recaptchaAnl": "false",
                "regiao": regiao,
            },
        )

    def list_faturas(
        self,
        uc: str,
        documento: str,
        protocolo: str,
        distribuidora: str = "COELBA",
        tipo_perfil: int = 1,
    ) -> dict[str, Any]:
        return self._get(
            "/multilogin/2.0.0/servicos/faturas/ucs/faturas",
            {
                "codigo": uc,
                "documento": documento,
                "canalSolicitante": "AGC",
                "usuario": "WSO2_CONEXAO",
                "protocolo": protocolo,
                "tipificacao": "",
                "byPassActiv": "X",
                "documentoSolicitante": documento,
                "documentoCliente": documento,
                "distribuidora": distribuidora,
                "tipoPerfil": str(tipo_perfil),
            },
        )

    def get_historico_consumo(
        self,
        uc: str,
        documento: str,
        protocolo: str,
        distribuidora: str = "COELBA",
        tipo_perfil: int = 1,
    ) -> dict[str, Any]:
        return self._get(
            f"/multilogin/2.0.0/servicos/historicos/ucs/{uc}/consumos",
            {
                "canalSolicitante": "AGC",
                "usuario": "WSO2_CONEXAO",
                "dataInicioPeriodoCalc": "2021-04-18T00:00:00",
                "protocoloSonda": protocolo,
                "opcaoSSOS": "N",
                "protocolo": protocolo,
                "documentoSolicitante": documento,
                "byPassAtiv": "X",
                "distribuidora": distribuidora,
                "tipoPerfil": str(tipo_perfil),
                "codigo": uc,
            },
        )

    def get_data_certa(self, uc: str, distribuidora: str = "COELBA", tipo_perfil: int = 1) -> dict[str, Any]:
        return self._get(
            f"/multilogin/2.0.0/servicos/datacerta/ucs/{uc}/datacerta",
            {
                "codigo": uc,
                "canalSolicitante": "AGC",
                "usuario": "WSO2_CONEXAO",
                "operacao": "CON",
                "tipoPerfil": str(tipo_perfil),
                "documentoSolicitante": "",
                "distribuidora": distribuidora,
            },
        )

    def get_fatura_digital(self, uc: str, distribuidora: str = "COELBA", tipo_perfil: int = 1) -> dict[str, Any]:
        return self._get(
            "/multilogin/2.0.0/servicos/fatura-digital/ucs/fatura-digital",
            {
                "codigo": uc,
                "canalSolicitante": "AGC",
                "usuario": "WSO2_CONEXAO",
                "distribuidora": distribuidora,
                "tipoPerfil": str(tipo_perfil),
            },
        )

    def get_debito_automatico(self, uc: str, cod_cliente: str, distribuidora: str = "COELBA", tipo_perfil: int = 1) -> dict[str, Any]:
        return self._get(
            "/multilogin/2.0.0/servicos/debito-automatico/conta-cadastrada-debito",
            {
                "codigo": uc,
                "codCliente": cod_cliente,
                "canalSolicitante": "AGC",
                "usuario": "WSO2_CONEXAO",
                "valida": "",
                "distribuidora": distribuidora,
                "tipoPerfil": str(tipo_perfil),
                "documentoSolicitante": "",
            },
        )

    def get_motivos_segunda_via(self, uc: str, documento: str, distribuidora: str = "COELBA", tipo_perfil: int = 1) -> dict[str, Any]:
        return self._get(
            "/multilogin/2.0.0/servicos/faturas/lista-motivo-segundavia",
            {
                "usuario": "WSO2_CONEXAO",
                "canalSolicitante": "AGC",
                "distribuidora": distribuidora,
                "regiao": "NE",
                "tipoPerfil": str(tipo_perfil),
                "documentoSolicitante": documento,
                "codigo": uc,
            },
        )

    def get_dados_pagamento(
        self,
        uc: str,
        numero_fatura: str,
        documento: str,
        protocolo: str,
        distribuidora: str = "COELBA",
        tipo_perfil: int = 1,
    ) -> dict[str, Any]:
        return self._get(
            f"/multilogin/2.0.0/servicos/faturas/{numero_fatura}/dados-pagamento",
            {
                "codigo": uc,
                "protocolo": protocolo,
                "usuario": "WSO2_CONEXAO",
                "canalSolicitante": "AGC",
                "distribuidora": distribuidora,
                "regiao": "NE",
                "tipoPerfil": str(tipo_perfil),
                "byPassActiv": "X",
                "documentoSolicitante": documento,
                "documento": documento,
            },
        )

    def get_fatura_pdf(
        self,
        uc: str,
        numero_fatura: str,
        documento: str,
        protocolo: str,
        motivo: str,
        distribuidora: str = "COELBA",
        tipo_perfil: int = 1,
    ) -> dict[str, Any]:
        return self._get(
            f"/multilogin/2.0.0/servicos/faturas/{numero_fatura}/pdf",
            {
                "codigo": uc,
                "protocolo": protocolo,
                "tipificacao": "1031602",
                "usuario": "WSO2_CONEXAO",
                "canalSolicitante": "AGC",
                "motivo": motivo,
                "distribuidora": distribuidora,
                "regiao": "NE",
                "tipoPerfil": str(tipo_perfil),
                "documento": documento,
                "documentoSolicitante": documento,
                "byPassActiv": "",
            },
        )

    def registrar_log_atividade(
        self,
        uc: str,
        numero_fatura: str,
        documento: str,
        protocolo: str,
        distribuidora: str = "COELBA",
        tipo_perfil: int = 1,
    ) -> dict[str, Any]:
        return self._post(
            "/multilogin/2.0.0/servicos/log-atividade/registra-log-atividade",
            {
                "protocolo": protocolo,
                "numeroFatura": numero_fatura,
                "documentoFiscal": documento,
                "documentoSolicitante": documento,
                "codigo": uc,
                "canalSolicitante": "AGC",
                "tipificacao": "1010809",
                "usuario": "WSO2_CONEXAO",
                "distribuidora": distribuidora,
                "regiao": "NE",
                "recaptchaAnl": False,
                "recaptcha": "",
                "tipoPerfil": tipo_perfil,
                "semContaContrato": True,
            },
        )


async def bootstrap_session_with_playwright(config: AppConfig, cliente: ClienteConfig) -> NeoenergiaSession:
    async with CoelbaClient(config.portal) as client:
        await client.login(
            cpf_cnpj=cliente.cpf_cnpj,
            senha=cliente.senha_portal,
            uf=cliente.uf,
            tipo_acesso=cliente.tipo_acesso,
        )
        storage_state = await client.context.storage_state()
        origin = next(
            (
                item
                for item in storage_state.get("origins", [])
                if item.get("origin") == "https://agenciavirtual.neoenergia.com"
            ),
            None,
        )
        if not origin:
            raise LoginError("Storage state sem origin da Agência Virtual.")
        local_storage = {item["name"]: item["value"] for item in origin.get("localStorage", [])}
        raw_token = local_storage.get("token")
        token_ne_se = _try_json_loads(local_storage.get("tokenNeSe", "{}")) or {}
        token_value = _try_json_loads(raw_token)
        token = token_value if isinstance(token_value, str) else token_ne_se.get("ne")
        if not token:
            raise LoginError("Token Bearer não encontrado no localStorage após login.")
        cookies = await client.context.cookies()
        return NeoenergiaSession(
            token=token,
            token_ne_se=token_ne_se,
            local_storage=local_storage,
            cookies=cookies,
        )


async def bootstrap_session_and_capture_with_playwright(
    config: AppConfig,
    cliente: ClienteConfig,
) -> tuple[NeoenergiaSession, list[ApiExchange]]:
    requests: dict[int, dict[str, Any]] = {}
    exchanges: list[ApiExchange] = []
    pending: set[asyncio.Task[Any]] = set()

    def schedule(coro: Any) -> None:
        task = asyncio.create_task(coro)
        pending.add(task)
        task.add_done_callback(pending.discard)

    async def on_request(request: Request) -> None:
        if not _is_interesting_url(request.url):
            return
        parsed = urlparse(request.url)
        requests[id(request)] = {
            "method": request.method,
            "url": request.url,
            "path": parsed.path,
            "query": dict(parse_qsl(parsed.query, keep_blank_values=True)),
            "headers": dict(request.headers),
            "body": _try_json_loads(request.post_data),
        }

    async def on_response(response: Response) -> None:
        request = response.request
        if not _is_interesting_url(request.url):
            return
        req_data = requests.get(id(request), {})
        content_type = (response.headers.get("content-type") or "").lower()
        body: Any
        try:
            if "application/json" in content_type:
                body = await response.json()
            elif "application/pdf" in content_type:
                body = "<binary pdf>"
            else:
                body = (await response.text())[:4000]
        except Exception as exc:
            body = f"<unreadable: {exc}>"
        exchanges.append(
            ApiExchange(
                method=req_data.get("method", request.method),
                url=req_data.get("url", request.url),
                path=req_data.get("path", urlparse(request.url).path),
                query=req_data.get("query", {}),
                request_headers=req_data.get("headers", {}),
                request_body=req_data.get("body"),
                status_code=response.status,
                response_headers=dict(response.headers),
                response_body=body,
            )
        )

    async with CoelbaClient(config.portal) as client:
        client.page.on("request", lambda req: schedule(on_request(req)))
        client.page.on("response", lambda resp: schedule(on_response(resp)))

        await client.login(
            cpf_cnpj=cliente.cpf_cnpj,
            senha=cliente.senha_portal,
            uf=cliente.uf,
            tipo_acesso=cliente.tipo_acesso,
        )

        await asyncio.sleep(2)
        if pending:
            await asyncio.gather(*pending, return_exceptions=True)

        storage_state = await client.context.storage_state()
        origin = next(
            (
                item
                for item in storage_state.get("origins", [])
                if item.get("origin") == "https://agenciavirtual.neoenergia.com"
            ),
            None,
        )
        if not origin:
            raise LoginError("Storage state sem origin da Agência Virtual.")
        local_storage = {item["name"]: item["value"] for item in origin.get("localStorage", [])}
        raw_token = local_storage.get("token")
        token_ne_se = _try_json_loads(local_storage.get("tokenNeSe", "{}")) or {}
        token_value = _try_json_loads(raw_token)
        token = token_value if isinstance(token_value, str) else token_ne_se.get("ne")
        if not token:
            raise LoginError("Token Bearer não encontrado no localStorage após login.")
        cookies = await client.context.cookies()
        session = NeoenergiaSession(
            token=token,
            token_ne_se=token_ne_se,
            local_storage=local_storage,
            cookies=cookies,
        )
        return session, exchanges


def bootstrap_session(config: AppConfig, cliente: ClienteConfig) -> NeoenergiaSession:
    return asyncio.run(bootstrap_session_with_playwright(config, cliente))


def bootstrap_session_and_capture(
    config: AppConfig,
    cliente: ClienteConfig,
) -> tuple[NeoenergiaSession, list[ApiExchange]]:
    return asyncio.run(bootstrap_session_and_capture_with_playwright(config, cliente))
