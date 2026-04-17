# Clientes + UCs + Faturas — Tasks

**Design**: `.specs/features/clientes-ucs-faturas/design.md`
**Status**: In Progress

## Progress
- [x] T01 — deps + vitest (commits: 594630f, 727a2f0)
- [x] T02 — .env.local.example (commit: 49065b4)
- [x] T03 — prisma schema (commit: f12bf04)
- [x] T04 — prisma.config.ts + lib/db.ts (commit: 5b3de7d) ⚠️ migration pendente: requer PostgreSQL rodando
- [x] T05 — go-client.ts + 4 testes (commit: 914ed0a)
- [x] T06 — types/clientes.ts + 5 testes (commit: 17a203e)
- [x] T07–T16 — BFF routes (Phase 2+3)
- [x] T17–T23 — UI components (Phase 4)
- [x] T24–T31 — Pages (Phase 5)
- [x] T32 — seed + docs (Phase 6)

## Desvios registrados
- Prisma 7.7 não aceita `url` em schema.prisma → criado `prisma.config.ts` com `defineConfig`
- `client_id` opcional em `CreateCredentialSchema` (vem do URL path, não do body)
- Zod v4: `.partial()` não funciona em schemas com `.superRefine()` → padrão de base schema

---

## Execution Plan

### Phase 1: Foundation (Sequential — must complete before any other phase)

```
T01 → T02 → T03 → T04 → T05 → T06
```

### Phase 2+3+4: Parallel Execution (after Phase 1 complete)

```
Phase 1 done
    │
    ├── Phase 2: BFF Local CRUD ──────────────────┐
    │   T07 [P] T08 [P] T09 [P] T10 [P] T11 [P]  │
    │                                              │
    ├── Phase 3: BFF Integration ─────────────────┤──► Phase 5
    │   T12 [P] T13 [P] T14 [P] T15 [P] T16 [P]  │
    │                                              │
    └── Phase 4: UI Components ───────────────────┘
        T17 [P] T18 [P] T19 [P] T20 [P]
        T21 [P] T22 [P] T23 [P]
```

### Phase 5: Pages + Nav (Sequential, after Phase 2+3+4)

```
T24 → T25 → T26 → T27 → T28 → T29 → T30 → T31
```

### Phase 6: Documentation

```
T32 (after T31)
```

---

## Task Breakdown

### T01: Project setup — install dependencies + Vitest

**What**: Install Prisma, React Hook Form, Zod, TanStack React Table, React Query, Vitest; add pnpm scripts
**Where**: `package.json`, `vitest.config.ts`
**Depends on**: None
**Reuses**: Nothing — greenfield additions
**Requirement**: All (foundation)

**Commands**:
```bash
pnpm add prisma @prisma/client
pnpm add react-hook-form @hookform/resolvers zod
pnpm add @tanstack/react-table @tanstack/react-query
pnpm add -D vitest @vitest/ui
```

Add to `package.json` scripts:
```json
"test": "vitest run",
"test:watch": "vitest",
"db:migrate": "prisma migrate dev",
"db:seed": "tsx prisma/seed.ts",
"db:studio": "prisma studio"
```

Create `vitest.config.ts`:
```ts
import { defineConfig } from 'vitest/config'
export default defineConfig({
  test: { environment: 'node' }
})
```

**Done when**:
- [ ] `pnpm typecheck` passes
- [ ] `pnpm test` runs (0 tests, no error)
- [ ] All packages appear in `node_modules`

**Tests**: none
**Gate**: typecheck

---

### T02: Environment template + .env.local

**What**: Create `.env.local.example` with required vars and document them
**Where**: `web/.env.local.example`
**Depends on**: T01
**Reuses**: Nothing

**Content**:
```bash
# Local PostgreSQL
DATABASE_URL="postgresql://user:password@localhost:5432/fatura_dev"

# Backend Go integration (server-only — NEVER use NEXT_PUBLIC_)
BACKEND_GO_URL="https://api5g.numbro.app"
```

**Done when**:
- [ ] `.env.local.example` created
- [ ] `.env.local.example` is NOT in `.gitignore` (it's a template)
- [ ] `.env.local` (actual secrets) IS in `.gitignore` — verify

**Tests**: none
**Gate**: none

---

### T03: Prisma schema — all 5 models

**What**: Create `prisma/schema.prisma` with Client, ClientAddress, ConsumerUnit, CommercialData, IntegrationCredential models and all enums
**Where**: `prisma/schema.prisma`
**Depends on**: T02
**Reuses**: Models defined in `design.md` (Data Models section)
**Requirement**: CLNT-01, CLNT-02, UC-01, CRED-01

**Schema must include**:
- `generator client` + `datasource db` pointing to `DATABASE_URL`
- All 5 models with exact fields from design.md
- `@unique` on `cpf_cnpj`, `uc_code`, `go_credential_id`
- Cascade deletes on all child records (Address, UC, Commercial, Credential → Client)
- Enums: `TipoPessoa`, `ClientStatus`, `TipoCliente`

**Done when**:
- [ ] `pnpm prisma validate` passes with no errors
- [ ] `pnpm prisma format` produces clean output

**Tests**: none
**Gate**: `pnpm prisma validate`

---

### T04: Initial migration + `lib/db.ts`

**What**: Run initial migration to create tables; create Prisma Client singleton
**Where**: `prisma/migrations/`, `lib/db.ts`
**Depends on**: T03
**Reuses**: Pattern: Next.js Prisma singleton to avoid connection exhaustion in dev

**`lib/db.ts`**:
```ts
import { PrismaClient } from '@prisma/client'

const globalForPrisma = globalThis as unknown as { prisma: PrismaClient }

export const db = globalForPrisma.prisma || new PrismaClient()

if (process.env.NODE_ENV !== 'production') globalForPrisma.prisma = db
```

**Commands**:
```bash
pnpm prisma migrate dev --name init
```

**Done when**:
- [ ] Migration file created in `prisma/migrations/`
- [ ] `lib/db.ts` exports `db` (PrismaClient instance)
- [ ] `pnpm typecheck` passes (Prisma types generated)
- [ ] `db.client.findMany()` compiles without error

**Tests**: none
**Gate**: typecheck

---

### T05: `lib/go-client.ts` — typed fetch wrapper + unit tests

**What**: Server-only utility for calling Backend Go with timeout, typed errors, and JSON parsing
**Where**: `lib/go-client.ts`, `lib/go-client.test.ts`
**Depends on**: T01
**Reuses**: Native `fetch` API
**Requirement**: SYNC-01, CRED-01, INV-01, SYNC-02

**Interface**:
```ts
export class GoApiError extends Error {
  constructor(
    public status: number,
    public message: string,
    public path: string
  ) { super(message) }
}

export async function goFetch<T>(
  path: string,
  options?: RequestInit & { timeoutMs?: number }
): Promise<T>
```

**Behavior**:
- Reads `BACKEND_GO_URL` from `process.env` (throws if missing)
- Default timeout: 30s; sync endpoints should pass `timeoutMs: 60000`
- Non-2xx response → throws `GoApiError` with status + parsed `error` field from JSON body
- Adds `Content-Type: application/json` header on POST/PATCH
- Never logs request body (security — may contain credentials)

**Tests** (`lib/go-client.test.ts`):
1. Successful GET returns parsed JSON
2. Non-2xx response throws `GoApiError` with correct status
3. Network timeout throws `GoApiError` with status 504
4. Missing `BACKEND_GO_URL` throws configuration error

**Done when**:
- [ ] `pnpm test` passes with 4 tests
- [ ] `pnpm typecheck` passes

**Tests**: unit
**Gate**: quick (`pnpm test`)

---

### T06: `types/clientes.ts` — Zod schemas + TypeScript types + unit tests

**What**: All shared Zod schemas and inferred TypeScript types used by BFF and client components
**Where**: `types/clientes.ts`, `types/clientes.test.ts`
**Depends on**: T01
**Reuses**: Zod
**Requirement**: CLNT-02, UC-01, CRED-01

**Schemas to define**:

*Local CRUD (BFF input validation):*
- `CreateClientSchema` — all Client fields, CPF/CNPJ validated by `tipo_pessoa`
- `UpdateClientSchema` — partial of CreateClientSchema
- `CreateUcSchema` — uc_code, distribuidora, apelido, classe_consumo, endereco_unidade, cidade, uf, credential_id (optional)
- `UpdateUcSchema` — partial minus uc_code (immutable)
- `CreateCredentialSchema` — label, documento, senha (min 6 chars), uf, tipo_acesso

*Pagination:*
- `ClientListQuerySchema` — page, pageSize, search, status, tipo_cliente, archived, orderBy, order

*Go API response shapes (TypeScript interfaces only — no Zod needed for responses):*
- `GoConsumerUnit`, `GoInvoice`, `GoBillingRecord`, `GoDocumentRecord`, `GoSyncRun`, `GoCredential`

**Tests** (`types/clientes.test.ts`):
1. `CreateClientSchema` validates valid PF with CPF-length documento
2. `CreateClientSchema` rejects PF with 14-char documento (CNPJ)
3. `CreateClientSchema` validates valid PJ with CNPJ-length documento
4. `CreateCredentialSchema` rejects when `senha` is empty string
5. `ClientListQuerySchema` applies default values correctly

**Done when**:
- [ ] `pnpm test` passes with 5+ tests
- [ ] `pnpm typecheck` passes

**Tests**: unit
**Gate**: quick (`pnpm test`)

---

### T07: BFF — `GET /api/clients` + `POST /api/clients` [P]

**What**: Route handler for listing clients (paginated, filtered) and creating a new client with address
**Where**: `app/api/clients/route.ts`
**Depends on**: T04, T06
**Reuses**: `db` from `lib/db.ts`, `CreateClientSchema` from `types/clientes.ts`
**Requirement**: CLNT-01, CLNT-02

**GET /api/clients**:
- Parse query with `ClientListQuerySchema` (Zod safeParse)
- Prisma query: include `_count: { ucs: true }`, `address: true`
- Filter: `archived_at` null unless `archived=true`
- Return: `{ data, total, page, pageSize }`

**POST /api/clients**:
- Validate body with `CreateClientSchema`
- Prisma create with optional nested `address` upsert
- Return 201 with created client
- On unique constraint violation (cpf_cnpj) → 409

**Done when**:
- [ ] GET `/api/clients` returns `{ data: [], total: 0, ... }` on empty DB
- [ ] POST creates client and returns 201
- [ ] POST with duplicate cpf_cnpj returns 409
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T08: BFF — `GET /api/clients/[id]` + `PATCH /api/clients/[id]` [P]

**What**: Route handler for fetching a single client and updating it
**Where**: `app/api/clients/[id]/route.ts`
**Depends on**: T04, T06
**Reuses**: `db`, `UpdateClientSchema`
**Requirement**: CLNT-02, CLNT-03

**GET /api/clients/[id]**:
- Include: `address`, `ucs`, `commercial_data`, `credentials`
- Not found → 404

**PATCH /api/clients/[id]**:
- Validate body with `UpdateClientSchema`
- Upsert address if address fields present
- Upsert commercial_data if commercial fields present
- Not found → 404

**Done when**:
- [ ] GET seed client by ID returns full object with relations
- [ ] PATCH updates `nome_razao` correctly
- [ ] PATCH to non-existent ID returns 404
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T09: BFF — `POST /api/clients/[id]/archive` [P]

**What**: Route handler for archiving a client (sets archived_at + status inativo)
**Where**: `app/api/clients/[id]/archive/route.ts`
**Depends on**: T04
**Reuses**: `db`
**Requirement**: CLNT-04

**POST /api/clients/[id]/archive**:
- Set `archived_at = new Date()` and `status = 'inativo'`
- Idempotent (archive already archived → 200)
- Not found → 404
- Return updated client

**Done when**:
- [ ] Archive sets `archived_at` to non-null
- [ ] GET list after archive excludes client by default
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T10: BFF — `GET /api/clients/[id]/ucs` + `POST /api/clients/[id]/ucs` [P]

**What**: Route handler for listing UCs of a client and creating a new UC
**Where**: `app/api/clients/[id]/ucs/route.ts`
**Depends on**: T04, T06
**Reuses**: `db`, `CreateUcSchema`
**Requirement**: UC-01

**GET**: Return all UCs for client, include credential (id, label, documento_masked)

**POST**:
- Validate with `CreateUcSchema`
- On `uc_code` duplicate → 409
- Return 201

**Done when**:
- [ ] POST creates UC linked to client
- [ ] Duplicate uc_code → 409
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T11: BFF — `PATCH /api/ucs/[id]` [P]

**What**: Route handler for updating a UC's local fields
**Where**: `app/api/ucs/[id]/route.ts`
**Depends on**: T04, T06
**Reuses**: `db`, `UpdateUcSchema`
**Requirement**: UC-02

**PATCH**:
- Validate with `UpdateUcSchema` (excludes uc_code)
- Not found → 404
- Return updated UC

**Done when**:
- [ ] PATCH updates `apelido` correctly
- [ ] PATCH body containing `uc_code` → 422 (field not in UpdateUcSchema)
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T12: BFF — `POST /api/integration/credentials` + `POST /api/integration/credentials/[id]/session` [P]

**What**: BFF proxy for creating Go credentials + establishing session
**Where**: `app/api/integration/credentials/route.ts`, `app/api/integration/credentials/[id]/session/route.ts`
**Depends on**: T04, T05, T06
**Reuses**: `goFetch`, `db`, `CreateCredentialSchema`
**Requirement**: CRED-01

**POST /api/integration/credentials**:
1. Validate body with `CreateCredentialSchema`
2. Call `goFetch<GoCredential>('POST', '/v1/credentials', { label, documento, senha, uf, tipo_acesso })`
3. On Go success: persist to local DB `IntegrationCredential` with `go_credential_id`, `documento_masked` (from Go response), `label`, `uf`, `tipo_acesso`, `client_id` (from query param or body)
4. Return local record (NEVER return `go_credential_id` to client — return local `id`)
5. On Go error: return Go's error status/message, do NOT persist locally

**POST /api/integration/credentials/[id]/session**:
- Lookup local credential by `id`, get `go_credential_id`
- Call `goFetch<GoSession>('POST', `/v1/credentials/${go_credential_id}/session`)`
- Return session info (no local persistence needed)

**Security**: senha never logged, never returned, never saved

**Done when**:
- [ ] `pnpm typecheck` passes
- [ ] POST body without `senha` → 422
- [ ] Go error → local DB NOT updated (verify via Prisma)

**Tests**: none (manual — requires Go backend)
**Gate**: typecheck

---

### T13: BFF — `POST /api/integration/ucs/[uc]/sync` [P]

**What**: BFF proxy for triggering UC sync on Go backend
**Where**: `app/api/integration/ucs/[uc]/sync/route.ts`
**Depends on**: T05, T06
**Reuses**: `goFetch`
**Requirement**: SYNC-01

**POST**:
- Body: `{ credential_id: string }` (local credential id)
- Lookup local credential → get `go_credential_id`
- Call `goFetch('POST', `/v1/consumer-units/${uc}/sync`, { credential_id: go_credential_id, include_pdf: true, include_extraction: true }, { timeoutMs: 60000 })`
- Return `{ sync_run_id, invoice_id, status }` from persistence block
- 504 on timeout, 503 on Go unavailable

**Done when**:
- [ ] Route compiles without TypeScript errors
- [ ] Unknown local credential_id → 404 before calling Go
- [ ] `pnpm typecheck` passes

**Tests**: none (manual — requires Go backend)
**Gate**: typecheck

---

### T14: BFF — `GET /api/integration/ucs` + `GET /api/integration/ucs/[uc]` [P]

**What**: BFF proxies for listing all UCs and getting a single UC from Go backend
**Where**: `app/api/integration/ucs/route.ts`, `app/api/integration/ucs/[uc]/route.ts`
**Depends on**: T05
**Reuses**: `goFetch`
**Requirement**: UC-01, SYNC-01

**GET /api/integration/ucs**: Forward query params `limit`, `status`
**GET /api/integration/ucs/[uc]**: Forward `uc` path param

**Done when**:
- [ ] Both routes compile without errors
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T15: BFF — `GET /api/integration/ucs/[uc]/invoices` + `GET /api/integration/invoices/[id]` [P]

**What**: BFF proxies for listing invoices and getting invoice detail from Go backend
**Where**: `app/api/integration/ucs/[uc]/invoices/route.ts`, `app/api/integration/invoices/[id]/route.ts`
**Depends on**: T05
**Reuses**: `goFetch`
**Requirement**: INV-01, INV-02

**GET .../invoices**: Forward `limit`, `status` query params
**GET /api/integration/invoices/[id]**: Return full invoice with billing_record + document_record + items

**Done when**:
- [ ] Both routes compile without errors
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T16: BFF — `GET /api/integration/sync-runs/[id]` [P]

**What**: BFF proxy for fetching sync audit trail from Go backend
**Where**: `app/api/integration/sync-runs/[id]/route.ts`
**Depends on**: T05
**Reuses**: `goFetch`
**Requirement**: SYNC-02

**GET**: Return full sync run including status, error_message, raw_response

**Done when**:
- [ ] Route compiles without errors
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T17: UI — `status-badge.tsx` + `completude-badge.tsx` [P]

**What**: Two small display components: status badge for client/UC/sync states; completude badge for invoice extraction quality
**Where**: `components/clientes/status-badge.tsx`, `components/clientes/completude-badge.tsx`
**Depends on**: T06
**Reuses**: `Badge` from `components/ui/badge.tsx`, `cn()` from `lib/utils.ts`
**Requirement**: CLNT-01, INV-01, INV-02

**`StatusBadge`**:
```tsx
type StatusBadgeProps = {
  status: 'ativo' | 'inativo' | 'prospecto' | 'succeeded' | 'partial' | 'failed' | 'running'
}
```
Color map: ativo/succeeded → green; inativo/failed → red; prospecto/partial → yellow; running → blue

**`CompletudeBadge`**:
```tsx
type CompletudeBadgeProps = { status: 'complete' | 'partial' | 'failed' }
```
Same color convention + icon.

**Done when**:
- [ ] Both components render without errors
- [ ] All status variants have distinct colors
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T18: UI — `client-table.tsx` (TanStack Table) [P]

**What**: Client list table with search input, status/tipo filter dropdowns, pagination, sortable columns, row click navigation
**Where**: `components/clientes/client-table.tsx`
**Depends on**: T06, T17
**Reuses**: TanStack React Table, `Table`, `Badge`, `Button`, `Skeleton`, `Input`, `Select`, `Pagination`, `DropdownMenu` from shadcn/ui
**Requirement**: CLNT-01

**Columns**: nome_razao (sortable), cpf_cnpj, tipo_cliente, status (StatusBadge), qtd UCs (`_count.ucs`), cidade/UF, created_at (formatted pt-BR), ações (Visualizar)

**Features**:
- `useSearchParams` + `useRouter` for URL-synced filters
- Debounced search (300ms) via `setTimeout`
- Skeleton: 5 rows of `Skeleton` components during `isLoading`
- Empty state: `Empty` component (from `components/ui/empty.tsx`) with CTA
- "Importar CSV" button: disabled with `Badge` "Em breve" beside it

**Done when**:
- [ ] Table renders with seed data (requires running app)
- [ ] Search filters rows client-side (or updates URL params)
- [ ] Skeleton visible during loading state
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T19: UI — `client-form.tsx` (React Hook Form + Zod) [P]

**What**: Unified create/edit form for client with 5 sections; validates with Zod; works for both create and edit
**Where**: `components/clientes/client-form.tsx`
**Depends on**: T06
**Reuses**: React Hook Form, `CreateClientSchema` / `UpdateClientSchema`, shadcn/ui `Input`, `Select`, `Textarea`, `Button`, `Separator`, `Card`
**Requirement**: CLNT-02

**Sections (use `<section>` + `<Separator>` between them)**:
1. **Dados Principais**: tipo_pessoa (Radio: PF/PJ), nome_razao, nome_fantasia, cpf_cnpj, tipo_cliente, status
2. **Contato**: email, telefone
3. **Endereço**: cep, logradouro, numero, complemento, bairro, cidade, uf (Select com estados BR)
4. **Comercial**: tipo_contrato, data_inicio, data_fim, status_contrato, observacoes_comerciais
5. **Observações**: textarea livre

**Props**:
```tsx
type ClientFormProps = {
  defaultValues?: Partial<CreateClientInput>
  onSubmit: (data: CreateClientInput) => Promise<void>
  isLoading?: boolean
  mode: 'create' | 'edit'
}
```

**Done when**:
- [ ] Form renders all 5 sections
- [ ] Inline field errors appear on submit with invalid data
- [ ] `tipo_pessoa` PF/PJ switches CPF/CNPJ mask hint
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T20: UI — `uc-form.tsx` + `credential-form.tsx` [P]

**What**: Form for creating/editing a ConsumerUnit; form for creating an IntegrationCredential
**Where**: `components/clientes/uc-form.tsx`, `components/clientes/credential-form.tsx`
**Depends on**: T06
**Reuses**: React Hook Form, Zod schemas, shadcn/ui components
**Requirement**: UC-01, UC-02, CRED-01

**`UcForm`**:
- Fields: uc_code (disabled on edit), distribuidora, apelido, classe_consumo, endereco_unidade, cidade, uf, ativa (Switch)
- `credential_id` → `<Select>` populated by `credentials: IntegrationCredential[]` prop
- "Adicionar credencial" link beside select (opens `CredentialForm` in Sheet)

**`CredentialForm`**:
- Fields: label, documento, senha (type=password, never stored), uf, tipo_acesso
- Warning text: "A senha é enviada diretamente para a concessionária e nunca é armazenada"
- Submit → POST /api/integration/credentials

**Done when**:
- [ ] UcForm renders with credential select populated
- [ ] CredentialForm has password field with warning
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T21: UI — `uc-list.tsx` + `uc-sync-button.tsx` + `use-sync-polling.ts` [P]

**What**: UC list component, sync trigger button with state machine, polling hook for sync run status
**Where**: `components/clientes/uc-list.tsx`, `components/clientes/uc-sync-button.tsx`, `hooks/use-sync-polling.ts`
**Depends on**: T06, T17
**Reuses**: `Spinner`, `Button`, `Badge`, `StatusBadge`; React Query for polling
**Requirement**: SYNC-01, UC-01

**`useSyncPolling(syncRunId: string | null)`**:
- Uses `useQuery` with `refetchInterval: 2000` while `status === 'running'`
- Returns `{ data: GoSyncRun | undefined, isPolling: boolean }`

**`UcSyncButton`**:
```
State: idle → syncing → polling(syncRunId) → success | error
```
- `idle`: "Sincronizar" button enabled
- `syncing`: POST in flight → spinner + "Sincronizando..."
- `polling`: `useSyncPolling` active → spinner + "Verificando..."
- `success`: toast "Sincronização concluída", call `onSuccess()`
- `error`: toast with `error_message`, return to idle

**`UcList`**:
- Card per UC showing: código, apelido, distribuidora, status ativa, `StatusBadge` for last sync
- Actions: "Sincronizar", "Ver faturas" button

**Done when**:
- [ ] UcSyncButton cycles through states correctly
- [ ] Polling stops when status leaves 'running'
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T22: UI — `invoice-table.tsx` + `invoice-detail.tsx` [P]

**What**: Invoice list table and invoice detail display components
**Where**: `components/clientes/invoice-table.tsx`, `components/clientes/invoice-detail.tsx`
**Depends on**: T06, T17
**Reuses**: TanStack Table, `CompletudeBadge`, shadcn/ui `Table`, `Badge`, `Collapsible`
**Requirement**: INV-01, INV-02

**`InvoiceTable`**:
- Columns: numero_fatura, mes_referencia, valor_total, data_vencimento, status, completeness (CompletudeBadge), updated_at
- Filter: status dropdown, date range (optional)
- Row click → navigate to detail

**`InvoiceDetail`**:
- Section "Fatura" → `billing_record` fields in a description list
- Section "Documento Extraído" → `document_record` fields with source_map badges (API/PDF/OCR)
- Section "Itens" → table of invoice items
- `CompletudeBadge` in header

**Done when**:
- [ ] Both components compile and render mock data
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T23: UI — `sync-run-detail.tsx` [P]

**What**: Sync audit detail component showing status, timing, error, and collapsible raw_response
**Where**: `components/clientes/sync-run-detail.tsx`
**Depends on**: T06, T17
**Reuses**: `Collapsible`, `StatusBadge`, `Button`, `Alert`
**Requirement**: SYNC-02

**Layout**:
- Header: status (StatusBadge) + created_at formatted
- Error section (shown only when `error_message` present): Alert destructive variant
- Raw response: `<Collapsible>` closed by default, shows JSON.stringify(raw_response, null, 2) in `<pre>`
- CTA "Reprocessar" (only when status === 'failed') → calls sync endpoint again

**Done when**:
- [ ] Component renders with failed sync run mock data
- [ ] Raw response collapsible works
- [ ] "Reprocessar" only visible when failed
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T24: `app/(admin)/layout.tsx` — route group layout

**What**: Layout file for (admin) route group that wraps all client/UC/invoice pages in `SidebarLayout`
**Where**: `app/(admin)/layout.tsx`
**Depends on**: T17 (Phase 4 should be at least started)
**Reuses**: `SidebarLayout` from `components/sidebar-layout.tsx`
**Requirement**: All telas

```tsx
import { SidebarLayout } from '@/components/sidebar-layout'
export default function AdminLayout({ children }: { children: React.ReactNode }) {
  return <SidebarLayout>{children}</SidebarLayout>
}
```

**Done when**:
- [ ] Navigating to `/clientes` renders within the sidebar layout
- [ ] No duplicate sidebars (root `app/page.tsx` unaffected)
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T25: Update `app-sidebar.tsx` — add Clientes nav group

**What**: Add "Clientes" navigation item to existing sidebar with subitems
**Where**: `components/app-sidebar.tsx`
**Depends on**: T24
**Reuses**: Existing nav group pattern in `app-sidebar.tsx`
**Requirement**: CLNT-01

Add nav group after main items:
```
Clientes
  ├── Todos os clientes  → /clientes
  └── Novo cliente       → /clientes/novo
```

Use `UserGroupIcon` (or similar) from HugeIcons.

**Done when**:
- [ ] "Clientes" appears in sidebar
- [ ] Link to `/clientes` works
- [ ] Existing nav items unchanged
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T26: Telas 1 + 2 — `/clientes` list + `/clientes/novo` + `/clientes/[id]/editar`

**What**: Three pages: client list, create client, edit client
**Where**: `app/(admin)/clientes/page.tsx`, `app/(admin)/clientes/novo/page.tsx`, `app/(admin)/clientes/[id]/editar/page.tsx`
**Depends on**: T07, T08, T18, T19, T24, T25
**Reuses**: `ClientTable`, `ClientForm`, `SidebarLayout` (via (admin) layout)
**Requirement**: CLNT-01, CLNT-02

**`/clientes/page.tsx`** (Tela 1):
- Server Component fetches initial data from `/api/clients`
- Passes data to `<ClientTable>` (client component)
- React Query `QueryClientProvider` wraps for client state
- Breadcrumb: Clientes

**`/clientes/novo/page.tsx`** (Tela 2 create):
- Renders `<ClientForm mode="create">`
- On submit → POST `/api/clients` → redirect to `/clientes/:id`
- Breadcrumb: Clientes > Novo cliente

**`/clientes/[id]/editar/page.tsx`** (Tela 2 edit):
- Server Component fetches existing client data
- Renders `<ClientForm mode="edit" defaultValues={client}>`
- On submit → PATCH `/api/clients/:id` → redirect to `/clientes/:id`
- Breadcrumb: Clientes > [nome] > Editar

**Done when**:
- [ ] List page shows seed data in table
- [ ] Filters update URL and table
- [ ] Create form creates client and redirects
- [ ] Edit form loads existing data and saves
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T27: Tela 3 — `/clientes/[id]` — client detail with tabs

**What**: Client detail page with tabbed layout: Dados, Endereço, UCs, Comercial, Integração
**Where**: `app/(admin)/clientes/[id]/page.tsx`
**Depends on**: T08, T10, T17, T19, T20, T24
**Reuses**: `Tabs`, `TabsContent`, `TabsList`, `TabsTrigger`, `ClientForm` (read-only sections), `UcList`
**Requirement**: CLNT-03

**Header**: nome, `StatusBadge`, tipo_cliente, cpf_cnpj, actions row (Editar, Arquivar)

**Tabs**:
- **Dados** → display-only fields from client
- **Endereço** → address fields display
- **UCs** → `<UcList>` component with add UC button
- **Comercial** → commercial_data display
- **Integração** → list of `IntegrationCredential` (masked doc, label, Go ID hidden), add credential button

**Arquivar action**: `AlertDialog` confirmation before POST /api/clients/:id/archive

**Done when**:
- [ ] All 5 tabs render without errors
- [ ] Archive dialog appears and works
- [ ] Add UC opens `UcForm` in Sheet
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T28: Tela 4 — `/clientes/[id]/ucs` — UC panel

**What**: Full UC panel page for a specific client showing all UCs with sync status
**Where**: `app/(admin)/clientes/[id]/ucs/page.tsx`
**Depends on**: T10, T14, T21, T24
**Reuses**: `UcList`, `UcSyncButton`, React Query for live sync status
**Requirement**: UC-01, SYNC-01

- Fetch local UCs from GET /api/clients/:id/ucs
- For each UC with `uc_code`, fetch Go status from GET /api/integration/ucs/:uc (client-side, React Query)
- Merge: local UC data + Go status (latest_invoice, latest_sync_run)
- Render `<UcList>` with merged data
- "Sincronizar agora" per UC triggers `<UcSyncButton>`

**Done when**:
- [ ] Page lists UCs with local + Go data merged
- [ ] Sync button functional (requires Go backend)
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

### T29: Tela 5 — `/clientes/[id]/ucs/[ucId]/faturas` — invoices list

**What**: Invoice list page for a specific UC, data from Go backend
**Where**: `app/(admin)/clientes/[id]/ucs/[ucId]/faturas/page.tsx`
**Depends on**: T15, T22, T24
**Reuses**: `InvoiceTable`, React Query
**Requirement**: INV-01

- `ucId` is the local UC id; resolve to `uc_code` via GET /api/clients/:id/ucs
- Fetch invoices from GET /api/integration/ucs/:uc/invoices
- Render `<InvoiceTable>`
- Empty state with CTA "Sincronizar UC" (link to UC panel)
- Breadcrumb: Clientes > [nome] > UCs > [uc_code] > Faturas

**Done when**:
- [ ] Page renders (requires Go backend for data)
- [ ] Empty state shows CTA when no invoices
- [ ] `pnpm typecheck` passes

**Tests**: none (manual — requires Go backend)
**Gate**: typecheck

---

### T30: Tela 6 — `/clientes/[id]/ucs/[ucId]/faturas/[faturaId]` — invoice detail

**What**: Invoice detail page showing billing_record, document_record, items and completeness
**Where**: `app/(admin)/clientes/[id]/ucs/[ucId]/faturas/[faturaId]/page.tsx`
**Depends on**: T15, T22, T24
**Reuses**: `InvoiceDetail`
**Requirement**: INV-02

- `faturaId` is the Go invoice UUID
- Fetch from GET /api/integration/invoices/:faturaId
- Render `<InvoiceDetail>`
- Breadcrumb: ... > Faturas > [numero_fatura]

**Done when**:
- [ ] Page compiles without errors
- [ ] Renders with Go backend data
- [ ] `pnpm typecheck` passes

**Tests**: none (manual — requires Go backend)
**Gate**: typecheck

---

### T31: Tela 7 — `/clientes/[id]/ucs/[ucId]/sync/[syncId]` — sync audit

**What**: Sync run audit page with collapsible raw_response and reprocess CTA
**Where**: `app/(admin)/clientes/[id]/ucs/[ucId]/sync/[syncId]/page.tsx`
**Depends on**: T16, T23, T24
**Reuses**: `SyncRunDetail`
**Requirement**: SYNC-02

- Fetch from GET /api/integration/sync-runs/:syncId
- Render `<SyncRunDetail>`
- "Reprocessar" CTA triggers POST /api/integration/ucs/:uc/sync with same credential

**Done when**:
- [ ] Page compiles without errors
- [ ] Renders with Go backend data
- [ ] Reprocessar button calls sync endpoint
- [ ] `pnpm typecheck` passes

**Tests**: none (manual — requires Go backend)
**Gate**: typecheck

---

### T32: Seed + README update + prisma/seed.ts

**What**: Development seed with realistic data; README section for this module
**Where**: `prisma/seed.ts`, `README.md`
**Depends on**: T04
**Reuses**: `db`, Prisma models
**Requirement**: Dev environment

**Seed content**:
- 30 clients (mix of PF/PJ, all status variants, all tipo_cliente variants)
- Each client with address
- 3 clients with commercial data
- 15 UCs spread across 10 clients (multiple UCs per client for some)
- 5 IntegrationCredential records (masked docs, fake go_credential_id UUIDs)
- 10 UCs linked to credentials

**README section** ("Módulo Clientes - Desenvolvimento"):
- Prerequisites: PostgreSQL running, Go backend for Telas 5/6/7
- Setup: `pnpm db:migrate && pnpm db:seed`
- Start: `DATABASE_URL=... BACKEND_GO_URL=https://api5g.numbro.app pnpm dev`

**Done when**:
- [ ] `pnpm db:seed` completes without errors
- [ ] GET /api/clients returns 30 clients
- [ ] README section added
- [ ] `pnpm typecheck` passes

**Tests**: none (manual)
**Gate**: typecheck

---

## Parallel Execution Map

```
Phase 1 (Sequential):
  T01 → T02 → T03 → T04 ──┬── T05 → T06
                           │
Phase 2+3+4 (all parallel after Phase 1):

  After T04+T06:
    T07 [P] ─────────────────────────────────────────┐
    T08 [P] ─────────────────────────────────────────┤
    T09 [P] ─────────────────────────────────────────┤
    T10 [P] ─────────────────────────────────────────┤
    T11 [P] ─────────────────────────────────────────┤
                                                      │
  After T05+T06:                                      │
    T12 [P] ─────────────────────────────────────────┤
    T13 [P] ─────────────────────────────────────────┤
    T14 [P] ─────────────────────────────────────────┤
    T15 [P] ─────────────────────────────────────────┤
    T16 [P] ─────────────────────────────────────────┤
                                                      │
  After T06:                                          │
    T17 [P] ─────────────────────────────────────────┤
    T18 [P] ─────────────────────────────────────────┤
    T19 [P] ─────────────────────────────────────────┤
    T20 [P] ─────────────────────────────────────────┤
    T21 [P] ─────────────────────────────────────────┤
    T22 [P] ─────────────────────────────────────────┤
    T23 [P] ─────────────────────────────────────────┘
                           │
Phase 5 (sequential after all Phase 2+3+4):
  T24 → T25 → T26 → T27 → T28 → T29 → T30 → T31

Phase 6:
  T31 → T32
```

---

## Task Granularity Check

| Task | Scope | Status |
|------|-------|--------|
| T01: Install deps + Vitest | 1 setup operation | ✅ |
| T02: .env.local.example | 1 file | ✅ |
| T03: schema.prisma | 1 file | ✅ |
| T04: migration + db.ts | 1 migration + 1 utility | ✅ cohesive |
| T05: go-client.ts + tests | 1 utility + tests | ✅ |
| T06: types/clientes.ts + tests | 1 types file + tests | ✅ |
| T07-T11: BFF local routes | 1 route.ts per task | ✅ |
| T12-T16: BFF integration routes | 1-2 related routes per task | ✅ cohesive |
| T17: 2 badge components | 2 tiny display components | ✅ cohesive |
| T18: client-table.tsx | 1 complex component | ✅ |
| T19: client-form.tsx | 1 form component | ✅ |
| T20: uc-form + credential-form | 2 related forms | ✅ cohesive |
| T21: uc-list + sync-button + hook | 1 feature cluster | ✅ cohesive |
| T22: invoice-table + detail | 2 related components | ✅ cohesive |
| T23: sync-run-detail | 1 component | ✅ |
| T24-T25: layout + nav | 2 infrastructure files | ✅ cohesive |
| T26-T31: Pages | 1 page per task (or 3 in T26) | ✅ |
| T32: seed + README | 1 dev setup task | ✅ cohesive |

---

## Diagram-Definition Cross-Check

| Task | Depends On (body) | Diagram Shows | Status |
|------|------------------|---------------|--------|
| T01 | None | Start of Phase 1 | ✅ |
| T02 | T01 | T01 → T02 | ✅ |
| T03 | T02 | T02 → T03 | ✅ |
| T04 | T03 | T03 → T04 | ✅ |
| T05 | T01 | After Phase 1 | ✅ |
| T06 | T01 | After Phase 1 | ✅ |
| T07 [P] | T04, T06 | Phase 2 parallel | ✅ |
| T08 [P] | T04, T06 | Phase 2 parallel | ✅ |
| T09 [P] | T04 | Phase 2 parallel | ✅ |
| T10 [P] | T04, T06 | Phase 2 parallel | ✅ |
| T11 [P] | T04, T06 | Phase 2 parallel | ✅ |
| T12 [P] | T04, T05, T06 | Phase 3 parallel | ✅ |
| T13 [P] | T05, T06 | Phase 3 parallel | ✅ |
| T14 [P] | T05 | Phase 3 parallel | ✅ |
| T15 [P] | T05 | Phase 3 parallel | ✅ |
| T16 [P] | T05 | Phase 3 parallel | ✅ |
| T17 [P] | T06 | Phase 4 parallel | ✅ |
| T18 [P] | T06, T17 | Phase 4 parallel | ✅ |
| T19 [P] | T06 | Phase 4 parallel | ✅ |
| T20 [P] | T06 | Phase 4 parallel | ✅ |
| T21 [P] | T06, T17 | Phase 4 parallel | ✅ |
| T22 [P] | T06, T17 | Phase 4 parallel | ✅ |
| T23 [P] | T06, T17 | Phase 4 parallel | ✅ |
| T24 | Phase 4 done | Phase 5 start | ✅ |
| T25 | T24 | T24 → T25 | ✅ |
| T26 | T07, T08, T18, T19, T24, T25 | Phase 5 seq | ✅ |
| T27 | T08, T10, T17, T19, T20, T24 | Phase 5 seq | ✅ |
| T28 | T10, T14, T21, T24 | Phase 5 seq | ✅ |
| T29 | T15, T22, T24 | Phase 5 seq | ✅ |
| T30 | T15, T22, T24 | Phase 5 seq | ✅ |
| T31 | T16, T23, T24 | Phase 5 seq | ✅ |
| T32 | T04 (+ T31 logically) | Phase 6 | ✅ |

---

## Test Co-location Validation

| Task | Layer Created | Matrix Requires | Task Says | Status |
|------|--------------|-----------------|-----------|--------|
| T05 | `lib/go-client.ts` | unit | unit | ✅ |
| T06 | `types/clientes.ts` | unit | unit | ✅ |
| T07-T16 | Route Handlers | none | none | ✅ |
| T17-T23 | UI Components | none | none | ✅ |
| T24-T31 | Pages | none | none | ✅ |
| T32 | Seed | none | none | ✅ |
