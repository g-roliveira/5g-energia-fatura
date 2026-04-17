import { NextRequest, NextResponse } from 'next/server'
import { db } from '@/lib/db'
import { CreateUcSchema } from '@/types/clientes'
import { Prisma } from '@prisma/client'

type Params = { params: Promise<{ id: string }> }

export async function GET(_req: NextRequest, { params }: Params) {
  const { id } = await params

  const ucs = await db.consumerUnit.findMany({
    where: { client_id: id },
    include: {
      credential: {
        select: { id: true, label: true, documento_masked: true },
      },
    },
    orderBy: { created_at: 'desc' },
  })

  return NextResponse.json(ucs)
}

export async function POST(req: NextRequest, { params }: Params) {
  const { id } = await params
  const body = await req.json()

  const parsed = CreateUcSchema.safeParse(body)
  if (!parsed.success) {
    return NextResponse.json({ error: 'Dados inválidos', details: parsed.error.issues }, { status: 422 })
  }

  const existing = await db.client.findUnique({ where: { id } })
  if (!existing) return NextResponse.json({ error: 'Cliente não encontrado' }, { status: 404 })
  if (existing.archived_at) return NextResponse.json({ error: 'Cliente arquivado não pode receber novas UCs' }, { status: 422 })

  try {
    const uc = await db.consumerUnit.create({
      data: { ...parsed.data, client_id: id },
    })
    return NextResponse.json(uc, { status: 201 })
  } catch (err) {
    if (err instanceof Prisma.PrismaClientKnownRequestError && err.code === 'P2002') {
      return NextResponse.json({ error: 'UC já cadastrada' }, { status: 409 })
    }
    throw err
  }
}
