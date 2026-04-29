package catalog

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// testDBURL lê a URL do banco de teste da env var.
func testDBURL() string {
	if u := os.Getenv("TEST_DATABASE_URL"); u != "" {
		return u
	}
	// Fallback para banco de teste do redesign
	return "postgresql://azi:azi@localhost:5434/azi_billing"
}

// setupPool cria um pool de conexões para testes.
func setupPool(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, testDBURL())
	if err != nil {
		t.Fatalf("connect to test db: %v", err)
	}
	t.Cleanup(func() { pool.Close() })

	// Verifica conectividade
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping test db: %v", err)
	}

	return pool
}

// cleanTables limpa as tabelas usadas nos testes.
func strPtr(s string) *string { return &s }

func cleanTables(t *testing.T, pool *pgxpool.Pool) {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
		DELETE FROM public.billing_calculation WHERE 1=1;
		DELETE FROM public.contract WHERE 1=1;
		DELETE FROM public.consumer_unit WHERE 1=1;
		DELETE FROM public.customer WHERE 1=1;
	`)
	if err != nil {
		t.Fatalf("clean tables: %v", err)
	}
}

func TestCustomerStore(t *testing.T) {
	pool := setupPool(t)
	store := NewStore(pool)
	ctx := context.Background()

	t.Run("create and get customer", func(t *testing.T) {
		cleanTables(t, pool)

		c := &Customer{
			ID:          uuid.New(),
			TipoPessoa:  "PF",
			NomeRazao:   "João Silva",
			CPFCNPJ:     "12345678901",
			Status:      "prospect",
			TipoCliente: "residencial",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := store.CreateCustomer(ctx, c); err != nil {
			t.Fatalf("create customer: %v", err)
		}

		got, err := store.GetCustomer(ctx, c.ID)
		if err != nil {
			t.Fatalf("get customer: %v", err)
		}
		if got.NomeRazao != c.NomeRazao {
			t.Errorf("nome = %q, want %q", got.NomeRazao, c.NomeRazao)
		}
		if got.CPFCNPJ != c.CPFCNPJ {
			t.Errorf("cpf = %q, want %q", got.CPFCNPJ, c.CPFCNPJ)
		}
	})

	t.Run("get customer by cpf_cnpj", func(t *testing.T) {
		cleanTables(t, pool)

		c := &Customer{
			ID:          uuid.New(),
			TipoPessoa:  "PJ",
			NomeRazao:   "Empresa XYZ",
			CPFCNPJ:     "12345678000195",
			Status:      "active",
			TipoCliente: "empresa",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := store.CreateCustomer(ctx, c); err != nil {
			t.Fatalf("create customer: %v", err)
		}

		got, err := store.GetCustomerByCPFCNPJ(ctx, c.CPFCNPJ)
		if err != nil {
			t.Fatalf("get by cpf: %v", err)
		}
		if got.ID != c.ID {
			t.Errorf("id = %v, want %v", got.ID, c.ID)
		}
	})

	t.Run("list customers", func(t *testing.T) {
		cleanTables(t, pool)

		for i := 0; i < 3; i++ {
			c := &Customer{
				ID:          uuid.New(),
				TipoPessoa:  "PF",
				NomeRazao:   "Cliente " + string(rune('A'+i)),
				CPFCNPJ:     "1111111111" + string(rune('0'+i)),
				Status:      "active",
				TipoCliente: "residencial",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			if err := store.CreateCustomer(ctx, c); err != nil {
				t.Fatalf("create customer %d: %v", i, err)
			}
		}

		customers, nextCursor, err := store.ListCustomers(ctx, CustomerFilter{Limit: 2})
		if err != nil {
			t.Fatalf("list customers: %v", err)
		}
		if len(customers) != 2 {
			t.Errorf("len = %d, want 2", len(customers))
		}
		if nextCursor == "" {
			t.Error("expected next_cursor for pagination")
		}
	})

	t.Run("update customer", func(t *testing.T) {
		cleanTables(t, pool)

		c := &Customer{
			ID:          uuid.New(),
			TipoPessoa:  "PF",
			NomeRazao:   "Original",
			CPFCNPJ:     "99999999999",
			Status:      "prospect",
			TipoCliente: "residencial",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := store.CreateCustomer(ctx, c); err != nil {
			t.Fatalf("create: %v", err)
		}

		newName := "Atualizado"
		if err := store.UpdateCustomer(ctx, c.ID, CustomerPatch{NomeRazao: &newName}); err != nil {
			t.Fatalf("update: %v", err)
		}

		got, err := store.GetCustomer(ctx, c.ID)
		if err != nil {
			t.Fatalf("get after update: %v", err)
		}
		if got.NomeRazao != newName {
			t.Errorf("nome = %q, want %q", got.NomeRazao, newName)
		}
	})

	t.Run("archive customer", func(t *testing.T) {
		cleanTables(t, pool)

		c := &Customer{
			ID:          uuid.New(),
			TipoPessoa:  "PF",
			NomeRazao:   "Para Arquivar",
			CPFCNPJ:     "88888888888",
			Status:      "active",
			TipoCliente: "residencial",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := store.CreateCustomer(ctx, c); err != nil {
			t.Fatalf("create: %v", err)
		}

		if err := store.ArchiveCustomer(ctx, c.ID); err != nil {
			t.Fatalf("archive: %v", err)
		}

		got, err := store.GetCustomer(ctx, c.ID)
		if err != nil {
			t.Fatalf("get after archive: %v", err)
		}
		if got.Status != "archived" {
			t.Errorf("status = %q, want inativo", got.Status)
		}
		if got.ArchivedAt == nil {
			t.Error("expected archived_at to be set")
		}
	})
}

func TestConsumerUnitStore(t *testing.T) {
	pool := setupPool(t)
	store := NewStore(pool)
	ctx := context.Background()

	// Cria um customer base
	customer := &Customer{
		ID:          uuid.New(),
		TipoPessoa:  "PF",
		NomeRazao:   "Cliente Teste",
		CPFCNPJ:     "77777777777",
		Status:      "active",
		TipoCliente: "residencial",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := store.CreateCustomer(ctx, customer); err != nil {
		t.Fatalf("create customer: %v", err)
	}

	t.Run("create and get unit", func(t *testing.T) {
		cleanTables(t, pool)
		// Recria customer
		if err := store.CreateCustomer(ctx, customer); err != nil {
			t.Fatalf("create customer: %v", err)
		}

		u := &ConsumerUnit{
			ID:            uuid.New(),
			CustomerID:    customer.ID,
			UCCode:        "007085489032",
			Distribuidora: strPtr("neoenergia_ba"),
			Ativa:         true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := store.CreateUnit(ctx, u); err != nil {
			t.Fatalf("create unit: %v", err)
		}

		got, err := store.GetUnit(ctx, u.ID)
		if err != nil {
			t.Fatalf("get unit: %v", err)
		}
		if got.UCCode != u.UCCode {
			t.Errorf("uc_code = %q, want %q", got.UCCode, u.UCCode)
		}
	})

	t.Run("get unit by code", func(t *testing.T) {
		cleanTables(t, pool)
		if err := store.CreateCustomer(ctx, customer); err != nil {
			t.Fatalf("create customer: %v", err)
		}

		u := &ConsumerUnit{
			ID:         uuid.New(),
			CustomerID: customer.ID,
			UCCode:        "007085489033",
			Distribuidora: strPtr("neoenergia_ba"),
			Ativa:      true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if err := store.CreateUnit(ctx, u); err != nil {
			t.Fatalf("create unit: %v", err)
		}

		got, err := store.GetUnitByCode(ctx, u.UCCode)
		if err != nil {
			t.Fatalf("get by code: %v", err)
		}
		if got.ID != u.ID {
			t.Errorf("id = %v, want %v", got.ID, u.ID)
		}
	})

	t.Run("link unit to customer", func(t *testing.T) {
		cleanTables(t, pool)
		if err := store.CreateCustomer(ctx, customer); err != nil {
			t.Fatalf("create customer: %v", err)
		}

		u := &ConsumerUnit{
			ID:         uuid.New(),
			CustomerID: customer.ID,
			UCCode:        "007085489034",
			Distribuidora: strPtr("neoenergia_ba"),
			Ativa:      true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if err := store.CreateUnit(ctx, u); err != nil {
			t.Fatalf("create unit: %v", err)
		}

		// Cria novo customer
		newCustomer := &Customer{
			ID:          uuid.New(),
			TipoPessoa:  "PF",
			NomeRazao:   "Novo Cliente",
			CPFCNPJ:     "66666666666",
			Status:      "active",
			TipoCliente: "residencial",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := store.CreateCustomer(ctx, newCustomer); err != nil {
			t.Fatalf("create new customer: %v", err)
		}

		if err := store.LinkUnitToCustomer(ctx, u.ID, newCustomer.ID); err != nil {
			t.Fatalf("link: %v", err)
		}

		got, err := store.GetUnit(ctx, u.ID)
		if err != nil {
			t.Fatalf("get after link: %v", err)
		}
		if got.CustomerID != newCustomer.ID {
			t.Errorf("customer_id = %v, want %v", got.CustomerID, newCustomer.ID)
		}
	})
}

func TestContractStore(t *testing.T) {
	pool := setupPool(t)
	store := NewStore(pool)
	ctx := context.Background()

	// Cria customer e UC base
	customer := &Customer{
		ID:          uuid.New(),
		TipoPessoa:  "PF",
		NomeRazao:   "Cliente Contrato",
		CPFCNPJ:     "55555555555",
		Status:      "active",
		TipoCliente: "residencial",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := store.CreateCustomer(ctx, customer); err != nil {
		t.Fatalf("create customer: %v", err)
	}

	unit := &ConsumerUnit{
		ID:         uuid.New(),
		CustomerID: customer.ID,
		UCCode:        "007085489035",
			Distribuidora: strPtr("neoenergia_ba"),
		Ativa:      true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := store.CreateUnit(ctx, unit); err != nil {
		t.Fatalf("create unit: %v", err)
	}

	t.Run("create and get contract", func(t *testing.T) {
		c := &Contract{
			ID:                 uuid.New(),
			CustomerID:         customer.ID,
			ConsumerUnitID:     unit.ID,
			VigenciaInicio:     time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
			DescontoPercentual: "0.85",
			IPFaturamentoMode:  "fixed",
			IPFaturamentoValor: "10.00",
			IPFaturamentoPercent: "0",
			Status:             "active",
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}

		if err := store.CreateContract(ctx, c); err != nil {
			t.Fatalf("create contract: %v", err)
		}

		got, err := store.GetContract(ctx, c.ID)
		if err != nil {
			t.Fatalf("get contract: %v", err)
		}
		// Postgres NUMERIC pode adicionar zeros à direita
		if got.DescontoPercentual != c.DescontoPercentual && got.DescontoPercentual != c.DescontoPercentual+"00" {
			t.Errorf("desconto = %q, want %q", got.DescontoPercentual, c.DescontoPercentual)
		}
	})

	t.Run("get active contract", func(t *testing.T) {
		got, err := store.GetActiveContract(ctx, unit.ID)
		if err != nil {
			t.Fatalf("get active: %v", err)
		}
		if got == nil {
			t.Fatal("expected active contract")
		}
		if got.Status != "active" {
			t.Errorf("status = %q, want active", got.Status)
		}
	})

	t.Run("close contract", func(t *testing.T) {
		closeDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
		if err := store.CloseContract(ctx, unit.ID, closeDate); err != nil {
			t.Fatalf("close contract: %v", err)
		}

		_, err := store.GetActiveContract(ctx, unit.ID)
		// Deve retornar erro (não encontrado) ou nil
		if err == nil {
			t.Error("expected no active contract after close")
		}
	})
}
