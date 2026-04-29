#!/usr/bin/env bash
set -euo pipefail

# Cores
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}🚀 Iniciando ambiente de desenvolvimento 5G Energia Fatura...${NC}"

# 1. Subir infraestrutura (Postgres + Extractor)
echo -e "${YELLOW}📦 Subindo Docker Compose (Postgres + Extractor)...${NC}"
docker compose up -d

# Aguardar Postgres ficar saudável
echo -e "${YELLOW}⏳ Aguardando Postgres...${NC}"
until docker exec azi-billing-postgres pg_isready -U azi -d azi_billing >/dev/null 2>&1; do
  sleep 1
done
echo -e "${GREEN}✅ Postgres pronto${NC}"

# 2. Verificar se .env.local existe no web
if [ ! -f "web/.env.local" ]; then
  echo -e "${YELLOW}⚠️  web/.env.local não encontrado. Copiando de .env.example...${NC}"
  if [ -f "web/.env.example" ]; then
    cp web/.env.example web/.env.local
  else
    echo -e "${YELLOW}⚠️  Criando web/.env.local básico...${NC}"
    cat > web/.env.local <<EOF
DATABASE_URL="postgresql://azi:azi@localhost:5434/azi_billing"
BACKEND_GO_URL="http://localhost:8080"
EOF
  fi
fi

# 3. Iniciar Go backend em background
echo -e "${YELLOW}🔧 Iniciando Go backend...${NC}"
cd services/backend-go
go run ./cmd/api &
GO_PID=$!
cd ../..

echo -e "${GREEN}✅ Go backend iniciado (PID: $GO_PID)${NC}"

# 4. Iniciar Next.js
echo -e "${YELLOW}🌐 Iniciando Next.js...${NC}"
cd web
pnpm dev --turbopack &
WEB_PID=$!
cd ..

echo -e "${GREEN}✅ Next.js iniciado (PID: $WEB_PID)${NC}"

echo ""
echo -e "${GREEN}═══════════════════════════════════════════════${NC}"
echo -e "${GREEN}  🎉 Ambiente de dev pronto!${NC}"
echo -e "${GREEN}  • Frontend:  http://localhost:3000${NC}"
echo -e "${GREEN}  • Backend:   http://localhost:8080${NC}"
echo -e "${GREEN}  • Postgres:  localhost:5434${NC}"
echo -e "${GREEN}  • Extractor: http://localhost:8090${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════${NC}"
echo ""
echo "Pressione Ctrl+C para encerrar todos os serviços"

# Trap para matar os processos ao sair
cleanup() {
  echo ""
  echo -e "${YELLOW}🛑 Encerrando serviços...${NC}"
  kill $WEB_PID $GO_PID 2>/dev/null || true
  docker compose down
  echo -e "${GREEN}✅ Tudo encerrado${NC}"
}
trap cleanup INT TERM EXIT

wait
