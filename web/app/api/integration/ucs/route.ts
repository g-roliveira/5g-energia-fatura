import { NextRequest, NextResponse } from 'next/server'
import { goFetch, GoApiError } from '@/lib/go-client'
import type { GoConsumerUnit } from '@/types/clientes'

export async function GET(req: NextRequest) {
  const { searchParams } = req.nextUrl
  const query = new URLSearchParams()
  if (searchParams.get('limit')) query.set('limit', searchParams.get('limit')!)
  if (searchParams.get('status')) query.set('status', searchParams.get('status')!)

  try {
    const data = await goFetch<GoConsumerUnit[]>(`/v1/consumer-units?${query}`)
    return NextResponse.json(data)
  } catch (err) {
    if (err instanceof GoApiError) {
      return NextResponse.json({ error: err.message }, { status: err.status })
    }
    return NextResponse.json({ error: 'Serviço de integração indisponível' }, { status: 503 })
  }
}
