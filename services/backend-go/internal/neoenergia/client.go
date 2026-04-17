package neoenergia

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultCanalSolicitante = "AGC"
	defaultUsuarioSonda     = "WSO2_CONEXAO"
	defaultDistribuidora    = "COELBA"
	defaultRegiao           = "NE"
	defaultTipoPerfil       = "1"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 45 * time.Second,
		},
	}
}

type RequestContext struct {
	BearerToken string
	Documento   string
}

type ErrorResponse struct {
	StatusCode int             `json:"status_code"`
	Method     string          `json:"method"`
	Path       string          `json:"path"`
	Body       json.RawMessage `json:"body,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%s %s returned %d", e.Method, e.Path, e.StatusCode)
}

type GrupoClienteResponse struct {
	Dados []struct {
		Codigo          int    `json:"codigo"`
		Nome            string `json:"nome"`
		Ativo           bool   `json:"ativo"`
		Funcionalidades []struct {
			Codigo string `json:"codigo"`
		} `json:"funcionalidades"`
	} `json:"dados"`
}

type MinhaContaResponse struct {
	Nome                string `json:"nome"`
	UsuarioAcesso       string `json:"usuarioAcesso"`
	Email               string `json:"email"`
	Celular             string `json:"celular"`
	DtUltimaAtualizacao string `json:"dtUltimaAtualizacao"`
}

type MinhaContaLegadoResponse struct {
	NomeTitular     string `json:"nomeTitular"`
	DtNascimento    string `json:"dtNascimento"`
	EmailCadastro   string `json:"emailCadastro"`
	TelefoneContato string `json:"telefoneContato"`
}

type UCsResponse struct {
	UCs []UC `json:"ucs"`
}

type UC struct {
	Status      string `json:"status"`
	UC          string `json:"uc"`
	NomeCliente string `json:"nomeCliente"`
	Instalacao  string `json:"instalacao"`
	GrupoTensao string `json:"grupoTensao"`
	Contrato    string `json:"contrato"`
	DtInicio    string `json:"dt_inicio"`
	DtFim       string `json:"dt_fim"`
	Local       struct {
		Endereco  string `json:"endereco"`
		Bairro    string `json:"bairro"`
		Municipio string `json:"municipio"`
		CEP       string `json:"cep"`
		UF        string `json:"uf"`
	} `json:"local"`
}

type ImovelResponse struct {
	Codigo      string `json:"codigo"`
	Instalacao  string `json:"instalacao"`
	Medidor     string `json:"medidor"`
	Fase        string `json:"fase"`
	DataLigacao string `json:"dataLigacao"`
	Situacao    struct {
		Codigo         string `json:"codigo"`
		Descricao      string `json:"descricao"`
		DataSituacaoUC string `json:"dataSituacaoUC"`
	} `json:"situacao"`
	Cliente struct {
		Codigo string `json:"codigo"`
		Nome   string `json:"nome"`
	} `json:"cliente"`
}

type ProtocoloResponse struct {
	ProtocoloSalesforceStr string `json:"protocoloSalesforceStr"`
	ProtocoloLegadoStr     string `json:"protocoloLegadoStr"`
}

type FaturasResponse struct {
	EntregaFaturas map[string]any `json:"entregaFaturas"`
	Faturas        []Fatura       `json:"faturas"`
}

type Fatura struct {
	NumeroFatura      string `json:"numeroFatura"`
	MesReferencia     string `json:"mesReferencia"`
	StatusFatura      string `json:"statusFatura"`
	DataCompetencia   string `json:"dataCompetencia"`
	DataEmissao       string `json:"dataEmissao"`
	DataPagamento     string `json:"dataPagamento"`
	DataVencimento    string `json:"dataVencimento"`
	DataInicioPeriodo string `json:"dataInicioPeriodo"`
	DataFimPeriodo    string `json:"dataFimPeriodo"`
	ValorEmissao      string `json:"valorEmissao"`
	UC                string `json:"uc"`
	TipoDoc           string `json:"tipoDoc"`
	OrigemFatura      string `json:"origemFatura"`
	TipoEntrega       string `json:"tipoEntrega"`
	TipoLeitura       string `json:"tipoLeitura"`
	AceitaPix         string `json:"aceitaPix"`
}

type HistoricoConsumoResponse struct {
	HistoricoConsumo []HistoricoConsumo `json:"historicoConsumo"`
	MediaMensal      string             `json:"mediamensal"`
}

type HistoricoConsumo struct {
	DataPagamento     string `json:"dataPagamento"`
	DataVencimento    string `json:"dataVencimento"`
	DataLeitura       string `json:"dataLeitura"`
	ConsumoKW         string `json:"consumoKw"`
	MesReferencia     string `json:"mesReferencia"`
	NumeroLeitura     string `json:"numeroLeitura"`
	TipoLeitura       string `json:"tipoLeitura"`
	DataInicioPeriodo string `json:"dataInicioPeriodoCalc"`
	DataFimPeriodo    string `json:"dataFimPeriodoCalc"`
	DataProxLeitura   string `json:"dataProxLeitura"`
	ValorFatura       string `json:"valorFatura"`
	StatusFatura      string `json:"statusFatura"`
	NumeroFatura      string `json:"numeroFatura"`
}

type DataCertaResponse struct {
	PossuiDataBoa string `json:"possuiDataBoa"`
	DataAtual     string `json:"dataAtual"`
}

type FaturaDigitalResponse struct {
	PossuiFaturaDigital any    `json:"PossuiFaturaDigital"`
	EmailFatura         any    `json:"emailFatura"`
	EmailCadastro       string `json:"emailCadastro"`
}

type DebitoAutomaticoResponse struct {
	Retorno map[string]any `json:"retorno"`
}

type MotivosSegundaViaResponse struct {
	Motivos []struct {
		IDMotivo  string `json:"idMotivo"`
		Descricao string `json:"descricao"`
	} `json:"motivos"`
}

type DadosPagamentoResponse struct {
	CodBarras string `json:"codBarras"`
}

type FaturaPDFResponse struct {
	FileName      string `json:"fileName"`
	FileSize      string `json:"fileSize"`
	FileData      string `json:"fileData"`
	FileExtension string `json:"fileExtension"`
}

func (c *Client) doJSON(ctx context.Context, reqCtx RequestContext, method string, path string, query url.Values, payload any, out any) error {
	fullURL := c.baseURL + path
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewReader(raw)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "pt-BR")
	req.Header.Set("Origin", "https://agenciavirtual.neoenergia.com")
	req.Header.Set("Referer", "https://agenciavirtual.neoenergia.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("Authorization", "Bearer "+reqCtx.BearerToken)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return &ErrorResponse{
			StatusCode: resp.StatusCode,
			Method:     method,
			Path:       path,
			Body:       raw,
		}
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(raw, out)
}

func (c *Client) GetGrupoCliente(ctx context.Context, reqCtx RequestContext) (GrupoClienteResponse, error) {
	var out GrupoClienteResponse
	q := url.Values{"tipoPerfil": []string{"0"}}
	err := c.doJSON(ctx, reqCtx, http.MethodGet, fmt.Sprintf("/multilogin/2.0.0/agv/cliente/%s/%s/grupo-de-cliente", reqCtx.Documento, defaultDistribuidora), q, nil, &out)
	return out, err
}

func (c *Client) GetMinhaConta(ctx context.Context, reqCtx RequestContext) (MinhaContaResponse, error) {
	var out MinhaContaResponse
	q := url.Values{
		"canalSolicitante":     []string{defaultCanalSolicitante},
		"distribuidora":        []string{defaultDistribuidora},
		"usuario":              []string{reqCtx.Documento},
		"tipoPerfil":           []string{defaultTipoPerfil},
		"documentoSolicitante": []string{reqCtx.Documento},
	}
	err := c.doJSON(ctx, reqCtx, http.MethodGet, "/multilogin/2.0.0/servicos/minha-conta", q, nil, &out)
	return out, err
}

func (c *Client) GetMinhaContaLegado(ctx context.Context, reqCtx RequestContext) (MinhaContaLegadoResponse, error) {
	var out MinhaContaLegadoResponse
	q := url.Values{
		"canalSolicitante":     []string{defaultCanalSolicitante},
		"usuario":              []string{reqCtx.Documento},
		"usuarioSap":           []string{defaultUsuarioSonda},
		"usuarioSonda":         []string{defaultUsuarioSonda},
		"distribuidora":        []string{defaultDistribuidora},
		"tipoPerfil":           []string{defaultTipoPerfil},
		"documentoSolicitante": []string{reqCtx.Documento},
	}
	err := c.doJSON(ctx, reqCtx, http.MethodGet, "/multilogin/2.0.0/servicos/minha-conta/minha-conta-legado", q, nil, &out)
	return out, err
}

func (c *Client) ListUCs(ctx context.Context, reqCtx RequestContext) (UCsResponse, error) {
	var out UCsResponse
	q := url.Values{
		"documento":        []string{reqCtx.Documento},
		"canalSolicitante": []string{defaultCanalSolicitante},
		"distribuidora":    []string{defaultDistribuidora},
		"usuario":          []string{defaultUsuarioSonda},
		"indMaisUcs":       []string{"X"},
		"protocolo":        []string{"123"},
		"opcaoSSOS":        []string{"S"},
		"tipoPerfil":       []string{defaultTipoPerfil},
	}
	err := c.doJSON(ctx, reqCtx, http.MethodGet, fmt.Sprintf("/imoveis/1.1.0/clientes/%s/ucs", reqCtx.Documento), q, nil, &out)
	return out, err
}

func (c *Client) GetImovel(ctx context.Context, reqCtx RequestContext, uc string) (ImovelResponse, error) {
	var out ImovelResponse
	q := url.Values{
		"usuario":          []string{defaultUsuarioSonda},
		"canalSolicitante": []string{defaultCanalSolicitante},
		"distribuidora":    []string{defaultDistribuidora},
		"protocolo":        []string{"123"},
		"tipoPerfil":       []string{defaultTipoPerfil},
		"opcaoSSOS":        []string{"N"},
	}
	err := c.doJSON(ctx, reqCtx, http.MethodGet, fmt.Sprintf("/multilogin/2.0.0/servicos/imoveis/ucs/%s", uc), q, nil, &out)
	return out, err
}

func (c *Client) GetProtocolo(ctx context.Context, reqCtx RequestContext, codCliente string) (ProtocoloResponse, error) {
	var out ProtocoloResponse
	q := url.Values{
		"distribuidora":    []string{"COEL"},
		"canalSolicitante": []string{defaultCanalSolicitante},
		"documento":        []string{reqCtx.Documento},
		"codCliente":       []string{codCliente},
		"recaptchaAnl":     []string{"false"},
		"regiao":           []string{defaultRegiao},
	}
	err := c.doJSON(ctx, reqCtx, http.MethodGet, "/protocolo/1.1.0/obterProtocolo", q, nil, &out)
	return out, err
}

func (c *Client) ListFaturas(ctx context.Context, reqCtx RequestContext, uc string, protocolo string) (FaturasResponse, error) {
	var out FaturasResponse
	q := url.Values{
		"codigo":               []string{uc},
		"documento":            []string{reqCtx.Documento},
		"canalSolicitante":     []string{defaultCanalSolicitante},
		"usuario":              []string{defaultUsuarioSonda},
		"protocolo":            []string{protocolo},
		"tipificacao":          []string{""},
		"byPassActiv":          []string{"X"},
		"documentoSolicitante": []string{reqCtx.Documento},
		"documentoCliente":     []string{reqCtx.Documento},
		"distribuidora":        []string{defaultDistribuidora},
		"tipoPerfil":           []string{defaultTipoPerfil},
	}
	err := c.doJSON(ctx, reqCtx, http.MethodGet, "/multilogin/2.0.0/servicos/faturas/ucs/faturas", q, nil, &out)
	return out, err
}

func (c *Client) GetHistoricoConsumo(ctx context.Context, reqCtx RequestContext, uc string, protocolo string) (HistoricoConsumoResponse, error) {
	var out HistoricoConsumoResponse
	q := url.Values{
		"canalSolicitante":      []string{defaultCanalSolicitante},
		"usuario":               []string{defaultUsuarioSonda},
		"dataInicioPeriodoCalc": []string{"2021-04-18T00:00:00"},
		"protocoloSonda":        []string{protocolo},
		"opcaoSSOS":             []string{"N"},
		"protocolo":             []string{protocolo},
		"documentoSolicitante":  []string{reqCtx.Documento},
		"byPassAtiv":            []string{"X"},
		"distribuidora":         []string{defaultDistribuidora},
		"tipoPerfil":            []string{defaultTipoPerfil},
		"codigo":                []string{uc},
	}
	err := c.doJSON(ctx, reqCtx, http.MethodGet, fmt.Sprintf("/multilogin/2.0.0/servicos/historicos/ucs/%s/consumos", uc), q, nil, &out)
	return out, err
}

func (c *Client) GetDataCerta(ctx context.Context, reqCtx RequestContext, uc string) (DataCertaResponse, error) {
	var out DataCertaResponse
	q := url.Values{
		"codigo":               []string{uc},
		"canalSolicitante":     []string{defaultCanalSolicitante},
		"usuario":              []string{defaultUsuarioSonda},
		"operacao":             []string{"CON"},
		"tipoPerfil":           []string{defaultTipoPerfil},
		"documentoSolicitante": []string{""},
		"distribuidora":        []string{defaultDistribuidora},
	}
	err := c.doJSON(ctx, reqCtx, http.MethodGet, fmt.Sprintf("/multilogin/2.0.0/servicos/datacerta/ucs/%s/datacerta", uc), q, nil, &out)
	return out, err
}

func (c *Client) GetFaturaDigital(ctx context.Context, reqCtx RequestContext, uc string) (FaturaDigitalResponse, error) {
	var out FaturaDigitalResponse
	q := url.Values{
		"codigo":           []string{uc},
		"canalSolicitante": []string{defaultCanalSolicitante},
		"usuario":          []string{defaultUsuarioSonda},
		"distribuidora":    []string{defaultDistribuidora},
		"tipoPerfil":       []string{defaultTipoPerfil},
	}
	err := c.doJSON(ctx, reqCtx, http.MethodGet, "/multilogin/2.0.0/servicos/fatura-digital/ucs/fatura-digital", q, nil, &out)
	return out, err
}

func (c *Client) GetDebitoAutomatico(ctx context.Context, reqCtx RequestContext, uc string, codCliente string) (DebitoAutomaticoResponse, error) {
	var out DebitoAutomaticoResponse
	q := url.Values{
		"codigo":               []string{uc},
		"codCliente":           []string{codCliente},
		"canalSolicitante":     []string{defaultCanalSolicitante},
		"usuario":              []string{defaultUsuarioSonda},
		"valida":               []string{""},
		"distribuidora":        []string{defaultDistribuidora},
		"tipoPerfil":           []string{defaultTipoPerfil},
		"documentoSolicitante": []string{""},
	}
	err := c.doJSON(ctx, reqCtx, http.MethodGet, "/multilogin/2.0.0/servicos/debito-automatico/conta-cadastrada-debito", q, nil, &out)
	return out, err
}

func (c *Client) GetMotivosSegundaVia(ctx context.Context, reqCtx RequestContext, uc string) (MotivosSegundaViaResponse, error) {
	var out MotivosSegundaViaResponse
	q := url.Values{
		"usuario":              []string{defaultUsuarioSonda},
		"canalSolicitante":     []string{defaultCanalSolicitante},
		"distribuidora":        []string{defaultDistribuidora},
		"regiao":               []string{defaultRegiao},
		"tipoPerfil":           []string{defaultTipoPerfil},
		"documentoSolicitante": []string{reqCtx.Documento},
		"codigo":               []string{uc},
	}
	err := c.doJSON(ctx, reqCtx, http.MethodGet, "/multilogin/2.0.0/servicos/faturas/lista-motivo-segundavia", q, nil, &out)
	return out, err
}

func (c *Client) GetDadosPagamento(ctx context.Context, reqCtx RequestContext, uc string, numeroFatura string, protocolo string) (DadosPagamentoResponse, error) {
	var out DadosPagamentoResponse
	q := url.Values{
		"codigo":               []string{uc},
		"protocolo":            []string{protocolo},
		"usuario":              []string{defaultUsuarioSonda},
		"canalSolicitante":     []string{defaultCanalSolicitante},
		"distribuidora":        []string{defaultDistribuidora},
		"regiao":               []string{defaultRegiao},
		"tipoPerfil":           []string{defaultTipoPerfil},
		"byPassActiv":          []string{"X"},
		"documentoSolicitante": []string{reqCtx.Documento},
		"documento":            []string{reqCtx.Documento},
	}
	err := c.doJSON(ctx, reqCtx, http.MethodGet, fmt.Sprintf("/multilogin/2.0.0/servicos/faturas/%s/dados-pagamento", numeroFatura), q, nil, &out)
	return out, err
}

func (c *Client) GetFaturaPDF(ctx context.Context, reqCtx RequestContext, uc string, numeroFatura string, protocolo string, motivo string) (FaturaPDFResponse, error) {
	var out FaturaPDFResponse
	q := url.Values{
		"codigo":               []string{uc},
		"protocolo":            []string{protocolo},
		"tipificacao":          []string{"1031602"},
		"usuario":              []string{defaultUsuarioSonda},
		"canalSolicitante":     []string{defaultCanalSolicitante},
		"motivo":               []string{motivo},
		"distribuidora":        []string{defaultDistribuidora},
		"regiao":               []string{defaultRegiao},
		"tipoPerfil":           []string{defaultTipoPerfil},
		"documento":            []string{reqCtx.Documento},
		"documentoSolicitante": []string{reqCtx.Documento},
		"byPassActiv":          []string{""},
	}
	err := c.doJSON(ctx, reqCtx, http.MethodGet, fmt.Sprintf("/multilogin/2.0.0/servicos/faturas/%s/pdf", numeroFatura), q, nil, &out)
	return out, err
}
