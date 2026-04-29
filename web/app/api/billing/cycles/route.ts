import { NextRequest, NextResponse } from 'next/server'
import { goFetch, GoApiError } from '@/lib/go-client'
import { CreateCycleSchema, CycleListQuerySchema } from '@/types/billing'
import type { GoBillingCycle, GoCycleListResponse } from '@/types/billing'

export async function GET(req: NextRequest) {
  const params = Object.fromEntries(req.nextUrl.searchParams)
  const parsed = CycleListQuerySchema.safeParse(params)
  if (!parsed.success) {
    return NextResponse.json(
      { error: 'Parâmetros inválidos', details: parsed.error.issues },
      { status: 422 },
    )
  }

  const { year, status, limit, offset } = parsed.data
  const query = new URLSearchParams()
  if (year) query.set('year', String(year))
  if (status) query.set('status', status)
  query.set('limit', String(limit))
  query.set('offset', String(offset))

  try {
    const result = await goFetch<GoCycleListResponse>(`/v1/billing/cycles?${query.toString()}`)
    return NextResponse.json(result)
  } catch (err) {
    if (err instanceof GoApiError) {
      return NextResponse.json({ error: err.message }, { status: err.status })
    }
    return NextResponse.json({ error: 'Serviço de faturamento indisponível' }, { status: 503 })
  }
}

export async function POST(req: NextRequest) {
  const body = await req.json()
  const parsed = CreateCycleSchema.safeParse(body)
  if (!parsed.success) {
    return NextResponse.json(
      { error: 'Dados inválidos', details: parsed.error.issues },
      { status: 422 },
    )
  }

  try {
    const cycle = await goFetch<GoBillingCycle>('/v1/billing/cycles', {
      method: 'POST',
      body: JSON.stringify(parsed.data),
    })
    return NextResponse.json(cycle, { status: 201 })
  } catch (err) {
    if (err instanceof GoApiError) {
      return NextResponse.json({ error: err.message }, { status: err.status })
    }
    return NextResponse.json({ error: 'Serviço de faturamento indisponível' }, { status: 503 })
  }
}
