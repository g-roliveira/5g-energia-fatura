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
    limites_tensao: str | None = None
    endereco: str | None = None


class Consumo(BaseModel):
    medidor: str = ""
    constante: str | None = None
    leitura_anterior: str | None = None
    leitura_atual: str | None = None
    leitura_anterior_data: str | None = None
    leitura_data: str | None = None
    leitura_proxima_data: str | None = None


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


class SceeSummary(BaseModel):
    """Resumo SCEE/MMGD extraído do rodapé da fatura Coelba."""

    excedente_kwh: Decimal | None = None
    creditos_utilizados_kwh: Decimal | None = None
    saldo_proximo_ciclo_kwh: Decimal | None = None
    energia_injetada_kwh: Decimal | None = None
    saldo_credito_kwh: Decimal | None = None
    texto_original: str = ""
    confianca: float = 1.0

    @field_validator(
        "excedente_kwh",
        "creditos_utilizados_kwh",
        "saldo_proximo_ciclo_kwh",
        "energia_injetada_kwh",
        "saldo_credito_kwh",
        mode="before",
    )
    @classmethod
    def parse_kwh(cls, v: str | Decimal | None) -> Decimal | None:
        return normalizar_decimal_br(v)


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
    emissao_data: str | None = None
    controle_n: str | None = None
    aviso: str | None = None
    informacoes_gerais: str | None = None
    cliente: Cliente = Cliente()
    consumo: Consumo | None = None
    historico_energia: list[HistoricoConsumo] = []
    composicao: ComposicaoFornecimento | None = None
    itens_fatura: list[ItemFatura] = []
    scee_summary: SceeSummary | None = None
    pdf_path: str | None = None
    parsed_at: datetime | None = None

    @field_validator("valor", mode="before")
    @classmethod
    def parse_valor(cls, v: str | Decimal) -> Decimal:
        result = normalizar_decimal_br(v)
        return result if result is not None else Decimal("0")


def formatar_decimal_br(valor: Decimal | None) -> str:
    if valor is None:
        return ""
    quantizado = valor.quantize(Decimal("0.01"))
    inteiro, frac = f"{quantizado:.2f}".split(".")
    inteiro = f"{int(inteiro):,}".replace(",", ".")
    return f"{inteiro},{frac}"


def formatar_data_br(valor: date | str | None) -> str:
    if valor is None:
        return ""
    if isinstance(valor, date):
        return valor.strftime("%d/%m/%Y")
    return str(valor)


def conta_para_ocr_payload(conta: ContaDistribuidora) -> dict:
    composicao = conta.composicao
    return {
        "mes": conta.mes,
        "ano": conta.ano,
        "valor": f"R$ {formatar_decimal_br(conta.valor)}",
        "normalizado_valor": float(conta.valor),
        "vencimento": formatar_data_br(conta.vencimento),
        "leitura_anterior_data": conta.consumo.leitura_anterior_data if conta.consumo else "",
        "leitura_data": conta.consumo.leitura_data if conta.consumo else "",
        "leitura_proxima_data": conta.consumo.leitura_proxima_data if conta.consumo else "",
        "emissao_data": conta.emissao_data or "",
        "controle_n": conta.controle_n or "",
        "numero_dias": conta.numero_dias,
        "codigo_barras": conta.codigo_barras or "",
        "aviso": conta.aviso or "",
        "nota_fiscal": {
            "numero_serie": conta.nota_fiscal.numero_serie if conta.nota_fiscal else "",
            "apresentacao_data": conta.nota_fiscal.apresentacao_data if conta.nota_fiscal else "",
        },
        "cliente": {
            "codigo": conta.cliente.codigo or "",
            "cpf": conta.cliente.cpf or "",
            "cnpj": conta.cliente.cnpj or "",
            "nome": conta.cliente.nome or "",
            "classificacao": conta.cliente.classificacao or "",
            "tensao_nominal": conta.cliente.tensao_nominal or "",
            "limites_tensao": conta.cliente.limites_tensao or "",
            "endereco": conta.cliente.endereco or "",
        },
        "consumo": {
            "medidor": conta.consumo.medidor if conta.consumo else "",
            "constante": conta.consumo.constante if conta.consumo else "",
            "leitura_anterior": conta.consumo.leitura_anterior if conta.consumo else "",
            "leitura_atual": conta.consumo.leitura_atual if conta.consumo else "",
        },
        "energia": {
            "historico_consumo": [
                {"periodo": item.periodo, "kwh": str(item.kwh)} for item in conta.historico_energia
            ]
        },
        "composicao_fornecimento": {
            "energia": f"R$ {formatar_decimal_br(composicao.energia)}" if composicao else "",
            "normalizado_energia": float(composicao.energia) if composicao else None,
            "encargos": f"R$ {formatar_decimal_br(composicao.encargos)}" if composicao else "",
            "normalizado_encargos": float(composicao.encargos) if composicao else None,
            "distribuicao": f"R$ {formatar_decimal_br(composicao.distribuicao)}" if composicao else "",
            "normalizado_distribuicao": float(composicao.distribuicao) if composicao else None,
            "tributos": f"R$ {formatar_decimal_br(composicao.tributos)}" if composicao else "",
            "normalizado_tributos": float(composicao.tributos) if composicao else None,
            "transmissao": f"R$ {formatar_decimal_br(composicao.transmissao)}" if composicao else "",
            "normalizado_transmissao": float(composicao.transmissao) if composicao else None,
            "perdas": f"R$ {formatar_decimal_br(composicao.perdas)}" if composicao else "",
            "normalizado_perdas": float(composicao.perdas) if composicao else None,
        },
        "informacoes_gerais": conta.informacoes_gerais or "",
        "scee_summary": {
            "excedente_kwh": formatar_decimal_br(conta.scee_summary.excedente_kwh) if conta.scee_summary else "",
            "creditados_utilizados_kwh": formatar_decimal_br(conta.scee_summary.creditos_utilizados_kwh) if conta.scee_summary else "",
            "saldo_proximo_ciclo_kwh": formatar_decimal_br(conta.scee_summary.saldo_proximo_ciclo_kwh) if conta.scee_summary else "",
            "energia_injetada_kwh": formatar_decimal_br(conta.scee_summary.energia_injetada_kwh) if conta.scee_summary else "",
            "saldo_credito_kwh": formatar_decimal_br(conta.scee_summary.saldo_credito_kwh) if conta.scee_summary else "",
        },
        "itens_fatura": [
            {
                "codigo": item.codigo or "",
                "descricao": item.descricao or "",
                "quantidade": formatar_decimal_br(item.quantidade) if item.quantidade is not None else "",
                "quantidade_residual": formatar_decimal_br(item.quantidade_residual) if item.quantidade_residual is not None else "",
                "quantidade_faturada": formatar_decimal_br(item.quantidade_faturada) if item.quantidade_faturada is not None else "",
                "tarifa": formatar_decimal_br(item.tarifa) if item.tarifa is not None else "",
                "valor": formatar_decimal_br(item.valor) if item.valor is not None else "",
                "base_icms": formatar_decimal_br(item.base_icms) if item.base_icms is not None else "",
                "aliq_icms": item.aliq_icms or "",
                "icms": formatar_decimal_br(item.icms) if item.icms is not None else "",
                "valor_total": formatar_decimal_br(item.valor_total) if item.valor_total is not None else "",
            }
            for item in conta.itens_fatura
        ],
    }
