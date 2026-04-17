package extractor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type Request struct {
	SchemaVersion string         `json:"schema_version"`
	JobID         string         `json:"job_id"`
	UC            string         `json:"uc"`
	Documento     string         `json:"documento"`
	NumeroFatura  string         `json:"numero_fatura"`
	MesReferencia string         `json:"mes_referencia,omitempty"`
	PDF           PDFPayload     `json:"pdf"`
	Requested     []string       `json:"requested_fields,omitempty"`
	APISnapshot   map[string]any `json:"api_snapshot,omitempty"`
}

type PDFPayload struct {
	Mode     string `json:"mode"`
	Base64   string `json:"base64,omitempty"`
	Path     string `json:"path,omitempty"`
	FileName string `json:"file_name,omitempty"`
}

type Response struct {
	SchemaVersion string             `json:"schema_version"`
	JobID         string             `json:"job_id"`
	Status        string             `json:"status"`
	Fields        map[string]any     `json:"fields"`
	SourceMap     map[string]string  `json:"source_map"`
	ConfidenceMap map[string]float64 `json:"confidence_map"`
	Warnings      []string           `json:"warnings"`
	Artifacts     map[string]any     `json:"artifacts"`
}

func (c *Client) Extract(ctx context.Context, in Request) (Response, error) {
	raw, err := json.Marshal(in)
	if err != nil {
		return Response{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/extract", bytes.NewReader(raw))
	if err != nil {
		return Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return Response{}, fmt.Errorf("extractor returned status %d", resp.StatusCode)
	}
	var out Response
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return Response{}, err
	}
	return out, nil
}
