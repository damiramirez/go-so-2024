package algorithm

import (
	"fmt"
	"strings"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

func Fifo(){
	global.Logger.Log("Arranca FIFO", log.DEBUG)
	for {
		global.Logger.Log(fmt.Sprintf("global.SemExecute: %d\n", len(global.SemExecute)), log.DEBUG)
		global.SemExecute <- 0
		
		if !global.WorkingPlani {
			global.Logger.Log("TERMINO CON FIFO", log.DEBUG)
			break
		}

		global.Logger.Log(fmt.Sprintf("Cola READY: %d", global.ReadyState.Len()), log.DEBUG)

		if global.ReadyState.Len() != 0 {
			global.Logger.Log(fmt.Sprintf("PCB a execute: %+v", global.ReadyState.Front().Value), log.DEBUG)

			global.MutexReadyState.Lock()
			pcb := global.ReadyState.Front().Value.(*model.PCB)
			global.ReadyState.Remove(global.ReadyState.Front())
			global.MutexReadyState.Unlock()

			global.MutexExecuteState.Lock()
			global.ExecuteState.PushBack(pcb)
			global.MutexExecuteState.Unlock()

			updatePCB, _ := utils.PCBToCPU(pcb)

			global.Logger.Log(fmt.Sprintf("Recibi de CPU: %+v", updatePCB), log.DEBUG)

			global.MutexExecuteState.Lock()
			global.ExecuteState.Remove(global.ExecuteState.Front())
			global.MutexExecuteState.Unlock()

			// No hay mas instrucciones - EXIT
			if updatePCB.Instruction.Operation == "EXIT" {
				updatePCB.State = "EXIT"
				global.MutexExitState.Lock()
				global.ExitState.PushBack(updatePCB)
				global.MutexExitState.Unlock()
			}

			// Lo mando a BLOCK si instruccion tiene IO
			if strings.Contains(updatePCB.Instruction.Operation, "IO") {
				updatePCB.State = "BLOCK"
				global.MutexBlockState.Lock()
				global.BlockedState.PushBack(updatePCB)
				global.MutexBlockState.Unlock()
			}

			// VER ESTE SEMAFORO - DONDE VA?
			<- global.SemExecute
		}
	}
}
