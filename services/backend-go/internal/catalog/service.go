package catalog

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Service expõe as operações do domínio catalog.
type Service struct {
	store Store
}

// NewService cria um novo Service.
func NewService(store Store) *Service {
	return &Service{store: store}
}

// --- Customer ---

func (s *Service) CreateCustomer(ctx context.Context, input CustomerInput) (*Customer, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	existing, _ := s.store.GetCustomerByCPFCNPJ(ctx, input.CPFCNPJ)
	if existing != nil {
		return nil, fmt.Errorf("customer with CPF/CNPJ %s already exists", input.CPFCNPJ)
	}

	now := time.Now()
	c := &Customer{
		ID:           uuid.New(),
		TipoPessoa:   input.TipoPessoa,
		NomeRazao:    input.NomeRazao,
		NomeFantasia: input.NomeFantasia,
		CPFCNPJ:      input.CPFCNPJ,
		Email:        input.Email,
		Telefone:     input.Telefone,
		Status:       "prospect",
		TipoCliente:  input.TipoCliente,
		Notes:        input.Notes,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.store.CreateCustomer(ctx, c); err != nil {
		return nil, fmt.Errorf("create customer: %w", err)
	}
	return c, nil
}

func (s *Service) GetCustomer(ctx context.Context, id uuid.UUID) (*Customer, error) {
	c, err := s.store.GetCustomer(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get customer: %w", err)
	}
	return c, nil
}

func (s *Service) ListCustomers(ctx context.Context, filter CustomerFilter) ([]Customer, string, error) {
	return s.store.ListCustomers(ctx, filter)
}

func (s *Service) UpdateCustomer(ctx context.Context, id uuid.UUID, patch CustomerPatch) error {
	if err := s.store.UpdateCustomer(ctx, id, patch); err != nil {
		return fmt.Errorf("update customer: %w", err)
	}
	return nil
}

func (s *Service) ArchiveCustomer(ctx context.Context, id uuid.UUID) error {
	if err := s.store.ArchiveCustomer(ctx, id); err != nil {
		return fmt.Errorf("archive customer: %w", err)
	}
	return nil
}

// --- Consumer Unit ---

func (s *Service) CreateUnit(ctx context.Context, input UnitInput) (*ConsumerUnit, error) {
	if input.UCCode == "" {
		return nil, fmt.Errorf("uc_code is required")
	}

	existing, _ := s.store.GetUnitByCode(ctx, input.UCCode)
	if existing != nil {
		return nil, fmt.Errorf("unit with code %s already exists", input.UCCode)
	}

	now := time.Now()
	u := &ConsumerUnit{
		ID:            uuid.New(),
		CustomerID:    input.CustomerID,
		UCCode:        input.UCCode,
		Distribuidora: input.Distribuidora,
		Apelido:       input.Apelido,
		ClasseConsumo: input.ClasseConsumo,
		Endereco:      input.Endereco,
		Cidade:        input.Cidade,
		UF:            input.UF,
		Ativa:         true,
		CredentialID:  input.CredentialID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.store.CreateUnit(ctx, u); err != nil {
		return nil, fmt.Errorf("create unit: %w", err)
	}
	return u, nil
}

func (s *Service) GetUnit(ctx context.Context, id uuid.UUID) (*ConsumerUnit, error) {
	return s.store.GetUnit(ctx, id)
}

func (s *Service) ListUnits(ctx context.Context, filter UnitFilter) ([]ConsumerUnit, string, error) {
	return s.store.ListUnits(ctx, filter)
}

func (s *Service) ListUnitsByCustomer(ctx context.Context, customerID uuid.UUID) ([]ConsumerUnit, error) {
	return s.store.ListUnitsByCustomer(ctx, customerID)
}

func (s *Service) LinkUnitToCustomer(ctx context.Context, unitID, customerID uuid.UUID) error {
	return s.store.LinkUnitToCustomer(ctx, unitID, customerID)
}

// --- Contract ---

// CreateContract fecha o contrato vigente (se existir) e cria um novo.
// Essa é a operação crítica de versionamento — nunca UPDATE, sempre INSERT.
func (s *Service) CreateContract(ctx context.Context, input ContractInput) (*Contract, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	vigenciaInicio, err := time.Parse("2006-01-02", input.VigenciaInicio)
	if err != nil {
		return nil, fmt.Errorf("vigencia_inicio inválida: %w", err)
	}

	now := time.Now()
	contract := &Contract{
		ID:                                uuid.New(),
		CustomerID:                        input.CustomerID,
		ConsumerUnitID:                    input.ConsumerUnitID,
		VigenciaInicio:                    vigenciaInicio,
		DescontoPercentual:                input.DescontoPercentual,
		IPFaturamentoMode:                 input.IPFaturamentoMode,
		IPFaturamentoValor:                input.IPFaturamentoValor,
		IPFaturamentoPercent:              input.IPFaturamentoPercent,
		BandeiraComDesconto:               input.BandeiraComDesconto,
		CustoDisponibilidadeSempreCobrado: input.CustoDisponibilidadeSempreCobrado,
		ConsumoMinimoKWh:                  defaultString(input.ConsumoMinimoKWh, "30.0"),
		Notes:                             input.Notes,
		Status:                            "active",
		CreatedAt:                         now,
		UpdatedAt:                         now,
	}
	if input.CreatedBy != nil {
		contract.CreatedBy = input.CreatedBy
	}

	// Fecha contrato anterior e insere novo numa transação
	err = s.store.WithTx(ctx, func(tx pgx.Tx) error {
		// Fechar contrato vigente
		closeDate := vigenciaInicio.AddDate(0, 0, -1)
		if _, err := tx.Exec(ctx, `
			UPDATE public.contract
			SET vigencia_fim = $1, status = 'ended', updated_at = NOW()
			WHERE consumer_unit_id = $2 AND vigencia_fim IS NULL AND status = 'active'
		`, closeDate, input.ConsumerUnitID); err != nil {
			return fmt.Errorf("close existing contract: %w", err)
		}

		// Inserir novo
		if _, err := tx.Exec(ctx, `
			INSERT INTO public.contract (id, customer_id, consumer_unit_id, vigencia_inicio, desconto_percentual, ip_faturamento_mode, ip_faturamento_valor, ip_faturamento_percent, bandeira_com_desconto, custo_disponibilidade_sempre_cobrado, consumo_minimo_kwh, notes, status, created_at, created_by, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		`,
			contract.ID, contract.CustomerID, contract.ConsumerUnitID,
			contract.VigenciaInicio, contract.DescontoPercentual,
			contract.IPFaturamentoMode, contract.IPFaturamentoValor,
			contract.IPFaturamentoPercent, contract.BandeiraComDesconto,
			contract.CustoDisponibilidadeSempreCobrado, contract.ConsumoMinimoKWh, contract.Notes,
			contract.Status, contract.CreatedAt, contract.CreatedBy, contract.UpdatedAt,
		); err != nil {
			return fmt.Errorf("insert contract: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return contract, nil
}

func (s *Service) GetContract(ctx context.Context, id uuid.UUID) (*Contract, error) {
	return s.store.GetContract(ctx, id)
}

func (s *Service) GetActiveContract(ctx context.Context, unitID uuid.UUID) (*Contract, error) {
	return s.store.GetActiveContract(ctx, unitID)
}

func (s *Service) ListContracts(ctx context.Context, filter ContractFilter) ([]Contract, error) {
	return s.store.ListContracts(ctx, filter)
}

// --- Input types ---

type CustomerInput struct {
	TipoPessoa   string  `json:"tipo_pessoa"`
	NomeRazao    string  `json:"nome_razao"`
	NomeFantasia *string `json:"nome_fantasia,omitempty"`
	CPFCNPJ      string  `json:"cpf_cnpj"`
	Email        *string `json:"email,omitempty"`
	Telefone     *string `json:"telefone,omitempty"`
	TipoCliente  string  `json:"tipo_cliente"`
	Notes        *string `json:"notes,omitempty"`
}

func (i CustomerInput) Validate() error {
	if i.NomeRazao == "" {
		return fmt.Errorf("nome_razao is required")
	}
	if i.CPFCNPJ == "" {
		return fmt.Errorf("cpf_cnpj is required")
	}
	if i.TipoPessoa != "PF" && i.TipoPessoa != "PJ" {
		return fmt.Errorf("tipo_pessoa must be PF or PJ")
	}
	return nil
}

type UnitInput struct {
	CustomerID    uuid.UUID `json:"customer_id"`
	UCCode        string    `json:"uc_code"`
	Distribuidora *string   `json:"distribuidora,omitempty"`
	Apelido       *string   `json:"apelido,omitempty"`
	ClasseConsumo *string   `json:"classe_consumo,omitempty"`
	Endereco      *string   `json:"endereco,omitempty"`
	Cidade        *string   `json:"cidade,omitempty"`
	UF            *string   `json:"uf,omitempty"`
	CredentialID  *string   `json:"credential_id,omitempty"`
}

type ContractInput struct {
	CustomerID                        uuid.UUID  `json:"customer_id"`
	ConsumerUnitID                    uuid.UUID  `json:"consumer_unit_id"`
	VigenciaInicio                    string     `json:"vigencia_inicio"` // YYYY-MM-DD
	DescontoPercentual                string     `json:"desconto_percentual"`
	IPFaturamentoMode                 string     `json:"ip_faturamento_mode"`
	IPFaturamentoValor                string     `json:"ip_faturamento_valor"`
	IPFaturamentoPercent              string     `json:"ip_faturamento_percent"`
	BandeiraComDesconto               bool       `json:"bandeira_com_desconto"`
	CustoDisponibilidadeSempreCobrado bool       `json:"custo_disponibilidade_sempre_cobrado"`
	ConsumoMinimoKWh                  string     `json:"consumo_minimo_kwh,omitempty"`
	Notes                             *string    `json:"notes,omitempty"`
	CreatedBy                         *uuid.UUID `json:"created_by,omitempty"`
}

func (i ContractInput) Validate() error {
	if i.ConsumerUnitID == uuid.Nil {
		return fmt.Errorf("consumer_unit_id is required")
	}
	if i.DescontoPercentual == "" {
		return fmt.Errorf("desconto_percentual is required")
	}
	if i.IPFaturamentoMode != "fixed" && i.IPFaturamentoMode != "percent" {
		return fmt.Errorf("ip_faturamento_mode must be fixed or percent")
	}
	if i.VigenciaInicio == "" {
		return fmt.Errorf("vigencia_inicio is required")
	}
	return nil
}

// defaultString retorna def se s for vazio, senão retorna s.
func defaultString(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
