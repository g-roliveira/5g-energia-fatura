# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
pnpm dev            # starts Docker DB + Next.js (Turbopack)
pnpm typecheck      # tsc --noEmit — run after every change
pnpm test           # vitest run (unit tests only)
pnpm test:watch     # vitest watch mode

pnpm db:migrate     # prisma migrate dev (requires DB running)
pnpm db:seed        # tsx prisma/seed.ts (30 clients, 15 UCs, 5 credentials)
pnpm db:studio      # prisma studio

pnpm lint
pnpm format         # prettier --write
```

`db:migrate`, `db:seed`, `db:studio` use `dotenv -e .env.local` to bridge the `.env.local` file to Prisma CLI. The `.env.local` file must exist (copy from `.env.local.example`).

## Architecture

This is a **dual-backend** SaaS. The browser only ever calls Next.js internal routes (`/api/...`). Those routes act as BFF (Backend for Frontend) and fan out to either the local PostgreSQL or the Go backend.

```
Browser → Next.js App Router
              ├── app/api/clients/**         → Prisma → PostgreSQL (local, port 5433)
              ├── app/api/ucs/**             → Prisma → PostgreSQL
              └── app/api/integration/**    → goFetch → Go backend (BACKEND_GO_URL)
```

**The browser must never call `BACKEND_GO_URL` directly.** It is a server-only env var (no `NEXT_PUBLIC_` prefix).

### Backends

| Backend | Responsibility |
|---------|---------------|
| PostgreSQL (local) | Clients, addresses, consumer units, commercial data, credential metadata |
| Go (`https://api5g.numbro.app`) | Credential encryption, session/token, UC sync ETL, invoices, sync audit |

### Key utility files

- `lib/go-client.ts` — `goFetch<T>(path, options?)` + `GoApiError`. Default timeout 30s; sync calls use 60s. Never call Go from the browser.
- `lib/db.ts` — Prisma singleton (`db`) with `globalThis` cache for dev HMR. Uses `@prisma/adapter-pg` (driver adapter pattern).
- `types/clientes.ts` — Zod schemas (`CreateClientSchema`, `UpdateClientSchema`, `CreateUcSchema`, etc.) and TypeScript interfaces for Go API responses (`GoInvoice`, `GoSyncRun`, etc.).
- `contexts/breadcrumb.tsx` — Auto-breadcrumb system. Exports `useSetBreadcrumbTitle(key, title)` for pages to register dynamic labels.

## Prisma 7.7 — breaking change

Database URL is **not** in `schema.prisma`. It lives in `prisma.config.ts`:

```ts
import { defineConfig, env } from 'prisma/config'
export default defineConfig({ datasource: { url: env('DATABASE_URL') } })
```

## Known workaround: Zod v4 + @hookform/resolvers

`zodResolver(Schema as any)` is required in all three form components. This is a type-level mismatch between `@hookform/resolvers` v5.2.2 (compiled against `@zod/core` with `version.minor = 0`) and Zod v4.3.6 (`minor = 3`). The cast is intentional — do not remove it.

Also: `UpdateClientSchema` cannot use `.partial()` on a schema with `.superRefine()`. The pattern used is an unexported base object schema; `CreateClientSchema` adds `superRefine()`, `UpdateClientSchema` calls `.partial()` on the base.

## Component conventions

**Always use shadcn/ui components — never raw HTML elements.** The project has 59 components in `components/ui/`.

| Instead of | Use |
|-----------|-----|
| `<label>` | `Label` from `@/components/ui/label` |
| `<input type="radio">` | `Controller` + `RadioGroup` + `RadioGroupItem` |
| `<select>` in a form | `NativeSelect` + `NativeSelectOption` (works with `register()`) |
| `<select>` as standalone filter | `Select` + `SelectTrigger` + `SelectContent` + `SelectItem` (Radix, styled popover) |
| `<button>` | `Button` |

Icons: always `HugeiconsIcon` from `@hugeicons/react` with icons from `@hugeicons/core-free-icons`. Do not use lucide-react icons for new UI.

Forms that open from another page or panel must use `Dialog` (modal), not `Sheet` (side panel).

## Breadcrumb system

The `(admin)` route group layout (`app/(admin)/layout.tsx`) provides `BreadcrumbProvider` and renders `BreadcrumbSidebarLayout`. Static segments are auto-labeled (`clientes` → "Clientes", `ucs` → "UCs", etc.).

For dynamic segments (client id, UC id, etc.), each page registers its label:

```tsx
import { useSetBreadcrumbTitle } from '@/contexts/breadcrumb'

const { id } = useParams<{ id: string }>()
const { data: client } = useQuery(...)
useSetBreadcrumbTitle(id, client?.nome_razao) // call unconditionally — hook handles null internally
```

## Route group `(admin)`

All client/UC/invoice pages live under `app/(admin)/`. The layout there provides:
- `QueryClientProvider` (React Query)
- `BreadcrumbProvider` + `BreadcrumbSidebarLayout` (sidebar + automatic breadcrumbs)

Individual admin pages must **not** call `SidebarLayout` directly.

## Environment variables

```bash
DATABASE_URL="postgresql://user:password@localhost:5433/fatura_dev"
BACKEND_GO_URL="https://api5g.numbro.app"   # server-only, never NEXT_PUBLIC_
```

PostgreSQL runs in Docker via `docker-compose.yml` on port **5433** (not the default 5432).

## Testing

Unit tests exist only for `lib/go-client.ts` and `types/clientes.ts`. All other layers (BFF routes, UI components, pages) are verified by `pnpm typecheck` only — no mocks for Go or Prisma.
