package longterm

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

var working bool

// Mover procesos a READY mientras sean < que la multiprogramacion
// Tiene que estar corriendo todo el tiempo en un hilo?
func InitLongTermPlani() {
	working = true

	for working {
		global.SemMulti <- 0
		sendPCBToReady()
		global.Logger.Log(fmt.Sprintf("PCB to READY - Semaforo %d - Multi: %d", len(global.SemMulti), global.KernelConfig.Multiprogramming), log.DEBUG)
	}
}

func sendPCBToReady() {
	global.MutexNewState.Lock()
	pcbToReady := global.NewState.Remove(global.NewState.Front()).(*model.PCB)
	global.MutexNewState.Unlock()

	pcbToReady.State = "READY"

	global.MutexReadyState.Lock()
	global.ReadyState.PushBack(pcbToReady)
	global.MutexReadyState.Unlock()
}
