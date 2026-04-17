import { NextRequest, NextResponse } from 'next/server'
import { db } from '@/lib/db'
import { goFetch, GoApiError } from '@/lib/go-client'
import type { GoPersistence } from '@/types/clientes'

type Params = { params: Promise<{ uc: string }> }

interface SyncResponse {
  uc: string
  persistence: GoPersistence
}

export async function POST(req: NextRequest, { params }: Params) {
  const { uc } = await params
  const body = await req.json()
  const { credential_id } = body as { credential_id?: string }

  if (!credential_id) {
    return NextResponse.json({ error: 'credential_id é obrigatório' }, { status: 422 })
  }

  const credential = await db.integrationCredential.findUnique({ where: { id: credential_id } })
  if (!credential) return NextResponse.json({ error: 'Credencial não encontrada' }, { status: 404 })

  try {
    const result = await goFetch<SyncResponse>(
      `/v1/consumer-units/${uc}/sync`,
      {
        method: 'POST',
        body: JSON.stringify({
          credential_id: credential.go_credential_id,
          include_pdf: true,
          include_extraction: true,
        }),
        timeoutMs: 60_000,
      },
    )
    return NextResponse.json(result.persistence)
  } catch (err) {
    if (err instanceof GoApiError) {
      return NextResponse.json({ error: err.message }, { status: err.status })
    }
    return NextResponse.json({ error: 'Serviço de integração indisponível' }, { status: 503 })
  }
}
