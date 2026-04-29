import { NextRequest, NextResponse } from 'next/server'
import { db } from '@/lib/db'
import { Prisma } from '@prisma/client'
import { z } from 'zod'

type Params = { params: Promise<{ id: string }> }

const BulkUcItemSchema = z.object({
  uc_code: z.string().min(1),
  distribuidora: z.string().optional(),
  endereco_unidade: z.string().optional(),
  cidade: z.string().optional(),
  uf: z.string().length(2).optional(),
  credential_id: z.string().optional(),
  ativa: z.boolean().default(true),
})

const BulkImportSchema = z.object({
  ucs: z.array(BulkUcItemSchema).min(1),
})

export type BulkUcResult = {
  uc_code: string
  status: 'imported' | 'already_exists' | 'error'
  id?: string
  error?: string
}

export async function POST(req: NextRequest, { params }: Params) {
  const { id } = await params

  const existing = await db.client.findUnique({ where: { id } })
  if (!existing) return NextResponse.json({ error: 'Cliente não encontrado' }, { status: 404 })
  if (existing.archived_at) return NextResponse.json({ error: 'Cliente arquivado' }, { status: 422 })

  const body = await req.json()
  const parsed = BulkImportSchema.safeParse(body)
  if (!parsed.success) {
    return NextResponse.json({ error: 'Dados inválidos', details: parsed.error.issues }, { status: 422 })
  }

  const results: BulkUcResult[] = await Promise.all(
    parsed.data.ucs.map(async (ucData) => {
      try {
        const uc = await db.consumerUnit.create({
          data: { ...ucData, client_id: id },
        })
        return { uc_code: ucData.uc_code, status: 'imported' as const, id: uc.id }
      } catch (err) {
        if (err instanceof Prisma.PrismaClientKnownRequestError && err.code === 'P2002') {
          return { uc_code: ucData.uc_code, status: 'already_exists' as const }
        }
        return { uc_code: ucData.uc_code, status: 'error' as const, error: 'Erro interno' }
      }
    }),
  )

  return NextResponse.json({ results })
}
