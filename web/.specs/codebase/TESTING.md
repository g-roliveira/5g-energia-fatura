# Testing Strategy

## Scope
Unit tests with Vitest for critical utilities only. All other validation via TypeScript typecheck.

## Test Types

| Layer | Test Type | When |
|-------|-----------|------|
| `lib/go-client.ts` | unit | Created/modified |
| `types/clientes.ts` (Zod schemas) | unit | Created/modified |
| Route Handlers | none (manual validation) | — |
| UI Components | none (manual validation) | — |
| Pages | none (manual validation) | — |

## Test Commands

| Gate | Command | When to Run |
|------|---------|-------------|
| quick | `pnpm test` | After creating/modifying go-client or types |
| typecheck | `pnpm typecheck` | After any TypeScript change |
| build | `pnpm build` | Before declaring a phase done |

## Parallelism Assessment

| Test Type | Parallel-Safe | Reason |
|-----------|--------------|--------|
| unit (Vitest) | Yes | No shared state, in-memory only |
| typecheck | Yes | Read-only, no side effects |
| build | No | Shared .next/ output directory |

## Setup Required
```bash
pnpm add -D vitest @vitest/ui
```

Add to package.json scripts:
```json
"test": "vitest run",
"test:watch": "vitest"
```

## Dev Environment for Go-backed Screens (Telas 5/6/7)
Requires Go backend running locally with real data.
No mock/stub in this phase.
Document in README: `BACKEND_GO_URL=https://api5g.numbro.app pnpm dev`
