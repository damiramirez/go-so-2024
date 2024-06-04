package algorithm

import (
	"container/list"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	resource "github.com/sisoputnfrba/tp-golang/kernel/internal/resources"

	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

var updateChan = make(chan *model.PCB)
var displaceChan = make(chan *model.PCB)
var DisplaceList *list.List = list.New()
var MutexDisplaceList sync.Mutex

func VirtualRoundRobin() {
	
	global.Logger.Log(fmt.Sprintf("Semaforo de SemReadyList INICIO: %d", len(global.SemReadyList)), log.DEBUG)
	var pcb *model.PCB
	for {
		global.Logger.Log(fmt.Sprintf("Semaforo de SemReadyList antesss: %d", len(global.SemReadyList)), log.DEBUG)
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
				global.Logger.Log(fmt.Sprintf("en medio de semaforos: %d", len(interruptTimer)), log.DEBUG)
				global.Logger.Log(fmt.Sprintf("antes de chan: %d", len(displaceChan)), log.DEBUG)
				//displaceChan <- updatePCB
				MutexDisplaceList.Lock()
				DisplaceList.PushBack(updatePCB)
				MutexDisplaceList.Unlock()
				global.Logger.Log(fmt.Sprintf("despues de chan: %d", len(displaceChan)), log.DEBUG)
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

func VRRDisplaceFunction(interruptTimer chan int, OldPcb *model.PCB) {

	<-global.SemInterrupt
	quantumTime := time.Duration(OldPcb.RemainingQuantum) * time.Millisecond
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
		global.Logger.Log("antes de displace chan: ", log.DEBUG)
		pcb := DisplaceList.Front().Value.(*model.PCB)
		global.ReadyState.Remove(DisplaceList.Front())

		global.Logger.Log(fmt.Sprintf("PCB EN DISPLACE: %+v", pcb), log.DEBUG)

		
		if pcb.DisplaceReason == "BLOCKED" {

			// Transformar el tiempo a segundos para redondearlo y despues pasarlo a ms
			// Asi uso los ms en la PCB
			elapsedTime := time.Since(startTime)
			elapsedSeconds := math.Round(elapsedTime.Seconds())
			elapsedMillisRounded := int64(elapsedSeconds * 1000)
			global.Logger.Log("estoy dentro de block", log.DEBUG)
			remainingQuantum := quantumTime - elapsedTime
			remainingSeconds := math.Round(remainingQuantum.Seconds())
			remainingMillisRounded := int64(remainingSeconds * 1000)

			global.Logger.Log(fmt.Sprintf("Rounded ElapsedTime: %d ms", elapsedMillisRounded), log.DEBUG)
			global.Logger.Log(fmt.Sprintf("Rounded RemainingTime: %d ms", remainingMillisRounded), log.DEBUG)

			if remainingMillisRounded > 0 {
				pcb.RemainingQuantum = int(remainingMillisRounded)
			} else {
				pcb.RemainingQuantum = global.KernelConfig.Quantum
			}

			global.Logger.Log(fmt.Sprintf("PCB BLOQUEADA: %+v ", pcb), log.DEBUG)

		}
	}
}
