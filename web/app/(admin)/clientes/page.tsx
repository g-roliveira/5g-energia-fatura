'use client'

import Link from 'next/link'
import { useSearchParams } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import { Button } from '@/components/ui/button'
import { ClientTable } from '@/components/clientes/client-table'

export default function ClientesPage() {
  const searchParams = useSearchParams()
  const page = Number(searchParams.get('page') ?? 1)
  const pageSize = Number(searchParams.get('pageSize') ?? 20)
  const search = searchParams.get('search') ?? ''
  const status = searchParams.get('status') ?? ''
  const tipo_cliente = searchParams.get('tipo_cliente') ?? ''

  const { data, isLoading } = useQuery({
    queryKey: ['clients', { page, pageSize, search, status, tipo_cliente }],
    queryFn: async () => {
      const qs = new URLSearchParams({
        page: String(page),
        pageSize: String(pageSize),
        ...(search ? { search } : {}),
        ...(status ? { status } : {}),
        ...(tipo_cliente ? { tipo_cliente } : {}),
      })
      const res = await fetch(`/api/clients?${qs}`)
      if (!res.ok) throw new Error('Erro ao carregar clientes')
      return res.json() as Promise<{ data: any[]; total: number; page: number; pageSize: number }>
    },
  })

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-semibold">Clientes</h1>
        <Button asChild>
          <Link href="/clientes/novo">Novo cliente</Link>
        </Button>
      </div>

      <ClientTable
        data={data?.data ?? []}
        total={data?.total ?? 0}
        page={page}
        pageSize={pageSize}
        isLoading={isLoading}
      />
    </div>
  )
}
