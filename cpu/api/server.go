package api

import (
	"net/http"

	global "github.com/sisoputnfrba/tp-golang/cpu/global"
	"github.com/sisoputnfrba/tp-golang/utils/server"
)

func CreateServer() *server.Server {

	configServer := server.Config{
		Port: global.CPUConfig.Port,
		Handlers: map[string]http.HandlerFunc{
			"GET /ping": func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("msg: Se conecto a CPU")) },
		},
	}
	return server.NewServer(configServer)
}
