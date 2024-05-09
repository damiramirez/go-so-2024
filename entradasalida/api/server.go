package api

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/entradasalida/api/handlers"
	global "github.com/sisoputnfrba/tp-golang/entradasalida/global"
	"github.com/sisoputnfrba/tp-golang/utils/server"
)

func CreateServer() *server.Server {

	configServer := server.Config{
		Port: global.IOConfig.Port,
		Handlers: map[string]http.HandlerFunc{
			"PUT /IO_GEN_SLEEP": handlers.Sleep,
		},
	}
	return server.NewServer(configServer)
}
