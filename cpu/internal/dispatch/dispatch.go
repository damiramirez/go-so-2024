package internal

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/cpu/global"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/execute"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/fetch"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

func Dispatch(pcb *model.PCB) (*model.PCB, error) {
	global.Logger.Log(fmt.Sprintf("Recibi PCB %+v", pcb), log.DEBUG)

	executing := true

	for executing {
		instruction, err := fetch.Fetch(pcb)
		if err != nil {
			return nil, err
		}

		execute.Execute(pcb, instruction)
	}
	
	return pcb, nil
}