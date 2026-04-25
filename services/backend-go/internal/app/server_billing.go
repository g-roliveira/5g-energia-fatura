package app

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/billing/adjustment"
	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/billing/contract"
	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/billing/cycle"
	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/billing/repo"
)

// BillingDeps holds the Postgres-backed dependencies that the billing
// handlers need. Wired up in NewServer when the Postgres pool is present.
type BillingDeps struct {
	Pool        *pgxpool.Pool
	Contract    *contract.Service
	Calculation *repo.CalculationRepo
	Cycle       *cycle.Service
	Adjustment  *adjustment.Service
}

// NewBillingDeps wires the repos and services from a pgx pool.
// Call this from NewServer after pgstore.Open succeeds.
func NewBillingDeps(pool *pgxpool.Pool) *BillingDeps {
	contractRepo := repo.NewContractRepo(pool)
	calcRepo := repo.NewCalculationRepo(pool)
	return &BillingDeps{
		Pool:        pool,
		Contract:    contract.NewService(contractRepo),
		Calculation: calcRepo,
		Cycle:       cycle.NewService(pool),
		Adjustment:  adjustment.NewService(pool),
	}
}

// RegisterBillingRoutes attaches the billing HTTP routes to the existing
// mux and route catalog. Call this from NewServer.
//
// Integration — in NewServer() you add:
//
//	if cfg.BackofficePGURL != "" {
//	    pool, err := pgstore.Open(context.Background(), pgstore.LoadConfigFromEnv())
//	    if err != nil {
//	        return nil, err
//	    }
//	    deps := NewBillingDeps(pool)
//	    RegisterBillingRoutes(mux, docs, deps, logger)
//	}
//
// The routes show up automatically in /docs and /openapi.json because
// we use the same routeCatalog + mux from server.go.
func RegisterBillingRoutes(
	mux *http.ServeMux,
	docs *routeCatalog,
	deps *BillingDeps,
	logger *slog.Logger,
) {
	// --- CYCLES ------------------------------------------------------
	cycleHandler := cycle.NewHandler(deps.Cycle, logger)
	cycleHandler.RegisterRoutes(mux)


	// --- CONTRACTS ---------------------------------------------------

	docs.add(http.MethodPost, "/v1/billing/contracts", "Create contract (nova versão fecha anterior)", []string{"billing", "contracts"}, http.StatusCreated)
	docs.add(http.MethodGet, "/v1/billing/contracts/{id}", "Get contract by id", []string{"billing", "contracts"}, http.StatusOK)
	mux.HandleFunc("/v1/billing/contracts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handleContractCreate(w, r, deps, logger)
	})
	mux.HandleFunc("/v1/billing/contracts/", func(w http.ResponseWriter, r *http.Request) {
		parts := splitPath(r.URL.Path)
		if len(parts) != 4 || parts[0] != "v1" || parts[1] != "billing" || parts[2] != "contracts" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		id, err := uuid.Parse(parts[3])
		if err != nil {
			writeClientError(w, http.StatusBadRequest, "id inválido")
			return
		}
		c, err := deps.Contract.Get(r.Context(), id)
		if errors.Is(err, repo.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			writeInternalError(w, logger, "contract_get", err)
			return
		}
		writeJSON(w, http.StatusOK, contractView(c))
	})

	// --- ACTIVE CONTRACT FOR UC --------------------------------------

	docs.add(http.MethodGet, "/v1/billing/consumer-units/{uc_id}/active-contract", "Get active contract for UC", []string{"billing", "contracts"}, http.StatusOK)
	mux.HandleFunc("/v1/billing/consumer-units/", func(w http.ResponseWriter, r *http.Request) {
		parts := splitPath(r.URL.Path)
		// /v1/billing/consumer-units/{id}/active-contract
		if len(parts) != 5 || parts[0] != "v1" || parts[1] != "billing" ||
			parts[2] != "consumer-units" || parts[4] != "active-contract" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		ucID, err := uuid.Parse(parts[3])
		if err != nil {
			writeClientError(w, http.StatusBadRequest, "uc_id inválido")
			return
		}
		c, err := deps.Contract.GetActiveForUC(r.Context(), ucID)
		if errors.Is(err, repo.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			writeInternalError(w, logger, "contract_active", err)
			return
		}
		writeJSON(w, http.StatusOK, contractView(c))
	})

	// --- CALCULATION READ --------------------------------------------

	docs.add(http.MethodGet, "/v1/billing/calculations/{id}", "Get billing calculation", []string{"billing", "calculations"}, http.StatusOK)
	mux.HandleFunc("/v1/billing/calculations/", func(w http.ResponseWriter, r *http.Request) {
		parts := splitPath(r.URL.Path)
		if len(parts) < 4 || parts[0] != "v1" || parts[1] != "billing" || parts[2] != "calculations" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		id, err := uuid.Parse(parts[3])
		if err != nil {
			writeClientError(w, http.StatusBadRequest, "id inválido")
			return
		}

		// /v1/billing/calculations/{id}/adjust
		if len(parts) == 5 && parts[4] == "adjust" {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			var req adjustment.ApplyRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeClientError(w, http.StatusBadRequest, "invalid_json")
				return
			}
			req.CalculationID = id
			adj, err := deps.Adjustment.Apply(r.Context(), req)
			if err != nil {
				writeInternalError(w, logger, "adjustment_apply", err)
				return
			}
			writeJSON(w, http.StatusCreated, adj)
			return
		}

		// /v1/billing/calculations/{id}/adjustments
		if len(parts) == 5 && parts[4] == "adjustments" {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			adjs, err := deps.Adjustment.List(r.Context(), id)
			if err != nil {
				writeInternalError(w, logger, "adjustment_list", err)
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"items": adjs, "count": len(adjs)})
			return
		}

		// /v1/billing/calculations/{id}
		if len(parts) == 4 {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			c, err := deps.Calculation.GetByID(r.Context(), id)
			if errors.Is(err, repo.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if err != nil {
				writeInternalError(w, logger, "calculation_get", err)
				return
			}
			writeJSON(w, http.StatusOK, calcView(c))
			return
		}

		w.WriteHeader(http.StatusNotFound)
	})
}

// -----------------------------------------------------------------
// HTTP handler internals
// -----------------------------------------------------------------

type createContractBody struct {
	CustomerID                        string `json:"customer_id"`
	ConsumerUnitID                    string `json:"consumer_unit_id"`
	VigenciaInicio                    string `json:"vigencia_inicio"` // YYYY-MM-DD
	DescontoPercentual                string `json:"desconto_percentual"`
	IPFaturamentoMode                 string `json:"ip_faturamento_mode"`
	IPFaturamentoValor                string `json:"ip_faturamento_valor"`
	IPFaturamentoPercent              string `json:"ip_faturamento_percent"`
	BandeiraComDesconto               bool   `json:"bandeira_com_desconto"`
	CustoDisponibilidadeSempreCobrado bool   `json:"custo_disponibilidade_sempre_cobrado"`
	Notes                             string `json:"notes"`
	CreatedBy                         string `json:"created_by"`
}

func handleContractCreate(w http.ResponseWriter, r *http.Request, deps *BillingDeps, logger *slog.Logger) {
	var body createContractBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeClientError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	customerID, err := uuid.Parse(body.CustomerID)
	if err != nil {
		writeClientError(w, http.StatusBadRequest, "customer_id inválido")
		return
	}
	ucID, err := uuid.Parse(body.ConsumerUnitID)
	if err != nil {
		writeClientError(w, http.StatusBadRequest, "consumer_unit_id inválido")
		return
	}
	vig, err := parseDate(body.VigenciaInicio)
	if err != nil {
		writeClientError(w, http.StatusBadRequest, "vigencia_inicio deve estar em YYYY-MM-DD")
		return
	}
	desc, err := decimal.NewFromString(body.DescontoPercentual)
	if err != nil {
		writeClientError(w, http.StatusBadRequest, "desconto_percentual inválido")
		return
	}

	in := contract.CreateInput{
		CustomerID:                        customerID,
		ConsumerUnitID:                    ucID,
		VigenciaInicio:                    vig,
		DescontoPercentual:                desc,
		IPFaturamentoMode:                 repo.IPMode(body.IPFaturamentoMode),
		BandeiraComDesconto:               body.BandeiraComDesconto,
		CustoDisponibilidadeSempreCobrado: body.CustoDisponibilidadeSempreCobrado,
	}
	if body.IPFaturamentoValor != "" {
		if v, err := decimal.NewFromString(body.IPFaturamentoValor); err == nil {
			in.IPFaturamentoValor = v
		}
	}
	if body.IPFaturamentoPercent != "" {
		if v, err := decimal.NewFromString(body.IPFaturamentoPercent); err == nil {
			in.IPFaturamentoPercent = v
		}
	}
	if body.Notes != "" {
		in.Notes = &body.Notes
	}
	if body.CreatedBy != "" {
		if u, err := uuid.Parse(body.CreatedBy); err == nil {
			in.CreatedBy = &u
		}
	}

	c, err := deps.Contract.Create(r.Context(), in)
	if err != nil {
		// Validation errors come back as plain strings
		writeClientError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, contractView(c))
}

// -----------------------------------------------------------------
// View mappers — translate DB rows into stable JSON shapes.
// -----------------------------------------------------------------

func contractView(c *repo.Contract) map[string]any {
	v := map[string]any{
		"id":                                   c.ID.String(),
		"customer_id":                          c.CustomerID.String(),
		"consumer_unit_id":                     c.ConsumerUnitID.String(),
		"vigencia_inicio":                      c.VigenciaInicio.Format("2006-01-02"),
		"desconto_percentual":                  c.DescontoPercentual.String(),
		"ip_faturamento_mode":                  string(c.IPFaturamentoMode),
		"ip_faturamento_valor":                 c.IPFaturamentoValor.String(),
		"ip_faturamento_percent":               c.IPFaturamentoPercent.String(),
		"bandeira_com_desconto":                c.BandeiraComDesconto,
		"custo_disponibilidade_sempre_cobrado": c.CustoDisponibilidadeSempreCobrado,
		"status":                               string(c.Status),
		"created_at":                           c.CreatedAt,
		"updated_at":                           c.UpdatedAt,
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

func calcView(c *repo.BillingCalculation) map[string]any {
	v := map[string]any{
		"id":                     c.ID.String(),
		"utility_invoice_ref_id": c.UtilityInvoiceRefID.String(),
		"billing_cycle_id":       c.BillingCycleID.String(),
		"consumer_unit_id":       c.ConsumerUnitID.String(),
		"contract_id":            c.ContractID.String(),
		"total_sem_desconto":     c.TotalSemDesconto.String(),
		"total_com_desconto":     c.TotalComDesconto.String(),
		"economia_rs":            c.EconomiaRS.String(),
		"economia_pct":           c.EconomiaPct.String(),
		"status":                 string(c.Status),
		"version":                c.Version,
		"calculated_at":          c.CalculatedAt,
		"needs_review_reasons":   c.NeedsReviewReasons,
	}
	if len(c.ContractSnapshotJSON) > 0 {
		var j any
		if err := json.Unmarshal(c.ContractSnapshotJSON, &j); err == nil {
			v["contract_snapshot"] = j
		}
	}
	if len(c.InputsSnapshotJSON) > 0 {
		var j any
		if err := json.Unmarshal(c.InputsSnapshotJSON, &j); err == nil {
			v["inputs_snapshot"] = j
		}
	}
	if len(c.ResultSnapshotJSON) > 0 {
		var j any
		if err := json.Unmarshal(c.ResultSnapshotJSON, &j); err == nil {
			v["result_snapshot"] = j
		}
	}
	if c.ApprovedAt != nil {
		v["approved_at"] = *c.ApprovedAt
	}
	if c.ApprovedBy != nil {
		v["approved_by"] = c.ApprovedBy.String()
	}
	return v
}

func parseDate(s string) (timeOnly, error) {
	return parseISODate(s)
}
