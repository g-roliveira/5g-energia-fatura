from datetime import date, datetime
from decimal import Decimal, InvalidOperation

from pydantic import BaseModel, field_validator


def normalizar_decimal_br(valor: str | Decimal | None) -> Decimal | None:
    """Converte formato brasileiro '1.234,56' para Decimal('1234.56')."""
    if valor is None:
        return None
    if isinstance(valor, Decimal):
        return valor
    if isinstance(valor, (int, float)):
        return Decimal(str(valor))
    valor = str(valor).strip()
    if not valor:
        return None
    valor = valor.replace("R$", "").replace(" ", "").strip()
    if "," in valor:
        valor = valor.replace(".", "").replace(",", ".")
    try:
        return Decimal(valor)
    except InvalidOperation:
        return None


class Cliente(BaseModel):
    codigo: str = ""
    cpf: str | None = None
    cnpj: str | None = None
    nome: str = ""
    classificacao: str | None = None
    tensao_nominal: str | None = None
    endereco: str | None = None


class Consumo(BaseModel):
    medidor: str = ""
    constante: str | None = None
    leitura_anterior: str | None = None
    leitura_atual: str | None = None


class HistoricoConsumo(BaseModel):
    periodo: str
    kwh: Decimal

    @field_validator("kwh", mode="before")
    @classmethod
    def parse_kwh(cls, v: str | Decimal) -> Decimal:
        result = normalizar_decimal_br(v)
        return result if result is not None else Decimal("0")


class ComposicaoFornecimento(BaseModel):
    energia: Decimal = Decimal("0")
    encargos: Decimal = Decimal("0")
    distribuicao: Decimal = Decimal("0")
    tributos: Decimal = Decimal("0")
    transmissao: Decimal = Decimal("0")
    perdas: Decimal = Decimal("0")

    @field_validator("*", mode="before")
    @classmethod
    def parse_valores(cls, v: str | Decimal) -> Decimal:
        result = normalizar_decimal_br(v)
        return result if result is not None else Decimal("0")


class ItemFatura(BaseModel):
    codigo: str = ""
    descricao: str = ""
    quantidade: Decimal | None = None
    quantidade_residual: Decimal | None = None
    quantidade_faturada: Decimal | None = None
    tarifa: Decimal | None = None
    valor: Decimal | None = None
    base_icms: Decimal | None = None
    aliq_icms: str | None = None
    icms: Decimal | None = None
    valor_total: Decimal | None = None

    @field_validator(
        "quantidade",
        "quantidade_residual",
        "quantidade_faturada",
        "tarifa",
        "valor",
        "base_icms",
        "icms",
        "valor_total",
        mode="before",
    )
    @classmethod
    def parse_decimal(cls, v: str | Decimal | None) -> Decimal | None:
        return normalizar_decimal_br(v)


class NotaFiscal(BaseModel):
    numero_serie: str | None = None
    apresentacao_data: str | None = None


class ContaDistribuidora(BaseModel):
    """Fatura completa da distribuidora parseada."""

    uc: str
    mes: int
    ano: int
    valor: Decimal
    vencimento: date
    numero_dias: int | None = None
    codigo_barras: str | None = None
    nota_fiscal: NotaFiscal | None = None
    cliente: Cliente = Cliente()
    consumo: Consumo | None = None
    historico_energia: list[HistoricoConsumo] = []
    composicao: ComposicaoFornecimento | None = None
    itens_fatura: list[ItemFatura] = []
    pdf_path: str | None = None
    parsed_at: datetime | None = None

    @field_validator("valor", mode="before")
    @classmethod
    def parse_valor(cls, v: str | Decimal) -> Decimal:
        result = normalizar_decimal_br(v)
        return result if result is not None else Decimal("0")
