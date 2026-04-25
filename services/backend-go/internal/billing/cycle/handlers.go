package cycle

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// Handler expõe endpoints HTTP para ciclos de faturamento.
type Handler struct {
	svc    *Service
	logger *slog.Logger
}

// NewHandler cria um novo Handler.
func NewHandler(svc *Service, logger *slog.Logger) *Handler {
	return &Handler{svc: svc, logger: logger}
}

// RegisterRoutes registra as rotas de ciclo no mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/v1/billing/cycles", h.handleCycles)
	mux.HandleFunc("/v1/billing/cycles/", h.handleCycleDetail)
	mux.HandleFunc("/v1/billing/events/cycles/", h.handleSSE)
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

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// --- /v1/billing/cycles ---

func (h *Handler) handleCycles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createCycle(w, r)
	case http.MethodGet:
		h.listCycles(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

func (h *Handler) createCycle(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}
	if req.Year == 0 || req.Month == 0 {
		writeError(w, http.StatusBadRequest, "year e month são obrigatórios")
		return
	}

	cycle, err := h.svc.Create(r.Context(), req)
	if err != nil {
		h.logger.Error("create_cycle_failed", "error", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, cycle)
}

func (h *Handler) listCycles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	req := ListCyclesRequest{Limit: 50}
	if y := q.Get("year"); y != "" {
		v, _ := strconv.ParseInt(y, 10, 16)
		req.Year = int16(v)
	}
	req.Status = q.Get("status")
	if l := q.Get("limit"); l != "" {
		v, _ := strconv.Atoi(l)
		req.Limit = v
	}
	if o := q.Get("offset"); o != "" {
		v, _ := strconv.Atoi(o)
		req.Offset = v
	}

	cycles, err := h.svc.List(r.Context(), req)
	if err != nil {
		h.logger.Error("list_cycles_failed", "error", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items":  cycles,
		"count":  len(cycles),
		"limit":  req.Limit,
		"offset": req.Offset,
	})
}

// --- /v1/billing/cycles/{id} / {id}/rows / {id}/close ---

func (h *Handler) handleCycleDetail(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path)
	if len(parts) < 4 || parts[0] != "v1" || parts[1] != "billing" || parts[2] != "cycles" {
		writeError(w, http.StatusNotFound, "not_found")
		return
	}
	id, err := parseUUID(parts[3])
	if err != nil {
		writeError(w, http.StatusBadRequest, "id inválido")
		return
	}

	// /v1/billing/cycles/{id}/rows
	if len(parts) == 5 && parts[4] == "rows" {
		switch r.Method {
		case http.MethodGet:
			h.getCycleRows(w, r, id)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		}
		return
	}

	// /v1/billing/cycles/{id}/close
	if len(parts) == 5 && parts[4] == "close" {
		switch r.Method {
		case http.MethodPost:
			h.closeCycle(w, r, id)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		}
		return
	}

	// /v1/billing/cycles/{id}
	if len(parts) == 4 {
		switch r.Method {
		case http.MethodGet:
			h.getCycle(w, r, id)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		}
		return
	}

	writeError(w, http.StatusNotFound, "not_found")
}

func (h *Handler) getCycle(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	cycle, err := h.svc.Get(r.Context(), id)
	if err != nil {
		if err.Error() == "not_found" {
			writeError(w, http.StatusNotFound, "not_found")
			return
		}
		h.logger.Error("get_cycle_failed", "error", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cycle)
}

func (h *Handler) getCycleRows(w http.ResponseWriter, r *http.Request, cycleID uuid.UUID) {
	q := r.URL.Query()
	req := ListRowsRequest{CycleID: cycleID, Limit: 100}
	req.Q = q.Get("q")
	req.SyncStatus = q.Get("sync_status")
	req.CalcStatus = q.Get("calc_status")
	req.NeedsReviewOnly = q.Get("needs_review_only") == "true"
	if l := q.Get("limit"); l != "" {
		v, _ := strconv.Atoi(l)
		req.Limit = v
	}
	if o := q.Get("offset"); o != "" {
		v, _ := strconv.Atoi(o)
		req.Offset = v
	}

	rows, err := h.svc.ListRows(r.Context(), req)
	if err != nil {
		h.logger.Error("list_cycle_rows_failed", "error", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items":   rows,
		"count":   len(rows),
		"cycle_id": cycleID.String(),
	})
}

func (h *Handler) closeCycle(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	var req CloseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}
	if err := h.svc.Close(r.Context(), id, req); err != nil {
		h.logger.Error("close_cycle_failed", "error", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "closed"})
}

// --- SSE ---

func (h *Handler) handleSSE(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}
	parts := splitPath(r.URL.Path)
	if len(parts) != 5 || parts[0] != "v1" || parts[1] != "billing" || parts[2] != "events" || parts[3] != "cycles" {
		writeError(w, http.StatusNotFound, "not_found")
		return
	}
	cycleID, err := parseUUID(parts[4])
	if err != nil {
		writeError(w, http.StatusBadRequest, "id inválido")
		return
	}

	// Simples SSE com polling (LISTEN/NOTIFY vem depois)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "streaming não suportado")
		return
	}

	// Enviar evento inicial
	fmt.Fprintf(w, "event: connected\ndata: %s\n\n", fmt.Sprintf(`{"cycle_id":"%s"}`, cycleID))
	flusher.Flush()

	// TODO: implementar LISTEN/NOTIFY realtime
	// Por enquanto, o cliente deve reconectar para obter atualizações
	fmt.Fprintf(w, "event: done\ndata: {}\n\n")
	flusher.Flush()
}
