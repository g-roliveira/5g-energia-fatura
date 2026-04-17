package app

import (
	"sort"
	"strconv"
	"strings"
)

type routeDoc struct {
	Method      string
	Path        string
	Summary     string
	Tags        []string
	SuccessCode int
	OperationID string
	Description string
}

type routeCatalog struct {
	routes []routeDoc
}

func newRouteCatalog() *routeCatalog {
	return &routeCatalog{routes: make([]routeDoc, 0, 32)}
}

func (c *routeCatalog) add(method, path, summary string, tags []string, successCode int) {
	method = strings.ToUpper(strings.TrimSpace(method))
	if successCode <= 0 {
		successCode = 200
	}
	c.routes = append(c.routes, routeDoc{
		Method:      method,
		Path:        path,
		Summary:     summary,
		Tags:        append([]string(nil), tags...),
		SuccessCode: successCode,
		OperationID: makeOperationID(method, path),
	})
}

func (c *routeCatalog) spec() map[string]any {
	paths := map[string]map[string]any{}
	for _, r := range c.routes {
		method := strings.ToLower(r.Method)
		item, ok := paths[r.Path]
		if !ok {
			item = map[string]any{}
			paths[r.Path] = item
		}
		item[method] = map[string]any{
			"tags":        r.Tags,
			"summary":     r.Summary,
			"operationId": r.OperationID,
			"responses": map[string]any{
				statusKey(r.SuccessCode): map[string]any{
					"description": "Success",
				},
				"default": map[string]any{
					"description": "Error",
				},
			},
		}
	}

	// Keep tags sorted for deterministic docs output.
	seen := map[string]bool{}
	tags := make([]string, 0, 8)
	for _, r := range c.routes {
		for _, tag := range r.Tags {
			if !seen[tag] {
				seen[tag] = true
				tags = append(tags, tag)
			}
		}
	}
	sort.Strings(tags)
	tagObjs := make([]map[string]string, 0, len(tags))
	for _, tag := range tags {
		tagObjs = append(tagObjs, map[string]string{"name": tag})
	}

	return map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title":       "5G Energia Fatura API",
			"version":     "1.0.0",
			"description": "Backend Go API para sincronizacao de faturas Neoenergia e consulta via banco.",
		},
		"servers": []map[string]string{
			{"url": "/"},
		},
		"paths": paths,
		"tags":  tagObjs,
	}
}

func makeOperationID(method, path string) string {
	path = strings.TrimSpace(path)
	path = strings.Trim(path, "/")
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, "{", "")
	path = strings.ReplaceAll(path, "}", "")
	if path == "" {
		path = "root"
	}
	return strings.ToLower(method) + "_" + path
}

func statusKey(code int) string {
	if code <= 0 {
		return "200"
	}
	return strconv.Itoa(code)
}
