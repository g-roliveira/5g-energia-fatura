'use client'

import Link from 'next/link'
import { useState } from 'react'
import { useRouter, useParams } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import { HugeiconsIcon } from '@hugeicons/react'
import { ArrowLeft01Icon } from '@hugeicons/core-free-icons'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { ClientForm } from '@/components/clientes/client-form'
import type { CreateClientInput } from '@/types/clientes'
import { useSetBreadcrumbTitle } from '@/contexts/breadcrumb'

export default function EditarClientePage() {
  const router = useRouter()
  const { id } = useParams<{ id: string }>()
  const [isLoading, setIsLoading] = useState(false)

  const { data: client, isLoading: isFetching } = useQuery({
    queryKey: ['client', id],
    queryFn: async () => {
      const res = await fetch(`/api/clients/${id}`)
      if (!res.ok) throw new Error('Erro ao carregar cliente')
      return res.json()
    },
  })

  useSetBreadcrumbTitle(id, client?.nome_razao)

  async function handleSubmit(data: CreateClientInput) {
    setIsLoading(true)
    try {
      const res = await fetch(`/api/clients/${id}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      })
      if (!res.ok) throw new Error('Erro ao salvar cliente')
      router.push(`/clientes/${id}`)
    } finally {
      setIsLoading(false)
    }
  }

  // Flatten nested address and commercial_data into the form's flat shape
  const defaultValues: Partial<CreateClientInput> | undefined = client
    ? {
        tipo_pessoa: client.tipo_pessoa,
        nome_razao: client.nome_razao,
        nome_fantasia: client.nome_fantasia ?? undefined,
        cpf_cnpj: client.cpf_cnpj,
        email: client.email ?? undefined,
        telefone: client.telefone ?? undefined,
        status: client.status,
        tipo_cliente: client.tipo_cliente,
        observacoes: client.observacoes ?? undefined,
        // address (nested)
        cep: client.address?.cep ?? undefined,
        logradouro: client.address?.logradouro ?? undefined,
        numero: client.address?.numero ?? undefined,
        complemento: client.address?.complemento ?? undefined,
        bairro: client.address?.bairro ?? undefined,
        cidade: client.address?.cidade ?? undefined,
        uf: client.address?.uf ?? undefined,
        // commercial_data (nested)
        tipo_contrato: client.commercial_data?.tipo_contrato ?? undefined,
        data_inicio: client.commercial_data?.data_inicio
          ? new Date(client.commercial_data.data_inicio).toISOString().slice(0, 10)
          : undefined,
        data_fim: client.commercial_data?.data_fim
          ? new Date(client.commercial_data.data_fim).toISOString().slice(0, 10)
          : undefined,
        status_contrato: client.commercial_data?.status_contrato ?? undefined,
        observacoes_comerciais: client.commercial_data?.observacoes_comerciais ?? undefined,
      }
    : undefined

  return (
    <div className="flex flex-col gap-6 p-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Editar Cliente</h1>
          <p className="text-sm text-muted-foreground">Atualize os dados do cliente</p>
        </div>
        <Button variant="outline" asChild>
          <Link href={`/clientes/${id}`}>
            <HugeiconsIcon icon={ArrowLeft01Icon} strokeWidth={2} />
            Voltar
          </Link>
        </Button>
      </div>

      {isFetching ? (
        <div className="space-y-4">
          <Skeleton className="h-8 w-64" />
          <Skeleton className="h-8 w-full" />
          <Skeleton className="h-8 w-full" />
          <Skeleton className="h-8 w-3/4" />
        </div>
      ) : (
        <ClientForm
          mode="edit"
          defaultValues={defaultValues}
          onSubmit={handleSubmit}
          isLoading={isLoading}
        />
      )}
    </div>
  )
}
