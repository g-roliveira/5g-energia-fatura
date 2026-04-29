package e2e

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/integration"
	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/platform/worker"
)

// TestIntegrationCredentialCRUD valida criação e leitura de credenciais.
func TestIntegrationCredentialCRUD(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	// 1. Criar credencial
	resp := suite.POST(t, "/v1/integration/credentials", map[string]any{
		"label":            "test-cred",
		"documento_cipher": "cipher123",
		"documento_nonce":  "nonce123",
		"senha_cipher":     "senha456",
		"senha_nonce":      "nonce456",
		"uf":               "BA",
		"tipo_acesso":      "normal",
		"key_version":      "v1",
	})
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("create credential: status=%d body=%s", resp.StatusCode, body)
	}
	var created map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode created: %v", err)
	}
	credID := created["id"].(string)
	t.Logf("✓ Credential criada: %s", credID)

	// 2. Ler credencial
	resp2 := suite.GET(t, "/v1/integration/credentials/"+credID)
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp2.Body)
		t.Fatalf("get credential: status=%d body=%s", resp2.StatusCode, body)
	}
	var got map[string]any
	if err := json.NewDecoder(resp2.Body).Decode(&got); err != nil {
		t.Fatalf("decode got: %v", err)
	}
	if got["label"] != "test-cred" {
		t.Errorf("label = %q, want %q", got["label"], "test-cred")
	}
	if got["uf"] != "BA" {
		t.Errorf("uf = %q, want %q", got["uf"], "BA")
	}
	t.Logf("✓ Credential lida com sucesso")
}

// TestIntegrationSyncRunCRUD valida criação e leitura de sync runs.
func TestIntegrationSyncRunCRUD(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	// 1. Criar sync run
	resp := suite.POST(t, "/v1/integration/sync-runs", map[string]any{
		"documento": "12345678901",
		"uc":        "007085489032",
	})
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("create sync run: status=%d body=%s", resp.StatusCode, body)
	}
	var created map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode created: %v", err)
	}
	syncID := created["id"].(string)
	if created["status"] != "pending" {
		t.Errorf("status = %q, want pending", created["status"])
	}
	t.Logf("✓ SyncRun criado: %s", syncID)

	// 2. Ler sync run
	resp2 := suite.GET(t, "/v1/integration/sync-runs/"+syncID)
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp2.Body)
		t.Fatalf("get sync run: status=%d body=%s", resp2.StatusCode, body)
	}
	var got map[string]any
	if err := json.NewDecoder(resp2.Body).Decode(&got); err != nil {
		t.Fatalf("decode got: %v", err)
	}
	if got["uc"] != "007085489032" {
		t.Errorf("uc = %q, want 007085489032", got["uc"])
	}
	t.Logf("✓ SyncRun lido com sucesso")
}

// TestIntegrationJobQueue valida enqueue e worker pool claim.
func TestIntegrationJobQueue(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	// 1. Enfileirar job
	resp := suite.POST(t, "/v1/integration/jobs", map[string]any{
		"job_type": "sync_uc",
		"payload": map[string]any{
			"uc": "007085489032",
		},
	})
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("enqueue job: status=%d body=%s", resp.StatusCode, body)
	}
	var created map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode created: %v", err)
	}
	jobID := created["id"].(string)
	if created["status"] != "pending" {
		t.Errorf("status = %q, want pending", created["status"])
	}
	t.Logf("✓ Job enfileirado: %s", jobID)

	// 2. Claim via worker pool (direto no store, não via HTTP)
	store := integration.NewStore(suite.Pool)
	pool := worker.NewPool(store, 1, 100*time.Millisecond, suite.Logger)

	claimed := make(chan bool, 1)
	pool.RegisterHandler("sync_uc", func(ctx context.Context, job *integration.Job) error {
		t.Logf("✓ Job processado: %s payload=%v", job.ID, job.Payload)
		claimed <- true
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool.Start(ctx)

	select {
	case <-claimed:
		t.Log("✓ Worker pool processou job com sucesso")
	case <-ctx.Done():
		t.Fatal("worker pool não processou job a tempo")
	}

	pool.Stop()

	// 3. Verificar que job está completed
	job, err := store.ClaimNextJob(context.Background(), "worker-test")
	if err == nil && job != nil {
		t.Fatalf("esperava que não houvesse mais jobs pendentes, mas encontrou: %s", job.ID)
	}
	t.Log("✓ Nenhum job pendente restante")
}

// TestIntegrationConsumerUnitUpsert valida upsert e listagem de UCs.
func TestIntegrationConsumerUnitUpsert(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	// Inserir UC diretamente via service
	store := integration.NewStore(suite.Pool)
	svc := integration.NewService(store)
	ctx := context.Background()

	status := "active"
	nome := "João Silva"
	err := svc.SyncConsumerUnit(ctx, &integration.ConsumerUnit{
		UC:          "007085489099",
		Status:      &status,
		NomeCliente: &nome,
		Endereco:    map[string]any{"cidade": "Salvador", "uf": "BA"},
	})
	if err != nil {
		t.Fatalf("sync consumer unit: %v", err)
	}
	t.Logf("✓ ConsumerUnit upserted: 007085489099")

	// Listar UCs via HTTP
	resp := suite.GET(t, "/v1/integration/consumer-units")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("list consumer units: status=%d body=%s", resp.StatusCode, body)
	}
	var list map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	items := list["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	item := items[0].(map[string]any)
	if item["uc"] != "007085489099" {
		t.Errorf("uc = %q, want 007085489099", item["uc"])
	}
	t.Logf("✓ ConsumerUnit listado via HTTP")

	// Get UC por código
	resp2 := suite.GET(t, "/v1/integration/consumer-units/007085489099")
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp2.Body)
		t.Fatalf("get consumer unit: status=%d body=%s", resp2.StatusCode, body)
	}
	var got map[string]any
	if err := json.NewDecoder(resp2.Body).Decode(&got); err != nil {
		t.Fatalf("decode got: %v", err)
	}
	if got["nome_cliente"] != "João Silva" {
		t.Errorf("nome_cliente = %q, want João Silva", got["nome_cliente"])
	}
	t.Logf("✓ ConsumerUnit recuperado via HTTP")
}
