import { describe, it, expect } from 'vitest'
import {
  CreateClientSchema,
  CreateCredentialSchema,
  ClientListQuerySchema,
} from './clientes'

describe('CreateClientSchema', () => {
  const baseClient = {
    tipo_pessoa: 'PF' as const,
    nome_razao: 'João Silva',
    cpf_cnpj: '12345678901',
    tipo_cliente: 'residencial' as const,
  }

  it('validates a valid PF client with 11-digit cpf_cnpj', () => {
    const result = CreateClientSchema.safeParse(baseClient)
    expect(result.success).toBe(true)
  })

  it('rejects PF client with 14-digit cpf_cnpj (CNPJ length)', () => {
    const result = CreateClientSchema.safeParse({
      ...baseClient,
      cpf_cnpj: '12345678000195',
    })
    expect(result.success).toBe(false)
    expect(result.error?.issues[0].path).toContain('cpf_cnpj')
  })

  it('validates a valid PJ client with 14-digit cpf_cnpj', () => {
    const result = CreateClientSchema.safeParse({
      ...baseClient,
      tipo_pessoa: 'PJ',
      cpf_cnpj: '12345678000195',
    })
    expect(result.success).toBe(true)
  })
})

describe('CreateCredentialSchema', () => {
  it('rejects when senha is empty string', () => {
    const result = CreateCredentialSchema.safeParse({
      label: 'neo-test',
      documento: '12345678901',
      senha: '',
      uf: 'BA',
      tipo_acesso: 'normal',
    })
    expect(result.success).toBe(false)
    expect(result.error?.issues[0].path).toContain('senha')
  })
})

describe('ClientListQuerySchema', () => {
  it('applies default values when no params provided', () => {
    const result = ClientListQuerySchema.parse({})
    expect(result.page).toBe(1)
    expect(result.pageSize).toBe(20)
    expect(result.archived).toBe(false)
    expect(result.order).toBe('desc')
    expect(result.orderBy).toBe('created_at')
  })
})
