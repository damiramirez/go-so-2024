package internal

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/cpu/global"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/execute"
	internal "github.com/sisoputnfrba/tp-golang/cpu/internal/fetch"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

func Dispatch(pcb *model.PCB) (*model.PCB, error) {
	global.Logger.Log(fmt.Sprintf("Recibi PCB %+v", pcb), log.DEBUG)

	global.ExecuteMutex.Lock()
	global.Execute=true
	global.ExecuteMutex.Unlock()
	
	for global.Execute {
		instruction, err := internal.Fetch(pcb)
		if err != nil {
			return nil, err
		}

		exec_result := execute.Execute(pcb, instruction)
		if exec_result == execute.RETURN_CONTEXT{
			global.Execute = false
		}
	}
	
	global.Logger.Log(fmt.Sprintf("PCB Actualizada %+v", pcb), log.DEBUG)
	return pcb, nil
}