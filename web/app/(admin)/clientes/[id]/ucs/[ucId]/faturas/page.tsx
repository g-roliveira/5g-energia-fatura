'use client'

import { useParams } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import Link from 'next/link'
import { InvoiceTable } from '@/components/clientes/invoice-table'

export default function FaturasPage() {
  const params = useParams()
  const id = params.id as string
  const ucId = params.ucId as string

  const { data: ucs } = useQuery({
    queryKey: ['clients', id, 'ucs'],
    queryFn: async () => {
      const res = await fetch(`/api/clients/${id}/ucs`)
      if (!res.ok) throw new Error('Falha ao carregar UCs')
      return res.json() as Promise<Array<{ id: string; uc_code: string }>>
    },
  })

  const uc = ucs?.find((u) => u.id === ucId)
  const ucCode = uc?.uc_code

  const { data: invoices, isLoading } = useQuery({
    queryKey: ['ucs', ucCode, 'invoices'],
    enabled: !!ucCode,
    queryFn: async () => {
      const res = await fetch(`/api/integration/ucs/${ucCode}/invoices`)
      if (!res.ok) throw new Error('Falha ao carregar faturas')
      return res.json()
    },
  })

  const invoiceList = invoices ?? []

  return (
    <div className="p-6">
      <div className="mb-6">
        <div className="mb-1 text-sm text-muted-foreground">
          <Link href={`/clientes/${id}`} className="hover:underline">
            Cliente
          </Link>
          {' / '}
          <Link href={`/clientes/${id}/ucs`} className="hover:underline">
            ← UCs
          </Link>
        </div>
        <h1 className="text-2xl font-semibold">Faturas</h1>
      </div>

      {!isLoading && invoiceList.length === 0 ? (
        <div className="py-12 text-center text-muted-foreground">
          <p>Nenhuma fatura encontrada.</p>
          <p className="mt-2 text-sm">
            <Link
              href={`/clientes/${id}/ucs`}
              className="underline hover:text-foreground"
            >
              Sincronizar UC
            </Link>
          </p>
        </div>
      ) : (
        <InvoiceTable
          data={invoiceList}
          isLoading={isLoading}
          clientId={id}
          ucId={ucId}
        />
      )}
    </div>
  )
}
