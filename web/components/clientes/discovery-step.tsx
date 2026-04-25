'use client'

import { useState } from 'react'
import { HugeiconsIcon } from '@hugeicons/react'
import { CheckmarkCircle01Icon, Alert01Icon, Cancel01Icon } from '@hugeicons/core-free-icons'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import type { GoDiscoveryResult, GoDiscoveryUC } from '@/types/clientes'

type ImportStatus = 'imported' | 'already_exists' | 'error'

type DiscoveryStepProps = {
  clientId: string
  credentialId: string
  discovery: GoDiscoveryResult
  currentClient: { nome_razao: string; email?: string | null; telefone?: string | null }
  onImportComplete: () => void
  onClose: () => void
}

function FieldCompare({ label, current, discovered }: { label: string; current?: string | null; discovered?: string }) {
  if (!discovered) return null
  return (
    <div className="grid grid-cols-3 items-center gap-2 py-1 text-sm">
      <span className="text-muted-foreground">{label}</span>
      <span className="text-muted-foreground line-through">{current || '—'}</span>
      <span className="font-medium">{discovered}</span>
    </div>
  )
}

export function DiscoveryStep({
  clientId,
  credentialId,
  discovery,
  currentClient,
  onImportComplete,
  onClose,
}: DiscoveryStepProps) {
  const mc = discovery.minha_conta
  const ucs = discovery.ucs

  const [selected, setSelected] = useState<Set<string>>(new Set(ucs.map((u) => u.uc)))
  const [importResults, setImportResults] = useState<Record<string, ImportStatus>>({})
  const [isImporting, setIsImporting] = useState(false)
  const [importDone, setImportDone] = useState(false)
  const [isUpdatingClient, setIsUpdatingClient] = useState(false)
  const [clientUpdated, setClientUpdated] = useState(false)

  const allSelected = selected.size === ucs.length
  const noneSelected = selected.size === 0

  function toggleAll() {
    if (allSelected) {
      setSelected(new Set())
    } else {
      setSelected(new Set(ucs.map((u) => u.uc)))
    }
  }

  function toggleUc(uc: string) {
    setSelected((prev) => {
      const next = new Set(prev)
      if (next.has(uc)) next.delete(uc)
      else next.add(uc)
      return next
    })
  }

  async function handleImport() {
    const toImport = ucs.filter((u) => selected.has(u.uc))
    if (toImport.length === 0) return

    setIsImporting(true)
    try {
      const res = await fetch(`/api/clients/${clientId}/ucs/bulk`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          ucs: toImport.map((u) => ({
            uc_code: u.uc,
            distribuidora: 'Neoenergia',
            endereco_unidade: u.local.endereco,
            cidade: u.local.municipio,
            uf: u.local.uf.toUpperCase().slice(0, 2),
            credential_id: credentialId,
            ativa: true,
          })),
        }),
      })
      const data = await res.json()
      const resultMap: Record<string, ImportStatus> = {}
      for (const r of data.results ?? []) {
        resultMap[r.uc_code] = r.status
      }
      setImportResults(resultMap)
      setImportDone(true)
      onImportComplete()
    } finally {
      setIsImporting(false)
    }
  }

  async function handleUpdateClient() {
    if (!mc) return
    setIsUpdatingClient(true)
    try {
      const patch: Record<string, string> = {}
      if (mc.nome) patch.nome_razao = mc.nome
      if (mc.email) patch.email = mc.email
      if (mc.celular) patch.telefone = mc.celular
      await fetch(`/api/clients/${clientId}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(patch),
      })
      setClientUpdated(true)
      onImportComplete()
    } finally {
      setIsUpdatingClient(false)
    }
  }

  const importedCount = Object.values(importResults).filter((s) => s === 'imported').length
  const existsCount = Object.values(importResults).filter((s) => s === 'already_exists').length

  return (
    <div className="flex flex-col gap-5">
      {/* Client data section */}
      {mc && (
        <div className="space-y-3">
          <div>
            <p className="text-sm font-medium">Dados do titular</p>
            <p className="text-xs text-muted-foreground">Comparação: valor atual → descoberto</p>
          </div>
          <div className="rounded-lg border p-3">
            <div className="grid grid-cols-3 gap-2 pb-1 text-xs font-medium text-muted-foreground">
              <span>Campo</span><span>Atual</span><span>Descoberto</span>
            </div>
            <Separator className="mb-2" />
            <FieldCompare label="Nome" current={currentClient.nome_razao} discovered={mc.nome} />
            <FieldCompare label="E-mail" current={currentClient.email} discovered={mc.email} />
            <FieldCompare label="Telefone" current={currentClient.telefone} discovered={mc.celular} />
          </div>
          {clientUpdated ? (
            <p className="flex items-center gap-1.5 text-sm text-green-600 dark:text-green-400">
              <HugeiconsIcon icon={CheckmarkCircle01Icon} strokeWidth={2} className="size-4" />
              Dados do cliente atualizados
            </p>
          ) : (
            <Button
              variant="outline"
              size="sm"
              onClick={handleUpdateClient}
              disabled={isUpdatingClient}
            >
              {isUpdatingClient ? 'Atualizando...' : 'Atualizar dados do cliente'}
            </Button>
          )}
        </div>
      )}

      {!mc && !discovery.errors?.minha_conta && null}
      {!mc && discovery.errors?.minha_conta && (
        <p className="text-sm text-muted-foreground">
          Não foi possível carregar os dados do titular para este tipo de acesso.
        </p>
      )}

      <Separator />

      {/* UC list section */}
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm font-medium">Unidades Consumidoras encontradas</p>
            <p className="text-xs text-muted-foreground">{ucs.length} UC{ucs.length !== 1 ? 's' : ''} no portal</p>
          </div>
          {!importDone && (
            <Button variant="ghost" size="sm" onClick={toggleAll}>
              {allSelected ? 'Desmarcar todas' : 'Selecionar todas'}
            </Button>
          )}
        </div>

        {ucs.length === 0 ? (
          <p className="py-6 text-center text-sm text-muted-foreground">Nenhuma UC encontrada para esta credencial.</p>
        ) : (
          <div className="max-h-[40vh] overflow-y-auto rounded-lg border divide-y">
            {ucs.map((uc) => {
              const result = importResults[uc.uc]
              const isActive = uc.status?.toLowerCase().includes('ativ') || uc.status === 'LIGADA'
              return (
                <div key={uc.uc} className="flex items-center gap-3 px-4 py-3">
                  {importDone ? (
                    <div className="flex size-5 items-center justify-center">
                      {result === 'imported' && (
                        <HugeiconsIcon icon={CheckmarkCircle01Icon} strokeWidth={2} className="size-5 text-green-600" />
                      )}
                      {result === 'already_exists' && (
                        <HugeiconsIcon icon={Alert01Icon} strokeWidth={2} className="size-5 text-amber-500" />
                      )}
                      {result === 'error' && (
                        <HugeiconsIcon icon={Cancel01Icon} strokeWidth={2} className="size-5 text-destructive" />
                      )}
                      {!result && (
                        <div className="size-5" />
                      )}
                    </div>
                  ) : (
                    <Checkbox
                      id={`uc-${uc.uc}`}
                      checked={selected.has(uc.uc)}
                      onCheckedChange={() => toggleUc(uc.uc)}
                    />
                  )}

                  <Label
                    htmlFor={`uc-${uc.uc}`}
                    className="flex flex-1 cursor-pointer items-start justify-between gap-3"
                  >
                    <div className="space-y-0.5">
                      <div className="flex items-center gap-2">
                        <span className="font-mono text-sm font-medium">{uc.uc}</span>
                        <Badge variant={isActive ? 'default' : 'secondary'} className="text-[10px]">
                          {uc.status}
                        </Badge>
                        {result === 'already_exists' && (
                          <Badge variant="outline" className="border-amber-400 text-amber-600 text-[10px]">
                            Já cadastrada
                          </Badge>
                        )}
                      </div>
                      {uc.local.endereco && (
                        <p className="text-xs text-muted-foreground">
                          {uc.local.endereco}{uc.local.municipio ? `, ${uc.local.municipio}/${uc.local.uf}` : ''}
                        </p>
                      )}
                    </div>
                    {uc.grupoTensao && (
                      <span className="shrink-0 text-xs text-muted-foreground">{uc.grupoTensao}</span>
                    )}
                  </Label>
                </div>
              )
            })}
          </div>
        )}
      </div>

      {/* Import summary */}
      {importDone && (
        <Alert className="border-green-200 bg-green-50 text-green-800 dark:border-green-800 dark:bg-green-950 dark:text-green-200">
          <HugeiconsIcon icon={CheckmarkCircle01Icon} strokeWidth={2} className="size-4" />
          <AlertDescription>
            {importedCount > 0 && <span>{importedCount} UC{importedCount !== 1 ? 's' : ''} importada{importedCount !== 1 ? 's' : ''}. </span>}
            {existsCount > 0 && <span>{existsCount} já existia{existsCount !== 1 ? 'm' : ''} no sistema.</span>}
          </AlertDescription>
        </Alert>
      )}

      {/* Actions */}
      <div className="flex justify-between">
        <Button variant="outline" onClick={onClose}>
          {importDone ? 'Fechar' : 'Pular'}
        </Button>
        {!importDone && ucs.length > 0 && (
          <Button
            onClick={handleImport}
            disabled={noneSelected || isImporting}
          >
            {isImporting
              ? 'Importando...'
              : `Importar ${selected.size} UC${selected.size !== 1 ? 's' : ''} selecionada${selected.size !== 1 ? 's' : ''}`}
          </Button>
        )}
      </div>
    </div>
  )
}
