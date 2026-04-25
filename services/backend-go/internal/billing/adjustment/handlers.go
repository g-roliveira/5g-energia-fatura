package adjustment

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// Handler expõe endpoints HTTP para ajustes manuais.
type Handler struct {
	svc    *Service
	logger *slog.Logger
}

// NewHandler cria um novo Handler.
func NewHandler(svc *Service, logger *slog.Logger) *Handler {
	return &Handler{svc: svc, logger: logger}
}

// RegisterRoutes registra as rotas de ajuste no mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/v1/billing/calculations/", h.handleCalculation)
}

func splitPath(path string) []string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "/")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func (h *Handler) handleCalculation(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path)
	if len(parts) < 4 || parts[0] != "v1" || parts[1] != "billing" || parts[2] != "calculations" {
		writeError(w, http.StatusNotFound, "not_found")
		return
	}
	calcID, err := uuid.Parse(parts[3])
	if err != nil {
		writeError(w, http.StatusBadRequest, "id inválido")
		return
	}

	// /v1/billing/calculations/{id}/adjust
	if len(parts) == 5 && parts[4] == "adjust" {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
			return
		}
		h.applyAdjustment(w, r, calcID)
		return
	}

	// /v1/billing/calculations/{id}/adjustments
	if len(parts) == 5 && parts[4] == "adjustments" {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
			return
		}
		h.listAdjustments(w, r, calcID)
		return
	}

	writeError(w, http.StatusNotFound, "not_found")
}

func (h *Handler) applyAdjustment(w http.ResponseWriter, r *http.Request, calcID uuid.UUID) {
	var req ApplyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}
	req.CalculationID = calcID
	if req.FieldPath == "" || req.Reason == "" {
		writeError(w, http.StatusBadRequest, "field_path e reason são obrigatórios")
		return
	}

	adj, err := h.svc.Apply(r.Context(), req)
	if err != nil {
		h.logger.Error("apply_adjustment_failed", "error", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, adj)
}

func (h *Handler) listAdjustments(w http.ResponseWriter, r *http.Request, calcID uuid.UUID) {
	adjs, err := h.svc.List(r.Context(), calcID)
	if err != nil {
		h.logger.Error("list_adjustments_failed", "error", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items": adjs,
		"count": len(adjs),
	})
}
