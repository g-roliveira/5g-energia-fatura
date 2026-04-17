'use client'

import { useRouter } from 'next/navigation'
import {
  useReactTable,
  getCoreRowModel,
  flexRender,
  type ColumnDef,
} from '@tanstack/react-table'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { CompletudeBadge } from './completude-badge'
import type { GoInvoice } from '@/types/clientes'

type InvoiceTableProps = {
  data: GoInvoice[]
  isLoading?: boolean
  clientId: string
  ucId: string
}

export function InvoiceTable({ data, isLoading, clientId, ucId }: InvoiceTableProps) {
  const router = useRouter()

  const columns: ColumnDef<GoInvoice>[] = [
    {
      id: 'numero',
      header: 'Número',
      cell: ({ row }) => (
        <span className="font-mono text-sm">{row.original.billing_record?.numero_fatura ?? '—'}</span>
      ),
    },
    {
      id: 'referencia',
      header: 'Referência',
      cell: ({ row }) => row.original.billing_record?.mes_referencia ?? '—',
    },
    {
      id: 'valor',
      header: 'Valor',
      cell: ({ row }) => {
        const v = row.original.billing_record?.valor_total
        return v ? `R$ ${v}` : '—'
      },
    },
    {
      id: 'vencimento',
      header: 'Vencimento',
      cell: ({ row }) => {
        const d = row.original.billing_record?.data_vencimento
        if (!d) return '—'
        try {
          return format(new Date(d), 'dd/MM/yyyy', { locale: ptBR })
        } catch {
          return d
        }
      },
    },
    {
      id: 'completude',
      header: 'Completude',
      cell: ({ row }) => {
        const s = row.original.billing_record?.completeness?.status ?? row.original.completeness_status
        if (!s) return '—'
        return <CompletudeBadge status={s as 'complete' | 'partial' | 'failed'} />
      },
    },
    {
      id: 'atualizado',
      header: 'Atualizado',
      cell: ({ row }) => {
        try {
          return format(new Date(row.original.updated_at), 'dd/MM/yyyy HH:mm', { locale: ptBR })
        } catch {
          return row.original.updated_at
        }
      },
    },
    {
      id: 'acoes',
      header: '',
      cell: ({ row }) => (
        <Button
          variant="ghost"
          size="sm"
          onClick={() =>
            router.push(`/clientes/${clientId}/ucs/${ucId}/faturas/${row.original.id}`)
          }
        >
          Detalhe
        </Button>
      ),
    },
  ]

  const table = useReactTable({ data, columns, getCoreRowModel: getCoreRowModel() })

  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((hg) => (
            <TableRow key={hg.id}>
              {hg.headers.map((h) => (
                <TableHead key={h.id}>
                  {flexRender(h.column.columnDef.header, h.getContext())}
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {isLoading ? (
            Array.from({ length: 5 }).map((_, i) => (
              <TableRow key={i}>
                {columns.map((_, j) => (
                  <TableCell key={j}><Skeleton className="h-4 w-full" /></TableCell>
                ))}
              </TableRow>
            ))
          ) : table.getRowModel().rows.length === 0 ? (
            <TableRow>
              <TableCell colSpan={columns.length} className="h-32 text-center text-muted-foreground">
                Nenhuma fatura encontrada.
              </TableCell>
            </TableRow>
          ) : (
            table.getRowModel().rows.map((row) => (
              <TableRow key={row.id} className="cursor-pointer hover:bg-muted/50">
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  )
}
