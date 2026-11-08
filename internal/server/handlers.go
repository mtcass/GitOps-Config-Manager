package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// handleDashboard renders the main HTML shell of the application
func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
<html lang="en" class="dark">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GitOps Config Manager</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-slate-950 text-slate-100 min-h-screen font-sans">
    <div class="max-w-6xl mx-auto p-8">
        <header class="border-b border-slate-800 pb-6 mb-8 flex justify-between items-center">
            <div>
                <h1 class="text-2xl font-bold tracking-tight">GitOps Config Manager</h1>
                <p class="text-sm text-slate-400">Zero-database configuration management powered by Git</p>
            </div>
            <div class="text-xs bg-slate-900 border border-slate-800 px-3 py-1.5 rounded-lg text-slate-400">
                Repo: <span class="text-slate-200 font-medium">` + s.cfg.GitHubOwner + `/` + s.cfg.GitHubRepo + `</span>
            </div>
        </header>

        <main class="grid gap-6">
            <div class="bg-slate-900 border border-slate-800 rounded-xl p-6 shadow-xl">
                <h2 class="text-lg font-semibold mb-2">Configuration Files</h2>
                <p class="text-sm text-slate-400 mb-4">Files tracked in directory: <code class="text-indigo-400">` + s.cfg.ConfigDirPath + `</code></p>
                
                <!-- HTMX container that automatically fetches files on page load -->
                <div id="config-list" hx-get="/api/configs" hx-trigger="load" class="space-y-3">
                    <div class="text-sm text-slate-500 animate-pulse">Loading configuration files from GitHub...</div>
                </div>
            </div>
        </main>
    </div>
</body>
</html>`

	w.Write([]byte(html))
}

// handleOpenAPIYAML serves the raw OpenAPI specification file
func (s *Server) handleOpenAPIYAML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-yaml; charset=utf-8")
	// Lê o arquivo openapi.yaml diretamente do diretório api/
	http.ServeFile(w, r, "api/openapi.yaml")
}

// handleSwaggerUI renders the Swagger UI interface page
func (s *Server) handleSwaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>GitOps Config Manager - API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
</head>
<body style="margin: 0; background: #fafafa;">
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js"></script>
    <script>
        window.onload = () => {
            window.ui = SwaggerUIBundle({
                url: '/openapi.yaml',
                dom_id: '#swagger-ui',
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIBundle.StandalonePreset
                ],
            });
        };
    </script>
</body>
</html>`

	w.Write([]byte(html))
}

// handleListConfigs fetches files from Git provider and returns a clean array of relative filenames
func (s *Server) handleListConfigs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	files, err := s.client.ListConfigFiles(s.cfg.ConfigDirPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("failed to list configuration files: %v", err),
		})
		return
	}

	// Remove config dir prefix (CONFIG_DIR_PATH) from the file name
	paths := make([]string, 0, len(files))
	for _, f := range files {

		relPath := strings.TrimPrefix(f.Path, s.cfg.ConfigDirPath)
		relPath = strings.TrimLeft(relPath, "/")

		if relPath != "" {
			paths = append(paths, relPath)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paths)
}
