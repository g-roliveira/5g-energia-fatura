'use client'

import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { CreateClientSchema, type CreateClientInput } from '@/types/clientes'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Separator } from '@/components/ui/separator'

type ClientFormProps = {
  defaultValues?: Partial<CreateClientInput>
  onSubmit: (data: CreateClientInput) => Promise<void>
  isLoading?: boolean
  mode: 'create' | 'edit'
}

export function ClientForm({ defaultValues, onSubmit, isLoading, mode }: ClientFormProps) {
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<CreateClientInput>({
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    resolver: zodResolver(CreateClientSchema as any),
    defaultValues: { status: 'prospecto', tipo_pessoa: 'PF', ...defaultValues },
  })

  const tipoPessoa = watch('tipo_pessoa')

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-8">
      {/* Seção 1: Dados Principais */}
      <section className="space-y-4">
        <h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wide">
          Dados Principais
        </h3>

        <div className="flex gap-4">
          <label className="flex items-center gap-2 cursor-pointer">
            <input type="radio" value="PF" {...register('tipo_pessoa')} className="accent-primary" />
            <span className="text-sm">Pessoa Física</span>
          </label>
          <label className="flex items-center gap-2 cursor-pointer">
            <input type="radio" value="PJ" {...register('tipo_pessoa')} className="accent-primary" />
            <span className="text-sm">Pessoa Jurídica</span>
          </label>
        </div>

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="space-y-1.5">
            <label className="text-sm font-medium">Nome / Razão Social *</label>
            <Input {...register('nome_razao')} placeholder="Nome completo ou razão social" />
            {errors.nome_razao && <p className="text-xs text-destructive">{errors.nome_razao.message}</p>}
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-medium">Nome Fantasia</label>
            <Input {...register('nome_fantasia')} placeholder="Nome fantasia (opcional)" />
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-medium">{tipoPessoa === 'PF' ? 'CPF' : 'CNPJ'} *</label>
            <Input
              {...register('cpf_cnpj')}
              placeholder={tipoPessoa === 'PF' ? '000.000.000-00' : '00.000.000/0000-00'}
            />
            {errors.cpf_cnpj && <p className="text-xs text-destructive">{errors.cpf_cnpj.message}</p>}
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-medium">Tipo de Cliente *</label>
            <select
              {...register('tipo_cliente')}
              className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
            >
              <option value="">Selecionar...</option>
              <option value="residencial">Residencial</option>
              <option value="condominio">Condomínio</option>
              <option value="empresa">Empresa</option>
              <option value="imobiliaria">Imobiliária</option>
              <option value="outro">Outro</option>
            </select>
            {errors.tipo_cliente && <p className="text-xs text-destructive">{errors.tipo_cliente.message}</p>}
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-medium">Status</label>
            <select
              {...register('status')}
              className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
            >
              <option value="prospecto">Prospecto</option>
              <option value="ativo">Ativo</option>
              <option value="inativo">Inativo</option>
            </select>
          </div>
        </div>
      </section>

      <Separator />

      {/* Seção 2: Contato */}
      <section className="space-y-4">
        <h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wide">Contato</h3>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="space-y-1.5">
            <label className="text-sm font-medium">Email</label>
            <Input type="email" {...register('email')} placeholder="email@exemplo.com" />
            {errors.email && <p className="text-xs text-destructive">{errors.email.message}</p>}
          </div>
          <div className="space-y-1.5">
            <label className="text-sm font-medium">Telefone</label>
            <Input {...register('telefone')} placeholder="(00) 00000-0000" />
          </div>
        </div>
      </section>

      <Separator />

      {/* Seção 3: Endereço */}
      <section className="space-y-4">
        <h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wide">Endereço</h3>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <div className="space-y-1.5">
            <label className="text-sm font-medium">CEP</label>
            <Input {...register('cep')} placeholder="00000-000" />
          </div>
          <div className="space-y-1.5 sm:col-span-2">
            <label className="text-sm font-medium">Logradouro</label>
            <Input {...register('logradouro')} placeholder="Rua, Avenida..." />
          </div>
          <div className="space-y-1.5">
            <label className="text-sm font-medium">Número</label>
            <Input {...register('numero')} placeholder="123" />
          </div>
          <div className="space-y-1.5">
            <label className="text-sm font-medium">Complemento</label>
            <Input {...register('complemento')} placeholder="Apto, Sala..." />
          </div>
          <div className="space-y-1.5">
            <label className="text-sm font-medium">Bairro</label>
            <Input {...register('bairro')} placeholder="Bairro" />
          </div>
          <div className="space-y-1.5">
            <label className="text-sm font-medium">Cidade</label>
            <Input {...register('cidade')} placeholder="Cidade" />
          </div>
          <div className="space-y-1.5">
            <label className="text-sm font-medium">UF</label>
            <Input {...register('uf')} placeholder="BA" maxLength={2} className="uppercase" />
            {errors.uf && <p className="text-xs text-destructive">{errors.uf.message}</p>}
          </div>
        </div>
      </section>

      <Separator />

      {/* Seção 4: Comercial */}
      <section className="space-y-4">
        <h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wide">
          Dados Comerciais
        </h3>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="space-y-1.5">
            <label className="text-sm font-medium">Tipo de Contrato</label>
            <Input {...register('tipo_contrato')} placeholder="Ex: Mensal, Anual..." />
          </div>
          <div className="space-y-1.5">
            <label className="text-sm font-medium">Status do Contrato</label>
            <Input {...register('status_contrato')} placeholder="Ex: Ativo, Encerrado..." />
          </div>
          <div className="space-y-1.5">
            <label className="text-sm font-medium">Data Início</label>
            <Input type="date" {...register('data_inicio')} />
          </div>
          <div className="space-y-1.5">
            <label className="text-sm font-medium">Data Fim</label>
            <Input type="date" {...register('data_fim')} />
          </div>
        </div>
      </section>

      <Separator />

      {/* Seção 5: Observações */}
      <section className="space-y-4">
        <h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wide">Observações</h3>
        <Textarea
          {...register('observacoes')}
          placeholder="Observações gerais sobre o cliente..."
          className="min-h-[100px]"
        />
        <Textarea
          {...register('observacoes_comerciais')}
          placeholder="Observações comerciais e contratuais..."
          className="min-h-[80px]"
        />
      </section>

      <div className="flex gap-3 pt-2">
        <Button type="submit" disabled={isLoading}>
          {isLoading ? 'Salvando...' : mode === 'create' ? 'Criar cliente' : 'Salvar alterações'}
        </Button>
      </div>
    </form>
  )
}
