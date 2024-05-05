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

	// TODO: Mover codigo
	for {
		// TODO: ESPERA ACTIVA? BUCLE INFINITO - VER SEMAFOROS
		global.SemReadyList <- 0

		global.SemExecute <- 0

		if !global.WorkingPlani {
			global.Logger.Log("TERMINO CON FIFO", log.DEBUG)
			break
		}

		if global.ReadyState.Len() != 0 {
			global.Logger.Log(fmt.Sprintf("PCB a execute: %+v", global.ReadyState.Front().Value), log.DEBUG)

			// Sacar de ready
			global.MutexReadyState.Lock()
			pcb := global.ReadyState.Front().Value.(*model.PCB)
			global.ReadyState.Remove(global.ReadyState.Front())
			global.MutexReadyState.Unlock()

			// Pasar a execute
			global.MutexExecuteState.Lock()
			global.ExecuteState.PushBack(pcb)
			global.MutexExecuteState.Unlock()

			// Enviar a execute
			updatePCB, _ := utils.PCBToCPU(pcb)
			global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: READY - Estado Actual: %s", pcb.PID, pcb.State), log.INFO)
			global.Logger.Log(fmt.Sprintf("Recibi de CPU: %+v", updatePCB), log.DEBUG)
			

			// Sacar de execute
			global.MutexExecuteState.Lock()
			global.ExecuteState.Remove(global.ExecuteState.Front())
			global.MutexExecuteState.Unlock()

			// EXIT - Agregar a exit
			if updatePCB.Instruction.Operation == "EXIT" {
				updatePCB.State = "EXIT"
				global.MutexExitState.Lock()
				global.ExitState.PushBack(updatePCB)
				global.MutexExitState.Unlock()
				global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: EXEC - Estado Actual: %s", updatePCB.PID, updatePCB.State), log.INFO)
			}

			// Agregar a block
			if strings.Contains(updatePCB.Instruction.Operation, "IO") {
				updatePCB.State = "BLOCK"
				global.MutexBlockState.Lock()
				global.BlockedState.PushBack(updatePCB)
				global.MutexBlockState.Unlock()
				global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: EXEC - Estado Actual: %s", updatePCB.PID, updatePCB.State), log.INFO)
			}

		}
		// VER ESTE SEMAFORO - DONDE VA?
		<- global.SemExecute
	}
}
