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
			"GET /ping":         func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("msg: Se conecto a Memoria")) },
			"PUT /stdin_read":   handlers.Stdin_read,
			"PUT /stdout_write": handlers.Stdout_write,
		},
	}
	return server.NewServer(configServer)
}
