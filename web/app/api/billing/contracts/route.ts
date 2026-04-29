import { NextRequest, NextResponse } from 'next/server'
import { goFetch, GoApiError } from '@/lib/go-client'
import type { GoContract } from '@/types/billing'

export async function GET(req: NextRequest) {
  const ucId = req.nextUrl.searchParams.get('uc_id')
  if (!ucId) {
    return NextResponse.json({ error: 'uc_id é obrigatório' }, { status: 422 })
  }

  try {
    const result = await goFetch<{ items: GoContract[]; count: number }>(
      `/v1/billing/contracts?uc_id=${encodeURIComponent(ucId)}`,
    )
    return NextResponse.json(result)
  } catch (err) {
    if (err instanceof GoApiError) {
      return NextResponse.json({ error: err.message }, { status: err.status })
    }
    return NextResponse.json({ error: 'Serviço de faturamento indisponível' }, { status: 503 })
  }
}
