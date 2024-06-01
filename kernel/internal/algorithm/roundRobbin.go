package algorithm

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	resource "github.com/sisoputnfrba/tp-golang/kernel/internal/Resource"

	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

func RoundRobbin() {
	global.Logger.Log(fmt.Sprintf("Semaforo de SemReadyList INICIO: %d", len(global.SemReadyList)), log.DEBUG)

	for {

		<-global.SemReadyList
		global.SemExecute <- 0

		if !global.WorkingPlani {
			global.Logger.Log("TERMINO CON ROUND ROBIN", log.DEBUG)
			<-global.SemExecute
			break
		}

		if global.ReadyState.Len() != 0 {
			global.Logger.Log(fmt.Sprintf("PCB a execute: %+v", global.ReadyState.Front().Value), log.DEBUG)

			pcb := utils.PCBReadytoExec()
			// Enviar a execute
			updateChan := make(chan *model.PCB)
			InterruptTimer := make(chan int, 1)

			go DisplaceFunction(InterruptTimer)

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
				global.Logger.Log(fmt.Sprintf("EXIT - Antes de Interrupt. Semaforo: %d", len(InterruptTimer)), log.DEBUG)
				InterruptTimer <- 0
				global.Logger.Log(fmt.Sprintf("EXIT - Despues de Interrupt. Semaforo: %d", len(InterruptTimer)), log.DEBUG)
				utils.PCBtoExit(updatePCB)
			}
			// Agregar a block
			if updatePCB.DisplaceReason == "BLOCKED" {
				global.Logger.Log(fmt.Sprintf("BLOCKED - Antes de Interrupt. Semaforo: %d", len(InterruptTimer)), log.DEBUG)
				InterruptTimer <- 0
				global.Logger.Log(fmt.Sprintf("BLOCKED - Despues de Interrupt. Semaforo: %d", len(InterruptTimer)), log.DEBUG)
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

func DisplaceFunction(InterruptTimer chan int) {

	<-global.SemInterrupt

	Time := time.Duration(global.KernelConfig.Quantum)
	//time.Sleep(Time * time.Millisecond)
	timer := time.NewTimer(Time * time.Millisecond)
	defer timer.Stop()

	select {
	case <-timer.C:
		global.Logger.Log("Displace - Termino timer.C", log.DEBUG)
		url := fmt.Sprintf("http://%s:%d/%s", global.KernelConfig.IPCPU, global.KernelConfig.PortCPU, "interrupt")
		_, err := http.Get(url)
		if err != nil {
			return
		}
	case <-InterruptTimer:
		global.Logger.Log(fmt.Sprintf("Displace - Interrupt: Semaforo: %d", len(InterruptTimer)), log.DEBUG)
		timer.Stop()
}
}