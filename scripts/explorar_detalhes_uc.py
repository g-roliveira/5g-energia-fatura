"""Exploração da tela de detalhes da UC após login.
Faz login, seleciona Bahia, clica na UC e explora a tela de detalhes/faturas.

Uso: python scripts/explorar_detalhes_uc.py
"""

import asyncio
from pathlib import Path

import yaml
from playwright.async_api import async_playwright

SCREENSHOTS_DIR = Path("downloads/_exploracao")
URL = "https://agenciavirtual.neoenergia.com"
USER_AGENT = (
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 "
    "(KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
)

with open("config.yaml") as f:
    cfg = yaml.safe_load(f)
    cliente = cfg["clientes"][0]
    CPF_CNPJ = "".join(c for c in cliente["cpf_cnpj"] if c.isdigit())
    SENHA = cliente["senha_portal"]
    UC = cliente["uc"]


async def dump_page(page, step, name):
    SCREENSHOTS_DIR.mkdir(parents=True, exist_ok=True)
    fname = f"{step:02d}_{name}"
    await page.screenshot(path=SCREENSHOTS_DIR / f"{fname}.png", full_page=True)
    print(f"\n[STEP {step}] {name} | URL: {page.url}")

    # Dump todos os elementos visíveis relevantes
    all_text = await page.evaluate("""() => {
        const els = document.querySelectorAll('h1,h2,h3,h4,h5,p,span,a,button,label,mat-card-title,mat-card-subtitle,td,th');
        const seen = new Set();
        const results = [];
        for (const el of els) {
            if (el.offsetParent === null) continue;
            const t = el.textContent.trim().replace(/\\s+/g, ' ');
            if (t.length > 2 && t.length < 150 && !seen.has(t)) {
                seen.add(t);
                results.push({tag: el.tagName, text: t, class: el.className?.substring?.(0,40) || ''});
            }
        }
        return results.slice(0, 50);
    }""")
    for item in all_text:
        print(f"  <{item['tag']:20s}> {item['text']}")

    # Inputs
    inputs = await page.eval_on_selector_all("input:not([type=hidden]), select, textarea", """els => els.filter(e => e.offsetParent !== null).map(e => ({
        type: e.type, id: e.id, placeholder: e.placeholder,
        formControl: e.getAttribute('formcontrolname'),
    }))""")
    if inputs:
        print(f"  --- Inputs ({len(inputs)}) ---")
        for i in inputs:
            print(f"    type={i['type']} id={i['id']} placeholder='{i['placeholder']}' formcontrolname={i.get('formControl','')}")

    # Salvar HTML
    html = await page.content()
    Path(SCREENSHOTS_DIR / f"{fname}.html").write_text(html)


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

        # LOGIN
        print("[*] Login...")
        await page.goto(f"{URL}/#/login", wait_until="domcontentloaded")
        await page.wait_for_selector("#userId", state="visible", timeout=20_000)
        await page.fill("#userId", CPF_CNPJ)
        await page.fill("#password", SENHA)
        await asyncio.sleep(1)
        await page.locator('button:has-text("ENTRAR")').click()
        await page.wait_for_url(lambda u: "#/login" not in u, timeout=30_000)
        await page.wait_for_load_state("networkidle")
        print(f"  Login OK -> {page.url}")

        # SELECIONAR BAHIA
        if "selecionar-estado" in page.url:
            print("[*] Selecionando Bahia...")
            await page.locator("text=Bahia").first.click()
            await page.wait_for_url("**/meus-imoveis**", timeout=15_000)
            await page.wait_for_load_state("networkidle")
            await asyncio.sleep(2)
            print(f"  Bahia selecionada -> {page.url}")

        await dump_page(page, 1, "meus_imoveis")

        # CLICAR NA UC - usar a seta > (chevron_right) ao lado da UC
        print(f"\n[*] Procurando UC {UC} e clicando na seta >...")

        # Abordagem: encontrar o card que contém o texto da UC e clicar na seta dentro dele
        uc_card = page.locator(f"*:has-text('{UC}')").locator("..").locator("..").locator("mat-icon:text('chevron_right'), button:has(mat-icon)")

        # Alternativa mais robusta: procurar todos os cards e encontrar o da UC
        cards = page.locator(".card-wrapper, [class*='card-wrapper'], mat-card")
        count = await cards.count()
        print(f"  Cards encontrados: {count}")

        clicked = False
        for i in range(count):
            card = cards.nth(i)
            text = await card.text_content()
            if UC in text:
                print(f"  UC encontrada no card {i}: {text[:80]}")
                # Clicar na seta > dentro deste card
                chevron = card.locator("mat-icon:has-text('chevron_right'), [class*='chevron'], button").last
                try:
                    await chevron.click(timeout=5000)
                    clicked = True
                except Exception:
                    # Fallback: clicar no card inteiro
                    print("  Tentando clicar no card inteiro...")
                    await card.click()
                    clicked = True
                break

        if not clicked:
            # Fallback: clicar no texto da UC diretamente
            print(f"  Fallback: clicando no texto {UC}...")
            await page.locator(f"text={UC}").first.click()

        # Aguardar loading "Aguarde um instante..." desaparecer
        print("  Aguardando carregamento...")
        try:
            await page.wait_for_selector("text=Aguarde um instante", state="hidden", timeout=30_000)
        except Exception:
            pass
        await page.wait_for_load_state("networkidle")
        await asyncio.sleep(3)

        await dump_page(page, 2, "detalhes_uc")

        # Explorar a tela de detalhes
        print(f"\n[*] Tela de detalhes da UC. URL: {page.url}")

        # Verificar se há tabs ou menu lateral
        tabs = await page.eval_on_selector_all("[role='tab'], mat-tab, .mat-tab-label", """els => els.map(e => ({
            text: e.textContent.trim().substring(0,40),
            selected: e.getAttribute('aria-selected'),
        }))""")
        if tabs:
            print(f"\n  Tabs encontradas ({len(tabs)}):")
            for t in tabs:
                print(f"    [{t.get('selected','')}] \"{t['text']}\"")

        # Verificar menus laterais
        side_links = await page.eval_on_selector_all("mat-nav-list a, mat-list-item, [class*='sidebar'] a, [class*='menu-lateral'] a", """els => els.filter(e => e.offsetParent !== null).map(e => ({
            text: e.textContent.trim().substring(0,40),
            href: e.href,
        }))""")
        if side_links:
            print(f"\n  Menu lateral ({len(side_links)}):")
            for l in side_links:
                print(f"    \"{l['text']}\" -> {l['href']}")

        # Verificar se há link/botão de "2ª via" ou "faturas" agora
        fatura_links = await page.eval_on_selector_all("a, button, mat-card", """els => els.filter(e => {
            const t = e.textContent.toLowerCase();
            return e.offsetParent !== null && (t.includes('2ª via') || t.includes('segunda via') || t.includes('fatura') || t.includes('débito') || t.includes('financeiro') || t.includes('histórico') || t.includes('pagamento'));
        }).map(e => ({
            tag: e.tagName, text: e.textContent.trim().substring(0,60),
            href: e.href || '',
        }))""")
        if fatura_links:
            print(f"\n  Links relacionados a faturas ({len(fatura_links)}):")
            for l in fatura_links:
                print(f"    <{l['tag']}> \"{l['text']}\" -> {l['href']}")

        # Se a URL mudou para algo com "detalhes" ou "imovel", navegar pelos menus
        if "meus-imoveis" not in page.url:
            await dump_page(page, 3, "dentro_da_uc")

            # Tentar clicar em "2ª via" ou "faturas"
            for sel in ["text=/2[ªa] via/i", "text=/fatura/i", "text=/débito/i", "text=/financeiro/i"]:
                try:
                    elem = page.locator(sel).first
                    if await elem.is_visible(timeout=2000):
                        print(f"\n[*] Clicando em: {sel}")
                        await elem.click()
                        await asyncio.sleep(3)
                        await page.wait_for_load_state("networkidle")
                        await dump_page(page, 4, "tela_faturas")
                        break
                except Exception:
                    continue

        # Manter aberto
        print("\n[*] Browser aberto para inspeção. Pressione Enter para fechar...")
        await asyncio.get_event_loop().run_in_executor(None, input)
        await browser.close()


if __name__ == "__main__":
    asyncio.run(explorar())
