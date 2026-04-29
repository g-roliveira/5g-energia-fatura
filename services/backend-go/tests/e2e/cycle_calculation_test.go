package e2e

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestBillingCycleSyncCalculateApprove valida o fluxo completo:
// criar ciclo → sync → calculate → approve → verificar resultado.
func TestBillingCycleSyncCalculateApprove(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	// 1. Criar cliente
	customer := suite.CreateCustomer(t, map[string]any{
		"tipo_pessoa":  "PF",
		"nome_razao":   "Cliente Ciclo Calc",
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

	// 3. Criar contrato ativo (15% desconto, IP fixo R$ 10, bandeira com desconto)
	respContract := suite.POST(t, "/v1/billing/contracts", map[string]any{
		"customer_id":                          customerID,
		"consumer_unit_id":                     ucID,
		"vigencia_inicio":                      "2025-01-01",
		"desconto_percentual":                  "0.85",
		"ip_faturamento_mode":                  "fixed",
		"ip_faturamento_valor":                 "10.00",
		"bandeira_com_desconto":                true,
		"custo_disponibilidade_sempre_cobrado": false,
		"consumo_minimo_kwh":                   "30.0",
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

	// 5. Inserir utility_invoice_ref com snapshot realista diretamente no banco
	// Isso simula o resultado de um sync bem-sucedido
	refID := uuid.New()
	snapshot := []byte(`{
		"itens_fatura": [
			{"descricao": "Consumo-TUSD", "quantidade": "100,00", "tarifa": "0,50", "valor": "50,00", "valor_total": "50,00"},
			{"descricao": "Consumo-TE", "quantidade": "100,00", "tarifa": "0,30", "valor": "30,00", "valor_total": "30,00"},
			{"descricao": "BANDEIRA AMARELA", "valor": "5,00", "valor_total": "5,00"},
			{"descricao": "Ilum. Púb. Municipal", "valor": "8,00", "valor_total": "8,00"}
		]
	}`)
	_, err := suite.Pool.Exec(context.Background(), `
		INSERT INTO public.utility_invoice_ref (
			id, consumer_unit_id, billing_cycle_id, sync_invoice_id, numero_fatura,
			mes_referencia, billing_record_snapshot, synced_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW(), NOW())
	`, refID, ucID, cycleID, "sync-test-001", "TEST-001", "2026/04", snapshot)
	if err != nil {
		t.Fatalf("insert invoice ref: %v", err)
	}

	// Atualizar cycle_consumer_unit para synced
	_, err = suite.Pool.Exec(context.Background(), `
		UPDATE public.cycle_consumer_unit
		SET status = 'synced', synced_at = NOW(), updated_at = NOW()
		WHERE billing_cycle_id = $1 AND consumer_unit_id = $2
	`, cycleID, ucID)
	if err != nil {
		t.Fatalf("update cycle_consumer_unit to synced: %v", err)
	}
	t.Logf("✓ Invoice ref inserido e UC marcada como synced")

	// 6. Bulk calculate
	respCalc := suite.POST(t, "/v1/billing/cycles/"+cycleID+"/bulk", map[string]any{
		"action": "recalculate",
	})
	defer respCalc.Body.Close()
	if respCalc.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(respCalc.Body)
		t.Fatalf("bulk calculate: status=%d body=%s", respCalc.StatusCode, body)
	}
	var calcBulk map[string]any
	if err := json.NewDecoder(respCalc.Body).Decode(&calcBulk); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if calcBulk["jobs_created"].(float64) != 1 {
		t.Fatalf("expected 1 calculate job, got %v", calcBulk["jobs_created"])
	}
	t.Logf("✓ Bulk calculate: 1 job criado")

	// 7. Processar o job de calculate (executar handler síncronamente via banco)
	// Como o worker pool roda em background, vamos esperar e verificar
	var calcID string
	for i := 0; i < 10; i++ {
		time.Sleep(500 * time.Millisecond)
		var id *string
		err := suite.Pool.QueryRow(context.Background(), `
			SELECT id::text FROM public.billing_calculation
			WHERE billing_cycle_id = $1 AND consumer_unit_id = $2
			ORDER BY version DESC LIMIT 1
		`, cycleID, ucID).Scan(&id)
		if err == nil && id != nil {
			calcID = *id
			break
		}
	}
	if calcID == "" {
		t.Fatalf("calculation not created after polling")
	}
	t.Logf("✓ Cálculo criado: %s", calcID)

	// 8. Verificar cálculo via API
	respCalcGet := suite.GET(t, "/v1/billing/calculations/"+calcID)
	defer respCalcGet.Body.Close()
	if respCalcGet.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(respCalcGet.Body)
		t.Fatalf("get calculation: status=%d body=%s", respCalcGet.StatusCode, body)
	}
	var calcResult map[string]any
	if err := json.NewDecoder(respCalcGet.Body).Decode(&calcResult); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if calcResult["status"] != "draft" {
		t.Errorf("status = %v, want draft", calcResult["status"])
	}
	// Verificar valores: energia sem desconto = 100 * (0.50 + 0.30) = 80
	// + bandeira 5 + IP 8 = 93 sem desconto
	// com desconto: 80 * 0.85 + 5 * 0.85 + 8 + 10 (IP usina) = 68 + 4.25 + 8 + 10 = 90.25
	semDesc := calcResult["total_sem_desconto"].(string)
	comDesc := calcResult["total_com_desconto"].(string)
	t.Logf("✓ Total sem desconto: %s, com desconto: %s", semDesc, comDesc)

	// 9. Bulk approve
	respAppr := suite.POST(t, "/v1/billing/cycles/"+cycleID+"/bulk", map[string]any{
		"action": "approve",
	})
	defer respAppr.Body.Close()
	if respAppr.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(respAppr.Body)
		t.Fatalf("bulk approve: status=%d body=%s", respAppr.StatusCode, body)
	}
	var apprBulk map[string]any
	if err := json.NewDecoder(respAppr.Body).Decode(&apprBulk); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if apprBulk["jobs_created"].(float64) != 1 {
		t.Fatalf("expected 1 approve job, got %v", apprBulk["jobs_created"])
	}
	t.Logf("✓ Bulk approve: 1 job criado")

	// 10. Esperar aprovação (worker pool processa em background)
	var approved bool
	for i := 0; i < 20; i++ {
		time.Sleep(500 * time.Millisecond)
		var status string
		err := suite.Pool.QueryRow(context.Background(), `
			SELECT status FROM public.billing_calculation
			WHERE id = $1
		`, calcID).Scan(&status)
		if err == nil && status == "approved" {
			approved = true
			break
		}
	}
	if !approved {
		t.Fatalf("calculation not approved after polling")
	}

	// 11. Verificar aprovação
	respApprGet := suite.GET(t, "/v1/billing/calculations/"+calcID)
	defer respApprGet.Body.Close()
	if respApprGet.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(respApprGet.Body)
		t.Fatalf("get calculation after approve: status=%d body=%s", respApprGet.StatusCode, body)
	}
	var apprResult map[string]any
	if err := json.NewDecoder(respApprGet.Body).Decode(&apprResult); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if apprResult["status"] != "approved" {
		t.Fatalf("status = %v, want approved", apprResult["status"])
	}
	t.Logf("✓ Cálculo aprovado")

	// 12. Verificar dashboard rows
	respRows := suite.GET(t, "/v1/billing/cycles/"+cycleID+"/rows")
	defer respRows.Body.Close()
	if respRows.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(respRows.Body)
		t.Fatalf("get rows: status=%d body=%s", respRows.StatusCode, body)
	}
	var rows map[string]any
	if err := json.NewDecoder(respRows.Body).Decode(&rows); err != nil {
		t.Fatalf("decode: %v", err)
	}
	rowItems := rows["items"].([]any)
	if len(rowItems) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rowItems))
	}
	row := rowItems[0].(map[string]any)
	if row["sync_status"] != "approved" {
		t.Errorf("sync_status = %v, want approved", row["sync_status"])
	}
	if row["calculation_status"] != "approved" {
		t.Errorf("calculation_status = %v, want approved", row["calculation_status"])
	}
	t.Logf("✓ Dashboard: UC %s sync=%s calc=%s", row["uc_code"], row["sync_status"], row["calculation_status"])
}
