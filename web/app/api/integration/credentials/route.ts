import { NextRequest, NextResponse } from 'next/server'
import { db } from '@/lib/db'
import { goFetch, GoApiError } from '@/lib/go-client'
import { CreateCredentialSchema } from '@/types/clientes'
import type { GoCredential } from '@/types/clientes'

export async function POST(req: NextRequest) {
  const body = await req.json()
  const parsed = CreateCredentialSchema.safeParse(body)
  if (!parsed.success) {
    return NextResponse.json({ error: 'Dados inválidos', details: parsed.error.issues }, { status: 422 })
  }

  const { client_id, label, documento, senha, uf, tipo_acesso } = parsed.data

  let goCredential: GoCredential
  try {
    goCredential = await goFetch<GoCredential>('/v1/credentials', {
      method: 'POST',
      body: JSON.stringify({ label, documento, senha, uf, tipo_acesso }),
    })
  } catch (err) {
    if (err instanceof GoApiError) {
      return NextResponse.json({ error: err.message }, { status: err.status })
    }
    return NextResponse.json({ error: 'Serviço de integração indisponível' }, { status: 503 })
  }

  // Persist locally — never store senha
  const credential = await db.integrationCredential.create({
    data: {
      client_id: client_id ?? '',
      label,
      documento_masked: goCredential.documento,
      uf,
      tipo_acesso,
      go_credential_id: goCredential.id,
    },
  })

  return NextResponse.json(credential, { status: 201 })
}
