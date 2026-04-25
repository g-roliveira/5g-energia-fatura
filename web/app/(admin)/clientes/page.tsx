'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useSearchParams } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import { HugeiconsIcon } from '@hugeicons/react'
import { PlusSignIcon, Download01Icon } from '@hugeicons/core-free-icons'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { ClientTable } from '@/components/clientes/client-table'
import { ClientFormMinimal } from '@/components/clientes/client-form-minimal'
import type { CreateClientMinimalInput } from '@/types/clientes'

export default function ClientesPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const page = Number(searchParams.get('page') ?? 1)
  const pageSize = Number(searchParams.get('pageSize') ?? 20)
  const search = searchParams.get('search') ?? ''
  const status = searchParams.get('status') ?? ''
  const tipo_cliente = searchParams.get('tipo_cliente') ?? ''

  const [dialogOpen, setDialogOpen] = useState(false)
  const [isCreating, setIsCreating] = useState(false)

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

  async function handleCreateClient(data: CreateClientMinimalInput) {
    setIsCreating(true)
    try {
      const res = await fetch('/api/clients', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      })
      if (!res.ok) throw new Error('Erro ao criar cliente')
      const result = await res.json()
      setDialogOpen(false)
      router.push(`/clientes/${result.id}`)
    } finally {
      setIsCreating(false)
    }
  }

  return (
    <div className="flex flex-col gap-6 p-6">
      {/* Page header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Clientes</h1>
          <p className="text-sm text-muted-foreground">Gerencie seus clientes cadastrados</p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" disabled>
            <HugeiconsIcon icon={Download01Icon} strokeWidth={2} />
            Importar CSV
            <Badge variant="secondary" className="ml-1 text-[10px]">Em breve</Badge>
          </Button>
          <Button onClick={() => setDialogOpen(true)}>
            <HugeiconsIcon icon={PlusSignIcon} strokeWidth={2} />
            Novo cliente
          </Button>
        </div>
      </div>

      {/* Table card */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Lista de clientes</CardTitle>
              <CardDescription>
                {isLoading ? 'Carregando...' : `${data?.total ?? 0} cliente${(data?.total ?? 0) !== 1 ? 's' : ''} encontrado${(data?.total ?? 0) !== 1 ? 's' : ''}`}
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <ClientTable
            data={data?.data ?? []}
            total={data?.total ?? 0}
            page={page}
            pageSize={pageSize}
            isLoading={isLoading}
          />
        </CardContent>
      </Card>

      {/* New client dialog */}
      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Novo Cliente</DialogTitle>
          </DialogHeader>
          <ClientFormMinimal onSubmit={handleCreateClient} isLoading={isCreating} />
        </DialogContent>
      </Dialog>
    </div>
  )
}
