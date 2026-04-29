package integration

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// Handler expõe endpoints HTTP para o domínio integration.
type Handler struct {
	svc    *Service
	logger *slog.Logger
}

// NewHandler cria um novo Handler.
func NewHandler(svc *Service, logger *slog.Logger) *Handler {
	return &Handler{svc: svc, logger: logger}
}

// RegisterRoutes registra todos os endpoints do domínio integration no mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Credentials
	mux.HandleFunc("/v1/integration/credentials", h.handleCredentials)
	mux.HandleFunc("/v1/integration/credentials/", h.handleCredentialByID)

	// Consumer Units
	mux.HandleFunc("/v1/integration/consumer-units", h.ListConsumerUnits)
	mux.HandleFunc("/v1/integration/consumer-units/", h.handleConsumerUnitByUC)

	// Sync Runs
	mux.HandleFunc("/v1/integration/sync-runs", h.CreateSyncRun)
	mux.HandleFunc("/v1/integration/sync-runs/", h.handleSyncRunByID)

	// Jobs
	mux.HandleFunc("/v1/integration/jobs", h.EnqueueJob)
}

// --- Credentials ---

func (h *Handler) handleCredentials(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateCredential(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

func (h *Handler) handleCredentialByID(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path)
	if len(parts) != 4 || parts[0] != "v1" || parts[1] != "integration" || parts[2] != "credentials" {
		writeError(w, http.StatusNotFound, "not_found")
		return
	}
	r = r.WithContext(r.Context())
	// Hack: pass id via query since we can't use PathValue with splitPath pattern
	q := r.URL.Query()
	q.Set("__id", parts[3])
	r.URL.RawQuery = q.Encode()

	switch r.Method {
	case http.MethodGet:
		h.GetCredential(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

// CreateCredentialRequest é o payload para criar credencial.
type CreateCredentialRequest struct {
	Label           string `json:"label"`
	DocumentoCipher string `json:"documento_cipher"`
	DocumentoNonce  string `json:"documento_nonce"`
	SenhaCipher     string `json:"senha_cipher"`
	SenhaNonce      string `json:"senha_nonce"`
	UF              string `json:"uf"`
	TipoAcesso      string `json:"tipo_acesso"`
	KeyVersion      string `json:"key_version"`
}

// CreateCredential cria uma nova credencial.
func (h *Handler) CreateCredential(w http.ResponseWriter, r *http.Request) {
	var req CreateCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}
	if req.Label == "" || req.DocumentoCipher == "" || req.SenhaCipher == "" {
		writeError(w, http.StatusBadRequest, "label, documento_cipher e senha_cipher são obrigatórios")
		return
	}

	c := &Credential{
		Label:           req.Label,
		DocumentoCipher: req.DocumentoCipher,
		DocumentoNonce:  req.DocumentoNonce,
		SenhaCipher:     req.SenhaCipher,
		SenhaNonce:      req.SenhaNonce,
		UF:              req.UF,
		TipoAcesso:      req.TipoAcesso,
		KeyVersion:      req.KeyVersion,
	}
	created, err := h.svc.CreateCredential(r.Context(), c)
	if err != nil {
		h.logger.Error("create_credential_failed", "error", err)
		writeError(w, http.StatusInternalServerError, "create_credential_failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":         created.ID.String(),
		"label":      created.Label,
		"uf":         created.UF,
		"created_at": created.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// GetCredential retorna uma credencial por ID (sem dados sensíveis).
func (h *Handler) GetCredential(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("__id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id é obrigatório")
		return
	}
	c, err := h.svc.GetCredential(r.Context(), id)
	if err != nil {
		h.logger.Error("get_credential_failed", "error", err, "id", id)
		writeError(w, http.StatusNotFound, "not_found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":          c.ID.String(),
		"label":       c.Label,
		"uf":          c.UF,
		"tipo_acesso": c.TipoAcesso,
		"created_at":  c.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updated_at":  c.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// --- Consumer Units ---

func (h *Handler) handleConsumerUnitByUC(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path)
	if len(parts) != 4 || parts[0] != "v1" || parts[1] != "integration" || parts[2] != "consumer-units" {
		writeError(w, http.StatusNotFound, "not_found")
		return
	}
	q := r.URL.Query()
	q.Set("__uc", parts[3])
	r.URL.RawQuery = q.Encode()

	switch r.Method {
	case http.MethodGet:
		h.GetConsumerUnit(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

// ListConsumerUnits lista UCs descobertas.
func (h *Handler) ListConsumerUnits(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}
	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	status := r.URL.Query().Get("status")

	units, err := h.svc.ListConsumerUnits(r.Context(), limit, status)
	if err != nil {
		h.logger.Error("list_consumer_units_failed", "error", err)
		writeError(w, http.StatusInternalServerError, "list_consumer_units_failed")
		return
	}

	items := make([]map[string]any, len(units))
	for i, u := range units {
		items[i] = toConsumerUnitResponse(&u)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"items": items,
		"count": len(items),
	})
}

// GetConsumerUnit retorna uma UC por código.
func (h *Handler) GetConsumerUnit(w http.ResponseWriter, r *http.Request) {
	uc := r.URL.Query().Get("__uc")
	if uc == "" {
		writeError(w, http.StatusBadRequest, "uc é obrigatório")
		return
	}
	u, err := h.svc.GetConsumerUnit(r.Context(), uc)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toConsumerUnitResponse(u))
}

// --- Sync Runs ---

func (h *Handler) handleSyncRunByID(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path)
	if len(parts) != 4 || parts[0] != "v1" || parts[1] != "integration" || parts[2] != "sync-runs" {
		writeError(w, http.StatusNotFound, "not_found")
		return
	}
	q := r.URL.Query()
	q.Set("__id", parts[3])
	r.URL.RawQuery = q.Encode()

	switch r.Method {
	case http.MethodGet:
		h.GetSyncRun(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

// CreateSyncRun inicia um novo sync run.
func (h *Handler) CreateSyncRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}
	var req struct {
		CredentialID string `json:"credential_id,omitempty"`
		Documento    string `json:"documento"`
		UC           string `json:"uc"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}
	if req.Documento == "" || req.UC == "" {
		writeError(w, http.StatusBadRequest, "documento e uc são obrigatórios")
		return
	}

	sr := &SyncRun{
		Documento: req.Documento,
		UC:        req.UC,
		Status:    "pending",
	}
	if req.CredentialID != "" {
		uid, err := uuid.Parse(req.CredentialID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "credential_id inválido")
			return
		}
		sr.CredentialID = &uid
	}

	created, err := h.svc.RecordSyncRun(r.Context(), sr)
	if err != nil {
		h.logger.Error("create_sync_run_failed", "error", err)
		writeError(w, http.StatusInternalServerError, "create_sync_run_failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":     created.ID.String(),
		"status": created.Status,
		"uc":     created.UC,
	})
}

// GetSyncRun retorna um sync run por ID.
func (h *Handler) GetSyncRun(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("__id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id é obrigatório")
		return
	}
	sr, err := h.svc.GetSyncRun(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toSyncRunResponse(sr))
}

// --- Jobs ---

// EnqueueJob adiciona um job à fila.
func (h *Handler) EnqueueJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}
	var req struct {
		JobType string         `json:"job_type"`
		Payload map[string]any `json:"payload"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}
	if req.JobType == "" {
		writeError(w, http.StatusBadRequest, "job_type é obrigatório")
		return
	}

	job, err := h.svc.EnqueueJob(r.Context(), req.JobType, req.Payload)
	if err != nil {
		h.logger.Error("enqueue_job_failed", "error", err)
		writeError(w, http.StatusInternalServerError, "enqueue_job_failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":         job.ID.String(),
		"job_type":   job.JobType,
		"status":     job.Status,
		"created_at": job.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// --- Helpers ---

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func splitPath(path string) []string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "/")
}

func toSyncRunResponse(sr *SyncRun) map[string]any {
	credID := ""
	if sr.CredentialID != nil {
		credID = sr.CredentialID.String()
	}
	step := ""
	if sr.Step != nil {
		step = *sr.Step
	}
	errMsg := ""
	if sr.ErrorMessage != nil {
		errMsg = *sr.ErrorMessage
	}
	return map[string]any{
		"id":            sr.ID.String(),
		"credential_id": credID,
		"documento":     sr.Documento,
		"uc":            sr.UC,
		"status":        sr.Status,
		"step":          step,
		"error_message": errMsg,
		"created_at":    sr.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toConsumerUnitResponse(u *ConsumerUnit) map[string]any {
	credID := ""
	if u.CredentialID != nil {
		credID = u.CredentialID.String()
	}
	status := ""
	if u.Status != nil {
		status = *u.Status
	}
	nome := ""
	if u.NomeCliente != nil {
		nome = *u.NomeCliente
	}
	inst := ""
	if u.Instalacao != nil {
		inst = *u.Instalacao
	}
	contrato := ""
	if u.Contrato != nil {
		contrato = *u.Contrato
	}
	grupo := ""
	if u.GrupoTensao != nil {
		grupo = *u.GrupoTensao
	}
	return map[string]any{
		"uc":            u.UC,
		"credential_id": credID,
		"status":        status,
		"nome_cliente":  nome,
		"instalacao":    inst,
		"contrato":      contrato,
		"grupo_tensao":  grupo,
		"endereco":      u.Endereco,
		"imovel":        u.Imovel,
		"created_at":    u.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updated_at":    u.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
