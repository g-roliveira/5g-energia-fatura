import { NextRequest, NextResponse } from 'next/server'
import { goFetch, GoApiError } from '@/lib/go-client'
import type { GoInvoice } from '@/types/clientes'

type Params = { params: Promise<{ uc: string }> }

export async function GET(req: NextRequest, { params }: Params) {
  const { uc } = await params
  const { searchParams } = req.nextUrl
  const query = new URLSearchParams()
  if (searchParams.get('limit')) query.set('limit', searchParams.get('limit')!)
  if (searchParams.get('status')) query.set('status', searchParams.get('status')!)

  try {
    const data = await goFetch<GoInvoice[]>(`/v1/consumer-units/${uc}/invoices?${query}`)
    return NextResponse.json(data)
  } catch (err) {
    if (err instanceof GoApiError) {
      return NextResponse.json({ error: err.message }, { status: err.status })
    }
    return NextResponse.json({ error: 'Serviço de integração indisponível' }, { status: 503 })
  }
}
