package api

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/api/handlers"
	global "github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/utils/server"
)

func CreateServer() *server.Server {

	configServer := server.Config{
		Port: global.KernelConfig.Port,
		Handlers: map[string]http.HandlerFunc{
			"GET /process":          handlers.ListProcessHandler,
			"GET /process/{pid}":    handlers.ProcessByIdHandler,
			"PUT /process":          handlers.InitProcessHandler,
			"DELETE /process/{pid}": handlers.EndProcessHandler,
			"PUT /plani":            handlers.InitPlanningHandler,
			"DELETE /plani":         handlers.StopPlanningHandler,
			"PUT /newio":            handlers.NewIO,
			"GET /interrupt":        handlers.Interrupt,
		},
	}
	return server.NewServer(configServer)
}
