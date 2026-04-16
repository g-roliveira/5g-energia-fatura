from pathlib import Path

from fastapi.testclient import TestClient

from fatura.api import create_app
from fatura.config import AppConfig, DatabaseConfig, PortalConfig, ServiceConfig
from fatura.repository import SqliteFaturaRepository
from fatura.service_models import BatchItemResult


def build_config(tmp_path: Path) -> AppConfig:
    return AppConfig(
        portal=PortalConfig(download_dir=str(tmp_path / "downloads"), headless=True),
        database=DatabaseConfig(url=f"sqlite:///{tmp_path / 'api.db'}"),
        service=ServiceConfig(api_key="secret", max_concurrent_jobs=1),
        clientes=[],
    )


def test_api_creates_job_and_exposes_status_and_result(tmp_path: Path):
    config = build_config(tmp_path)
    app = create_app(config=config)
    repo = SqliteFaturaRepository(config.database.url)
    app.state.runtime._repo = repo

    async def fake_run_job(job_id: str) -> None:
        request = repo.carregar_job_request(job_id)
        repo.marcar_job_em_execucao(job_id)
        for target in request.ucs:
            repo.marcar_item_em_execucao(job_id, target.uc, attempts=1)
            repo.salvar_resultado_item(
                job_id,
                BatchItemResult(
                    uc=target.uc,
                    nome=target.nome,
                    status="sucesso",
                    mensagem="ok",
                    pdf_path=f"/tmp/{target.uc}.pdf",
                    mes=12,
                    ano=2024,
                    attempts=1,
                ),
            )
        repo.finalizar_job(job_id, "succeeded")

    app.state.runtime.run_job = fake_run_job

    with TestClient(app) as client:
        payload = {
            "cpf_cnpj": "12345678901",
            "senha_portal": "senha",
            "uf": "BA",
            "tipo_acesso": "normal",
            "mes_ano": "122024",
            "ucs": [{"uc": "1001", "nome": "UC API"}],
        }

        create_response = client.post(
            "/jobs/faturas",
            headers={"X-API-Key": "secret"},
            json=payload,
        )
        assert create_response.status_code == 200
        job_id = create_response.json()["job_id"]

        status_response = client.get(f"/jobs/{job_id}", headers={"X-API-Key": "secret"})
        assert status_response.status_code == 200
        assert status_response.json()["status"] == "succeeded"

        result_response = client.get(f"/jobs/{job_id}/result", headers={"X-API-Key": "secret"})
        assert result_response.status_code == 200
        body = result_response.json()
        assert body["status"] == "succeeded"
        assert body["items"][0]["uc"] == "1001"
        assert body["items"][0]["status"] == "sucesso"


def test_api_requires_key_when_configured(tmp_path: Path):
    app = create_app(config=build_config(tmp_path))

    with TestClient(app) as client:
        response = client.get("/jobs")
        assert response.status_code == 401
