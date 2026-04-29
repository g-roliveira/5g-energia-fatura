import { NextRequest, NextResponse } from 'next/server'
import { goFetch, GoApiError } from '@/lib/go-client'
import type { GoCycleRowsResponse } from '@/types/billing'

type Params = { params: Promise<{ id: string }> }

export async function GET(req: NextRequest, { params }: Params) {
  const { id } = await params
  const search = req.nextUrl.searchParams

  const query = new URLSearchParams()
  if (search.has('q')) query.set('q', search.get('q')!)
  if (search.has('sync_status')) query.set('sync_status', search.get('sync_status')!)
  if (search.has('calc_status')) query.set('calc_status', search.get('calc_status')!)
  if (search.has('needs_review_only')) query.set('needs_review_only', search.get('needs_review_only')!)
  if (search.has('limit')) query.set('limit', search.get('limit')!)
  if (search.has('offset')) query.set('offset', search.get('offset')!)

  try {
    const result = await goFetch<GoCycleRowsResponse>(
      `/v1/billing/cycles/${id}/rows?${query.toString()}`,
    )
    return NextResponse.json(result)
  } catch (err) {
    if (err instanceof GoApiError) {
      return NextResponse.json({ error: err.message }, { status: err.status })
    }
    return NextResponse.json({ error: 'Serviço de faturamento indisponível' }, { status: 503 })
  }
}
