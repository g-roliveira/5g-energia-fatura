import os
from pathlib import Path

import pytest

from fatura.coelba_client import CoelbaClient
from fatura.config import load_config

_MES_PARA_NUMERO = {
    "JANEIRO": "01",
    "FEVEREIRO": "02",
    "MARÇO": "03",
    "ABRIL": "04",
    "MAIO": "05",
    "JUNHO": "06",
    "JULHO": "07",
    "AGOSTO": "08",
    "SETEMBRO": "09",
    "OUTUBRO": "10",
    "NOVEMBRO": "11",
    "DEZEMBRO": "12",
}


def _should_run_real_tests() -> bool:
    return os.getenv("RUN_REAL_PORTAL_TESTS") == "1"


def _load_real_config():
    config_path = Path("config.yaml")
    if not config_path.exists():
        pytest.skip("config.yaml não encontrado para testes reais")
    config = load_config(config_path)
    ativos = [cliente for cliente in config.clientes if cliente.ativo]
    if not ativos:
        pytest.skip("nenhum cliente ativo em config.yaml para testes reais")
    return config, ativos[0]


def _referencia_para_mes_ano(referencia: str) -> str:
    mes, ano = referencia.split("/", maxsplit=1)
    return f"{_MES_PARA_NUMERO[mes.upper()]}{ano}"


pytestmark = [
    pytest.mark.real_portal,
    pytest.mark.skipif(not _should_run_real_tests(), reason="RUN_REAL_PORTAL_TESTS=1 não definido"),
]


@pytest.mark.asyncio
async def test_real_portal_login_form_and_authentication():
    config, cliente = _load_real_config()

    async with CoelbaClient(config.portal) as client:
        await client.login(
            cpf_cnpj=cliente.cpf_cnpj,
            senha=cliente.senha_portal,
            uf=cliente.uf,
            tipo_acesso=cliente.tipo_acesso,
        )

        assert "#/login" not in client.page.url


@pytest.mark.asyncio
async def test_real_portal_listar_faturas():
    config, cliente = _load_real_config()

    async with CoelbaClient(config.portal) as client:
        await client.login(
            cpf_cnpj=cliente.cpf_cnpj,
            senha=cliente.senha_portal,
            uf=cliente.uf,
            tipo_acesso=cliente.tipo_acesso,
        )
        faturas = await client.listar_faturas(cliente.uc)

        assert len(faturas) >= 1
        assert any(fatura.referencia for fatura in faturas)


@pytest.mark.asyncio
async def test_real_portal_baixar_pdf_fatura_disponivel(tmp_path):
    config, cliente = _load_real_config()
    portal_config = config.portal.model_copy(
        update={"download_dir": str(tmp_path / "downloads")}
    )

    async with CoelbaClient(portal_config) as client:
        await client.login(
            cpf_cnpj=cliente.cpf_cnpj,
            senha=cliente.senha_portal,
            uf=cliente.uf,
            tipo_acesso=cliente.tipo_acesso,
        )
        faturas = await client.listar_faturas(cliente.uc)
        assert faturas

        mes_ano = _referencia_para_mes_ano(faturas[0].referencia)
        pdf_path = await client.baixar_fatura(
            uc=cliente.uc,
            mes_ano=mes_ano,
            destino_dir=str(tmp_path / "pdfs"),
        )

        assert pdf_path.exists()
        assert pdf_path.stat().st_size > 0
        assert pdf_path.read_bytes().startswith(b"%PDF")
