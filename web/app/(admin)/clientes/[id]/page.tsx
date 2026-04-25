'use client'

import { useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { ArrowLeft01Icon, PencilEdit01Icon } from '@hugeicons/core-free-icons'
import { HugeiconsIcon } from '@hugeicons/react'

import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
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
import { DiscoveryStep } from '@/components/clientes/discovery-step'
import type { CreateUcInput, CreateCredentialInput, GoDiscoveryResult } from '@/types/clientes'
import { useSetBreadcrumbTitle } from '@/contexts/breadcrumb'

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

function PageSkeleton() {
  return (
    <div className="flex flex-col gap-6 p-6">
      <Skeleton className="h-8 w-32" />
      <div className="space-y-2">
        <Skeleton className="h-8 w-64" />
        <Skeleton className="h-4 w-48" />
      </div>
      <Skeleton className="h-10 w-48" />
      <Skeleton className="h-64 w-full" />
    </div>
  )
}

// ─── Page ─────────────────────────────────────────────────────────────────────

export default function ClientDetailPage() {
  const { id } = useParams<{ id: string }>()
  const router = useRouter()
  const queryClient = useQueryClient()

  const [ucSheetOpen, setUcSheetOpen] = useState(false)
  const [credSheetOpen, setCredSheetOpen] = useState(false)
  const [ucLoading, setUcLoading] = useState(false)
  const [credLoading, setCredLoading] = useState(false)

  // ── Credential discovery dialog state ────────────────────────────────────────
  const [credDialogStep, setCredDialogStep] = useState<'form' | 'discovery'>('form')
  const [createdCredentialId, setCreatedCredentialId] = useState<string | null>(null)
  const [discoveryData, setDiscoveryData] = useState<GoDiscoveryResult | null>(null)
  const [discoveryLoading, setDiscoveryLoading] = useState(false)
  const [discoveryError, setDiscoveryError] = useState<string | null>(null)

  const { data: client, isLoading } = useQuery<ClientDetail>({
    queryKey: ['client', id],
    queryFn: async () => {
      const res = await fetch(`/api/clients/${id}`)
      if (!res.ok) throw new Error('Erro ao carregar cliente')
      return res.json()
    },
  })

  useSetBreadcrumbTitle(id, client?.nome_razao)

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
      const credential = await res.json()
      await queryClient.invalidateQueries({ queryKey: ['client', id] })

      // Transition to discovery step
      setCreatedCredentialId(credential.id)
      setCredDialogStep('discovery')
      setDiscoveryLoading(true)
      setDiscoveryError(null)

      try {
        const discRes = await fetch(`/api/integration/credentials/${credential.id}/discover`)
        if (!discRes.ok) throw new Error('Erro ao consultar dados da concessionária')
        setDiscoveryData(await discRes.json())
      } catch (e) {
        setDiscoveryError((e as Error).message)
      } finally {
        setDiscoveryLoading(false)
      }
    } finally {
      setCredLoading(false)
    }
  }

  function handleCredDialogClose(open: boolean) {
    if (!open) {
      setCredSheetOpen(false)
      setCredDialogStep('form')
      setCreatedCredentialId(null)
      setDiscoveryData(null)
      setDiscoveryError(null)
    }
  }

  async function handleArchive() {
    await fetch(`/api/clients/${id}/archive`, { method: 'POST' })
    router.push('/clientes')
  }

  // ── Loading ──────────────────────────────────────────────────────────────────

  if (isLoading || !client) return <PageSkeleton />

  // ── Render ──────────────────────────────────────────────────────────────────

  return (
    <div className="flex flex-col gap-6 p-6">
      {/* Back button */}
      <div>
        <Button variant="outline" size="sm" asChild>
          <Link href="/clientes">
            <HugeiconsIcon icon={ArrowLeft01Icon} strokeWidth={2} />
            Clientes
          </Link>
        </Button>
      </div>

      {/* Header */}
      <div className="space-y-3">
        <div className="flex flex-wrap items-center gap-3">
          <h1 className="text-2xl font-semibold">{client.nome_razao}</h1>
          <StatusBadge status={client.status} />
          <Badge variant="secondary" className="capitalize">{client.tipo_cliente}</Badge>
          <span className="font-mono text-sm text-muted-foreground">{client.cpf_cnpj}</span>
        </div>

        {client.email && (
          <p className="text-sm text-muted-foreground">{client.email}</p>
        )}

        <div className="flex gap-2">
          <Button variant="outline" size="sm" asChild>
            <Link href={`/clientes/${id}/editar`}>
              <HugeiconsIcon icon={PencilEdit01Icon} strokeWidth={2} />
              Editar
            </Link>
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
              <CardDescription>Informações cadastrais e de contato</CardDescription>
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
              <CardDescription>Localização do cliente</CardDescription>
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
                <div>
                  <CardTitle>Unidades consumidoras</CardTitle>
                  <CardDescription>UCs vinculadas a este cliente</CardDescription>
                </div>
                <div className="flex items-center gap-2">
                  <Button variant="outline" size="sm" asChild>
                    <Link href={`/clientes/${id}/ucs`}>
                      Ver todas
                    </Link>
                  </Button>
                  <Button size="sm" onClick={() => setUcSheetOpen(true)}>
                    Adicionar UC
                  </Button>
                </div>
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

          <Dialog open={ucSheetOpen} onOpenChange={setUcSheetOpen}>
            <DialogContent className="max-w-2xl">
              <DialogHeader>
                <DialogTitle>Adicionar UC</DialogTitle>
              </DialogHeader>
              <UcForm
                mode="create"
                credentials={client.credentials}
                onSubmit={handleCreateUc}
                isLoading={ucLoading}
              />
            </DialogContent>
          </Dialog>
        </TabsContent>

        {/* Tab 4: Comercial */}
        <TabsContent value="comercial">
          <Card>
            <CardHeader>
              <CardTitle>Dados comerciais</CardTitle>
              <CardDescription>Contrato e histórico comercial</CardDescription>
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
                <div>
                  <CardTitle>Credenciais de integração</CardTitle>
                  <CardDescription>Acesso ao portal da concessionária</CardDescription>
                </div>
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

          <Dialog open={credSheetOpen} onOpenChange={handleCredDialogClose}>
            <DialogContent className="max-w-3xl">
              <DialogHeader>
                <DialogTitle>
                  {credDialogStep === 'form' ? 'Adicionar Credencial' : 'Dados Descobertos'}
                </DialogTitle>
              </DialogHeader>

              {credDialogStep === 'form' && (
                <CredentialForm
                  clientId={id}
                  onSubmit={handleCreateCredential}
                  isLoading={credLoading}
                />
              )}

              {credDialogStep === 'discovery' && (
                <>
                  {discoveryLoading && (
                    <div className="space-y-3 py-10 text-center">
                      <div className="space-y-2">
                        <Skeleton className="mx-auto h-3 w-48" />
                        <Skeleton className="mx-auto h-3 w-36" />
                      </div>
                      <p className="text-sm text-muted-foreground">
                        Consultando a Neoenergia, aguarde...
                      </p>
                    </div>
                  )}
                  {discoveryError && !discoveryLoading && (
                    <div className="space-y-4 py-4">
                      <p className="text-sm text-destructive">{discoveryError}</p>
                      <Button variant="outline" onClick={() => handleCredDialogClose(false)}>
                        Fechar
                      </Button>
                    </div>
                  )}
                  {discoveryData && !discoveryLoading && (
                    <DiscoveryStep
                      clientId={id}
                      credentialId={createdCredentialId!}
                      discovery={discoveryData}
                      currentClient={{
                        nome_razao: client.nome_razao,
                        email: client.email,
                        telefone: client.telefone,
                      }}
                      onImportComplete={() => queryClient.invalidateQueries({ queryKey: ['client', id] })}
                      onClose={() => handleCredDialogClose(false)}
                    />
                  )}
                </>
              )}
            </DialogContent>
          </Dialog>
        </TabsContent>
      </Tabs>
    </div>
  )
}
