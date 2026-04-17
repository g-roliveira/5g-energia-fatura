from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime
from decimal import Decimal

from pydantic import BaseModel, ConfigDict, Field, field_validator

from fatura.config import TipoAcesso


class JobTargetRequest(BaseModel):
    model_config = ConfigDict(
        json_schema_extra={"example": {"uc": "007085489032", "nome": "Paula Fernandes"}}
    )

    uc: str = Field(description="Código da unidade consumidora.")
    nome: str = Field(default="", description="Nome amigável da UC para rastreamento.")


class FaturaJobRequest(BaseModel):
    model_config = ConfigDict(
        json_schema_extra={
            "example": {
                "cpf_cnpj": "12345678901",
                "senha_portal": "minha_senha_portal",
                "uf": "BA",
                "tipo_acesso": "normal",
                "mes_ano": "122024",
                "force": False,
                "ucs": [
                    {"uc": "007085489032", "nome": "Paula Fernandes"},
                    {"uc": "001234567890", "nome": "UC Reserva"},
                ],
            }
        }
    )

    cpf_cnpj: str = Field(description="CPF ou CNPJ do titular usado no login do portal.")
    senha_portal: str = Field(description="Senha da Agência Virtual Neoenergia.")
    uf: str = Field(default="BA", description="UF da distribuidora dentro do portal.")
    tipo_acesso: TipoAcesso = Field(
        default=TipoAcesso.NORMAL,
        description="Perfil de acesso no portal. Ex.: normal, imobiliaria.",
    )
    ucs: list[JobTargetRequest] = Field(
        min_length=1,
        description="Lista de unidades consumidoras a processar dentro do mesmo login.",
    )
    mes_ano: str | None = Field(
        default=None,
        description="Competência no formato MMAAAA. Se omitido, usa a mais recente disponível.",
    )
    force: bool = Field(
        default=False,
        description="Quando true, reprocessa mesmo que a fatura já exista no banco.",
    )

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
    total: int = Field(default=0, description="Quantidade total de UCs no job.")
    completed: int = Field(default=0, description="Quantidade já finalizada.")
    success: int = Field(default=0, description="Quantidade concluída com sucesso.")
    error: int = Field(default=0, description="Quantidade concluída com erro.")


class JobStatusResponse(BaseModel):
    job_id: str = Field(description="Identificador único do job.")
    status: str = Field(description="Estado do job: queued, running, succeeded, partial_failure, failed.")
    created_at: datetime = Field(description="Data/hora de criação do job.")
    started_at: datetime | None = Field(default=None, description="Data/hora de início do processamento.")
    finished_at: datetime | None = Field(default=None, description="Data/hora de conclusão do processamento.")
    progress_total: int = Field(description="Total de itens/UCs no job.")
    progress_done: int = Field(description="Quantidade já finalizada.")
    summary: JobSummary = Field(description="Resumo consolidado do job.")


class JobItemResponse(BaseModel):
    uc: str = Field(description="Código da unidade consumidora.")
    nome: str = Field(default="", description="Nome amigável da UC.")
    status: str = Field(description="Status do item: sucesso, erro_download, erro_login etc.")
    mensagem: str = Field(default="", description="Mensagem final do processamento do item.")
    erro_tipo: str | None = Field(default=None, description="Classe/tipo de erro persistido.")
    pdf_path: str | None = Field(default=None, description="Caminho do PDF salvo localmente.")
    screenshot_path: str | None = Field(default=None, description="Screenshot de evidência em caso de erro.")
    html_path: str | None = Field(default=None, description="HTML de evidência em caso de erro.")
    mes: int | None = Field(default=None, description="Mês da fatura obtida.")
    ano: int | None = Field(default=None, description="Ano da fatura obtida.")
    valor: str | None = Field(default=None, description="Valor textual extraído da fatura.")
    data_vencimento: str | None = Field(default=None, description="Data de vencimento formatada.")
    normalizado_valor: float | None = Field(default=None, description="Valor numérico normalizado.")
    ocr: dict | None = Field(default=None, description="Payload enriquecido derivado do PDF.")
    attempts: int = Field(default=0, description="Quantidade de tentativas executadas para a UC.")


class JobResultResponse(BaseModel):
    job_id: str = Field(description="Identificador único do job.")
    status: str = Field(description="Status consolidado do job.")
    items: list[JobItemResponse] = Field(description="Resultado detalhado por UC.")


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
    data_vencimento: str | None = None
    normalizado_valor: float | None = None
    ocr_data: dict | None = None
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
