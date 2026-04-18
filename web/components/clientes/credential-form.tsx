'use client'

import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { CreateCredentialSchema, type CreateCredentialInput } from '@/types/clientes'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { NativeSelect, NativeSelectOption } from '@/components/ui/native-select'
import { HugeiconsIcon } from '@hugeicons/react'
import { SecurityLockIcon } from '@hugeicons/core-free-icons'

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
      <Alert className="border-amber-200 bg-amber-50 text-amber-800 dark:border-amber-800 dark:bg-amber-950 dark:text-amber-200">
        <HugeiconsIcon icon={SecurityLockIcon} strokeWidth={2} className="size-4" />
        <AlertDescription>
          A senha é enviada diretamente para a concessionária e <strong>nunca é armazenada</strong> no sistema.
        </AlertDescription>
      </Alert>

      <div className="space-y-1.5">
        <Label htmlFor="label">Identificação / Label *</Label>
        <Input id="label" {...register('label')} placeholder="Ex: neo-joao-silva" />
        {errors.label && <p className="text-xs text-destructive">{errors.label.message}</p>}
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="documento">CPF / CNPJ do portal *</Label>
        <Input id="documento" {...register('documento')} placeholder="CPF ou CNPJ usado no portal" />
        {errors.documento && <p className="text-xs text-destructive">{errors.documento.message}</p>}
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="senha">Senha do portal *</Label>
        <Input id="senha" type="password" {...register('senha')} placeholder="Senha do portal da concessionária" autoComplete="new-password" />
        {errors.senha && <p className="text-xs text-destructive">{errors.senha.message}</p>}
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-1.5">
          <Label htmlFor="uf">UF *</Label>
          <Input id="uf" {...register('uf')} placeholder="BA" maxLength={2} className="uppercase" />
          {errors.uf && <p className="text-xs text-destructive">{errors.uf.message}</p>}
        </div>
        <div className="space-y-1.5">
          <Label htmlFor="tipo_acesso">Tipo de Acesso</Label>
          <NativeSelect id="tipo_acesso" {...register('tipo_acesso')} className="w-full">
            <NativeSelectOption value="normal">Normal</NativeSelectOption>
            <NativeSelectOption value="procurador">Procurador</NativeSelectOption>
          </NativeSelect>
        </div>
      </div>

      <input type="hidden" {...register('client_id')} />

      <Button type="submit" disabled={isLoading}>
        {isLoading ? 'Criando...' : 'Criar credencial'}
      </Button>
    </form>
  )
}
