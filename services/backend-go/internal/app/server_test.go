package app

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestDocsEndpoints(t *testing.T) {
	cfg := Config{
		Host:               "127.0.0.1",
		Port:               "8080",
		ExtractorBaseURL:   "http://127.0.0.1:8090",
		NeoenergiaBaseURL:  "http://127.0.0.1:9999",
		DatabaseURL:        "file::memory:?cache=shared",
		EncryptionKey:      "test-secret",
		BootstrapPythonBin: "/bin/false",
		BootstrapScript:    "scripts/bootstrap_neoenergia_token.py",
	}
	server, err := NewServer(cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	apiServer := httptest.NewServer(server.mux)
	defer apiServer.Close()

	htmlResp, err := http.Get(apiServer.URL + "/docs")
	if err != nil {
		t.Fatal(err)
	}
	defer htmlResp.Body.Close()
	if htmlResp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected /docs status: %d", htmlResp.StatusCode)
	}
	htmlBody, _ := io.ReadAll(htmlResp.Body)
	if !strings.Contains(string(htmlBody), "Backend Go API") {
		t.Fatalf("expected docs title in /docs response: %s", string(htmlBody))
	}

	mdResp, err := http.Get(apiServer.URL + "/docs.md")
	if err != nil {
		t.Fatal(err)
	}
	defer mdResp.Body.Close()
	if mdResp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected /docs.md status: %d", mdResp.StatusCode)
	}
	mdBody, _ := io.ReadAll(mdResp.Body)
	if !strings.Contains(string(mdBody), "/v1/consumer-units/{uc}/sync") {
		t.Fatalf("expected sync endpoint in /docs.md response: %s", string(mdBody))
	}
}

func TestSyncUCEndpoint(t *testing.T) {
	neoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/multilogin/2.0.0/agv/cliente/03021937586/COELBA/grupo-de-cliente":
			io.WriteString(w, `{"dados":[{"codigo":1,"nome":"Comum","ativo":true,"funcionalidades":[{"codigo":"F001"}]}]}`)
		case "/multilogin/2.0.0/servicos/minha-conta":
			io.WriteString(w, `{"nome":"Paula","usuarioAcesso":"Paula","email":"a@b.com","celular":"(74) 99999-9999","dtUltimaAtualizacao":"2026-04-16T22:00:00-03:00"}`)
		case "/multilogin/2.0.0/servicos/minha-conta/minha-conta-legado":
			io.WriteString(w, `{"nomeTitular":"Paula","dtNascimento":"1994-06-26","emailCadastro":"a@b.com","telefoneContato":"74-99999-9999"}`)
		case "/imoveis/1.1.0/clientes/03021937586/ucs":
			io.WriteString(w, `{"ucs":[{"status":"LIGADA","uc":"007098175908","nomeCliente":"PAULA","instalacao":"0001","grupoTensao":"B","contrato":"0023","dt_inicio":"2026-02-20","dt_fim":"9999-12-31","local":{"endereco":"Rua X","bairro":"Centro","municipio":"Lapao","cep":"44905-000","uf":"BA"}}]}`)
		case "/multilogin/2.0.0/servicos/imoveis/ucs/007098175908":
			io.WriteString(w, `{"codigo":"007098175908","instalacao":"0001","medidor":"1204300816","fase":"MONOFASE","dataLigacao":"20260220","situacao":{"codigo":"LG","descricao":"LIGADA","dataSituacaoUC":"20260220"},"cliente":{"codigo":"1014945628","nome":"PAULA"}}`)
		case "/protocolo/1.1.0/obterProtocolo":
			io.WriteString(w, `{"protocoloSalesforceStr":"20260416280227390","protocoloLegadoStr":"20260416280227390"}`)
		case "/multilogin/2.0.0/servicos/historicos/ucs/007098175908/consumos":
			io.WriteString(w, `{"historicoConsumo":[{"numeroFatura":"339800707843","mesReferencia":"04/2026","valorFatura":"521,53"}],"mediamensal":"394"}`)
		case "/multilogin/2.0.0/servicos/faturas/ucs/faturas":
			io.WriteString(w, `{"entregaFaturas":{"dataVencimento":"06"},"faturas":[{"numeroFatura":"339800707843","mesReferencia":"2026/04","statusFatura":"A Vencer","dataEmissao":"2026-04-14","dataPagamento":"0000-00-00","dataVencimento":"2026-05-06","valorEmissao":"521.53","uc":"007098175908"}]}`)
		case "/multilogin/2.0.0/servicos/datacerta/ucs/007098175908/datacerta":
			io.WriteString(w, `{"possuiDataBoa":"X","dataAtual":"06"}`)
		case "/multilogin/2.0.0/servicos/fatura-digital/ucs/fatura-digital":
			io.WriteString(w, `{"PossuiFaturaDigital":null,"emailFatura":null,"emailCadastro":"a@b.com"}`)
		case "/multilogin/2.0.0/servicos/faturas/lista-motivo-segundavia":
			io.WriteString(w, `{"motivos":[{"idMotivo":"02","descricao":"NÃO RECEBI - CLIENTE"}]}`)
		case "/multilogin/2.0.0/servicos/debito-automatico/conta-cadastrada-debito":
			io.WriteString(w, `{"retorno":{"mensagem":"UC não possui débito automático cadastrado"}}`)
		case "/multilogin/2.0.0/servicos/faturas/339800707843/dados-pagamento":
			io.WriteString(w, `{"codBarras":"838700000052215300300078098175908215056066285939"}`)
		case "/multilogin/2.0.0/servicos/faturas/339800707843/pdf":
			io.WriteString(w, `{"fileName":"007098175908","fileSize":"60147","fileData":"JVBERi0xLjQ=","fileExtension":".pdf"}`)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer neoServer.Close()

	cfg := Config{
		Host:               "127.0.0.1",
		Port:               "8080",
		ExtractorBaseURL:   "http://127.0.0.1:8090",
		NeoenergiaBaseURL:  neoServer.URL,
		DatabaseURL:        "file::memory:?cache=shared",
		EncryptionKey:      "test-secret",
		BootstrapPythonBin: "/bin/false",
		BootstrapScript:    "scripts/bootstrap_neoenergia_token.py",
	}
	server, err := NewServer(cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	apiServer := httptest.NewServer(server.mux)
	defer apiServer.Close()

	body := bytes.NewBufferString(`{"bearer_token":"abc","documento":"03021937586","uc":"007098175908","include_pdf":true}`)
	resp, err := http.Post(apiServer.URL+"/v1/sync/uc", "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}

	if payload["uc"] != "007098175908" {
		t.Fatalf("unexpected uc: %v", payload["uc"])
	}

	faturas := payload["faturas"].(map[string]any)
	if faturas["data"] == nil {
		t.Fatalf("expected faturas data")
	}

	pdf := payload["pdf_primeira_fatura"].(map[string]any)
	if pdf["data"] == nil {
		t.Fatalf("expected pdf data")
	}

	billingRecord := payload["billing_record"].(map[string]any)
	if billingRecord["numero_fatura"] != "339800707843" {
		t.Fatalf("unexpected billing record invoice: %v", billingRecord["numero_fatura"])
	}
	if billingRecord["codigo_barras"] != "838700000052215300300078098175908215056066285939" {
		t.Fatalf("unexpected barcode: %v", billingRecord["codigo_barras"])
	}
}

func TestCredentialAndSessionEndpointsDoNotExposeSecrets(t *testing.T) {
	pythonBin, err := exec.LookPath("python3")
	if err != nil {
		pythonBin, err = exec.LookPath("python")
		if err != nil {
			t.Skip("python not available")
		}
	}

	scriptPath := filepath.Join(t.TempDir(), "bootstrap_fake.py")
	script := `import json; print(json.dumps({"token":"secret-bearer-token","token_ne_se":{"ne":"x"},"local_storage":{"token":"x"}}))`
	if err := os.WriteFile(scriptPath, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	cfg := Config{
		Host:               "127.0.0.1",
		Port:               "8080",
		ExtractorBaseURL:   "http://127.0.0.1:8090",
		NeoenergiaBaseURL:  "http://127.0.0.1:9999",
		DatabaseURL:        "file::memory:?cache=shared",
		EncryptionKey:      "test-secret",
		BootstrapPythonBin: pythonBin,
		BootstrapScript:    scriptPath,
	}
	server, err := NewServer(cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	apiServer := httptest.NewServer(server.mux)
	defer apiServer.Close()

	credBody := bytes.NewBufferString(`{"label":"neo-paula","documento":"03021937586","senha":"MinhaSenha@123","uf":"BA","tipo_acesso":"normal"}`)
	credResp, err := http.Post(apiServer.URL+"/v1/credentials", "application/json", credBody)
	if err != nil {
		t.Fatal(err)
	}
	defer credResp.Body.Close()
	if credResp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected status: %d", credResp.StatusCode)
	}
	rawCredResp, _ := io.ReadAll(credResp.Body)
	credText := string(rawCredResp)
	if strings.Contains(credText, "03021937586") || strings.Contains(credText, "MinhaSenha@123") {
		t.Fatalf("credential response leaked secret: %s", credText)
	}

	var credPayload map[string]any
	if err := json.Unmarshal(rawCredResp, &credPayload); err != nil {
		t.Fatal(err)
	}
	credentialID := credPayload["id"].(string)
	if credentialID == "" {
		t.Fatal("expected credential id")
	}

	sessReq, err := http.NewRequest(http.MethodPost, apiServer.URL+"/v1/credentials/"+credentialID+"/session", bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	sessResp, err := http.DefaultClient.Do(sessReq)
	if err != nil {
		t.Fatal(err)
	}
	defer sessResp.Body.Close()
	if sessResp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", sessResp.StatusCode)
	}
	rawSessResp, _ := io.ReadAll(sessResp.Body)
	sessText := string(rawSessResp)
	if strings.Contains(sessText, "secret-bearer-token") || strings.Contains(sessText, "MinhaSenha@123") {
		t.Fatalf("session response leaked secret: %s", sessText)
	}
}

func TestSyncUCEndpointWithExtraction(t *testing.T) {
	extractorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/extract" {
			t.Fatalf("unexpected extractor path: %s", r.URL.Path)
		}
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		pdf, _ := req["pdf"].(map[string]any)
		if pdf["mode"] != "base64" {
			t.Fatalf("expected base64 pdf mode, got %v", pdf["mode"])
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"schema_version": "1.0.0",
			"job_id":         req["job_id"],
			"status":         "ok",
			"fields": map[string]any{
				"valor": "521.53",
			},
			"source_map": map[string]string{
				"valor": "pymupdf",
			},
			"confidence_map": map[string]float64{
				"valor": 0.9,
			},
			"warnings":  []string{},
			"artifacts": map[string]any{"pdf_path": "/tmp/fake.pdf"},
		})
	}))
	defer extractorServer.Close()

	neoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/multilogin/2.0.0/agv/cliente/03021937586/COELBA/grupo-de-cliente":
			io.WriteString(w, `{"dados":[{"codigo":1,"nome":"Comum","ativo":true,"funcionalidades":[{"codigo":"F001"}]}]}`)
		case "/multilogin/2.0.0/servicos/minha-conta":
			io.WriteString(w, `{"nome":"Paula","usuarioAcesso":"Paula","email":"a@b.com","celular":"(74) 99999-9999","dtUltimaAtualizacao":"2026-04-16T22:00:00-03:00"}`)
		case "/multilogin/2.0.0/servicos/minha-conta/minha-conta-legado":
			io.WriteString(w, `{"nomeTitular":"Paula","dtNascimento":"1994-06-26","emailCadastro":"a@b.com","telefoneContato":"74-99999-9999"}`)
		case "/imoveis/1.1.0/clientes/03021937586/ucs":
			io.WriteString(w, `{"ucs":[{"status":"LIGADA","uc":"007098175908","nomeCliente":"PAULA","instalacao":"0001","grupoTensao":"B","contrato":"0023","dt_inicio":"2026-02-20","dt_fim":"9999-12-31","local":{"endereco":"Rua X","bairro":"Centro","municipio":"Lapao","cep":"44905-000","uf":"BA"}}]}`)
		case "/multilogin/2.0.0/servicos/imoveis/ucs/007098175908":
			io.WriteString(w, `{"codigo":"007098175908","instalacao":"0001","medidor":"1204300816","fase":"MONOFASE","dataLigacao":"20260220","situacao":{"codigo":"LG","descricao":"LIGADA","dataSituacaoUC":"20260220"},"cliente":{"codigo":"1014945628","nome":"PAULA"}}`)
		case "/protocolo/1.1.0/obterProtocolo":
			io.WriteString(w, `{"protocoloSalesforceStr":"20260416280227390","protocoloLegadoStr":"20260416280227390"}`)
		case "/multilogin/2.0.0/servicos/historicos/ucs/007098175908/consumos":
			io.WriteString(w, `{"historicoConsumo":[{"numeroFatura":"339800707843","mesReferencia":"04/2026","valorFatura":"521,53"}],"mediamensal":"394"}`)
		case "/multilogin/2.0.0/servicos/faturas/ucs/faturas":
			io.WriteString(w, `{"entregaFaturas":{"dataVencimento":"06"},"faturas":[{"numeroFatura":"339800707843","mesReferencia":"2026/04","statusFatura":"A Vencer","dataEmissao":"2026-04-14","dataPagamento":"0000-00-00","dataVencimento":"2026-05-06","valorEmissao":"521.53","uc":"007098175908"}]}`)
		case "/multilogin/2.0.0/servicos/datacerta/ucs/007098175908/datacerta":
			io.WriteString(w, `{"possuiDataBoa":"X","dataAtual":"06"}`)
		case "/multilogin/2.0.0/servicos/fatura-digital/ucs/fatura-digital":
			io.WriteString(w, `{"PossuiFaturaDigital":null,"emailFatura":null,"emailCadastro":"a@b.com"}`)
		case "/multilogin/2.0.0/servicos/faturas/lista-motivo-segundavia":
			io.WriteString(w, `{"motivos":[{"idMotivo":"02","descricao":"NÃO RECEBI - CLIENTE"}]}`)
		case "/multilogin/2.0.0/servicos/debito-automatico/conta-cadastrada-debito":
			io.WriteString(w, `{"retorno":{"mensagem":"UC não possui débito automático cadastrado"}}`)
		case "/multilogin/2.0.0/servicos/faturas/339800707843/dados-pagamento":
			io.WriteString(w, `{"codBarras":"838700000052215300300078098175908215056066285939"}`)
		case "/multilogin/2.0.0/servicos/faturas/339800707843/pdf":
			io.WriteString(w, `{"fileName":"007098175908","fileSize":"60147","fileData":"JVBERi0xLjQ=","fileExtension":".pdf"}`)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer neoServer.Close()

	cfg := Config{
		Host:               "127.0.0.1",
		Port:               "8080",
		ExtractorBaseURL:   extractorServer.URL,
		NeoenergiaBaseURL:  neoServer.URL,
		DatabaseURL:        "file::memory:?cache=shared",
		EncryptionKey:      "test-secret",
		BootstrapPythonBin: "/bin/false",
		BootstrapScript:    "scripts/bootstrap_neoenergia_token.py",
	}
	server, err := NewServer(cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	apiServer := httptest.NewServer(server.mux)
	defer apiServer.Close()

	body := bytes.NewBufferString(`{"bearer_token":"abc","documento":"03021937586","uc":"007098175908","include_pdf":true,"include_extraction":true}`)
	resp, err := http.Post(apiServer.URL+"/v1/sync/uc", "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}

	extraction := payload["extraction"].(map[string]any)
	if extraction["data"] == nil {
		t.Fatalf("expected extraction data")
	}
	billingRecord := payload["billing_record"].(map[string]any)
	if billingRecord["extractor_status"] != "ok" {
		t.Fatalf("expected extractor status in billing record, got %v", billingRecord["extractor_status"])
	}
	sourceMap := billingRecord["source_map"].(map[string]any)
	if sourceMap["valor_total"] != "api" {
		t.Fatalf("expected api source for valor_total, got %v", sourceMap["valor_total"])
	}
	documentRecord := payload["document_record"].(map[string]any)
	if documentRecord["uc"] != "007098175908" {
		t.Fatalf("unexpected document record uc: %v", documentRecord["uc"])
	}
	ocr := documentRecord["ocr"].(map[string]any)
	if ocr["codigo_barras"] != "838700000052215300300078098175908215056066285939" {
		t.Fatalf("expected barcode merged into document record ocr, got %v", ocr["codigo_barras"])
	}
	persistence := payload["persistence"].(map[string]any)
	if persistence["invoice_id"] == "" || persistence["sync_run_id"] == "" {
		t.Fatalf("expected persistence ids, got %v", persistence)
	}

	invoiceResp, err := http.Get(apiServer.URL + "/v1/invoices/" + persistence["invoice_id"].(string))
	if err != nil {
		t.Fatal(err)
	}
	defer invoiceResp.Body.Close()
	if invoiceResp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected invoice status: %d", invoiceResp.StatusCode)
	}
	var invoicePayload map[string]any
	if err := json.NewDecoder(invoiceResp.Body).Decode(&invoicePayload); err != nil {
		t.Fatal(err)
	}
	if invoicePayload["numero_fatura"] != "339800707843" {
		t.Fatalf("unexpected invoice numero_fatura: %v", invoicePayload["numero_fatura"])
	}

	ucInvoicesResp, err := http.Get(apiServer.URL + "/v1/consumer-units/007098175908/invoices")
	if err != nil {
		t.Fatal(err)
	}
	defer ucInvoicesResp.Body.Close()
	if ucInvoicesResp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected uc invoices status: %d", ucInvoicesResp.StatusCode)
	}
	var ucInvoicesPayload map[string]any
	if err := json.NewDecoder(ucInvoicesResp.Body).Decode(&ucInvoicesPayload); err != nil {
		t.Fatal(err)
	}
	if len(ucInvoicesPayload["items"].([]any)) != 1 {
		t.Fatalf("expected one invoice for uc, got %v", len(ucInvoicesPayload["items"].([]any)))
	}

	filteredInvoicesResp, err := http.Get(apiServer.URL + "/v1/consumer-units/007098175908/invoices?status=A+Vencer&limit=1")
	if err != nil {
		t.Fatal(err)
	}
	defer filteredInvoicesResp.Body.Close()
	if filteredInvoicesResp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected filtered invoices status: %d", filteredInvoicesResp.StatusCode)
	}
	var filteredInvoicesPayload map[string]any
	if err := json.NewDecoder(filteredInvoicesResp.Body).Decode(&filteredInvoicesPayload); err != nil {
		t.Fatal(err)
	}
	if int(filteredInvoicesPayload["limit"].(float64)) != 1 {
		t.Fatalf("expected limit 1, got %v", filteredInvoicesPayload["limit"])
	}

	ucDetailsResp, err := http.Get(apiServer.URL + "/v1/consumer-units/007098175908")
	if err != nil {
		t.Fatal(err)
	}
	defer ucDetailsResp.Body.Close()
	if ucDetailsResp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected consumer unit details status: %d", ucDetailsResp.StatusCode)
	}
	var ucDetailsPayload map[string]any
	if err := json.NewDecoder(ucDetailsResp.Body).Decode(&ucDetailsPayload); err != nil {
		t.Fatal(err)
	}
	if ucDetailsPayload["uc"] != "007098175908" {
		t.Fatalf("unexpected consumer unit uc: %v", ucDetailsPayload["uc"])
	}
	if ucDetailsPayload["latest_invoice"] == nil || ucDetailsPayload["latest_sync_run"] == nil {
		t.Fatalf("expected latest_invoice and latest_sync_run in consumer unit details")
	}

	syncRunResp, err := http.Get(apiServer.URL + "/v1/sync-runs/" + persistence["sync_run_id"].(string))
	if err != nil {
		t.Fatal(err)
	}
	defer syncRunResp.Body.Close()
	if syncRunResp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected sync run status: %d", syncRunResp.StatusCode)
	}
	var syncRunPayload map[string]any
	if err := json.NewDecoder(syncRunResp.Body).Decode(&syncRunPayload); err != nil {
		t.Fatal(err)
	}
	if syncRunPayload["uc"] != "007098175908" {
		t.Fatalf("unexpected sync run uc: %v", syncRunPayload["uc"])
	}
}
