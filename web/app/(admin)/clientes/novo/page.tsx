'use client'

import Link from 'next/link'
import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { HugeiconsIcon } from '@hugeicons/react'
import { ArrowLeft01Icon } from '@hugeicons/core-free-icons'
import { Button } from '@/components/ui/button'
import { ClientForm } from '@/components/clientes/client-form'
import type { CreateClientInput } from '@/types/clientes'

export default function NovoClientePage() {
  const router = useRouter()
  const [isLoading, setIsLoading] = useState(false)

  async function handleSubmit(data: CreateClientInput) {
    setIsLoading(true)
    try {
      const res = await fetch('/api/clients', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      })
      if (!res.ok) throw new Error('Erro ao criar cliente')
      const result = await res.json()
      router.push(`/clientes/${result.id}`)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="flex flex-col gap-6 p-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Novo Cliente</h1>
          <p className="text-sm text-muted-foreground">Preencha os dados para cadastrar um novo cliente</p>
        </div>
        <Button variant="outline" asChild>
          <Link href="/clientes">
            <HugeiconsIcon icon={ArrowLeft01Icon} strokeWidth={2} />
            Voltar
          </Link>
        </Button>
      </div>
      <ClientForm mode="create" onSubmit={handleSubmit} isLoading={isLoading} />
    </div>
  )
}
