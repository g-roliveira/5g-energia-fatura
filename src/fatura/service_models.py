from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime
from decimal import Decimal

from pydantic import BaseModel, Field, field_validator

from fatura.config import TipoAcesso


class JobTargetRequest(BaseModel):
    uc: str
    nome: str = ""


class FaturaJobRequest(BaseModel):
    cpf_cnpj: str
    senha_portal: str
    uf: str = "BA"
    tipo_acesso: TipoAcesso = TipoAcesso.NORMAL
    ucs: list[JobTargetRequest] = Field(min_length=1)
    mes_ano: str | None = None
    force: bool = False

    @field_validator("cpf_cnpj")
    @classmethod
    def validar_cpf_cnpj(cls, v: str) -> str:
        digitos = "".join(c for c in v if c.isdigit())
        if len(digitos) not in (11, 14):
            raise ValueError("CPF/CNPJ inválido")
        return digitos

    @field_validator("mes_ano")
    @classmethod
    def validar_mes_ano(cls, v: str | None) -> str | None:
        if v is None:
            return None
        if len(v) != 6 or not v.isdigit():
            raise ValueError("mes_ano deve usar o formato MMAAAA")
        return v


class JobSummary(BaseModel):
    total: int = 0
    completed: int = 0
    success: int = 0
    error: int = 0


class JobStatusResponse(BaseModel):
    job_id: str
    status: str
    created_at: datetime
    started_at: datetime | None = None
    finished_at: datetime | None = None
    progress_total: int
    progress_done: int
    summary: JobSummary


class JobItemResponse(BaseModel):
    uc: str
    nome: str = ""
    status: str
    mensagem: str = ""
    erro_tipo: str | None = None
    pdf_path: str | None = None
    screenshot_path: str | None = None
    html_path: str | None = None
    mes: int | None = None
    ano: int | None = None
    valor: str | None = None
    attempts: int = 0


class JobResultResponse(BaseModel):
    job_id: str
    status: str
    items: list[JobItemResponse]


@dataclass(slots=True)
class BatchTarget:
    uc: str
    nome: str = ""


@dataclass(slots=True)
class BatchSpec:
    cpf_cnpj: str
    senha_portal: str
    uf: str
    tipo_acesso: TipoAcesso
    targets: list[BatchTarget]
    mes_ano: str | None = None
    force: bool = False


@dataclass(slots=True)
class BatchItemResult:
    uc: str
    nome: str
    status: str
    mensagem: str = ""
    error_type: str | None = None
    pdf_path: str | None = None
    screenshot_path: str | None = None
    html_path: str | None = None
    mes: int | None = None
    ano: int | None = None
    valor: Decimal | None = None
    conta_id: int | None = None
    attempts: int = 0


@dataclass(slots=True)
class BatchRunResult:
    status: str
    total: int = 0
    completed: int = 0
    success: int = 0
    error: int = 0
    items: list[BatchItemResult] = field(default_factory=list)
