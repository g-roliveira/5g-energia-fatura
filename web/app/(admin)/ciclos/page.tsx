'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useSearchParams } from 'next/navigation'
import { HugeiconsIcon } from '@hugeicons/react'
import { PlusSignIcon, InvoiceIcon } from '@hugeicons/core-free-icons'
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
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useCycles, useCreateCycle } from '@/hooks/use-billing'
import type { GoBillingCycle } from '@/types/billing'

const monthNames = [
  'Janeiro', 'Fevereiro', 'Março', 'Abril', 'Maio', 'Junho',
  'Julho', 'Agosto', 'Setembro', 'Outubro', 'Novembro', 'Dezembro',
]

function statusLabel(status: string) {
  const map: Record<string, string> = {
    open: 'Aberto',
    syncing: 'Sincronizando',
    processing: 'Processando',
    review: 'Revisão',
    approved: 'Aprovado',
    closed: 'Fechado',
  }
  return map[status] ?? status
}

function statusBadgeVariant(status: string) {
  const map: Record<string, string> = {
    open: 'secondary',
    syncing: 'outline',
    processing: 'outline',
    review: 'destructive',
    approved: 'default',
    closed: 'secondary',
  }
  return map[status] ?? 'secondary'
}

export default function CiclosPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const yearFilter = searchParams.get('year') ? Number(searchParams.get('year')) : undefined
  const statusFilter = searchParams.get('status') ?? undefined

  const [dialogOpen, setDialogOpen] = useState(false)
  const [newYear, setNewYear] = useState(new Date().getFullYear())
  const [newMonth, setNewMonth] = useState(new Date().getMonth() + 1)

  const { data, isLoading } = useCycles(yearFilter, statusFilter)
  const createCycle = useCreateCycle()

  async function handleCreate() {
    try {
      const cycle = await createCycle.mutateAsync({
        year: newYear,
        month: newMonth,
        include_all_active: true,
      })
      setDialogOpen(false)
      router.push(`/ciclos/${cycle.id}`)
    } catch {
      // erro mostrado pelo toast se houver
    }
  }

  return (
    <div className="flex flex-col gap-6 p-6">
      {/* Page header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Ciclos de Faturamento</h1>
          <p className="text-sm text-muted-foreground">
            Gerencie os ciclos de faturamento mensal
          </p>
        </div>
        <Button onClick={() => setDialogOpen(true)}>
          <HugeiconsIcon icon={PlusSignIcon} strokeWidth={2} />
          Novo ciclo
        </Button>
      </div>

      {/* Filters */}
      <div className="flex items-center gap-2">
        <Select
          value={statusFilter ?? 'all'}
          onValueChange={(v) => {
            const url = new URL(window.location.href)
            if (v === 'all') url.searchParams.delete('status')
            else url.searchParams.set('status', v)
            router.push(url.pathname + url.search)
          }}
        >
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Todos os status</SelectItem>
            <SelectItem value="open">Aberto</SelectItem>
            <SelectItem value="syncing">Sincronizando</SelectItem>
            <SelectItem value="processing">Processando</SelectItem>
            <SelectItem value="review">Revisão</SelectItem>
            <SelectItem value="approved">Aprovado</SelectItem>
            <SelectItem value="closed">Fechado</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Cycles list */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Lista de ciclos</CardTitle>
              <CardDescription>
                {isLoading
                  ? 'Carregando...'
                  : `${data?.count ?? 0} ciclo${(data?.count ?? 0) !== 1 ? 's' : ''} encontrado${(data?.count ?? 0) !== 1 ? 's' : ''}`}
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-2">
              {Array.from({ length: 5 }).map((_, i) => (
                <div key={i} className="h-12 bg-muted rounded animate-pulse" />
              ))}
            </div>
          ) : data?.items.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-12 text-muted-foreground">
              <HugeiconsIcon icon={InvoiceIcon} strokeWidth={2} className="size-12 mb-4 opacity-40" />
              <p className="text-sm">Nenhum ciclo encontrado.</p>
              <p className="text-xs">Clique em "Novo ciclo" para começar.</p>
            </div>
          ) : (
            <div className="divide-y">
              {data?.items.map((cycle: GoBillingCycle) => (
                <div
                  key={cycle.id}
                  className="flex items-center justify-between py-4 cursor-pointer hover:bg-muted/50 rounded px-2 transition-colors"
                  onClick={() => router.push(`/ciclos/${cycle.id}`)}
                >
                  <div className="flex items-center gap-4">
                    <div className="flex size-10 items-center justify-center rounded-lg bg-primary/10 text-primary">
                      <HugeiconsIcon icon={InvoiceIcon} strokeWidth={2} className="size-5" />
                    </div>
                    <div>
                      <p className="font-medium">
                        {monthNames[cycle.month - 1]} / {cycle.year}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        Ref: {new Date(cycle.reference_date).toLocaleDateString('pt-BR')}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-4">
                    <div className="flex gap-3 text-xs text-muted-foreground">
                      <span>UCs: {cycle.total_ucs ?? 0}</span>
                      <span>Sync: {cycle.synced_count ?? 0}</span>
                      <span>Calc: {cycle.calculated_count ?? 0}</span>
                      <span>Aprov: {cycle.approved_count ?? 0}</span>
                    </div>
                    <Badge variant={statusBadgeVariant(cycle.status) as any}>
                      {statusLabel(cycle.status)}
                    </Badge>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* New cycle dialog */}
      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>Novo Ciclo de Faturamento</DialogTitle>
          </DialogHeader>
          <div className="space-y-4 py-2">
            <div className="space-y-2">
              <Label htmlFor="month">Mês</Label>
              <Select value={String(newMonth)} onValueChange={(v) => setNewMonth(Number(v))}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {monthNames.map((name, idx) => (
                    <SelectItem key={idx + 1} value={String(idx + 1)}>
                      {name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="year">Ano</Label>
              <Input
                id="year"
                type="number"
                value={newYear}
                onChange={(e) => setNewYear(Number(e.target.value))}
                min={2000}
                max={2100}
              />
            </div>
            <Button
              className="w-full"
              onClick={handleCreate}
              disabled={createCycle.isPending}
            >
              {createCycle.isPending ? 'Criando...' : 'Criar ciclo'}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}
