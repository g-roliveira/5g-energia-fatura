"""Exploração da tela de faturas.
Login -> Bahia -> UC -> Faturas e 2ª Via de Faturas -> dump completo.
"""

import asyncio
from pathlib import Path

import yaml
from playwright.async_api import async_playwright

DIR = Path("downloads/_exploracao")
URL = "https://agenciavirtual.neoenergia.com"
UA = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"

with open("config.yaml") as f:
    cfg = yaml.safe_load(f)
    c = cfg["clientes"][0]
    CPF = "".join(ch for ch in c["cpf_cnpj"] if ch.isdigit())
    SENHA = c["senha_portal"]
    UC = c["uc"]


async def dump(page, step, name):
    DIR.mkdir(parents=True, exist_ok=True)
    f = f"{step:02d}_{name}"
    await page.screenshot(path=DIR / f"{f}.png", full_page=True)
    html = await page.content()
    (DIR / f"{f}.html").write_text(html)
    print(f"  [{step}] {name} | {page.url} | screenshot saved")


async def main():
    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=False, args=["--disable-blink-features=AutomationControlled"])
        ctx = await browser.new_context(viewport={"width": 1280, "height": 900}, locale="pt-BR", user_agent=UA)
        page = await ctx.new_page()
        await page.add_init_script('Object.defineProperty(navigator, "webdriver", {get: () => undefined})')

        # 1. Login
        print("[1] Login...")
        await page.goto(f"{URL}/#/login", wait_until="domcontentloaded")
        await page.wait_for_selector("#userId", state="visible", timeout=20_000)
        await page.fill("#userId", CPF)
        await page.fill("#password", SENHA)
        await asyncio.sleep(1)
        await page.locator('button:has-text("ENTRAR")').click()
        await page.wait_for_url(lambda u: "#/login" not in u, timeout=30_000)
        await page.wait_for_load_state("networkidle")
        print(f"  OK -> {page.url}")

        # 2. Selecionar Bahia
        if "selecionar-estado" in page.url:
            print("[2] Selecionando Bahia...")
            await page.locator("text=Bahia").first.click()
            await page.wait_for_url("**/meus-imoveis**", timeout=15_000)
            await page.wait_for_load_state("networkidle")
            await asyncio.sleep(2)

        # 3. Clicar na UC
        print(f"[3] Clicando na UC {UC}...")
        await page.locator(f"text={UC}").first.click()
        try:
            await page.wait_for_selector("text=Aguarde um instante", state="visible", timeout=5_000)
            await page.wait_for_selector("text=Aguarde um instante", state="hidden", timeout=30_000)
        except Exception:
            pass
        await page.wait_for_load_state("networkidle")
        await asyncio.sleep(3)
        print(f"  OK -> {page.url}")
        await dump(page, 3, "dashboard_uc")

        # 4. Clicar em "Faturas e 2ª Via de Faturas"
        print("[4] Navegando para Faturas e 2ª Via de Faturas...")
        # Tentar pelo link do sidebar primeiro
        faturas_link = page.locator('a:has-text("Faturas e 2ª Via de Faturas")').first
        try:
            await faturas_link.click(timeout=5_000)
        except Exception:
            # Fallback: tentar pelo card
            print("  Sidebar link não clicou, tentando card...")
            faturas_card = page.locator('mat-card:has-text("Faturas e 2ª via de faturas")').first
            await faturas_card.click(timeout=5_000)

        # Aguardar carregamento
        try:
            await page.wait_for_selector("text=Aguarde um instante", state="visible", timeout=3_000)
            await page.wait_for_selector("text=Aguarde um instante", state="hidden", timeout=30_000)
        except Exception:
            pass
        await page.wait_for_load_state("networkidle")
        await asyncio.sleep(3)
        print(f"  OK -> {page.url}")
        await dump(page, 4, "tela_faturas")

        # 5. Dump detalhado da tela de faturas
        print("\n[5] Análise da tela de faturas:")

        # Textos visíveis
        all_text = await page.evaluate("""() => {
            const els = document.querySelectorAll('h1,h2,h3,h4,h5,p,span,a,button,label,td,th,mat-card-title,mat-card-subtitle');
            const seen = new Set();
            const results = [];
            for (const el of els) {
                if (el.offsetParent === null) continue;
                const t = el.textContent.trim().replace(/\\s+/g, ' ');
                if (t.length > 2 && t.length < 120 && !seen.has(t)) {
                    seen.add(t);
                    results.push({tag: el.tagName, text: t});
                }
            }
            return results.slice(0, 60);
        }""")
        for item in all_text:
            print(f"  <{item['tag']:20s}> {item['text']}")

        # Tabelas
        tables = await page.eval_on_selector_all("table", """els => els.map((t, i) => {
            const rows = Array.from(t.querySelectorAll('tr')).map(tr =>
                Array.from(tr.querySelectorAll('th, td')).map(c => c.textContent.trim().substring(0,30))
            );
            return {index: i, rows: rows.slice(0, 10)};
        })""")
        if tables:
            print(f"\n  Tabelas ({len(tables)}):")
            for t in tables:
                print(f"    Tabela {t['index']}:")
                for row in t['rows']:
                    print(f"      {row}")

        # Ícones de ação (download, print, etc)
        icons = await page.eval_on_selector_all("mat-icon, [class*='icon']", """els => els.filter(e => e.offsetParent !== null).map(e => ({
            text: e.textContent.trim(), class: e.className.substring(0,40),
            parent: e.parentElement?.tagName,
            parentClass: e.parentElement?.className?.substring(0,40),
            tooltip: e.parentElement?.getAttribute('mattooltip') || e.getAttribute('mattooltip') || '',
        })).filter(e => e.text.length > 0 && e.text.length < 30)""")
        if icons:
            print(f"\n  Ícones ({len(icons)}):")
            for ic in icons:
                print(f"    '{ic['text']}' parent=<{ic['parent']}> tooltip='{ic['tooltip']}'")

        # Botões
        buttons = await page.eval_on_selector_all("button", """els => els.filter(e => e.offsetParent !== null).map(e => ({
            text: e.textContent.trim().substring(0,50),
            tooltip: e.getAttribute('mattooltip') || '',
            disabled: e.disabled,
        }))""")
        if buttons:
            print(f"\n  Botões visíveis ({len(buttons)}):")
            for b in buttons:
                dis = " [DISABLED]" if b["disabled"] else ""
                tip = f" (tooltip: {b['tooltip']})" if b["tooltip"] else ""
                print(f"    \"{b['text']}\"{dis}{tip}")

        # Links
        links = await page.eval_on_selector_all("a", """els => els.filter(e => e.offsetParent !== null && e.textContent.trim()).map(e => ({
            text: e.textContent.trim().substring(0,50), href: e.href,
        }))""")
        if links:
            print(f"\n  Links ({len(links)}):")
            for l in links:
                print(f"    \"{l['text']}\" -> {l['href']}")

        print("\n[*] Browser aberto. Pressione Enter para fechar...")
        await asyncio.get_event_loop().run_in_executor(None, input)
        await browser.close()


if __name__ == "__main__":
    asyncio.run(main())
