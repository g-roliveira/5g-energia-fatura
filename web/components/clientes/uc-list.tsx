'use client'

import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { StatusBadge } from './status-badge'
import { UcSyncButton } from './uc-sync-button'

type UcListItem = {
  id: string
  uc_code: string
  apelido?: string | null
  distribuidora?: string | null
  ativa: boolean
  credential?: { id: string; label: string; documento_masked: string } | null
  lastSyncStatus?: 'running' | 'succeeded' | 'partial' | 'failed'
}

type UcListProps = {
  clientId: string
  ucs: UcListItem[]
  onSyncSuccess?: () => void
}

export function UcList({ clientId, ucs, onSyncSuccess }: UcListProps) {
  const router = useRouter()

  if (ucs.length === 0) {
    return (
      <div className="py-12 text-center text-muted-foreground">
        <p>Nenhuma UC vinculada.</p>
        <p className="mt-1 text-sm">Adicione uma UC para começar a sincronizar.</p>
      </div>
    )
  }

  return (
    <div className="space-y-3">
      {ucs.map((uc) => (
        <div
          key={uc.id}
          className="flex items-center justify-between rounded-lg border p-4"
        >
          <div className="space-y-1">
            <div className="flex items-center gap-2">
              <span className="font-mono text-sm font-medium">{uc.uc_code}</span>
              {uc.apelido && (
                <span className="text-sm text-muted-foreground">— {uc.apelido}</span>
              )}
              <StatusBadge status={uc.ativa ? 'ativo' : 'inativo'} />
            </div>
            {uc.distribuidora && (
              <p className="text-xs text-muted-foreground">{uc.distribuidora}</p>
            )}
            {uc.lastSyncStatus && (
              <div className="flex items-center gap-1 text-xs text-muted-foreground">
                <span>Último sync:</span>
                <StatusBadge status={uc.lastSyncStatus} />
              </div>
            )}
          </div>

          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => router.push(`/clientes/${clientId}/ucs/${uc.id}/faturas`)}
            >
              Faturas
            </Button>
            {uc.credential && (
              <UcSyncButton
                ucCode={uc.uc_code}
                credentialId={uc.credential.id}
                onSuccess={onSyncSuccess}
              />
            )}
          </div>
        </div>
      ))}
    </div>
  )
}
