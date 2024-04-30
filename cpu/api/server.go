package api

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/cpu/api/handlers"
	global "github.com/sisoputnfrba/tp-golang/cpu/global"
	"github.com/sisoputnfrba/tp-golang/utils/server"
)

func CreateServer() *server.Server {

	configServer := server.Config{
		Port: global.CPUConfig.Port,
		Handlers: map[string]http.HandlerFunc{
			"PUT /process":  handlers.PCBreciever,
			"PUT /dispatch": handlers.Dispatch,
		},
	}
	return server.NewServer(configServer)
}
