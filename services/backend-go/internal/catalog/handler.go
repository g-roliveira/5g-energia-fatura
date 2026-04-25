package catalog

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// RegisterHandlers registra as rotas de catalog no mux.
func RegisterHandlers(mux *http.ServeMux, svc *Service, logger *slog.Logger) {
	// Customers
	mux.HandleFunc("/v1/catalog/customers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleCreateCustomer(w, r, svc, logger)
		case http.MethodGet:
			handleListCustomers(w, r, svc, logger)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/v1/catalog/customers/", func(w http.ResponseWriter, r *http.Request) {
		parts := splitPath(r.URL.Path)
		if len(parts) != 4 || parts[0] != "v1" || parts[1] != "catalog" || parts[2] != "customers" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		id, err := uuid.Parse(parts[3])
		if err != nil {
			writeClientError(w, http.StatusBadRequest, "invalid id")
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetCustomer(w, r, svc, logger, id)
		case http.MethodPatch:
			handleUpdateCustomer(w, r, svc, logger, id)
		case http.MethodDelete:
			handleArchiveCustomer(w, r, svc, logger, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Consumer Units
	mux.HandleFunc("/v1/catalog/units", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleCreateUnit(w, r, svc, logger)
		case http.MethodGet:
			handleListUnits(w, r, svc, logger)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/v1/catalog/units/", func(w http.ResponseWriter, r *http.Request) {
		parts := splitPath(r.URL.Path)
		if len(parts) < 4 || parts[0] != "v1" || parts[1] != "catalog" || parts[2] != "units" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if len(parts) == 4 {
			id, err := uuid.Parse(parts[3])
			if err != nil {
				writeClientError(w, http.StatusBadRequest, "invalid id")
				return
			}
			if r.Method == http.MethodGet {
				handleGetUnit(w, r, svc, logger, id)
				return
			}
		}

		if len(parts) == 5 && parts[4] == "link" {
			if r.Method == http.MethodPost {
				handleLinkUnit(w, r, svc, logger, parts[3])
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
	})

	// Contracts
	mux.HandleFunc("/v1/catalog/contracts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handleCreateContract(w, r, svc, logger)
	})

	mux.HandleFunc("/v1/catalog/contracts/", func(w http.ResponseWriter, r *http.Request) {
		parts := splitPath(r.URL.Path)
		if len(parts) != 4 || parts[0] != "v1" || parts[1] != "catalog" || parts[2] != "contracts" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		id, err := uuid.Parse(parts[3])
		if err != nil {
			writeClientError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handleGetContract(w, r, svc, logger, id)
	})

	mux.HandleFunc("/v1/catalog/units/", func(w http.ResponseWriter, r *http.Request) {
		parts := splitPath(r.URL.Path)
		// /v1/catalog/units/{ucCode}/active-contract
		if len(parts) != 5 || parts[0] != "v1" || parts[1] != "catalog" || parts[2] != "units" || parts[4] != "active-contract" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		ucCode := parts[3]
		handleGetActiveContract(w, r, svc, logger, ucCode)
	})
}

// --- Customer handlers ---

func handleCreateCustomer(w http.ResponseWriter, r *http.Request, svc *Service, logger *slog.Logger) {
	var input CustomerInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeClientError(w, http.StatusBadRequest, "invalid json")
		return
	}

	c, err := svc.CreateCustomer(r.Context(), input)
	if err != nil {
		writeClientError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, customerView(c))
}

func handleGetCustomer(w http.ResponseWriter, r *http.Request, svc *Service, logger *slog.Logger, id uuid.UUID) {
	c, err := svc.GetCustomer(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, customerView(c))
}

func handleListCustomers(w http.ResponseWriter, r *http.Request, svc *Service, logger *slog.Logger) {
	q := r.URL.Query()
	filter := CustomerFilter{Limit: 50}

	if s := q.Get("status"); s != "" {
		filter.Status = &s
	}
	if s := q.Get("q"); s != "" {
		filter.Query = &s
	}
	if c := q.Get("cursor"); c != "" {
		filter.Cursor = &c
	}
	if l := q.Get("limit"); l != "" {
		// parse limit
	}

	customers, nextCursor, err := svc.ListCustomers(r.Context(), filter)
	if err != nil {
		writeInternalError(w, logger, "list_customers", err)
		return
	}

	items := make([]map[string]any, len(customers))
	for i, c := range customers {
		items[i] = customerView(&c)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items":       items,
		"has_more":    nextCursor != "",
		"next_cursor": nextCursor,
	})
}

func handleUpdateCustomer(w http.ResponseWriter, r *http.Request, svc *Service, logger *slog.Logger, id uuid.UUID) {
	var patch CustomerPatch
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeClientError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if err := svc.UpdateCustomer(r.Context(), id, patch); err != nil {
		writeInternalError(w, logger, "update_customer", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func handleArchiveCustomer(w http.ResponseWriter, r *http.Request, svc *Service, logger *slog.Logger, id uuid.UUID) {
	if err := svc.ArchiveCustomer(r.Context(), id); err != nil {
		writeInternalError(w, logger, "archive_customer", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Unit handlers ---

func handleCreateUnit(w http.ResponseWriter, r *http.Request, svc *Service, logger *slog.Logger) {
	var input UnitInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeClientError(w, http.StatusBadRequest, "invalid json")
		return
	}

	u, err := svc.CreateUnit(r.Context(), input)
	if err != nil {
		writeClientError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, unitView(u))
}

func handleGetUnit(w http.ResponseWriter, r *http.Request, svc *Service, logger *slog.Logger, id uuid.UUID) {
	u, err := svc.GetUnit(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, unitView(u))
}

func handleListUnits(w http.ResponseWriter, r *http.Request, svc *Service, logger *slog.Logger) {
	q := r.URL.Query()
	filter := UnitFilter{Limit: 50}

	if cid := q.Get("customer_id"); cid != "" {
		if id, err := uuid.Parse(cid); err == nil {
			filter.CustomerID = &id
		}
	}
	if q.Get("active_only") == "true" {
		filter.ActiveOnly = true
	}

	units, nextCursor, err := svc.ListUnits(r.Context(), filter)
	if err != nil {
		writeInternalError(w, logger, "list_units", err)
		return
	}

	items := make([]map[string]any, len(units))
	for i, u := range units {
		items[i] = unitView(&u)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items":       items,
		"has_more":    nextCursor != "",
		"next_cursor": nextCursor,
	})
}

func handleLinkUnit(w http.ResponseWriter, r *http.Request, svc *Service, logger *slog.Logger, unitIDStr string) {
	unitID, err := uuid.Parse(unitIDStr)
	if err != nil {
		writeClientError(w, http.StatusBadRequest, "invalid unit id")
		return
	}

	var payload struct {
		CustomerID string `json:"customer_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeClientError(w, http.StatusBadRequest, "invalid json")
		return
	}

	customerID, err := uuid.Parse(payload.CustomerID)
	if err != nil {
		writeClientError(w, http.StatusBadRequest, "invalid customer id")
		return
	}

	if err := svc.LinkUnitToCustomer(r.Context(), unitID, customerID); err != nil {
		writeInternalError(w, logger, "link_unit", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Contract handlers ---

func handleCreateContract(w http.ResponseWriter, r *http.Request, svc *Service, logger *slog.Logger) {
	var input ContractInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeClientError(w, http.StatusBadRequest, "invalid json")
		return
	}

	c, err := svc.CreateContract(r.Context(), input)
	if err != nil {
		writeClientError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, contractView(c))
}

func handleGetContract(w http.ResponseWriter, r *http.Request, svc *Service, logger *slog.Logger, id uuid.UUID) {
	c, err := svc.GetContract(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, contractView(c))
}

func handleGetActiveContract(w http.ResponseWriter, r *http.Request, svc *Service, logger *slog.Logger, ucCode string) {
	// Primeiro busca a UC pelo código
	// Depois busca o contrato ativo
	// Simplificação: este handler precisa de GetUnitByCode no service
	w.WriteHeader(http.StatusNotImplemented)
}

// --- View mappers ---

func customerView(c *Customer) map[string]any {
	v := map[string]any{
		"id":           c.ID.String(),
		"tipo_pessoa":  c.TipoPessoa,
		"nome_razao":   c.NomeRazao,
		"cpf_cnpj":     c.CPFCNPJ,
		"status":       c.Status,
		"tipo_cliente": c.TipoCliente,
		"created_at":   c.CreatedAt,
		"updated_at":   c.UpdatedAt,
	}
	if c.NomeFantasia != nil {
		v["nome_fantasia"] = *c.NomeFantasia
	}
	if c.Email != nil {
		v["email"] = *c.Email
	}
	if c.Telefone != nil {
		v["telefone"] = *c.Telefone
	}
	if c.Notes != nil {
		v["notes"] = *c.Notes
	}
	if c.ArchivedAt != nil {
		v["archived_at"] = *c.ArchivedAt
	}
	return v
}

func unitView(u *ConsumerUnit) map[string]any {
	v := map[string]any{
		"id":            u.ID.String(),
		"customer_id":   u.CustomerID.String(),
		"uc_code":       u.UCCode,
		"ativa":         u.Ativa,
		"created_at":    u.CreatedAt,
		"updated_at":    u.UpdatedAt,
	}
	if u.Distribuidora != nil {
		v["distribuidora"] = *u.Distribuidora
	}
	if u.Apelido != nil {
		v["apelido"] = *u.Apelido
	}
	if u.ClasseConsumo != nil {
		v["classe_consumo"] = *u.ClasseConsumo
	}
	if u.CredentialID != nil {
		v["credential_id"] = *u.CredentialID
	}
	return v
}

func contractView(c *Contract) map[string]any {
	v := map[string]any{
		"id":                                   c.ID.String(),
		"customer_id":                          c.CustomerID.String(),
		"consumer_unit_id":                     c.ConsumerUnitID.String(),
		"vigencia_inicio":                      c.VigenciaInicio.Format("2006-01-02"),
		"desconto_percentual":                  c.DescontoPercentual,
		"ip_faturamento_mode":                  c.IPFaturamentoMode,
		"ip_faturamento_valor":                 c.IPFaturamentoValor,
		"ip_faturamento_percent":               c.IPFaturamentoPercent,
		"bandeira_com_desconto":                c.BandeiraComDesconto,
		"custo_disponibilidade_sempre_cobrado": c.CustoDisponibilidadeSempreCobrado,
		"status":                               c.Status,
		"created_at":                           c.CreatedAt,
	}
	if c.VigenciaFim != nil {
		v["vigencia_fim"] = c.VigenciaFim.Format("2006-01-02")
	}
	if c.Notes != nil {
		v["notes"] = *c.Notes
	}
	if c.CreatedBy != nil {
		v["created_by"] = c.CreatedBy.String()
	}
	return v
}

// --- HTTP helpers (copiados de app package para evitar circular dependency) ---

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeClientError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeInternalError(w http.ResponseWriter, logger *slog.Logger, op string, err error) {
	logger.Error(op, "error", err)
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal"})
}

func splitPath(path string) []string {
	path = path
	for len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	for len(path) > 0 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	if path == "" {
		return nil
	}
	var parts []string
	start := 0
	for i := 0; i <= len(path); i++ {
		if i == len(path) || path[i] == '/' {
			if i > start {
				parts = append(parts, path[start:i])
			}
			start = i + 1
		}
	}
	return parts
}
