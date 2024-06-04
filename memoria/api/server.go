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
			"GET  /ping":         func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("msg: Se conecto a Memoria")) },
			"PUT /process":       handlers.CodeReciever,
			"PUT /process/{pid}": handlers.SendInstruction,
			"PUT /mov_in":        handlers.Mov_in,
			"PUT /mov_out":       handlers.Mov_out,
		},
	}
	return server.NewServer(configServer)
}
