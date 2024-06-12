package algorithm

import (
	"container/list"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	resource "github.com/sisoputnfrba/tp-golang/kernel/internal/resources"

	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

var updateChan = make(chan *model.PCB)
var DisplaceList *list.List = list.New()
var MutexDisplaceList sync.Mutex
var DisplaceChan = make(chan *model.PCB)

func VirtualRoundRobin() {


	var pcb *model.PCB
	for {
		global.Logger.Log("Log antes de SemReadyList", log.DEBUG)
		<-global.SemReadyList
		global.SemExecute <- 0

		if !global.WorkingPlani {
			global.Logger.Log("TERMINO CON VRR", log.DEBUG)
			<-global.SemExecute
			break
		}

		if global.ReadyState.Len() != 0 || global.ReadyPlus.Len() != 0 {

			if global.ReadyPlus.Len() != 0 {
				global.Logger.Log(fmt.Sprintf("PCB a execute: %+v", global.ReadyPlus.Front().Value), log.DEBUG)
				pcb = utils.VrrPCBtoEXEC()

			} else {
				global.Logger.Log(fmt.Sprintf("PCB a execute: %+v", global.ReadyState.Front().Value), log.DEBUG)
				pcb = utils.PCBReadytoExec()
			}

			interruptTimer := make(chan int, 1)

			go VRRDisplaceFunction(interruptTimer, pcb)

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
			if updatePCB.Instruction.Operation == "EXIT" {
				// global.Logger.Log(fmt.Sprintf("EXIT - Antes de Interrupt. Semaforo: %d", len(interruptTimer)), log.DEBUG)
				interruptTimer <- 0
				DisplaceChan <-updatePCB
				// global.Logger.Log(fmt.Sprintf("EXIT - Despues de Interrupt. Semaforo: %d", len(interruptTimer)), log.DEBUG)
				utils.PCBtoExit(updatePCB)
			}
			if updatePCB.DisplaceReason == "BLOCKED" {
				interruptTimer <- 0
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
				
				interruptTimer <- 0

				global.Logger.Log("antes de displace chan", log.DEBUG)
				DisplaceChan <-updatePCB
				global.Logger.Log("despues de displace chan", log.DEBUG)

				resource.Wait(updatePCB)
			}
			if updatePCB.DisplaceReason == "SIGNAL" {
				interruptTimer <- 0

				global.Logger.Log("antes de displace chan", log.DEBUG)
				DisplaceChan <-updatePCB
				global.Logger.Log("despues de displace chan", log.DEBUG)

				resource.Signal(updatePCB)
			}

		}

		<-global.SemExecute
	}
}

func VRRDisplaceFunction(interruptTimer chan int, OldPcb *model.PCB) {

	<-global.SemInterrupt
	
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
	case <-interruptTimer:

		timer.Stop()

		pcb := <-DisplaceChan

		if pcb.Instruction.Operation=="EXIT" {
			return
		}
		global.Logger.Log(fmt.Sprintf("PCB EN DISPLACE: %+v", pcb), log.DEBUG)

		// Transformar el tiempo a segundos para redondearlo y despues pasarlo a ms
		// Asi uso los ms en la PCB
		remainingMillisRounded:=utils.TimeCalc(startTime,quantumTime,pcb)

		if remainingMillisRounded > 0 {
			pcb.RemainingQuantum = remainingMillisRounded
		} else {
			pcb.RemainingQuantum = global.KernelConfig.Quantum
		}

	
	}
}
