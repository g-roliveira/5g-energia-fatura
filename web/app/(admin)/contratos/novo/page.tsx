'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { HugeiconsIcon } from '@hugeicons/react'
import { ArrowLeft01Icon } from '@hugeicons/core-free-icons'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Checkbox } from '@/components/ui/checkbox'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { useQuery } from '@tanstack/react-query'
import { useCreateContract } from '@/hooks/use-billing'
import { toast } from '@/hooks/use-toast'

export default function NovoContratoPage() {
  const router = useRouter()
  const createContract = useCreateContract()

  const [clientId, setClientId] = useState('')
  const [ucId, setUcId] = useState('')
  const [vigenciaInicio, setVigenciaInicio] = useState('')
  const [desconto, setDesconto] = useState('0.85')
  const [valorIpComDesc, setValorIpComDesc] = useState('13.35')
  const [ipMode, setIpMode] = useState<'fixed' | 'percent'>('fixed')
  const [ipValor, setIpValor] = useState('10.00')
  const [ipPercent, setIpPercent] = useState('0.10')
  const [bandeiraDesc, setBandeiraDesc] = useState(false)
  const [dispSempre, setDispSempre] = useState(false)
  const [notes, setNotes] = useState('')

  const { data: clients } = useQuery({
    queryKey: ['clients-minimal'],
    queryFn: async () => {
      const res = await fetch('/api/clients?page=1&pageSize=100')
      if (!res.ok) throw new Error('Erro ao carregar clientes')
      return res.json() as Promise<{ data: Array<{ id: string; nome_razao: string }> }>
    },
  })

  const { data: ucs } = useQuery({
    queryKey: ['ucs', clientId],
    queryFn: async () => {
      const res = await fetch(`/api/clients/${clientId}/ucs`)
      if (!res.ok) throw new Error('Erro ao carregar UCs')
      return res.json() as Promise<Array<{ id: string; uc_code: string; apelido?: string }>>
    },
    enabled: !!clientId,
  })

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!clientId || !ucId || !vigenciaInicio) {
      toast({ title: 'Erro', description: 'Preencha todos os campos obrigatórios', variant: 'destructive' })
      return
    }

    try {
      await createContract.mutateAsync({
        customer_id: clientId,
        consumer_unit_id: ucId,
        vigencia_inicio: vigenciaInicio,
        fator_repasse_energia: desconto,
	valor_ip_com_desconto: valorIpComDesc,
        ip_faturamento_mode: ipMode,
        ip_faturamento_valor: ipMode === 'fixed' ? ipValor : undefined,
        ip_faturamento_percent: ipMode === 'percent' ? ipPercent : undefined,
        bandeira_com_desconto: bandeiraDesc,
        custo_disponibilidade_sempre_cobrado: dispSempre,
        notes: notes || undefined,
      })
      toast({ title: 'Contrato criado com sucesso' })
      router.push('/contratos')
    } catch (err: any) {
      toast({ title: 'Erro ao criar contrato', description: err.message, variant: 'destructive' })
    }
  }

  return (
    <div className="flex flex-col gap-6 p-6 max-w-3xl mx-auto">
      <div className="flex items-center gap-2">
        <Button variant="ghost" size="icon" onClick={() => router.push('/contratos')}>
          <HugeiconsIcon icon={ArrowLeft01Icon} strokeWidth={2} className="size-4" />
        </Button>
        <div>
          <h1 className="text-2xl font-semibold">Novo Contrato</h1>
          <p className="text-sm text-muted-foreground">Crie um novo contrato de faturamento</p>
        </div>
      </div>

      <form onSubmit={handleSubmit}>
        <Card>
          <CardHeader>
            <CardTitle>Dados do contrato</CardTitle>
            <CardDescription>Informe os termos do contrato de faturamento</CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            {/* Cliente e UC */}
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="client">Cliente *</Label>
                <Select value={clientId} onValueChange={(v) => { setClientId(v); setUcId('') }}>
                  <SelectTrigger id="client">
                    <SelectValue placeholder="Selecione" />
                  </SelectTrigger>
                  <SelectContent>
                    {clients?.data.map((c) => (
                      <SelectItem key={c.id} value={c.id}>{c.nome_razao}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label htmlFor="uc">Unidade Consumidora *</Label>
                <Select value={ucId} onValueChange={setUcId} disabled={!clientId}>
                  <SelectTrigger id="uc">
                    <SelectValue placeholder="Selecione" />
                  </SelectTrigger>
                  <SelectContent>
                    {ucs?.map((uc) => (
                      <SelectItem key={uc.id} value={uc.id}>
                        {uc.uc_code} {uc.apelido ? `— ${uc.apelido}` : ''}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>

            {/* Vigência e Desconto */}
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="vigencia">Vigência Início *</Label>
                <Input
                  id="vigencia"
                  type="date"
                  value={vigenciaInicio}
                  onChange={(e) => setVigenciaInicio(e.target.value)}
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="desconto">Fator de Repasse *</Label>
                <Input
                  id="desconto"
                  type="number"
                  step="0.01"
                  min="0.01"
                  max="1"
                  value={desconto}
                  onChange={(e) => setDesconto(e.target.value)}
                  required
                />
                <p className="text-xs text-muted-foreground">Ex: 0.85 = repasse de 85% (5G fica com a diferença)</p>
              </div>
            </div>

            {/* IP Faturamento */}
            <div className="space-y-4 rounded-lg border p-4">
              <Label>IP Faturamento (Iluminação Pública da Usina)</Label>
              <RadioGroup
                value={ipMode}
                onValueChange={(v: string) => setIpMode(v as 'fixed' | 'percent')}
                className="flex items-center gap-4"
              >
                <div className="flex items-center gap-2">
                  <RadioGroupItem value="fixed" id="ip-fixed" />
                  <Label htmlFor="ip-fixed" className="font-normal">Valor fixo</Label>
                </div>
                <div className="flex items-center gap-2">
                  <RadioGroupItem value="percent" id="ip-percent" />
                  <Label htmlFor="ip-percent" className="font-normal">Percentual</Label>
                </div>
              </RadioGroup>

              {ipMode === 'fixed' ? (
                <div className="space-y-2">
                  <Label htmlFor="ip-valor">Valor (R$)</Label>
                  <Input
                    id="ip-valor"
                    type="number"
                    step="0.01"
                    min="0"
                    value={ipValor}
                    onChange={(e) => setIpValor(e.target.value)}
                  />
                </div>
              ) : (
                <div className="space-y-2">
                  <Label htmlFor="ip-percent">Percentual</Label>
                  <Input
                    id="ip-percent"
                    type="number"
                    step="0.01"
                    min="0"
                    max="1"
                    value={ipPercent}
                    onChange={(e) => setIpPercent(e.target.value)}
                  />
                </div>
              )}
            </div>

            {/* Checkboxes */}
            <div className="flex items-center gap-6">
              <div className="flex items-center gap-2">
                <Checkbox
                  id="bandeira"
                  checked={bandeiraDesc}
                  onCheckedChange={(v) => setBandeiraDesc(v === true)}
                />
                <Label htmlFor="bandeira" className="font-normal">Bandeira com desconto</Label>
              </div>
              <div className="flex items-center gap-2">
                <Checkbox
                  id="disp"
                  checked={dispSempre}
                  onCheckedChange={(v) => setDispSempre(v === true)}
                />
                <Label htmlFor="disp" className="font-normal">Custo disponibilidade sempre cobrado</Label>
              </div>
            </div>

            {/* Notes */}
            <div className="space-y-2">
              <Label htmlFor="notes">Observações</Label>
              <Textarea
                id="notes"
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                placeholder="Observações adicionais..."
                rows={3}
              />
            </div>

            <div className="flex justify-end gap-2">
              <Button type="button" variant="outline" onClick={() => router.push('/contratos')}>
                Cancelar
              </Button>
              <Button type="submit" disabled={createContract.isPending}>
                {createContract.isPending ? 'Criando...' : 'Criar contrato'}
              </Button>
            </div>
          </CardContent>
        </Card>
      </form>
    </div>
  )
}
