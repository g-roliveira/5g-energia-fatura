import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type {
  GoBillingCycle,
  GoCycleListResponse,
  GoCycleRowsResponse,
  GoContract,
  GoBillingCalculation,
  GoBulkActionResponse,
} from '@/types/billing'

// ─── Cycles ───────────────────────────────────────────────────────────────────

export function useCycles(year?: number, status?: string, limit = 50, offset = 0) {
  return useQuery({
    queryKey: ['cycles', { year, status, limit, offset }],
    queryFn: async (): Promise<GoCycleListResponse> => {
      const qs = new URLSearchParams()
      if (year) qs.set('year', String(year))
      if (status) qs.set('status', status)
      qs.set('limit', String(limit))
      qs.set('offset', String(offset))
      const res = await fetch(`/api/billing/cycles?${qs}`)
      if (!res.ok) throw new Error('Erro ao carregar ciclos')
      return res.json()
    },
  })
}

export function useCycle(id: string) {
  return useQuery({
    queryKey: ['cycle', id],
    queryFn: async (): Promise<GoBillingCycle> => {
      const res = await fetch(`/api/billing/cycles/${id}`)
      if (!res.ok) throw new Error('Erro ao carregar ciclo')
      return res.json()
    },
    enabled: !!id,
  })
}

export function useCycleRows(
  id: string,
  opts?: {
    q?: string
    sync_status?: string
    calc_status?: string
    needs_review_only?: boolean
    limit?: number
    offset?: number
  },
) {
  return useQuery({
    queryKey: ['cycle-rows', id, opts],
    queryFn: async (): Promise<GoCycleRowsResponse> => {
      const qs = new URLSearchParams()
      if (opts?.q) qs.set('q', opts.q)
      if (opts?.sync_status) qs.set('sync_status', opts.sync_status)
      if (opts?.calc_status) qs.set('calc_status', opts.calc_status)
      if (opts?.needs_review_only) qs.set('needs_review_only', 'true')
      qs.set('limit', String(opts?.limit ?? 100))
      qs.set('offset', String(opts?.offset ?? 0))
      const res = await fetch(`/api/billing/cycles/${id}/rows?${qs}`)
      if (!res.ok) throw new Error('Erro ao carregar linhas do ciclo')
      return res.json()
    },
    enabled: !!id,
    refetchInterval: (query) => {
      const data = query.state.data as GoCycleRowsResponse | undefined
      const hasPending = data?.items.some(
        (r) => r.sync_status === 'pending' || r.sync_status === 'syncing' || r.calculation_status === 'draft',
      )
      return hasPending ? 2000 : false
    },
  })
}

// ─── Mutations ────────────────────────────────────────────────────────────────

export function useCreateCycle() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (body: { year: number; month: number; include_all_active?: boolean }) => {
      const res = await fetch('/api/billing/cycles', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })
      if (!res.ok) {
        const err = await res.json().catch(() => ({ error: 'Erro ao criar ciclo' }))
        throw new Error(err.error)
      }
      return res.json() as Promise<GoBillingCycle>
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cycles'] })
    },
  })
}

export function useBulkAction() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({
      cycleId,
      body,
    }: {
      cycleId: string
      body: { action: string; uc_codes?: string[]; force_all?: boolean }
    }) => {
      const res = await fetch(`/api/billing/cycles/${cycleId}/bulk`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })
      if (!res.ok) {
        const err = await res.json().catch(() => ({ error: 'Erro na ação em massa' }))
        throw new Error(err.error)
      }
      return res.json() as Promise<GoBulkActionResponse>
    },
    onSuccess: (_, { cycleId }) => {
      queryClient.invalidateQueries({ queryKey: ['cycle-rows', cycleId] })
      queryClient.invalidateQueries({ queryKey: ['cycle', cycleId] })
    },
  })
}

export function useCloseCycle() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (cycleId: string) => {
      const res = await fetch(`/api/billing/cycles/${cycleId}/close`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({}),
      })
      if (!res.ok) {
        const err = await res.json().catch(() => ({ error: 'Erro ao fechar ciclo' }))
        throw new Error(err.error)
      }
      return res.json() as Promise<{ status: string }>
    },
    onSuccess: (_, cycleId) => {
      queryClient.invalidateQueries({ queryKey: ['cycle', cycleId] })
      queryClient.invalidateQueries({ queryKey: ['cycles'] })
    },
  })
}

// ─── Contracts ────────────────────────────────────────────────────────────────

export function useContracts(ucId: string) {
  return useQuery({
    queryKey: ['contracts', ucId],
    queryFn: async (): Promise<{ items: GoContract[]; count: number }> => {
      const res = await fetch(`/api/billing/contracts?uc_id=${encodeURIComponent(ucId)}`)
      if (!res.ok) throw new Error('Erro ao carregar contratos')
      return res.json()
    },
    enabled: !!ucId,
  })
}

export function useContract(id: string) {
  return useQuery({
    queryKey: ['contract', id],
    queryFn: async (): Promise<GoContract> => {
      const res = await fetch(`/api/billing/contracts/${id}`)
      if (!res.ok) throw new Error('Erro ao carregar contrato')
      return res.json()
    },
    enabled: !!id,
  })
}

export function useActiveContract(ucId: string) {
  return useQuery({
    queryKey: ['active-contract', ucId],
    queryFn: async (): Promise<GoContract> => {
      const res = await fetch(`/api/billing/consumer-units/${ucId}/active-contract`)
      if (!res.ok) throw new Error('Erro ao carregar contrato ativo')
      return res.json()
    },
    enabled: !!ucId,
  })
}

export function useCreateContract() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (body: Record<string, unknown>) => {
      const res = await fetch('/api/billing/contracts', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })
      if (!res.ok) {
        const err = await res.json().catch(() => ({ error: 'Erro ao criar contrato' }))
        throw new Error(err.error)
      }
      return res.json() as Promise<GoContract>
    },
    onSuccess: (_, vars) => {
      const ucId = vars.consumer_unit_id as string
      queryClient.invalidateQueries({ queryKey: ['active-contract', ucId] })
    },
  })
}

// ─── Calculations ─────────────────────────────────────────────────────────────

export function useCalculation(id: string) {
  return useQuery({
    queryKey: ['calculation', id],
    queryFn: async (): Promise<GoBillingCalculation> => {
      const res = await fetch(`/api/billing/calculations/${id}`)
      if (!res.ok) throw new Error('Erro ao carregar cálculo')
      return res.json()
    },
    enabled: !!id,
  })
}
