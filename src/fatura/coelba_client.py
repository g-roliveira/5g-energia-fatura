import asyncio
import random
import re
from dataclasses import dataclass
from datetime import datetime
from pathlib import Path

import structlog
from playwright.async_api import Browser, BrowserContext, Page, async_playwright
from tenacity import retry, retry_if_exception_type, stop_after_attempt, wait_exponential

from fatura.config import PortalConfig, TipoAcesso
from fatura.exceptions import (
    CaptchaError,
    DownloadError,
    LayoutChangedError,
    LoginError,
    SessionExpiredError,
)

logger = structlog.get_logger()

_USER_AGENT = (
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 "
    "(KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
)

_MESES = {
    1: "JANEIRO", 2: "FEVEREIRO", 3: "MARÇO", 4: "ABRIL",
    5: "MAIO", 6: "JUNHO", 7: "JULHO", 8: "AGOSTO",
    9: "SETEMBRO", 10: "OUTUBRO", 11: "NOVEMBRO", 12: "DEZEMBRO",
}


def _mes_ano_para_label(mes_ano: str) -> str:
    """Converte '042026' → 'ABRIL/2026' (formato exibido no portal)."""
    mes = int(mes_ano[:2])
    ano = mes_ano[2:]
    return f"{_MESES[mes]}/{ano}"


@dataclass
class FaturaInfo:
    referencia: str        # ex: "DEZEMBRO/2024"
    vencimento: str        # ex: "08/01/25"
    valor: str             # ex: "R$ 113,12"
    situacao: str          # ex: "Pago"
    data_pagamento: str    # ex: "06/02/25"


def _normalizar_texto(texto: str | None) -> str:
    if not texto:
        return ""
    return " ".join(texto.replace("\xa0", " ").split())


class CoelbaClient:
    """Cliente Playwright para o portal Agência Virtual Neoenergia.

    Portal: https://agenciavirtual.neoenergia.com
    Framework: Angular Material (SPA com hash routing)
    Anti-bot: Akamai (requer headed mode ou user-agent real + webdriver flag removida)

    Fluxo de navegação confirmado:
        /#/login
          → preencher #userId (CPF/CNPJ sem pontuação) + #password → botão ENTRAR
        /#/home/selecionar-estado
          → clicar no card do estado (ex: "Bahia")
        /#/home/meus-imoveis
          → clicar no texto da UC desejada
        /#/home  (dashboard da UC)
          → clicar no card "Faturas e 2ª via de faturas"
        /#/home/servicos/consultar-debitos
          → selecionar checkbox da fatura desejada
          → clicar button#btn-baixar-faturas
          → capturar download via page.expect_download()
    """

    def __init__(self, config: PortalConfig):
        self._config = config
        self._playwright = None
        self._browser: Browser | None = None
        self._context: BrowserContext | None = None
        self._page: Page | None = None
        self._logged_in = False
        self._uf_selecionada: str | None = None
        self._last_step_name: str | None = None
        self._last_screenshot_path: str | None = None
        self._last_html_path: str | None = None

    async def __aenter__(self) -> "CoelbaClient":
        self._playwright = await async_playwright().start()
        launch_kwargs = {
            "headless": self._config.headless,
            "args": [
                "--disable-blink-features=AutomationControlled",
                "--disable-dev-shm-usage",
                "--no-sandbox",
            ],
        }
        if self._config.browser_channel:
            launch_kwargs["channel"] = self._config.browser_channel
        if self._config.browser_executable_path:
            launch_kwargs["executable_path"] = self._config.browser_executable_path

        self._browser = await self._playwright.chromium.launch(
            **launch_kwargs,
        )
        download_path = Path(self._config.download_dir)
        download_path.mkdir(parents=True, exist_ok=True)
        self._context = await self._browser.new_context(
            accept_downloads=True,
            viewport={"width": 1280, "height": 900},
            locale="pt-BR",
            user_agent=_USER_AGENT,
        )
        self._context.set_default_timeout(self._config.timeout_ms)
        self._page = await self._context.new_page()
        # Remove webdriver flag para contornar Akamai
        await self._page.add_init_script(
            'Object.defineProperty(navigator, "webdriver", {get: () => undefined})'
        )
        return self

    async def __aexit__(self, *args) -> None:
        if self._context:
            await self._context.close()
        if self._browser:
            await self._browser.close()
        if self._playwright:
            await self._playwright.stop()
        self._logged_in = False

    @property
    def page(self) -> Page:
        if not self._page:
            raise RuntimeError("CoelbaClient não inicializado. Use 'async with'.")
        return self._page

    @property
    def context(self) -> BrowserContext:
        if not self._context:
            raise RuntimeError("CoelbaClient não inicializado. Use 'async with'.")
        return self._context

    @property
    def last_error_context(self) -> dict[str, str | None]:
        return {
            "step_name": self._last_step_name,
            "screenshot_path": self._last_screenshot_path,
            "html_path": self._last_html_path,
        }

    # -------------------------------------------------------------------------
    # Público
    # -------------------------------------------------------------------------

    async def login(
        self,
        cpf_cnpj: str,
        senha: str,
        uf: str = "BA",
        tipo_acesso: TipoAcesso = TipoAcesso.NORMAL,
    ) -> None:
        self._set_step("login")
        log = logger.bind(uf=uf, tipo_acesso=tipo_acesso.value)
        log.info("iniciando_login")

        try:
            login_url = self._config.url_base.rstrip("/") + "/#/login"
            await self.page.goto(login_url, wait_until="domcontentloaded")
            await self._garantir_formulario_login_visivel()
            log.debug("formulario_login_visivel")

            # Preencher CPF/CNPJ (apenas dígitos)
            digitos = "".join(c for c in cpf_cnpj if c.isdigit())
            await self.page.fill("#userId", digitos)

            # Preencher senha
            await self.page.fill("#password", senha)
            await asyncio.sleep(0.5)

            # Aguardar botão ENTRAR habilitar
            entrar = self.page.locator('button:has-text("ENTRAR")')
            await entrar.wait_for(state="visible", timeout=5_000)
            if await entrar.is_disabled():
                await asyncio.sleep(2)
                if await entrar.is_disabled():
                    await self._screenshot_on_error("entrar_desabilitado")
                    raise LoginError(
                        "Botão ENTRAR permanece desabilitado. "
                        "Verifique se CPF/CNPJ e senha estão corretos."
                    )

            await entrar.click()
            log.debug("botao_entrar_clicado")

            # Aguardar navegação pós-login (sai de /#/login)
            await self.page.wait_for_url(
                lambda u: "#/login" not in u, timeout=30_000
            )
            await self._wait_stable()

            # Selecionar estado se necessário
            if "selecionar-estado" in self.page.url:
                await self._selecionar_estado(uf)

            # Verificar erros de login
            await self._verificar_erro_login()

            self._logged_in = True
            self._uf_selecionada = uf
            log.info("login_sucesso", url=self.page.url)

        except (CaptchaError, LoginError):
            raise
        except Exception as e:
            await self._screenshot_on_error("login_erro")
            raise LoginError(f"Falha no login: {e}") from e

    async def listar_faturas(self, uc: str) -> list[FaturaInfo]:
        """Retorna a lista de faturas disponíveis para a UC."""
        if not self._logged_in:
            raise LoginError("Não autenticado. Chame login() primeiro.")

        self._set_step("listar_faturas")
        log = logger.bind(uc=uc)
        log.info("listando_faturas")

        try:
            await self._navegar_para_uc(uc)
            await self._abrir_consultar_debitos()
            faturas = await self._extrair_lista_faturas()
            log.info("faturas_encontradas", quantidade=len(faturas))
            return faturas

        except (LoginError, LayoutChangedError):
            raise
        except Exception as e:
            await self._screenshot_on_error(f"listar_faturas_{uc}")
            raise LayoutChangedError(f"Erro ao listar faturas: {e}", uc=uc) from e

    @retry(
        stop=stop_after_attempt(3),
        wait=wait_exponential(multiplier=2, min=4, max=30),
        retry=retry_if_exception_type((DownloadError, SessionExpiredError)),
    )
    async def baixar_fatura(
        self,
        uc: str,
        mes_ano: str | None = None,
        destino_dir: str | None = None,
    ) -> Path:
        """Baixa o PDF da fatura para a UC e competência indicadas.

        Args:
            uc: Código da unidade consumidora.
            mes_ano: Competência no formato MMAAAA (ex: '042026').
                     Se None, baixa a fatura mais recente disponível.
            destino_dir: Diretório de destino. Padrão: config.download_dir/uc/.

        Returns:
            Path para o PDF baixado.
        """
        if not self._logged_in:
            raise LoginError("Não autenticado. Chame login() primeiro.")

        self._set_step("baixar_fatura")
        log = logger.bind(uc=uc, mes_ano=mes_ano)
        log.info("baixando_fatura")

        dest_dir = Path(destino_dir or self._config.download_dir) / uc
        dest_dir.mkdir(parents=True, exist_ok=True)

        try:
            await self._navegar_para_uc(uc)
            await self._abrir_consultar_debitos()
            await self._selecionar_fatura(mes_ano)
            self._set_step("confirmar_download_fatura")

            filename = mes_ano or datetime.now().strftime("%m%Y")
            dest_path = dest_dir / f"{filename}.pdf"

            download_task = asyncio.create_task(
                self.page.wait_for_event("download", timeout=60_000)
            )
            try:
                await self.page.locator("#btn-baixar-faturas").click()
                await self._confirmar_modal_segunda_via_se_necessario()
                download = await download_task
            except Exception:
                if not download_task.done():
                    download_task.cancel()
                raise

            await download.save_as(dest_path)
            log.info("fatura_baixada", path=str(dest_path), size_bytes=dest_path.stat().st_size)
            return dest_path

        except (DownloadError, SessionExpiredError):
            raise
        except Exception as e:
            await self._screenshot_on_error(f"download_{uc}_{mes_ano}")
            raise DownloadError(f"Falha ao baixar fatura: {e}", uc=uc) from e

    # -------------------------------------------------------------------------
    # Navegação interna (fluxo real confirmado por exploração)
    # -------------------------------------------------------------------------

    async def _selecionar_estado(self, uf: str) -> None:
        """Seleciona o estado na tela /#/home/selecionar-estado."""
        self._set_step("selecionar_estado")
        uf_para_nome = {
            "BA": "Bahia",
            "PE": "Pernambuco",
            "RN": "Rio Grande do Norte",
            "SP": "São Paulo",
            "MS": "Mato Grosso do Sul",
        }
        nome_estado = uf_para_nome.get(uf, "Bahia")
        logger.debug("selecionando_estado", uf=uf, nome=nome_estado)

        await self.page.locator(f"text={nome_estado}").first.click()
        await self.page.wait_for_url("**/meus-imoveis**", timeout=15_000)
        await self._wait_stable()
        logger.debug("estado_selecionado", url=self.page.url)

    async def _navegar_para_uc(self, uc: str) -> None:
        """A partir da lista de imóveis, seleciona a UC desejada.

        Se já estiver no dashboard da UC correta, não faz nada.
        """
        self._set_step("navegar_para_uc")
        # Se não estiver na lista de imóveis, voltar para lá
        if "meus-imoveis" not in self.page.url and "home/selecionar-estado" not in self.page.url:
            # Já pode estar no dashboard de uma UC — verificar se é a certa
            uc_atual = await self._obter_uc_atual()
            if uc_atual == uc:
                logger.debug("ja_na_uc_correta", uc=uc)
                return
            # Clicar em "Minhas Unidades" para voltar
            await self._voltar_para_lista_ucs()

        if "selecionar-estado" in self.page.url:
            await self._selecionar_estado(self._uf_selecionada or "BA")

        logger.debug("navegando_para_uc", uc=uc, url_atual=self.page.url)

        # Clicar no texto da UC na lista
        await self.page.locator(f"text={uc}").first.click()

        # Aguardar loading spinner desaparecer
        try:
            await self.page.wait_for_selector(
                "text=Aguarde um instante", state="visible", timeout=5_000
            )
            await self.page.wait_for_selector(
                "text=Aguarde um instante", state="hidden", timeout=30_000
            )
        except Exception:
            pass

        await self._wait_stable()
        logger.debug("uc_selecionada", uc=uc, url=self.page.url)

    async def _abrir_consultar_debitos(self) -> None:
        """Abre a tela de consulta de débitos/faturas (/#/home/servicos/consultar-debitos)."""
        self._set_step("abrir_consultar_debitos")
        if "consultar-debitos" in self.page.url:
            return

        logger.debug("abrindo_consultar_debitos", url_atual=self.page.url)

        # Clicar no card "Faturas e 2ª via de faturas" no dashboard
        try:
            card = self.page.locator(
                'mat-card:has-text("Faturas e 2ª via de faturas")'
            ).first
            if await card.is_visible(timeout=5_000):
                await card.click()
                await self._aguardar_loading()
                if "consultar-debitos" in self.page.url:
                    logger.debug("consultar_debitos_aberto_via_card")
                    return
        except Exception:
            pass

        # Fallback: clicar no link do sidebar
        try:
            link = self.page.locator(
                'a:has-text("Faturas e 2ª Via de Faturas")'
            ).first
            if await link.is_visible(timeout=3_000):
                await link.click()
                await self._aguardar_loading()
                logger.debug("consultar_debitos_aberto_via_sidebar")
                return
        except Exception:
            pass

        await self._screenshot_on_error("consultar_debitos_nao_encontrado")
        raise LayoutChangedError(
            "Não encontrou link para 'Faturas e 2ª Via de Faturas'."
        )

    async def _selecionar_fatura(self, mes_ano: str | None) -> None:
        """Seleciona o checkbox da fatura desejada na lista.

        Se mes_ano for None, seleciona a primeira (mais recente).
        Aguarda o botão #btn-baixar-faturas ficar habilitado.
        """
        self._set_step("selecionar_fatura")
        await self.page.wait_for_selector(
            "text=LISTA DE FATURAS", state="visible", timeout=15_000
        )
        await self._wait_stable()

        paineis = await self._listar_paineis_faturas()

        if mes_ano:
            label = _mes_ano_para_label(mes_ano)
            logger.debug("selecionando_fatura", label=label)
            painel = None
            for item in paineis:
                referencia = (await self._extrair_campos_fatura(item)).get("REFERÊNCIA", "")
                if referencia == label:
                    painel = item
                    break

            if painel is None:
                await self._screenshot_on_error(f"fatura_nao_encontrada_{mes_ano}")
                raise DownloadError(f"Fatura {label} não encontrada na lista de faturas.")

        else:
            logger.debug("selecionando_fatura_mais_recente")
            painel = paineis[0]

        checkbox = painel.locator("mat-checkbox label").first
        await checkbox.click()

        # Aguardar botão Download habilitar
        dl_btn = self.page.locator("#btn-baixar-faturas")
        await asyncio.sleep(1)
        is_disabled = await dl_btn.is_disabled()
        if is_disabled:
            await self._screenshot_on_error("download_btn_desabilitado")
            raise DownloadError(
                "Botão de download permanece desabilitado após selecionar a fatura."
            )

        logger.debug("fatura_selecionada_download_pronto")

    async def _confirmar_modal_segunda_via_se_necessario(self) -> None:
        """Confirma o modal de motivo da 2ª via quando o portal exige justificativa."""
        dialog = self.page.locator("#dialog-baixar-segunda-via").first
        try:
            await dialog.wait_for(state="visible", timeout=5_000)
        except Exception:
            return

        self._set_step("confirmar_modal_segunda_via")
        logger.debug("confirmando_modal_segunda_via")

        motivo_preferencial = dialog.locator(
            'mat-radio-button:has-text("Não Estou Com Fatura Em Mãos") label'
        ).first
        if await motivo_preferencial.count() > 0:
            await motivo_preferencial.click()
        else:
            await dialog.locator("mat-radio-button label").first.click()

        baixar = dialog.locator('button:has-text("BAIXAR")').first
        await baixar.wait_for(state="visible", timeout=10_000)
        for _ in range(20):
            if not await baixar.is_disabled():
                break
            await asyncio.sleep(0.25)

        if await baixar.is_disabled():
            await self._screenshot_on_error("modal_segunda_via_baixar_desabilitado")
            raise DownloadError(
                "Botão BAIXAR do modal de 2ª via permaneceu desabilitado."
            )

        await baixar.click()

    async def _extrair_lista_faturas(self) -> list[FaturaInfo]:
        """Extrai as informações das faturas visíveis na lista."""
        self._set_step("extrair_lista_faturas")
        await self.page.wait_for_selector("text=LISTA DE FATURAS", state="visible", timeout=15_000)
        await self._wait_stable()
        faturas = []
        paineis = await self._listar_paineis_faturas()

        for i, painel in enumerate(paineis):
            try:
                campos = await self._extrair_campos_fatura(painel)
                referencia = campos.get("REFERÊNCIA", "")
                vencimento = campos.get("VENCIMENTO", "")
                valor = campos.get("VALOR FATURA", "")
                situacao = campos.get("SITUAÇÃO", "")
                data_pagamento = campos.get("DATA PAGAMENTO", "")

                if referencia:
                    faturas.append(FaturaInfo(
                        referencia=referencia,
                        vencimento=vencimento,
                        valor=valor,
                        situacao=situacao,
                        data_pagamento=data_pagamento,
                    ))
            except Exception as e:
                logger.warning("erro_extrair_fatura_item", indice=i, erro=str(e))

        return faturas

    async def _listar_paineis_faturas(self) -> list:
        accordion = self.page.locator("#accordion-fatura")
        try:
            await accordion.wait_for(state="visible", timeout=15_000)
        except Exception as exc:
            await self._screenshot_on_error("accordion_faturas_nao_encontrado")
            raise LayoutChangedError(
                "Não encontrou o accordion principal da lista de faturas."
            ) from exc

        paineis = accordion.locator("> mat-expansion-panel")
        count = await paineis.count()
        if count == 0:
            corpo = _normalizar_texto(await self.page.locator("body").text_content())
            if "NENHUMA FATURA" in corpo.upper():
                return []
            await self._screenshot_on_error("nenhum_painel_fatura_encontrado")
            raise LayoutChangedError(
                "Accordion de faturas encontrado, mas sem painéis de fatura."
            )

        return [paineis.nth(i) for i in range(count)]

    async def _extrair_campos_fatura(self, painel) -> dict[str, str]:
        header = painel.locator("mat-expansion-panel-header").first
        blocos = header.locator(".fatura-situacao")
        campos: dict[str, str] = {}

        for i in range(await blocos.count()):
            textos = [
                _normalizar_texto(texto)
                for texto in await blocos.nth(i).locator("span").all_text_contents()
            ]
            textos = [texto for texto in textos if texto]
            if len(textos) < 2:
                continue

            rotulo = textos[0].upper()
            valor = textos[-1]
            campos[rotulo] = valor

        if "REFERÊNCIA" not in campos:
            texto_header = _normalizar_texto(await header.text_content())
            referencia_match = re.search(
                r"(JANEIRO|FEVEREIRO|MARÇO|ABRIL|MAIO|JUNHO|JULHO|AGOSTO|SETEMBRO|OUTUBRO|NOVEMBRO|DEZEMBRO)/\d{4}",
                texto_header,
                re.I,
            )
            if referencia_match:
                campos["REFERÊNCIA"] = referencia_match.group(0)

        return campos

    # -------------------------------------------------------------------------
    # Utilitários internos
    # -------------------------------------------------------------------------

    async def _obter_uc_atual(self) -> str | None:
        """Tenta ler o código da UC exibida no dashboard."""
        try:
            uc_span = self.page.locator("span.mat-icon.home_work + span, h2 + span").first
            return await uc_span.text_content(timeout=2_000)
        except Exception:
            return None

    async def _voltar_para_lista_ucs(self) -> None:
        """Volta para a lista de imóveis/UCs."""
        logger.debug("voltando_para_lista_ucs")
        try:
            btn = self.page.locator('button:has-text("Minhas Unidades"), a:has-text("Minhas Unidades Consumidoras")').first
            if await btn.is_visible(timeout=3_000):
                await btn.click()
                await self.page.wait_for_url("**/meus-imoveis**", timeout=15_000)
                await self._wait_stable()
                return
        except Exception:
            pass
        # Fallback: navegar direto
        base = self._config.url_base.rstrip("/")
        await self.page.goto(f"{base}/#/home/meus-imoveis", wait_until="domcontentloaded")
        await self._wait_stable()

    async def _aguardar_loading(self) -> None:
        """Aguarda spinner 'Aguarde um instante...' desaparecer."""
        try:
            await self.page.wait_for_selector(
                "text=Aguarde um instante", state="visible", timeout=3_000
            )
            await self.page.wait_for_selector(
                "text=Aguarde um instante", state="hidden", timeout=30_000
            )
        except Exception:
            pass
        await self._wait_stable()

    async def _verificar_erro_login(self) -> None:
        """Detecta mensagens de erro exibidas pelo portal após tentativa de login."""
        seletores = [
            "mat-snack-bar-container",
            "simple-snack-bar",
            ".mat-error",
        ]
        for sel in seletores:
            try:
                elem = self.page.locator(sel).first
                if await elem.is_visible(timeout=2_000):
                    texto = await elem.text_content()
                    if texto and texto.strip():
                        raise LoginError(f"Portal retornou erro: {texto.strip()}")
            except LoginError:
                raise
            except Exception:
                continue

        # Mensagens de erro por texto
        try:
            err = self.page.locator(
                "text=/senha incorreta|usuário não encontrado|credenciais inválidas|acesso negado/i"
            ).first
            if await err.is_visible(timeout=2_000):
                texto = await err.text_content()
                raise LoginError(f"Erro de autenticação: {texto}")
        except LoginError:
            raise
        except Exception:
            pass

        body_text = (await self.page.locator("body").text_content(timeout=2_000) or "").strip()
        if "Access Denied" in body_text or "You don't have permission" in body_text:
            raise LoginError("Portal bloqueou o acesso antes do formulário de login.")

    async def _wait_stable(self) -> None:
        """Aguarda SPA Angular estabilizar."""
        try:
            await self.page.wait_for_load_state("networkidle", timeout=10_000)
        except Exception:
            pass
        await asyncio.sleep(1)

    async def _screenshot_on_error(self, name: str) -> None:
        """Salva screenshot e HTML reduzido para debugging."""
        try:
            error_dir = Path(self._config.download_dir) / "_errors"
            error_dir.mkdir(parents=True, exist_ok=True)
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            screenshot_path = error_dir / f"{name}_{timestamp}.png"
            html_path = error_dir / f"{name}_{timestamp}.html"
            await self.page.screenshot(path=screenshot_path, full_page=True)
            html = await self.page.locator("body").inner_html(timeout=5_000)
            html_path.write_text(html, encoding="utf-8")
            self._last_screenshot_path = str(screenshot_path)
            self._last_html_path = str(html_path)
            logger.info(
                "evidencia_erro_salva",
                screenshot_path=str(screenshot_path),
                html_path=str(html_path),
                step_name=self._last_step_name,
            )
        except Exception as e:
            logger.warning("screenshot_falhou", erro=str(e))

    async def delay_entre_clientes(self) -> None:
        """Delay com jitter para evitar rate limiting."""
        base = self._config.delay_between_clients_s
        jitter = random.uniform(0, base * 0.3)
        await asyncio.sleep(base + jitter)

    def _set_step(self, name: str) -> None:
        self._last_step_name = name

    async def _garantir_formulario_login_visivel(self) -> None:
        """Abre o formulário de login mesmo quando a rota /#/login carrega a landing page."""
        try:
            await self.page.wait_for_selector("#userId", state="visible", timeout=8_000)
            return
        except Exception:
            pass

        login_button = self.page.locator(
            'button[aria-label*="Conectar-se"], button:has-text("LOGIN")'
        ).first
        if await login_button.count() > 0:
            try:
                if await login_button.is_visible(timeout=5_000):
                    await login_button.click()
            except Exception:
                pass

        await self.page.wait_for_selector("#userId", state="visible", timeout=30_000)
        await self.page.wait_for_selector("#password", state="visible", timeout=10_000)
