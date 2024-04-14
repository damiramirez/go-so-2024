package api

import (
	"fmt"
	"net/http"

	handler "github.com/sisoputnfrba/tp-golang/kernel/api/handler"
	"github.com/sisoputnfrba/tp-golang/kernel/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

type Server struct {
	router *http.ServeMux
	logger *log.LoggerStruct
}

// Crea una nueva instancia de Server
func NewServer() *Server {
	return &Server{
		router: http.NewServeMux(),
		logger: global.Logger,
	}
}

// Configura las rutas de la API
func (s *Server) ConfigureRoutes() {
	s.router.HandleFunc("GET /process", handler.ListProcessHandler)
	s.router.HandleFunc("GET /process/{pid}", handler.ProcessStateHandler)
	s.router.HandleFunc("PUT /process", handler.InitProcessHandler)
	s.router.HandleFunc("DELETE /process/{pid}", handler.EndProcessHandler)
	s.router.HandleFunc("PUT /plani", handler.InitPlanningHandler)
	s.router.HandleFunc("DELETE /plani", handler.StopPlanningHandler)
}

// Inicia el servidor HTTP
func (s *Server) Start() error {
	s.ConfigureRoutes()

	port := global.KernelConfig.Port

	s.logger.Log(fmt.Sprintf("Starting kernel API server on port %d", port), log.INFO)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), s.router)
	if err != nil {
		return err
	}

	return nil
}
