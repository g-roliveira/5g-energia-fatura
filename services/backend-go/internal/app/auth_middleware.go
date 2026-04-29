package app

import (
	"net/http"
	"strings"
)

func withOptionalAPIKeyAuth(apiKey string, next http.Handler) http.Handler {
	trimmed := strings.TrimSpace(apiKey)
	if trimmed == "" {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Infra endpoints remain public even when API key protection is enabled.
		if !strings.HasPrefix(r.URL.Path, "/v1/") {
			next.ServeHTTP(w, r)
			return
		}

		if strings.TrimSpace(r.Header.Get("X-API-Key")) != trimmed {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
