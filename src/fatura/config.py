import os
from enum import Enum
from pathlib import Path

import yaml
from pydantic import BaseModel, field_validator

from fatura.exceptions import ConfigError


class TipoAcesso(str, Enum):
    NORMAL = "normal"
    IMOBILIARIA = "imobiliaria"


class ClienteConfig(BaseModel):
    nome: str
    uc: str
    cpf_cnpj: str
    senha_portal: str
    uf: str = "BA"
    tipo_acesso: TipoAcesso = TipoAcesso.NORMAL
    ativo: bool = True

    @field_validator("cpf_cnpj")
    @classmethod
    def validar_cpf_cnpj(cls, v: str) -> str:
        digitos = "".join(c for c in v if c.isdigit())
        if len(digitos) not in (11, 14):
            raise ValueError(f"CPF/CNPJ inválido (esperado 11 ou 14 dígitos): {v}")
        return v

    @property
    def is_cnpj(self) -> bool:
        digitos = "".join(c for c in self.cpf_cnpj if c.isdigit())
        return len(digitos) == 14


class PortalConfig(BaseModel):
    url_base: str = "https://agenciavirtual.neoenergia.com"
    headless: bool = True
    browser_channel: str | None = "chrome"
    browser_executable_path: str | None = None
    timeout_ms: int = 60_000
    max_retries: int = 3
    delay_between_clients_s: float = 5.0
    download_dir: str = "./downloads"


class DatabaseConfig(BaseModel):
    url: str = "postgresql+psycopg://backoffice:backoffice@127.0.0.1:5432/backoffice"


class ParserConfig(BaseModel):
    engine: str = "pymupdf"
    enable_mistral_fallback: bool = True
    validate_new_pdfs_with_mistral: bool = False
    mistral_api_key: str = ""
    mistral_model: str = "mistral-ocr-latest"
    mistral_timeout_ms: int = 45_000


class ServiceConfig(BaseModel):
    host: str = "127.0.0.1"
    port: int = 8000
    api_key: str = ""
    max_concurrent_jobs: int = 1
    artifacts_dir: str = "./downloads/_artifacts"
    reset_incomplete_jobs_on_startup: bool = True


class AppConfig(BaseModel):
    portal: PortalConfig = PortalConfig()
    database: DatabaseConfig = DatabaseConfig()
    parser: ParserConfig = ParserConfig()
    service: ServiceConfig = ServiceConfig()
    clientes: list[ClienteConfig] = []


def _load_dotenv_if_present(path: str | Path = ".env") -> None:
    dotenv_path = Path(path)
    if not dotenv_path.exists():
        return
    for raw_line in dotenv_path.read_text(encoding="utf-8").splitlines():
        line = raw_line.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        key, value = line.split("=", 1)
        key = key.strip()
        value = value.strip().strip('"').strip("'")
        os.environ.setdefault(key, value)


def load_config(path: str | Path = "config.yaml") -> AppConfig:
    """Lê o arquivo YAML de configuração e retorna AppConfig validado."""
    _load_dotenv_if_present()
    config_path = Path(path)
    if not config_path.exists():
        raise ConfigError(f"Arquivo de configuração não encontrado: {config_path}")

    try:
        with open(config_path) as f:
            raw = yaml.safe_load(f)
    except yaml.YAMLError as e:
        raise ConfigError(f"Erro ao parsear YAML: {e}") from e

    if not raw or not isinstance(raw, dict):
        raise ConfigError("Arquivo de configuração vazio ou inválido")

    try:
        config = AppConfig(**raw)
    except Exception as e:
        raise ConfigError(f"Erro de validação na configuração: {e}") from e

    if not config.service.api_key:
        config.service.api_key = os.getenv("FATURA_API_KEY", "")
    if not config.parser.mistral_api_key:
        config.parser.mistral_api_key = os.getenv("MISTRAL_API_KEY", "")

    return config
