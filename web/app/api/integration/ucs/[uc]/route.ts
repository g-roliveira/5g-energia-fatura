import { NextRequest, NextResponse } from 'next/server'
import { goFetch, GoApiError } from '@/lib/go-client'
import type { GoConsumerUnit } from '@/types/clientes'

type Params = { params: Promise<{ uc: string }> }

export async function GET(_req: NextRequest, { params }: Params) {
  const { uc } = await params

  try {
    const data = await goFetch<GoConsumerUnit>(`/v1/consumer-units/${uc}`)
    return NextResponse.json(data)
  } catch (err) {
    if (err instanceof GoApiError) {
      return NextResponse.json({ error: err.message }, { status: err.status })
    }
    return NextResponse.json({ error: 'Serviço de integração indisponível' }, { status: 503 })
  }
}
