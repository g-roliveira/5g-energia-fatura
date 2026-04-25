'use client'

import * as React from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { CompletudeBadge } from './completude-badge'
import type { GoInvoice } from '@/types/clientes'

type InvoiceDetailProps = {
  invoice: GoInvoice
}

function FieldRow({ label, value }: { label: string; value?: unknown }) {
  if (value === undefined || value === null) return null
  return (
    <div className="flex justify-between border-b py-2 text-sm last:border-0">
      <span className="text-muted-foreground">{label}</span>
      <span className="font-medium">{String(value)}</span>
    </div>
  )
}

export function InvoiceDetail({ invoice }: InvoiceDetailProps) {
  const billing = invoice.billing_record
  const completeness = billing?.completeness?.status ?? invoice.completeness_status

  return (
    <div className="space-y-6">
      {/* Header summary */}
      <div className="flex items-center justify-between">
        <div>
          <p className="text-lg font-semibold">
            Fatura {billing?.numero_fatura ?? invoice.id}
          </p>
          <p className="text-sm text-muted-foreground">{billing?.mes_referencia}</p>
        </div>
        {completeness && (
          <CompletudeBadge status={completeness as 'complete' | 'partial' | 'failed'} />
        )}
      </div>

      {/* Billing Record */}
      <Card>
        <CardHeader>
          <CardTitle>Dados da Fatura</CardTitle>
        </CardHeader>
        <CardContent>
          <FieldRow label="Número" value={billing?.numero_fatura} />
          <FieldRow label="Referência" value={billing?.mes_referencia} />
          <FieldRow label="Valor Total" value={billing?.valor_total ? `R$ ${billing.valor_total}` : undefined} />
          <FieldRow label="Vencimento" value={billing?.data_vencimento} />
          {billing && Object.entries(billing).map(([k, v]) => {
            if (['numero_fatura', 'mes_referencia', 'valor_total', 'data_vencimento', 'completeness'].includes(k)) return null
            return <FieldRow key={k} label={k} value={typeof v === 'object' ? JSON.stringify(v) : v as string} />
          })}
        </CardContent>
      </Card>

      {/* Document Record */}
      {invoice.document_record && Object.keys(invoice.document_record).length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Documento Extraído</CardTitle>
          </CardHeader>
          <CardContent>
            {Object.entries(invoice.document_record).map(([k, v]) => (
              <FieldRow key={k} label={k} value={typeof v === 'object' ? JSON.stringify(v) : v as string} />
            ))}
          </CardContent>
        </Card>
      )}

      {/* Items */}
      {invoice.items && invoice.items.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Itens da Fatura</CardTitle>
          </CardHeader>
          <CardContent className="p-0">
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b bg-muted/50">
                    <th className="p-3 text-left font-medium">Descrição</th>
                    <th className="p-3 text-right font-medium">Qtd</th>
                    <th className="p-3 text-right font-medium">Valor Unit.</th>
                    <th className="p-3 text-right font-medium">Total</th>
                  </tr>
                </thead>
                <tbody>
                  {invoice.items.map((item, i) => (
                    <tr key={i} className="border-b last:border-0">
                      <td className="p-3">{item.descricao ?? '—'}</td>
                      <td className="p-3 text-right">{item.quantidade ?? '—'}</td>
                      <td className="p-3 text-right">{item.valor_unitario ? `R$ ${item.valor_unitario}` : '—'}</td>
                      <td className="p-3 text-right">{item.valor_total ? `R$ ${item.valor_total}` : '—'}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
