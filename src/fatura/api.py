from __future__ import annotations

import os
from contextlib import asynccontextmanager

import uvicorn
from fastapi import BackgroundTasks, Depends, FastAPI, Header, HTTPException, Query, status
from pydantic import BaseModel

from fatura.config import AppConfig, load_config
from fatura.logging_config import setup_logging
from fatura.runtime import FaturaJobRuntime
from fatura.service_models import FaturaJobRequest, JobResultResponse, JobStatusResponse


class HealthResponse(BaseModel):
    status: str
    reset_jobs: int = 0


def _load_app_config(config_path: str | None = None, config: AppConfig | None = None) -> AppConfig:
    if config is not None:
        return config
    return load_config(config_path or os.getenv("FATURA_CONFIG", "config.yaml"))


def _build_auth_dependency(app: FastAPI):
    async def require_api_key(x_api_key: str | None = Header(default=None)) -> None:
        expected = app.state.config.service.api_key
        if not expected:
            return
        if x_api_key != expected:
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="X-API-Key inválido",
            )

    return require_api_key


def create_app(config_path: str | None = None, config: AppConfig | None = None) -> FastAPI:
    setup_logging(1)
    app_config = _load_app_config(config_path=config_path, config=config)
    runtime = FaturaJobRuntime(app_config)

    @asynccontextmanager
    async def lifespan(app: FastAPI):
        app.state.config = app_config
        app.state.runtime = runtime
        app.state.reset_jobs = runtime.prepare()
        yield

    app = FastAPI(
        title="5G Energia Fatura API",
        version="0.1.0",
        lifespan=lifespan,
    )
    app.state.config = app_config
    app.state.runtime = runtime
    app.state.reset_jobs = 0

    require_api_key = _build_auth_dependency(app)

    @app.get("/health", response_model=HealthResponse)
    async def health() -> HealthResponse:
        return HealthResponse(status="ok", reset_jobs=app.state.reset_jobs)

    @app.post("/jobs/faturas", response_model=JobStatusResponse, dependencies=[Depends(require_api_key)])
    async def create_fatura_job(
        request: FaturaJobRequest,
        background_tasks: BackgroundTasks,
    ) -> JobStatusResponse:
        job_id = app.state.runtime.create_job(request)
        background_tasks.add_task(app.state.runtime.run_job, job_id)
        status_response = app.state.runtime.get_status(job_id)
        assert status_response is not None
        return status_response

    @app.get("/jobs", response_model=list[JobStatusResponse], dependencies=[Depends(require_api_key)])
    async def list_jobs(
        limit: int = Query(default=50, ge=1, le=200),
    ) -> list[JobStatusResponse]:
        return app.state.runtime.list_jobs(limit=limit)

    @app.get("/jobs/{job_id}", response_model=JobStatusResponse, dependencies=[Depends(require_api_key)])
    async def get_job(job_id: str) -> JobStatusResponse:
        job = app.state.runtime.get_status(job_id)
        if not job:
            raise HTTPException(status_code=404, detail="Job não encontrado")
        return job

    @app.get(
        "/jobs/{job_id}/result",
        response_model=JobResultResponse,
        dependencies=[Depends(require_api_key)],
    )
    async def get_job_result(job_id: str) -> JobResultResponse:
        result = app.state.runtime.get_result(job_id)
        if not result:
            raise HTTPException(status_code=404, detail="Job não encontrado")
        return result

    return app


def main() -> None:
    config = _load_app_config()
    app = create_app(config=config)
    uvicorn.run(
        app,
        host=config.service.host,
        port=config.service.port,
    )

