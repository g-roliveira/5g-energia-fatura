'use client'

import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { CreateUcSchema, type CreateUcInput } from '@/types/clientes'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

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
          <label className="text-sm font-medium">Código UC *</label>
          <Input
            {...register('uc_code')}
            placeholder="Ex: 007098175908"
            disabled={mode === 'edit'}
          />
          {errors.uc_code && <p className="text-xs text-destructive">{errors.uc_code.message}</p>}
        </div>

        <div className="space-y-1.5">
          <label className="text-sm font-medium">Distribuidora</label>
          <Input {...register('distribuidora')} placeholder="Ex: Neoenergia" />
        </div>

        <div className="space-y-1.5">
          <label className="text-sm font-medium">Apelido</label>
          <Input {...register('apelido')} placeholder="Ex: Sede comercial" />
        </div>

        <div className="space-y-1.5">
          <label className="text-sm font-medium">Classe de Consumo</label>
          <Input {...register('classe_consumo')} placeholder="Ex: Comercial B3" />
        </div>

        <div className="space-y-1.5 sm:col-span-2">
          <label className="text-sm font-medium">Endereço da Unidade</label>
          <Input {...register('endereco_unidade')} placeholder="Endereço completo da UC" />
        </div>

        <div className="space-y-1.5">
          <label className="text-sm font-medium">Cidade</label>
          <Input {...register('cidade')} placeholder="Cidade" />
        </div>

        <div className="space-y-1.5">
          <label className="text-sm font-medium">UF</label>
          <Input {...register('uf')} placeholder="BA" maxLength={2} className="uppercase" />
        </div>

        <div className="space-y-1.5 sm:col-span-2">
          <label className="text-sm font-medium">Credencial de Integração</label>
          <select
            {...register('credential_id')}
            className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
          >
            <option value="">Sem credencial</option>
            {credentials.map((c) => (
              <option key={c.id} value={c.id}>
                {c.label} — {c.documento_masked}
              </option>
            ))}
          </select>
          <p className="text-xs text-muted-foreground">
            Selecione uma credencial para habilitar a sincronização automática.
          </p>
        </div>
      </div>

      <div className="flex items-center gap-2 pt-2">
        <input type="checkbox" id="ativa" {...register('ativa')} className="accent-primary" />
        <label htmlFor="ativa" className="text-sm">UC ativa</label>
      </div>

      <Button type="submit" disabled={isLoading}>
        {isLoading ? 'Salvando...' : mode === 'create' ? 'Criar UC' : 'Salvar'}
      </Button>
    </form>
  )
}
