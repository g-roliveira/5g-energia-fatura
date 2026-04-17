import { PrismaPg } from '@prisma/adapter-pg'
import { PrismaClient } from '@prisma/client'
import { Pool } from 'pg'

const connectionString = process.env.DATABASE_URL

if (!connectionString) {
  throw new Error('DATABASE_URL is not set')
}

const globalForPrisma = globalThis as unknown as {
  pool: Pool | undefined
  prisma: PrismaClient | undefined
}

const pool = globalForPrisma.pool ?? new Pool({ connectionString })
const adapter = new PrismaPg(pool)

export const db = globalForPrisma.prisma ?? new PrismaClient({ adapter })

if (process.env.NODE_ENV !== 'production') {
  globalForPrisma.pool = pool
  globalForPrisma.prisma = db
}
