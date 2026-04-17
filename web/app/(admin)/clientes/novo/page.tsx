'use client'

import Link from 'next/link'
import { useState } from 'react'
import { useRouter } from 'next/navigation'
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
    <div className="p-6 max-w-4xl">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-semibold">Novo Cliente</h1>
        <Button variant="outline" asChild>
          <Link href="/clientes">Voltar</Link>
        </Button>
      </div>

      <ClientForm mode="create" onSubmit={handleSubmit} isLoading={isLoading} />
    </div>
  )
}
