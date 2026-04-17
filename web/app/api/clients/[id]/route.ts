import { NextRequest, NextResponse } from 'next/server'
import { db } from '@/lib/db'
import { UpdateClientSchema } from '@/types/clientes'
import { Prisma } from '@prisma/client'

type Params = { params: Promise<{ id: string }> }

export async function GET(_req: NextRequest, { params }: Params) {
  const { id } = await params

  const client = await db.client.findUnique({
    where: { id },
    include: {
      address: true,
      ucs: {
        include: { credential: true },
      },
      commercial_data: true,
      credentials: true,
    },
  })

  if (!client) return NextResponse.json({ error: 'Cliente não encontrado' }, { status: 404 })
  return NextResponse.json(client)
}

export async function PATCH(req: NextRequest, { params }: Params) {
  const { id } = await params
  const body = await req.json()

  const parsed = UpdateClientSchema.safeParse(body)
  if (!parsed.success) {
    return NextResponse.json({ error: 'Dados inválidos', details: parsed.error.issues }, { status: 422 })
  }

  const {
    cep, logradouro, numero, complemento, bairro, cidade, uf,
    tipo_contrato, data_inicio, data_fim, status_contrato, observacoes_comerciais,
    ...clientData
  } = parsed.data

  const addressFields = { cep, logradouro, numero, complemento, bairro, cidade, uf }
  const hasAddress = Object.values(addressFields).some((v) => v !== undefined)

  const commercialFields = { tipo_contrato, data_inicio, data_fim, status_contrato, observacoes_comerciais }
  const hasCommercial = Object.values(commercialFields).some((v) => v !== undefined)

  try {
    const existing = await db.client.findUnique({ where: { id } })
    if (!existing) return NextResponse.json({ error: 'Cliente não encontrado' }, { status: 404 })

    const client = await db.client.update({
      where: { id },
      data: {
        ...clientData,
        ...(hasAddress ? {
          address: {
            upsert: {
              create: addressFields,
              update: addressFields,
            },
          },
        } : {}),
        ...(hasCommercial ? {
          commercial_data: {
            upsert: {
              create: { ...commercialFields, data_inicio: commercialFields.data_inicio ? new Date(commercialFields.data_inicio) : undefined, data_fim: commercialFields.data_fim ? new Date(commercialFields.data_fim) : undefined },
              update: { ...commercialFields, data_inicio: commercialFields.data_inicio ? new Date(commercialFields.data_inicio) : undefined, data_fim: commercialFields.data_fim ? new Date(commercialFields.data_fim) : undefined },
            },
          },
        } : {}),
      },
      include: { address: true, commercial_data: true },
    })

    return NextResponse.json(client)
  } catch (err) {
    if (err instanceof Prisma.PrismaClientKnownRequestError && err.code === 'P2002') {
      return NextResponse.json({ error: 'CPF/CNPJ já cadastrado' }, { status: 409 })
    }
    throw err
  }
}
