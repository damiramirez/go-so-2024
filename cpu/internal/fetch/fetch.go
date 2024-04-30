package internal

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/cpu/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
)

func Fetch(pcb *model.PCB) (*model.Instruction, error) {
	global.Logger.Log(fmt.Sprintf("PID: %d -> Buscando instruccion: %d", pcb.PID, pcb.PC), log.DEBUG)
	instruction, err := getInstruction(pcb.PID, pcb.PC)
	if err != nil {
		global.Logger.Log("Error al obtener la instruccion: "+err.Error(), log.ERROR)
		return nil, err
	}

	pcb.PC++

	global.Logger.Log(fmt.Sprintf("Actualizo PCB: %+v", pcb), log.DEBUG)

	return instruction, err
}

func getInstruction(id, address int) (*model.Instruction, error) {
	path := fmt.Sprintf("process/%d/instructions/%d", id, address)
	instruction, err := requests.GetHTTP[model.Instruction](
		global.CPUConfig.IPMemory,
		global.CPUConfig.PortKernel,
		path,
	)

	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al solicitar instrucción desde memoria: %v", err), log.ERROR)
		return nil, err
	}

	global.Logger.Log(fmt.Sprintf("Instruction: %+v", instruction), log.DEBUG)
	return instruction, nil
}
