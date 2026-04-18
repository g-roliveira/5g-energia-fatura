package app

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/extractor"
	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/neoenergia"
	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/security"
	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/session"
	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/store"
	syncsvc "github.com/gustavo/5g-energia-fatura/services/backend-go/internal/sync"
)

type Server struct {
	cfg    Config
	logger *slog.Logger
	mux    http.Handler
}

func NewServer(cfg Config, logger *slog.Logger) (*Server, error) {
	mux := http.NewServeMux()
	docs := newRouteCatalog()
	apiClient := neoenergia.NewClient(cfg.NeoenergiaBaseURL)
	extractorClient := extractor.NewClient(cfg.ExtractorBaseURL)
	sqliteStore, err := store.OpenSQLite(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	syncService := syncsvc.NewService(apiClient, extractorClient, sqliteStore)
	cipher, err := security.NewCipher(cfg.EncryptionKey)
	if err != nil {
		return nil, err
	}
	sessionManager := session.NewManager(
		sqliteStore,
		cipher,
		session.BootstrapRunner{
			PythonBin:  cfg.BootstrapPythonBin,
			ScriptPath: cfg.BootstrapScript,
		},
	)

	docs.add(http.MethodGet, "/healthz", "Health check", []string{"infra"}, http.StatusOK)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "backend-go"})
	})
	docs.add(http.MethodGet, "/openapi.json", "OpenAPI schema", []string{"infra"}, http.StatusOK)
	mux.HandleFunc("/openapi.json", openAPIJSONHandler(docs.spec))
	docs.add(http.MethodGet, "/docs", "Swagger UI", []string{"infra"}, http.StatusOK)
	mux.HandleFunc("/docs", docsHTMLHandler)
	docs.add(http.MethodGet, "/docs.md", "Markdown docs", []string{"infra"}, http.StatusOK)
	mux.HandleFunc("/docs.md", docsMarkdownHandler)

	docs.add(http.MethodPost, "/v1/credentials", "Create credential", []string{"credentials"}, http.StatusCreated)
	mux.HandleFunc("/v1/credentials", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req session.CredentialInput
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeClientError(w, http.StatusBadRequest, "invalid_json")
			return
		}
		if err := sessionManager.RequireCredentialFields(req); err != nil {
			writeClientError(w, http.StatusBadRequest, err.Error())
			return
		}
		created, err := sessionManager.CreateCredential(r.Context(), req)
		if err != nil {
			writeInternalError(w, logger, "create_credential", err)
			return
		}
		writeJSON(w, http.StatusCreated, created)
	})

	docs.add(http.MethodPost, "/v1/credentials/{id}/session", "Create session from credential", []string{"credentials"}, http.StatusOK)
	docs.add(http.MethodGet, "/v1/credentials/{id}/discover", "Discover profile and UCs from credential", []string{"credentials"}, http.StatusOK)
	mux.HandleFunc("/v1/credentials/", func(w http.ResponseWriter, r *http.Request) {
		parts := splitPath(r.URL.Path)
		if len(parts) != 4 || parts[0] != "v1" || parts[1] != "credentials" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		credentialID, action := parts[2], parts[3]

		switch action {
		case "session":
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			sessionView, _, err := sessionManager.CreateSessionFromCredential(r.Context(), credentialID)
			if err != nil {
				writeInternalError(w, logger, "create_session_from_credential", err)
				return
			}
			writeJSON(w, http.StatusOK, sessionView)
		case "discover":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			handleDiscover(w, r, credentialID, sessionManager, apiClient, logger)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	docs.add(http.MethodPost, "/v1/sync/uc", "Sync consumer unit", []string{"sync"}, http.StatusOK)
	mux.HandleFunc("/v1/sync/uc", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var req syncsvc.SyncUCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeClientError(w, http.StatusBadRequest, "invalid_json")
			return
		}
		if req.UC == "" {
			writeClientError(w, http.StatusBadRequest, "uc é obrigatória")
			return
		}
		if req.BearerToken == "" && req.CredentialID == "" {
			writeClientError(w, http.StatusBadRequest, "bearer_token ou credential_id é obrigatório")
			return
		}

		if req.BearerToken == "" && req.CredentialID != "" {
			resolved, err := sessionManager.ResolveToken(r.Context(), req.CredentialID)
			if err != nil {
				writeInternalError(w, logger, "resolve_token", err)
				return
			}
			req.BearerToken = resolved.Token
			req.Documento = resolved.Documento
		}
		if req.BearerToken != "" && req.Documento == "" {
			writeClientError(w, http.StatusBadRequest, "documento é obrigatório quando bearer_token é enviado manualmente")
			return
		}

		result := syncService.SyncUC(r.Context(), req)
		writeJSON(w, http.StatusOK, result)
	})

	docs.add(http.MethodGet, "/v1/consumer-units", "List consumer units", []string{"consumer-units"}, http.StatusOK)
	mux.HandleFunc("/v1/consumer-units", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		limit := parsePositiveInt(r.URL.Query().Get("limit"), 100)
		status := r.URL.Query().Get("status")
		items, err := sqliteStore.ListConsumerUnits(limit, status)
		if err != nil {
			writeInternalError(w, logger, "list_consumer_units", err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": items, "limit": limit, "status": status})
	})

	docs.add(http.MethodGet, "/v1/consumer-units/{uc}", "Get consumer unit by UC", []string{"consumer-units"}, http.StatusOK)
	docs.add(http.MethodGet, "/v1/consumer-units/{uc}/invoices", "List invoices by UC", []string{"invoices"}, http.StatusOK)
	docs.add(http.MethodGet, "/v1/consumer-units/{uc}/latest-invoice", "Get latest invoice by UC", []string{"invoices"}, http.StatusOK)
	docs.add(http.MethodPost, "/v1/consumer-units/{uc}/sync", "Sync specific UC", []string{"sync"}, http.StatusOK)
	mux.HandleFunc("/v1/consumer-units/", func(w http.ResponseWriter, r *http.Request) {
		parts := splitPath(r.URL.Path)
		if len(parts) < 3 || parts[0] != "v1" || parts[1] != "consumer-units" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		uc := parts[2]
		if uc == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if len(parts) == 4 && parts[3] == "invoices" {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			limit := parsePositiveInt(r.URL.Query().Get("limit"), 100)
			status := r.URL.Query().Get("status")
			items, err := sqliteStore.ListInvoicesByUC(uc, limit, status)
			if err != nil {
				writeInternalError(w, logger, "list_invoices_by_uc", err)
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"uc": uc, "items": items, "limit": limit, "status": status})
			return
		}

		if len(parts) == 3 {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			item, err := sqliteStore.GetConsumerUnitByUC(uc)
			if err != nil {
				writeInternalError(w, logger, "get_consumer_unit_by_uc", err)
				return
			}
			if item == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			writeJSON(w, http.StatusOK, item)
			return
		}

		if len(parts) == 4 && parts[3] == "latest-invoice" {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			item, err := sqliteStore.GetLatestInvoiceByUC(uc)
			if err != nil {
				writeInternalError(w, logger, "get_latest_invoice_by_uc", err)
				return
			}
			if item == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			writeJSON(w, http.StatusOK, item)
			return
		}

		if len(parts) == 4 && parts[3] == "sync" {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			var req syncsvc.SyncUCRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
				writeClientError(w, http.StatusBadRequest, "invalid_json")
				return
			}
			req.UC = uc
			if req.BearerToken == "" && req.CredentialID == "" {
				writeClientError(w, http.StatusBadRequest, "bearer_token ou credential_id é obrigatório")
				return
			}
			if req.BearerToken == "" && req.CredentialID != "" {
				resolved, err := sessionManager.ResolveToken(r.Context(), req.CredentialID)
				if err != nil {
					writeInternalError(w, logger, "resolve_token", err)
					return
				}
				req.BearerToken = resolved.Token
				req.Documento = resolved.Documento
			}
			if req.BearerToken != "" && req.Documento == "" {
				writeClientError(w, http.StatusBadRequest, "documento é obrigatório quando bearer_token é enviado manualmente")
				return
			}
			result := syncService.SyncUC(r.Context(), req)
			writeJSON(w, http.StatusOK, result)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	})

	docs.add(http.MethodGet, "/v1/invoices/{id}", "Get invoice by id", []string{"invoices"}, http.StatusOK)
	mux.HandleFunc("/v1/invoices/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		parts := splitPath(r.URL.Path)
		if len(parts) != 3 || parts[0] != "v1" || parts[1] != "invoices" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		item, err := sqliteStore.GetInvoiceByID(parts[2])
		if err != nil {
			writeInternalError(w, logger, "get_invoice_by_id", err)
			return
		}
		if item == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, item)
	})

	docs.add(http.MethodGet, "/v1/sync-runs/{id}", "Get sync run by id", []string{"sync"}, http.StatusOK)
	mux.HandleFunc("/v1/sync-runs/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		parts := splitPath(r.URL.Path)
		if len(parts) != 3 || parts[0] != "v1" || parts[1] != "sync-runs" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		item, err := sqliteStore.GetSyncRunByID(parts[2])
		if err != nil {
			writeInternalError(w, logger, "get_sync_run_by_id", err)
			return
		}
		if item == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, item)
	})

	docs.add(http.MethodGet, "/v1/extractor/contracts", "List extractor contracts", []string{"extractor"}, http.StatusOK)
	mux.HandleFunc("/v1/extractor/contracts", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"extractor_request":  "packages/contracts/extractor-request.schema.json",
			"extractor_response": "packages/contracts/extractor-response.schema.json",
		})
	})

	rootHandler := withRequestLogging(logger, mux)
	return &Server{
		cfg:    cfg,
		logger: logger,
		mux:    rootHandler,
	}, nil
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)
	s.logger.Info("server_start", "addr", addr, "extractor_base_url", s.cfg.ExtractorBaseURL)
	return http.ListenAndServe(addr, s.mux)
}

func splitPath(path string) []string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "/")
}

func parsePositiveInt(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
