'use client'

import { useRouter } from 'next/navigation'
import { HugeiconsIcon } from '@hugeicons/react'
import { FileViewIcon, Alert01Icon, PrinterIcon } from '@hugeicons/core-free-icons'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Skeleton } from '@/components/ui/skeleton'
import type { GoCycleRow } from '@/types/billing'

interface CycleRowTableProps {
  rows: GoCycleRow[]
  isLoading: boolean
}

function statusBadge(status: string) {
  const map: Record<string, string> = {
    pending: 'secondary',
    synced: 'default',
    calculated: 'outline',
    approved: 'default',
    open: 'secondary',
    syncing: 'outline',
    processing: 'outline',
    review: 'destructive',
    closed: 'secondary',
    draft: 'secondary',
    needs_review: 'destructive',
  }
  return map[status] ?? 'secondary'
}

export function CycleRowTable({ rows, isLoading }: CycleRowTableProps) {
  const router = useRouter()

  if (isLoading) {
    return (
      <div className="space-y-2">
        {Array.from({ length: 5 }).map((_, i) => (
          <Skeleton key={i} className="h-10 w-full" />
        ))}
      </div>
    )
  }

  if (rows.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 text-muted-foreground">
        <p className="text-sm">Nenhuma unidade consumidora encontrada neste ciclo.</p>
      </div>
    )
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>UC</TableHead>
          <TableHead>Cliente</TableHead>
          <TableHead>Sync</TableHead>
          <TableHead>Cálculo</TableHead>
          <TableHead className="text-right">Valor Ref. (R$)</TableHead>
          <TableHead className="text-right">Repasse 5G (R$)</TableHead>
          <TableHead className="text-right">Economia</TableHead>
          <TableHead className="w-[100px]">Ações</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {rows.map((row) => (
          <TableRow key={row.consumer_unit_id}>
            <TableCell className="font-medium">{row.uc_code}</TableCell>
            <TableCell>{row.customer_name}</TableCell>
            <TableCell>
              <Badge variant={statusBadge(row.sync_status) as any}>{row.sync_status}</Badge>
            </TableCell>
            <TableCell>
              <div className="flex items-center gap-1">
                <Badge variant={statusBadge(row.calculation_status ?? '') as any}>
                  {row.calculation_status ?? '—'}
                </Badge>
                {row.needs_review_reasons && row.needs_review_reasons.length > 0 && (
                  <HugeiconsIcon icon={Alert01Icon} strokeWidth={2} className="size-4 text-amber-500" />
                )}
              </div>
            </TableCell>
            <TableCell className="text-right font-mono text-sm">
              {row.valor_azi_sem_desconto ? `R$ ${row.valor_azi_sem_desconto}` : '—'}
            </TableCell>
            <TableCell className="text-right font-mono text-sm">
              {row.valor_azi_com_desconto ? `R$ ${row.valor_azi_com_desconto}` : '—'}
            </TableCell>
            <TableCell className="text-right font-mono text-sm">
              {row.economia_rs ? `R$ ${row.economia_rs}` : '—'}
              {row.economia_pct && (
                <span className="text-muted-foreground ml-1">({row.economia_pct}%)</span>
              )}
            </TableCell>
            <TableCell>
              <div className="flex items-center gap-1">
                {row.pdf_generated && (
                  <HugeiconsIcon icon={PrinterIcon} strokeWidth={2} className="size-4 text-green-600" />
                )}
                <Button
                  variant="ghost"
                  size="icon"
                  className="size-7"
                  onClick={() => {
                    // TODO: abrir modal de detalhe do cálculo
                    // router.push(`/calculos/${row.calculation_id}`)
                  }}
                  disabled
                >
                  <HugeiconsIcon icon={FileViewIcon} strokeWidth={2} className="size-4" />
                </Button>
              </div>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}
