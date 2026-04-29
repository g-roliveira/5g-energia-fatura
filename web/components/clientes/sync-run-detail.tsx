'use client'

import * as React from 'react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Alert, AlertDescription } from '@/components/ui/alert'
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
      <Card>
        <CardHeader>
          <CardTitle>Resumo</CardTitle>
        </CardHeader>
        <CardContent className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <StatusBadge status={syncRun.status} />
            <span className="text-sm text-muted-foreground">{createdAt}</span>
          </div>
          <span className="font-mono text-xs text-muted-foreground">{syncRun.id}</span>
        </CardContent>
      </Card>

      {syncRun.error_message && (
        <Alert variant="destructive">
          <AlertDescription>
            <p className="font-medium">Erro</p>
            <p className="mt-1">{syncRun.error_message}</p>
          </AlertDescription>
        </Alert>
      )}

      {syncRun.raw_response !== undefined && (
        <Card>
          <CardHeader className="cursor-pointer" onClick={() => setRawOpen((o) => !o)}>
            <div className="flex items-center justify-between">
              <CardTitle>Resposta bruta</CardTitle>
              <Button variant="ghost" size="sm" onClick={(e) => { e.stopPropagation(); setRawOpen((o) => !o) }}>
                {rawOpen ? '▲' : '▼'}
              </Button>
            </div>
          </CardHeader>
          {rawOpen && (
            <CardContent className="p-0">
              <pre className="overflow-x-auto rounded-b-lg bg-muted/30 p-4 text-xs">
                {JSON.stringify(syncRun.raw_response, null, 2)}
              </pre>
            </CardContent>
          )}
        </Card>
      )}

      {syncRun.status === 'failed' && onReprocess && (
        <Button variant="outline" onClick={onReprocess}>
          Reprocessar
        </Button>
      )}
    </div>
  )
}
