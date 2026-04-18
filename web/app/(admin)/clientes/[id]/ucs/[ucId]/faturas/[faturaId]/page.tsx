'use client'

import { useParams } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import Link from 'next/link'
import { Skeleton } from '@/components/ui/skeleton'
import { InvoiceDetail } from '@/components/clientes/invoice-detail'
import { useSetBreadcrumbTitle } from '@/contexts/breadcrumb'

export default function FaturaDetailPage() {
  const params = useParams()
  const id = params.id as string
  const ucId = params.ucId as string
  const faturaId = params.faturaId as string

  const { data: clientData } = useQuery({
    queryKey: ['client', id],
    queryFn: async () => {
      const res = await fetch(`/api/clients/${id}`)
      if (!res.ok) throw new Error('Falha ao carregar cliente')
      return res.json() as Promise<{ nome_razao: string }>
    },
    staleTime: 5 * 60 * 1000,
  })

  const { data: ucs } = useQuery({
    queryKey: ['clients', id, 'ucs'],
    queryFn: async () => {
      const res = await fetch(`/api/clients/${id}/ucs`)
      if (!res.ok) throw new Error('Falha ao carregar UCs')
      return res.json() as Promise<Array<{ id: string; uc_code: string }>>
    },
    staleTime: 5 * 60 * 1000,
  })

  const { data: invoice, isLoading } = useQuery({
    queryKey: ['invoices', faturaId],
    queryFn: async () => {
      const res = await fetch(`/api/integration/invoices/${faturaId}`)
      if (!res.ok) throw new Error('Falha ao carregar fatura')
      return res.json()
    },
  })

  const uc = ucs?.find((u) => u.id === ucId)

  useSetBreadcrumbTitle(id, clientData?.nome_razao)
  useSetBreadcrumbTitle(ucId, uc?.uc_code)
  useSetBreadcrumbTitle(faturaId, invoice?.numero_fatura)

  return (
    <div className="p-6">
      <div className="mb-6">
        <div className="mb-1 text-sm text-muted-foreground">
          <Link href={`/clientes/${id}/ucs/${ucId}/faturas`} className="hover:underline">
            ← Faturas
          </Link>
        </div>
        <h1 className="text-2xl font-semibold">Detalhe da Fatura</h1>
      </div>

      {isLoading ? (
        <div className="space-y-4">
          <Skeleton className="h-8 w-48" />
          <Skeleton className="h-64 w-full rounded-md" />
        </div>
      ) : invoice ? (
        <InvoiceDetail invoice={invoice} />
      ) : null}
    </div>
  )
}
