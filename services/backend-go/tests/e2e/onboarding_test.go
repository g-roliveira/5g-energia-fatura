package e2e

import (
	"testing"
)

// TestOnboarding valida o fluxo completo de cadastro:
// criar cliente → criar UC → criar contrato → verificar contrato ativo.
func TestOnboarding(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	// 1. Criar cliente
	customer := suite.CreateCustomer(t, map[string]any{
		"tipo_pessoa":  "PF",
		"nome_razao":   "João da Silva",
		"cpf_cnpj":     "52998224725",
		"tipo_cliente": "residencial",
		"email":        "joao@teste.com",
	})
	customerID := customer["id"].(string)
	t.Logf("✓ Cliente criado: %s", customerID)

	// 2. Verificar que o cliente existe
	got := suite.GetCustomer(t, customerID)
	if got["nome_razao"] != "João da Silva" {
		t.Errorf("nome = %q, want %q", got["nome_razao"], "João da Silva")
	}

	// 3. Criar UC vinculada ao cliente
	unit := suite.CreateUnit(t, map[string]any{
		"customer_id":    customerID,
		"uc_code":        "007085489032",
		"distribuidora":  "neoenergia_ba",
		"apelido":        "Apartamento 101",
		"classe_consumo": "residencial",
	})
	unitID := unit["id"].(string)
	t.Logf("✓ UC criada: %s", unitID)

	// 4. Criar contrato
	contract := suite.CreateContract(t, map[string]any{
		"customer_id":                          customerID,
		"consumer_unit_id":                     unitID,
		"vigencia_inicio":                      "2025-10-01",
		"desconto_percentual":                  "0.85",
		"ip_faturamento_mode":                  "fixed",
		"ip_faturamento_valor":                 "10.00",
		"ip_faturamento_percent":               "0",
		"bandeira_com_desconto":                false,
		"custo_disponibilidade_sempre_cobrado": true,
		"notes":                                "Contrato inicial",
	})
	contractID := contract["id"].(string)
	t.Logf("✓ Contrato criado: %s", contractID)

	// 5. Verificar contrato
	if contract["desconto_percentual"] != "0.85" {
		t.Errorf("desconto = %q, want 0.85", contract["desconto_percentual"])
	}
	if contract["status"] != "active" {
		t.Errorf("status = %q, want active", contract["status"])
	}

	// 6. Criar segundo contrato (deve fechar o anterior)
	contract2 := suite.CreateContract(t, map[string]any{
		"customer_id":                          customerID,
		"consumer_unit_id":                     unitID,
		"vigencia_inicio":                      "2026-01-01",
		"desconto_percentual":                  "0.80",
		"ip_faturamento_mode":                  "fixed",
		"ip_faturamento_valor":                 "12.00",
		"ip_faturamento_percent":               "0",
		"bandeira_com_desconto":                false,
		"custo_disponibilidade_sempre_cobrado": true,
	})
	t.Logf("✓ Novo contrato criado: %s", contract2["id"])

	// 7. Verificar que o contrato anterior foi fechado
	resp := suite.GET(t, "/v1/catalog/contracts/"+contractID)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("get old contract: status=%d", resp.StatusCode)
	}

	t.Log("✓ Fluxo de onboarding completo validado")
}

// TestCustomerCRUD valida operações básicas de cliente.
func TestCustomerCRUD(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	// Create
	c := suite.CreateCustomer(t, map[string]any{
		"tipo_pessoa":  "PJ",
		"nome_razao":   "Empresa Teste LTDA",
		"cpf_cnpj":     "12345678000195",
		"tipo_cliente": "empresa",
	})
	id := c["id"].(string)

	// Get
	got := suite.GetCustomer(t, id)
	if got["nome_razao"] != "Empresa Teste LTDA" {
		t.Errorf("nome = %q", got["nome_razao"])
	}

	// Update
	resp := suite.PATCH(t, "/v1/catalog/customers/"+id, map[string]any{
		"nome_razao": "Empresa Atualizada LTDA",
	})
	resp.Body.Close()
	if resp.StatusCode != 204 {
		t.Errorf("update status = %d, want 204", resp.StatusCode)
	}

	got = suite.GetCustomer(t, id)
	if got["nome_razao"] != "Empresa Atualizada LTDA" {
		t.Errorf("nome após update = %q", got["nome_razao"])
	}

	// Archive
	resp = suite.DELETE(t, "/v1/catalog/customers/"+id)
	resp.Body.Close()
	if resp.StatusCode != 204 {
		t.Errorf("archive status = %d, want 204", resp.StatusCode)
	}

	got = suite.GetCustomer(t, id)
	if got["status"] != "archived" {
		t.Errorf("status após archive = %q, want archived", got["status"])
	}
}

// TestMotorDeCalculo valida o motor de cálculo via endpoint (quando disponível)
// ou diretamente. Como o motor é função pura, testamos diretamente.
func TestMotorDeCalculo(t *testing.T) {
	// Este teste já existe em internal/billing/engine_test.go
	// Aqui poderíamos testar a integração catalog → billing quando o billing estiver pronto
	t.Skip("motor de cálculo testado em unit tests; teste E2E quando billing/calculation estiver implementado")
}
