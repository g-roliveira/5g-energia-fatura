'use client'

import { useParams } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import Link from 'next/link'
import { HugeiconsIcon } from '@hugeicons/react'
import { ArrowLeft01Icon } from '@hugeicons/core-free-icons'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
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
  useSetBreadcrumbTitle(faturaId, invoice?.billing_record?.numero_fatura ?? invoice?.numero_fatura)

  return (
    <div className="flex flex-col gap-6 p-6">
      <div>
        <Button variant="outline" size="sm" asChild>
          <Link href={`/clientes/${id}/ucs/${ucId}/faturas`}>
            <HugeiconsIcon icon={ArrowLeft01Icon} strokeWidth={2} />
            Faturas
          </Link>
        </Button>
      </div>

      <div>
        <h1 className="text-2xl font-semibold">Detalhe da Fatura</h1>
        <p className="text-sm text-muted-foreground">
          {uc?.uc_code ? `UC ${uc.uc_code}` : 'Carregando...'}
        </p>
      </div>

      {isLoading ? (
        <Card>
          <CardHeader>
            <Skeleton className="h-6 w-48" />
          </CardHeader>
          <CardContent className="space-y-3">
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-3/4" />
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-2/3" />
          </CardContent>
        </Card>
      ) : invoice ? (
        <InvoiceDetail invoice={invoice} />
      ) : (
        <Card>
          <CardContent className="py-12 text-center text-muted-foreground">
            Fatura não encontrada.
          </CardContent>
        </Card>
      )}
    </div>
  )
}
