import json
import os
import socket
import subprocess
import sys
import time
from pathlib import Path

import httpx
import pytest
import yaml

from fatura.config import load_config


def _should_run_real_api_e2e() -> bool:
    return os.getenv("RUN_REAL_API_E2E") == "1"


def _find_free_port() -> int:
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        sock.bind(("127.0.0.1", 0))
        sock.listen(1)
        return int(sock.getsockname()[1])


def _load_real_client():
    config_path = Path("config.yaml")
    if not config_path.exists():
        pytest.skip("config.yaml não encontrado para teste E2E real da API")
    config = load_config(config_path)
    ativos = [cliente for cliente in config.clientes if cliente.ativo]
    if not ativos:
        pytest.skip("nenhum cliente ativo em config.yaml para teste E2E real da API")
    return ativos[0]


def _write_temp_config(tmp_path: Path, port: int, api_key: str) -> Path:
    config_data = yaml.safe_load(Path("config.yaml").read_text(encoding="utf-8"))
    config_data["service"] = {
        **config_data.get("service", {}),
        "host": "127.0.0.1",
        "port": port,
        "api_key": api_key,
        "max_concurrent_jobs": 1,
    }
    config_data["database"] = {
        "url": f"sqlite:///{tmp_path / 'api-e2e.db'}",
    }
    config_data["portal"] = {
        **config_data.get("portal", {}),
        "download_dir": str(tmp_path / "downloads"),
    }
    temp_config = tmp_path / "config.e2e.yaml"
    temp_config.write_text(yaml.safe_dump(config_data, sort_keys=False), encoding="utf-8")
    return temp_config


def _wait_for_server(base_url: str, timeout_s: float = 30.0) -> None:
    deadline = time.monotonic() + timeout_s
    while time.monotonic() < deadline:
        try:
            response = httpx.get(f"{base_url}/health", timeout=5.0)
            if response.status_code == 200:
                return
        except httpx.HTTPError:
            pass
        time.sleep(0.5)
    raise AssertionError("API não ficou pronta dentro do tempo esperado.")


pytestmark = [
    pytest.mark.real_api_e2e,
    pytest.mark.skipif(
        not _should_run_real_api_e2e(),
        reason="RUN_REAL_API_E2E=1 não definido",
    ),
]


def test_real_api_e2e_job_flow(tmp_path: Path):
    cliente = _load_real_client()
    port = _find_free_port()
    api_key = "e2e-secret"
    temp_config = _write_temp_config(tmp_path, port=port, api_key=api_key)
    base_url = f"http://127.0.0.1:{port}"

    process = subprocess.Popen(
        [sys.executable, "-c", "from fatura.api import main; main()"],
        cwd=Path.cwd(),
        env={**os.environ, "FATURA_CONFIG": str(temp_config)},
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
    )

    try:
        _wait_for_server(base_url)

        with httpx.Client(timeout=30.0) as client:
            docs_response = client.get(f"{base_url}/docs")
            assert docs_response.status_code == 200

            openapi_response = client.get(f"{base_url}/openapi.json")
            assert openapi_response.status_code == 200

            create_response = client.post(
                f"{base_url}/jobs/faturas",
                headers={"X-API-Key": api_key},
                json={
                    "cpf_cnpj": cliente.cpf_cnpj,
                    "senha_portal": cliente.senha_portal,
                    "uf": cliente.uf,
                    "tipo_acesso": cliente.tipo_acesso.value,
                    "mes_ano": "122024",
                    "ucs": [{"uc": cliente.uc, "nome": cliente.nome}],
                },
            )
            assert create_response.status_code == 200, create_response.text

            job = create_response.json()
            job_id = job["job_id"]
            assert job["status"] in {"queued", "running"}

            deadline = time.monotonic() + 180
            final_status = None
            while time.monotonic() < deadline:
                status_response = client.get(
                    f"{base_url}/jobs/{job_id}",
                    headers={"X-API-Key": api_key},
                )
                assert status_response.status_code == 200, status_response.text
                final_status = status_response.json()
                if final_status["status"] in {"succeeded", "partial_failure", "failed"}:
                    break
                time.sleep(5)

            assert final_status is not None
            assert final_status["status"] == "succeeded", json.dumps(final_status, ensure_ascii=False)

            result_response = client.get(
                f"{base_url}/jobs/{job_id}/result",
                headers={"X-API-Key": api_key},
            )
            assert result_response.status_code == 200, result_response.text

            result = result_response.json()
            assert result["status"] == "succeeded"
            assert len(result["items"]) == 1
            item = result["items"][0]
            assert item["uc"] == cliente.uc
            assert item["status"] == "sucesso"
            assert item["pdf_path"]
            assert item["ocr"] is not None
            assert item["ocr"]["mes"] == 12
            assert item["ocr"]["ano"] == 2024
            assert item["ocr"]["cliente"]["nome"]
            assert item["ocr"]["consumo"]["medidor"]

            pdf_path = Path(item["pdf_path"])
            assert pdf_path.exists()
            assert pdf_path.stat().st_size > 0
            assert pdf_path.read_bytes().startswith(b"%PDF")
    finally:
        process.terminate()
        try:
            process.wait(timeout=10)
        except subprocess.TimeoutExpired:
            process.kill()
            process.wait(timeout=10)
