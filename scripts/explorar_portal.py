"""Script de exploração do portal Neoenergia para mapear seletores.
Uso: python scripts/explorar_portal.py
"""

import asyncio
from pathlib import Path

from playwright.async_api import async_playwright

SCREENSHOTS_DIR = Path("downloads/_exploracao")
URL = "https://agenciavirtual.neoenergia.com"


async def explorar():
    SCREENSHOTS_DIR.mkdir(parents=True, exist_ok=True)

    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=False)
        context = await browser.new_context(
            viewport={"width": 1280, "height": 900},
            locale="pt-BR",
        )
        page = await context.new_page()

        # 1. Página inicial
        print(f"[1] Navegando para {URL} ...")
        await page.goto(URL, wait_until="domcontentloaded")
        await page.wait_for_load_state("networkidle")
        await asyncio.sleep(3)

        await page.screenshot(path=SCREENSHOTS_DIR / "01_pagina_inicial.png", full_page=True)
        print(f"    URL atual: {page.url}")
        print(f"    Título: {await page.title()}")

        # Dump de todos os links
        links = await page.eval_on_selector_all(
            "a[href]",
            "els => els.map(e => ({text: e.textContent.trim(), href: e.href}))"
        )
        print(f"\n[2] Links encontrados ({len(links)}):")
        for link in links:
            if link["text"]:
                print(f"    {link['text'][:60]:60s} -> {link['href']}")

        # Dump de todos os botões
        buttons = await page.eval_on_selector_all(
            "button, input[type='submit'], input[type='button']",
            "els => els.map(e => ({text: e.textContent?.trim() || e.value, type: e.type, id: e.id, class: e.className}))"
        )
        print(f"\n[3] Botões encontrados ({len(buttons)}):")
        for btn in buttons:
            print(f"    [{btn['type']}] {btn['text'][:50]} id={btn['id']} class={btn['class'][:50]}")

        # Dump de todos os inputs
        inputs = await page.eval_on_selector_all(
            "input, select, textarea",
            "els => els.map(e => ({tag: e.tagName, type: e.type, name: e.name, id: e.id, placeholder: e.placeholder, class: e.className}))"
        )
        print(f"\n[4] Inputs encontrados ({len(inputs)}):")
        for inp in inputs:
            print(f"    <{inp['tag']} type={inp['type']} name={inp['name']} id={inp['id']} placeholder='{inp['placeholder']}'>")

        # Verificar se é SPA (Angular, React, etc)
        spa_check = await page.evaluate("""() => {
            return {
                angular: !!window.ng || !!document.querySelector('[ng-app]') || !!document.querySelector('[data-ng-app]'),
                react: !!document.querySelector('[data-reactroot]') || !!document.querySelector('#root'),
                vue: !!document.querySelector('[data-v-]') || !!window.__VUE__,
                hash_routing: window.location.hash.length > 1,
                url: window.location.href,
            }
        }""")
        print(f"\n[5] SPA check: {spa_check}")

        # Tentar navegar para #/login
        print(f"\n[6] Navegando para {URL}/#/login ...")
        await page.goto(f"{URL}/#/login", wait_until="domcontentloaded")
        await asyncio.sleep(3)
        await page.screenshot(path=SCREENSHOTS_DIR / "02_login_page.png", full_page=True)
        print(f"    URL atual: {page.url}")

        # Re-dump inputs na página de login
        inputs_login = await page.eval_on_selector_all(
            "input, select, textarea",
            "els => els.map(e => ({tag: e.tagName, type: e.type, name: e.name, id: e.id, placeholder: e.placeholder, class: e.className, 'aria-label': e.getAttribute('aria-label')}))"
        )
        print(f"\n[7] Inputs na página de login ({len(inputs_login)}):")
        for inp in inputs_login:
            print(f"    <{inp['tag']} type={inp['type']} name={inp['name']} id={inp['id']} placeholder='{inp['placeholder']}' aria-label='{inp.get('aria-label', '')}'>")

        # Dump HTML do body para análise
        body_html = await page.content()
        html_path = SCREENSHOTS_DIR / "login_page.html"
        html_path.write_text(body_html)
        print(f"\n[8] HTML salvo em: {html_path}")

        # Manter browser aberto para inspeção manual
        print("\n[*] Browser aberto. Pressione Enter no terminal para fechar...")
        await asyncio.get_event_loop().run_in_executor(None, input)

        await browser.close()


if __name__ == "__main__":
    asyncio.run(explorar())
