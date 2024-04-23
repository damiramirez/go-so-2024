package api

import (
	"net/http"
	handlers "github.com/sisoputnfrba/tp-golang/memoria/api/handlers"
	global "github.com/sisoputnfrba/tp-golang/memoria/global"
	server "github.com/sisoputnfrba/tp-golang/utils/server"
)
var Memory=global.NewMemory()
func CreateServer() *server.Server {

	configServer := server.Config{
		Port: global.MemoryConfig.Port,
		Handlers: map[string]http.HandlerFunc{
			"GET  /ping": func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("msg: Se conecto a Memoria")) },
			"POST /process": func(w http.ResponseWriter, r *http.Request) {handlers.CodeReciever(w,r,Memory)},
		},
	}
	return server.NewServer(configServer)
}
