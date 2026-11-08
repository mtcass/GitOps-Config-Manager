package server

import (
	"fmt"
	"net/http"

	"github.com/mtcass/gitops-config-manager/internal/config"
	"github.com/mtcass/gitops-config-manager/internal/gitprovider"
)

type Server struct {
	cfg    *config.Config
	client gitprovider.GitProvider
}

func NewServer(cfg *config.Config, client gitprovider.GitProvider) *Server {
	return &Server{
		cfg:    cfg,
		client: client,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Register HTTP routes
	mux.HandleFunc("GET /", s.handleDashboard)
	mux.HandleFunc("GET /openapi.yaml", s.handleOpenAPIYAML)
	mux.HandleFunc("GET /docs", s.handleSwaggerUI)
	mux.HandleFunc("GET /api/configs", s.handleListConfigs)

	addr := fmt.Sprintf(":%s", s.cfg.Port)
	fmt.Printf("Starting HTTP server on http://localhost%s...\n", addr)

	return http.ListenAndServe(addr, mux)
}
