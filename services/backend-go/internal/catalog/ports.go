package catalog

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CustomerStore define as operações de persistência para clientes.
type CustomerStore interface {
	CreateCustomer(ctx context.Context, c *Customer) error
	GetCustomer(ctx context.Context, id uuid.UUID) (*Customer, error)
	GetCustomerByCPFCNPJ(ctx context.Context, cpfCnpj string) (*Customer, error)
	ListCustomers(ctx context.Context, filter CustomerFilter) ([]Customer, string, error)
	UpdateCustomer(ctx context.Context, id uuid.UUID, patch CustomerPatch) error
	ArchiveCustomer(ctx context.Context, id uuid.UUID) error
}

// ConsumerUnitStore define as operações de persistência para UCs.
type ConsumerUnitStore interface {
	CreateUnit(ctx context.Context, u *ConsumerUnit) error
	GetUnit(ctx context.Context, id uuid.UUID) (*ConsumerUnit, error)
	GetUnitByCode(ctx context.Context, code string) (*ConsumerUnit, error)
	ListUnits(ctx context.Context, filter UnitFilter) ([]ConsumerUnit, string, error)
	ListUnitsByCustomer(ctx context.Context, customerID uuid.UUID) ([]ConsumerUnit, error)
	LinkUnitToCustomer(ctx context.Context, unitID, customerID uuid.UUID) error
}

// ContractStore define as operações de persistência para contratos.
type ContractStore interface {
	CreateContract(ctx context.Context, c *Contract) error
	GetContract(ctx context.Context, id uuid.UUID) (*Contract, error)
	GetActiveContract(ctx context.Context, unitID uuid.UUID) (*Contract, error)
	ListContracts(ctx context.Context, filter ContractFilter) ([]Contract, error)
	CloseContract(ctx context.Context, unitID uuid.UUID, closeDate interface{}) error
}

// UserStore define as operações de persistência para usuários.
type UserStore interface {
	CreateUser(ctx context.Context, u *User) error
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	ListUsers(ctx context.Context) ([]User, error)
}

// Store agrupa todas as interfaces de persistência.
type Store interface {
	CustomerStore
	ConsumerUnitStore
	ContractStore
	UserStore
	WithTx(ctx context.Context, fn func(pgx.Tx) error) error
}

// pgxStore implementa Store usando pgxpool.
type pgxStore struct {
	pool *pgxpool.Pool
}

// NewStore cria uma nova instância de Store.
func NewStore(pool *pgxpool.Pool) Store {
	return &pgxStore{pool: pool}
}

func (s *pgxStore) WithTx(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
