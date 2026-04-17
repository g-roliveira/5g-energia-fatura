import { NextRequest, NextResponse } from 'next/server'
import { db } from '@/lib/db'
import { UpdateUcSchema } from '@/types/clientes'

type Params = { params: Promise<{ id: string }> }

export async function PATCH(req: NextRequest, { params }: Params) {
  const { id } = await params
  const body = await req.json()

  const parsed = UpdateUcSchema.safeParse(body)
  if (!parsed.success) {
    return NextResponse.json({ error: 'Dados inválidos', details: parsed.error.issues }, { status: 422 })
  }

  const existing = await db.consumerUnit.findUnique({ where: { id } })
  if (!existing) return NextResponse.json({ error: 'UC não encontrada' }, { status: 404 })

  const uc = await db.consumerUnit.update({
    where: { id },
    data: parsed.data,
  })

  return NextResponse.json(uc)
}
