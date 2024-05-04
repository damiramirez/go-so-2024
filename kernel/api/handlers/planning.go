package handlers

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/longterm"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/shortterm"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

// Este mensaje se encargará de retomar (en caso que se encuentre pausada)
// la planificación de corto y largo plazo. En caso que la planificación no
// se encuentre pausada, se debe ignorar el mensaje.
func InitPlanningHandler(w http.ResponseWriter, r *http.Request) {

	global.Logger.Log("Init plani", log.DEBUG)
	
	
	startShortTerm := make(chan bool)  // Crear un canal para sincronizar el inicio de shortterm

	// Lanzar la planificación a largo plazo en una goroutine
	go func() {
		longterm.InitLongTermPlani()
		startShortTerm <- true  // Enviar señal para iniciar shortterm después de empezar longterm
	}()
	
	// Lanzar la planificación a corto plazo en otra goroutine
	go func() {
		<-startShortTerm  // Esperar la señal para empezar
		shortterm.InitShortTermPlani()
	}()

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
