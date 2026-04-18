Quero que você atue como engenheiro de software sênior e product engineer especialista em Next.js, PostgreSQL, shadcn/ui e sistemas SaaS.

## Objetivo

Construir o módulo de Clientes + UCs + Sincronização/Faturas, integrado ao backend Go já existente.

Este frontend será implementado por outro agente, então preciso de uma implementação madura, completa e consistente com produto real.

## Contexto de arquitetura (obrigatório)

Existem dois backends com responsabilidades diferentes:

1. Backend principal do front (Next.js + PostgreSQL local)
- fonte de verdade para cadastro de clientes, endereços, UCs e dados comerciais.

2. Backend Go de integração com concessionária
- credenciais criptografadas da concessionária
- sessão/token
- sincronização de UC
- persistência de invoices
- itens de fatura
- extração documental
- auditoria de sync

Regra:
- O frontend web (browser) NÃO chama backend Go direto.
- O frontend chama somente rotas internas Next (`app/api/...`).
- As rotas internas Next fazem server-to-server com backend Go.

## Documentação e fontes oficiais para consulta

Use obrigatoriamente estas referências no projeto:

- `/home/gustavo/Projetos/5g-energia-fatura/docs/backend-go-api.md`
- `/home/gustavo/Projetos/5g-energia-fatura/docs/architecture/backend-monorepo.md`
- `/home/gustavo/Projetos/5g-energia-fatura/services/backend-go/internal/app/openapi_routes.go`
- Backend Go runtime:
  - `GET /docs`
  - `GET /openapi.json`

## Base visual já existente

O projeto já possui template/layout/admin UI em `web/` com:
- sidebar/layout base
- componentes shadcn/ui
- estrutura visual pronta

Portanto:
- NÃO recriar layout base do zero.
- Reaproveitar componentes e estrutura de navegação existentes.
- Criar somente componentes novos necessários.

## Stack obrigatória

- Next.js App Router
- TypeScript
- shadcn/ui
- React Hook Form
- Zod
- TanStack Table
- PostgreSQL (sem Supabase)

## Estratégia de persistência

Persistência do módulo de cadastro no banco local PostgreSQL do front:
- clientes
- endereços
- unidades consumidoras
- dados comerciais/contratuais
- vínculo com credencial de integração

Integração externa via backend Go:
- sessão de integração
- sync runs
- invoices
- detalhe de invoice e auditoria

## Modelagem mínima no PostgreSQL local

### Cliente
- `id`
- `tipo_pessoa` (`PF` | `PJ`)
- `nome_razao_social`
- `nome_fantasia`
- `cpf_cnpj` (único)
- `email`
- `telefone`
- `status` (`ativo` | `inativo` | `prospecto`)
- `tipo_cliente` (`residencial` | `condominio` | `empresa` | `imobiliaria` | `outro`)
- `observacoes`
- `created_at`
- `updated_at`
- `archived_at`

### Endereço do cliente
- `id`
- `client_id`
- `cep`
- `logradouro`
- `numero`
- `complemento`
- `bairro`
- `cidade`
- `uf`
- `created_at`
- `updated_at`

### Unidade consumidora (local)
- `id`
- `client_id`
- `uc_code` (único)
- `distribuidora`
- `apelido`
- `classe_consumo`
- `endereco_unidade`
- `cidade`
- `uf`
- `ativa`
- `credential_id` (referência local da credencial de integração)
- `created_at`
- `updated_at`

### Dados comerciais
- `id`
- `client_id`
- `tipo_contrato`
- `data_inicio`
- `data_fim`
- `status_contrato`
- `observacoes_comerciais`
- `created_at`
- `updated_at`

### Credencial de integração (local)
- `id`
- `client_id` (ou escopo por UC, justificar)
- `label`
- `documento_masked`
- `uf`
- `tipo_acesso`
- `go_credential_id` (id retornado por `POST /v1/credentials`)
- `created_at`
- `updated_at`

Observação:
- Senha da concessionária não deve ficar em texto puro na UI.
- Se persistir localmente, deve ser criptografada no backend local.
- Nunca trafegar senha em query string.

## Rotas internas Next (BFF) obrigatórias

Essas rotas serão consumidas pela UI:

### CRUD local
- `POST /api/clients`
- `GET /api/clients`
- `GET /api/clients/:id`
- `PATCH /api/clients/:id`
- `POST /api/clients/:id/archive`
- `POST /api/clients/:id/ucs`
- `PATCH /api/ucs/:id`

### Integração com backend Go (server-to-server)
- `POST /api/integration/credentials` -> `POST /v1/credentials`
- `POST /api/integration/credentials/:id/session` -> `POST /v1/credentials/{id}/session`
- `POST /api/integration/ucs/:uc/sync` -> `POST /v1/consumer-units/{uc}/sync`
- `GET /api/integration/ucs` -> `GET /v1/consumer-units`
- `GET /api/integration/ucs/:uc` -> `GET /v1/consumer-units/{uc}`
- `GET /api/integration/ucs/:uc/invoices` -> `GET /v1/consumer-units/{uc}/invoices`
- `GET /api/integration/invoices/:id` -> `GET /v1/invoices/{id}`
- `GET /api/integration/sync-runs/:id` -> `GET /v1/sync-runs/{id}`

## Telas obrigatórias

### Tela 1: Lista de clientes
- tabela com busca, filtros, paginação, ordenação
- filtros por status e tipo de cliente
- colunas:
  - nome/razão social
  - CPF/CNPJ
  - tipo cliente
  - status
  - quantidade de UCs
  - cidade/UF
  - data de cadastro
  - ações
- botão `Novo cliente`
- botão `Importar CSV` desabilitado com badge `Em breve`

### Tela 2: Cadastro/edição de cliente
- formulário por seções:
  - dados principais
  - contato
  - endereço
  - unidades consumidoras
  - observações comerciais
- reutilizar mesmo form para criar e editar
- validação com Zod + React Hook Form

### Tela 3: Detalhe do cliente
- header com dados principais
- seções/abas:
  - dados cadastrais
  - endereço
  - UCs
  - observações comerciais
  - integração
- ações rápidas:
  - editar cliente
  - arquivar
  - adicionar UC
  - sincronizar UC

### Tela 4: Painel de UCs do cliente
- lista de UCs vinculadas
- status local da UC
- status da última sincronização (se houver)
- ação `Sincronizar agora` por UC
- ação para abrir invoices da UC

### Tela 5: Invoices por UC
- origem dos dados: backend Go persistido (não consultar Neoenergia direto)
- filtros por status/período
- colunas:
  - número da fatura
  - referência
  - valor
  - vencimento
  - status
  - completude
  - atualizado em
- ação de detalhe

### Tela 6: Detalhe da invoice
- exibir:
  - `billing_record`
  - `document_record`
  - itens da fatura
  - completude (`complete|partial|failed`)
  - origem dos campos quando disponível

### Tela 7: Auditoria de sincronização
- detalhe de `sync_run`
- mostrar:
  - status
  - horário de execução
  - erro (quando existir)
  - `raw_response` resumido/colapsável
- CTA para reprocessar

## Requisitos de UX

- padrão SaaS moderno
- estados vazios
- skeleton loading
- feedback de sucesso/erro
- confirmação para ações destrutivas
- responsivo
- acessível

## Requisitos técnicos

1. Implementar no `web/` reaproveitando a base visual.
2. Definir modelagem e migrations SQL PostgreSQL.
3. Implementar camada de acesso a dados desacoplada (Prisma ou Drizzle, justificar).
4. Implementar seeds de desenvolvimento.
5. Implementar BFF no Next para integração com backend Go.
6. Tratar timeout/erro nas integrações.
7. Garantir que segredos não vazem para client.

## Critérios de aceite

1. CRUD de cliente e UC funcional com persistência local.
2. Fluxo de sincronização por UC funcional via backend Go.
3. Leitura de invoices e detalhe de invoice a partir de endpoints persistidos.
4. Tela de auditoria de sync funcional.
5. Nenhuma chamada direta do browser para backend Go.
6. Código componentizado e consistente com layout existente.
7. Documentação curta de:
   - decisões de arquitetura
   - modelo de dados
   - integração BFF
   - pontos preparados para evolução.

## Entregáveis

1. proposta de arquitetura final do módulo
2. modelagem de dados + SQL
3. estratégia de acesso a dados com justificativa
4. estrutura de pastas
5. lista de componentes reutilizados e novos
6. páginas implementadas
7. rotas BFF implementadas
8. seed/mock de desenvolvimento
9. resumo de integração com backend Go
10. checklist de validação funcional
