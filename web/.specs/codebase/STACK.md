# Stack

## Runtime
- Next.js 16.1.6 (App Router)
- React 19.2.4
- TypeScript 5.9 (strict mode)
- Node.js (pnpm)

## Styling
- TailwindCSS v4 (PostCSS plugin, OKLCH color space)
- shadcn/ui (style: radix-lyra, 59 components installed)
- tw-animate-css
- HugeIcons (free tier, @hugeicons/react)

## Font
- Oxanium (Google Font, CSS var: --font-sans)

## Utilities
- clsx + tailwind-merge via `cn()` in `@/lib/utils`
- date-fns 4.1.0
- next-themes (dark mode with "D" hotkey)

## ORM / DB
- NOT YET INSTALLED — must add Prisma (decision documented in design.md)
- PostgreSQL target (no Supabase)

## Forms / Validation
- NOT YET INSTALLED — must add React Hook Form + Zod

## Tables
- NOT YET INSTALLED — must add @tanstack/react-table

## Data Fetching (client)
- NOT YET INSTALLED — must add @tanstack/react-query

## Package Manager
- pnpm (pnpm-lock.yaml present)

## Build
- `pnpm dev` (Turbopack)
- `pnpm build` / `pnpm typecheck` / `pnpm lint`
- TypeScript errors currently ignored in build (next.config.mjs)
