from __future__ import annotations

import asyncio

import structlog

from fatura.config import AppConfig
from fatura.jobs import executar_job_persistido
from fatura.repository import SqliteFaturaRepository
from fatura.service_models import FaturaJobRequest, JobResultResponse, JobStatusResponse

logger = structlog.get_logger()


class FaturaJobRuntime:
    def __init__(self, config: AppConfig, repo: SqliteFaturaRepository | None = None) -> None:
        self._config = config
        self._repo = repo or SqliteFaturaRepository(config.database.url)
        self._semaphore = asyncio.Semaphore(max(1, config.service.max_concurrent_jobs))
        self._request_cache: dict[str, FaturaJobRequest] = {}

    @property
    def repo(self) -> SqliteFaturaRepository:
        return self._repo

    def prepare(self) -> int:
        if self._config.service.reset_incomplete_jobs_on_startup:
            return self._repo.marcar_jobs_incompletos_como_falhos(
                "Job interrompido por reinício do serviço."
            )
        return 0

    def create_job(self, request: FaturaJobRequest) -> str:
        job_id = self._repo.criar_job(request)
        self._request_cache[job_id] = request
        return job_id

    def get_status(self, job_id: str) -> JobStatusResponse | None:
        return self._repo.obter_status_job(job_id)

    def list_jobs(self, limit: int = 50) -> list[JobStatusResponse]:
        return self._repo.listar_jobs(limit=limit)

    def get_result(self, job_id: str) -> JobResultResponse | None:
        return self._repo.obter_resultado_job(job_id)

    async def run_job(self, job_id: str) -> None:
        async with self._semaphore:
            log = logger.bind(job_id=job_id)
            log.info("job_worker_iniciado")
            request = self._request_cache.get(job_id)
            try:
                await executar_job_persistido(
                    config=self._config,
                    job_id=job_id,
                    repo=self._repo,
                    request=request,
                )
            except Exception as exc:  # pragma: no cover - proteção final do worker
                log.exception("job_worker_falhou", erro=str(exc))
                self._repo.falhar_itens_pendentes_do_job(job_id, str(exc))
                self._repo.finalizar_job(job_id, "failed")
            else:
                log.info("job_worker_concluido")
            finally:
                self._request_cache.pop(job_id, None)
