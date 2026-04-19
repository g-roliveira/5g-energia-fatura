package app

import (
	"io"
	"log/slog"
	"net/http"
	"testing"
)

func TestRegisterBillingRoutesAddsOpenAPIPaths(t *testing.T) {
	mux := http.NewServeMux()
	docs := newRouteCatalog()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	RegisterBillingRoutes(mux, docs, &BillingDeps{}, logger)

	spec := docs.spec()
	paths, ok := spec["paths"].(map[string]any)
	if !ok {
		t.Fatalf("paths missing or invalid in spec: %T", spec["paths"])
	}

	expected := []struct {
		path   string
		method string
	}{
		{path: "/v1/billing/contracts", method: "post"},
		{path: "/v1/billing/contracts/{id}", method: "get"},
		{path: "/v1/billing/consumer-units/{uc_id}/active-contract", method: "get"},
		{path: "/v1/billing/calculations/{id}", method: "get"},
	}

	for _, item := range expected {
		pathNode, ok := paths[item.path].(map[string]any)
		if !ok {
			t.Fatalf("path %s missing from openapi", item.path)
		}
		if _, ok := pathNode[item.method]; !ok {
			t.Fatalf("method %s missing for path %s", item.method, item.path)
		}
	}
}
