'use client'

import * as React from 'react'
import { useRouter, useSearchParams, usePathname } from 'next/navigation'
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
import { Input } from '@/components/ui/input'
import { NativeSelect } from '@/components/ui/native-select'
import { Skeleton } from '@/components/ui/skeleton'
import { Badge } from '@/components/ui/badge'
import { StatusBadge } from './status-badge'
import type { ClientStatus, TipoCliente } from '@prisma/client'

type ClientRow = {
  id: string
  nome_razao: string
  cpf_cnpj: string
  tipo_cliente: TipoCliente
  status: ClientStatus
  cidade?: string | null
  uf?: string | null
  created_at: Date
  address?: { cidade?: string | null; uf?: string | null } | null
  _count: { ucs: number }
}

type ClientTableProps = {
  data: ClientRow[]
  total: number
  page: number
  pageSize: number
  isLoading?: boolean
}

export function ClientTable({ data, total, page, pageSize, isLoading }: ClientTableProps) {
  const router = useRouter()
  const pathname = usePathname()
  const searchParams = useSearchParams()
  const [search, setSearch] = React.useState(searchParams.get('search') ?? '')
  const debounceRef = React.useRef<ReturnType<typeof setTimeout> | null>(null)

  function updateParam(key: string, value: string) {
    const params = new URLSearchParams(searchParams.toString())
    if (value) params.set(key, value)
    else params.delete(key)
    params.set('page', '1')
    router.push(`${pathname}?${params}`)
  }

  function handleSearch(value: string) {
    setSearch(value)
    if (debounceRef.current) clearTimeout(debounceRef.current)
    debounceRef.current = setTimeout(() => updateParam('search', value), 300)
  }

  const columns: ColumnDef<ClientRow>[] = [
    {
      accessorKey: 'nome_razao',
      header: 'Nome / Razão Social',
      cell: ({ row }) => <span className="font-medium">{row.original.nome_razao}</span>,
    },
    { accessorKey: 'cpf_cnpj', header: 'CPF / CNPJ' },
    {
      accessorKey: 'tipo_cliente',
      header: 'Tipo',
      cell: ({ row }) => <Badge variant="secondary" className="capitalize">{row.original.tipo_cliente}</Badge>,
    },
    {
      accessorKey: 'status',
      header: 'Status',
      cell: ({ row }) => <StatusBadge status={row.original.status} />,
    },
    {
      id: 'ucs',
      header: 'UCs',
      cell: ({ row }) => row.original._count.ucs,
    },
    {
      id: 'localizacao',
      header: 'Cidade / UF',
      cell: ({ row }) => {
        const addr = row.original.address
        if (addr?.cidade && addr?.uf) return `${addr.cidade} / ${addr.uf}`
        return '—'
      },
    },
    {
      accessorKey: 'created_at',
      header: 'Cadastrado em',
      cell: ({ row }) =>
        format(new Date(row.original.created_at), 'dd/MM/yyyy', { locale: ptBR }),
    },
    {
      id: 'acoes',
      header: '',
      cell: ({ row }) => (
        <Button variant="ghost" size="sm" onClick={() => router.push(`/clientes/${row.original.id}`)}>
          Visualizar
        </Button>
      ),
    },
  ]

  const table = useReactTable({ data, columns, getCoreRowModel: getCoreRowModel() })

  const totalPages = Math.ceil(total / pageSize)

  return (
    <div className="space-y-4">
      {/* Filters */}
      <div className="flex flex-wrap gap-3">
        <Input
          placeholder="Buscar por nome ou CPF/CNPJ..."
          value={search}
          onChange={(e) => handleSearch(e.target.value)}
          className="max-w-sm"
        />
        <NativeSelect
          value={searchParams.get('status') ?? ''}
          onChange={(e) => updateParam('status', e.target.value)}
          className="w-40"
        >
          <option value="">Todos os status</option>
          <option value="ativo">Ativo</option>
          <option value="inativo">Inativo</option>
          <option value="prospecto">Prospecto</option>
        </NativeSelect>
        <NativeSelect
          value={searchParams.get('tipo_cliente') ?? ''}
          onChange={(e) => updateParam('tipo_cliente', e.target.value)}
          className="w-44"
        >
          <option value="">Todos os tipos</option>
          <option value="residencial">Residencial</option>
          <option value="condominio">Condomínio</option>
          <option value="empresa">Empresa</option>
          <option value="imobiliaria">Imobiliária</option>
          <option value="outro">Outro</option>
        </NativeSelect>
      </div>

      {/* Table */}
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
                  Nenhum cliente encontrado.
                </TableCell>
              </TableRow>
            ) : (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  className="cursor-pointer hover:bg-muted/50"
                  onClick={() => router.push(`/clientes/${row.original.id}`)}
                >
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

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-between text-sm text-muted-foreground">
          <span>{total} clientes no total</span>
          <div className="flex gap-2">
            <Button
              variant="outline" size="sm"
              disabled={page <= 1}
              onClick={() => updateParam('page', String(page - 1))}
            >
              Anterior
            </Button>
            <span className="flex items-center px-2">
              {page} / {totalPages}
            </span>
            <Button
              variant="outline" size="sm"
              disabled={page >= totalPages}
              onClick={() => updateParam('page', String(page + 1))}
            >
              Próximo
            </Button>
          </div>
        </div>
      )}
    </div>
  )
}
