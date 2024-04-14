package api

import (
	"net/http"

	handlers "github.com/sisoputnfrba/tp-golang/kernel/api/handlers"
	global "github.com/sisoputnfrba/tp-golang/kernel/global"
	server "github.com/sisoputnfrba/tp-golang/utils/server"
)

func NewServerConfig() server.Config {
	return server.Config{
		Port: global.KernelConfig.Port,
		Handlers: map[string]http.HandlerFunc{
			"GET /process":          handlers.ListProcessHandler,
			"GET /process/{pid}":    handlers.ProcessByIdHandler,
			"PUT /process":          handlers.InitProcessHandler,
			"DELETE /process/{pid}": handlers.EndProcessHandler,
			"PUT /plani":            handlers.InitPlanningHandler,
			"DELETE /plani":         handlers.StopPlanningHandler,
		},
	}
}
