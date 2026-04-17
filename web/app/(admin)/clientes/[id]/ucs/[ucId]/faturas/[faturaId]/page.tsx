'use client'

import { useParams } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import Link from 'next/link'
import { Skeleton } from '@/components/ui/skeleton'
import { InvoiceDetail } from '@/components/clientes/invoice-detail'

export default function FaturaDetailPage() {
  const params = useParams()
  const id = params.id as string
  const ucId = params.ucId as string
  const faturaId = params.faturaId as string

  const { data: invoice, isLoading } = useQuery({
    queryKey: ['invoices', faturaId],
    queryFn: async () => {
      const res = await fetch(`/api/integration/invoices/${faturaId}`)
      if (!res.ok) throw new Error('Falha ao carregar fatura')
      return res.json()
    },
  })

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
