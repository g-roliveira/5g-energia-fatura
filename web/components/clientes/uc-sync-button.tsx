'use client'

import * as React from 'react'
import { Button } from '@/components/ui/button'
import { useSyncPolling } from '@/hooks/use-sync-polling'

type SyncState = 'idle' | 'syncing' | 'polling' | 'success' | 'error'

type UcSyncButtonProps = {
  ucCode: string
  credentialId: string
  onSuccess?: () => void
}

export function UcSyncButton({ ucCode, credentialId, onSuccess }: UcSyncButtonProps) {
  const [state, setState] = React.useState<SyncState>('idle')
  const [syncRunId, setSyncRunId] = React.useState<string | null>(null)
  const [errorMsg, setErrorMsg] = React.useState<string | null>(null)

  const { data: syncRun } = useSyncPolling(state === 'polling' ? syncRunId : null)

  React.useEffect(() => {
    if (state !== 'polling' || !syncRun) return
    if (syncRun.status === 'succeeded') {
      setState('success')
      onSuccess?.()
    } else if (syncRun.status === 'failed' || syncRun.status === 'partial') {
      setState('error')
      setErrorMsg(syncRun.error_message ?? 'Falha na sincronização')
    }
  }, [syncRun, state, onSuccess])

  async function handleSync() {
    setState('syncing')
    setErrorMsg(null)
    try {
      const res = await fetch(`/api/integration/ucs/${ucCode}/sync`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ credential_id: credentialId }),
      })
      if (!res.ok) {
        const data = await res.json()
        throw new Error(data.error ?? 'Erro ao iniciar sync')
      }
      const data = await res.json()
      setSyncRunId(data.sync_run_id)
      setState('polling')
    } catch (err) {
      setState('error')
      setErrorMsg(err instanceof Error ? err.message : 'Erro desconhecido')
    }
  }

  function reset() {
    setState('idle')
    setSyncRunId(null)
    setErrorMsg(null)
  }

  if (state === 'error') {
    return (
      <div className="flex items-center gap-2">
        <span className="text-xs text-destructive">{errorMsg}</span>
        <Button variant="outline" size="sm" onClick={reset}>
          Tentar novamente
        </Button>
      </div>
    )
  }

  return (
    <Button
      size="sm"
      onClick={handleSync}
      disabled={state === 'syncing' || state === 'polling'}
    >
      {state === 'idle' && 'Sincronizar'}
      {state === 'syncing' && 'Iniciando...'}
      {state === 'polling' && 'Verificando...'}
      {state === 'success' && 'Concluído ✓'}
    </Button>
  )
}
