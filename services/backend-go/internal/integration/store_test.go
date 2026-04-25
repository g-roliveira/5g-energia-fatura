package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func testPool(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		url = "postgresql://azi:azi@localhost:5434/azi_billing"
	}
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() { pool.Close() })
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping: %v", err)
	}
	return pool
}

func cleanIntegration(t *testing.T, pool *pgxpool.Pool) {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
		DELETE FROM integration.jobs WHERE 1=1;
		DELETE FROM integration.sync_runs WHERE 1=1;
		DELETE FROM integration.raw_invoices WHERE 1=1;
		DELETE FROM integration.consumer_units WHERE 1=1;
		DELETE FROM integration.sessions WHERE 1=1;
		DELETE FROM integration.credentials WHERE 1=1;
	`)
	if err != nil {
		t.Fatalf("clean integration: %v", err)
	}
}

func TestCredentialStore(t *testing.T) {
	pool := testPool(t)
	store := NewStore(pool)
	ctx := context.Background()

	t.Run("create and get credential", func(t *testing.T) {
		cleanIntegration(t, pool)

		c := &Credential{
			ID:              uuid.New(),
			Label:           "neo-paula",
			DocumentoCipher: "cipher123",
			DocumentoNonce:  "nonce123",
			SenhaCipher:     "senha456",
			SenhaNonce:      "nonce456",
			UF:              "BA",
			TipoAcesso:      "normal",
			KeyVersion:      "v1",
		}

		if err := store.InsertCredential(ctx, c); err != nil {
			t.Fatalf("insert: %v", err)
		}

		got, err := store.GetCredentialByID(ctx, c.ID.String())
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if got.Label != c.Label {
			t.Errorf("label = %q, want %q", got.Label, c.Label)
		}
	})
}

func TestSessionStore(t *testing.T) {
	pool := testPool(t)
	store := NewStore(pool)
	ctx := context.Background()

	t.Run("create and get session", func(t *testing.T) {
		cleanIntegration(t, pool)

		cred := &Credential{
			ID:              uuid.New(),
			Label:           "test",
			DocumentoCipher: "c",
			DocumentoNonce:  "n",
			SenhaCipher:     "s",
			SenhaNonce:      "n",
			UF:              "BA",
			TipoAcesso:      "normal",
			KeyVersion:      "v1",
		}
		if err := store.InsertCredential(ctx, cred); err != nil {
			t.Fatalf("insert cred: %v", err)
		}

		sess := &Session{
			ID:                uuid.New(),
			CredentialID:      cred.ID,
			BearerTokenCipher: "token-cipher",
			BearerTokenNonce:  "token-nonce",
		}
		if err := store.InsertSession(ctx, sess); err != nil {
			t.Fatalf("insert session: %v", err)
		}

		got, err := store.GetLatestSessionByCredentialID(ctx, cred.ID.String())
		if err != nil {
			t.Fatalf("get session: %v", err)
		}
		if got.BearerTokenCipher != sess.BearerTokenCipher {
			t.Errorf("token = %q, want %q", got.BearerTokenCipher, sess.BearerTokenCipher)
		}
	})
}

func TestConsumerUnitStore(t *testing.T) {
	pool := testPool(t)
	store := NewStore(pool)
	ctx := context.Background()

	t.Run("upsert and get consumer unit", func(t *testing.T) {
		cleanIntegration(t, pool)

		u := &ConsumerUnit{
			UC:          "007085489032",
			Status:      strPtr("active"),
			NomeCliente: strPtr("João Silva"),
			Instalacao:  strPtr("123456"),
			Contrato:    strPtr("C-001"),
			GrupoTensao: strPtr("bifasico"),
			Endereco:    map[string]any{"cidade": "Salvador", "uf": "BA"},
		}

		if err := store.UpsertConsumerUnit(ctx, u); err != nil {
			t.Fatalf("upsert: %v", err)
		}

		got, err := store.GetConsumerUnitByUC(ctx, u.UC)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if *got.NomeCliente != "João Silva" {
			t.Errorf("nome = %q, want %q", *got.NomeCliente, "João Silva")
		}

		// Update
		u.NomeCliente = strPtr("João Atualizado")
		if err := store.UpsertConsumerUnit(ctx, u); err != nil {
			t.Fatalf("upsert update: %v", err)
		}

		got, _ = store.GetConsumerUnitByUC(ctx, u.UC)
		if *got.NomeCliente != "João Atualizado" {
			t.Errorf("nome após update = %q", *got.NomeCliente)
		}
	})
}

func TestSyncRunStore(t *testing.T) {
	pool := testPool(t)
	store := NewStore(pool)
	ctx := context.Background()

	t.Run("insert and get sync run", func(t *testing.T) {
		cleanIntegration(t, pool)

		sr := &SyncRun{
			ID:       uuid.New(),
			Documento: "12345678901",
			UC:        "007085489032",
			Status:    "success",
			Step:      strPtr("extract"),
		}
		now := time.Now()
		sr.StartedAt = &now
		sr.FinishedAt = &now

		if err := store.InsertSyncRun(ctx, sr); err != nil {
			t.Fatalf("insert: %v", err)
		}

		got, err := store.GetSyncRunByID(ctx, sr.ID.String())
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if got.Status != "success" {
			t.Errorf("status = %q, want success", got.Status)
		}
	})
}

func TestJobQueue(t *testing.T) {
	pool := testPool(t)
	store := NewStore(pool)
	ctx := context.Background()

	t.Run("enqueue and claim job", func(t *testing.T) {
		cleanIntegration(t, pool)

		job, err := store.EnqueueJob(ctx, "sync_uc", map[string]any{"uc": "007085489032"})
		if err != nil {
			t.Fatalf("enqueue: %v", err)
		}
		if job.Status != "pending" {
			t.Errorf("status = %q, want pending", job.Status)
		}

		// Claim
		claimed, err := store.ClaimNextJob(ctx, "worker-1")
		if err != nil {
			t.Fatalf("claim: %v", err)
		}
		if claimed == nil {
			t.Fatal("expected claimed job")
		}
		if claimed.Status != "running" {
			t.Errorf("status = %q, want running", claimed.Status)
		}
		if claimed.ClaimedBy == nil || *claimed.ClaimedBy != "worker-1" {
			t.Error("expected claimed_by = worker-1")
		}

		// Complete
		if err := store.CompleteJob(ctx, claimed.ID, map[string]any{"result": "ok"}); err != nil {
			t.Fatalf("complete: %v", err)
		}

		// Verify no more pending jobs
		_, err = store.ClaimNextJob(ctx, "worker-2")
		if err == nil {
			t.Error("expected no more pending jobs")
		}
	})

	t.Run("claim with no pending jobs", func(t *testing.T) {
		cleanIntegration(t, pool)

		_, err := store.ClaimNextJob(ctx, "worker-1")
		if err == nil {
			t.Error("expected error when no pending jobs")
		}
	})
}

func strPtr(s string) *string { return &s }

func TestRawInvoiceNumericScan(t *testing.T) {
	pool := testPool(t)
	store := NewStore(pool)
	ctx := context.Background()
	
	// Cleanup
	pool.Exec(ctx, "DELETE FROM integration.raw_invoices WHERE uc = 'TEST123'")
	
	inv := &RawInvoice{
		UC:            "TEST123",
		NumeroFatura:  "999",
		MesReferencia: "04/2026",
		StatusFatura:  strPtr("A Vencer"),
		ValorTotal:    strPtr("521.53"),
	}
	_, err := store.UpsertRawInvoice(ctx, inv)
	if err != nil {
		t.Fatalf("upsert: %v", err)
	}
	
	got, err := store.GetRawInvoiceByID(ctx, inv.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ValorTotal == nil || *got.ValorTotal != "521.53" {
		t.Fatalf("valor_total = %v, want 521.53", got.ValorTotal)
	}
}
