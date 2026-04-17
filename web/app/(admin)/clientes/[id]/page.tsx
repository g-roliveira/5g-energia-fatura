'use client'

import { useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import { useQuery, useQueryClient } from '@tanstack/react-query'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Sheet, SheetContent, SheetHeader, SheetTitle } from '@/components/ui/sheet'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { StatusBadge } from '@/components/clientes/status-badge'
import { UcList } from '@/components/clientes/uc-list'
import { UcForm } from '@/components/clientes/uc-form'
import { CredentialForm } from '@/components/clientes/credential-form'
import type { CreateUcInput, CreateCredentialInput } from '@/types/clientes'

// ─── Local types (avoids importing Prisma in a client component) ──────────────

type Address = {
  cep?: string | null
  logradouro?: string | null
  numero?: string | null
  complemento?: string | null
  bairro?: string | null
  cidade?: string | null
  uf?: string | null
}

type CommercialData = {
  tipo_contrato?: string | null
  data_inicio?: string | null
  data_fim?: string | null
  status_contrato?: string | null
  observacoes_comerciais?: string | null
}

type Credential = {
  id: string
  label: string
  documento_masked: string
  uf: string
  tipo_acesso: string
}

type ConsumerUnit = {
  id: string
  uc_code: string
  apelido?: string | null
  distribuidora?: string | null
  ativa: boolean
  credential?: { id: string; label: string; documento_masked: string } | null
}

type ClientDetail = {
  id: string
  tipo_pessoa: string
  nome_razao: string
  nome_fantasia?: string | null
  cpf_cnpj: string
  email?: string | null
  telefone?: string | null
  status: 'ativo' | 'inativo' | 'prospecto'
  tipo_cliente: string
  observacoes?: string | null
  address?: Address | null
  ucs: ConsumerUnit[]
  commercial_data?: CommercialData | null
  credentials: Credential[]
}

// ─── Helper ───────────────────────────────────────────────────────────────────

function Field({ label, value }: { label: string; value?: string | null }) {
  return (
    <div>
      <dt className="text-xs font-medium text-muted-foreground uppercase tracking-wide">{label}</dt>
      <dd className="mt-0.5 text-sm">{value || <span className="text-muted-foreground">—</span>}</dd>
    </div>
  )
}

function formatDate(iso?: string | null) {
  if (!iso) return null
  return new Date(iso).toLocaleDateString('pt-BR')
}

// ─── Skeleton ─────────────────────────────────────────────────────────────────

function Skeleton() {
  return (
    <div className="p-6 space-y-6 animate-pulse">
      <div className="h-8 w-64 rounded bg-muted" />
      <div className="h-4 w-48 rounded bg-muted" />
      <div className="h-40 rounded bg-muted" />
    </div>
  )
}

// ─── Page ─────────────────────────────────────────────────────────────────────

export default function ClientDetailPage() {
  const params = useParams()
  const id = params.id as string
  const router = useRouter()
  const queryClient = useQueryClient()

  const [ucSheetOpen, setUcSheetOpen] = useState(false)
  const [credSheetOpen, setCredSheetOpen] = useState(false)
  const [ucLoading, setUcLoading] = useState(false)
  const [credLoading, setCredLoading] = useState(false)

  const { data: client, isLoading } = useQuery<ClientDetail>({
    queryKey: ['client', id],
    queryFn: async () => {
      const res = await fetch(`/api/clients/${id}`)
      if (!res.ok) throw new Error('Erro ao carregar cliente')
      return res.json()
    },
  })

  // ── Handlers ────────────────────────────────────────────────────────────────

  async function handleCreateUc(data: CreateUcInput) {
    setUcLoading(true)
    try {
      const res = await fetch(`/api/clients/${id}/ucs`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      })
      if (!res.ok) throw new Error('Erro ao criar UC')
      await queryClient.invalidateQueries({ queryKey: ['client', id] })
      setUcSheetOpen(false)
    } finally {
      setUcLoading(false)
    }
  }

  async function handleCreateCredential(data: CreateCredentialInput) {
    setCredLoading(true)
    try {
      const res = await fetch('/api/integration/credentials', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ...data, client_id: id }),
      })
      if (!res.ok) throw new Error('Erro ao criar credencial')
      await queryClient.invalidateQueries({ queryKey: ['client', id] })
      setCredSheetOpen(false)
    } finally {
      setCredLoading(false)
    }
  }

  async function handleArchive() {
    await fetch(`/api/clients/${id}/archive`, { method: 'POST' })
    router.push('/clientes')
  }

  // ── Loading ──────────────────────────────────────────────────────────────────

  if (isLoading || !client) return <Skeleton />

  // ── Render ──────────────────────────────────────────────────────────────────

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="space-y-3">
        <div className="flex flex-wrap items-center gap-3">
          <h1 className="text-2xl font-semibold">{client.nome_razao}</h1>
          <StatusBadge status={client.status} />
          <span className="text-sm text-muted-foreground capitalize">{client.tipo_cliente}</span>
          <span className="font-mono text-sm text-muted-foreground">{client.cpf_cnpj}</span>
        </div>

        <div className="flex gap-2">
          <Button variant="outline" size="sm" asChild>
            <Link href={`/clientes/${id}/editar`}>Editar</Link>
          </Button>

          <AlertDialog>
            <AlertDialogTrigger asChild>
              <Button variant="destructive" size="sm">Arquivar</Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Arquivar cliente</AlertDialogTitle>
                <AlertDialogDescription>
                  Tem certeza que deseja arquivar este cliente? Ele não será excluído, mas ficará inativo.
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel>Cancelar</AlertDialogCancel>
                <AlertDialogAction onClick={handleArchive}>Arquivar</AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </div>
      </div>

      {/* Tabs */}
      <Tabs defaultValue="dados">
        <TabsList>
          <TabsTrigger value="dados">Dados</TabsTrigger>
          <TabsTrigger value="endereco">Endereço</TabsTrigger>
          <TabsTrigger value="ucs">UCs</TabsTrigger>
          <TabsTrigger value="comercial">Comercial</TabsTrigger>
          <TabsTrigger value="integracao">Integração</TabsTrigger>
        </TabsList>

        {/* Tab 1: Dados */}
        <TabsContent value="dados">
          <Card>
            <CardHeader>
              <CardTitle>Dados do cliente</CardTitle>
            </CardHeader>
            <CardContent>
              <dl className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <Field label="Tipo de pessoa" value={client.tipo_pessoa} />
                <Field label="Nome / Razão social" value={client.nome_razao} />
                <Field label="Nome fantasia" value={client.nome_fantasia} />
                <Field label="CPF / CNPJ" value={client.cpf_cnpj} />
                <Field label="E-mail" value={client.email} />
                <Field label="Telefone" value={client.telefone} />
                <Field label="Status" value={client.status} />
                <Field label="Tipo de cliente" value={client.tipo_cliente} />
              </dl>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Tab 2: Endereço */}
        <TabsContent value="endereco">
          <Card>
            <CardHeader>
              <CardTitle>Endereço</CardTitle>
            </CardHeader>
            <CardContent>
              {client.address ? (
                <dl className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                  <Field label="CEP" value={client.address.cep} />
                  <Field
                    label="Logradouro"
                    value={
                      [client.address.logradouro, client.address.numero]
                        .filter(Boolean)
                        .join(', ') || null
                    }
                  />
                  <Field label="Complemento" value={client.address.complemento} />
                  <Field label="Bairro" value={client.address.bairro} />
                  <Field label="Cidade" value={client.address.cidade} />
                  <Field label="UF" value={client.address.uf} />
                </dl>
              ) : (
                <p className="text-sm text-muted-foreground">Endereço não cadastrado.</p>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Tab 3: UCs */}
        <TabsContent value="ucs">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>Unidades consumidoras</CardTitle>
                <Button size="sm" onClick={() => setUcSheetOpen(true)}>
                  Adicionar UC
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <UcList
                clientId={id}
                ucs={client.ucs}
                onSyncSuccess={() => queryClient.invalidateQueries({ queryKey: ['client', id] })}
              />
            </CardContent>
          </Card>

          <Sheet open={ucSheetOpen} onOpenChange={setUcSheetOpen}>
            <SheetContent className="overflow-y-auto">
              <SheetHeader>
                <SheetTitle>Nova UC</SheetTitle>
              </SheetHeader>
              <div className="mt-4">
                <UcForm
                  mode="create"
                  credentials={client.credentials}
                  onSubmit={handleCreateUc}
                  isLoading={ucLoading}
                />
              </div>
            </SheetContent>
          </Sheet>
        </TabsContent>

        {/* Tab 4: Comercial */}
        <TabsContent value="comercial">
          <Card>
            <CardHeader>
              <CardTitle>Dados comerciais</CardTitle>
            </CardHeader>
            <CardContent>
              {client.commercial_data ? (
                <dl className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                  <Field label="Tipo de contrato" value={client.commercial_data.tipo_contrato} />
                  <Field label="Data de início" value={formatDate(client.commercial_data.data_inicio)} />
                  <Field label="Data de fim" value={formatDate(client.commercial_data.data_fim)} />
                  <Field label="Status do contrato" value={client.commercial_data.status_contrato} />
                  <div className="sm:col-span-2">
                    <Field label="Observações comerciais" value={client.commercial_data.observacoes_comerciais} />
                  </div>
                </dl>
              ) : (
                <p className="text-sm text-muted-foreground">Dados comerciais não cadastrados.</p>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Tab 5: Integração */}
        <TabsContent value="integracao">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>Credenciais de integração</CardTitle>
                <Button size="sm" onClick={() => setCredSheetOpen(true)}>
                  Adicionar credencial
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              {client.credentials.length === 0 ? (
                <p className="text-sm text-muted-foreground">Nenhuma credencial cadastrada.</p>
              ) : (
                <div className="space-y-3">
                  {client.credentials.map((cred) => (
                    <div key={cred.id} className="flex items-center justify-between rounded-lg border p-4">
                      <div className="space-y-0.5">
                        <p className="text-sm font-medium">{cred.label}</p>
                        <p className="font-mono text-xs text-muted-foreground">{cred.documento_masked}</p>
                        <p className="text-xs text-muted-foreground">
                          {cred.uf} · {cred.tipo_acesso}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>

          <Sheet open={credSheetOpen} onOpenChange={setCredSheetOpen}>
            <SheetContent className="overflow-y-auto">
              <SheetHeader>
                <SheetTitle>Nova credencial</SheetTitle>
              </SheetHeader>
              <div className="mt-4">
                <CredentialForm
                  clientId={id}
                  onSubmit={handleCreateCredential}
                  isLoading={credLoading}
                />
              </div>
            </SheetContent>
          </Sheet>
        </TabsContent>
      </Tabs>
    </div>
  )
}
