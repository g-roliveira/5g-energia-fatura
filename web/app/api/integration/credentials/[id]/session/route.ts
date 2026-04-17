import { NextRequest, NextResponse } from 'next/server'
import { db } from '@/lib/db'
import { goFetch, GoApiError } from '@/lib/go-client'
import type { GoSession } from '@/types/clientes'

type Params = { params: Promise<{ id: string }> }

export async function POST(_req: NextRequest, { params }: Params) {
  const { id } = await params

  const credential = await db.integrationCredential.findUnique({ where: { id } })
  if (!credential) return NextResponse.json({ error: 'Credencial não encontrada' }, { status: 404 })

  try {
    const session = await goFetch<GoSession>(
      `/v1/credentials/${credential.go_credential_id}/session`,
      { method: 'POST', body: '{}' },
    )
    return NextResponse.json(session)
  } catch (err) {
    if (err instanceof GoApiError) {
      return NextResponse.json({ error: err.message }, { status: err.status })
    }
    return NextResponse.json({ error: 'Serviço de integração indisponível' }, { status: 503 })
  }
}
