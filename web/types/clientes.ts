import { z } from 'zod'

// ─── Enums ────────────────────────────────────────────────────────────────────

export const TipoPessoaEnum = z.enum(['PF', 'PJ'])
export const ClientStatusEnum = z.enum(['ativo', 'inativo', 'prospecto'])
export const TipoClienteEnum = z.enum([
  'residencial',
  'condominio',
  'empresa',
  'imobiliaria',
  'outro',
])

// ─── Client schemas ───────────────────────────────────────────────────────────

const ClientBaseSchema = z.object({
  tipo_pessoa: TipoPessoaEnum,
  nome_razao: z.string().min(1),
  nome_fantasia: z.string().optional(),
  cpf_cnpj: z.string().min(1),
  email: z.string().email().optional().or(z.literal('')),
  telefone: z.string().optional(),
  status: ClientStatusEnum.default('prospecto'),
  tipo_cliente: TipoClienteEnum,
  observacoes: z.string().optional(),
  // Address (optional nested)
  cep: z.string().optional(),
  logradouro: z.string().optional(),
  numero: z.string().optional(),
  complemento: z.string().optional(),
  bairro: z.string().optional(),
  cidade: z.string().optional(),
  uf: z.string().length(2).optional(),
  // Commercial (optional nested)
  tipo_contrato: z.string().optional(),
  data_inicio: z.string().optional(),
  data_fim: z.string().optional(),
  status_contrato: z.string().optional(),
  observacoes_comerciais: z.string().optional(),
})

export const CreateClientSchema = ClientBaseSchema.superRefine((data, ctx) => {
  if (data.tipo_pessoa === 'PF' && data.cpf_cnpj.replace(/\D/g, '').length !== 11) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: 'CPF deve ter 11 dígitos',
      path: ['cpf_cnpj'],
    })
  }
  if (data.tipo_pessoa === 'PJ' && data.cpf_cnpj.replace(/\D/g, '').length !== 14) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: 'CNPJ deve ter 14 dígitos',
      path: ['cpf_cnpj'],
    })
  }
})

export type CreateClientInput = z.infer<typeof CreateClientSchema>

export const UpdateClientSchema = ClientBaseSchema.partial()
export type UpdateClientInput = z.infer<typeof UpdateClientSchema>

// ─── UC schemas ───────────────────────────────────────────────────────────────

export const CreateUcSchema = z.object({
  uc_code: z.string().min(1),
  distribuidora: z.string().optional(),
  apelido: z.string().optional(),
  classe_consumo: z.string().optional(),
  endereco_unidade: z.string().optional(),
  cidade: z.string().optional(),
  uf: z.string().length(2).optional(),
  ativa: z.boolean().default(true),
  credential_id: z.string().optional(),
})

export type CreateUcInput = z.infer<typeof CreateUcSchema>

export const UpdateUcSchema = CreateUcSchema.omit({ uc_code: true }).partial()
export type UpdateUcInput = z.infer<typeof UpdateUcSchema>

// ─── Credential schemas ───────────────────────────────────────────────────────

export const CreateCredentialSchema = z.object({
  client_id: z.string().min(1).optional(),
  label: z.string().min(1),
  documento: z.string().min(1),
  senha: z.string().min(6, 'Senha deve ter pelo menos 6 caracteres'),
  uf: z.string().length(2),
  tipo_acesso: z.string().default('normal'),
})

export type CreateCredentialInput = z.infer<typeof CreateCredentialSchema>

// ─── Pagination ───────────────────────────────────────────────────────────────

export const ClientListQuerySchema = z.object({
  page: z.coerce.number().int().positive().default(1),
  pageSize: z.coerce.number().int().min(1).max(100).default(20),
  search: z.string().optional(),
  status: ClientStatusEnum.optional(),
  tipo_cliente: TipoClienteEnum.optional(),
  archived: z.coerce.boolean().default(false),
  orderBy: z.string().default('created_at'),
  order: z.enum(['asc', 'desc']).default('desc'),
})

export type ClientListQuery = z.infer<typeof ClientListQuerySchema>

// ─── Go API response types (TypeScript interfaces — no Zod validation needed) ─

export interface GoCredential {
  id: string
  label: string
  documento: string
  uf: string
  tipo_acesso: string
  created_at: string
}

export interface GoSession {
  id: string
  credential_id: string
  created_at: string
}

export interface GoConsumerUnit {
  uc: string
  status: string
  nome_cliente?: string
  latest_invoice?: {
    numero_fatura: string
    mes_referencia: string
    completeness_status: string
  }
  latest_sync_run?: {
    status: string
  }
}

export interface GoBillingRecord {
  numero_fatura?: string
  mes_referencia?: string
  valor_total?: string
  data_vencimento?: string
  completeness?: {
    status: 'complete' | 'partial' | 'failed'
  }
  [key: string]: unknown
}

export interface GoDocumentRecord {
  [key: string]: unknown
}

export interface GoInvoiceItem {
  descricao?: string
  quantidade?: number
  valor_unitario?: string
  valor_total?: string
  [key: string]: unknown
}

export interface GoInvoice {
  id: string
  billing_record: GoBillingRecord
  document_record?: GoDocumentRecord
  items?: GoInvoiceItem[]
  completeness_status: 'complete' | 'partial' | 'failed'
  updated_at: string
}

export interface GoSyncRun {
  id: string
  status: 'running' | 'succeeded' | 'partial' | 'failed'
  error_message?: string
  raw_response?: unknown
  created_at: string
}

export interface GoPersistence {
  sync_run_id: string
  invoice_id: string
  status: string
}
