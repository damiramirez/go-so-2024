package algorithm

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	resource "github.com/sisoputnfrba/tp-golang/kernel/internal/resources"

	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

func RoundRobbin() {
	global.Logger.Log(fmt.Sprintf("Semaforo de SemReadyList INICIO: %d", len(global.SemReadyList)), log.DEBUG)
	displaceMap = make(map[int]*model.PCB)

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

			go DisplaceFunction(InterruptTimer, pcb)

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

			if updatePCB.Instruction.Operation == "EXIT" {
				// global.Logger.Log(fmt.Sprintf("EXIT - Antes de Interrupt. Semaforo: %d", len(interruptTimer)), log.DEBUG)
				InterruptTimer <- 0
				// global.Logger.Log(fmt.Sprintf("EXIT - Despues de Interrupt. Semaforo: %d", len(InterruptTimer)), log.DEBUG)
				utils.PCBtoExit(updatePCB)
			}
			if updatePCB.DisplaceReason == "BLOCKED" {
				InterruptTimer <- 0
				DisplaceChan <- updatePCB
				utils.PCBtoBlock(updatePCB)
			} else if updatePCB.DisplaceReason == "QUANTUM" && updatePCB.Instruction.Operation != "EXIT" {
				if updatePCB.Instruction.Operation == "SIGNAL" {
					resource.Signal(updatePCB)
				} else if updatePCB.Instruction.Operation == "WAIT" {
					resource.Wait(updatePCB)
				} else if strings.Contains(updatePCB.Instruction.Operation, "IO") {
					utils.PCBtoBlock(updatePCB)
				} else {
					utils.PCBExectoReady(updatePCB)
				}
			}

			if updatePCB.DisplaceReason == "WAIT" {
				
				InterruptTimer <- 0

				global.Logger.Log("antes de displace chan", log.DEBUG)
				DisplaceChan <-updatePCB
				global.Logger.Log("despues de displace chan", log.DEBUG)

				resource.Wait(updatePCB)
			}
			if updatePCB.DisplaceReason == "SIGNAL" {
				InterruptTimer <- 0

				global.Logger.Log("antes de displace chan", log.DEBUG)
				DisplaceChan <-updatePCB
				global.Logger.Log("despues de displace chan", log.DEBUG)

				resource.Signal(updatePCB)
			}
		}

		<-global.SemExecute
	}
}

func DisplaceFunction(InterruptTimer chan int, OldPcb *model.PCB) {


	<-global.SemInterrupt
	global.Logger.Log(fmt.Sprintf("pcb antes de select: %+v", OldPcb), log.DEBUG)


	quantumTime := time.Duration(OldPcb.RemainingQuantum) * time.Millisecond

	timer := time.NewTimer(quantumTime)

	defer timer.Stop()

	startTime := time.Now()

	select {
	case <-timer.C:

		global.Logger.Log(fmt.Sprintf("PID: %d Displace - Termino timer.C", OldPcb.PID), log.DEBUG)
		url := fmt.Sprintf("http://%s:%d/%s", global.KernelConfig.IPCPU, global.KernelConfig.PortCPU, "interrupt")
		_, err := http.Get(url)
		if err != nil {
			global.Logger.Log(fmt.Sprintf("Error al enviar la interrupción: %v", err), log.ERROR)
			return
		}
	case <-InterruptTimer:

		timer.Stop()

		pcb := <-DisplaceChan
		// Transformar el tiempo a segundos para redondearlo y despues pasarlo a ms
		// Asi uso los ms en la PCB

		if pcb.Instruction.Operation=="WAIT"||pcb.Instruction.Operation=="SIGNAL" {
			
		remainingMillisRounded:=utils.TimeCalc(startTime,quantumTime,pcb)
		

		if remainingMillisRounded > 0 {
			pcb.RemainingQuantum = remainingMillisRounded
		} else {
			pcb.RemainingQuantum = global.KernelConfig.Quantum
		}
	}
}

}

