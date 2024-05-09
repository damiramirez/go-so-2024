package algorithm

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/block"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"

)

func RoundRobbin() {
	global.Logger.Log("Arranca RoundRobbin", log.DEBUG)

	// TODO: Mover codigo

	for {

		// TODO: ESPERA ACTIVA? BUCLE INFINITO - VER SEMAFOROS
		global.Logger.Log("LOG ANTES DE SEMREADYLIST", log.DEBUG)
		<-global.SemReadyList
		global.Logger.Log(fmt.Sprintf("READY: %d", global.ReadyState.Len()), log.DEBUG)

		global.SemExecute <- 0

		if !global.WorkingPlani {
			global.Logger.Log("TERMINO CON ROUND ROBIN", log.DEBUG)
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
			updateChan := make(chan *model.PCB)
			go DisplaceFunction()

			go func() {
				global.Logger.Log("se bloquea funcion",log.DEBUG)
				global.SemInterrupt <- 0
				global.Logger.Log("se libera funcion",log.DEBUG)
				updatePCB, _ = utils.PCBToCPU(pcb)
				updateChan <- updatePCB
			}()
			
			updatePCB = <-updateChan
			global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: READY - Estado Actual: %s", pcb.PID, pcb.State), log.INFO)
			global.Logger.Log(fmt.Sprintf("Recibi de CPU: %+v", updatePCB), log.DEBUG)

			// Sacar de execute
			global.MutexExecuteState.Lock()
			global.ExecuteState.Remove(global.ExecuteState.Front())
			global.MutexExecuteState.Unlock()

			// EXIT - Agregar a exit
			if updatePCB.DisplaceReason == "EXIT" {
				updatePCB.State = "EXIT"
				global.MutexExitState.Lock()
				global.ExitState.PushBack(updatePCB)
				global.MutexExitState.Unlock()
				global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: EXEC - Estado Actual: %s", updatePCB.PID, updatePCB.State), log.INFO)
				<-global.SemMulti
			}

			// Agregar a block
			if 	updatePCB.DisplaceReason=="BLOCKED"{
				updatePCB.State = "BLOCK"
				global.MutexBlockState.Lock()
				global.BlockedState.PushBack(updatePCB)
				global.MutexBlockState.Unlock()
				global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: EXEC - Estado Actual: %s", updatePCB.PID, updatePCB.State), log.INFO)

				go block.ProcessToIO()
			}
			if updatePCB.DisplaceReason=="QUANTUM"{
				//se guarda en ready	
				updatePCB.State = "READY"
				global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: EXEC - Estado Actual: %s", updatePCB.PID, updatePCB.State), log.INFO)
				global.MutexReadyState.Lock()
				global.ReadyState.PushBack(updatePCB)
				global.MutexReadyState.Unlock()
				global.SemReadyList <- struct{}{}
			}

		}
		// VER ESTE SEMAFORO - DONDE VA?
		<-global.SemExecute
	}
}

func DisplaceFunction(){
	global.Logger.Log("EJECUTE DISPLACE",log.DEBUG)
	<-global.SemInterrupt
	Time:=time.Duration(global.KernelConfig.Quantum)
	time.Sleep(Time * time.Millisecond)
	
	url := fmt.Sprintf("http://%s:%d/%s", global.KernelConfig.IPCPU, global.KernelConfig.PortCPU, "interrupt")
	_, err := http.Get(url)
	if err != nil {
		return 
	}
	
} 