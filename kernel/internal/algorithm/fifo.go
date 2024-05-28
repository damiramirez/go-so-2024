package algorithm

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/kernel/global"

	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

var updatePCB *model.PCB

func Fifo() {
	global.Logger.Log("Arranca FIFO", log.DEBUG)

	for {

		global.Logger.Log("LOG ANTES DE SEMREADYLIST", log.DEBUG)
		<-global.SemReadyList

		global.SemExecute <- 0

		if !global.WorkingPlani {
			global.Logger.Log("TERMINO CON FIFO", log.DEBUG)
			break
		}

		if global.ReadyState.Len() != 0 {
			global.Logger.Log(fmt.Sprintf("PCB a execute: %+v", global.ReadyState.Front().Value), log.DEBUG)

			pcb := utils.PCBReadytoExec()

			updateChan := make(chan *model.PCB)
			go func() {
				updatePCB, _ = utils.PCBToCPU(pcb)
				updateChan <- updatePCB
			}()
			updatePCB = <-updateChan
			global.Logger.Log(fmt.Sprintf("Recibi de CPU: %+v", updatePCB), log.DEBUG)

			// Sacar de execute
			global.MutexExecuteState.Lock()
			global.ExecuteState.Remove(global.ExecuteState.Front())
			global.MutexExecuteState.Unlock()

			// EXIT - Agregar a exit
			if updatePCB.DisplaceReason == "EXIT" {
				utils.PCBtoExit(updatePCB)
			}

			// Agregar a block
			if updatePCB.DisplaceReason == "BLOCKED" {
				utils.PCBtoBlock(updatePCB)
			}

		}

		<-global.SemExecute
	}
}
