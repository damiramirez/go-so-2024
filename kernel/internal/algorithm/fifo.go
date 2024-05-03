package algorithm

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

func Fifo() (*model.PCB, error) {
	global.Logger.Log("Arranca FIFO", log.DEBUG)
	for {
		global.SemExecute <- 0
		if !global.WorkingPlani {
			break
		}

		if global.ReadyState.Len() != 0 {
			global.Logger.Log(fmt.Sprintf("PCB a execute: %+v", global.ReadyState.Front().Value), log.DEBUG)
			global.MutexReadyState.Lock()
			pcb := global.ReadyState.Front().Value.(*model.PCB)
			global.ReadyState.Remove(global.ReadyState.Front())
			global.MutexReadyState.Unlock()
			updatePCB, err := utils.PCBToCPU(pcb)
			if err != nil {
				return nil, err
			}
			return updatePCB, nil
		}
	}

	return nil, nil
}