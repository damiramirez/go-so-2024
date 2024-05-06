package utils

import (
	"container/list"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
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
	resp, err := requests.PutHTTPwithBody[*model.PCB, model.PCB](
		global.KernelConfig.IPCPU, global.KernelConfig.PortCPU, "dispatch", pcb)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
