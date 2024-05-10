package algorithm

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

func RoundRobbin() {
	global.Logger.Log("Arranca RoundRobbin", log.DEBUG)

	// TODO: Mover codigo

	for {

		global.Logger.Log("LOG ANTES DE SEMREADYLIST", log.DEBUG)
		<-global.SemReadyList
		global.SemExecute <- 0

		if !global.WorkingPlani {
			global.Logger.Log("TERMINO CON ROUND ROBIN", log.DEBUG)
			break
		}

		if global.ReadyState.Len() != 0 {
			global.Logger.Log(fmt.Sprintf("PCB a execute: %+v", global.ReadyState.Front().Value), log.DEBUG)

			pcb:= utils.PCBReadytoExec()
			// Enviar a execute
			updateChan := make(chan *model.PCB)
			InterruptTimer := make(chan int)

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
				global.Logger.Log("antes del timerinterrupt",log.DEBUG)
				InterruptTimer <- 0
				global.Logger.Log("despues del timerinterrupt",log.DEBUG)
				utils.PCBtoExit(updatePCB)
				
			}

			// Agregar a block
			if updatePCB.DisplaceReason == "BLOCKED" {
				InterruptTimer <- 0 
				
				utils.PCBtoBlock(updatePCB)
				
			}
			if updatePCB.DisplaceReason == "QUANTUM" {
				utils.PCBExectoReady(updatePCB)
				}

		}
		// VER ESTE SEMAFORO - DONDE VA?
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
		global.Logger.Log("EJECUTE DISPLACE", log.DEBUG)
		url := fmt.Sprintf("http://%s:%d/%s", global.KernelConfig.IPCPU, global.KernelConfig.PortCPU, "interrupt")
		_, err := http.Get(url)
		if err != nil {
			return
		}
	case <-InterruptTimer:
		global.Logger.Log("CORTE EL TIMER", log.DEBUG)
		
	}
	
}
