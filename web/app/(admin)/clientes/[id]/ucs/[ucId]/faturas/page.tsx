'use client'

import { useParams } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import Link from 'next/link'
import { HugeiconsIcon } from '@hugeicons/react'
import { ArrowLeft01Icon } from '@hugeicons/core-free-icons'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { InvoiceTable } from '@/components/clientes/invoice-table'
import { useSetBreadcrumbTitle } from '@/contexts/breadcrumb'

export default function FaturasPage() {
  const params = useParams()
  const id = params.id as string
  const ucId = params.ucId as string

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
  })

  const uc = ucs?.find((u) => u.id === ucId)
  const ucCode = uc?.uc_code

  useSetBreadcrumbTitle(id, clientData?.nome_razao)
  useSetBreadcrumbTitle(ucId, ucCode)

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
        <h1 className="text-2xl font-semibold">Faturas</h1>
        <p className="text-sm text-muted-foreground">
          {ucCode ? `UC ${ucCode}` : 'Carregando...'}
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Histórico de faturas</CardTitle>
          <CardDescription>
            {isLoading
              ? 'Carregando...'
              : `${invoiceList.length} fatura${invoiceList.length !== 1 ? 's' : ''} encontrada${invoiceList.length !== 1 ? 's' : ''}`}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <InvoiceTable
            data={invoiceList}
            isLoading={isLoading}
            clientId={id}
            ucId={ucId}
          />
        </CardContent>
      </Card>
    </div>
  )
}
