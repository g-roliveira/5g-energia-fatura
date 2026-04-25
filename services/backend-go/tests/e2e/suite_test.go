package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/app"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

var testDBURL = func() string {
	if u := os.Getenv("TEST_DATABASE_URL"); u != "" {
		return u
	}
	return "postgresql://azi:azi@localhost:5434/azi_billing"
}()

// TestSuite encapsula servidor HTTP + banco para testes E2E.
type TestSuite struct {
	Server *httptest.Server
	Pool   *pgxpool.Pool
	Client *http.Client
}

// NewTestSuite cria um servidor HTTP de teste com banco real.
func NewTestSuite(t *testing.T) *TestSuite {
	ctx := context.Background()

	// Conecta ao banco
	pool, err := pgxpool.New(ctx, testDBURL)
	if err != nil {
		t.Fatalf("connect to test db: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping test db: %v", err)
	}

	// Limpa dados antes dos testes
	cleanAll(t, pool)

	// Cria servidor
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := app.Config{
		Host:            "127.0.0.1",
		Port:            "0",
		BackofficePGURL: testDBURL,
		ExtractorBaseURL: "http://127.0.0.1:8090",
		NeoenergiaBaseURL: "https://apineprd.neoenergia.com",
		EncryptionKey:   "test-key-32-bytes-long-for-aes!!",
		BootstrapPythonBin: "python3",
		BootstrapScript:    "scripts/bootstrap_neoenergia_token.py",
	}

	srv, err := app.NewServer(cfg, logger)
	if err != nil {
		t.Fatalf("create server: %v", err)
	}

	// httptest server
	ts := httptest.NewServer(srv.Mux())

	return &TestSuite{
		Server: ts,
		Pool:   pool,
		Client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *TestSuite) Close() {
	s.Server.Close()
	s.Pool.Close()
}

// cleanAll limpa todas as tabelas de teste.
func cleanAll(t *testing.T, pool *pgxpool.Pool) {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
		DELETE FROM public.manual_adjustment WHERE 1=1;
		DELETE FROM public.generated_document WHERE 1=1;
		DELETE FROM public.billing_calculation WHERE 1=1;
		DELETE FROM public.utility_invoice_scee WHERE 1=1;
		DELETE FROM public.utility_invoice_item WHERE 1=1;
		DELETE FROM public.utility_invoice_ref WHERE 1=1;
		DELETE FROM public.sync_job WHERE 1=1;
		DELETE FROM public.audit_log WHERE 1=1;
		DELETE FROM public.contract WHERE 1=1;
		DELETE FROM public.billing_cycle WHERE 1=1;
		DELETE FROM public.credential_link WHERE 1=1;
		DELETE FROM public.address WHERE 1=1;
		DELETE FROM public.consumer_unit WHERE 1=1;
		DELETE FROM public.customer WHERE 1=1;
	`)
	if err != nil {
		t.Fatalf("clean all tables: %v", err)
	}
}

// HTTP helpers

func (s *TestSuite) POST(t *testing.T, path string, body any) *http.Response {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		t.Fatalf("encode body: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, s.Server.URL+path, &buf)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.Client.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	return resp
}

func (s *TestSuite) GET(t *testing.T, path string) *http.Response {
	resp, err := s.Client.Get(s.Server.URL + path)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	return resp
}

func (s *TestSuite) PATCH(t *testing.T, path string, body any) *http.Response {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		t.Fatalf("encode body: %v", err)
	}
	req, err := http.NewRequest(http.MethodPatch, s.Server.URL+path, &buf)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.Client.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	return resp
}

func (s *TestSuite) DELETE(t *testing.T, path string) *http.Response {
	req, err := http.NewRequest(http.MethodDelete, s.Server.URL+path, nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	resp, err := s.Client.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	return resp
}

// Factory helpers

func (s *TestSuite) CreateCustomer(t *testing.T, input map[string]any) map[string]any {
	resp := s.POST(t, "/v1/catalog/customers", input)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("create customer: status=%d body=%s", resp.StatusCode, body)
	}
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode customer response: %v", err)
	}
	return result
}

func (s *TestSuite) CreateUnit(t *testing.T, input map[string]any) map[string]any {
	resp := s.POST(t, "/v1/catalog/units", input)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("create unit: status=%d body=%s", resp.StatusCode, body)
	}
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode unit response: %v", err)
	}
	return result
}

func (s *TestSuite) CreateContract(t *testing.T, input map[string]any) map[string]any {
	resp := s.POST(t, "/v1/catalog/contracts", input)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("create contract: status=%d body=%s", resp.StatusCode, body)
	}
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode contract response: %v", err)
	}
	return result
}

func (s *TestSuite) GetCustomer(t *testing.T, id string) map[string]any {
	resp := s.GET(t, fmt.Sprintf("/v1/catalog/customers/%s", id))
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("get customer: status=%d body=%s", resp.StatusCode, body)
	}
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode customer: %v", err)
	}
	return result
}
