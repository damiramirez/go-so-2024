package handlers

import (
	"context"
	"net/http"
	"sync"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/longterm"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/shortterm"
)

// Este mensaje se encargará de retomar (en caso que se encuentre pausada)
// la planificación de corto y largo plazo. En caso que la planificación no
// se encuentre pausada, se debe ignorar el mensaje.

var (
	ctx       context.Context
	ctxCancel context.CancelFunc
)

func InitPlanningHandler(w http.ResponseWriter, r *http.Request) {

	ctx, ctxCancel = context.WithCancel(context.Background())

	global.MutexPlani.Lock()
	global.WorkingPlani = true
	global.MutexPlani.Unlock()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		longterm.InitLongTermPlani(ctx)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		shortterm.InitShortTermPlani()
		wg.Done()
	}()

	wg.Wait()

	w.WriteHeader(http.StatusNoContent)
}

// Este mensaje se encargará de pausar la planificación de corto y largo plazo.
// El proceso que se encuentra en ejecución NO es desalojado, pero una vez que salga
// de EXEC se va a pausar el manejo de su motivo de desalojo. De la misma forma,
// los procesos bloqueados van a pausar su transición a la cola de Ready.
func StopPlanningHandler(w http.ResponseWriter, r *http.Request) {

	global.MutexPlani.Lock()
	ctxCancel()
	global.WorkingPlani = false
	global.MutexPlani.Unlock()

	w.WriteHeader(http.StatusNoContent)
}
