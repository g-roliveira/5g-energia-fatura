'use client'

import { useParams } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import Link from 'next/link'
import { HugeiconsIcon } from '@hugeicons/react'
import { ArrowLeft01Icon } from '@hugeicons/core-free-icons'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { SyncRunDetail } from '@/components/clientes/sync-run-detail'
import { useSetBreadcrumbTitle } from '@/contexts/breadcrumb'

export default function SyncRunPage() {
  const params = useParams()
  const id = params.id as string
  const ucId = params.ucId as string
  const syncId = params.syncId as string

  const { data: clientData } = useQuery({
    queryKey: ['client', id],
    queryFn: async () => {
      const res = await fetch(`/api/clients/${id}`)
      if (!res.ok) throw new Error('Falha ao carregar cliente')
      return res.json() as Promise<{ nome_razao: string }>
    },
    staleTime: 5 * 60 * 1000,
  })

  const { data: syncRun, isLoading: syncLoading } = useQuery({
    queryKey: ['sync-runs', syncId],
    queryFn: async () => {
      const res = await fetch(`/api/integration/sync-runs/${syncId}`)
      if (!res.ok) throw new Error('Falha ao carregar sync run')
      return res.json()
    },
  })

  const { data: ucs } = useQuery({
    queryKey: ['clients', id, 'ucs'],
    queryFn: async () => {
      const res = await fetch(`/api/clients/${id}/ucs`)
      if (!res.ok) throw new Error('Falha ao carregar UCs')
      return res.json() as Promise<Array<{ id: string; uc_code: string; credential?: { id: string } | null }>>
    },
  })

  const uc = ucs?.find((u) => u.id === ucId)
  const ucCode = uc?.uc_code ?? ''
  const credentialId = uc?.credential?.id

  useSetBreadcrumbTitle(id, clientData?.nome_razao)
  useSetBreadcrumbTitle(ucId, uc?.uc_code)
  useSetBreadcrumbTitle(syncId, `Sync ${syncRun?.status ?? ''}`.trim())

  return (
    <div className="flex flex-col gap-6 p-6">
      <div>
        <Button variant="outline" size="sm" asChild>
          <Link href={`/clientes/${id}/ucs`}>
            <HugeiconsIcon icon={ArrowLeft01Icon} strokeWidth={2} />
            Unidades Consumidoras
          </Link>
        </Button>
      </div>

      <div>
        <h1 className="text-2xl font-semibold">Auditoria de Sincronização</h1>
        <p className="text-sm text-muted-foreground">
          {ucCode ? `UC ${ucCode}` : 'Carregando...'}
        </p>
      </div>

      {syncLoading ? (
        <Card>
          <CardHeader>
            <Skeleton className="h-6 w-48" />
          </CardHeader>
          <CardContent className="space-y-3">
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-32 w-full rounded-md" />
          </CardContent>
        </Card>
      ) : syncRun ? (
        <SyncRunDetail
          syncRun={syncRun}
          ucCode={ucCode}
          credentialId={credentialId}
        />
      ) : (
        <Card>
          <CardContent className="py-12 text-center text-muted-foreground">
            Registro de sincronização não encontrado.
          </CardContent>
        </Card>
      )}
    </div>
  )
}
