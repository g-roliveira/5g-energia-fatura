package catalog

import (
	"time"

	"github.com/google/uuid"
)

// Customer (corresponde ao model Client no Prisma)
type Customer struct {
	ID           uuid.UUID
	TipoPessoa   string // PF | PJ
	NomeRazao    string
	NomeFantasia *string
	CPFCNPJ      string
	Email        *string
	Telefone     *string
	Status       string // ativo | inativo | prospecto | archived
	TipoCliente  string // residencial | condominio | empresa | imobiliaria | outro
	Notes        *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ArchivedAt   *time.Time
}

// Address (corresponde ao model ClientAddress no Prisma)
type Address struct {
	ID          uuid.UUID
	CustomerID  uuid.UUID
	CEP         *string
	Logradouro  *string
	Numero      *string
	Complemento *string
	Bairro      *string
	Cidade      *string
	UF          *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ConsumerUnit (corresponde ao model ConsumerUnit no Prisma)
type ConsumerUnit struct {
	ID            uuid.UUID
	CustomerID    uuid.UUID
	UCCode        string
	Distribuidora *string
	Apelido       *string
	ClasseConsumo *string
	Endereco      *string
	Cidade        *string
	UF            *string
	Ativa         bool
	CredentialID  *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Contract (versionado — nunca editado, só criado e fechado)
type Contract struct {
	ID                                uuid.UUID
	CustomerID                        uuid.UUID
	ConsumerUnitID                    uuid.UUID
	VigenciaInicio                    time.Time
	VigenciaFim                       *time.Time
	FatorRepasseEnergia               string // decimal como string para precisão
	ValorIPComDesconto                string // decimal; valor contratual da IP no cenário COM repasse
	IPFaturamentoMode                 string // fixed | percent
	IPFaturamentoValor                string
	IPFaturamentoPercent              string
	BandeiraComDesconto               bool
	CustoDisponibilidadeSempreCobrado bool
	ConsumoMinimoKWh                  string
	Notes                             *string
	Status                            string // draft | active | ended
	CreatedAt                         time.Time
	CreatedBy                         *uuid.UUID
	UpdatedAt                         time.Time
}

// User (corresponde ao model AppUser no Prisma — quando implementado)
type User struct {
	ID           uuid.UUID
	Email        string
	Name         string
	PasswordHash string
	Role         string // admin | operator | reviewer
	Active       bool
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CustomerFilter para listagem
type CustomerFilter struct {
	Status *string
	Query  *string
	Limit  int
	Cursor *string
}

// ContractFilter para listagem
type ContractFilter struct {
	ConsumerUnitID *uuid.UUID
	ActiveOnly     bool
}
