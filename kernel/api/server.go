package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sisoputnfrba/tp-golang/kernel/api/handler"
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
	s.router.HandleFunc("/process/{pid}", func(w http.ResponseWriter, r *http.Request) {
		handler.ProcessHandler(w,r, s.logger)
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
