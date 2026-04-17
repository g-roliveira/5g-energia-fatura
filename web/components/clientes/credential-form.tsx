'use client'

import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { CreateCredentialSchema, type CreateCredentialInput } from '@/types/clientes'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

type CredentialFormProps = {
  clientId: string
  onSubmit: (data: CreateCredentialInput) => Promise<void>
  isLoading?: boolean
}

export function CredentialForm({ clientId, onSubmit, isLoading }: CredentialFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateCredentialInput>({
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    resolver: zodResolver(CreateCredentialSchema as any),
    defaultValues: { client_id: clientId, tipo_acesso: 'normal' },
  })

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="rounded-md border border-yellow-200 bg-yellow-50 p-3 text-sm text-yellow-800">
        A senha é enviada diretamente para a concessionária e <strong>nunca é armazenada</strong> no sistema.
      </div>

      <div className="space-y-1.5">
        <label className="text-sm font-medium">Identificação / Label *</label>
        <Input {...register('label')} placeholder="Ex: neo-joao-silva" />
        {errors.label && <p className="text-xs text-destructive">{errors.label.message}</p>}
      </div>

      <div className="space-y-1.5">
        <label className="text-sm font-medium">CPF / CNPJ do portal *</label>
        <Input {...register('documento')} placeholder="CPF ou CNPJ usado no portal" />
        {errors.documento && <p className="text-xs text-destructive">{errors.documento.message}</p>}
      </div>

      <div className="space-y-1.5">
        <label className="text-sm font-medium">Senha do portal *</label>
        <Input type="password" {...register('senha')} placeholder="Senha do portal da concessionária" autoComplete="new-password" />
        {errors.senha && <p className="text-xs text-destructive">{errors.senha.message}</p>}
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-1.5">
          <label className="text-sm font-medium">UF *</label>
          <Input {...register('uf')} placeholder="BA" maxLength={2} className="uppercase" />
          {errors.uf && <p className="text-xs text-destructive">{errors.uf.message}</p>}
        </div>
        <div className="space-y-1.5">
          <label className="text-sm font-medium">Tipo de Acesso</label>
          <select
            {...register('tipo_acesso')}
            className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm"
          >
            <option value="normal">Normal</option>
          </select>
        </div>
      </div>

      <input type="hidden" {...register('client_id')} />

      <Button type="submit" disabled={isLoading}>
        {isLoading ? 'Criando...' : 'Criar credencial'}
      </Button>
    </form>
  )
}
