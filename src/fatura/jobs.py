from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime
from itertools import groupby

import structlog

from fatura.coelba_client import CoelbaClient
from fatura.config import AppConfig, ClienteConfig
from fatura.exceptions import (
    CaptchaError,
    DownloadError,
    FaturaError,
    LayoutChangedError,
    LoginError,
    ParserError,
    SessionExpiredError,
)
from fatura.parser_pdf import CoelbaPdfParser
from fatura.repository import SqliteFaturaRepository
from fatura.models import conta_para_ocr_payload, formatar_data_br
from fatura.service_models import (
    BatchItemResult,
    BatchRunResult,
    BatchSpec,
    BatchTarget,
    FaturaJobRequest,
)

logger = structlog.get_logger()


@dataclass
class DetalheProcessamento:
    uc: str
    nome: str
    status: str
    mensagem: str = ""


@dataclass
class ProcessamentoResult:
    total: int = 0
    sucesso: int = 0
    erro: int = 0
    pulado: int = 0
    detalhes: list[DetalheProcessamento] = field(default_factory=list)


def _somente_digitos(valor: str | None) -> str:
    return "".join(ch for ch in (valor or "") if ch.isdigit())


def _ucs_equivalentes(uc_a: str | None, uc_b: str | None) -> bool:
    digitos_a = _somente_digitos(uc_a).lstrip("0") or "0"
    digitos_b = _somente_digitos(uc_b).lstrip("0") or "0"
    return digitos_a == digitos_b


def _resolver_mes_ano(mes_ano: str | None) -> tuple[int | None, int | None, str | None]:
    if mes_ano:
        if len(mes_ano) != 6 or not mes_ano.isdigit():
            raise ValueError(f"Formato inválido para mes_ano: '{mes_ano}'. Use MMAAAA (ex: 032026)")
        return int(mes_ano[:2]), int(mes_ano[2:]), mes_ano
    return None, None, None


def _error_type(exc: Exception) -> str:
    return exc.__class__.__name__


def _final_job_status(result: BatchRunResult) -> str:
    if result.error == 0:
        return "succeeded"
    if result.success > 0:
        return "partial_failure"
    return "failed"


class BatchProcessor:
    def __init__(
        self,
        config: AppConfig,
        repo: SqliteFaturaRepository | None = None,
        parser: CoelbaPdfParser | None = None,
        client_factory: type[CoelbaClient] = CoelbaClient,
    ) -> None:
        self._config = config
        self._repo = repo or SqliteFaturaRepository(config.database.url)
        self._parser = parser or CoelbaPdfParser(config=config.parser)
        self._client_factory = client_factory

    async def run_batch(self, spec: BatchSpec, job_id: str | None = None) -> BatchRunResult:
        result = BatchRunResult(status="running", total=len(spec.targets))
        log = logger.bind(
            cpf_cnpj=spec.cpf_cnpj[:6] + "***",
            total_targets=len(spec.targets),
            mes_ano=spec.mes_ano,
            job_id=job_id,
        )
        log.info("batch_iniciado")

        if job_id:
            self._repo.marcar_job_em_execucao(job_id)

        try:
            async with self._client_factory(self._config.portal) as client:
                await client.login(
                    cpf_cnpj=spec.cpf_cnpj,
                    senha=spec.senha_portal,
                    uf=spec.uf,
                    tipo_acesso=spec.tipo_acesso,
                )

                for target in spec.targets:
                    if job_id:
                        self._repo.marcar_item_em_execucao(job_id, target.uc, attempts=1)

                    item_result = await self._process_target(client, spec, target)
                    result.items.append(item_result)
                    result.completed += 1
                    if item_result.status in ("sucesso", "pulado"):
                        result.success += 1
                    else:
                        result.error += 1

                    if job_id:
                        self._repo.salvar_resultado_item(job_id, item_result)

                    await client.delay_entre_clientes()

        except (LoginError, CaptchaError) as exc:
            log.error("batch_falhou_antes_itens", error_type=_error_type(exc), erro=str(exc))
            for target in spec.targets[len(result.items):]:
                item_result = BatchItemResult(
                    uc=target.uc,
                    nome=target.nome,
                    status=f"erro_{_error_type(exc).replace('Error', '').lower()}",
                    mensagem=str(exc),
                    error_type=_error_type(exc),
                    attempts=1,
                )
                result.items.append(item_result)
                result.completed += 1
                result.error += 1
                if job_id:
                    self._repo.salvar_resultado_item(job_id, item_result)
                self._registrar_log_erro(target.uc, spec.mes_ano, item_result.status, item_result.mensagem)

        result.status = _final_job_status(result)
        if job_id:
            self._repo.finalizar_job(job_id, result.status)

        log.info(
            "batch_concluido",
            status=result.status,
            completed=result.completed,
            success=result.success,
            error=result.error,
        )
        return result

    async def _process_target(
        self,
        client: CoelbaClient,
        spec: BatchSpec,
        target: BatchTarget,
    ) -> BatchItemResult:
        mes, ano, mes_ano_str = _resolver_mes_ano(spec.mes_ano)
        log = logger.bind(uc=target.uc, nome=target.nome)

        if mes and ano and not spec.force and self._repo.conta_existe(target.uc, mes, ano):
            log.info("conta_ja_existe_pulando")
            return BatchItemResult(
                uc=target.uc,
                nome=target.nome,
                status="pulado",
                mensagem="Já existe no banco",
                mes=mes,
                ano=ano,
                attempts=1,
            )

        try:
            pdf_path = await client.baixar_fatura(uc=target.uc, mes_ano=mes_ano_str)
            conta = self._parser.parse(pdf_path)

            if not conta.uc:
                conta.uc = target.uc
            elif _ucs_equivalentes(conta.uc, target.uc):
                conta.uc = target.uc
            if mes and ano:
                conta.mes = mes
                conta.ano = ano

            if not spec.force and self._repo.conta_existe(conta.uc, conta.mes, conta.ano):
                log.info("conta_ja_existe_pulando", mes=conta.mes, ano=conta.ano)
                return BatchItemResult(
                    uc=target.uc,
                    nome=target.nome,
                    status="pulado",
                    mensagem=f"Já existe: {conta.mes:02d}/{conta.ano}",
                    mes=conta.mes,
                    ano=conta.ano,
                    attempts=1,
                )

            conta_id = self._repo.salvar_conta(conta)
            self._repo.registrar_log(conta.uc, conta.mes, conta.ano, "sucesso")
            ocr_data = conta_para_ocr_payload(conta)

            return BatchItemResult(
                uc=target.uc,
                nome=target.nome,
                status="sucesso",
                mensagem=f"{conta.mes:02d}/{conta.ano} — R$ {conta.valor}",
                pdf_path=conta.pdf_path,
                mes=conta.mes,
                ano=conta.ano,
                valor=conta.valor,
                data_vencimento=formatar_data_br(conta.vencimento),
                normalizado_valor=float(conta.valor),
                ocr_data=ocr_data,
                conta_id=conta_id,
                attempts=1,
            )

        except (DownloadError, ParserError, LayoutChangedError, SessionExpiredError, FaturaError) as exc:
            log.error("target_falhou", error_type=_error_type(exc), erro=str(exc))
            status = self._status_for_exception(exc)
            context = client.last_error_context
            self._registrar_log_erro(target.uc, spec.mes_ano, status, str(exc))
            return BatchItemResult(
                uc=target.uc,
                nome=target.nome,
                status=status,
                mensagem=str(exc),
                error_type=_error_type(exc),
                screenshot_path=context["screenshot_path"],
                html_path=context["html_path"],
                attempts=1,
            )

    def _registrar_log_erro(
        self, uc: str, mes_ano: str | None, status: str, mensagem: str
    ) -> None:
        mes, ano, _ = _resolver_mes_ano(mes_ano)
        self._repo.registrar_log(
            uc=uc,
            mes=mes or datetime.now().month,
            ano=ano or datetime.now().year,
            status=status,
            mensagem=mensagem,
        )

    def _status_for_exception(self, exc: Exception) -> str:
        if isinstance(exc, DownloadError):
            return "erro_download"
        if isinstance(exc, ParserError):
            return "erro_parser"
        if isinstance(exc, LayoutChangedError):
            return "erro_layout"
        if isinstance(exc, SessionExpiredError):
            return "erro_sessao"
        if isinstance(exc, CaptchaError):
            return "erro_captcha"
        if isinstance(exc, LoginError):
            return "erro_login"
        return "erro"


async def executar_job_persistido(
    config: AppConfig,
    job_id: str,
    repo: SqliteFaturaRepository | None = None,
    client_factory: type[CoelbaClient] = CoelbaClient,
    request: FaturaJobRequest | None = None,
) -> BatchRunResult:
    repo = repo or SqliteFaturaRepository(config.database.url)
    request = request or repo.carregar_job_request(job_id)
    spec = BatchSpec(
        cpf_cnpj=request.cpf_cnpj,
        senha_portal=request.senha_portal,
        uf=request.uf,
        tipo_acesso=request.tipo_acesso,
        targets=[BatchTarget(uc=item.uc, nome=item.nome) for item in request.ucs],
        mes_ano=request.mes_ano,
        force=request.force,
    )
    processor = BatchProcessor(config=config, repo=repo, client_factory=client_factory)
    return await processor.run_batch(spec, job_id=job_id)


async def processar_faturas_mes(
    config: AppConfig,
    clientes: list[ClienteConfig],
    mes_ano: str | None = None,
    force: bool = False,
) -> ProcessamentoResult:
    resultado = ProcessamentoResult(total=len(clientes))
    processor = BatchProcessor(config=config)

    sorted_clientes = sorted(
        clientes,
        key=lambda cliente: (
            cliente.cpf_cnpj,
            cliente.senha_portal,
            cliente.uf,
            cliente.tipo_acesso.value,
        ),
    )
    grupos = groupby(
        sorted_clientes,
        key=lambda cliente: (
            cliente.cpf_cnpj,
            cliente.senha_portal,
            cliente.uf,
            cliente.tipo_acesso,
        ),
    )

    for (cpf_cnpj, senha_portal, uf, tipo_acesso), grupo_clientes in grupos:
        clientes_grupo = list(grupo_clientes)
        spec = BatchSpec(
            cpf_cnpj=cpf_cnpj,
            senha_portal=senha_portal,
            uf=uf,
            tipo_acesso=tipo_acesso,
            targets=[BatchTarget(uc=cliente.uc, nome=cliente.nome) for cliente in clientes_grupo],
            mes_ano=mes_ano,
            force=force,
        )
        batch_result = await processor.run_batch(spec)

        for item in batch_result.items:
            resultado.detalhes.append(
                DetalheProcessamento(
                    uc=item.uc,
                    nome=item.nome,
                    status=item.status,
                    mensagem=item.mensagem,
                )
            )
            if item.status == "sucesso":
                resultado.sucesso += 1
            elif item.status == "pulado":
                resultado.pulado += 1
            else:
                resultado.erro += 1

    return resultado


def criar_job_request_de_clientes(
    cpf_cnpj: str,
    senha_portal: str,
    uf: str,
    tipo_acesso,
    clientes: list[ClienteConfig],
    mes_ano: str | None,
    force: bool,
) -> FaturaJobRequest:
    return FaturaJobRequest(
        cpf_cnpj=cpf_cnpj,
        senha_portal=senha_portal,
        uf=uf,
        tipo_acesso=tipo_acesso,
        ucs=[{"uc": cliente.uc, "nome": cliente.nome} for cliente in clientes],
        mes_ano=mes_ano,
        force=force,
    )
