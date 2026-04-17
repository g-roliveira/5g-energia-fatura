package app

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
)

//go:embed docs/backend-go-api.md
var backendDocsMarkdown string

func docsMarkdownHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	_, _ = w.Write([]byte(backendDocsMarkdown))
}

func docsHTMLHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, `<!doctype html>
<html lang="pt-BR">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Swagger UI - Backend Go API</title>
  <style>
    html, body { margin: 0; padding: 0; height: 100%%; }
    #swagger-ui { height: 100%%; }
    .top-links {
      position: fixed;
      z-index: 99999;
      right: 16px;
      top: 12px;
      font-family: system-ui, -apple-system, "Segoe UI", sans-serif;
      font-size: 12px;
      background: rgba(0,0,0,.7);
      border-radius: 6px;
      padding: 6px 10px;
    }
    .top-links a { color: #fff; text-decoration: none; }
  </style>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div class="top-links"><a href="/docs.md">docs.md</a></div>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.ui = SwaggerUIBundle({
      url: "/openapi.json",
      dom_id: "#swagger-ui",
      deepLinking: true,
      displayRequestDuration: true
    });
  </script>
</body>
</html>`)
}

func openAPIJSONHandler(specFn func() map[string]any) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		spec := specFn()
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		_ = encoder.Encode(spec)
	}
}
