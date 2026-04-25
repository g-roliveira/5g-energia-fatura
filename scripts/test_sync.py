#!/usr/bin/env python3
"""
Script de teste real de sync com a Neoenergia.

Lê credenciais do config.yaml na raiz do projeto e executa o fluxo completo
via API local (localhost:8080).

Uso:
    python3 scripts/test_sync.py

Pré-requisitos:
    - API rodando: cd services/backend-go && go run ./cmd/api
    - Postgres rodando em localhost:5434
    - config.yaml na raiz do projeto com credenciais
"""

import json
import re
import sys
from pathlib import Path

import requests
import yaml

BASE_URL = "http://localhost:8091"


def load_config():
    """Lê config.yaml da raiz do projeto."""
    paths = [
        Path("config.yaml"),
        Path("../config.yaml"),
        Path("../../config.yaml"),
    ]
    for p in paths:
        if p.exists():
            with open(p) as f:
                return yaml.safe_load(f)
    print("ERRO: config.yaml não encontrado")
    sys.exit(1)


def clean_doc(doc):
    """Remove pontuação do CPF/CNPJ."""
    return re.sub(r"[^0-9]", "", doc)


def health_check():
    """Verifica se a API está respondendo."""
    r = requests.get(f"{BASE_URL}/healthz", timeout=5)
    r.raise_for_status()
    print("✓ API health check OK")


def create_credential(cliente):
    """Cria credencial na API (envia plain text, backend criptografa)."""
    payload = {
        "label": cliente["nome"],
        "documento": clean_doc(cliente["cpf_cnpj"]),
        "senha": cliente["senha_portal"],
        "uf": cliente["uf"],
        "tipo_acesso": cliente.get("tipo_acesso", "normal"),
    }

    r = requests.post(f"{BASE_URL}/v1/credentials", json=payload, timeout=10)
    if r.status_code != 201:
        print(f"ERRO ao criar credencial: {r.status_code} - {r.text}")
        sys.exit(1)

    data = r.json()
    cred_id = data.get("id", data.get("credential_id", ""))
    print(f"✓ Credential criada: {cred_id}")
    return cred_id


def create_session(cred_id):
    """Cria sessão (login real na Neoenergia)."""
    print("→ Criando sessão (login real na Neoenergia)...")
    r = requests.post(f"{BASE_URL}/v1/credentials/{cred_id}/session", timeout=60)
    if r.status_code != 200:
        print(f"ERRO: {r.status_code} - {r.text}")
        print("\nNota: se falhou com erro de login, verifique o CPF/senha no config.yaml")
        sys.exit(1)
    data = r.json()
    print("✓ Sessão criada")
    return data


def sync_uc(cred_id, uc, doc):
    """Executa sync da UC."""
    print(f"→ Executando sync para UC {uc}...")
    payload = {
        "credential_id": cred_id,
        "uc": uc,
        "documento": doc,
    }
    r = requests.post(f"{BASE_URL}/v1/sync/uc", json=payload, timeout=120)
    if r.status_code != 200:
        print(f"ERRO no sync: {r.status_code} - {r.text}")
        sys.exit(1)

    data = r.json()
    print("✓ Sync executado:")
    print(json.dumps(data, indent=2, ensure_ascii=False))
    return data


def list_consumer_units():
    """Lista consumer units."""
    print("→ Listando consumer units...")
    r = requests.get(f"{BASE_URL}/v1/consumer-units", timeout=10)
    r.raise_for_status()
    data = r.json()
    items = data.get("items", [])
    print(f"✓ {len(items)} consumer unit(s) encontrada(s)")
    return items


def list_invoices(uc):
    """Lista invoices da UC."""
    print(f"→ Listando invoices para UC {uc}...")
    r = requests.get(f"{BASE_URL}/v1/consumer-units/{uc}/invoices", timeout=10)
    r.raise_for_status()
    data = r.json()
    items = data.get("items", [])
    print(f"✓ {len(items)} invoice(s) encontrada(s)")
    for inv in items:
        print(f"  - Fatura {inv.get('numero_fatura', '?')} | Mês: {inv.get('mes_referencia', '?')} | Valor: {inv.get('valor_total', '?')}")
    return items


def main():
    print("=== Teste de Sync Real ===\n")

    cfg = load_config()
    clientes = cfg.get("clientes", [])
    if not clientes:
        print("ERRO: nenhum cliente no config.yaml")
        sys.exit(1)

    cliente = clientes[0]
    print(f"Cliente: {cliente['nome']}")
    print(f"UC: {cliente['uc']}")
    print(f"UF: {cliente['uf']}")
    print()

    health_check()
    cred_id = create_credential(cliente)
    create_session(cred_id)
    sync_uc(cred_id, cliente["uc"], clean_doc(cliente["cpf_cnpj"]))
    list_consumer_units()
    list_invoices(cliente["uc"])

    print("\n=== Teste concluído com sucesso ===")


if __name__ == "__main__":
    main()
