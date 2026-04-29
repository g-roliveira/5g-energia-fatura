import { NextRequest, NextResponse } from 'next/server'
import { goFetch, GoApiError } from '@/lib/go-client'
import { CloseCycleSchema } from '@/types/billing'

type Params = { params: Promise<{ id: string }> }

export async function POST(req: NextRequest, { params }: Params) {
  const { id } = await params
  const body = await req.json().catch(() => ({}))

  const parsed = CloseCycleSchema.safeParse(body)
  if (!parsed.success) {
    return NextResponse.json(
      { error: 'Dados inválidos', details: parsed.error.issues },
      { status: 422 },
    )
  }

  try {
    const result = await goFetch<{ status: string }>(`/v1/billing/cycles/${id}/close`, {
      method: 'POST',
      body: JSON.stringify(parsed.data),
    })
    return NextResponse.json(result)
  } catch (err) {
    if (err instanceof GoApiError) {
      return NextResponse.json({ error: err.message }, { status: err.status })
    }
    return NextResponse.json({ error: 'Serviço de faturamento indisponível' }, { status: 503 })
  }
}
