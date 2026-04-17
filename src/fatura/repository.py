import json
import uuid
from abc import ABC, abstractmethod
from datetime import datetime

from sqlalchemy import create_engine, inspect, select, text
from sqlalchemy.orm import Session, sessionmaker

from fatura.db.schema import (
    Base,
    ClienteDB,
    ContaDB,
    ItemFaturaDB,
    JobDB,
    JobItemDB,
    ProcessamentoLogDB,
)
from fatura.exceptions import RepositoryError
from fatura.models import ContaDistribuidora
from fatura.models import conta_para_ocr_payload
from fatura.service_models import (
    BatchItemResult,
    FaturaJobRequest,
    JobItemResponse,
    JobResultResponse,
    JobStatusResponse,
    JobSummary,
)


class FaturaRepository(ABC):
    @abstractmethod
    def salvar_conta(self, conta: ContaDistribuidora) -> int: ...

    @abstractmethod
    def buscar_conta(self, uc: str, mes: int, ano: int) -> ContaDistribuidora | None: ...

    @abstractmethod
    def listar_contas(
        self, uc: str | None = None, ano: int | None = None
    ) -> list[ContaDistribuidora]: ...

    @abstractmethod
    def conta_existe(self, uc: str, mes: int, ano: int) -> bool: ...

    @abstractmethod
    def registrar_log(
        self, uc: str, mes: int, ano: int, status: str, mensagem: str | None = None
    ) -> None: ...


class SqliteFaturaRepository(FaturaRepository):
    def __init__(self, db_url: str = "sqlite:///data/faturas.db"):
        self._engine = create_engine(db_url, echo=False)
        Base.metadata.create_all(self._engine)
        self._ensure_runtime_columns()
        self._Session = sessionmaker(bind=self._engine)

    def _get_session(self) -> Session:
        return self._Session()

    def _ensure_runtime_columns(self) -> None:
        inspector = inspect(self._engine)
        conta_columns = {column["name"] for column in inspector.get_columns("contas")}
        with self._engine.begin() as connection:
            if "ocr_json" not in conta_columns:
                connection.execute(text("ALTER TABLE contas ADD COLUMN ocr_json TEXT"))

    def salvar_conta(self, conta: ContaDistribuidora) -> int:
        try:
            with self._get_session() as session:
                cliente_db = self._upsert_cliente(session, conta)
                conta_db = ContaDB(
                    uc=conta.uc,
                    mes=conta.mes,
                    ano=conta.ano,
                    valor=str(conta.valor),
                    vencimento=conta.vencimento,
                    numero_dias=conta.numero_dias,
                    codigo_barras=conta.codigo_barras,
                    pdf_path=conta.pdf_path,
                    parsed_at=conta.parsed_at or datetime.now(),
                    cliente_id=cliente_db.id if cliente_db else None,
                    composicao_json=(
                        conta.composicao.model_dump_json() if conta.composicao else None
                    ),
                    consumo_json=conta.consumo.model_dump_json() if conta.consumo else None,
                    energia_json=(
                        json.dumps([h.model_dump() for h in conta.historico_energia], default=str)
                        if conta.historico_energia
                        else None
                    ),
                    nota_fiscal_json=(
                        conta.nota_fiscal.model_dump_json() if conta.nota_fiscal else None
                    ),
                    ocr_json=json.dumps(conta_para_ocr_payload(conta), ensure_ascii=False),
                )

                for item in conta.itens_fatura:
                    conta_db.itens.append(
                        ItemFaturaDB(
                            codigo=item.codigo,
                            descricao=item.descricao,
                            quantidade=str(item.quantidade) if item.quantidade else None,
                            tarifa=str(item.tarifa) if item.tarifa else None,
                            valor=str(item.valor) if item.valor else None,
                            base_icms=str(item.base_icms) if item.base_icms else None,
                            aliq_icms=item.aliq_icms,
                            icms=str(item.icms) if item.icms else None,
                            valor_total=str(item.valor_total) if item.valor_total else None,
                        )
                    )

                session.add(conta_db)
                session.commit()
                return conta_db.id
        except Exception as e:
            raise RepositoryError(f"Erro ao salvar conta: {e}") from e

    def _upsert_cliente(self, session: Session, conta: ContaDistribuidora) -> ClienteDB | None:
        if not conta.cliente.nome:
            return None

        stmt = select(ClienteDB).where(ClienteDB.nome == conta.cliente.nome)
        existing = session.execute(stmt).scalar_one_or_none()

        if existing:
            existing.cpf = conta.cliente.cpf
            existing.cnpj = conta.cliente.cnpj
            existing.classificacao = conta.cliente.classificacao
            existing.updated_at = datetime.now()
            return existing

        cliente_db = ClienteDB(
            codigo=conta.cliente.codigo,
            cpf=conta.cliente.cpf,
            cnpj=conta.cliente.cnpj,
            nome=conta.cliente.nome,
            classificacao=conta.cliente.classificacao,
            tensao_nominal=conta.cliente.tensao_nominal,
            endereco=conta.cliente.endereco,
        )
        session.add(cliente_db)
        session.flush()
        return cliente_db

    def buscar_conta(self, uc: str, mes: int, ano: int) -> ContaDistribuidora | None:
        with self._get_session() as session:
            stmt = select(ContaDB).where(
                ContaDB.uc == uc, ContaDB.mes == mes, ContaDB.ano == ano
            )
            conta_db = session.execute(stmt).scalar_one_or_none()
            if not conta_db:
                return None
            return self._db_to_model(conta_db)

    def listar_contas(
        self, uc: str | None = None, ano: int | None = None
    ) -> list[ContaDistribuidora]:
        with self._get_session() as session:
            stmt = select(ContaDB)
            if uc:
                stmt = stmt.where(ContaDB.uc == uc)
            if ano:
                stmt = stmt.where(ContaDB.ano == ano)
            stmt = stmt.order_by(ContaDB.ano.desc(), ContaDB.mes.desc())
            rows = session.execute(stmt).scalars().all()
            return [self._db_to_model(r) for r in rows]

    def conta_existe(self, uc: str, mes: int, ano: int) -> bool:
        with self._get_session() as session:
            stmt = select(ContaDB.id).where(
                ContaDB.uc == uc, ContaDB.mes == mes, ContaDB.ano == ano
            )
            return session.execute(stmt).scalar_one_or_none() is not None

    def registrar_log(
        self, uc: str, mes: int, ano: int, status: str, mensagem: str | None = None
    ) -> None:
        with self._get_session() as session:
            log = ProcessamentoLogDB(
                uc=uc, mes=mes, ano=ano, status=status, mensagem=mensagem
            )
            session.add(log)
            session.commit()

    def criar_job(self, request: FaturaJobRequest) -> str:
        job_id = str(uuid.uuid4())
        request_public = request.model_dump()
        request_public["cpf_cnpj"] = "***"
        request_public["senha_portal"] = "***"
        try:
            with self._get_session() as session:
                job = JobDB(
                    id=job_id,
                    status="queued",
                    request_json=json.dumps(request_public, ensure_ascii=False),
                    total_items=len(request.ucs),
                )
                session.add(job)
                session.flush()

                for target in request.ucs:
                    session.add(
                        JobItemDB(
                            job_id=job_id,
                            uc=target.uc,
                            nome=target.nome,
                            status="queued",
                        )
                    )

                session.commit()
                return job_id
        except Exception as e:
            raise RepositoryError(f"Erro ao criar job: {e}") from e

    def carregar_job_request(self, job_id: str) -> FaturaJobRequest:
        with self._get_session() as session:
            job = session.get(JobDB, job_id)
            if not job:
                raise RepositoryError(f"Job não encontrado: {job_id}")
            try:
                request_data = json.loads(job.request_json)
            except json.JSONDecodeError as e:
                raise RepositoryError(f"request_json inválido para job {job_id}: {e}") from e

            if request_data.get("cpf_cnpj") in (None, "", "***") or request_data.get(
                "senha_portal"
            ) in (None, "", "***"):
                raise RepositoryError(
                    "Credenciais do job não estão disponíveis no banco. "
                    "O job precisa ser executado com as credenciais mantidas em memória pelo runtime."
                )

            return FaturaJobRequest.model_validate(request_data)

    def marcar_job_em_execucao(self, job_id: str) -> None:
        with self._get_session() as session:
            job = session.get(JobDB, job_id)
            if not job:
                raise RepositoryError(f"Job não encontrado: {job_id}")
            job.status = "running"
            job.started_at = job.started_at or datetime.now()
            session.commit()

    def finalizar_job(self, job_id: str, status: str) -> None:
        with self._get_session() as session:
            job = session.get(JobDB, job_id)
            if not job:
                raise RepositoryError(f"Job não encontrado: {job_id}")
            self._recompute_job_counts(session, job)
            job.status = status
            job.finished_at = datetime.now()
            session.commit()

    def falhar_itens_pendentes_do_job(self, job_id: str, mensagem: str) -> None:
        with self._get_session() as session:
            job = session.get(JobDB, job_id)
            if not job:
                raise RepositoryError(f"Job não encontrado: {job_id}")
            now = datetime.now()
            for item in job.items:
                if item.status in ("queued", "running"):
                    item.status = "error"
                    item.error_type = "WorkerUnhandledError"
                    item.mensagem = mensagem
                    item.finished_at = now
            self._recompute_job_counts(session, job)
            session.commit()

    def marcar_jobs_incompletos_como_falhos(self, mensagem: str) -> int:
        with self._get_session() as session:
            stmt = select(JobDB).where(JobDB.status.in_(("queued", "running")))
            jobs = session.execute(stmt).scalars().all()
            now = datetime.now()
            count = 0

            for job in jobs:
                job.status = "failed"
                job.finished_at = now
                for item in job.items:
                    if item.status in ("queued", "running"):
                        item.status = "error"
                        item.error_type = "ServiceRestartError"
                        item.mensagem = mensagem
                        item.finished_at = now
                self._recompute_job_counts(session, job)
                count += 1

            session.commit()
            return count

    def marcar_item_em_execucao(
        self, job_id: str, uc: str, attempts: int, step_name: str | None = None
    ) -> None:
        with self._get_session() as session:
            item = self._buscar_item(session, job_id, uc)
            item.status = "running"
            item.attempts = attempts
            item.step_name = step_name
            item.started_at = item.started_at or datetime.now()
            session.commit()

    def salvar_resultado_item(self, job_id: str, result: BatchItemResult) -> None:
        with self._get_session() as session:
            item = self._buscar_item(session, job_id, result.uc)
            item.nome = result.nome
            item.status = result.status
            item.mensagem = result.mensagem
            item.error_type = result.error_type
            item.pdf_path = result.pdf_path
            item.screenshot_path = result.screenshot_path
            item.html_path = result.html_path
            item.mes = result.mes
            item.ano = result.ano
            item.valor = str(result.valor) if result.valor is not None else None
            item.conta_id = result.conta_id
            item.attempts = result.attempts
            item.finished_at = datetime.now()
            item.result_json = json.dumps(
                {
                    "uc": result.uc,
                    "status": result.status,
                    "mensagem": result.mensagem,
                    "error_type": result.error_type,
                    "pdf_path": result.pdf_path,
                    "mes": result.mes,
                    "ano": result.ano,
                    "valor": str(result.valor) if result.valor is not None else None,
                    "data_vencimento": result.data_vencimento,
                    "normalizado_valor": result.normalizado_valor,
                    "ocr": result.ocr_data,
                },
                ensure_ascii=False,
            )

            job = session.get(JobDB, job_id)
            if not job:
                raise RepositoryError(f"Job não encontrado: {job_id}")
            self._recompute_job_counts(session, job)
            session.commit()

    def obter_status_job(self, job_id: str) -> JobStatusResponse | None:
        with self._get_session() as session:
            job = session.get(JobDB, job_id)
            if not job:
                return None
            self._recompute_job_counts(session, job)
            session.commit()
            return self._job_to_status(job)

    def listar_jobs(self, limit: int = 50) -> list[JobStatusResponse]:
        with self._get_session() as session:
            stmt = select(JobDB).order_by(JobDB.created_at.desc()).limit(limit)
            jobs = session.execute(stmt).scalars().all()
            for job in jobs:
                self._recompute_job_counts(session, job)
            session.commit()
            return [self._job_to_status(job) for job in jobs]

    def obter_resultado_job(self, job_id: str) -> JobResultResponse | None:
        with self._get_session() as session:
            job = session.get(JobDB, job_id)
            if not job:
                return None
            items = [
                JobItemResponse(
                    uc=item.uc,
                    nome=item.nome,
                    status=item.status,
                    mensagem=item.mensagem or "",
                    erro_tipo=item.error_type,
                    pdf_path=item.pdf_path,
                    screenshot_path=item.screenshot_path,
                    html_path=item.html_path,
                    mes=item.mes,
                    ano=item.ano,
                    valor=item.valor,
                    data_vencimento=(json.loads(item.result_json).get("data_vencimento") if item.result_json else None),
                    normalizado_valor=(json.loads(item.result_json).get("normalizado_valor") if item.result_json else None),
                    ocr=(json.loads(item.result_json).get("ocr") if item.result_json else None),
                    attempts=item.attempts,
                )
                for item in sorted(job.items, key=lambda value: value.id)
            ]
            return JobResultResponse(job_id=job.id, status=job.status, items=items)

    def _buscar_item(self, session: Session, job_id: str, uc: str) -> JobItemDB:
        stmt = select(JobItemDB).where(JobItemDB.job_id == job_id, JobItemDB.uc == uc)
        item = session.execute(stmt).scalar_one_or_none()
        if not item:
            raise RepositoryError(f"Item do job não encontrado: {job_id}/{uc}")
        return item

    def _recompute_job_counts(self, session: Session, job: JobDB) -> None:
        _ = session
        items = list(job.items)
        job.total_items = len(items)
        job.completed_items = sum(1 for item in items if item.status not in ("queued", "running"))
        job.success_items = sum(1 for item in items if item.status in ("sucesso", "pulado"))
        job.error_items = sum(1 for item in items if item.status.startswith("erro") or item.status == "error")

    def _job_to_status(self, job: JobDB) -> JobStatusResponse:
        return JobStatusResponse(
            job_id=job.id,
            status=job.status,
            created_at=job.created_at,
            started_at=job.started_at,
            finished_at=job.finished_at,
            progress_total=job.total_items,
            progress_done=job.completed_items,
            summary=JobSummary(
                total=job.total_items,
                completed=job.completed_items,
                success=job.success_items,
                error=job.error_items,
            ),
        )

    def _db_to_model(self, conta_db: ContaDB) -> ContaDistribuidora:
        from fatura.models import (
            Cliente,
            ComposicaoFornecimento,
            Consumo,
            HistoricoConsumo,
            ItemFatura,
            NotaFiscal,
        )

        cliente = Cliente()
        if conta_db.cliente:
            cliente = Cliente(
                codigo=conta_db.cliente.codigo,
                cpf=conta_db.cliente.cpf,
                cnpj=conta_db.cliente.cnpj,
                nome=conta_db.cliente.nome,
                classificacao=conta_db.cliente.classificacao,
                tensao_nominal=conta_db.cliente.tensao_nominal,
                endereco=conta_db.cliente.endereco,
            )

        composicao = None
        if conta_db.composicao_json:
            composicao = ComposicaoFornecimento.model_validate_json(conta_db.composicao_json)

        consumo = None
        if conta_db.consumo_json:
            consumo = Consumo.model_validate_json(conta_db.consumo_json)

        historico = []
        if conta_db.energia_json:
            historico = [HistoricoConsumo(**h) for h in json.loads(conta_db.energia_json)]

        nota_fiscal = None
        if conta_db.nota_fiscal_json:
            nota_fiscal = NotaFiscal.model_validate_json(conta_db.nota_fiscal_json)

        itens = [
            ItemFatura(
                codigo=i.codigo,
                descricao=i.descricao,
                quantidade=i.quantidade,
                tarifa=i.tarifa,
                valor=i.valor,
                base_icms=i.base_icms,
                aliq_icms=i.aliq_icms,
                icms=i.icms,
                valor_total=i.valor_total,
            )
            for i in conta_db.itens
        ]

        return ContaDistribuidora(
            uc=conta_db.uc,
            mes=conta_db.mes,
            ano=conta_db.ano,
            valor=conta_db.valor,
            vencimento=conta_db.vencimento,
            numero_dias=conta_db.numero_dias,
            codigo_barras=conta_db.codigo_barras,
            pdf_path=conta_db.pdf_path,
            parsed_at=conta_db.parsed_at,
            cliente=cliente,
            consumo=consumo,
            historico_energia=historico,
            composicao=composicao,
            itens_fatura=itens,
            nota_fiscal=nota_fiscal,
        )
