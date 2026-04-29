'use client'

import { Controller, useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { CreateClientMinimalSchema, type CreateClientMinimalInput } from '@/types/clientes'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { NativeSelect, NativeSelectOption } from '@/components/ui/native-select'

type ClientFormMinimalProps = {
  onSubmit: (data: CreateClientMinimalInput) => Promise<void>
  isLoading?: boolean
}

export function ClientFormMinimal({ onSubmit, isLoading }: ClientFormMinimalProps) {
  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
  } = useForm<CreateClientMinimalInput>({
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    resolver: zodResolver(CreateClientMinimalSchema as any),
    defaultValues: { tipo_pessoa: 'PF', tipo_cliente: 'outro', status: 'prospecto' },
  })

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <p className="text-sm text-muted-foreground">
        Endereço, contato e UCs serão importados automaticamente ao cadastrar as credenciais do portal.
      </p>

      <div className="space-y-2">
        <Label>Tipo de pessoa *</Label>
        <Controller
          name="tipo_pessoa"
          control={control}
          render={({ field }) => (
            <RadioGroup value={field.value} onValueChange={field.onChange} className="flex gap-4">
              <div className="flex items-center gap-2">
                <RadioGroupItem value="PF" id="tipo_pf" />
                <Label htmlFor="tipo_pf">Pessoa Física</Label>
              </div>
              <div className="flex items-center gap-2">
                <RadioGroupItem value="PJ" id="tipo_pj" />
                <Label htmlFor="tipo_pj">Pessoa Jurídica</Label>
              </div>
            </RadioGroup>
          )}
        />
        {errors.tipo_pessoa && <p className="text-xs text-destructive">{errors.tipo_pessoa.message}</p>}
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="nome_razao">Nome / Razão Social *</Label>
        <Input id="nome_razao" {...register('nome_razao')} placeholder="Nome completo ou razão social" />
        {errors.nome_razao && <p className="text-xs text-destructive">{errors.nome_razao.message}</p>}
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="cpf_cnpj">CPF / CNPJ *</Label>
        <Input id="cpf_cnpj" {...register('cpf_cnpj')} placeholder="Somente números" />
        {errors.cpf_cnpj && <p className="text-xs text-destructive">{errors.cpf_cnpj.message}</p>}
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="tipo_cliente">Tipo de cliente</Label>
        <NativeSelect id="tipo_cliente" {...register('tipo_cliente')} className="w-full">
          <NativeSelectOption value="residencial">Residencial</NativeSelectOption>
          <NativeSelectOption value="condominio">Condomínio</NativeSelectOption>
          <NativeSelectOption value="empresa">Empresa</NativeSelectOption>
          <NativeSelectOption value="imobiliaria">Imobiliária</NativeSelectOption>
          <NativeSelectOption value="outro">Outro</NativeSelectOption>
        </NativeSelect>
      </div>

      <Button type="submit" className="w-full" disabled={isLoading}>
        {isLoading ? 'Criando...' : 'Criar cliente'}
      </Button>
    </form>
  )
}
