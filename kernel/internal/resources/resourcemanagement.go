package resources

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

func Wait(Pcb *model.PCB) {
	resource := global.ResourceMap[Pcb.Instruction.Parameters[0]]
	resource.Count -= 1
	resource.PidList = append(resource.PidList, Pcb.PID)
	global.Logger.Log(fmt.Sprintf("Recurso: %s - Cantidad instancias: %d", resource.Name, resource.Count), log.DEBUG)
	
	if resource.Count < 0 {
		resource.MutexList.Lock()
		resource.BlockedList.PushBack(Pcb)
		resource.MutexList.Unlock()

		//poner en listar procesos
		global.Logger.Log(fmt.Sprintf("Bloqueo proceso: %d", Pcb.PID), log.DEBUG)
	}
	if Pcb.DisplaceReason=="QUANTUM" {
		Pcb.RemainingQuantum=global.KernelConfig.Quantum
		global.MutexReadyState.Lock()
		global.ReadyState.PushBack(Pcb)
		global.MutexReadyState.Unlock()
		global.Logger.Log(fmt.Sprintf("Envio PID %d ultimo a ready", Pcb.PID), log.DEBUG)
	}else {
		global.MutexReadyState.Lock()
		global.ReadyState.PushFront(Pcb)
		global.MutexReadyState.Unlock()

		global.Logger.Log(fmt.Sprintf("Envio PID %d primero a Ready", Pcb.PID), log.DEBUG)
		
	}
	global.SemReadyList <- struct{}{}
}

func Signal(PcbExec *model.PCB) {
	resource := global.ResourceMap[PcbExec.Instruction.Parameters[0]]
	resource.Count += 1
	global.Logger.Log(fmt.Sprintf("%s %d", resource.Name, resource.Count), log.DEBUG)

	if resource.BlockedList.Len() > 0 {
		resource.MutexList.Lock()
		PCBBlock := resource.BlockedList.Front().Value.(*model.PCB)
		resource.BlockedList.Remove(resource.BlockedList.Front())
		resource.MutexList.Unlock()

		global.MutexReadyState.Lock()
		global.ReadyState.PushBack(PCBBlock)
		global.MutexReadyState.Unlock()

		global.Logger.Log(fmt.Sprintf("Envio PID %d al fondo de ready", PCBBlock.PID), log.DEBUG)
		global.SemReadyList <- struct{}{}
	}

	value := checkInArray(resource.PidList, PcbExec.PID)

	if value != -1 {
		resource.PidList = removeAt(resource.PidList, value)
	}
	if PcbExec.DisplaceReason=="QUANTUM" {
		PcbExec.RemainingQuantum=global.KernelConfig.Quantum
		global.MutexReadyState.Lock()
		global.ReadyState.PushBack(PcbExec)
		global.MutexReadyState.Unlock()
		global.Logger.Log(fmt.Sprintf("Envio PID %d ultimo a ready", PcbExec.PID), log.DEBUG)
	}else {
		global.MutexReadyState.Lock()
		global.ReadyState.PushFront(PcbExec)
		global.MutexReadyState.Unlock()

		global.Logger.Log(fmt.Sprintf("Envio PID %d primero a Ready", PcbExec.PID), log.DEBUG)
		
	}
	global.SemReadyList <- struct{}{}
}

func checkInArray(resourcesIdds []int, pid int) int {
	for i, value := range resourcesIdds {
		if value == pid {
			return i
		}
	}
	return -1
}

func removeAt(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...)
}
