# PRD Front-end — Backoffice Azi Dourado

**Versão:** 1.0
**Data:** abril/2026
**Escopo:** aplicação web interna que a equipe da Azi Dourado (founder + operador + futuros reviewers) usa para gerenciar clientes, UCs, competências mensais e faturamento de energia compartilhada.
**Fora de escopo neste PRD:** portal do cliente final (vira PRD separado depois). Mobile native. Login social.

Este documento é a especificação completa. Um agente de front-end (Claude Code ou equivalente) deve conseguir implementar todas as telas listadas aqui sem precisar inferir nada sobre o backend.

---

## Sumário

1. Arquitetura e decisões fundamentais
2. Integrações com o backend
3. Padrões globais de UX
4. Estrutura de rotas (App Router)
5. Componentes compartilhados
6. Inventário completo de telas
7. Fluxos end-to-end
8. Sistema de notificações
9. O que NÃO está no MVP (com justificativa)

---

## 1. Arquitetura e decisões fundamentais

### 1.1 Stack

| Camada | Tecnologia | Por quê |
|---|---|---|
| Framework | Next.js 14+ (App Router) | Já escolhido no `prompt-front.md` |
| UI | shadcn/ui + Tailwind | Base visual já pronta em `web/` |
| Formulários | React Hook Form + Zod | Já escolhido; validação tipada |
| Tabelas | TanStack Table v8 | Ordenação/filtro/paginação no cliente |
| Server state | TanStack Query v5 | Cache, invalidação, optimistic updates, SSE-friendly |
| Datas | date-fns (locale pt-BR) | Formatação BR sem surpresa |
| Decimal | decimal.js | Evita erros de ponto flutuante em R$ |
| Charts (mínimos) | recharts | Se necessário para dashboard simples |
| Toasts | sonner (já no shadcn) | Feedback de ação |
| Icons | lucide-react | Já no shadcn |

### 1.2 Dois backends, duas fontes de verdade

O front conversa com **dois backends distintos**, que o Next.js App Router unifica através do BFF em `app/api/*`.

1. **Backend Go** (`services/backend-go`) — porta 8080.
   - Responsável por integração com a concessionária (Neoenergia/Coelba).
   - Dados: credenciais, sessões, UCs da API, faturas originais (PDF + JSON), sync runs.
   - Persiste em SQLite interno.
   - Expõe: `/v1/credentials`, `/v1/consumer-units`, `/v1/invoices`, `/v1/sync/uc`, etc.
   - Expõe também o **módulo billing** (novo): `/v1/billing/contracts`, `/v1/billing/calculations`, `/v1/billing/cycles` (no futuro).

2. **Postgres do backoffice** — banco direto do Next.
   - Responsável por cadastro (clientes, endereços, UCs cadastrais, credenciais) e **faturamento** (contracts, cycles, calculations, documentos gerados).
   - Schema `core` = cadastro. Schema `billing` = faturamento.
   - Acessado exclusivamente pelo BFF do Next (server actions e route handlers). **O browser NUNCA fala direto com Postgres.**

### 1.3 Regra de ouro

> **Browser nunca chama backend Go direto. Browser nunca chama Postgres direto.**
>
> Sempre através de `app/api/*` (BFF do Next). Isso garante:
> - Credenciais/segredos ficam no server.
> - Políticas de RBAC aplicadas num lugar só.
> - Fácil trocar backend Go por outro amanhã sem mexer em client code.

### 1.4 Papéis de usuário

Três roles no `core.app_user`:

- **admin** — tudo: criar usuário, fechar competência, alterar contrato.
- **operator** — dia a dia: criar cliente, vincular UC, disparar sync, revisar cálculo, gerar PDF. **Não** aprova cálculo, **não** fecha competência, **não** altera contrato.
- **reviewer** — só leitura + aprovar cálculo. **Não** cria, **não** edita.

O PRD marca em cada tela "quem vê" e "quem pode agir".

---

## 2. Integrações com o backend

### 2.1 Autenticação

**Entre o browser e o Next (BFF):** sessão baseada em cookie HttpOnly (iron-session ou next-auth em modo credentials). Tabela `core.app_user`. Senha com bcrypt.

**Entre o Next e o backend Go:** por ora, header `X-API-Key` configurado em env `BACKEND_GO_API_KEY`. Todo request server-to-server do Next pro Go inclui esse header. O browser nunca vê essa chave.

### 2.2 Endpoints do backend-go que o BFF consome

Retiro direto do `openapi.json` do backend-go + das rotas do módulo billing que o PR 1 acabou de adicionar.

**Credenciais e sessão Neoenergia:**
- `POST /v1/credentials` — cria credencial criptografada
- `POST /v1/credentials/{id}/session` — força criação de sessão (bootstrap Playwright)
- `GET  /v1/credentials/{id}/discover` — descobre UCs da credencial

**UCs e sincronização:**
- `GET  /v1/consumer-units` — lista UCs persistidas (paginação: `limit`, `status`)
- `GET  /v1/consumer-units/{uc}` — detalhe + última fatura + último sync
- `POST /v1/consumer-units/{uc}/sync` — dispara sync (body: `{credential_id, include_pdf, include_extraction}`)
- `GET  /v1/consumer-units/{uc}/invoices` — lista faturas da UC
- `GET  /v1/consumer-units/{uc}/latest-invoice` — última fatura

**Faturas:**
- `GET  /v1/invoices/{id}` — detalhe da fatura Coelba (inclui `billing_record` + `document_record` + PDF em base64 em `invoice_documents`)

**Sync runs:**
- `GET  /v1/sync-runs/{id}` — detalhe de uma execução (com evidências se falhou)

**Módulo billing (do PR 1):**
- `POST /v1/billing/contracts` — cria contrato (nova versão fecha anterior)
- `GET  /v1/billing/contracts/{id}` — detalhe
- `GET  /v1/billing/consumer-units/{uc_id}/active-contract` — contrato vigente de uma UC
- `GET  /v1/billing/calculations/{id}` — detalhe de um cálculo (com 3 snapshots)

**Módulo billing (próximos PRs — o front já pode desenhar contra eles):**
- `POST /v1/billing/cycles` — abre competência
- `GET  /v1/billing/cycles` — lista competências
- `GET  /v1/billing/cycles/{id}` — detalhe + métricas
- `GET  /v1/billing/cycles/{id}/rows` — tabela principal: linha por UC com estado
- `POST /v1/billing/cycles/{id}/sync-all` — dispara sync em massa
- `POST /v1/billing/cycles/{id}/close` — fecha
- `POST /v1/billing/calculations/{id}/adjust` — aplica ajuste (cria nova version)
- `POST /v1/billing/calculations/{id}/recalculate` — reroda motor
- `POST /v1/billing/calculations/{id}/approve` — aprova
- `POST /v1/billing/calculations/{id}/generate-pdf` — gera PDF do cliente
- `GET  /v1/billing/calculations/{id}/pdf` — URL assinada do PDF
- `POST /v1/billing/calculations/bulk` — ações em massa
- `GET  /v1/billing/events/cycles/{id}` — **SSE** com progresso do ciclo

### 2.3 Tabelas do Postgres que o BFF lê direto (via Drizzle ORM)

O front pode ler direto de `core.*` para telas puramente cadastrais (sem passar pelo backend-go).

- `core.customer`, `core.address`, `core.consumer_unit`, `core.credential_link`, `core.app_user`

Para `billing.*`, preferir sempre chamar o backend-go — ele valida invariantes que o ORM não valida.

### 2.4 Padrão do BFF

Todas as rotas em `app/api/` seguem o mesmo padrão:

```ts
// app/api/billing/contracts/route.ts
import { goClient } from '@/lib/backend-go-client'
import { requireSession } from '@/lib/auth'

export async function POST(req: Request) {
  const session = await requireSession(req)         // valida cookie
  const body = await req.json()
  const result = await goClient.post(
    '/v1/billing/contracts',
    { ...body, created_by: session.userId }
  )
  return Response.json(result.data, { status: result.status })
}
```

O `goClient` é um cliente axios/fetch com:
- `baseURL` = `process.env.BACKEND_GO_URL`
- `X-API-Key` = `process.env.BACKEND_GO_API_KEY`
- timeout de 30s (syncs podem demorar)
- retry automático em 5xx (3 tentativas com backoff)

### 2.5 Tratamento de erro padronizado

O backend-go retorna erros como `{"error": "código_do_erro"}` com HTTP status semântico. O BFF normaliza pra:

```ts
type ApiError = {
  code: string          // 'not_found', 'validation_error', 'upstream_unavailable', etc
  message: string       // msg legível pt-BR
  field?: string        // pra erros de validação
}
```

Tabela de mapeamento fica em `lib/errors.ts`. Toast global exibe a `message`; forms usam `field` pra destacar campo.

---

## 3. Padrões globais de UX

Estes são **invariantes**. Toda tela obedece.

### 3.1 Estados obrigatórios

Toda tela que carrega dados tem 4 estados desenhados:

1. **Loading** — skeleton da tabela/card/form. NUNCA spinner no meio da tela vazia.
2. **Empty** — mensagem + CTA ("Nenhum cliente cadastrado ainda. **[+ Novo cliente]**")
3. **Error** — mensagem + retry. Se for 401, redireciona pra login. Se for 5xx, "Problema no servidor — tente de novo em 1 min".
4. **Success** — os dados.

### 3.2 Formatação BR

- **Moeda:** `R$ 1.234,56` (ponto milhar, vírgula decimal, sempre 2 casas).
- **Data curta:** `06/05/2026` (DD/MM/YYYY).
- **Data + hora:** `06/05/2026 14:22` (HH:MM, 24h).
- **Data relativa:** "há 2 horas", "ontem às 14:22" (date-fns locale pt).
- **Mês/ano:** `abril/2026` (minúsculo, nome completo) ou `04/2026` em tabelas compactas.
- **Números grandes:** `1.234.567` (só milhar).
- **Porcentagem:** `15%` ou `15,5%` (nunca com mais de 1 casa salvo necessidade explícita).
- **UC:** sempre com zeros à esquerda como vem da Coelba (`007098175908`, não `7098175908`).
- **CPF:** `123.456.789-01`. **CNPJ:** `12.345.678/0001-99`.

Componente `<BRL value={...} />`, `<DateShort value={...} />`, etc. em `components/format/`.

### 3.3 Ações destrutivas

- **Arquivar cliente**, **fechar competência**, **cancelar cálculo aprovado**: sempre modal de confirmação com texto digitado.
- **Deletar**: não existe. Sempre arquivar ou marcar como inativo.
- Modal padrão: "Digite o nome do cliente para confirmar" — evita clique acidental.

### 3.4 Feedback de ação

- **Ação bem-sucedida:** toast verde (sonner) no canto superior direito. Duração 4s. Texto curto: "Cliente criado", "Fatura sincronizada".
- **Ação com erro:** toast vermelho. Duração 7s. Texto: mensagem do backend + botão "Detalhes" que abre modal com stack do erro (só pra admin; operator vê mensagem curta).
- **Ação demorada (sync):** toast com spinner persistente + link "Ver progresso" que leva pra tela do ciclo/job.

### 3.5 Atalhos de teclado

- `/` — foca busca global
- `n` — "novo" (contextual: na tela de clientes abre modal de novo cliente)
- `esc` — fecha modal/drawer
- `cmd+k` — paleta de comandos (opcional, nice-to-have)

### 3.6 Permissões

Toda rota sensível verifica role em `middleware.ts`. Botões de ação escondem quando o role não tem permissão. **Não basta desabilitar** — esconde, pra não poluir UI de reviewer.

Componente `<RoleGate allow={['admin','operator']}>{children}</RoleGate>` envolve blocos.

### 3.7 Responsividade

Alvo: desktop 1280px+. Deve funcionar razoavelmente em tablet 768px+ (sidebar colapsa). Mobile é out-of-scope (backoffice, não app do cliente).

---

## 4. Estrutura de rotas (App Router)

```
web/
├── app/
│   ├── (auth)/
│   │   └── login/
│   │       └── page.tsx              → Tela 1
│   │
│   ├── (app)/                         ← layout com sidebar + topbar
│   │   ├── layout.tsx                 → verifica sessão, injeta sidebar
│   │   ├── page.tsx                   → Tela 2 (Dashboard home)
│   │   │
│   │   ├── customers/
│   │   │   ├── page.tsx               → Tela 3 (Lista de clientes)
│   │   │   ├── new/
│   │   │   │   └── page.tsx           → Tela 4 (Novo cliente — wizard ou form)
│   │   │   └── [id]/
│   │   │       ├── page.tsx           → Tela 5 (Detalhe do cliente)
│   │   │       └── edit/
│   │   │           └── page.tsx       → Tela 6 (Editar cliente)
│   │   │
│   │   ├── consumer-units/
│   │   │   ├── page.tsx               → Tela 7 (Lista de UCs)
│   │   │   └── [uc]/
│   │   │       ├── page.tsx           → Tela 8 (Detalhe da UC)
│   │   │       ├── invoices/
│   │   │       │   ├── page.tsx       → Tela 9 (Faturas da UC)
│   │   │       │   └── [invoiceId]/
│   │   │       │       └── page.tsx   → Tela 10 (Detalhe da fatura Coelba)
│   │   │       └── contract/
│   │   │           └── page.tsx       → Tela 11 (Contrato vigente + histórico)
│   │   │
│   │   ├── cycles/
│   │   │   ├── page.tsx               → Tela 12 (Lista de competências)
│   │   │   ├── new/
│   │   │   │   └── page.tsx           → Tela 13 (Abrir competência)
│   │   │   └── [id]/
│   │   │       ├── page.tsx           → Tela 14 (Dashboard da competência)
│   │   │       └── rows/[rowId]/
│   │   │           └── page.tsx       → Tela 15 (Detalhe do cálculo)
│   │   │
│   │   ├── sync-jobs/
│   │   │   ├── page.tsx               → Tela 16 (Fila de jobs)
│   │   │   └── [id]/
│   │   │       └── page.tsx           → Tela 17 (Detalhe do job + evidências)
│   │   │
│   │   ├── settings/
│   │   │   ├── page.tsx               → Tela 18 (Config geral)
│   │   │   ├── users/
│   │   │   │   └── page.tsx           → Tela 19 (Usuários do backoffice)
│   │   │   └── integrations/
│   │   │       └── page.tsx           → Tela 20 (Credenciais Neoenergia)
│   │   │
│   │   └── notifications/
│   │       └── page.tsx               → Tela 21 (Inbox de notificações)
│   │
│   ├── api/                            ← BFF (server-side)
│   │   ├── auth/[...]/route.ts
│   │   ├── customers/[...]/route.ts
│   │   ├── consumer-units/[...]/route.ts
│   │   ├── contracts/[...]/route.ts
│   │   ├── cycles/[...]/route.ts
│   │   ├── calculations/[...]/route.ts
│   │   ├── sync-jobs/[...]/route.ts
│   │   ├── events/cycles/[id]/route.ts ← SSE proxy
│   │   └── users/[...]/route.ts
│   │
│   └── globals.css
│
├── components/
│   ├── ui/                             ← shadcn primitives (já existentes)
│   ├── format/                         ← BRL, DateShort, UCCode, etc
│   ├── tables/                         ← DataTable base
│   ├── forms/                          ← AddressForm, etc
│   ├── layout/                         ← Sidebar, Topbar, NotificationsBell
│   ├── status/                         ← StatusBadge, SyncStatusDot
│   ├── billing/                        ← ContractForm, CalculationBreakdown
│   └── cycle/                          ← CycleProgressBar, CycleRowsTable
│
├── lib/
│   ├── backend-go-client.ts
│   ├── db/                             ← Drizzle schemas + client
│   ├── auth.ts
│   ├── errors.ts
│   └── format.ts
│
└── drizzle/                            ← migrations do front (se usar separadas)
```

**Total: 21 telas.** Nenhuma é opcional no MVP — todas cobrem um fluxo real. A próxima seção detalha uma por uma.

---

## 5. Componentes compartilhados

Antes de atacar as telas, esta é a lista de componentes que aparecem em várias. Implementar esses primeiro economiza tempo.

### 5.1 `<DataTable>` (TanStack Table)

Base universal de tabela. Colunas configuráveis, filtros, ordenação, paginação (server ou client), seleção múltipla opcional, row actions, estado vazio, skeleton.

Props: `columns`, `data`, `isLoading`, `onRowClick`, `selectable`, `bulkActions`, `pagination`.

### 5.2 `<StatusBadge>`

Badge colorido com estado tipado. Cores fixas:

- `draft` — cinza
- `syncing` / `processing` / `running` — azul (com pulse opcional)
- `review` / `needs_review` — amarelo
- `approved` — verde
- `closed` — verde escuro / slate
- `error` / `failed` — vermelho
- `superseded` — cinza claro (riscado)
- `archived` — cinza bem claro

Uso: `<StatusBadge status={cycle.status} />`.

### 5.3 `<UCCode>`

Formata código da UC com fonte monoespaçada e zeros à esquerda preservados. Copiável ao clicar.

### 5.4 `<BRL>`, `<DateShort>`, `<DateTime>`, `<RelativeTime>`

Formatadores BR. Já mencionado na §3.2.

### 5.5 `<ConfirmDestructive>`

Modal de confirmação com input "digite X pra confirmar". Props: `title`, `description`, `confirmText`, `confirmWord`, `onConfirm`.

### 5.6 `<EmptyState>`

Composição de ícone + título + descrição + CTA. Uso em listas vazias.

### 5.7 `<SourceIndicator>`

Mostra de onde o dado veio (API, pymupdf, mistral, manual) com ícone + tooltip. Vem do `source_map` do `billing_record` — importante no detalhe da fatura.

### 5.8 `<ConfidenceBar>`

Barra horizontal 0–100% colorida: vermelho <0.5, amarelo 0.5–0.85, verde >0.85. Usada em detalhe de fatura pra mostrar `extractor_confidence`.

### 5.9 `<NotificationsBell>` (topbar)

Sino com badge numérico. Clica → dropdown com últimas 10 notificações. Link "Ver todas" → `/notifications`. Detalhes na §8.

### 5.10 `<CycleProgressBar>`

Barra compacta com percentual + breakdown (synced / calculated / approved). Escuta SSE se o ciclo está em processamento.

### 5.11 `<JsonViewer>`

Viewer read-only com colapsar/expandir. Usado pra mostrar `contract_snapshot_json`, `inputs_snapshot_json`, `result_snapshot_json`, `raw_response` de sync.

### 5.12 `<AddressForm>`

Bloco reutilizável CEP + busca ViaCEP + logradouro/número/complemento/bairro/cidade/UF.

### 5.13 `<PdfDownloadButton>`

Botão com ícone de download. Props: `url` (opcional), `onFetch` (opcional — chama BFF, obtém URL assinada, abre em aba nova). Estados: idle, loading, error.

### 5.14 `<Drawer>` (shadcn sheet)

Painel lateral direito pra detalhes secundários sem navegar. Usado no detalhe do ciclo pra ver linha sem perder contexto.

---

## 6. Inventário completo de telas

### Template por tela

Cada tela abaixo segue este template:

- **Rota**
- **Objetivo** — a decisão/ação que o usuário toma
- **Quem vê** / **Quem age**
- **Dados carregados** — endpoints exatos
- **Layout** — regiões
- **Estados** — loading/empty/error/success
- **Ações** — cada botão/link, o que faz
- **Integrações críticas** — chamadas explícitas
- **Componentes usados**

---

### Tela 1 — Login

- **Rota:** `/login`
- **Objetivo:** autenticar operador do backoffice.
- **Quem vê:** qualquer visitante não logado.
- **Quem age:** qualquer um com credencial válida.
- **Dados carregados:** nenhum (form estático).

**Layout:**
Tela cheia, centralizada, card no meio. Logo da Azi Dourado no topo do card. Título "Entrar". Campos: e-mail, senha. Botão primário "Entrar". Link pequeno "Esqueci minha senha" — **desabilitado com tooltip "Em breve"** no MVP (operador pede reset pro admin direto).

**Estados:**
- Loading: botão com spinner, campos desabilitados.
- Error: mensagem vermelha acima do form ("E-mail ou senha incorretos"). Mantém e-mail preenchido.
- Success: redirect pra `/` (Tela 2).

**Ações:**
- `[Entrar]` → POST `/api/auth/login` com `{email, password}` → se 200, set cookie de sessão, redirect `/`.
- Enter no campo senha = mesmo que clicar Entrar.

**Integrações:**
- `POST /api/auth/login` — BFF valida contra `core.app_user`, hash bcrypt, cria sessão iron-session.

**Componentes:**
- Form, Input (shadcn), Button, Card.

**Regras:**
- 5 tentativas erradas em 5 minutos bloqueia o e-mail por 15 min (rate limit no BFF).
- Senha nunca é logada, nem em dev.

---

### Tela 2 — Dashboard home

- **Rota:** `/`
- **Objetivo:** dar ao usuário o panorama do que precisa de atenção hoje. É a primeira tela ao logar.
- **Quem vê:** todos os roles logados.
- **Quem age:** ninguém direto (tela é dashboard de navegação).

**Dados carregados:**
- `GET /api/cycles?status=processing,review` — competências em aberto
- `GET /api/cycles/current` — competência corrente (mês vigente)
- `GET /api/sync-jobs?status=running,failed&limit=5` — jobs em execução ou falhados nas últimas 24h
- `GET /api/notifications?unread=true&limit=3` — notificações não lidas

**Layout:**
Grid de cards:

1. **Card grande (topo, full width):** "Competência corrente — abril/2026"
   - Status com `<StatusBadge>`
   - `<CycleProgressBar>` com percentual
   - Stats: `X/Y UCs sincronizadas · Z cálculos aprovados · R$ W.WWW,WW total faturado`
   - Botão `[Abrir competência]` → Tela 14

2. **Card "Pendências" (meia largura):**
   - Lista dos 5 itens que mais precisam de atenção:
     - "3 cálculos aguardando revisão" → link pra Tela 14 com filtro
     - "1 UC com sync falhado" → link pra Tela 16 com filtro
     - "2 clientes sem contrato vigente" → link pra Tela 3 com filtro
   - Empty: "Nada pendente 🎉"

3. **Card "Atividade recente" (meia largura):**
   - Últimas 8 entradas de `billing.audit_log` formatadas em linha do tempo.
   - "Paula aprovou cálculo v2 de UC 007098175908 — há 12 min"
   - "João gerou PDF do cliente Condomínio X — há 1 h"
   - Link "Ver tudo" → `/notifications`.

4. **Card "Ações rápidas" (full width, 4 botões grandes):**
   - [+ Novo cliente] → Tela 4
   - [+ Abrir competência] → Tela 13
   - [Sincronizar tudo] → confirma, dispara bulk sync, volta toast
   - [Buscar UC/Cliente] → foca search global

**Estados:** os 4 padrão, aplicados por card (um card com erro não trava os outros).

**Integrações:**
- Cards são queries paralelas (React Query).
- Se tem ciclo em processamento, card 1 assina SSE `/api/events/cycles/{id}` pra atualizar progresso em tempo real.

**Componentes:** Card, `<CycleProgressBar>`, `<StatusBadge>`, `<RelativeTime>`.

---

### Tela 3 — Lista de clientes

- **Rota:** `/customers`
- **Objetivo:** ver todos os clientes, buscar, filtrar, criar novo.
- **Quem vê:** todos.
- **Quem age:** admin/operator criam e editam. Reviewer só lê.

**Dados carregados:**
- `GET /api/customers?q=&status=&tipo_cliente=&cursor=&limit=50`

Resposta:
```json
{
  "items": [{
    "id": "uuid",
    "nome_razao_social": "Cond. Absolut Ville",
    "tipo_pessoa": "PJ",
    "tipo_cliente": "condominio",
    "cpf_cnpj": "12.345.678/0001-99",
    "status": "active",
    "consumer_units_count": 3,
    "cidade_uf": "Salvador/BA",
    "created_at": "2025-10-01T00:00:00Z"
  }],
  "next_cursor": "...",
  "has_more": true
}
```

**Layout:**
- Header: título "Clientes" + botão primário `[+ Novo cliente]` (Tela 4).
- Barra de filtros: busca (nome ou CPF/CNPJ), select status (ativo/inativo/prospecto/arquivado), select tipo (residencial/condomínio/empresa/imobiliária/outro).
- DataTable com colunas:
  - Nome/razão social (bold, clicável → Tela 5)
  - CPF/CNPJ
  - Tipo cliente (badge colorido)
  - Status (`<StatusBadge>`)
  - Qtd UCs (número com ícone)
  - Cidade/UF
  - Cadastrado em (`<DateShort>`)
  - Ações (menu ⋮: Editar, Arquivar, Ver UCs)
- Footer: paginação (cursor-based), total de resultados.
- Botão secundário `[Importar CSV]` desabilitado com badge "Em breve" — conforme o `prompt-front.md`.

**Estados:**
- Empty: `<EmptyState>` "Nenhum cliente ainda" + CTA `[+ Novo cliente]`.
- Filtros sem resultado: "Nenhum cliente corresponde aos filtros. [Limpar filtros]".

**Ações:**
- Click em linha → Tela 5.
- `[+ Novo cliente]` → Tela 4.
- Menu ⋮ → Editar (Tela 6), Arquivar (confirm destrutivo), Ver UCs (Tela 5 aba UCs).

**Componentes:** `<DataTable>`, `<StatusBadge>`, `<EmptyState>`, `<DateShort>`.

---

### Tela 4 — Novo cliente

- **Rota:** `/customers/new`
- **Objetivo:** cadastrar um cliente novo com endereço e opcionalmente já vincular UCs existentes.
- **Quem vê:** admin/operator.
- **Quem age:** admin/operator.

**Dados carregados:**
- Nenhum no mount. Se `?from-uc=<ucCode>` (querystring), pré-carrega a UC pra vincular.

**Layout — form em 4 seções, com "stepper" visual no topo:**

**Seção 1 — Dados principais:**
- Tipo pessoa (radio: PF/PJ)
- Nome / Razão social (required)
- Nome fantasia (só se PJ)
- CPF ou CNPJ (máscara dinâmica baseada em tipo pessoa, valida dígito verificador)
- Tipo cliente (select: residencial/condomínio/empresa/imobiliária/outro)

**Seção 2 — Contato:**
- E-mail
- Telefone (com máscara)

**Seção 3 — Endereço (`<AddressForm>`):**
- CEP (busca ViaCEP ao sair do campo)
- Logradouro, número, complemento, bairro, cidade, UF

**Seção 4 — UCs (opcional):**
- Mini-lista de UCs cadastradas sem vínculo + checkbox pra vincular.
- Ou botão `[+ Adicionar UC manualmente]` que abre form inline com `uc_code, distribuidora, apelido`.
- Ou "Pular e vincular depois".

**Bottom bar fixo:** `[Cancelar]` (volta pra Tela 3) e `[Salvar cliente]`.

**Estados:**
- Loading submit: botão com spinner, form desabilitado.
- Validação: Zod + React Hook Form, erros por campo, highlight vermelho, mensagem abaixo.
- Success: redirect pra Tela 5 do novo cliente + toast "Cliente criado".
- Error (duplicate CPF/CNPJ): toast vermelho + destaque no campo.

**Integrações:**
- `POST /api/customers` com payload completo (cria customer + address + vínculos de UC numa transação no Postgres).
- `GET /api/utils/cep/{cep}` — proxy do BFF pra ViaCEP (evita CORS).

**Componentes:** Form, `<AddressForm>`, Stepper, Input com máscara, Button.

---

### Tela 5 — Detalhe do cliente

- **Rota:** `/customers/[id]`
- **Objetivo:** ver tudo sobre um cliente e agir.
- **Quem vê:** todos.
- **Quem age:** admin/operator.

**Dados carregados:**
- `GET /api/customers/{id}` — dados do cliente
- `GET /api/customers/{id}/consumer-units` — UCs vinculadas
- `GET /api/customers/{id}/contracts` — todos os contratos (de todas as UCs dele) ordenados por vigência desc
- `GET /api/customers/{id}/recent-calculations?limit=10` — últimos cálculos

**Layout — header + abas:**

**Header:**
- Nome/razão social (grande), CPF/CNPJ (monospace ao lado), badge de status, badge de tipo cliente.
- Botões: `[Editar]` → Tela 6, `[⋮ Mais]` (menu: Arquivar, Duplicar).
- Sub-header: endereço principal, e-mail, telefone.

**Abas:**

**Aba 1 — Visão geral (default):**
- Card "Resumo financeiro" (4 KPIs):
  - Total faturado últimos 12 meses: R$
  - Economia última competência: R$ + %
  - Qtd UCs ativas
  - Qtd faturas pagas x abertas
- Card "Próximos vencimentos" (até 5):
  - UC · Competência · Valor · Vencimento · Status
- Card "Últimas atividades" (timeline com 10 entradas de `audit_log`)

**Aba 2 — UCs (conta no badge da aba):**
- Mini-tabela: código UC, apelido, status local, status sync último, contrato vigente (sim/não), ações (Ver detalhe → Tela 8, Sincronizar agora).
- Botão `[+ Vincular UC existente]` e `[+ Adicionar UC manual]`.

**Aba 3 — Contratos:**
- Timeline vertical: um card por versão de contrato.
- Cada card: período de vigência, desconto %, IP mode + valor, flags, botão `[Ver detalhes]` → Tela 11.
- Contrato vigente tem destaque verde.

**Aba 4 — Faturas:**
- Lista consolidada de todas as faturas de todas as UCs do cliente.
- Colunas: UC, Competência, Valor Coelba, Valor Azi, Status, Vencimento, PDF original, PDF cliente.
- Link por linha → Tela 10.

**Aba 5 — Observações:**
- Campo texto livre `notes` editável (autosave com debounce 1s).

**Estados:** padrão por card/aba.

**Ações:**
- `[Editar]` → Tela 6.
- `[Arquivar]` → `<ConfirmDestructive>` → PATCH `/api/customers/{id}` com `status='archived'`.
- Aba UCs `[Sincronizar agora]` → dispara sync da UC.

**Componentes:** Tabs, Card, Timeline, `<DataTable>`, `<BRL>`, `<StatusBadge>`.

---

### Tela 6 — Editar cliente

- **Rota:** `/customers/[id]/edit`
- **Objetivo:** editar dados cadastrais do cliente.
- **Quem vê/age:** admin/operator.

Igual à Tela 4 (mesmo form, mesmas seções), mas pré-preenchido com os dados atuais. Botão "Salvar alterações" → PATCH `/api/customers/{id}` → volta pra Tela 5.

Campos bloqueados pós-criação: CPF/CNPJ (só admin pode mudar após confirmação extra — integridade referencial com contratos existentes). Se admin quiser mudar, modal explica "Isso não afeta contratos existentes porque snapshots estão congelados".

---

### Tela 7 — Lista de UCs

- **Rota:** `/consumer-units`
- **Objetivo:** ver todas as UCs do sistema (inclusive as ainda não vinculadas a cliente, vindas do discover).
- **Quem vê:** todos.
- **Quem age:** admin/operator.

**Dados carregados:**
- `GET /api/consumer-units?q=&status=&linked=&cursor=&limit=50`

Dual source: combina `core.consumer_unit` (cadastral) com `/v1/consumer-units` do backend-go (da API Neoenergia) pra mostrar tanto UCs vinculadas quanto "soltas" esperando vínculo.

**Layout:**
- Header: título "Unidades consumidoras" + botão `[+ Adicionar UC manual]` + botão `[↻ Descobrir UCs]` (que abre modal pra escolher credencial Neoenergia e dispara `/v1/credentials/{id}/discover`).
- Filtros:
  - Busca (código UC ou apelido)
  - Select: "Todas / Vinculadas / Sem cliente / Ativas / Inativas"
  - Select: distribuidora (só neoenergia_ba no MVP)
- DataTable:
  - Código UC (`<UCCode>`, copiável)
  - Apelido
  - Cliente (link; se não tiver, "— Sem vínculo [Vincular]")
  - Classe
  - Cidade/UF
  - Status local (ativa/inativa)
  - Status último sync (`<StatusBadge>` + `<RelativeTime>`)
  - Contrato vigente (✔/✘)
  - Ações: Ver detalhe → Tela 8, Sincronizar agora, Vincular cliente

**Ações:**
- `[↻ Descobrir UCs]` → modal: escolhe credencial → GET `/v1/credentials/{id}/discover` → mostra UCs retornadas vs já cadastradas → marca quais importar → POST `/api/consumer-units/bulk-import`.
- Vincular cliente (para UCs sem cliente) → modal: busca cliente existente ou `[+ Criar novo]` (abre Tela 4 com `?from-uc=<código>`).

**Componentes:** `<DataTable>`, `<UCCode>`, `<StatusBadge>`, `<RelativeTime>`.

---

*(Continua na Parte 2 deste PRD — Telas 8 a 21, Fluxos, Notificações, fora de escopo)*
### Tela 8 — Detalhe da UC

- **Rota:** `/consumer-units/[uc]`
- **Objetivo:** ver tudo sobre uma unidade consumidora.
- **Quem vê:** todos.
- **Quem age:** admin/operator.

**Dados carregados:**
- `GET /api/consumer-units/{uc}` — dados locais + espelho do backend-go
  Combina: `core.consumer_unit` + `/v1/consumer-units/{uc}` (inclui `latest_invoice`, `latest_sync_run`)
- `GET /api/contracts/consumer-unit/{ucId}/active` — contrato vigente
- `GET /api/consumer-units/{uc}/recent-calculations?limit=12` — últimos 12 meses

**Layout — header + 3 áreas:**

**Header:**
- `<UCCode>` grande
- Nome do cliente + link pra Tela 5
- Status da UC (ligada/desligada — vem do `imovel.situacao` da API Neoenergia) como badge
- Apelido, classe, endereço físico
- Botões:
  - `[Sincronizar agora]` (primário)
  - `[Ver faturas]` → Tela 9
  - `[Contrato]` → Tela 11
  - `[⋮ Mais]` (menu: Editar UC, Desvincular cliente)

**Área 1 — Painel de sincronização:**
- Card mostrando último sync: status + timestamp relativo + botão `[Ver detalhes]` → Tela 17 (sync-run).
- Se último sync falhou: alerta amarelo/vermelho com motivo resumido + link pra evidências.
- Botão `[Sincronizar agora]` grande, com modal de opções:
  - Credencial a usar (select das credenciais vinculadas a este cliente via `credential_link`)
  - Checkbox "Baixar PDF"
  - Checkbox "Extrair texto do PDF"
  - `[Iniciar]` → POST `/api/consumer-units/{uc}/sync` → redirect pra Tela 17 do novo sync-run.

**Área 2 — Contrato vigente (card destacado):**
- Mini-resumo do contrato ativo:
  - Vigência de X até (aberta)
  - Desconto: 15% (0,85)
  - IP: R$ 10,00 fixo (ou X% do subtotal)
  - Flags: "Bandeira com desconto: não" / "Custo disponibilidade: sempre cobrado"
- Botões: `[Ver contrato completo]` → Tela 11, `[Novo contrato]` (abre form) se role permite.
- Se não tem contrato vigente: alerta vermelho "UC sem contrato vigente — cálculos não podem ser feitos. [Criar contrato]".

**Área 3 — Histórico de cálculos (últimos 12 meses):**
- Tabela compacta:
  - Mês/ano
  - Competência status (`<StatusBadge>`)
  - Valor Coelba
  - Valor Azi
  - Economia (R$ e %)
  - Status cálculo (`<StatusBadge>`)
  - Ações: Ver detalhe → Tela 15, PDF cliente (se gerado), PDF original Coelba
- Mini-gráfico (recharts) à direita: barras com valor Azi + linha com consumo kWh, últimos 12 meses.

**Estados:** padrão.

**Integrações críticas:**
- "Sincronizar agora" é a principal ação de valor — confira formato do payload do backend-go.

**Componentes:** `<UCCode>`, `<StatusBadge>`, `<BRL>`, Card, `<DataTable>`, BarChart (recharts).

---

### Tela 9 — Faturas da UC

- **Rota:** `/consumer-units/[uc]/invoices`
- **Objetivo:** ver histórico completo de faturas Coelba para uma UC.
- **Quem vê:** todos.
- **Quem age:** admin/operator (dispara reprocessamento).

**Dados carregados:**
- `GET /api/consumer-units/{uc}/invoices?status=&limit=50&cursor=` → chama `/v1/consumer-units/{uc}/invoices` do Go.

Cada item:
```json
{
  "id": "...",                       // invoice_id no SQLite do Go
  "numero_fatura": "339800707843",
  "mes_referencia": "2026/04",
  "status_fatura": "A Vencer",
  "valor_total": "521.53",
  "codigo_barras": "...",
  "data_emissao": "2026-04-14",
  "data_vencimento": "2026-05-06",
  "completeness_status": "complete",
  "extractor_status": "ok",
  "extractor_confidence": 0.9,
  "created_at": "..."
}
```

**Layout:**
- Header: "Faturas da UC `<UCCode>`" + botão `[Sincronizar agora]` (mesmo modal da Tela 8).
- Filtros: busca por número, select status, range de datas de vencimento.
- DataTable:
  - Competência (mês/ano)
  - Nº fatura (monospace)
  - Emissão (`<DateShort>`)
  - Vencimento (`<DateShort>`, com destaque se atrasou)
  - Valor Coelba (`<BRL>`)
  - Status fatura (badge: A Vencer, Pago, Em atraso, etc — vem da Coelba)
  - Completude (chip: "complete" verde / "partial" amarelo com tooltip listando `completeness_missing` / "failed" vermelho)
  - Extração (ícone com tooltip mostrando `extractor_status` + `<ConfidenceBar>`)
  - Ações (menu):
    - Ver detalhe → Tela 10
    - Baixar PDF original
    - Abrir cálculo (se existe) → Tela 15

**Ações:**
- `[Baixar PDF original]` → GET `/api/invoices/{id}/pdf` → BFF busca `invoice_documents.file_data_base64` do Go, converte pra blob, download.

**Componentes:** `<DataTable>`, `<BRL>`, `<DateShort>`, `<StatusBadge>`, `<ConfidenceBar>`, `<PdfDownloadButton>`.

---

### Tela 10 — Detalhe da fatura Coelba

- **Rota:** `/consumer-units/[uc]/invoices/[invoiceId]`
- **Objetivo:** inspecionar uma fatura original da Coelba em detalhe — dados estruturados, itens, origem de cada campo, PDF. **Não é o cálculo da Azi.**
- **Quem vê:** todos.
- **Quem age:** admin/operator (reprocessar, corrigir manualmente).

**Dados carregados:**
- `GET /api/invoices/{id}` → `GET /v1/invoices/{id}` do Go.

Retorna: `billing_record` (JSON), `document_record` (JSON), `items` (array), PDF disponível separadamente.

**Layout — 2 colunas 60/40:**

**Coluna esquerda (60%):**

**Seção 1 — Cabeçalho da fatura:**
- UC, número, competência, status, valor total Coelba
- Emissão, vencimento, período
- Código de barras (monospace com botão copiar)

**Seção 2 — Itens (DataTable):**
- Descrição · Quantidade · Tarifa · Valor · ICMS
- Badge "ignorado no cálculo" (cinza) para itens com `ignored_in_calc=true` (IRRF, reativo, bandeira verde)
- Cada item mostra source de origem (`<SourceIndicator>` baseado em `billing_record.source_map`)
- Footer da tabela: soma dos itens + sanity check (delta vs `valor_total` da fatura; se diverge > 0.02, warning amarelo)

**Seção 3 — SCEE (se presente):**
- Card destacado: "Sistema de Compensação (SCEE/MMGD)"
- Badge do layout: `mmgd_legado` / `mmgd_transicao` / `scee_moderno`
- Campos: Energia injetada kWh, Excedente kWh, Créditos utilizados, Saldo próximo ciclo
- Texto original do rodapé, expansível.

**Seção 4 — Completude:**
- Status badge + se "partial", lista dos campos faltando (`completeness_missing`)
- Extractor status + confidence bar
- Warnings do extractor (se houver)

**Seção 5 — Histórico de consumo (últimos 12 meses):**
- Gráfico de barras do `billing_record.historico_consumo`.

**Coluna direita (40%):**

**Painel do PDF:**
- Viewer embed do PDF original da Coelba (iframe ou `react-pdf`).
- Botão `[Baixar original]` no topo — muito importante, você pediu explicitamente.
- Botão `[Abrir em nova aba]`.

**Painel "Metadata":**
- Criado em / Atualizado em
- Sync run que trouxe (link → Tela 17)
- Botão `[Reprocessar fatura]` — reroda o extractor. Confirm destrutivo porque invalida cálculos associados.
- Botão `[Ver JSON completo]` → modal com `<JsonViewer>` do `billing_record` inteiro.

**Estados:** padrão. Se PDF não foi baixado ainda, viewer vazio com CTA "Sincronizar com PDF".

**Ações:**
- `[Baixar original]` → idem Tela 9.
- `[Reprocessar fatura]` → POST `/api/invoices/{id}/reprocess` → confirm → dispara `sync_job` tipo `extract_pdf`.
- `[Abrir cálculo]` (se existe) → Tela 15.
- `[Ir para UC]` → Tela 8.

**Componentes:** `<BRL>`, `<DataTable>`, `<SourceIndicator>`, `<ConfidenceBar>`, `<PdfDownloadButton>`, `<JsonViewer>`, BarChart.

---

### Tela 11 — Contrato vigente + histórico

- **Rota:** `/consumer-units/[uc]/contract`
- **Objetivo:** ver, criar e versionar contrato de uma UC.
- **Quem vê:** todos.
- **Quem age:** **só admin** (política: mudança de contrato afeta cobrança do cliente, risco alto).

**Dados carregados:**
- `GET /api/contracts/consumer-unit/{ucId}/active`
- `GET /api/contracts/consumer-unit/{ucId}` — histórico completo (todas as versões)

**Layout:**

**Topo — Contrato vigente (card destacado com borda verde):**
- Vigência de DD/MM/YYYY até (aberta)
- Grid 2x2 com parâmetros:
  - Desconto: **15% (0,85)**
  - IP da usina: **R$ 10,00 fixo**
  - Bandeira com desconto: **Não** (badge info)
  - Custo de disponibilidade sempre cobrado: **Sim** (badge info)
- Notas (se houver)
- Criado por Fulano em data
- Botão `[+ Nova versão]` (só admin) → abre form (bottom sheet ou modal)
- Botão `[Ver JSON]` → `<JsonViewer>`

**Form "Nova versão" (modal):**
- Data de início da nova vigência (date picker, default = dia 1 do próximo mês)
- Desconto % (input com validação 0 < x <= 100)
- IP mode (radio: fixo / percentual)
- IP valor (condicional)
- Flags (switches)
- Notas
- Aviso vermelho: **"Ao criar, o contrato vigente será fechado em <data - 1 dia>. Cálculos já aprovados não são afetados (snapshots estão congelados). Cálculos de competências futuras usarão as novas regras."**
- Botões: [Cancelar] [Criar nova versão]
- POST `/api/contracts` → backend-go faz tudo em transação.

**Meio — Histórico (timeline vertical):**
- Um card por versão (incluindo a ativa, destacada).
- Ordem: mais recente no topo.
- Cada card:
  - Vigência (de / até)
  - Diff resumido em relação à versão anterior (ex: "Desconto mudou de 20% para 15%" — comparação calculada no cliente)
  - Status badge (active / ended)
  - Botões: `[Ver JSON]`, `[Duplicar em nova versão]` (só admin — preenche form com esses valores)

**Baixo — Impacto (info):**
- Card: "Este contrato é usado em X cálculos aprovados. Alterações não afetam cálculos já realizados."
  - X = count de `billing_calculation` onde `contract_snapshot_json.id = contract.id`.
- Importante explicar isso visualmente — é a razão de snapshots existirem.

**Ações:**
- Criar nova versão — fluxo acima.
- Não existe "editar" versão existente — imutável por design.

**Componentes:** Card, Form, `<JsonViewer>`, `<StatusBadge>`, Timeline.

**Regra crítica:** se a UC não tem contrato vigente, a tela mostra vazio com CTA grande "[+ Criar primeiro contrato]".

---

### Tela 12 — Lista de competências

- **Rota:** `/cycles`
- **Objetivo:** ver todas as competências mensais (abertas, fechadas, arquivadas).
- **Quem vê:** todos.
- **Quem age:** admin/operator.

**Dados carregados:**
- `GET /api/cycles?year=&status=&cursor=&limit=50`

**Layout:**
- Header: "Competências" + botão `[+ Abrir nova competência]` → Tela 13.
- Filtros: select ano, select status.
- Cards por mês (grid 3 colunas em desktop):

Cada card:
- Título: **abril/2026** (nome do mês grande, ano)
- Status badge (`open`, `syncing`, `processing`, `review`, `approved`, `closed`)
- `<CycleProgressBar>`
- Stats: `X/Y sincronizadas · Z aprovadas`
- Total faturado `<BRL>` (se já tem cálculos)
- Criado por · Criado em · Fechado em (se fechada)
- Botão `[Abrir]` → Tela 14

**Estados:**
- Empty: "Nenhuma competência aberta ainda" + CTA.

**Componentes:** Card, `<CycleProgressBar>`, `<StatusBadge>`, `<BRL>`, `<DateShort>`.

---

### Tela 13 — Abrir competência

- **Rota:** `/cycles/new`
- **Objetivo:** abrir competência nova e configurar escopo inicial.
- **Quem vê/age:** admin/operator.

**Layout — form simples:**
- Ano (select, default = ano corrente)
- Mês (select, default = mês corrente) — se já existe ciclo pra esse mês/ano, mostra erro: "Competência já existe. [Ver]"
- Checkbox "Incluir todas as UCs ativas com contrato vigente" (default: ✅)
- Se desmarcar: aparece lista de UCs onde escolher quais incluir.
- Preview do lado direito:
  - "X UCs serão incluídas nesta competência"
  - Lista resumida

**Botões:** [Cancelar] [Abrir competência]

**Ações:**
- POST `/api/cycles` → retorna ID → redirect pra Tela 14.
- Após abrir, status inicial = `open` (ainda não sincronizou nada).

**Componentes:** Form, Select, Checkbox.

---

### Tela 14 — Dashboard da competência

- **Rota:** `/cycles/[id]`
- **Objetivo:** a tela principal de trabalho do operador. Aqui ele sincroniza, revisa cálculos, aprova, fecha.
- **Quem vê:** todos.
- **Quem age:** admin/operator (revisa/aprova/fecha). Reviewer pode aprovar cálculo.

**Dados carregados:**
- `GET /api/cycles/{id}` — resumo da competência
- `GET /api/cycles/{id}/rows?q=&status=&needs_review_only=&cursor=&limit=100` — tabela principal
- **SSE** `/api/events/cycles/{id}` — progresso em tempo real

Cada row:
```json
{
  "consumer_unit_id": "uuid",
  "uc_code": "007098175908",
  "customer_name": "Paula Fernandes",
  "customer_id": "...",
  "sync_status": "synced" | "pending" | "syncing" | "failed",
  "sync_run_id": "...",
  "invoice_id": "...",
  "valor_coelba": "521.53",
  "completeness_status": "complete" | "partial",
  "calculation_id": "...",
  "calculation_status": "draft" | "needs_review" | "approved",
  "valor_azi_sem_desconto": "...",
  "valor_azi_com_desconto": "...",
  "economia_rs": "...",
  "economia_pct": "...",
  "pdf_generated": true,
  "needs_review_reasons": ["unclassified_item", "low_confidence"]
}
```

**Layout:**

**Header:**
- Título: "Competência **abril/2026**"
- Status badge grande
- Botão `[⋮ Ações]`: Fechar competência (admin only, condicional), Recalcular tudo, Exportar CSV (disabled MVP), Deletar (só se aberta e sem cálculos).

**Painel de progresso (topo, full width):**
- 4 barras empilhadas lado a lado:
  - UCs sincronizadas: X / Y
  - Cálculos gerados: X / Y
  - Cálculos aprovados: X / Y
  - PDFs gerados: X / Y
- Cada barra clicável → filtra a tabela abaixo pelo status correspondente.

**Barra de ações em massa (flutuante acima da tabela quando há seleção):**
- "[N selecionados] [Sincronizar] [Recalcular] [Aprovar] [Gerar PDF] [Cancelar seleção]"

**Filtros acima da tabela:**
- Busca: nome do cliente ou UC
- Select: status de sync
- Select: status de cálculo
- Toggle: "Apenas os que precisam revisão"

**Tabela principal (DataTable, linha por UC):**
- Checkbox de seleção
- UC (`<UCCode>`)
- Cliente (link → Tela 5)
- Sync status (ícone: ✅ ⏳ ❌ + timestamp relativo)
- Completude da fatura (chip)
- Valor Coelba (`<BRL>`)
- Cálculo status (`<StatusBadge>`) + motivos de revisão como badges pequenos se `needs_review`
- Valor Azi (com desconto, `<BRL>`)
- Economia (R$ e %)
- PDF cliente (ícone ✅ / "—" / botão [Gerar])
- Ações (menu ⋮):
  - Ver cálculo → Tela 15
  - Sincronizar esta UC
  - Recalcular
  - Aprovar (se needs_review/draft)
  - Ver fatura Coelba → Tela 10
  - Baixar PDF original (atalho)

**Bottom — resumo:**
- "Total faturado (com desconto): R$ XX.XXX,XX"
- "Economia total do mês: R$ X.XXX,XX (X%)"

**Estados:**
- Loading: skeleton da tabela + progress bar parada.
- Empty: "Nenhuma UC incluída nesta competência. [Adicionar UCs]"
- SSE reconecta automaticamente em caso de drop (tanstack-query + evento).

**Integrações críticas:**
- SSE tipos de evento:
  - `sync_progress`: `{uc_code, status, run_id}`
  - `calculation_progress`: `{uc_code, status, calculation_id, valor_com_desconto}`
  - `row_status_changed`: `{uc_code, ...row}` — substitui row inteira
  - `cycle_status_changed`: `{status}` — muda badge do header
- Bulk actions:
  - `POST /api/cycles/{id}/bulk` com `{action: 'sync'|'recalculate'|'approve'|'generate_pdf', uc_codes: [...]}`
  - Retorna `jobs_created: N` → toast "N jobs enfileirados. [Ver progresso]" → dropdown se atualiza via SSE.

**Ações individuais:**
- Sincronizar → modal compacto com "credential_id + include_pdf" → dispara.
- Aprovar → se `needs_review`, modal: "Há motivos para revisão: X, Y. Confirma aprovação?" → confirm → POST.
- Fechar competência:
  - Disabled se algum cálculo não está `approved`.
  - `<ConfirmDestructive>`: "Digite 'abril/2026' para confirmar. Ao fechar, cálculos ficam imutáveis e não podem mais receber ajustes."
  - POST `/api/cycles/{id}/close`.

**Componentes:** `<DataTable>`, `<CycleProgressBar>`, `<StatusBadge>`, `<UCCode>`, `<BRL>`, `<ConfirmDestructive>`.

---

### Tela 15 — Detalhe do cálculo

- **Rota:** `/cycles/[id]/rows/[calculationId]`
- **Objetivo:** inspecionar um cálculo em profundidade, aplicar ajustes, recalcular, aprovar, gerar PDF do cliente.
- **Quem vê:** todos.
- **Quem age:** admin/operator (ajuste, recalc, gerar PDF). Reviewer e acima aprovam.

**Dados carregados:**
- `GET /api/calculations/{id}` — cálculo atual (inclui 3 snapshots + `needs_review_reasons`)
- `GET /api/calculations/{id}/versions` — histórico de versões do mesmo invoice_ref
- `GET /api/calculations/{id}/adjustments` — ajustes manuais aplicados
- `GET /api/invoices/{sync_invoice_id}` — fatura Coelba correspondente

**Layout — 3 áreas verticais:**

**Área 1 — Cabeçalho:**
- UC (`<UCCode>`) · Cliente · Competência
- Status (`<StatusBadge>`) + Versão (`v2 de 3`)
- Botões:
  - `[Aprovar]` se draft/needs_review (primário verde)
  - `[Recalcular]` (secundário)
  - `[Ajustar manualmente]` (secundário)
  - `[Gerar PDF cliente]` — disabled se não aprovado
  - `[Baixar PDF Coelba]` (secundário — SEMPRE disponível, você pediu)
  - `[Ver fatura Coelba]` → Tela 10
  - `[⋮ Mais]`: Duplicar como rascunho, Desaprovar (admin)

**Área 2 — Comparativo (card grande):**
- 2 colunas:
  - **SEM DESCONTO** — valor Coelba lado a lado: linha energia + linha bandeira + linha IP → soma.
  - **COM DESCONTO** — linha energia (x 0.85) + linha bandeira + linha IP Coelba + linha IP Usina → soma.
- Footer destacado verde: **"Economia: R$ 83,45 (15%)"**

**Área 3 — Abas:**

**Aba 1 — Detalhamento (default):**
- Tabela linha a linha do `result_snapshot.linhas`:
  - Label · Quantidade (kWh quando aplicável) · Preço unitário · Valor sem desconto · Valor com desconto
- Warnings do motor (se houver) — amarelos.
- Motivos de revisão (`needs_review_reasons`) — vermelhos. Cada motivo explica o que fazer:
  - `unclassified_item` → "Há item na fatura que o classifier não reconheceu. [Ver fatura] [Adicionar classificação]"
  - `low_confidence` → "Extractor teve confidence < 0.7. Verifique o PDF original."
  - `missing_scee` → "SCEE esperado mas não detectado no rodapé. [Inserir manualmente]"
  - `sanity_check_failed` → "Soma dos itens diverge do valor da fatura em R$ X. [Ver delta]"

**Aba 2 — Contrato usado (snapshot):**
- `<JsonViewer>` do `contract_snapshot_json`.
- Texto explicativo: "Estes são os termos do contrato **exatamente como estavam** no momento do cálculo. Mesmo que o contrato mude, este cálculo continuará usando estes valores."
- Link "Ver contrato atual → Tela 11" (se diferente).

**Aba 3 — Inputs (o que entrou no motor):**
- `<JsonViewer>` do `inputs_snapshot_json` (itens classificados + SCEE + consumo mínimo).
- Resumo: X itens passaram, Y ignorados, Z não classificados.

**Aba 4 — Ajustes manuais:**
- Lista cronológica de todos os `manual_adjustment` deste cálculo (e das versões anteriores).
- Cada ajuste:
  - Campo alterado
  - Valor antes → valor depois
  - Motivo
  - Autor + timestamp
- Botão `[+ Aplicar novo ajuste]`:
  - Modal com dropdown de campo (inputs.itens[i].preco_unitario, inputs.ip_coelba, etc)
  - Novo valor
  - Motivo obrigatório (textarea, min 10 chars)
  - Aviso: "Isto cria uma nova versão do cálculo. A versão atual vira 'superseded'."

**Aba 5 — Versões:**
- Tabela de todas as versions do mesmo `utility_invoice_ref_id`:
  - v1 (superseded), v2 (superseded), v3 (approved)...
  - Cada row clicável → carrega essa versão nas abas 1-4 (readonly).
  - Diff resumido entre versões.

**Aba 6 — Documentos gerados:**
- Lista de `generated_document`:
  - customer_invoice_pdf v1 (gerado em ...)
  - customer_invoice_pdf v2 (gerado em ... — se reemitido)
- Botão por row: `[Baixar]` `[Preview]`.

**Ações:**
- `[Aprovar]` → se `needs_review`, modal com motivos e "Confirma mesmo assim?". POST `/api/calculations/{id}/approve`. Toast + refresh.
- `[Recalcular]` → POST `/api/calculations/{id}/recalculate`. Mostra loading. Retorna nova version. Substitui URL pra `[newId]`.
- `[Ajustar manualmente]` → fluxo da Aba 4.
- `[Gerar PDF cliente]` → POST `/api/calculations/{id}/generate-pdf`. Enfileira job. Toast com link pra job.
- `[Baixar PDF Coelba]` → GET `/api/invoices/{sync_invoice_id}/pdf`.
- `[Desaprovar]` (admin) → confirm destrutivo → POST `/api/calculations/{id}/unapprove` (endpoint só existe se cycle ainda não está closed).

**Componentes:** Tabs, `<JsonViewer>`, `<BRL>`, `<StatusBadge>`, `<ConfirmDestructive>`, `<PdfDownloadButton>`.

---

*(Continua na Parte 3 — Telas 16-21, Fluxos, Notificações, Fora de Escopo)*
### Tela 16 — Fila de jobs (sync-jobs)

- **Rota:** `/sync-jobs`
- **Objetivo:** transparência operacional — ver todos os jobs enfileirados (sync de UC, extração de PDF, cálculo, geração de PDF do cliente) e seus status.
- **Quem vê:** admin/operator.
- **Quem age:** admin/operator (retry).

**Dados carregados:**
- `GET /api/sync-jobs?type=&status=&since=&cursor=&limit=50`

**Layout:**
- Header: "Fila de jobs" + contador por status (Pending: X, Running: Y, Failed: Z)
- Filtros:
  - Select tipo (`sync_uc`, `extract_pdf`, `calculate`, `generate_pdf`, `recalculate_cycle`)
  - Select status
  - Date range "desde"
- DataTable:
  - ID (8 chars, monospace, tooltip com uuid completo)
  - Tipo (badge)
  - Status (`<StatusBadge>` + ícone animado se running)
  - Payload resumido (ex: "UC 007098175908" extraído do `payload_json`)
  - Retry count / max_retries
  - Criado em
  - Iniciado em · Finalizado em (se aplicável)
  - Ações: Ver detalhe → Tela 17, Retry (se failed)

**Componentes:** `<DataTable>`, `<StatusBadge>`, `<RelativeTime>`.

---

### Tela 17 — Detalhe do job / sync-run

- **Rota:** `/sync-jobs/[id]`
- **Objetivo:** debugar um job específico — ver payload completo, resposta, erros, evidências (screenshots Playwright em caso de sync falho).
- **Quem vê:** admin/operator.
- **Quem age:** admin/operator (retry, reprocess).

**Dados carregados:**
Duas fontes dependendo do tipo:
- Se é sync de UC: `GET /api/sync-runs/{id}` → `/v1/sync-runs/{id}` do Go. Retorna `raw_response`, `uc`, `status`, `error_message`, etc.
- Se é job interno de billing: `GET /api/sync-jobs/{id}` → lê de `billing.sync_job`. Retorna `payload_json`, `status`, `retry_count`, `error_message`.

**Layout — único, com painéis colapsáveis:**

**Painel 1 — Identificação:**
- Tipo + Status grande
- UC (se aplicável, link pra Tela 8)
- Criado em · Iniciado em · Finalizado em · Duração calculada
- Retry N/max

**Painel 2 — Payload (collapsável, aberto por padrão):**
- `<JsonViewer>` do `payload_json`.

**Painel 3 — Resultado / Erro:**
- Se success: `<JsonViewer>` compacto do resumo (para sync_uc: `billing_record`, `document_record`).
- Se failed: mensagem de erro + tipo de erro + stack (se admin).

**Painel 4 — Evidências (se sync_uc falhado):**
- Se o backend-go salvou screenshot/html de erro (vem em `last_error_context` quando Playwright falha):
  - Thumbnail do screenshot (clique → modal full-screen).
  - Link pra HTML salvo.
  - Step name onde falhou ("login", "selecionar_estado", "baixar_fatura", etc).

**Painel 5 — Ações:**
- `[Retry]` → POST `/api/sync-jobs/{id}/retry`. Disabled se job nunca foi failed.
- `[Ver UC]` → Tela 8.
- `[Copiar payload]` → clipboard.
- `[Abrir fatura criada]` → se aplicável, link pra Tela 10.

**Painel 6 — Resposta bruta (collapsável, fechado por padrão):**
- `<JsonViewer>` do `raw_response_json` completo.

**Componentes:** `<JsonViewer>`, `<StatusBadge>`, Collapsible.

---

### Tela 18 — Configurações gerais

- **Rota:** `/settings`
- **Objetivo:** configurações globais do backoffice.
- **Quem vê:** admin.
- **Quem age:** admin.

**Layout — abas:**

**Aba 1 — Empresa:**
- Nome da empresa (Azi Dourado)
- CNPJ
- Logo (upload, usado nos PDFs gerados)
- Endereço
- Dados bancários (PIX, etc — usados no PDF do cliente)
- Salvar → PATCH `/api/settings`.

**Aba 2 — Padrões de faturamento:**
- Desconto padrão sugerido (usado como default ao criar novos contratos)
- IP padrão sugerido
- Template de PDF do cliente (dropdown se houver múltiplos templates no futuro)

**Aba 3 — Notificações:**
- Toggle "Enviar e-mail quando sync falhar"
- Toggle "Alertar quando competência estiver pronta pra revisão"
- E-mails de destino (admin e/ou operator, lista editável).

**Componentes:** Tabs, Form, Toggle, Upload.

---

### Tela 19 — Usuários do backoffice

- **Rota:** `/settings/users`
- **Objetivo:** gerenciar os usuários internos (admin, operator, reviewer).
- **Quem vê:** admin.
- **Quem age:** admin.

**Dados carregados:**
- `GET /api/users?active=&role=`

**Layout:**
- Header: "Usuários" + botão `[+ Novo usuário]`
- DataTable:
  - Nome
  - E-mail
  - Role (badge)
  - Ativo (toggle)
  - Último login (`<RelativeTime>`)
  - Criado em
  - Ações: Editar (modal), Resetar senha (modal com nova senha gerada), Desativar

**Modal "Novo usuário":**
- Nome, e-mail, role (select), senha inicial (auto-gerada com botão "regenerar" + "copiar").
- Ao salvar: POST `/api/users`. Toast mostra senha gerada uma única vez + botão copiar.

**Modal "Editar":**
- Mesmos campos exceto senha (que tem botão "Resetar senha" separado).

**Regras:**
- Não pode deletar o próprio usuário.
- Sempre precisa existir pelo menos 1 admin ativo.

**Componentes:** `<DataTable>`, Form, Modal.

---

### Tela 20 — Credenciais de integração

- **Rota:** `/settings/integrations`
- **Objetivo:** gerenciar credenciais Neoenergia que o backend-go usa pra sincronizar.
- **Quem vê:** admin.
- **Quem age:** admin.

**Dados carregados:**
- `GET /api/credentials` — lista de `core.credential_link` + espelho (status) do backend-go

**Layout:**
- Header: "Credenciais Neoenergia" + botão `[+ Nova credencial]`.
- Lista de cards (uma credencial por card):
  - Label (ex: "neo-paula")
  - Documento mascarado (ex: `***.***.789-01`)
  - UF · Tipo acesso
  - Cliente vinculado (link pra Tela 5)
  - Status da última sessão (`<StatusBadge>` + relative time)
  - Botões:
    - `[Testar sessão]` → POST `/api/credentials/{id}/session` no Go → toast success/failure.
    - `[Descobrir UCs]` → GET `/v1/credentials/{id}/discover` → modal com UCs encontradas.
    - `[Desativar]` (destrutivo).

**Modal "Nova credencial":**
- Vincular a qual cliente (autocomplete)
- Label
- CPF/CNPJ
- Senha (input password, com toggle "mostrar")
- UF (default BA)
- Tipo acesso (normal / imobiliaria)
- Aviso: "A senha é criptografada no backend. Nunca é exibida de volta."
- POST `/api/credentials` (BFF chama `POST /v1/credentials` do Go e cria o `credential_link` no Postgres).

**Componentes:** Card, Form, Modal, `<StatusBadge>`.

---

### Tela 21 — Inbox de notificações

- **Rota:** `/notifications`
- **Objetivo:** ver histórico completo de eventos do sistema (o que o sininho na topbar mostra só as últimas 10).
- **Quem vê:** todos.
- **Quem age:** marcar como lida / arquivar.

**Dados carregados:**
- `GET /api/notifications?unread=&type=&cursor=&limit=50`

Fonte: `billing.audit_log` + tabela nova `core.notification` (criada nesta tela — ver §8).

**Layout:**
- Header: "Notificações" + filtros (todas / não lidas / só meus / por tipo)
- Lista vertical de notificações. Cada item:
  - Ícone do tipo (sync/calc/approval/error)
  - Título (ex: "Sync da UC 007098175908 falhou")
  - Subtítulo (descrição, ator, tempo)
  - Ações inline: [Ver detalhe] [Marcar como lida] [Arquivar]
  - Classe visual: bold se não lida, cinza se lida, riscada se arquivada.
- Paginação infinita (scroll).

**Ações:**
- Click em notificação → navega pro recurso (sync-run, calculation, cycle, etc).
- `[Marcar todas como lidas]` no topo.

**Componentes:** Lista vertical, Filter pills, Ícones.

---

## 7. Fluxos end-to-end

Cada fluxo mostra em que ordem o usuário usa as telas, pra garantir que o agente de front não se perde entre elas.

### Fluxo A — Onboarding de cliente novo até primeiro faturamento

1. **Tela 1** — Login.
2. **Tela 3** → `[+ Novo cliente]` → **Tela 4**: preenche dados, salva.
3. Redirect pra **Tela 5** (detalhe do cliente).
4. Aba "UCs" → `[+ Adicionar UC]` ou `[+ Vincular UC existente]`.
5. Se usar Descobrir UCs: vai pra **Tela 20** → cria credencial → volta pra **Tela 7** → `[↻ Descobrir]`.
6. UC vinculada → vai pra **Tela 8** (detalhe da UC).
7. **Tela 11** → `[+ Criar primeiro contrato]`.
8. Voltar pra **Tela 8** → `[Sincronizar agora]`.
9. Quando sync terminar (notificação), abrir **Tela 14** da competência corrente.
10. Cálculo automático aparece como `draft` → revisar → aprovar → gerar PDF.

### Fluxo B — Fechamento mensal (competência)

1. **Tela 2** (dashboard) → card de pendências mostra "competência X aberta".
2. **Tela 12** → seleciona competência → **Tela 14**.
3. `[Sincronizar tudo]` (bulk) → SSE mostra progresso em tempo real.
4. À medida que syncs terminam, cálculos geram automaticamente (trigger do backend, não do front).
5. Filtra "Apenas os que precisam revisão" → revisa um a um via **Tela 15**.
6. Em **Tela 15**, cada cálculo: verifica, aplica ajuste se necessário, aprova, gera PDF.
7. Volta pra **Tela 14**. Quando tudo aprovado: `[Fechar competência]` → `<ConfirmDestructive>`.
8. Competência fechada — status muda, cálculos viram imutáveis.

### Fluxo C — Fatura com ajuste manual

1. **Tela 14** → linha com `needs_review`.
2. **Tela 15** → aba "Detalhamento": vê o motivo (ex: "item não classificado").
3. Aba "Inputs" → ve que tem item faltando.
4. Aba "Ajustes manuais" → `[+ Aplicar ajuste]` → explica campo, valor novo, motivo.
5. Sistema cria v2. v1 vira superseded.
6. Aba "Detalhamento" da v2 — se ok, aprova.

### Fluxo D — Sync falhou, investigar e retry

1. **Tela 2** → card "Pendências" mostra "1 sync falhou".
2. Ou **Tela 21** (notificações).
3. Click na notificação → **Tela 17** (detalhe do job).
4. Vê screenshot do erro, HTML, step name → identifica que senha da Coelba mudou.
5. Vai pra **Tela 20** → atualiza credencial.
6. Volta pra **Tela 17** → `[Retry]`.

### Fluxo E — Emergência: cliente reclama do boleto

1. Busca global (topbar) pelo nome do cliente ou UC.
2. **Tela 5** (cliente) → aba "Faturas".
3. Click na fatura em questão → **Tela 10** (detalhe da fatura Coelba).
4. Verifica se o PDF original tem o valor correto — botão `[Baixar original]` abre na hora.
5. Se PDF ok, volta pra **Tela 15** do cálculo dessa fatura.
6. Aba "Detalhamento" — confere linha a linha.
7. Se errou, admin faz ajuste manual e recalcula → gera novo PDF → envia pro cliente.

---

## 8. Sistema de notificações

Três camadas de feedback ao usuário. Importante projetar as três juntas pra não poluir a UI.

### 8.1 Camada 1 — Toasts (feedback imediato)

Biblioteca: `sonner` (já integrada com shadcn).

**Quando usar:**
- Ação síncrona terminou (salvou cliente, aprovou cálculo)
- Erro de validação
- Ação disparou job assíncrono: toast com spinner persistente + link "Ver progresso"

**Quando NÃO usar:**
- Não usar pra eventos passivos (algo que aconteceu em outro tab/worker)
- Não usar pra informação persistente (vai pra notificação, não toast)

**Exemplos:**
```
✓ Cliente criado
✓ Cálculo aprovado
⚠ Há 2 itens não classificados — revise antes de aprovar
✗ Sync falhou: credencial expirada. [Ver detalhes]
⟳ 10 jobs enfileirados. [Acompanhar]
```

### 8.2 Camada 2 — Banner de ciclo ativo (topbar)

Quando há competência em estado `syncing` ou `processing`, aparece uma **barra fina abaixo da topbar** em todas as telas:

```
⏳ Competência abril/2026 processando — 8 de 30 UCs sincronizadas · [Ver]
```

Assina SSE `/api/events/cycles/{currentId}`. Atualiza em tempo real. Clicar → Tela 14.

Some quando status muda pra `review`, `approved` ou `closed`.

### 8.3 Camada 3 — Inbox de notificações (sino na topbar + Tela 21)

**Nova tabela `core.notification`** (adicionar à migration):

```sql
CREATE TABLE core.notification (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID REFERENCES core.app_user(id), -- NULL = para todos
    type         TEXT NOT NULL,   -- 'sync_failed', 'calc_needs_review', 'cycle_closed', etc
    title        TEXT NOT NULL,
    description  TEXT,
    entity_type  TEXT,
    entity_id    UUID,
    link         TEXT,            -- rota interna pra navegar
    severity     TEXT NOT NULL DEFAULT 'info' CHECK (severity IN ('info','success','warning','error')),
    read_at      TIMESTAMPTZ,
    archived_at  TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notification_user_unread
    ON core.notification(user_id, created_at DESC)
    WHERE read_at IS NULL AND archived_at IS NULL;
```

**Eventos que geram notificação:**

| Evento | Severity | Destinatário | Título exemplo |
|---|---|---|---|
| Sync UC falhado | error | todos | "Sync da UC 007098175908 falhou" |
| Cálculo `needs_review` | warning | operator | "Cálculo de Paula precisa de revisão" |
| Cálculo aprovado | info | admin | "João aprovou cálculo da UC X" |
| Competência aberta | info | todos | "Competência abril/2026 aberta" |
| Competência fechada | success | todos | "Competência março/2026 fechada" |
| PDF cliente gerado | info | operator | "PDF cliente pronto: Paula / abril/2026" |
| Credencial expirada | error | admin | "Credencial neo-paula expirou — renove" |
| Contrato novo criado | info | admin | "Novo contrato vigente em UC Y" |

**Origem dos eventos:**
- Backend-go escreve em `billing.audit_log`.
- Um pequeno worker no Next (ou via trigger no Postgres) transforma subset do audit_log em `core.notification` pra cada usuário relevante.
- Alternativamente: o backend-go já escreve direto em `core.notification` via segunda conexão Postgres.

**Sino na topbar (`<NotificationsBell>`):**
- Badge numérico com count de não lidas (query frequente, cache 30s + invalidação via Query on event).
- Click → Popover com últimas 10:
  - Ícone de severity
  - Título em bold
  - Descrição em cinza
  - Relative time
  - Click → navega pro `link` + marca como lida (PATCH `/api/notifications/{id}` com `read_at=NOW()`).
- Link "[Ver todas]" → Tela 21.

**Push em tempo real:**
- Mesma SSE do ciclo (`/api/events/notifications`) enviando eventos de notificação nova.
- Quando nova notificação chega, badge incrementa + toast discreto se severity >= warning.

---

## 9. O que NÃO está no MVP

Lista explícita pro agente de front não inventar trabalho. Cada item vem com justificativa.

| Feature | Por que ficou de fora |
|---|---|
| Portal do cliente final | Escopo diferente, segurança diferente, UX diferente. Vira PRD 2 depois. |
| Mobile native / PWA | Backoffice é uso desktop. Responsivo básico já atende tablet eventual. |
| Gráficos analíticos avançados (dashboards históricos com filtros cruzados) | Priorizar confiabilidade do faturamento antes de BI. Recharts simples em Tela 8 já cobre 80% do valor. |
| Exportação CSV/XLSX | Nice-to-have. Fica como "Em breve" nos lugares óbvios (Tela 3, Tela 14). |
| Importação CSV de clientes | Já marcado "Em breve" no `prompt-front.md`. |
| SSO (Google, Microsoft) | Over-engineering pra 2-5 usuários. Credentials simples bastam. |
| 2FA | Pode entrar no PRD 2. Operador pode resistir no MVP. |
| Trilha completa de auditoria com diff visual | `audit_log` grava tudo, mas o front só mostra timeline simples. Diff visual avançado no PRD 2. |
| Multi-tenant (múltiplas empresas) | Mono-tenant pro Azi Dourado. `core.customer` já é suficiente; não precisa de `tenant_id` por enquanto. |
| Internacionalização (i18n) | Só pt-BR. Azi Dourado opera no Brasil. |
| Dark mode | Nice-to-have. shadcn já vem com, só não vamos investir em custom tokens agora. |
| Editor WYSIWYG de template de PDF | Template único definido no backend, editável só por dev no MVP. |
| Envio automático de e-mail/WhatsApp ao cliente final com o boleto | PRD 2 (portal do cliente + canais). |
| Integração com banco/ERP (registrar boleto, conciliação) | PRD 2. Pediu pra um módulo inteiro dedicado. |
| Relatório fiscal / contábil | Fora de escopo. |
| Histórico de versões de contratos com assinatura eletrônica | Fora de escopo. |

---

## 10. Checklist de implementação (ordem sugerida)

Pro agente de front atacar em ordem que maximiza valor e minimiza refactor:

**Sprint 1 — Fundação (sem qual nada funciona):**
- [ ] Setup rotas, layout, sidebar, topbar (base visual já existe)
- [ ] Tela 1 (Login) + `lib/auth.ts` + middleware
- [ ] Tela 2 (Dashboard, versão vazia com cards placeholder)
- [ ] `backend-go-client.ts` com auth via X-API-Key
- [ ] Drizzle setup + schema `core.*`
- [ ] Componentes: `<DataTable>`, `<StatusBadge>`, `<BRL>`, `<DateShort>`, `<UCCode>`, `<EmptyState>`, `<ConfirmDestructive>`, `<PdfDownloadButton>`, `<JsonViewer>`

**Sprint 2 — Cadastro:**
- [ ] Telas 3, 4, 5, 6 (clientes)
- [ ] Tela 7 (lista UCs)
- [ ] Tela 20 (credenciais) + Discover UCs
- [ ] Tela 19 (usuários, mínimo)

**Sprint 3 — Integração:**
- [ ] Tela 8 (detalhe UC com sync)
- [ ] Tela 9 (lista faturas)
- [ ] Tela 10 (detalhe fatura Coelba) — PDF viewer
- [ ] Tela 16 (fila jobs)
- [ ] Tela 17 (detalhe sync-run)

**Sprint 4 — Contratos e cálculo:**
- [ ] Tela 11 (contrato)
- [ ] Tela 12 (lista ciclos)
- [ ] Tela 13 (abrir ciclo)

**Sprint 5 — Cerne do valor (faturamento):**
- [ ] Tela 14 (dashboard do ciclo, com SSE)
- [ ] Tela 15 (detalhe do cálculo, aba a aba)
- [ ] Bulk actions em Tela 14

**Sprint 6 — Polir e notificar:**
- [ ] Sistema de notificações (tabela + sino + push SSE)
- [ ] Tela 21 (inbox)
- [ ] Tela 2 completa com dados reais
- [ ] Tela 18 (settings)
- [ ] Refactor de estados de erro/empty em tudo
- [ ] QA completo dos 5 fluxos

---

## 11. Notas finais para o agente de front

- **Onboarding do agente:** antes de codificar, leia (a) este PRD inteiro, (b) o `BILLING_INTEGRATION.md` do PR 1 do backend, (c) `prompt-front.md` do repo.
- **Regra de ouro sobre mocks:** se não tem endpoint pronto no backend ainda, criar mock no BFF com shape correto — nunca mocar no client. Quando endpoint real ficar pronto, troca só o BFF.
- **Regra de ouro sobre formato:** nunca mostre datas/valores brutos ao usuário. Sempre via os helpers `<BRL>`, `<DateShort>`, etc.
- **Regra de ouro sobre ações:** antes de implementar `[botão que faz coisa destrutiva]`, decide se precisa de `<ConfirmDestructive>`. Quando em dúvida, use.
- **Regra de ouro sobre estados:** toda tela tem 4 estados. Desenhe os 4 antes de considerar a tela pronta.
- **Sobre cópia (textos):** todos em pt-BR. Tom direto, sem marketing-speak. "Cliente criado" é melhor que "Operação realizada com sucesso". "Valor Azi" é melhor que "Valor Calculado pelo Sistema".

Se algo neste PRD conflita com `prompt-front.md` ou com o `BILLING_INTEGRATION.md`, **este PRD ganha** (é mais recente e mais detalhado). Se algo ficou ambíguo, pergunte ao Gustavo antes de implementar.

---

**Fim do PRD.**
