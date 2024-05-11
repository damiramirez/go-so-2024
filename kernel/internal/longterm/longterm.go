package longterm

import (
	"container/list"
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
		global.Logger.Log(fmt.Sprintf("Semaforo de NEWLIST ANTES DE PLANO: %d", len(global.SemNewList)), log.DEBUG)

		<-global.SemNewList
		global.Logger.Log(fmt.Sprintf("Semaforo de NEWLIST DESPIUES DE PLANI: %d", len(global.SemNewList)), log.DEBUG)

		if global.NewState.Len() != 0 {
			global.Logger.Log(fmt.Sprintf("NEW LEN: %d", global.NewState.Len()), log.DEBUG)
			global.SemMulti <- 0
			sendPCBToReady()
			array := ConvertListToArray(global.ReadyState)
			global.Logger.Log(fmt.Sprintf("PCB to READY - Semaforo %d - Multi: %d", len(global.SemMulti), global.KernelConfig.Multiprogramming), log.DEBUG)
			global.Logger.Log(fmt.Sprintf("Cola Ready : %v", array), log.INFO)
		}
	}
}

// funcion que cree para agarrar una lista de tipo list list a slice de interface
func ConvertListToArray(l *list.List) []interface{} {
	array := make([]interface{}, l.Len())
	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		array[i] = e.Value.(*model.PCB).PID
		i++
	}
	return array
}
func sendPCBToReady() {

	pcbFront := global.NewState.Front()
	if pcbFront != nil {
		global.MutexNewState.Lock()
		pcbToReady := global.NewState.Remove(pcbFront).(*model.PCB)
		global.MutexNewState.Unlock()

		pcbToReady.State = "READY"

		global.MutexReadyState.Lock()
		global.ReadyState.PushBack(pcbToReady)
		global.MutexReadyState.Unlock()

		// <- global.SemReadyList
		global.SemReadyList <- struct{}{}
	} else {
		global.Logger.Log("No PCB available to move to READY", log.DEBUG)
	}
}
