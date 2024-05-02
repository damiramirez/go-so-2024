package algorithm

import (
	"github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

func Fifo() {
	for {
		global.SemExecute <- 0
		if !global.WorkingPlani {
			break
		}

		if global.ReadyState.Len() != 0 {
			global.MutexReadyState.Lock()
			pcb := global.ReadyState.Front().Value.(*model.PCB)
			global.ReadyState.Remove(global.ReadyState.Front())
			global.MutexReadyState.Unlock()
			utils.PCBToCPU(pcb)
		}
	}
}