package store

import (
	"context"
	"os"
	"testing"

	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/neoenergia"
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

func TestAdapterPersistAndGetInvoice(t *testing.T) {
	ctx := context.Background()
	pool := testPool(t)

	// Cleanup
	pool.Exec(ctx, "DELETE FROM public.integration_raw_invoice_items WHERE raw_invoice_id IN (SELECT id FROM public.integration_raw_invoices WHERE uc = 'TEST456')")
	pool.Exec(ctx, "DELETE FROM public.integration_raw_invoices WHERE uc = 'TEST456'")
	pool.Exec(ctx, "DELETE FROM public.integration_sync_runs WHERE uc = 'TEST456'")

	adapter, err := OpenIntegrationPostgres("postgresql://azi:azi@localhost:5434/azi_billing")
	if err != nil {
		t.Fatalf("open adapter: %v", err)
	}
	defer adapter.Close()

	result, err := adapter.PersistSyncResult(PersistSyncInput{
		Documento: "12345678901",
		UC:        "TEST456",
		Fatura: &neoenergia.Fatura{
			UC:             "TEST456",
			NumeroFatura:   "888",
			MesReferencia:  "04/2026",
			StatusFatura:   "A Vencer",
			ValorEmissao:   "521.53",
			DataEmissao:    "2026-04-14",
			DataVencimento: "2026-05-06",
			DataPagamento:  "0000-00-00",
		},
		BillingRecord: map[string]any{"completeness": map[string]any{"status": "complete", "missing_fields": []any{}}},
	})
	if err != nil {
		t.Fatalf("persist: %v", err)
	}

	inv, err := adapter.GetInvoiceByID(result.InvoiceID)
	if err != nil {
		t.Fatalf("get invoice: %v", err)
	}
	if inv == nil {
		t.Fatal("invoice not found")
	}
	if inv.NumeroFatura != "888" {
		t.Fatalf("numero_fatura = %s, want 888", inv.NumeroFatura)
	}
}
