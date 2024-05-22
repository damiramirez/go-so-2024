package api

import (
	"net/http"
	handlers "github.com/sisoputnfrba/tp-golang/memoria/api/handlers"
	global "github.com/sisoputnfrba/tp-golang/memoria/global"
	server "github.com/sisoputnfrba/tp-golang/utils/server"
)

func CreateServer() *server.Server {

	configServer := server.Config{
		Port: global.MemoryConfig.Port,
		Handlers: map[string]http.HandlerFunc{
			"GET /process":       handlers.DeleteProcess,
			"PUT /process":       handlers.CodeReciever,
			"PUT /process/{pid}": handlers.SendInstruction,
			"GET /resize":        handlers.Resize,
			"PUT /framenumber":   handlers.PageTableAccess,
			"PUT /memaccess":     handlers.MemoryAccess,
		},
	}
	return server.NewServer(configServer)
}
