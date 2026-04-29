import { NextRequest, NextResponse } from 'next/server'
import { goFetch, GoApiError } from '@/lib/go-client'
import type { GoContract } from '@/types/billing'

type Params = { params: Promise<{ id: string }> }

export async function GET(_req: NextRequest, { params }: Params) {
  const { id } = await params

  try {
    const contract = await goFetch<GoContract>(`/v1/billing/contracts/${id}`)
    return NextResponse.json(contract)
  } catch (err) {
    if (err instanceof GoApiError) {
      return NextResponse.json({ error: err.message }, { status: err.status })
    }
    return NextResponse.json({ error: 'Serviço de faturamento indisponível' }, { status: 503 })
  }
}
