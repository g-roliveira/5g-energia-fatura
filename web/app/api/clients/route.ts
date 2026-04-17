import { NextRequest, NextResponse } from 'next/server'
import { db } from '@/lib/db'
import { CreateClientSchema, ClientListQuerySchema } from '@/types/clientes'
import { Prisma } from '@prisma/client'

export async function GET(req: NextRequest) {
  const params = Object.fromEntries(req.nextUrl.searchParams)
  const parsed = ClientListQuerySchema.safeParse(params)
  if (!parsed.success) {
    return NextResponse.json({ error: 'Parâmetros inválidos', details: parsed.error.issues }, { status: 422 })
  }

  const { page, pageSize, search, status, tipo_cliente, archived, orderBy, order } = parsed.data
  const skip = (page - 1) * pageSize

  const where: Prisma.ClientWhereInput = {
    ...(archived ? {} : { archived_at: null }),
    ...(status ? { status } : {}),
    ...(tipo_cliente ? { tipo_cliente } : {}),
    ...(search
      ? {
          OR: [
            { nome_razao: { contains: search, mode: 'insensitive' } },
            { cpf_cnpj: { contains: search } },
          ],
        }
      : {}),
  }

  const [data, total] = await Promise.all([
    db.client.findMany({
      where,
      skip,
      take: pageSize,
      orderBy: { [orderBy]: order },
      include: {
        address: true,
        _count: { select: { ucs: true } },
      },
    }),
    db.client.count({ where }),
  ])

  return NextResponse.json({ data, total, page, pageSize })
}

export async function POST(req: NextRequest) {
  const body = await req.json()
  const parsed = CreateClientSchema.safeParse(body)
  if (!parsed.success) {
    return NextResponse.json({ error: 'Dados inválidos', details: parsed.error.issues }, { status: 422 })
  }

  const {
    cep, logradouro, numero, complemento, bairro, cidade, uf,
    tipo_contrato, data_inicio, data_fim, status_contrato, observacoes_comerciais,
    ...clientData
  } = parsed.data

  const addressFields = { cep, logradouro, numero, complemento, bairro, cidade, uf }
  const hasAddress = Object.values(addressFields).some(Boolean)

  try {
    const client = await db.client.create({
      data: {
        ...clientData,
        ...(hasAddress ? {
          address: { create: addressFields },
        } : {}),
      },
      include: { address: true },
    })
    return NextResponse.json(client, { status: 201 })
  } catch (err) {
    if (err instanceof Prisma.PrismaClientKnownRequestError && err.code === 'P2002') {
      return NextResponse.json({ error: 'CPF/CNPJ já cadastrado' }, { status: 409 })
    }
    throw err
  }
}
