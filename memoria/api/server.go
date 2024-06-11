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
			"DELETE /process/{pid}": handlers.DeleteProcess,
			"PUT /process":          handlers.CodeReciever,    //me envian el path y el pid
			"PUT /process/{pid}":    handlers.SendInstruction, //envio instruccion segun el pc
			"PUT /resize":           handlers.Resize,          //agrando o achico tamaño del proceso
			"PUT /framenumber":      handlers.PageTableAccess,
			"PUT /memIn":            handlers.MemoryAccessIn,  //leer en memoria
			"PUT /memOut":           handlers.MemoryAccessOut, //escribir en memoria
			"PUT /stdin_read":       handlers.Stdin_read, //escribe en memoria
			"PUT /stdout_write":     handlers.Stdout_write,//lee en memoria
			"PUT /mov_in":        handlers.Mov_in,
			"PUT /mov_out":       handlers.Mov_out,
			"PUT /resize":        handlers.Resize,
		},
	}
	return server.NewServer(configServer)
}
