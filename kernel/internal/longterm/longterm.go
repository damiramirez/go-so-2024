package longterm

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

// Mover procesos a READY mientras sean < que la multiprogramacion
// Tiene que estar corriendo todo el tiempo en un hilo?
func InitLongTermPlani() {
	global.WorkingPlani = true
	for global.WorkingPlani {
		if global.NewState.Len() != 0 {
			global.Logger.Log(fmt.Sprintf("NEW LEN: %d", global.NewState.Len()), log.DEBUG)
			global.SemMulti <- 0
			sendPCBToReady()
			global.Logger.Log(fmt.Sprintf("PCB to READY - Semaforo %d - Multi: %d", len(global.SemMulti), global.KernelConfig.Multiprogramming), log.DEBUG)
		}
	}
}

func sendPCBToReady() {

	global.MutexNewState.Lock()
	pcbFront := global.NewState.Front()
	if pcbFront != nil  {
		pcbToReady := global.NewState.Remove(pcbFront).(*model.PCB)
		pcbToReady.State = "READY"
		global.MutexReadyState.Lock()
		global.ReadyState.PushBack(pcbToReady)
		global.MutexReadyState.Unlock()
	} else {
			global.Logger.Log("No PCB available to move to READY", log.DEBUG)
	}
	global.MutexNewState.Unlock()
}
