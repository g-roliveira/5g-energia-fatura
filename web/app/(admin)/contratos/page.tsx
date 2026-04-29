'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { HugeiconsIcon } from '@hugeicons/react'
import { PlusSignIcon, ContractsIcon } from '@hugeicons/core-free-icons'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { useQuery } from '@tanstack/react-query'
import { useContracts } from '@/hooks/use-billing'
import type { GoContract } from '@/types/billing'

function statusLabel(status: string) {
  const map: Record<string, string> = { active: 'Ativo', ended: 'Encerrado', draft: 'Rascunho' }
  return map[status] ?? status
}

function statusBadgeVariant(status: string) {
  const map: Record<string, string> = { active: 'default', ended: 'secondary', draft: 'outline' }
  return map[status] ?? 'secondary'
}

export default function ContratosPage() {
  const router = useRouter()
  const [selectedClient, setSelectedClient] = useState<string>('')
  const [selectedUC, setSelectedUC] = useState<string>('')

  const { data: clients } = useQuery({
    queryKey: ['clients-minimal'],
    queryFn: async () => {
      const res = await fetch('/api/clients?page=1&pageSize=100')
      if (!res.ok) throw new Error('Erro ao carregar clientes')
      return res.json() as Promise<{ data: Array<{ id: string; nome_razao: string }> }>
    },
  })

  const { data: ucs } = useQuery({
    queryKey: ['ucs', selectedClient],
    queryFn: async () => {
      const res = await fetch(`/api/clients/${selectedClient}/ucs`)
      if (!res.ok) throw new Error('Erro ao carregar UCs')
      return res.json() as Promise<Array<{ id: string; uc_code: string; apelido?: string }>>
    },
    enabled: !!selectedClient,
  })

  const { data: contractsData, isLoading: contractsLoading } = useContracts(selectedUC)

  return (
    <div className="flex flex-col gap-6 p-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Contratos</h1>
          <p className="text-sm text-muted-foreground">Gerencie contratos de faturamento por UC</p>
        </div>
        <Button
          onClick={() => router.push('/contratos/novo')}
          disabled={!selectedUC}
        >
          <HugeiconsIcon icon={PlusSignIcon} strokeWidth={2} />
          Novo contrato
        </Button>
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex items-end gap-4">
            <div className="w-72 space-y-2">
              <label className="text-sm font-medium">Cliente</label>
              <Select value={selectedClient} onValueChange={(v) => { setSelectedClient(v); setSelectedUC('') }}>
                <SelectTrigger>
                  <SelectValue placeholder="Selecione um cliente" />
                </SelectTrigger>
                <SelectContent>
                  {clients?.data.map((c) => (
                    <SelectItem key={c.id} value={c.id}>{c.nome_razao}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="w-72 space-y-2">
              <label className="text-sm font-medium">Unidade Consumidora</label>
              <Select value={selectedUC} onValueChange={setSelectedUC} disabled={!selectedClient}>
                <SelectTrigger>
                  <SelectValue placeholder="Selecione uma UC" />
                </SelectTrigger>
                <SelectContent>
                  {ucs?.map((uc) => (
                    <SelectItem key={uc.id} value={uc.id}>
                      {uc.uc_code} {uc.apelido ? `— ${uc.apelido}` : ''}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Contracts list */}
      <Card>
        <CardHeader>
          <CardTitle>Histórico de contratos</CardTitle>
          <CardDescription>
            {contractsLoading
              ? 'Carregando...'
              : selectedUC
                ? `${contractsData?.count ?? 0} contrato(s) encontrado(s)`
                : 'Selecione um cliente e uma UC para ver os contratos'}
          </CardDescription>
        </CardHeader>
        <CardContent>
          {contractsLoading ? (
            <div className="space-y-2">
              {Array.from({ length: 3 }).map((_, i) => (
                <Skeleton key={i} className="h-12 w-full" />
              ))}
            </div>
          ) : !selectedUC ? (
            <div className="flex flex-col items-center justify-center py-12 text-muted-foreground">
              <HugeiconsIcon icon={ContractsIcon} strokeWidth={2} className="size-12 mb-4 opacity-40" />
              <p className="text-sm">Selecione um cliente e uma UC para visualizar os contratos.</p>
            </div>
          ) : contractsData?.items.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-12 text-muted-foreground">
              <HugeiconsIcon icon={ContractsIcon} strokeWidth={2} className="size-12 mb-4 opacity-40" />
              <p className="text-sm">Nenhum contrato encontrado para esta UC.</p>
              <Button
                variant="link"
                onClick={() => router.push('/contratos/novo')}
                className="mt-2"
              >
                Criar primeiro contrato
              </Button>
            </div>
          ) : (
            <div className="divide-y">
              {contractsData?.items.map((c: GoContract) => (
                <div
                  key={c.id}
                  className="flex items-center justify-between py-4 px-2 hover:bg-muted/50 rounded transition-colors"
                >
                  <div className="flex items-center gap-4">
                    <div className="flex size-10 items-center justify-center rounded-lg bg-primary/10 text-primary">
                      <HugeiconsIcon icon={ContractsIcon} strokeWidth={2} className="size-5" />
                    </div>
                    <div>
                      <p className="font-medium">
                        Desconto {Math.round(Number(c.desconto_percentual) * 100)}%
                        {' '}
                        {c.ip_faturamento_mode === 'fixed'
                          ? `— IP R$ ${c.ip_faturamento_valor}`
                          : `— IP ${Math.round(Number(c.ip_faturamento_percent) * 100)}%`}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        Vigência: {c.vigencia_inicio}
                        {c.vigencia_fim ? ` → ${c.vigencia_fim}` : ' (atual)'}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    {c.bandeira_com_desconto && (
                      <Badge variant="outline">Bandeira c/ desc.</Badge>
                    )}
                    {c.custo_disponibilidade_sempre_cobrado && (
                      <Badge variant="outline">Disp. sempre</Badge>
                    )}
                    <Badge variant={statusBadgeVariant(c.status) as any}>
                      {statusLabel(c.status)}
                    </Badge>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
