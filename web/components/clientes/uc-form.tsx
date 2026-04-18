'use client'

import { Controller, useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { CreateUcSchema, type CreateUcInput } from '@/types/clientes'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { NativeSelect, NativeSelectOption } from '@/components/ui/native-select'
import { Checkbox } from '@/components/ui/checkbox'

type UcCredentialOption = { id: string; label: string; documento_masked: string }

type UcFormProps = {
  defaultValues?: Partial<CreateUcInput>
  credentials: UcCredentialOption[]
  onSubmit: (data: CreateUcInput) => Promise<void>
  isLoading?: boolean
  mode: 'create' | 'edit'
}

export function UcForm({ defaultValues, credentials, onSubmit, isLoading, mode }: UcFormProps) {
  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
  } = useForm<CreateUcInput>({
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    resolver: zodResolver(CreateUcSchema as any),
    defaultValues: { ativa: true, ...defaultValues },
  })

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div className="space-y-1.5">
          <Label htmlFor="uc_code">Código UC *</Label>
          <Input
            id="uc_code"
            {...register('uc_code')}
            placeholder="Ex: 007098175908"
            disabled={mode === 'edit'}
          />
          {errors.uc_code && <p className="text-xs text-destructive">{errors.uc_code.message}</p>}
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="distribuidora">Distribuidora</Label>
          <Input id="distribuidora" {...register('distribuidora')} placeholder="Ex: Neoenergia" />
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="apelido">Apelido</Label>
          <Input id="apelido" {...register('apelido')} placeholder="Ex: Sede comercial" />
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="classe_consumo">Classe de Consumo</Label>
          <Input id="classe_consumo" {...register('classe_consumo')} placeholder="Ex: Comercial B3" />
        </div>

        <div className="space-y-1.5 sm:col-span-2">
          <Label htmlFor="endereco_unidade">Endereço da Unidade</Label>
          <Input
            id="endereco_unidade"
            {...register('endereco_unidade')}
            placeholder="Endereço completo da UC"
          />
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="cidade">Cidade</Label>
          <Input id="cidade" {...register('cidade')} placeholder="Cidade" />
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="uf">UF</Label>
          <Input id="uf" {...register('uf')} placeholder="BA" maxLength={2} className="uppercase" />
        </div>

        <div className="space-y-1.5 sm:col-span-2">
          <Label htmlFor="credential_id">Credencial de Integração</Label>
          <NativeSelect id="credential_id" {...register('credential_id')} className="w-full">
            <NativeSelectOption value="">Sem credencial</NativeSelectOption>
            {credentials.map((c) => (
              <NativeSelectOption key={c.id} value={c.id}>
                {c.label} — {c.documento_masked}
              </NativeSelectOption>
            ))}
          </NativeSelect>
          <p className="text-xs text-muted-foreground">
            Selecione uma credencial para habilitar a sincronização automática.
          </p>
        </div>
      </div>

      <div className="flex items-center gap-2 pt-2">
        <Controller
          name="ativa"
          control={control}
          render={({ field }) => (
            <div className="flex items-center gap-2">
              <Checkbox
                id="ativa"
                checked={field.value ?? false}
                onCheckedChange={field.onChange}
              />
              <Label htmlFor="ativa">UC ativa</Label>
            </div>
          )}
        />
      </div>

      <Button type="submit" disabled={isLoading}>
        {isLoading ? 'Salvando...' : mode === 'create' ? 'Criar UC' : 'Salvar'}
      </Button>
    </form>
  )
}
