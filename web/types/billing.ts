import { z } from 'zod'

// ─── Cycle schemas ────────────────────────────────────────────────────────────

export const CycleStatusEnum = z.enum([
  'open',
  'syncing',
  'processing',
  'review',
  'approved',
  'closed',
])

export const CreateCycleSchema = z.object({
  year: z.number().int().min(2000).max(2100),
  month: z.number().int().min(1).max(12),
  include_all_active: z.boolean().default(true),
  created_by: z.string().optional(),
})

export type CreateCycleInput = z.infer<typeof CreateCycleSchema>

export const CycleListQuerySchema = z.object({
  year: z.coerce.number().int().optional(),
  status: CycleStatusEnum.optional(),
  limit: z.coerce.number().int().min(1).max(100).default(50),
  offset: z.coerce.number().int().min(0).default(0),
})

export type CycleListQuery = z.infer<typeof CycleListQuerySchema>

// ─── Bulk action schema ───────────────────────────────────────────────────────

export const BulkActionSchema = z.object({
  action: z.enum(['sync', 'recalculate', 'generate_pdf', 'approve']),
  uc_codes: z.array(z.string()).optional(),
  force_all: z.boolean().default(false),
  created_by: z.string().optional(),
})

export type BulkActionInput = z.infer<typeof BulkActionSchema>

// ─── Close cycle schema ───────────────────────────────────────────────────────

export const CloseCycleSchema = z.object({
  closed_by: z.string().optional(),
})

export type CloseCycleInput = z.infer<typeof CloseCycleSchema>

// ─── Contract schemas ─────────────────────────────────────────────────────────

export const IPModeEnum = z.enum(['fixed', 'percent'])

export const CreateContractSchema = z.object({
  customer_id: z.string().uuid(),
  consumer_unit_id: z.string().uuid(),
  vigencia_inicio: z.string().regex(/^\d{4}-\d{2}-\d{2}$/, 'Formato YYYY-MM-DD'),
  fator_repasse_energia: z.string(),
  valor_ip_com_desconto: z.string().optional(),
  ip_faturamento_mode: IPModeEnum,
  ip_faturamento_valor: z.string().optional(),
  ip_faturamento_percent: z.string().optional(),
  bandeira_com_desconto: z.boolean().default(false),
  custo_disponibilidade_sempre_cobrado: z.boolean().default(false),
  notes: z.string().optional(),
  created_by: z.string().optional(),
})

export type CreateContractInput = z.infer<typeof CreateContractSchema>

// ─── Go API response types ────────────────────────────────────────────────────

export interface GoBillingCycle {
  id: string
  year: number
  month: number
  reference_date: string
  status: string
  created_at: string
  created_by?: string
  closed_at?: string
  closed_by?: string
  total_ucs?: number
  synced_count?: number
  calculated_count?: number
  approved_count?: number
}

export interface GoCycleListResponse {
  items: GoBillingCycle[]
  count: number
  limit: number
  offset: number
}

export interface GoCycleRow {
  consumer_unit_id: string
  uc_code: string
  customer_name: string
  sync_status: string
  calculation_status?: string
  numero_fatura?: string
  mes_referencia?: string
  valor_azi_sem_desconto?: string
  valor_azi_com_desconto?: string
  economia_rs?: string
  economia_pct?: string
  pdf_generated: boolean
  needs_review_reasons?: string[]
}

export interface GoCycleRowsResponse {
  items: GoCycleRow[]
  count: number
  cycle_id: string
}

export interface GoBulkActionResponse {
  jobs_created: number
  jobs_skipped: number
  skipped_reasons?: string[]
}

export interface GoContract {
  id: string
  customer_id: string
  consumer_unit_id: string
  vigencia_inicio: string
  vigencia_fim?: string
  fator_repasse_energia: string
  valor_ip_com_desconto: string
  ip_faturamento_mode: string
  ip_faturamento_valor: string
  ip_faturamento_percent: string
  bandeira_com_desconto: boolean
  custo_disponibilidade_sempre_cobrado: boolean
  status: string
  notes?: string
  created_at: string
  updated_at: string
}

export interface GoBillingCalculation {
  id: string
  utility_invoice_ref_id: string
  billing_cycle_id: string
  consumer_unit_id: string
  contract_id: string
  total_sem_desconto: string
  total_com_desconto: string
  economia_rs: string
  economia_pct: string
  status: string
  version: number
  calculated_at: string
  approved_at?: string
  approved_by?: string
  contract_snapshot?: unknown
  inputs_snapshot?: unknown
  result_snapshot?: unknown
  needs_review_reasons?: string[]
}
