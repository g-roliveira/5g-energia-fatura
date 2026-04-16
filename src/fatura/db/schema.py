from datetime import date, datetime

from sqlalchemy import ForeignKey, String, UniqueConstraint
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column, relationship


class Base(DeclarativeBase):
    pass


class ClienteDB(Base):
    __tablename__ = "clientes"

    id: Mapped[int] = mapped_column(primary_key=True)
    codigo: Mapped[str] = mapped_column(String(50), default="")
    cpf: Mapped[str | None] = mapped_column(String(20))
    cnpj: Mapped[str | None] = mapped_column(String(25))
    nome: Mapped[str] = mapped_column(String(200))
    classificacao: Mapped[str | None] = mapped_column(String(100))
    tensao_nominal: Mapped[str | None] = mapped_column(String(50))
    endereco: Mapped[str | None] = mapped_column(String(500))
    created_at: Mapped[datetime] = mapped_column(default=datetime.now)
    updated_at: Mapped[datetime] = mapped_column(default=datetime.now, onupdate=datetime.now)

    contas: Mapped[list["ContaDB"]] = relationship(back_populates="cliente")


class ContaDB(Base):
    __tablename__ = "contas"

    id: Mapped[int] = mapped_column(primary_key=True)
    uc: Mapped[str] = mapped_column(String(20), index=True)
    mes: Mapped[int]
    ano: Mapped[int]
    valor: Mapped[str] = mapped_column(String(20))
    vencimento: Mapped[date]
    numero_dias: Mapped[int | None]
    codigo_barras: Mapped[str | None] = mapped_column(String(60))
    pdf_path: Mapped[str | None] = mapped_column(String(500))
    parsed_at: Mapped[datetime | None]
    cliente_id: Mapped[int | None] = mapped_column(ForeignKey("clientes.id"))
    composicao_json: Mapped[str | None]
    consumo_json: Mapped[str | None]
    energia_json: Mapped[str | None]
    nota_fiscal_json: Mapped[str | None]
    created_at: Mapped[datetime] = mapped_column(default=datetime.now)

    cliente: Mapped[ClienteDB | None] = relationship(back_populates="contas")
    itens: Mapped[list["ItemFaturaDB"]] = relationship(
        back_populates="conta", cascade="all, delete-orphan"
    )

    __table_args__ = (
        UniqueConstraint("uc", "mes", "ano", name="uq_conta_uc_mes_ano"),
    )


class ItemFaturaDB(Base):
    __tablename__ = "itens_fatura"

    id: Mapped[int] = mapped_column(primary_key=True)
    conta_id: Mapped[int] = mapped_column(ForeignKey("contas.id", ondelete="CASCADE"))
    codigo: Mapped[str] = mapped_column(String(20), default="")
    descricao: Mapped[str] = mapped_column(String(300), default="")
    quantidade: Mapped[str | None] = mapped_column(String(20))
    tarifa: Mapped[str | None] = mapped_column(String(20))
    valor: Mapped[str | None] = mapped_column(String(20))
    base_icms: Mapped[str | None] = mapped_column(String(20))
    aliq_icms: Mapped[str | None] = mapped_column(String(20))
    icms: Mapped[str | None] = mapped_column(String(20))
    valor_total: Mapped[str | None] = mapped_column(String(20))

    conta: Mapped[ContaDB] = relationship(back_populates="itens")


class ProcessamentoLogDB(Base):
    __tablename__ = "processamento_log"

    id: Mapped[int] = mapped_column(primary_key=True)
    uc: Mapped[str] = mapped_column(String(20), index=True)
    mes: Mapped[int]
    ano: Mapped[int]
    status: Mapped[str] = mapped_column(String(30))
    mensagem: Mapped[str | None] = mapped_column(String(1000))
    tentativa: Mapped[int] = mapped_column(default=1)
    created_at: Mapped[datetime] = mapped_column(default=datetime.now)


class JobDB(Base):
    __tablename__ = "jobs"

    id: Mapped[str] = mapped_column(String(36), primary_key=True)
    kind: Mapped[str] = mapped_column(String(50), default="neoenergia_fatura")
    status: Mapped[str] = mapped_column(String(30), index=True)
    request_json: Mapped[str]
    created_at: Mapped[datetime] = mapped_column(default=datetime.now)
    started_at: Mapped[datetime | None]
    finished_at: Mapped[datetime | None]
    total_items: Mapped[int] = mapped_column(default=0)
    completed_items: Mapped[int] = mapped_column(default=0)
    success_items: Mapped[int] = mapped_column(default=0)
    error_items: Mapped[int] = mapped_column(default=0)
    updated_at: Mapped[datetime] = mapped_column(default=datetime.now, onupdate=datetime.now)

    items: Mapped[list["JobItemDB"]] = relationship(
        back_populates="job", cascade="all, delete-orphan"
    )


class JobItemDB(Base):
    __tablename__ = "job_items"

    id: Mapped[int] = mapped_column(primary_key=True)
    job_id: Mapped[str] = mapped_column(ForeignKey("jobs.id", ondelete="CASCADE"), index=True)
    uc: Mapped[str] = mapped_column(String(20), index=True)
    nome: Mapped[str] = mapped_column(String(200), default="")
    status: Mapped[str] = mapped_column(String(30), default="queued", index=True)
    mensagem: Mapped[str | None] = mapped_column(String(1000))
    error_type: Mapped[str | None] = mapped_column(String(100))
    pdf_path: Mapped[str | None] = mapped_column(String(500))
    screenshot_path: Mapped[str | None] = mapped_column(String(500))
    html_path: Mapped[str | None] = mapped_column(String(500))
    step_name: Mapped[str | None] = mapped_column(String(100))
    mes: Mapped[int | None]
    ano: Mapped[int | None]
    valor: Mapped[str | None] = mapped_column(String(50))
    conta_id: Mapped[int | None] = mapped_column(ForeignKey("contas.id"))
    attempts: Mapped[int] = mapped_column(default=0)
    started_at: Mapped[datetime | None]
    finished_at: Mapped[datetime | None]
    result_json: Mapped[str | None]
    created_at: Mapped[datetime] = mapped_column(default=datetime.now)
    updated_at: Mapped[datetime] = mapped_column(default=datetime.now, onupdate=datetime.now)

    job: Mapped[JobDB] = relationship(back_populates="items")

    __table_args__ = (
        UniqueConstraint("job_id", "uc", name="uq_job_item_job_uc"),
    )
