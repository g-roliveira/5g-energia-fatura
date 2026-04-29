package integration

import (
	"time"

	"github.com/google/uuid"
)

// Credential representa uma credencial de acesso à Neoenergia.
type Credential struct {
	ID              uuid.UUID
	Label           string
	DocumentoCipher string
	DocumentoNonce  string
	SenhaCipher     string
	SenhaNonce      string
	UF              string
	TipoAcesso      string
	KeyVersion      string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Session representa uma sessão ativa do Playwright.
type Session struct {
	ID                uuid.UUID
	CredentialID      uuid.UUID
	BearerTokenCipher string
	BearerTokenNonce  string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// ConsumerUnit representa uma UC descoberta via scraping.
type ConsumerUnit struct {
	UC           string
	CredentialID *uuid.UUID
	Status       *string
	NomeCliente  *string
	Instalacao   *string
	Contrato     *string
	GrupoTensao  *string
	Endereco     map[string]any
	Imovel       map[string]any
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// RawInvoice representa uma fatura bruta baixada.
type RawInvoice struct {
	ID                   uuid.UUID
	UC                   string
	NumeroFatura         string
	MesReferencia        string
	StatusFatura         *string
	ValorTotal           *string
	CodigoBarras         *string
	DataEmissao          *string
	DataVencimento       *string
	DataPagamento        *string
	DataInicioPeriodo    *string
	DataFimPeriodo       *string
	CompletenessStatus   *string
	CompletenessMissing  []string
	BillingRecordJSON    map[string]any
	DocumentRecordJSON   map[string]any
	PDFBytes             []byte
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// SyncRun representa uma execução de sync.
type SyncRun struct {
	ID                uuid.UUID
	CredentialID      *uuid.UUID
	Documento         string
	UC                string
	Status            string
	Step              *string
	ErrorMessage      *string
	ErrorContext      map[string]any
	RawResponseJSON   map[string]any
	StartedAt         *time.Time
	FinishedAt        *time.Time
	CreatedAt         time.Time
}

// Job representa um item na fila de jobs.
type Job struct {
	ID           uuid.UUID
	JobType      string
	Status       string
	Payload      map[string]any
	Result       map[string]any
	ErrorMessage *string
	RetryCount   int
	MaxRetries   int
	ClaimedBy    *string
	CreatedAt    time.Time
	ClaimedAt    *time.Time
	CompletedAt  *time.Time
}
