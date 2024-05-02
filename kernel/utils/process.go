package utils

import (
	"container/list"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

// Busca en todas las listas el PID
func findProcessInList(pid int) {
	queues := []*list.List{global.NewState, global.ReadyState, global.RunningState, global.BlockedState}

	for _, queue := range queues {
		pcb := findProcess(pid, queue)
	}
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
