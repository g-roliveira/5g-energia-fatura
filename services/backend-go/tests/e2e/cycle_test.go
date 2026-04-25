package e2e

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// TestBillingCycleCRUD valida criação, listagem e dashboard de ciclos.
func TestBillingCycleCRUD(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	// 1. Criar cliente
	customer := suite.CreateCustomer(t, map[string]any{
		"tipo_pessoa":  "PF",
		"nome_razao":   "Cliente Ciclo",
		"cpf_cnpj":     "52998224725",
		"tipo_cliente": "residencial",
	})
	customerID := customer["id"].(string)

	// 2. Criar UC
	uc := suite.CreateUnit(t, map[string]any{
		"customer_id":   customerID,
		"uc_code":       "007085489099",
		"distribuidora": "neoenergia_ba",
		"ativa":         true,
	})
	ucID := uc["id"].(string)
	_ = ucID

	// 3. Criar contrato ativo
	respContract := suite.POST(t, "/v1/billing/contracts", map[string]any{
		"customer_id":     customerID,
		"consumer_unit_id": uc["id"].(string),
		"vigencia_inicio": "2025-01-01",
		"desconto_percentual": "0.85",
		"ip_faturamento_mode": "fixed",
		"ip_faturamento_valor": "10.00",
	})
	defer respContract.Body.Close()
	if respContract.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(respContract.Body)
		t.Fatalf("create contract: status=%d body=%s", respContract.StatusCode, body)
	}

	// 4. Criar ciclo
	resp := suite.POST(t, "/v1/billing/cycles", map[string]any{
		"year":               2026,
		"month":              4,
		"include_all_active": true,
	})
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("create cycle: status=%d body=%s", resp.StatusCode, body)
	}
	var created map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode: %v", err)
	}
	cycleID := created["id"].(string)
	t.Logf("✓ Ciclo criado: %s", cycleID)

	// 5. Verificar ciclo
	resp2 := suite.GET(t, "/v1/billing/cycles/"+cycleID)
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp2.Body)
		t.Fatalf("get cycle: status=%d body=%s", resp2.StatusCode, body)
	}
	var got map[string]any
	if err := json.NewDecoder(resp2.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got["year"].(float64) != 2026 {
		t.Errorf("year = %v, want 2026", got["year"])
	}
	if got["status"] != "open" {
		t.Errorf("status = %v, want open", got["status"])
	}
	// Deve ter 1 UC associada
	if got["total_ucs"].(float64) != 1 {
		t.Errorf("total_ucs = %v, want 1", got["total_ucs"])
	}
	t.Logf("✓ Ciclo recuperado: %d UCs", int(got["total_ucs"].(float64)))

	// 6. Listar ciclos
	resp3 := suite.GET(t, "/v1/billing/cycles?year=2026")
	defer resp3.Body.Close()
	if resp3.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp3.Body)
		t.Fatalf("list cycles: status=%d body=%s", resp3.StatusCode, body)
	}
	var list map[string]any
	if err := json.NewDecoder(resp3.Body).Decode(&list); err != nil {
		t.Fatalf("decode: %v", err)
	}
	items := list["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("expected 1 cycle, got %d", len(items))
	}
	t.Logf("✓ %d ciclo(s) listado(s)", len(items))

	// 7. Dashboard (rows)
	resp4 := suite.GET(t, "/v1/billing/cycles/"+cycleID+"/rows")
	defer resp4.Body.Close()
	if resp4.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp4.Body)
		t.Fatalf("get rows: status=%d body=%s", resp4.StatusCode, body)
	}
	var rows map[string]any
	if err := json.NewDecoder(resp4.Body).Decode(&rows); err != nil {
		t.Fatalf("decode: %v", err)
	}
	rowItems := rows["items"].([]any)
	if len(rowItems) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rowItems))
	}
	row := rowItems[0].(map[string]any)
	if row["uc_code"] != "007085489099" {
		t.Errorf("uc_code = %v", row["uc_code"])
	}
	if row["sync_status"] != "pending" {
		t.Errorf("sync_status = %v, want pending", row["sync_status"])
	}
	t.Logf("✓ Dashboard: UC %s status=%s", row["uc_code"], row["sync_status"])

	// 8. Bulk action (sync) — deve enfileirar job
	resp5 := suite.POST(t, "/v1/billing/cycles/"+cycleID+"/bulk", map[string]any{
		"action": "sync",
	})
	defer resp5.Body.Close()
	if resp5.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp5.Body)
		t.Fatalf("bulk sync: status=%d body=%s", resp5.StatusCode, body)
	}
	var bulkResult map[string]any
	if err := json.NewDecoder(resp5.Body).Decode(&bulkResult); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if bulkResult["jobs_created"].(float64) != 1 {
		t.Errorf("jobs_created = %v, want 1", bulkResult["jobs_created"])
	}
	t.Logf("✓ Bulk sync: %v job(s) criado(s)", bulkResult["jobs_created"])
}
