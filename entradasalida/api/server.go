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
<<<<<<< HEAD
			"GET /Ping":   handlers.Ping,
			"PUT /Sleep": handlers.Sleep,
=======
			"POST /sleep": handlers.Sleep,
>>>>>>> 132a78c1d3d248f566e3a9683820fcc47bb8fe79
		},
	}
	return server.NewServer(configServer)
}
