"""Exploração interativa do fluxo completo do portal Neoenergia.
Navega passo a passo, captura screenshots e DOM em cada etapa.

Uso: python scripts/explorar_fluxo_completo.py
"""

import asyncio
import os
from pathlib import Path

from playwright.async_api import async_playwright

SCREENSHOTS_DIR = Path("downloads/_exploracao")
URL = "https://agenciavirtual.neoenergia.com"
USER_AGENT = (
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 "
    "(KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
)

# Credenciais do config.yaml (lê do arquivo)
import yaml
with open("config.yaml") as f:
    cfg = yaml.safe_load(f)
    cliente = cfg["clientes"][0]
    CPF_CNPJ = "".join(c for c in cliente["cpf_cnpj"] if c.isdigit())
    SENHA = cliente["senha_portal"]
    UC = cliente["uc"]


async def dump_page(page, step: int, name: str):
    """Captura screenshot, inputs, botões e links de uma página."""
    SCREENSHOTS_DIR.mkdir(parents=True, exist_ok=True)

    fname = f"{step:02d}_{name}"
    await page.screenshot(path=SCREENSHOTS_DIR / f"{fname}.png", full_page=True)

    print(f"\n{'='*60}")
    print(f"[STEP {step}] {name}")
    print(f"  URL: {page.url}")
    print(f"  Screenshot: {fname}.png")

    # Inputs
    inputs = await page.eval_on_selector_all("input:not([type=hidden]), select, textarea", """els => els.map(e => ({
        tag: e.tagName, type: e.type, name: e.name, id: e.id,
        placeholder: e.placeholder,
        formControl: e.getAttribute('formcontrolname'),
        value: e.value.substring(0,30),
        visible: e.offsetParent !== null,
    }))""")
    if inputs:
        print(f"  Inputs ({len(inputs)}):")
        for i in inputs:
            vis = "✓" if i["visible"] else "✗"
            print(f"    {vis} <{i['tag']} type={i['type']} id={i['id']} formcontrolname={i.get('formControl','')} value='{i['value']}''>")

    # Botões
    buttons = await page.eval_on_selector_all("button", """els => els.map(e => ({
        text: e.textContent.trim().substring(0,50), disabled: e.disabled,
        class: e.className.substring(0,60),
    }))""")
    if buttons:
        print(f"  Botões ({len(buttons)}):")
        for b in buttons:
            dis = " [DISABLED]" if b["disabled"] else ""
            print(f"    \"{b['text']}\"{dis}")

    # Links significativos
    links = await page.eval_on_selector_all("a[href]", """els => els.filter(e => e.textContent.trim() && e.offsetParent !== null).map(e => ({
        text: e.textContent.trim().substring(0,50), href: e.href,
    }))""")
    if links:
        print(f"  Links ({len(links)}):")
        for l in links:
            print(f"    \"{l['text']}\" -> {l['href']}")

    # Textos relevantes (headings, cards, labels)
    headings = await page.eval_on_selector_all("h1, h2, h3, h4, mat-card-title, .mat-card-title", """els => els.filter(e => e.offsetParent !== null).map(e => ({
        tag: e.tagName, text: e.textContent.trim().substring(0,80),
    }))""")
    if headings:
        print(f"  Headings/Títulos:")
        for h in headings:
            print(f"    <{h['tag']}> \"{h['text']}\"")

    # Cards ou list items que possam ser estados/UCs/faturas
    cards = await page.eval_on_selector_all("mat-card, .mat-card, mat-list-item, .card, [class*='card'], [class*='item']", """els => els.filter(e => e.offsetParent !== null && e.textContent.trim().length > 3 && e.textContent.trim().length < 200).map(e => ({
        tag: e.tagName, class: e.className.substring(0,50),
        text: e.textContent.trim().substring(0,120),
    })).slice(0, 20)""")
    if cards:
        print(f"  Cards/Items ({len(cards)}):")
        for c in cards:
            print(f"    [{c['class'][:30]}] \"{c['text']}\"")

    # Mat-select (dropdowns Angular Material)
    selects = await page.eval_on_selector_all("mat-select, [role='listbox'], [role='combobox']", """els => els.map(e => ({
        id: e.id, class: e.className.substring(0,50),
        text: e.textContent.trim().substring(0,50),
        ariaLabel: e.getAttribute('aria-label'),
    }))""")
    if selects:
        print(f"  Mat-Selects ({len(selects)}):")
        for s in selects:
            print(f"    id={s['id']} \"{s['text']}\" aria-label={s.get('ariaLabel','')}")

    # Radio buttons / checkboxes
    radios = await page.eval_on_selector_all("mat-radio-button, mat-checkbox, [role='radio'], [role='checkbox']", """els => els.map(e => ({
        text: e.textContent.trim().substring(0,50),
        checked: e.getAttribute('aria-checked'),
        class: e.className.substring(0,50),
    }))""")
    if radios:
        print(f"  Radios/Checkboxes ({len(radios)}):")
        for r in radios:
            print(f"    [{r.get('checked','')}] \"{r['text']}\"")

    # Salvar HTML
    html = await page.content()
    Path(SCREENSHOTS_DIR / f"{fname}.html").write_text(html)

    print(f"{'='*60}")


async def explorar():
    async with async_playwright() as p:
        browser = await p.chromium.launch(
            headless=False,
            args=["--disable-blink-features=AutomationControlled"],
        )
        context = await browser.new_context(
            viewport={"width": 1280, "height": 900},
            locale="pt-BR",
            user_agent=USER_AGENT,
        )
        page = await context.new_page()
        await page.add_init_script('Object.defineProperty(navigator, "webdriver", {get: () => undefined})')

        # STEP 1: Login
        print("[*] Navegando para login...")
        await page.goto(f"{URL}/#/login", wait_until="domcontentloaded")
        await page.wait_for_load_state("networkidle")
        await asyncio.sleep(3)
        await dump_page(page, 1, "login_page")

        # STEP 2: Preencher e submeter login
        print("\n[*] Fazendo login...")
        await page.fill("#userId", CPF_CNPJ)
        await page.fill("#password", SENHA)
        await asyncio.sleep(1)

        entrar = page.locator('button:has-text("ENTRAR")')
        await entrar.click()
        await asyncio.sleep(5)
        await page.wait_for_load_state("networkidle")
        await dump_page(page, 2, "pos_login")

        # STEP 3: Selecionar estado (se estiver na página de seleção)
        if "selecionar-estado" in page.url:
            print("\n[*] Página de seleção de estado detectada. Procurando Bahia/Coelba...")
            await dump_page(page, 3, "selecionar_estado")

            # Tentar clicar em Bahia/Coelba
            bahia_selectors = [
                "text=/Bahia/i",
                "text=/Coelba/i",
                "text=/BA/",
                "img[alt*='Bahia' i]",
                "img[alt*='Coelba' i]",
                "[class*='bahia' i]",
                "[class*='coelba' i]",
            ]
            clicked = False
            for sel in bahia_selectors:
                try:
                    elem = page.locator(sel).first
                    if await elem.is_visible(timeout=2000):
                        print(f"    Clicando em: {sel}")
                        await elem.click()
                        await asyncio.sleep(3)
                        await page.wait_for_load_state("networkidle")
                        clicked = True
                        break
                except Exception:
                    continue

            if not clicked:
                print("    [!] Não encontrou seletor de Bahia. Dump da página para análise.")

            await dump_page(page, 4, "apos_selecao_estado")

        # STEP 5: Página principal após seleção de estado
        print(f"\n[*] Página atual: {page.url}")
        await asyncio.sleep(2)
        await dump_page(page, 5, "pagina_principal_logada")

        # STEP 6: Procurar e navegar para segunda via / faturas
        print("\n[*] Procurando link para segunda via ou faturas...")
        nav_selectors = [
            "text=/2[ªa] via/i",
            "text=/segunda via/i",
            "text=/fatura/i",
            "text=/conta/i",
            "text=/histórico/i",
            "text=/débito/i",
            "text=/pagamento/i",
            "text=/financeiro/i",
        ]
        for sel in nav_selectors:
            try:
                items = page.locator(sel)
                count = await items.count()
                for idx in range(min(count, 3)):
                    item = items.nth(idx)
                    if await item.is_visible(timeout=1000):
                        text = await item.text_content()
                        print(f"    Encontrado: \"{text.strip()[:50]}\" ({sel})")
            except Exception:
                continue

        # Tentar clicar em "2ª via" ou similar
        via_clicked = False
        for sel in nav_selectors[:3]:
            try:
                elem = page.locator(sel).first
                if await elem.is_visible(timeout=2000):
                    await elem.click()
                    await asyncio.sleep(3)
                    await page.wait_for_load_state("networkidle")
                    via_clicked = True
                    break
            except Exception:
                continue

        await dump_page(page, 6, "apos_navegar_faturas")

        # STEP 7: Se tem UC para selecionar
        print(f"\n[*] Página atual: {page.url}")
        await asyncio.sleep(2)
        await dump_page(page, 7, "tela_faturas_ou_uc")

        # STEP 8: Procurar a UC específica
        print(f"\n[*] Procurando UC {UC}...")
        try:
            uc_elem = page.locator(f"text={UC}").first
            if await uc_elem.is_visible(timeout=5000):
                print(f"    UC encontrada! Clicando...")
                await uc_elem.click()
                await asyncio.sleep(3)
                await page.wait_for_load_state("networkidle")
        except Exception:
            print(f"    UC {UC} não encontrada na página.")

        await dump_page(page, 8, "apos_selecao_uc")

        # STEP 9: Manter aberto para inspeção
        print("\n" + "="*60)
        print("[*] Exploração concluída. Browser aberto para inspeção manual.")
        print(f"    URL final: {page.url}")
        print(f"    Screenshots em: {SCREENSHOTS_DIR}/")
        print("    Pressione Enter para fechar...")
        print("="*60)
        await asyncio.get_event_loop().run_in_executor(None, input)

        await browser.close()


if __name__ == "__main__":
    asyncio.run(explorar())
