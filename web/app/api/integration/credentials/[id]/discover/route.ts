import { NextRequest, NextResponse } from 'next/server'
import { db } from '@/lib/db'
import { goFetch, GoApiError } from '@/lib/go-client'
import type { GoDiscoveryResult } from '@/types/clientes'

type Params = { params: Promise<{ id: string }> }

export async function GET(_req: NextRequest, { params }: Params) {
  const { id } = await params

  const credential = await db.integrationCredential.findUnique({ where: { id } })
  if (!credential) {
    return NextResponse.json({ error: 'Credencial não encontrada' }, { status: 404 })
  }

  try {
    const result = await goFetch<GoDiscoveryResult>(
      `/v1/credentials/${credential.go_credential_id}/discover`,
      { timeoutMs: 60_000 },
    )
    return NextResponse.json(result)
  } catch (err) {
    if (err instanceof GoApiError) {
      return NextResponse.json({ error: err.message }, { status: err.status })
    }
    return NextResponse.json({ error: 'Serviço de integração indisponível' }, { status: 503 })
  }
}
