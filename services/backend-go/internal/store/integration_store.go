package store

// IntegrationStore persists and reads integration-domain data used by backend-go
// endpoints and sync/session workflows.
type IntegrationStore interface {
	Close() error
	InsertCredential(rec CredentialRecord) (CredentialRecord, error)
	GetCredentialByID(id string) (CredentialRecord, error)
	InsertSession(rec SessionRecord) (SessionRecord, error)
	GetLatestSessionByCredentialID(credentialID string) (SessionRecord, error)
	ListConsumerUnits(limit int, status string) ([]ConsumerUnitView, error)
	GetConsumerUnitByUC(uc string) (*ConsumerUnitDetailsView, error)
	ListInvoicesByUC(uc string, limit int, status string) ([]InvoiceView, error)
	GetInvoiceByID(id string) (*InvoiceView, error)
	GetLatestInvoiceByUC(uc string) (*InvoiceView, error)
	GetSyncRunByID(id string) (*SyncRunView, error)
	PersistSyncResult(in PersistSyncInput) (PersistSyncResult, error)
}
