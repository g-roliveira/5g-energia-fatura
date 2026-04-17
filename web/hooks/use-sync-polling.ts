import { useQuery } from '@tanstack/react-query'
import type { GoSyncRun } from '@/types/clientes'

export function useSyncPolling(syncRunId: string | null) {
  const { data, isLoading } = useQuery<GoSyncRun>({
    queryKey: ['sync-run', syncRunId],
    queryFn: async () => {
      const res = await fetch(`/api/integration/sync-runs/${syncRunId}`)
      if (!res.ok) throw new Error('Erro ao buscar status do sync')
      return res.json()
    },
    enabled: !!syncRunId,
    refetchInterval: (query) => {
      const status = query.state.data?.status
      if (!status || status === 'running') return 2000
      return false
    },
  })

  return {
    data,
    isPolling: isLoading || data?.status === 'running',
  }
}
