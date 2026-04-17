package app

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func withRequestLogging(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		logger.Info(
			"http_request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeClientError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeInternalError(w http.ResponseWriter, logger *slog.Logger, operation string, err error) {
	logger.Error("internal_error", "operation", operation, "error", sanitizeError(err))
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal_error"})
}

func sanitizeError(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	replacements := []string{"senha", "token", "bearer", "authorization", "documento"}
	lower := strings.ToLower(msg)
	for _, marker := range replacements {
		if strings.Contains(lower, marker) {
			return "redacted_error"
		}
	}
	return msg
}
