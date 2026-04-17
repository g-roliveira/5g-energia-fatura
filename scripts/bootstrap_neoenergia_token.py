from __future__ import annotations

import argparse
import asyncio
import contextlib
import io
import json
import sys

from fatura.config import PortalConfig, TipoAcesso
from fatura.neoenergia_private_api import bootstrap_session_with_playwright


async def _run(documento: str, senha: str, uf: str, tipo_acesso: str) -> dict[str, object]:
    config = type("Cfg", (), {"portal": PortalConfig()})()
    cliente = type(
        "Cliente",
        (),
        {
            "cpf_cnpj": documento,
            "senha_portal": senha,
            "uf": uf,
            "tipo_acesso": TipoAcesso(tipo_acesso),
        },
    )()
    session = await bootstrap_session_with_playwright(config, cliente)
    return {
        "token": session.token,
        "token_ne_se": session.token_ne_se,
        "local_storage": session.local_storage,
    }


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--documento", required=True)
    parser.add_argument("--senha", required=True)
    parser.add_argument("--uf", default="BA")
    parser.add_argument("--tipo-acesso", default="normal")
    args = parser.parse_args()

    silent_stdout = io.StringIO()
    with contextlib.redirect_stdout(silent_stdout):
        payload = asyncio.run(_run(args.documento, args.senha, args.uf, args.tipo_acesso))
    sys.stdout.write(json.dumps(payload, ensure_ascii=False))


if __name__ == "__main__":
    main()
