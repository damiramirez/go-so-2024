package algorithm

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	resource "github.com/sisoputnfrba/tp-golang/kernel/internal/Resource"

	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)


var updateChan = make(chan *model.PCB)


func VirtualRoundRobin(){
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
			
			InterruptTimer := make(chan int, 1)
			

			go VRRDisplaceFunction(InterruptTimer)

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

func VRRDisplaceFunction(InterruptTimer chan int) {

	<-global.SemInterrupt

	Time := time.Duration(global.KernelConfig.Quantum)
	global.Logger.Log(fmt.Sprintf("valor de time: %d ms",Time), log.DEBUG)

	timer := time.NewTimer(Time * time.Millisecond)
	StartTime:=time.Now()
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
		Pcb:=<-updateChan
		
		Elapsed:=time.Since(StartTime)
		ElapsedMili:=Elapsed.Milliseconds()
		TotalTime:=Time.Milliseconds()
		global.Logger.Log(fmt.Sprintf("Displace - Interrupt: Semaforo: %d", len(InterruptTimer)), log.DEBUG)
		if  Pcb.DisplaceReason=="BLOCKED"{
			
			global.Logger.Log(fmt.Sprintf("valor de elapsed: %d ms",ElapsedMili), log.DEBUG)
			

			RemainingTime:=TotalTime-ElapsedMili
			RoundedValue:=int(math.Round(float64(RemainingTime)))

			global.Logger.Log(fmt.Sprintf("valor total: %d ms",Time), log.DEBUG)
			global.Logger.Log(fmt.Sprintf("Displace - Interrupt: Blocked: %d ms",RoundedValue), log.DEBUG)
			global.Logger.Log(fmt.Sprintf("PCB BLOQUEADA: %+v ",Pcb), log.DEBUG)
		}
		timer.Stop()
		}
}
