'use client'

import * as React from 'react'
import { Button } from '@/components/ui/button'
import { StatusBadge } from './status-badge'
import type { GoSyncRun } from '@/types/clientes'

type SyncRunDetailProps = {
  syncRun: GoSyncRun
  ucCode: string
  credentialId?: string
  onReprocess?: () => void
}

export function SyncRunDetail({ syncRun, onReprocess }: SyncRunDetailProps) {
  const [rawOpen, setRawOpen] = React.useState(false)

  const createdAt = React.useMemo(() => {
    try {
      return new Intl.DateTimeFormat('pt-BR', {
        dateStyle: 'short',
        timeStyle: 'medium',
      }).format(new Date(syncRun.created_at))
    } catch {
      return syncRun.created_at
    }
  }, [syncRun.created_at])

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <StatusBadge status={syncRun.status} />
          <span className="text-sm text-muted-foreground">{createdAt}</span>
        </div>
        <span className="font-mono text-xs text-muted-foreground">{syncRun.id}</span>
      </div>

      {/* Error */}
      {syncRun.error_message && (
        <div className="rounded-md border border-destructive/30 bg-destructive/5 p-3 text-sm text-destructive">
          <p className="font-medium">Erro</p>
          <p className="mt-1">{syncRun.error_message}</p>
        </div>
      )}

      {/* Raw response */}
      {syncRun.raw_response !== undefined && (
        <div className="rounded-md border">
          <button
            onClick={() => setRawOpen((o) => !o)}
            className="flex w-full items-center justify-between p-3 text-sm font-medium hover:bg-muted/50"
          >
            <span>Resposta bruta</span>
            <span className="text-muted-foreground">{rawOpen ? '▲' : '▼'}</span>
          </button>
          {rawOpen && (
            <pre className="overflow-x-auto border-t bg-muted/30 p-3 text-xs">
              {JSON.stringify(syncRun.raw_response, null, 2)}
            </pre>
          )}
        </div>
      )}

      {/* Reprocess CTA */}
      {syncRun.status === 'failed' && onReprocess && (
        <Button variant="outline" onClick={onReprocess}>
          Reprocessar
        </Button>
      )}
    </div>
  )
}
