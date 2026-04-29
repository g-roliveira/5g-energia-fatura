// Package repo implements data access for the billing domain. All SQL
// lives here — no other billing subpackage talks to the database directly.
// Services consume repo via interfaces so they can be unit-tested with
// mocks.
package repo

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// -------------------------------------------------------------------
// CONTRACT
// -------------------------------------------------------------------

type IPMode string

const (
	IPModeFixed   IPMode = "fixed"
	IPModePercent IPMode = "percent"
)

type ContractStatus string

const (
	ContractStatusDraft  ContractStatus = "draft"
	ContractStatusActive ContractStatus = "active"
	ContractStatusEnded  ContractStatus = "ended"
)

// Contract is a row in public.contract. Versioned by (vigencia_inicio,
// vigencia_fim). Never UPDATE the business fields — always INSERT a new
// row and close the previous (set vigencia_fim on the old one).
type Contract struct {
	ID                                uuid.UUID       `db:"id"`
	CustomerID                        uuid.UUID       `db:"customer_id"`
	ConsumerUnitID                    uuid.UUID       `db:"consumer_unit_id"`
	VigenciaInicio                    time.Time       `db:"vigencia_inicio"`
	VigenciaFim                       *time.Time      `db:"vigencia_fim"`
	DescontoPercentual                decimal.Decimal `db:"desconto_percentual"`
	IPFaturamentoMode                 IPMode          `db:"ip_faturamento_mode"`
	IPFaturamentoValor                decimal.Decimal `db:"ip_faturamento_valor"`
	IPFaturamentoPercent              decimal.Decimal `db:"ip_faturamento_percent"`
	BandeiraComDesconto               bool            `db:"bandeira_com_desconto"`
	CustoDisponibilidadeSempreCobrado bool            `db:"custo_disponibilidade_sempre_cobrado"`
	ConsumoMinimoKWh                  decimal.Decimal `db:"consumo_minimo_kwh"`
	Notes                             *string         `db:"notes"`
	Status                            ContractStatus  `db:"status"`
	CreatedAt                         time.Time       `db:"created_at"`
	CreatedBy                         *uuid.UUID      `db:"created_by"`
	UpdatedAt                         time.Time       `db:"updated_at"`
}

// IsActive returns true if this contract is currently in force
// (no vigencia_fim and status=active).
func (c *Contract) IsActive() bool {
	return c.VigenciaFim == nil && c.Status == ContractStatusActive
}

// -------------------------------------------------------------------
// BILLING CYCLE
// -------------------------------------------------------------------

type CycleStatus string

const (
	CycleStatusOpen       CycleStatus = "open"
	CycleStatusSyncing    CycleStatus = "syncing"
	CycleStatusProcessing CycleStatus = "processing"
	CycleStatusReview     CycleStatus = "review"
	CycleStatusApproved   CycleStatus = "approved"
	CycleStatusClosed     CycleStatus = "closed"
)

type BillingCycle struct {
	ID            uuid.UUID   `db:"id"`
	Year          int16       `db:"year"`
	Month         int16       `db:"month"`
	ReferenceDate time.Time   `db:"reference_date"`
	Status        CycleStatus `db:"status"`
	CreatedAt     time.Time   `db:"created_at"`
	CreatedBy     *uuid.UUID  `db:"created_by"`
	ClosedAt      *time.Time  `db:"closed_at"`
	ClosedBy      *uuid.UUID  `db:"closed_by"`
}

// -------------------------------------------------------------------
// UTILITY INVOICE REF (ponteiro pra fatura no SQLite do backend-go)
// -------------------------------------------------------------------

type UtilityInvoiceRef struct {
	ID                    uuid.UUID        `db:"id"`
	ConsumerUnitID        uuid.UUID        `db:"consumer_unit_id"`
	BillingCycleID        uuid.UUID        `db:"billing_cycle_id"`
	SyncInvoiceID         string           `db:"sync_invoice_id"`
	SyncRunID             *string          `db:"sync_run_id"`
	NumeroFatura          *string          `db:"numero_fatura"`
	MesReferencia         *string          `db:"mes_referencia"`
	ValorTotalCoelba      *decimal.Decimal `db:"valor_total_coelba"`
	StatusFatura          *string          `db:"status_fatura"`
	DataEmissao           *time.Time       `db:"data_emissao"`
	DataVencimento        *time.Time       `db:"data_vencimento"`
	DataInicioPeriodo     *time.Time       `db:"data_inicio_periodo"`
	DataFimPeriodo        *time.Time       `db:"data_fim_periodo"`
	CompletenessStatus    *string          `db:"completeness_status"`
	CompletenessMissing   []string         `db:"completeness_missing"`
	ExtractorStatus       *string          `db:"extractor_status"`
	ExtractorConfidence   *decimal.Decimal `db:"extractor_confidence"`
	BillingRecordSnapshot []byte           `db:"billing_record_snapshot"` // JSONB raw
	SyncedAt              *time.Time       `db:"synced_at"`
	CreatedAt             time.Time        `db:"created_at"`
	UpdatedAt             time.Time        `db:"updated_at"`
}

// -------------------------------------------------------------------
// BILLING CALCULATION
// -------------------------------------------------------------------

type CalcStatus string

const (
	CalcStatusDraft       CalcStatus = "draft"
	CalcStatusNeedsReview CalcStatus = "needs_review"
	CalcStatusApproved    CalcStatus = "approved"
	CalcStatusSuperseded  CalcStatus = "superseded"
)

type BillingCalculation struct {
	ID                   uuid.UUID       `db:"id"`
	UtilityInvoiceRefID  uuid.UUID       `db:"utility_invoice_ref_id"`
	BillingCycleID       uuid.UUID       `db:"billing_cycle_id"`
	ConsumerUnitID       uuid.UUID       `db:"consumer_unit_id"`
	ContractID           uuid.UUID       `db:"contract_id"`
	ContractSnapshotJSON []byte          `db:"contract_snapshot_json"`
	InputsSnapshotJSON   []byte          `db:"inputs_snapshot_json"`
	ResultSnapshotJSON   []byte          `db:"result_snapshot_json"`
	TotalSemDesconto     decimal.Decimal `db:"total_sem_desconto"`
	TotalComDesconto     decimal.Decimal `db:"total_com_desconto"`
	EconomiaRS           decimal.Decimal `db:"economia_rs"`
	EconomiaPct          decimal.Decimal `db:"economia_pct"`
	Status               CalcStatus      `db:"status"`
	NeedsReviewReasons   []string        `db:"needs_review_reasons"`
	Version              int             `db:"version"`
	CalculatedAt         time.Time       `db:"calculated_at"`
	ApprovedAt           *time.Time      `db:"approved_at"`
	ApprovedBy           *uuid.UUID      `db:"approved_by"`
}
