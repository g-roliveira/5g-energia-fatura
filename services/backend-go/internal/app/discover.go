package app

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"

	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/neoenergia"
	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/session"
)

type DiscoveryResult struct {
	MinhaConta       *neoenergia.MinhaContaResponse       `json:"minha_conta,omitempty"`
	MinhaContaLegado *neoenergia.MinhaContaLegadoResponse `json:"minha_conta_legado,omitempty"`
	UCs              []neoenergia.UC                      `json:"ucs"`
	Errors           map[string]string                    `json:"errors,omitempty"`
}

func handleDiscover(
	w http.ResponseWriter,
	r *http.Request,
	credentialID string,
	sm *session.Manager,
	api *neoenergia.Client,
	logger *slog.Logger,
) {
	resolved, err := sm.ResolveToken(r.Context(), credentialID)
	if err != nil {
		writeInternalError(w, logger, "discover_resolve_token", err)
		return
	}

	reqCtx := neoenergia.RequestContext{
		BearerToken: resolved.Token,
		Documento:   resolved.Documento,
	}

	type resT[T any] struct {
		val T
		err error
	}
	minhaCh := make(chan resT[neoenergia.MinhaContaResponse], 1)
	legadoCh := make(chan resT[neoenergia.MinhaContaLegadoResponse], 1)
	ucsCh := make(chan resT[neoenergia.UCsResponse], 1)

	ctx := r.Context()
	go func() {
		v, e := api.GetMinhaConta(ctx, reqCtx)
		minhaCh <- resT[neoenergia.MinhaContaResponse]{v, e}
	}()
	go func() {
		v, e := api.GetMinhaContaLegado(ctx, reqCtx)
		legadoCh <- resT[neoenergia.MinhaContaLegadoResponse]{v, e}
	}()
	go func() {
		v, e := api.ListUCs(ctx, reqCtx)
		ucsCh <- resT[neoenergia.UCsResponse]{v, e}
	}()

	minhaR := <-minhaCh
	legadoR := <-legadoCh
	ucsR := <-ucsCh

	out := DiscoveryResult{
		UCs:    []neoenergia.UC{},
		Errors: map[string]string{},
	}

	if minhaR.err != nil {
		out.Errors["minha_conta"] = discoveryErrorCode(minhaR.err)
	} else {
		mc := minhaR.val
		out.MinhaConta = &mc
	}

	if legadoR.err != nil {
		out.Errors["minha_conta_legado"] = discoveryErrorCode(legadoR.err)
	} else {
		ml := legadoR.val
		out.MinhaContaLegado = &ml
	}

	if ucsR.err != nil {
		writeInternalError(w, logger, "discover_list_ucs", ucsR.err)
		return
	}
	out.UCs = ucsR.val.UCs

	if len(out.Errors) == 0 {
		out.Errors = nil
	}

	writeJSON(w, http.StatusOK, out)
}

func discoveryErrorCode(err error) string {
	if err == nil {
		return ""
	}

	if errors.Is(err, context.Canceled) {
		return "request_canceled"
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return "timeout"
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return "timeout"
	}

	var apiErr *neoenergia.ErrorResponse
	if errors.As(err, &apiErr) {
		switch apiErr.StatusCode {
		case http.StatusUnauthorized, http.StatusForbidden:
			return "unauthorized"
		case http.StatusNotFound:
			return "not_found"
		case http.StatusTooManyRequests:
			return "rate_limited"
		case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			return "upstream_unavailable"
		default:
			if apiErr.StatusCode >= 500 {
				return "upstream_unavailable"
			}
			return "upstream_error"
		}
	}

	return "internal_error"
}
