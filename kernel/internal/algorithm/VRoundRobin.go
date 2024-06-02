package algorithm

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	resource "github.com/sisoputnfrba/tp-golang/kernel/internal/resources"

	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)


var updateChan = make(chan *model.PCB)
var displaceChan = make(chan *model.PCB)


func VirtualRoundRobin(){
	global.Logger.Log(fmt.Sprintf("Semaforo de SemReadyList INICIO: %d", len(global.SemReadyList)), log.DEBUG)

	for {
		<-global.SemReadyList
		global.SemExecute <- 0

		if !global.WorkingPlani {
			global.Logger.Log("TERMINO CON VRR", log.DEBUG)
			<-global.SemExecute
			break
		}

		if global.ReadyState.Len() != 0 {
			global.Logger.Log(fmt.Sprintf("PCB a execute: %+v", global.ReadyState.Front().Value), log.DEBUG)

			pcb := utils.PCBReadytoExec()
			
			interruptTimer := make(chan int, 1)

			go VRRDisplaceFunction(interruptTimer)

			go func() {
				global.SemInterrupt <- 0
				updatePCB, _ = utils.PCBToCPU(pcb)

				updateChan <- updatePCB
			}()

			updatePCB = <-updateChan
			//LOG CAMBIO DE ESTADO
			global.Logger.Log(fmt.Sprintf("Recibi de CPU: %+v", updatePCB), log.DEBUG)

			// Sacar de execute
			global.MutexExecuteState.Lock()
			global.ExecuteState.Remove(global.ExecuteState.Front())
			global.MutexExecuteState.Unlock()

			// EXIT - Agregar a exit
			if updatePCB.DisplaceReason == "EXIT" {
				global.Logger.Log(fmt.Sprintf("EXIT - Antes de Interrupt. Semaforo: %d", len(interruptTimer)), log.DEBUG)
				interruptTimer <- 0
				global.Logger.Log(fmt.Sprintf("EXIT - Despues de Interrupt. Semaforo: %d", len(interruptTimer)), log.DEBUG)
				utils.PCBtoExit(updatePCB)
			}
			// Agregar a block
			if updatePCB.DisplaceReason == "BLOCKED" {
				global.Logger.Log(fmt.Sprintf("BLOCKED - Antes de Interrupt. Semaforo: %d", len(interruptTimer)), log.DEBUG)
				interruptTimer <- 0
				displaceChan <- updatePCB
				global.Logger.Log(fmt.Sprintf("BLOCKED - Despues de Interrupt. Semaforo: %d", len(interruptTimer)), log.DEBUG)
				utils.PCBtoBlock(updatePCB)
			}
			if updatePCB.DisplaceReason == "QUANTUM" {
				utils.PCBExectoReady(updatePCB)
			}

			if updatePCB.DisplaceReason == "WAIT" {
				resource.Wait(updatePCB)
			}
			if updatePCB.DisplaceReason == "SIGNAL" {
				resource.Signal(updatePCB)
			}
		}

		<-global.SemExecute
	}
}

func VRRDisplaceFunction(interruptTimer chan int) {

	<-global.SemInterrupt

	quantumTime := time.Duration(global.KernelConfig.Quantum) * time.Millisecond
	timer := time.NewTimer(quantumTime)
	defer timer.Stop()

	startTime := time.Now()

	select {
	case <-timer.C:
		global.Logger.Log("Displace - Termino timer.C", log.DEBUG)
		url := fmt.Sprintf("http://%s:%d/%s", global.KernelConfig.IPCPU, global.KernelConfig.PortCPU, "interrupt")
		_, err := http.Get(url)
		if err != nil {
			global.Logger.Log(fmt.Sprintf("Error al enviar la interrupción: %v", err), log.ERROR)
			return
		}

	case <-interruptTimer:
		timer.Stop()
		pcb := <-displaceChan
		global.Logger.Log(fmt.Sprintf("PCB EN DISPLACE: %+v", pcb), log.DEBUG)
		if pcb.DisplaceReason == "BLOCKED" {

			elapsedTime := time.Since(startTime)
			elapsedMillis := elapsedTime.Milliseconds()
			global.Logger.Log(fmt.Sprintf("ElapsedTime: %d ms", elapsedMillis), log.DEBUG)
			remainingQuantum := quantumTime - elapsedTime
			remainingMillis := remainingQuantum.Milliseconds()
			global.Logger.Log(fmt.Sprintf("RemainingTime: %d ms", remainingMillis), log.DEBUG)

			global.Logger.Log(fmt.Sprintf("PCB BLOQUEADA: %+v ", pcb), log.DEBUG)
		}
	}
}