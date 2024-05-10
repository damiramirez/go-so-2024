package utils

import (
	"container/list"
	"fmt"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/block"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/longterm"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
)

// Busca en todas las listas el PID
func FindProcessInList(pid int) *model.PCB {
	queues := []*list.List{
		global.NewState,
		global.ReadyState,
		global.ExecuteState,
		global.BlockedState,
		global.ExitState,
	}

	for _, queue := range queues {
		pcb := findProcess(pid, queue)
		if pcb != nil {
			return pcb
		}
	}

	return nil
}

func findProcess(pid int, list *list.List) *model.PCB {
	for e := list.Front(); e != nil; e = e.Next() {
		pcb := e.Value.(*model.PCB)
		if pid == pcb.PID {
			return pcb
		}
	}

	return nil
}

func GetAllProcess() []ProcessState {

	var allProcesses []ProcessState
	queues := []*list.List{
		global.NewState,
		global.ReadyState,
		global.ExecuteState,
		global.BlockedState,
		global.ExitState,
	}

	for _, queue := range queues {
		for e := queue.Front(); e != nil; e = e.Next() {
			pcb := e.Value.(*model.PCB)
			allProcesses = append(allProcesses, ProcessState{
				PID:   pcb.PID,
				State: pcb.State,
			},
			)
		}
	}

	return allProcesses
}

func RemoveProcessByPID(pid int) bool {

	queues := []*list.List{
		global.NewState,
		global.BlockedState,
		global.ExecuteState,
		global.ReadyState,
		global.ExitState,
	}

	for _, queue := range queues {
		for e := queue.Front(); e != nil; e = e.Next() {
			pcb := e.Value.(*model.PCB)

			if pcb.PID == pid {
				queue.Remove(e)
				<-global.SemMulti
				return true
			}
		}
	}

	return false
}

func PCBToCPU(pcb *model.PCB) (*model.PCB, error) {
	pcb.State = "EXEC"
	global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: READY - Estado Actual: %s", pcb.PID, pcb.State), log.INFO)

	resp, err := requests.PutHTTPwithBody[*model.PCB, model.PCB](
		global.KernelConfig.IPCPU, global.KernelConfig.PortCPU, "dispatch", pcb)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func PCBtoExit(updatePCB *model.PCB) {
	updatePCB.State = "EXIT"
	global.MutexExitState.Lock()
	global.ExitState.PushBack(updatePCB)
	global.MutexExitState.Unlock()
	//LOG CAMBIO DE ESTADO
	global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: EXEC - Estado Actual: %s", updatePCB.PID, updatePCB.State), log.INFO)
	global.Logger.Log(fmt.Sprintf("Finaliza el proceso %d - Motivo: SUCCESS ", updatePCB.PID), log.INFO)
	<-global.SemMulti
}

func PCBtoBlock(updatePCB *model.PCB) {
	updatePCB.State = "BLOCK"
	global.MutexBlockState.Lock()
	global.BlockedState.PushBack(updatePCB)
	global.MutexBlockState.Unlock()
	global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: EXEC - Estado Actual: %s", updatePCB.PID, updatePCB.State), log.INFO)
	global.Logger.Log(fmt.Sprintf("PID: %d - Bloqueado por: %s ", updatePCB.PID, updatePCB.Instruction.Parameters[0]), log.INFO)

	go block.ProcessToIO(updatePCB)
	
}

func PCBReadytoExec() *model.PCB {

	global.MutexReadyState.Lock()
	pcb := global.ReadyState.Front().Value.(*model.PCB)
	global.ReadyState.Remove(global.ReadyState.Front())
	global.MutexReadyState.Unlock()

	// Pasar a execute
	global.MutexExecuteState.Lock()
	global.ExecuteState.PushBack(pcb)
	global.MutexExecuteState.Unlock()
	//global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: READY - Estado Actual: %s", pcb.PID, pcb.State), log.INFO)

	return pcb
}

func PCBExectoReady(updatePCB *model.PCB) {
	//se guarda en ready
	updatePCB.State = "READY"
	//LOG CAMBIO DE ESTADO
	global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: EXEC - Estado Actual: %s", updatePCB.PID, updatePCB.State), log.INFO)

	//LOG COLA A READY CHEQUEAR EN ESTE CASO
	

	//LOG FIN DE QUANTUM
	global.Logger.Log(fmt.Sprintf("PID: %d - Desalojado por fin de Quantum ", updatePCB.PID), log.INFO)
	
	global.MutexReadyState.Lock()
	global.ReadyState.PushBack(updatePCB)

	global.MutexReadyState.Unlock()
	array := longterm.ConvertListToArray(global.ReadyState)
	global.Logger.Log(fmt.Sprintf("Cola Ready : %v", array), log.INFO)
	global.SemReadyList <- struct{}{}

}
