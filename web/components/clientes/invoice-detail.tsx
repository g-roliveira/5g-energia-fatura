'use client'

import * as React from 'react'
import { CompletudeBadge } from './completude-badge'
import type { GoInvoice } from '@/types/clientes'

type InvoiceDetailProps = {
  invoice: GoInvoice
}

function FieldRow({ label, value, source }: { label: string; value?: unknown; source?: string }) {
  if (value === undefined || value === null) return null
  return (
    <div className="flex justify-between border-b py-2 text-sm last:border-0">
      <span className="text-muted-foreground">{label}</span>
      <div className="flex items-center gap-2">
        <span className="font-medium">{String(value)}</span>
        {source && (
          <span className="rounded bg-muted px-1.5 py-0.5 text-xs text-muted-foreground">
            {source}
          </span>
        )}
      </div>
    </div>
  )
}

export function InvoiceDetail({ invoice }: InvoiceDetailProps) {
  const billing = invoice.billing_record
  const completeness = billing?.completeness?.status ?? invoice.completeness_status

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-lg font-semibold">
            Fatura {billing?.numero_fatura ?? invoice.id}
          </h2>
          <p className="text-sm text-muted-foreground">{billing?.mes_referencia}</p>
        </div>
        {completeness && (
          <CompletudeBadge status={completeness as 'complete' | 'partial' | 'failed'} />
        )}
      </div>

      {/* Billing Record */}
      <section>
        <h3 className="mb-3 text-sm font-semibold uppercase tracking-wide text-muted-foreground">
          Dados da Fatura
        </h3>
        <div className="rounded-md border">
          <FieldRow label="Número" value={billing?.numero_fatura} />
          <FieldRow label="Referência" value={billing?.mes_referencia} />
          <FieldRow label="Valor Total" value={billing?.valor_total ? `R$ ${billing.valor_total}` : undefined} />
          <FieldRow label="Vencimento" value={billing?.data_vencimento} />
          {billing && Object.entries(billing).map(([k, v]) => {
            if (['numero_fatura','mes_referencia','valor_total','data_vencimento','completeness'].includes(k)) return null
            return <FieldRow key={k} label={k} value={typeof v === 'object' ? JSON.stringify(v) : v as string} />
          })}
        </div>
      </section>

      {/* Document Record */}
      {invoice.document_record && Object.keys(invoice.document_record).length > 0 && (
        <section>
          <h3 className="mb-3 text-sm font-semibold uppercase tracking-wide text-muted-foreground">
            Documento Extraído
          </h3>
          <div className="rounded-md border">
            {Object.entries(invoice.document_record).map(([k, v]) => (
              <FieldRow key={k} label={k} value={typeof v === 'object' ? JSON.stringify(v) : v as string} />
            ))}
          </div>
        </section>
      )}

      {/* Items */}
      {invoice.items && invoice.items.length > 0 && (
        <section>
          <h3 className="mb-3 text-sm font-semibold uppercase tracking-wide text-muted-foreground">
            Itens da Fatura
          </h3>
          <div className="rounded-md border overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b bg-muted/50">
                  <th className="p-2 text-left font-medium">Descrição</th>
                  <th className="p-2 text-right font-medium">Qtd</th>
                  <th className="p-2 text-right font-medium">Valor Unit.</th>
                  <th className="p-2 text-right font-medium">Total</th>
                </tr>
              </thead>
              <tbody>
                {invoice.items.map((item, i) => (
                  <tr key={i} className="border-b last:border-0">
                    <td className="p-2">{item.descricao ?? '—'}</td>
                    <td className="p-2 text-right">{item.quantidade ?? '—'}</td>
                    <td className="p-2 text-right">{item.valor_unitario ? `R$ ${item.valor_unitario}` : '—'}</td>
                    <td className="p-2 text-right">{item.valor_total ? `R$ ${item.valor_total}` : '—'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>
      )}
    </div>
  )
}
