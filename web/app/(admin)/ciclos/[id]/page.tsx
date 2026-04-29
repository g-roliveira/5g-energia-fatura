'use client'

import { useParams, useRouter } from 'next/navigation'
import { HugeiconsIcon } from '@hugeicons/react'
import {
  RefreshIcon,
  CalculatorIcon,
  CheckmarkCircle01Icon,
  PrinterIcon,
  LockPasswordIcon,
  ArrowLeft01Icon,
} from '@hugeicons/core-free-icons'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { Skeleton } from '@/components/ui/skeleton'
import { useCycle, useCycleRows, useBulkAction, useCloseCycle } from '@/hooks/use-billing'
import { CycleRowTable } from '@/components/billing/cycle-row-table'
import { toast } from '@/hooks/use-toast'

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

export default function CicloDetailPage() {
  const params = useParams()
  const router = useRouter()
  const cycleId = params.id as string

  const { data: cycle, isLoading: cycleLoading } = useCycle(cycleId)
  const { data: rowsData, isLoading: rowsLoading } = useCycleRows(cycleId)
  const bulkAction = useBulkAction()
  const closeCycle = useCloseCycle()

  const rows = rowsData?.items ?? []
  const totalUCs = cycle?.total_ucs ?? 0
  const syncedCount = cycle?.synced_count ?? 0
  const calcCount = cycle?.calculated_count ?? 0
  const approvedCount = cycle?.approved_count ?? 0

  async function handleBulk(action: string) {
    try {
      const result = await bulkAction.mutateAsync({
        cycleId,
        body: { action },
      })
      toast({
        title: 'Ação em massa iniciada',
        description: `${result.jobs_created} job(s) criado(s), ${result.jobs_skipped} pulado(s).`,
      })
    } catch (err: any) {
      toast({
        title: 'Erro',
        description: err.message,
        variant: 'destructive',
      })
    }
  }

  async function handleClose() {
    try {
      await closeCycle.mutateAsync(cycleId)
      toast({ title: 'Ciclo fechado com sucesso' })
    } catch (err: any) {
      toast({
        title: 'Erro ao fechar ciclo',
        description: err.message,
        variant: 'destructive',
      })
    }
  }

  if (cycleLoading) {
    return (
      <div className="flex flex-col gap-6 p-6">
        <Skeleton className="h-8 w-64" />
        <Skeleton className="h-24 w-full" />
        <Skeleton className="h-64 w-full" />
      </div>
    )
  }

  if (!cycle) {
    return (
      <div className="flex flex-col gap-6 p-6">
        <p className="text-muted-foreground">Ciclo não encontrado.</p>
      </div>
    )
  }

  const canSync = cycle.status === 'open' || cycle.status === 'syncing'
  const canCalc = cycle.status === 'open' || cycle.status === 'syncing' || cycle.status === 'processing'
  const canApprove = cycle.status !== 'closed'
  const canClose = cycle.status === 'approved' || cycle.status === 'review'

  return (
    <div className="flex flex-col gap-6 p-6">
      {/* Back + Title */}
      <div className="flex items-center gap-2">
        <Button variant="ghost" size="icon" onClick={() => router.push('/ciclos')}>
          <HugeiconsIcon icon={ArrowLeft01Icon} strokeWidth={2} className="size-4" />
        </Button>
        <div>
          <h1 className="text-2xl font-semibold">
            {monthNames[cycle.month - 1]} / {cycle.year}
          </h1>
          <p className="text-sm text-muted-foreground">
            Ciclo de faturamento — {new Date(cycle.reference_date).toLocaleDateString('pt-BR')}
          </p>
        </div>
        <Badge variant={statusBadgeVariant(cycle.status) as any} className="ml-auto">
          {statusLabel(cycle.status)}
        </Badge>
      </div>

      {/* Stats cards */}
      <div className="grid grid-cols-4 gap-4">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Total UCs</CardDescription>
            <CardTitle className="text-3xl">{totalUCs}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Sincronizadas</CardDescription>
            <CardTitle className="text-3xl">{syncedCount}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Calculadas</CardDescription>
            <CardTitle className="text-3xl">{calcCount}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Aprovadas</CardDescription>
            <CardTitle className="text-3xl">{approvedCount}</CardTitle>
          </CardHeader>
        </Card>
      </div>

      {/* Progress bar */}
      <div className="space-y-1">
        <div className="flex justify-between text-xs text-muted-foreground">
          <span>Progresso</span>
          <span>{totalUCs > 0 ? Math.round((approvedCount / totalUCs) * 100) : 0}% aprovado</span>
        </div>
        <div className="h-2 w-full rounded-full bg-muted overflow-hidden">
          <div
            className="h-full bg-primary transition-all"
            style={{
              width: `${totalUCs > 0 ? (approvedCount / totalUCs) * 100 : 0}%`,
            }}
          />
        </div>
      </div>

      {/* Bulk actions */}
      <div className="flex items-center gap-2 flex-wrap">
        <Button
          variant="outline"
          onClick={() => handleBulk('sync')}
          disabled={!canSync || bulkAction.isPending}
        >
          <HugeiconsIcon icon={RefreshIcon} strokeWidth={2} />
          Sincronizar
        </Button>
        <Button
          variant="outline"
          onClick={() => handleBulk('recalculate')}
          disabled={!canCalc || bulkAction.isPending}
        >
          <HugeiconsIcon icon={CalculatorIcon} strokeWidth={2} />
          Calcular
        </Button>
        <Button
          variant="outline"
          onClick={() => handleBulk('approve')}
          disabled={!canApprove || bulkAction.isPending}
        >
          <HugeiconsIcon icon={CheckmarkCircle01Icon} strokeWidth={2} />
          Aprovar
        </Button>
        <Button
          variant="outline"
          onClick={() => handleBulk('generate_pdf')}
          disabled={bulkAction.isPending}
        >
          <HugeiconsIcon icon={PrinterIcon} strokeWidth={2} />
          Gerar PDFs
        </Button>
        <Separator orientation="vertical" className="h-6 mx-1" />
        <Button
          variant="default"
          onClick={handleClose}
          disabled={!canClose || closeCycle.isPending}
        >
          <HugeiconsIcon icon={LockPasswordIcon} strokeWidth={2} />
          Fechar ciclo
        </Button>
      </div>

      {/* Rows table */}
      <Card>
        <CardHeader>
          <CardTitle>Unidades Consumidoras</CardTitle>
          <CardDescription>
            {rowsLoading ? 'Carregando...' : `${rows.length} unidade(s)`}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <CycleRowTable rows={rows} isLoading={rowsLoading} />
        </CardContent>
      </Card>
    </div>
  )
}
