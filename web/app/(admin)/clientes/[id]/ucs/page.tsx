'use client'

import { useState } from 'react'
import { useParams } from 'next/navigation'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import Link from 'next/link'
import { HugeiconsIcon } from '@hugeicons/react'
import { ArrowLeft01Icon, PlusSignIcon } from '@hugeicons/core-free-icons'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { UcList } from '@/components/clientes/uc-list'
import { UcForm } from '@/components/clientes/uc-form'
import { useSetBreadcrumbTitle } from '@/contexts/breadcrumb'
import type { CreateUcInput } from '@/types/clientes'

type ClientSummary = {
  nome_razao: string
  credentials: { id: string; label: string; documento_masked: string }[]
}

export default function UcsPage() {
  const params = useParams()
  const id = params.id as string
  const queryClient = useQueryClient()
  const [dialogOpen, setDialogOpen] = useState(false)
  const [isCreating, setIsCreating] = useState(false)

  const { data: client } = useQuery<ClientSummary>({
    queryKey: ['client', id],
    queryFn: async () => {
      const res = await fetch(`/api/clients/${id}`)
      if (!res.ok) throw new Error('Falha ao carregar cliente')
      return res.json()
    },
    staleTime: 5 * 60 * 1000,
  })

  useSetBreadcrumbTitle(id, client?.nome_razao)

  const { data: ucs, isLoading } = useQuery({
    queryKey: ['clients', id, 'ucs'],
    queryFn: async () => {
      const res = await fetch(`/api/clients/${id}/ucs`)
      if (!res.ok) throw new Error('Falha ao carregar UCs')
      return res.json()
    },
  })

  async function handleCreateUc(data: CreateUcInput) {
    setIsCreating(true)
    try {
      const res = await fetch(`/api/clients/${id}/ucs`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      })
      if (!res.ok) throw new Error('Erro ao criar UC')
      await queryClient.invalidateQueries({ queryKey: ['clients', id, 'ucs'] })
      setDialogOpen(false)
    } finally {
      setIsCreating(false)
    }
  }

  const ucList = ucs ?? []

  return (
    <div className="flex flex-col gap-6 p-6">
      <div>
        <Button variant="outline" size="sm" asChild>
          <Link href={`/clientes/${id}`}>
            <HugeiconsIcon icon={ArrowLeft01Icon} strokeWidth={2} />
            {client?.nome_razao ?? 'Cliente'}
          </Link>
        </Button>
      </div>

      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Unidades Consumidoras</h1>
          <p className="text-sm text-muted-foreground">
            UCs vinculadas a {client?.nome_razao ?? 'este cliente'}
          </p>
        </div>
        <Button onClick={() => setDialogOpen(true)}>
          <HugeiconsIcon icon={PlusSignIcon} strokeWidth={2} />
          Adicionar UC
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Lista de UCs</CardTitle>
          <CardDescription>
            {isLoading
              ? 'Carregando...'
              : `${ucList.length} UC${ucList.length !== 1 ? 's' : ''} encontrada${ucList.length !== 1 ? 's' : ''}`}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <UcList
            clientId={id}
            ucs={ucList}
            onSyncSuccess={() => queryClient.invalidateQueries({ queryKey: ['clients', id, 'ucs'] })}
          />
        </CardContent>
      </Card>

      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Adicionar UC</DialogTitle>
          </DialogHeader>
          <UcForm
            mode="create"
            credentials={client?.credentials ?? []}
            onSubmit={handleCreateUc}
            isLoading={isCreating}
          />
        </DialogContent>
      </Dialog>
    </div>
  )
}
