package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"solace/internal/engine"
)

type Server struct {
	engine *engine.HardeningEngine
	port   string
}

func NewServer(e *engine.HardeningEngine, port string) *Server {
	return &Server{
		engine: e,
		port:   port,
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/api/audit", s.handleAudit)
	http.HandleFunc("/api/status", s.handleStatus)

	fmt.Printf("🌐 API Server running at http://localhost:%s\n", s.port)
	return http.ListenAndServe(":"+s.port, nil)
}

func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	// Enable CORS for frontend development
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	results := s.engine.EvaluateRules()
	json.NewEncoder(w).Encode(results)
}

// health check
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	
	status := map[string]string{"status": "online", "service": "Solace Security Toolkit"}
	json.NewEncoder(w).Encode(status)
}