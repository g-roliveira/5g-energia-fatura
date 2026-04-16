"""Testa o fluxo completo até capturar o download do PDF."""

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


async def main():
    DIR.mkdir(parents=True, exist_ok=True)

    async with async_playwright() as p:
        browser = await p.chromium.launch(
            headless=False,
            args=["--disable-blink-features=AutomationControlled"],
        )
        ctx = await browser.new_context(
            viewport={"width": 1280, "height": 900},
            locale="pt-BR",
            user_agent=UA,
        )
        page = await ctx.new_page()
        await page.add_init_script(
            'Object.defineProperty(navigator, "webdriver", {get: () => undefined})'
        )

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

        # 3. Selecionar UC
        print(f"[3] Selecionando UC {UC}...")
        await page.locator(f"text={UC}").first.click()
        try:
            await page.wait_for_selector("text=Aguarde um instante", state="visible", timeout=5_000)
            await page.wait_for_selector("text=Aguarde um instante", state="hidden", timeout=30_000)
        except Exception:
            pass
        await page.wait_for_load_state("networkidle")
        await asyncio.sleep(3)

        # 4. Navegar para Faturas
        print("[4] Abrindo Faturas e 2ª via...")
        faturas_card = page.locator('mat-card:has-text("Faturas e 2ª via de faturas")').first
        await faturas_card.click()
        try:
            await page.wait_for_selector("text=Aguarde um instante", state="hidden", timeout=30_000)
        except Exception:
            pass
        await page.wait_for_load_state("networkidle")
        await asyncio.sleep(3)
        print(f"  URL: {page.url}")
        await page.screenshot(path=DIR / "A_lista_faturas.png", full_page=True)

        # 5. Analisar estrutura das linhas
        print("\n[5] Estrutura da lista de faturas...")

        # Container da lista
        lista = await page.eval_on_selector_all(
            "h5, [class*='list'], [class*='fatura'], [class*='row'], [class*='item']",
            """els => els.filter(e => e.offsetParent !== null && e.textContent.trim().length > 2 && e.textContent.trim().length < 200).map(e => ({
                tag: e.tagName, class: e.className.substring(0,60),
                text: e.textContent.trim().replace(/\\s+/g, ' ').substring(0, 100)
            })).slice(0,30)""",
        )
        for item in lista:
            print(f"  [{item['tag']:15s}|{item['class'][:35]}] {item['text']}")

        # 6. Checkboxes
        print("\n[6] Checkboxes...")
        checkboxes = page.locator("mat-checkbox")
        count = await checkboxes.count()
        print(f"  Total mat-checkbox: {count}")
        for i in range(count):
            cb = checkboxes.nth(i)
            label = await cb.text_content()
            aria = await cb.get_attribute("aria-checked")
            print(f"  [{i}] label='{label.strip()[:40]}' aria-checked={aria}")

        # 7. Testar clique no checkbox de uma fatura (índice 1, pulando "selecionar todas")
        print("\n[7] Selecionando primeira fatura (índice 1)...")
        if count > 1:
            await checkboxes.nth(1).click()
            await asyncio.sleep(1)
            await page.screenshot(path=DIR / "B_checkbox_marcado.png", full_page=True)

            # Verificar botão Download
            dl_btn = page.locator('button:has-text("Download")')
            disabled = await dl_btn.is_disabled()
            print(f"  Botão Download disabled: {disabled}")

            if not disabled:
                print("  Clicando Download e aguardando arquivo...")
                async with page.expect_download(timeout=60_000) as dl_info:
                    await dl_btn.click()
                download = await dl_info.value
                dest = DIR / download.suggested_filename
                await download.save_as(dest)
                print(f"  DOWNLOAD OK: {dest}")
                print(f"  Suggested filename: {download.suggested_filename}")
            else:
                print("  Download desabilitado — verificando expand row...")
                # Tentar expandir a primeira linha
                expands = page.locator(
                    "mat-icon:has-text('expand_more'), "
                    "mat-icon:has-text('keyboard_arrow_down'), "
                    "button:has(mat-icon:has-text('expand_more'))"
                )
                expand_count = await expands.count()
                print(f"  Expand icons: {expand_count}")
                if expand_count > 0:
                    await expands.first.click()
                    await asyncio.sleep(2)
                    await page.screenshot(path=DIR / "C_row_expandida.png", full_page=True)

                    # Dump do que apareceu na linha expandida
                    novos = await page.eval_on_selector_all(
                        "button, a, mat-icon",
                        """els => els.filter(e => {
                            const t = e.textContent.trim();
                            const tt = e.getAttribute('mattooltip') || '';
                            return e.offsetParent !== null && (
                                /download|imprimir|pdf|baixar|print|picture/i.test(t + ' ' + tt)
                            );
                        }).map(e => ({
                            tag: e.tagName, text: e.textContent.trim(),
                            tooltip: e.getAttribute('mattooltip') || '',
                            class: e.className.substring(0,60)
                        }))""",
                    )
                    print(f"  Após expand - botões de ação ({len(novos)}):")
                    for n in novos:
                        print(f"    <{n['tag']}> '{n['text']}' tooltip='{n['tooltip']}'")

        # 8. Explorar o que está dentro de uma linha (HTML detalhado)
        print("\n[8] HTML da primeira linha da lista de faturas...")
        first_row_html = await page.evaluate("""() => {
            // Procurar pelo container da lista de faturas
            const headers = Array.from(document.querySelectorAll('h5, h4'));
            const listaHeader = headers.find(h => /lista de fatura/i.test(h.textContent));
            if (!listaHeader) return 'Header não encontrado';
            const container = listaHeader.closest('section, div, mat-card') || listaHeader.parentElement;
            if (!container) return 'Container não encontrado';
            return container.innerHTML.substring(0, 3000);
        }""")
        print(f"  HTML container lista: {first_row_html[:1000]}")

        await browser.close()
        print("\nConcluído.")


if __name__ == "__main__":
    asyncio.run(main())
