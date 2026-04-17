package app

import (
	_ "embed"
	"fmt"
	"html"
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
  <title>Backend Go API Docs</title>
  <style>
    :root {
      color-scheme: light dark;
      font-family: Inter, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
    }
    body {
      margin: 0;
      background: #0b1020;
      color: #e8edf7;
    }
    main {
      max-width: 980px;
      margin: 0 auto;
      padding: 32px 20px 48px;
    }
    a {
      color: #7dd3fc;
    }
    .actions {
      margin-bottom: 16px;
    }
    pre {
      white-space: pre-wrap;
      word-break: break-word;
      background: #11192f;
      border: 1px solid #223255;
      border-radius: 12px;
      padding: 20px;
      overflow: auto;
      line-height: 1.5;
      font-size: 14px;
    }
  </style>
</head>
<body>
  <main>
    <div class="actions">
      <a href="/docs.md">Ver Markdown bruto</a>
    </div>
    <pre>%s</pre>
  </main>
</body>
</html>`, html.EscapeString(backendDocsMarkdown))
}
