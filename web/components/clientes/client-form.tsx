'use client'

import { Controller, useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { CreateClientSchema, type CreateClientInput } from '@/types/clientes'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { NativeSelect, NativeSelectOption } from '@/components/ui/native-select'

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
    control,
    formState: { errors },
  } = useForm<CreateClientInput>({
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    resolver: zodResolver(CreateClientSchema as any),
    defaultValues: { status: 'prospecto', tipo_pessoa: 'PF', ...defaultValues },
  })

  const tipoPessoa = watch('tipo_pessoa')

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-6">
      {/* Seção 1: Dados Principais */}
      <Card>
        <CardHeader>
          <CardTitle>Dados Principais</CardTitle>
          <CardDescription>Identificação e tipo do cliente</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <Controller
            name="tipo_pessoa"
            control={control}
            render={({ field }) => (
              <RadioGroup
                value={field.value}
                onValueChange={field.onChange}
                className="flex gap-4"
              >
                <div className="flex items-center gap-2">
                  <RadioGroupItem value="PF" id="tipo-pf" />
                  <Label htmlFor="tipo-pf">Pessoa Física</Label>
                </div>
                <div className="flex items-center gap-2">
                  <RadioGroupItem value="PJ" id="tipo-pj" />
                  <Label htmlFor="tipo-pj">Pessoa Jurídica</Label>
                </div>
              </RadioGroup>
            )}
          />

          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div className="space-y-1.5">
              <Label htmlFor="nome_razao" className="text-sm font-medium">Nome / Razão Social *</Label>
              <Input id="nome_razao" {...register('nome_razao')} placeholder="Nome completo ou razão social" />
              {errors.nome_razao && <p className="text-xs text-destructive">{errors.nome_razao.message}</p>}
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="nome_fantasia" className="text-sm font-medium">Nome Fantasia</Label>
              <Input id="nome_fantasia" {...register('nome_fantasia')} placeholder="Nome fantasia (opcional)" />
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="cpf_cnpj" className="text-sm font-medium">{tipoPessoa === 'PF' ? 'CPF' : 'CNPJ'} *</Label>
              <Input
                id="cpf_cnpj"
                {...register('cpf_cnpj')}
                placeholder={tipoPessoa === 'PF' ? '000.000.000-00' : '00.000.000/0000-00'}
              />
              {errors.cpf_cnpj && <p className="text-xs text-destructive">{errors.cpf_cnpj.message}</p>}
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="tipo_cliente" className="text-sm font-medium">Tipo de Cliente *</Label>
              <NativeSelect id="tipo_cliente" {...register('tipo_cliente')} className="w-full">
                <NativeSelectOption value="">Selecionar...</NativeSelectOption>
                <NativeSelectOption value="residencial">Residencial</NativeSelectOption>
                <NativeSelectOption value="condominio">Condomínio</NativeSelectOption>
                <NativeSelectOption value="empresa">Empresa</NativeSelectOption>
                <NativeSelectOption value="imobiliaria">Imobiliária</NativeSelectOption>
                <NativeSelectOption value="outro">Outro</NativeSelectOption>
              </NativeSelect>
              {errors.tipo_cliente && <p className="text-xs text-destructive">{errors.tipo_cliente.message}</p>}
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="status" className="text-sm font-medium">Status</Label>
              <NativeSelect id="status" {...register('status')} className="w-full">
                <NativeSelectOption value="prospecto">Prospecto</NativeSelectOption>
                <NativeSelectOption value="ativo">Ativo</NativeSelectOption>
                <NativeSelectOption value="inativo">Inativo</NativeSelectOption>
              </NativeSelect>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Seção 2: Contato */}
      <Card>
        <CardHeader>
          <CardTitle>Contato</CardTitle>
          <CardDescription>Email e telefone de contato</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div className="space-y-1.5">
              <Label htmlFor="email" className="text-sm font-medium">Email</Label>
              <Input id="email" type="email" {...register('email')} placeholder="email@exemplo.com" />
              {errors.email && <p className="text-xs text-destructive">{errors.email.message}</p>}
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="telefone" className="text-sm font-medium">Telefone</Label>
              <Input id="telefone" {...register('telefone')} placeholder="(00) 00000-0000" />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Seção 3: Endereço */}
      <Card>
        <CardHeader>
          <CardTitle>Endereço</CardTitle>
          <CardDescription>Localização do cliente</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            <div className="space-y-1.5">
              <Label htmlFor="cep" className="text-sm font-medium">CEP</Label>
              <Input id="cep" {...register('cep')} placeholder="00000-000" />
            </div>
            <div className="space-y-1.5 sm:col-span-2">
              <Label htmlFor="logradouro" className="text-sm font-medium">Logradouro</Label>
              <Input id="logradouro" {...register('logradouro')} placeholder="Rua, Avenida..." />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="numero" className="text-sm font-medium">Número</Label>
              <Input id="numero" {...register('numero')} placeholder="123" />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="complemento" className="text-sm font-medium">Complemento</Label>
              <Input id="complemento" {...register('complemento')} placeholder="Apto, Sala..." />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="bairro" className="text-sm font-medium">Bairro</Label>
              <Input id="bairro" {...register('bairro')} placeholder="Bairro" />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="cidade" className="text-sm font-medium">Cidade</Label>
              <Input id="cidade" {...register('cidade')} placeholder="Cidade" />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="uf" className="text-sm font-medium">UF</Label>
              <Input id="uf" {...register('uf')} placeholder="BA" maxLength={2} className="uppercase" />
              {errors.uf && <p className="text-xs text-destructive">{errors.uf.message}</p>}
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Seção 4: Dados Comerciais */}
      <Card>
        <CardHeader>
          <CardTitle>Dados Comerciais</CardTitle>
          <CardDescription>Informações contratuais</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div className="space-y-1.5">
              <Label htmlFor="tipo_contrato" className="text-sm font-medium">Tipo de Contrato</Label>
              <Input id="tipo_contrato" {...register('tipo_contrato')} placeholder="Ex: Mensal, Anual..." />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="status_contrato" className="text-sm font-medium">Status do Contrato</Label>
              <Input id="status_contrato" {...register('status_contrato')} placeholder="Ex: Ativo, Encerrado..." />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="data_inicio" className="text-sm font-medium">Data Início</Label>
              <Input id="data_inicio" type="date" {...register('data_inicio')} />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="data_fim" className="text-sm font-medium">Data Fim</Label>
              <Input id="data_fim" type="date" {...register('data_fim')} />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Seção 5: Observações */}
      <Card>
        <CardHeader>
          <CardTitle>Observações</CardTitle>
          <CardDescription>Anotações internas</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-1.5">
            <Label htmlFor="observacoes" className="text-sm font-medium">Observações Gerais</Label>
            <Textarea
              id="observacoes"
              {...register('observacoes')}
              placeholder="Observações gerais sobre o cliente..."
              className="min-h-[100px]"
            />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="observacoes_comerciais" className="text-sm font-medium">Observações Comerciais</Label>
            <Textarea
              id="observacoes_comerciais"
              {...register('observacoes_comerciais')}
              placeholder="Observações comerciais e contratuais..."
              className="min-h-[80px]"
            />
          </div>
        </CardContent>
      </Card>

      <div className="flex gap-3 pt-2">
        <Button type="submit" disabled={isLoading}>
          {isLoading ? 'Salvando...' : mode === 'create' ? 'Criar cliente' : 'Salvar alterações'}
        </Button>
      </div>
    </form>
  )
}
