import { NextRequest, NextResponse } from 'next/server'
import { db } from '@/lib/db'

type Params = { params: Promise<{ id: string }> }

export async function POST(_req: NextRequest, { params }: Params) {
  const { id } = await params

  const existing = await db.client.findUnique({ where: { id } })
  if (!existing) return NextResponse.json({ error: 'Cliente não encontrado' }, { status: 404 })

  const client = await db.client.update({
    where: { id },
    data: {
      archived_at: existing.archived_at ?? new Date(),
      status: 'inativo',
    },
  })

  return NextResponse.json(client)
}
