# Clientes + UCs + Faturas — Design

**Spec**: `.specs/features/clientes-ucs-faturas/spec.md`
**Status**: Draft

---

## Architecture Overview

```
Browser
  │
  │  fetch /api/*  (HTTPS, internal only)
  ▼
Next.js Route Handlers (BFF)
  ├── /api/clients/*          ──► Prisma ──► PostgreSQL (local)
  ├── /api/ucs/*              ──► Prisma ──► PostgreSQL (local)
  └── /api/integration/*      ──► go-client.ts ──► Backend Go
                                                       │
                                                  PostgreSQL (Go)
                                                  (invoices, sync runs)
```

**Invariant**: O browser nunca faz chamadas para o Backend Go. Todo acesso passa pelo BFF Next.js.

---

## Tech Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| ORM | **Prisma** | Migrações declarativas, geração de tipos, `prisma migrate dev` — superior ao Drizzle para schemas com 5+ tabelas relacionadas. Drizzle seria mais adequado para edge runtime, que não é o caso aqui. |
| Data fetching client | **@tanstack/react-query** | Pares naturalmente com TanStack Table (já obrigatório). Cache, invalidation, polling para sync status — React Query é o padrão. |
| API pattern | Route Handlers (`app/api/`) | Conforme spec: browser → `/api/*` → Go. Server Actions seriam adequados para mutations simples mas não para o BFF proxy pattern. |
| Credential password | **Nunca persiste localmente** | Senha vai ao Go via POST body (HTTPS), nunca trafega em query string, nunca salva no PG local. Localmente: apenas `go_credential_id` + documento mascarado. |
| `credential_id` scope | **Por UC** (não por cliente) | Um cliente pode ter UCs em contas diferentes da concessionária. Vincular por UC é mais granular e correto. |
| Forms | React Hook Form + Zod | Conforme stack obrigatória. Zod schema shared entre client e BFF para double validation. |
| Navigation | Extend `app-sidebar.tsx` | Não recriar sidebar — adicionar item "Clientes" com subitens. |

---

## Folder Structure

```
web/
├── app/
│   ├── (admin)/                        # Route group preserving sidebar layout
│   │   ├── layout.tsx                  # Reuses SidebarLayout
│   │   ├── clientes/
│   │   │   ├── page.tsx                # Tela 1: Lista de clientes
│   │   │   ├── novo/
│   │   │   │   └── page.tsx            # Tela 2: Criar cliente
│   │   │   └── [id]/
│   │   │       ├── page.tsx            # Tela 3: Detalhe do cliente
│   │   │       ├── editar/
│   │   │       │   └── page.tsx        # Tela 2 (edit mode)
│   │   │       └── ucs/
│   │   │           ├── page.tsx        # Tela 4: Painel de UCs
│   │   │           └── [ucId]/
│   │   │               ├── faturas/
│   │   │               │   ├── page.tsx        # Tela 5: Invoices list
│   │   │               │   └── [faturaId]/
│   │   │               │       └── page.tsx    # Tela 6: Invoice detail
│   │   │               └── sync/
│   │   │                   └── [syncId]/
│   │   │                       └── page.tsx    # Tela 7: Sync audit
│   └── api/
│       ├── clients/
│       │   ├── route.ts                # GET list, POST create
│       │   └── [id]/
│       │       ├── route.ts            # GET one, PATCH update
│       │       ├── archive/
│       │       │   └── route.ts        # POST archive
│       │       └── ucs/
│       │           └── route.ts        # GET list, POST create UC
│       ├── ucs/
│       │   └── [id]/
│       │       └── route.ts            # PATCH update UC
│       └── integration/
│           ├── credentials/
│           │   ├── route.ts            # POST → Go POST /v1/credentials
│           │   └── [id]/
│           │       └── session/
│           │           └── route.ts    # POST → Go POST /v1/credentials/{id}/session
│           ├── ucs/
│           │   ├── route.ts            # GET → Go GET /v1/consumer-units
│           │   └── [uc]/
│           │       ├── route.ts        # GET → Go GET /v1/consumer-units/{uc}
│           │       ├── sync/
│           │       │   └── route.ts    # POST → Go POST /v1/consumer-units/{uc}/sync
│           │       └── invoices/
│           │           └── route.ts    # GET → Go GET /v1/consumer-units/{uc}/invoices
│           ├── invoices/
│           │   └── [id]/
│           │       └── route.ts        # GET → Go GET /v1/invoices/{id}
│           └── sync-runs/
│               └── [id]/
│                   └── route.ts        # GET → Go GET /v1/sync-runs/{id}
│
├── components/
│   └── clientes/                       # Feature components
│       ├── client-form.tsx             # Create/edit form (controlled by RHF+Zod)
│       ├── client-table.tsx            # TanStack Table with filters
│       ├── client-detail-tabs.tsx      # Tabs for detalhe page
│       ├── uc-form.tsx                 # UC create/edit form
│       ├── uc-list.tsx                 # UC panel list
│       ├── uc-sync-button.tsx          # Sync action with loading state
│       ├── credential-form.tsx         # Credential create form (password field)
│       ├── invoice-table.tsx           # Invoices TanStack Table
│       ├── invoice-detail.tsx          # Invoice detail sections
│       ├── sync-run-detail.tsx         # Sync audit display
│       ├── status-badge.tsx            # Reusable status badge (ativo/inativo/etc)
│       └── completude-badge.tsx        # complete/partial/failed badge
│
├── lib/
│   ├── db.ts                           # Prisma Client singleton
│   ├── go-client.ts                    # Typed fetch wrapper for Go backend
│   └── utils.ts                        # Existing cn() utility
│
├── types/
│   └── clientes.ts                     # Shared TypeScript types + Zod schemas
│
├── prisma/
│   ├── schema.prisma                   # Data model
│   ├── seed.ts                         # Dev seed (30+ clientes, UCs, invoices mock)
│   └── migrations/                     # Auto-generated by prisma migrate dev
│
└── .env.local.example                  # Template with required env vars
```

---

## Code Reuse Analysis

### Existing Components to Leverage

| Component | Location | How to Use |
|-----------|----------|------------|
| `SidebarLayout` | `components/sidebar-layout.tsx` | Wrap all (admin) pages |
| `AppSidebar` | `components/app-sidebar.tsx` | Add "Clientes" nav group |
| All 59 shadcn/ui components | `components/ui/` | Tabs, Table, Badge, Dialog, Sheet, Form, Input, Select, Button, Skeleton, Toast, etc. |
| `cn()` utility | `lib/utils.ts` | All component styling |
| `useToast` | `hooks/use-toast.ts` | Feedback toasts |
| HugeIcons | `@hugeicons/react` | Icons throughout UI |

### New Dependencies to Install

```bash
pnpm add prisma @prisma/client
pnpm add react-hook-form @hookform/resolvers zod
pnpm add @tanstack/react-table @tanstack/react-query
pnpm add -D prisma
```

---

## Data Models

### Prisma Schema

```prisma
model Client {
  id              String    @id @default(cuid())
  tipo_pessoa     TipoPessoa
  nome_razao      String
  nome_fantasia   String?
  cpf_cnpj        String    @unique
  email           String?
  telefone        String?
  status          ClientStatus @default(prospecto)
  tipo_cliente    TipoCliente
  observacoes     String?
  created_at      DateTime  @default(now())
  updated_at      DateTime  @updatedAt
  archived_at     DateTime?

  address         ClientAddress?
  ucs             ConsumerUnit[]
  commercial_data CommercialData?
  credentials     IntegrationCredential[]
}

enum TipoPessoa    { PF PJ }
enum ClientStatus  { ativo inativo prospecto }
enum TipoCliente   { residencial condominio empresa imobiliaria outro }

model ClientAddress {
  id          String   @id @default(cuid())
  client_id   String   @unique
  client      Client   @relation(fields: [client_id], references: [id], onDelete: Cascade)
  cep         String?
  logradouro  String?
  numero      String?
  complemento String?
  bairro      String?
  cidade      String?
  uf          String?  @db.Char(2)
  created_at  DateTime @default(now())
  updated_at  DateTime @updatedAt
}

model ConsumerUnit {
  id                String   @id @default(cuid())
  client_id         String
  client            Client   @relation(fields: [client_id], references: [id], onDelete: Cascade)
  uc_code           String   @unique
  distribuidora     String?
  apelido           String?
  classe_consumo    String?
  endereco_unidade  String?
  cidade            String?
  uf                String?  @db.Char(2)
  ativa             Boolean  @default(true)
  credential_id     String?
  credential        IntegrationCredential? @relation(fields: [credential_id], references: [id], onDelete: SetNull)
  created_at        DateTime @default(now())
  updated_at        DateTime @updatedAt
}

model CommercialData {
  id                    String   @id @default(cuid())
  client_id             String   @unique
  client                Client   @relation(fields: [client_id], references: [id], onDelete: Cascade)
  tipo_contrato         String?
  data_inicio           DateTime?
  data_fim              DateTime?
  status_contrato       String?
  observacoes_comerciais String?
  created_at            DateTime @default(now())
  updated_at            DateTime @updatedAt
}

model IntegrationCredential {
  id                String   @id @default(cuid())
  client_id         String
  client            Client   @relation(fields: [client_id], references: [id], onDelete: Cascade)
  label             String
  documento_masked  String
  uf                String   @db.Char(2)
  tipo_acesso       String   @default("normal")
  go_credential_id  String   @unique
  created_at        DateTime @default(now())
  updated_at        DateTime @updatedAt

  ucs               ConsumerUnit[]
}
```

### TypeScript Types (types/clientes.ts)

Key Zod schemas shared client↔BFF:
- `CreateClientSchema` — used in form + POST /api/clients
- `UpdateClientSchema` — used in form (edit mode) + PATCH /api/clients/:id
- `CreateUcSchema` — used in UC form + POST /api/clients/:id/ucs
- `CreateCredentialSchema` — used in credential form + POST /api/integration/credentials
- `SyncRequestSchema` — used in sync button + POST /api/integration/ucs/:uc/sync

Go API response types:
- `GoConsumerUnit`, `GoInvoice`, `GoSyncRun`, `GoBillingRecord`, `GoDocumentRecord`

---

## Components

### `client-form.tsx`
- **Purpose**: Formulário de criação e edição de cliente com seções colapsáveis
- **Interfaces**: `ClientFormProps { defaultValues?: Client; onSubmit: (data) => void; isLoading: boolean }`
- **Reuses**: Input, Select, Textarea, Button, Separator from shadcn/ui
- **Sections**: Dados Principais → Contato → Endereço → Comercial → Observações

### `client-table.tsx`
- **Purpose**: TanStack Table com busca, filtros URL-synced, paginação, ordenação
- **Interfaces**: `ClientTableProps { data: Client[]; total: number; isLoading: boolean }`
- **Reuses**: Table, Badge, Button, Skeleton, Input, Select from shadcn/ui
- **Features**: Column visibility toggle, row click → navigate to detail

### `uc-sync-button.tsx`
- **Purpose**: Botão de sync com estados: idle → syncing → success/error
- **Interfaces**: `UcSyncButtonProps { ucCode: string; credentialId: string; onSuccess: () => void }`
- **Reuses**: Button, Spinner, Toast
- **State machine**: `idle | syncing | success | error` (local useState)

### `go-client.ts`
- **Purpose**: Wrapper tipado sobre fetch para chamadas server-to-server ao Go backend
- **Interface**: `goFetch<T>(path, options): Promise<T>` — throws `GoApiError` on non-2xx
- **Timeout**: 30s default, 60s para sync
- **Error types**: `GoApiError { status, message, path }`

---

## Error Handling Strategy

| Scenario | BFF Response | UI Display |
|----------|-------------|------------|
| Go backend down | 503 `{ error: "Serviço indisponível" }` | Toast vermelho + retry |
| Sync timeout (>30s) | 504 `{ error: "Timeout na sincronização" }` | Toast + link para sync-runs |
| UC não encontrada no Go | 404 passado ao cliente | Toast + empty state |
| Validação Zod falha no BFF | 422 `{ errors: ZodError[] }` | Erros inline nos campos |
| CPF/CNPJ duplicado | 409 `{ error: "CPF/CNPJ já cadastrado" }` | Erro no campo cpf_cnpj |
| Prisma unique constraint | 409 mapeado | Mesmo tratamento acima |

---

## Pagination Contract

`GET /api/clients` query params:
```
page         number  default 1
pageSize     number  default 20 (max 100)
search       string  filter nome_razao OR cpf_cnpj (ilike)
status       string  ativo | inativo | prospecto
tipo_cliente string  residencial | condominio | empresa | imobiliaria | outro
archived     boolean default false (true = mostrar arquivados)
orderBy      string  default created_at
order        string  asc | desc
```
Response: `{ data: Client[], total: number, page: number, pageSize: number }`

---

## Sync Polling Strategy

`uc-sync-button.tsx` triggers sync then polls sync run status:
1. POST /api/integration/ucs/:uc/sync → returns `{ sync_run_id }`
2. React Query `refetchInterval: 2000` on GET /api/integration/sync-runs/:id
3. Stop polling when `status !== "running"` (i.e. succeeded | partial | failed)
4. Future: replace interval with server-sent event or WebSocket notification (documented as evolution point)

```tsx
// Sketch inside uc-sync-button.tsx
const { data: syncRun } = useQuery({
  queryKey: ['sync-run', syncRunId],
  queryFn: () => fetchSyncRun(syncRunId),
  enabled: !!syncRunId,
  refetchInterval: (data) => 
    (!data || data.status === 'running') ? 2000 : false,
})
```

---

## Route Group Strategy

- Keep `app/page.tsx` (dashboard) untouched
- New routes live in `app/(admin)/` with its own `layout.tsx` wrapping `SidebarLayout`
- `(admin)` is a route group — URL is `/clientes/...`, not `/admin/clientes/...`
- No breadcrumb conflict: `SidebarLayout` receives breadcrumbs as props

---

## Credential Selection in UC Form

1. Operator creates credential separately (via "Integrações" tab on client detail)
2. When creating/editing a UC, a `<Select>` lists existing `IntegrationCredential[]` for that client
3. UC can be saved without credential — `credential_id` is nullable
4. After linking, sync button becomes enabled

---

## Authentication

Out of scope for this module. No auth checks on `/api/*` or pages. Will be added in a separate auth module.

---

## Development Environment for Go-backed Screens

Telas 5/6/7 (invoices, invoice detail, sync audit) require the Go backend running.
Start with: `BACKEND_GO_URL=https://api5g.numbro.app pnpm dev`
No mock/stub in this phase. Document in README.

---

## Security Notes

1. `BACKEND_GO_URL` e `DATABASE_URL` — nunca `NEXT_PUBLIC_*`
2. Senha de credencial: never logged, never stored, never returned
3. `go-client.ts` roda apenas em Route Handlers (server) — `"use server"` não é necessário pois Route Handlers já são server-side
4. Zod validation no BFF para rejeitar payloads malformados antes de chegar ao Go

---

## Environment Variables

```bash
# .env.local.example
DATABASE_URL="postgresql://user:pass@localhost:5432/fatura_dev"
BACKEND_GO_URL="https://api5g.numbro.app"
```
