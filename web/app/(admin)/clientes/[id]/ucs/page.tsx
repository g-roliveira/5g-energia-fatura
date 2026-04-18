'use client'

import { useParams } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import Link from 'next/link'
import { Skeleton } from '@/components/ui/skeleton'
import { UcList } from '@/components/clientes/uc-list'
import { useSetBreadcrumbTitle } from '@/contexts/breadcrumb'

export default function UcsPage() {
  const params = useParams()
  const id = params.id as string

  const { data: clientData } = useQuery({
    queryKey: ['client', id],
    queryFn: async () => {
      const res = await fetch(`/api/clients/${id}`)
      if (!res.ok) throw new Error('Falha ao carregar cliente')
      return res.json() as Promise<{ nome_razao: string }>
    },
    staleTime: 5 * 60 * 1000,
  })

  useSetBreadcrumbTitle(id, clientData?.nome_razao)

  const { data: ucs, isLoading, refetch } = useQuery({
    queryKey: ['clients', id, 'ucs'],
    queryFn: async () => {
      const res = await fetch(`/api/clients/${id}/ucs`)
      if (!res.ok) throw new Error('Falha ao carregar UCs')
      return res.json()
    },
  })

  return (
    <div className="p-6">
      <div className="mb-6">
        <div className="mb-1 text-sm text-muted-foreground">
          <Link href={`/clientes/${id}`} className="hover:underline">
            ← Cliente
          </Link>
        </div>
        <h1 className="text-2xl font-semibold">Unidades Consumidoras</h1>
      </div>

      {isLoading ? (
        <div className="space-y-3">
          {Array.from({ length: 3 }).map((_, i) => (
            <Skeleton key={i} className="h-20 w-full rounded-lg" />
          ))}
        </div>
      ) : (
        <UcList
          clientId={id}
          ucs={ucs ?? []}
          onSyncSuccess={() => refetch()}
        />
      )}
    </div>
  )
}
