package handlers

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

// Este mensaje se encargará de retomar (en caso que se encuentre pausada)
// la planificación de corto y largo plazo. En caso que la planificación no
// se encuentre pausada, se debe ignorar el mensaje.
func InitPlanningHandler(w http.ResponseWriter, r *http.Request) {

	global.Logger.Log("Init plani", log.DEBUG)

	// Testeando
	resp, _ := http.Get("http://127.0.0.1:8003/ping")
	/*if err!=nil{
		http.Error(w,"error al conectarse al cpu ",500)
		return
	}*/
	global.Logger.Log("Resp de CPU: "+resp.Status, log.INFO)

	// TODO: Manejar planificacion

	w.WriteHeader(http.StatusNoContent)
}

// Este mensaje se encargará de pausar la planificación de corto y largo plazo.
// El proceso que se encuentra en ejecución NO es desalojado, pero una vez que salga
// de EXEC se va a pausar el manejo de su motivo de desalojo. De la misma forma,
// los procesos bloqueados van a pausar su transición a la cola de Ready.
func StopPlanningHandler(w http.ResponseWriter, r *http.Request) {
	global.Logger.Log("Stop plani", log.DEBUG)
	// TODO: Frenar planificacion

	w.WriteHeader(http.StatusNoContent)
}
