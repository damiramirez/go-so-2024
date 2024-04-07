package api

import (
	"fmt"
	"net/http"
	"os"

	handler "github.com/sisoputnfrba/tp-golang/kernel/api/handler"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

type Server struct {
	router *http.ServeMux
	logger log.Logger
}

// Crea una nueva instancia de Server
func NewServer(logger log.Logger) *Server {
	return &Server{
		router: http.NewServeMux(),
		logger: logger,
	}
}

// Configura las rutas de la API
func (s *Server) ConfigureRoutes() {
	s.router.HandleFunc("GET /process", func(w http.ResponseWriter, r *http.Request) {
		handler.ListProcessHandler(w, r, s.logger)
	})
	s.router.HandleFunc("GET /process/{pid}", func(w http.ResponseWriter, r *http.Request) {
		handler.ProcessStateHandler(w, r, s.logger)
	})
	s.router.HandleFunc("PUT /process", func(w http.ResponseWriter, r *http.Request) {
		handler.InitProcessHandler(w, r, s.logger)
	})
	s.router.HandleFunc("DELETE /process/{pid}", func(w http.ResponseWriter, r *http.Request) {
		handler.EndProcessHandler(w, r, s.logger)
	})
	s.router.HandleFunc("PUT /plani", func(w http.ResponseWriter, r *http.Request) {
		handler.InitPlanningHandler(w, r, s.logger)
	})
	s.router.HandleFunc("DELETE /plani", func(w http.ResponseWriter, r *http.Request) {
		handler.StopPlanningHandler(w, r, s.logger)
	})
}

// Inicia el servidor HTTP
func (s *Server) Start(port int) {
	s.ConfigureRoutes()

	s.logger.Log(fmt.Sprintf("Starting kernel API server on port %d", port), log.INFO)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), s.router)
	if err != nil {
		s.logger.Log(fmt.Sprintf("Failed to start kernel API server: %v", err), log.ERROR)
		os.Exit(1)
	}
}
