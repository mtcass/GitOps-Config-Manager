package server

import (
	"fmt"
	"net/http"
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

// handleListConfigs fetches files from provider and returns an HTML snippet for HTMX
func (s *Server) handleListConfigs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Call our Git client to list files
	files, err := s.client.ListConfigFiles(s.cfg.ConfigDirPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `<div class="text-sm text-red-400 bg-red-950/50 border border-red-900 p-3 rounded-lg">Failed to load files: %v</div>`, err)
		return
	}

	if len(files) == 0 {
		fmt.Fprint(w, `<div class="text-sm text-slate-500 italic">No configuration files found in the repository.</div>`)
		return
	}

	// Render a clean list of files as HTML snippet
	html := `<div class="divide-y divide-slate-800 border border-slate-800 rounded-lg overflow-hidden bg-slate-950">`
	for _, file := range files {
		html += fmt.Sprintf(`
			<div class="p-4 flex items-center justify-between hover:bg-slate-900/50 transition">
				<div>
					<span class="font-medium text-slate-200">%s</span>
					<div class="text-xs text-slate-500 mt-0.5">SHA: %s</div>
				</div>
				<button class="px-3 py-1.5 text-xs font-medium bg-indigo-600 hover:bg-indigo-500 text-white rounded-md transition">
					Edit File
				</button>
			</div>`, file.Path, file.SHA)
	}
	html += `</div>`

	w.Write([]byte(html))
}
