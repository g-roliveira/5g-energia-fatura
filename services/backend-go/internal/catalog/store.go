package catalog

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// pgxStore já declarado em ports.go

var _ Store = (*pgxStore)(nil)

// --- CustomerStore ---

func (s *pgxStore) CreateCustomer(ctx context.Context, c *Customer) error {
	query := `
		INSERT INTO public.customer (id, tipo_pessoa, nome_razao, nome_fantasia, cpf_cnpj, email, phone, status, tipo_cliente, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := s.pool.Exec(ctx, query,
		c.ID, c.TipoPessoa, c.NomeRazao, c.NomeFantasia, c.CPFCNPJ,
		c.Email, c.Telefone, c.Status, c.TipoCliente, c.Notes,
		c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (s *pgxStore) GetCustomer(ctx context.Context, id uuid.UUID) (*Customer, error) {
	query := `
		SELECT id, tipo_pessoa, nome_razao, nome_fantasia, cpf_cnpj, email, phone, status, tipo_cliente, notes, created_at, updated_at, archived_at
		FROM public.customer WHERE id = $1
	`
	row := s.pool.QueryRow(ctx, query, id.String())
	return scanCustomer(row)
}

func (s *pgxStore) GetCustomerByCPFCNPJ(ctx context.Context, cpfCnpj string) (*Customer, error) {
	query := `
		SELECT id, tipo_pessoa, nome_razao, nome_fantasia, cpf_cnpj, email, phone, status, tipo_cliente, notes, created_at, updated_at, archived_at
		FROM public.customer WHERE cpf_cnpj = $1
	`
	row := s.pool.QueryRow(ctx, query, cpfCnpj)
	return scanCustomer(row)
}

func (s *pgxStore) ListCustomers(ctx context.Context, filter CustomerFilter) ([]Customer, string, error) {
	var args []interface{}
	var conds []string
	argNum := 1

	if filter.Status != nil && *filter.Status != "" {
		conds = append(conds, fmt.Sprintf("status = $%d", argNum))
		args = append(args, *filter.Status)
		argNum++
	}
	if filter.Query != nil && *filter.Query != "" {
		conds = append(conds, fmt.Sprintf("(nome_razao ILIKE $%d OR cpf_cnpj ILIKE $%d)", argNum, argNum))
		args = append(args, "%"+*filter.Query+"%")
		argNum++
	}
	if filter.Cursor != nil && *filter.Cursor != "" {
		conds = append(conds, fmt.Sprintf("id > $%d", argNum))
		args = append(args, *filter.Cursor)
		argNum++
	}

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	query := fmt.Sprintf(`
		SELECT id, tipo_pessoa, nome_razao, nome_fantasia, cpf_cnpj, email, phone, status, tipo_cliente, notes, created_at, updated_at, archived_at
		FROM public.customer
		%s
		ORDER BY id
		LIMIT $%d
	`, where, argNum)
	args = append(args, limit+1)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var customers []Customer
	var nextCursor string
	count := 0
	for rows.Next() {
		c, err := scanCustomerRows(rows)
		if err != nil {
			return nil, "", err
		}
		if count >= limit {
			nextCursor = c.ID.String()
			break
		}
		customers = append(customers, *c)
		count++
	}

	return customers, nextCursor, rows.Err()
}

type CustomerPatch struct {
	NomeRazao    *string `json:"nome_razao,omitempty"`
	NomeFantasia *string `json:"nome_fantasia,omitempty"`
	Email        *string `json:"email,omitempty"`
	Telefone     *string `json:"telefone,omitempty"`
	Status       *string `json:"status,omitempty"`
	Notes        *string `json:"notes,omitempty"`
}

func (s *pgxStore) UpdateCustomer(ctx context.Context, id uuid.UUID, patch CustomerPatch) error {
	var sets []string
	var args []interface{}
	argNum := 1

	if patch.NomeRazao != nil {
		sets = append(sets, fmt.Sprintf("nome_razao = $%d", argNum))
		args = append(args, *patch.NomeRazao)
		argNum++
	}
	if patch.NomeFantasia != nil {
		sets = append(sets, fmt.Sprintf("nome_fantasia = $%d", argNum))
		args = append(args, *patch.NomeFantasia)
		argNum++
	}
	if patch.Email != nil {
		sets = append(sets, fmt.Sprintf("email = $%d", argNum))
		args = append(args, *patch.Email)
		argNum++
	}
	if patch.Telefone != nil {
		sets = append(sets, fmt.Sprintf("phone = $%d", argNum))
		args = append(args, *patch.Telefone)
		argNum++
	}
	if patch.Status != nil {
		sets = append(sets, fmt.Sprintf("status = $%d", argNum))
		args = append(args, *patch.Status)
		argNum++
	}
	if patch.Notes != nil {
		sets = append(sets, fmt.Sprintf("notes = $%d", argNum))
		args = append(args, *patch.Notes)
		argNum++
	}

	if len(sets) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		UPDATE public.customer
		SET %s, updated_at = NOW()
		WHERE id = $%d
	`, strings.Join(sets, ", "), argNum)
	args = append(args, id.String())

	_, err := s.pool.Exec(ctx, query, args...)
	return err
}

func (s *pgxStore) ArchiveCustomer(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE public.customer SET status = 'archived', archived_at = NOW(), updated_at = NOW() WHERE id = $1`
	_, err := s.pool.Exec(ctx, query, id.String())
	return err
}

// --- ConsumerUnitStore ---

func (s *pgxStore) CreateUnit(ctx context.Context, u *ConsumerUnit) error {
	query := `
		INSERT INTO public.consumer_unit (id, customer_id, uc_code, distribuidora, apelido, classe_consumo, endereco_unidade, cidade, uf, ativa, sync_credential_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := s.pool.Exec(ctx, query,
		u.ID, u.CustomerID, u.UCCode, u.Distribuidora, u.Apelido,
		u.ClasseConsumo, u.Endereco, u.Cidade, u.UF, u.Ativa, u.CredentialID,
		u.CreatedAt, u.UpdatedAt,
	)
	return err
}

func (s *pgxStore) GetUnit(ctx context.Context, id uuid.UUID) (*ConsumerUnit, error) {
	query := `
		SELECT id, customer_id, uc_code, distribuidora, apelido, classe_consumo, endereco_unidade, cidade, uf, ativa, sync_credential_id, created_at, updated_at
		FROM public.consumer_unit WHERE id = $1
	`
	row := s.pool.QueryRow(ctx, query, id.String())
	return scanConsumerUnit(row)
}

func (s *pgxStore) GetUnitByCode(ctx context.Context, code string) (*ConsumerUnit, error) {
	query := `
		SELECT id, customer_id, uc_code, distribuidora, apelido, classe_consumo, endereco_unidade, cidade, uf, ativa, sync_credential_id, created_at, updated_at
		FROM public.consumer_unit WHERE uc_code = $1
	`
	row := s.pool.QueryRow(ctx, query, code)
	return scanConsumerUnit(row)
}

type UnitFilter struct {
	CustomerID *uuid.UUID
	ActiveOnly bool
	Limit      int
	Cursor     *string
}

func (s *pgxStore) ListUnits(ctx context.Context, filter UnitFilter) ([]ConsumerUnit, string, error) {
	var args []interface{}
	var conds []string
	argNum := 1

	if filter.CustomerID != nil {
		conds = append(conds, fmt.Sprintf("customer_id = $%d", argNum))
		args = append(args, filter.CustomerID.String())
		argNum++
	}
	if filter.ActiveOnly {
		conds = append(conds, "ativa = true")
	}
	if filter.Cursor != nil && *filter.Cursor != "" {
		conds = append(conds, fmt.Sprintf("id > $%d", argNum))
		args = append(args, *filter.Cursor)
		argNum++
	}

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	query := fmt.Sprintf(`
		SELECT id, customer_id, uc_code, distribuidora, apelido, classe_consumo, endereco_unidade, cidade, uf, ativa, sync_credential_id, created_at, updated_at
		FROM public.consumer_unit
		%s
		ORDER BY id
		LIMIT $%d
	`, where, argNum)
	args = append(args, limit+1)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var units []ConsumerUnit
	var nextCursor string
	count := 0
	for rows.Next() {
		u, err := scanConsumerUnitRows(rows)
		if err != nil {
			return nil, "", err
		}
		if count >= limit {
			nextCursor = u.ID.String()
			break
		}
		units = append(units, *u)
		count++
	}

	return units, nextCursor, rows.Err()
}

func (s *pgxStore) ListUnitsByCustomer(ctx context.Context, customerID uuid.UUID) ([]ConsumerUnit, error) {
	query := `
		SELECT id, customer_id, uc_code, distribuidora, apelido, classe_consumo, endereco_unidade, cidade, uf, ativa, sync_credential_id, created_at, updated_at
		FROM public.consumer_unit WHERE customer_id = $1 ORDER BY created_at DESC
	`
	rows, err := s.pool.Query(ctx, query, customerID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []ConsumerUnit
	for rows.Next() {
		u, err := scanConsumerUnitRows(rows)
		if err != nil {
			return nil, err
		}
		units = append(units, *u)
	}
	return units, rows.Err()
}

func (s *pgxStore) LinkUnitToCustomer(ctx context.Context, unitID, customerID uuid.UUID) error {
	query := `UPDATE public.consumer_unit SET customer_id = $1, updated_at = NOW() WHERE id = $2`
	_, err := s.pool.Exec(ctx, query, customerID.String(), unitID.String())
	return err
}

// --- ContractStore ---

func (s *pgxStore) CreateContract(ctx context.Context, c *Contract) error {
	query := `
		INSERT INTO public.contract (id, customer_id, consumer_unit_id, vigencia_inicio, vigencia_fim, fator_repasse_energia, valor_ip_com_desconto, ip_faturamento_mode, ip_faturamento_valor, ip_faturamento_percent, bandeira_com_desconto, custo_disponibilidade_sempre_cobrado, consumo_minimo_kwh, notes, status, created_at, created_by, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`
	_, err := s.pool.Exec(ctx, query,
		c.ID, c.CustomerID, c.ConsumerUnitID, c.VigenciaInicio, c.VigenciaFim,
		c.FatorRepasseEnergia, c.ValorIPComDesconto, c.IPFaturamentoMode, c.IPFaturamentoValor,
		c.IPFaturamentoPercent, c.BandeiraComDesconto, c.CustoDisponibilidadeSempreCobrado,
		c.ConsumoMinimoKWh,
		c.Notes, c.Status, c.CreatedAt, c.CreatedBy, c.UpdatedAt,
	)
	return err
}

func (s *pgxStore) GetContract(ctx context.Context, id uuid.UUID) (*Contract, error) {
	query := `
		SELECT id, customer_id, consumer_unit_id, vigencia_inicio, vigencia_fim, fator_repasse_energia, valor_ip_com_desconto, ip_faturamento_mode, ip_faturamento_valor, ip_faturamento_percent, bandeira_com_desconto, custo_disponibilidade_sempre_cobrado, consumo_minimo_kwh, notes, status, created_at, created_by, updated_at
		FROM public.contract WHERE id = $1
	`
	row := s.pool.QueryRow(ctx, query, id)
	return scanContract(row)
}

func (s *pgxStore) GetActiveContract(ctx context.Context, unitID uuid.UUID) (*Contract, error) {
	query := `
		SELECT id, customer_id, consumer_unit_id, vigencia_inicio, vigencia_fim, fator_repasse_energia, valor_ip_com_desconto, ip_faturamento_mode, ip_faturamento_valor, ip_faturamento_percent, bandeira_com_desconto, custo_disponibilidade_sempre_cobrado, consumo_minimo_kwh, notes, status, created_at, created_by, updated_at
		FROM public.contract
		WHERE consumer_unit_id = $1 AND vigencia_fim IS NULL AND status = 'active'
		ORDER BY vigencia_inicio DESC
		LIMIT 1
	`
	row := s.pool.QueryRow(ctx, query, unitID)
	return scanContract(row)
}

func (s *pgxStore) ListContracts(ctx context.Context, filter ContractFilter) ([]Contract, error) {
	var args []interface{}
	var conds []string
	argNum := 1

	if filter.ConsumerUnitID != nil {
		conds = append(conds, fmt.Sprintf("consumer_unit_id = $%d", argNum))
		args = append(args, *filter.ConsumerUnitID)
		argNum++
	}
	if filter.ActiveOnly {
		conds = append(conds, "vigencia_fim IS NULL AND status = 'active'")
	}

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT id, customer_id, consumer_unit_id, vigencia_inicio, vigencia_fim, fator_repasse_energia, valor_ip_com_desconto, ip_faturamento_mode, ip_faturamento_valor, ip_faturamento_percent, bandeira_com_desconto, custo_disponibilidade_sempre_cobrado, consumo_minimo_kwh, notes, status, created_at, created_by, updated_at
		FROM public.contract
		%s
		ORDER BY vigencia_inicio DESC
	`, where)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contracts []Contract
	for rows.Next() {
		c, err := scanContractRows(rows)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, *c)
	}
	return contracts, rows.Err()
}

func (s *pgxStore) CloseContract(ctx context.Context, unitID uuid.UUID, closeDate interface{}) error {
	query := `
		UPDATE public.contract
		SET vigencia_fim = $1, status = 'ended', updated_at = NOW()
		WHERE consumer_unit_id = $2 AND vigencia_fim IS NULL AND status = 'active'
	`
	_, err := s.pool.Exec(ctx, query, closeDate, unitID)
	return err
}

// --- UserStore ---

func (s *pgxStore) CreateUser(ctx context.Context, u *User) error {
	// TODO: quando tabela app_user existir no Prisma
	return fmt.Errorf("not implemented")
}

func (s *pgxStore) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *pgxStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *pgxStore) ListUsers(ctx context.Context) ([]User, error) {
	return nil, fmt.Errorf("not implemented")
}

// --- Row scanners ---

func scanCustomer(row pgx.Row) (*Customer, error) {
	var c Customer
	var idStr string
	err := row.Scan(
		&idStr, &c.TipoPessoa, &c.NomeRazao, &c.NomeFantasia, &c.CPFCNPJ,
		&c.Email, &c.Telefone, &c.Status, &c.TipoCliente, &c.Notes,
		&c.CreatedAt, &c.UpdatedAt, &c.ArchivedAt,
	)
	if err != nil {
		return nil, err
	}
	c.ID, _ = uuid.Parse(idStr)
	return &c, nil
}

func scanCustomerRows(rows pgx.Rows) (*Customer, error) {
	var c Customer
	var idStr string
	err := rows.Scan(
		&idStr, &c.TipoPessoa, &c.NomeRazao, &c.NomeFantasia, &c.CPFCNPJ,
		&c.Email, &c.Telefone, &c.Status, &c.TipoCliente, &c.Notes,
		&c.CreatedAt, &c.UpdatedAt, &c.ArchivedAt,
	)
	if err != nil {
		return nil, err
	}
	c.ID, _ = uuid.Parse(idStr)
	return &c, nil
}

func scanConsumerUnit(row pgx.Row) (*ConsumerUnit, error) {
	var u ConsumerUnit
	var idStr, customerIDStr string
	err := row.Scan(
		&idStr, &customerIDStr, &u.UCCode, &u.Distribuidora, &u.Apelido,
		&u.ClasseConsumo, &u.Endereco, &u.Cidade, &u.UF, &u.Ativa,
		&u.CredentialID, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	u.ID, _ = uuid.Parse(idStr)
	u.CustomerID, _ = uuid.Parse(customerIDStr)
	return &u, nil
}

func scanConsumerUnitRows(rows pgx.Rows) (*ConsumerUnit, error) {
	var u ConsumerUnit
	var idStr, customerIDStr string
	err := rows.Scan(
		&idStr, &customerIDStr, &u.UCCode, &u.Distribuidora, &u.Apelido,
		&u.ClasseConsumo, &u.Endereco, &u.Cidade, &u.UF, &u.Ativa,
		&u.CredentialID, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	u.ID, _ = uuid.Parse(idStr)
	u.CustomerID, _ = uuid.Parse(customerIDStr)
	return &u, nil
}

func scanContract(row pgx.Row) (*Contract, error) {
	var c Contract
	var createdBy *uuid.UUID
	err := row.Scan(
		&c.ID, &c.CustomerID, &c.ConsumerUnitID, &c.VigenciaInicio, &c.VigenciaFim,
		&c.FatorRepasseEnergia, &c.ValorIPComDesconto, &c.IPFaturamentoMode, &c.IPFaturamentoValor,
		&c.IPFaturamentoPercent, &c.BandeiraComDesconto, &c.CustoDisponibilidadeSempreCobrado,
		&c.ConsumoMinimoKWh,
		&c.Notes, &c.Status, &c.CreatedAt, &createdBy, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	c.CreatedBy = createdBy
	return &c, nil
}

func scanContractRows(rows pgx.Rows) (*Contract, error) {
	var c Contract
	var createdBy *uuid.UUID
	err := rows.Scan(
		&c.ID, &c.CustomerID, &c.ConsumerUnitID, &c.VigenciaInicio, &c.VigenciaFim,
		&c.FatorRepasseEnergia, &c.ValorIPComDesconto, &c.IPFaturamentoMode, &c.IPFaturamentoValor,
		&c.IPFaturamentoPercent, &c.BandeiraComDesconto, &c.CustoDisponibilidadeSempreCobrado,
		&c.ConsumoMinimoKWh,
		&c.Notes, &c.Status, &c.CreatedAt, &createdBy, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	c.CreatedBy = createdBy
	return &c, nil
}
