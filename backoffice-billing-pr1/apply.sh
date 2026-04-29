#!/usr/bin/env bash
# apply.sh — importa o bundle do PR Billing como branch local no seu repo.
#
# Uso:
#   1. Coloque este script junto com backoffice-billing.bundle no seu
#      workspace (qualquer lugar — o script descobre o caminho via $0).
#   2. Rode dentro do seu clone do 5g-energia-fatura:
#
#        cd ~/Projetos/5g-energia-fatura
#        /caminho/para/apply.sh
#
#   3. O script faz:
#        - valida que você está dentro de um repo git do 5g-energia-fatura
#        - valida a integridade do bundle (git bundle verify)
#        - fetch do bundle pra branch local 'feat/backoffice-billing'
#        - faz checkout na branch
#        - mostra o log dos 8 commits que chegaram
#
#   4. Daí é só revisar e push:
#
#        git push origin feat/backoffice-billing
#
#      Aí o GitHub sugere "Compare & pull request" no banner.

set -euo pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
BUNDLE="${SCRIPT_DIR}/backoffice-billing.bundle"
BRANCH="feat/backoffice-billing"

# --- 1. sanity checks ------------------------------------------------

if [[ ! -f "$BUNDLE" ]]; then
    echo "ERRO: bundle não encontrado em $BUNDLE"
    echo "Coloque backoffice-billing.bundle junto deste script e rode de novo."
    exit 1
fi

if ! git rev-parse --git-dir >/dev/null 2>&1; then
    echo "ERRO: você não está dentro de um repo git."
    echo "Rode este script a partir do seu clone do 5g-energia-fatura."
    exit 1
fi

REPO_NAME=$(basename "$(git rev-parse --show-toplevel)")
if [[ "$REPO_NAME" != "5g-energia-fatura" ]]; then
    echo "AVISO: o repo atual se chama '$REPO_NAME', não '5g-energia-fatura'."
    echo "Se tem certeza que é o repo correto, pressione Enter. Caso contrário Ctrl+C."
    read -r
fi

# --- 2. valida bundle ------------------------------------------------

echo "→ Validando integridade do bundle..."
if ! git bundle verify "$BUNDLE" 2>&1 | tail -5; then
    echo "ERRO: bundle corrompido ou incompatível com este repo."
    exit 1
fi

# --- 3. fetch + checkout ---------------------------------------------

if git show-ref --verify --quiet "refs/heads/${BRANCH}"; then
    echo ""
    echo "AVISO: a branch '$BRANCH' já existe neste repo."
    echo "Escolha:"
    echo "  [r] remover e recriar do bundle (descarta branch antiga)"
    echo "  [a] abortar"
    read -r -p "> " choice
    case "$choice" in
        r|R)
            git branch -D "$BRANCH"
            ;;
        *)
            echo "Abortado."
            exit 0
            ;;
    esac
fi

echo "→ Importando bundle como branch local '$BRANCH'..."
git fetch "$BUNDLE" "main:${BRANCH}"

echo "→ Checking out '$BRANCH'..."
git checkout "$BRANCH"

# --- 4. resumo -------------------------------------------------------

echo ""
echo "✓ Pronto. Resumo do que chegou:"
echo ""
git log --oneline main..HEAD 2>/dev/null || git log --oneline -10
echo ""
echo "Arquivos alterados:"
git diff --stat main..HEAD 2>/dev/null | tail -20 || git diff --stat HEAD~8..HEAD | tail -20
echo ""
echo "Próximo passo:"
echo "  1. Revisar: leia services/backend-go/internal/app/BILLING_INTEGRATION.md"
echo "     (explica as 3 pequenas mudanças manuais no server.go existente)"
echo ""
echo "  2. Subir:   git push origin ${BRANCH}"
echo "              Depois abra PR no GitHub."
echo ""
