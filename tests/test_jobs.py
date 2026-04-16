from datetime import date
from decimal import Decimal
from pathlib import Path

import pytest

from fatura.config import AppConfig, ClienteConfig, DatabaseConfig, PortalConfig, ServiceConfig, TipoAcesso
from fatura.exceptions import DownloadError
from fatura.jobs import BatchProcessor, BatchSpec, BatchTarget, processar_faturas_mes
from fatura.models import Cliente, ContaDistribuidora
from fatura.repository import SqliteFaturaRepository


class FakeParser:
    def parse(self, pdf_path: Path | str) -> ContaDistribuidora:
        return ContaDistribuidora(
            uc="7085489032",
            mes=12,
            ano=2024,
            valor=Decimal("113.12"),
            vencimento=date(2025, 1, 8),
            cliente=Cliente(nome="Cliente Teste"),
            pdf_path=str(pdf_path),
        )


class FakeClient:
    failing_ucs: set[str] = set()

    def __init__(self, config: PortalConfig):
        self.config = config
        self._last_error_context = {
            "step_name": "baixar_fatura",
            "screenshot_path": str(Path(config.download_dir) / "_errors" / "fake.png"),
            "html_path": str(Path(config.download_dir) / "_errors" / "fake.html"),
        }

    async def __aenter__(self):
        return self

    async def __aexit__(self, *args):
        return None

    @property
    def last_error_context(self):
        return self._last_error_context

    async def login(self, cpf_cnpj: str, senha: str, uf: str = "BA", tipo_acesso: TipoAcesso = TipoAcesso.NORMAL):
        return None

    async def baixar_fatura(self, uc: str, mes_ano: str | None = None, destino_dir: str | None = None):
        if uc in self.failing_ucs:
            raise DownloadError(f"falha download {uc}", uc=uc)
        return Path(self.config.download_dir) / f"{uc}.pdf"

    async def delay_entre_clientes(self):
        return None


def build_config(tmp_path: Path) -> AppConfig:
    return AppConfig(
        portal=PortalConfig(download_dir=str(tmp_path / "downloads"), headless=True),
        database=DatabaseConfig(url=f"sqlite:///{tmp_path / 'faturas.db'}"),
        service=ServiceConfig(api_key="secret"),
        clientes=[],
    )


@pytest.mark.asyncio
async def test_batch_processor_persists_success_and_error(tmp_path: Path):
    config = build_config(tmp_path)
    repo = SqliteFaturaRepository(config.database.url)
    FakeClient.failing_ucs = {"2002"}
    processor = BatchProcessor(
        config=config,
        repo=repo,
        parser=FakeParser(),
        client_factory=FakeClient,
    )

    spec = BatchSpec(
        cpf_cnpj="12345678901",
        senha_portal="senha",
        uf="BA",
        tipo_acesso=TipoAcesso.NORMAL,
        targets=[BatchTarget(uc="1001", nome="UC 1"), BatchTarget(uc="2002", nome="UC 2")],
        mes_ano="122024",
    )

    result = await processor.run_batch(spec)

    assert result.status == "partial_failure"
    assert result.success == 1
    assert result.error == 1
    assert repo.conta_existe("7085489032", 12, 2024) is True
    failing = next(item for item in result.items if item.uc == "2002")
    assert failing.status == "erro_download"
    assert failing.screenshot_path is not None


@pytest.mark.asyncio
async def test_cli_processar_faturas_reuses_batch_orchestration(tmp_path: Path):
    config = build_config(tmp_path)
    FakeClient.failing_ucs = set()
    clientes = [
        ClienteConfig(
            nome="Cliente A",
            uc="1001",
            cpf_cnpj="12345678901",
            senha_portal="senha",
        )
    ]

    original_factory = BatchProcessor.__init__

    def init_with_fake(self, config, repo=None, parser=None, client_factory=FakeClient):
        return original_factory(self, config, repo=repo, parser=FakeParser(), client_factory=FakeClient)

    try:
        BatchProcessor.__init__ = init_with_fake  # type: ignore[method-assign]
        result = await processar_faturas_mes(config=config, clientes=clientes, mes_ano="122024", force=False)
    finally:
        BatchProcessor.__init__ = original_factory  # type: ignore[method-assign]

    assert result.total == 1
    assert result.sucesso == 1
    assert result.erro == 0
