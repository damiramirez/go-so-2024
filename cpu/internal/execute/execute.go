package execute

import (
	"fmt"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

// TODO: IO_GEN_SLEEP

const (
	CONTINUE       = 0
	RETURN_CONTEXT = 1
)

// Ejecuto -> sumo PC en dispatch?
func Execute(pcb *model.PCB, instruction *model.Instruction) int {
	result := 0
	switch instruction.Operation {
	case "SET":
		set(pcb, instruction)
		result = CONTINUE
	case "SUM":
		sum(pcb, instruction)
		result = CONTINUE
	case "SUB":
		sub(pcb, instruction)
		result = CONTINUE
	case "JNZ":
		jnz(pcb, instruction)
		result = CONTINUE
	case "IO_GEN_SLEEP":
		result = RETURN_CONTEXT
	case "IO_STDIN_READ":
		result = RETURN_CONTEXT
	case "EXIT":
		result = RETURN_CONTEXT
	}

	global.Logger.Log(
		fmt.Sprintf("PID: %d - Ejecutando: %s - %+v",
			pcb.PID,
			instruction.Operation,
			instruction.Parameters,
		),
		log.INFO)

	pcb.Instruction = *instruction

	return result
}

func set(pcb *model.PCB, instruction *model.Instruction) {
	value, _ := strconv.Atoi(instruction.Parameters[1])
	setRegister(instruction.Parameters[0], value, pcb)
}

func sum(pcb *model.PCB, instruction *model.Instruction) {

	destinationValue := getRegister(instruction.Parameters[0], pcb)
	sourceValue := getRegister(instruction.Parameters[1], pcb)
	destinationValue = destinationValue + sourceValue
	setRegister(instruction.Parameters[0], destinationValue, pcb)
}

func sub(pcb *model.PCB, instruction *model.Instruction) {

	destinationValue := getRegister(instruction.Parameters[0], pcb)
	sourceValue := getRegister(instruction.Parameters[1], pcb)
	destinationValue = destinationValue - sourceValue
	setRegister(instruction.Parameters[0], destinationValue, pcb)
}

func jnz(pcb *model.PCB, instruction *model.Instruction) {
	value := getRegister(instruction.Parameters[0], pcb)
	if value != 0 {
		newPC, _ := strconv.Atoi(instruction.Parameters[1])
		pcb.PC = newPC
	}
}

func getRegister(register string, pcb *model.PCB) int {
	switch register {
	case "AX":
		return pcb.Registers.AX
	case "BX":
		return pcb.Registers.BX
	case "CX":
		return pcb.Registers.CX
	case "DX":
		return pcb.Registers.DX
	case "EAX":
		return pcb.Registers.EAX
	case "EBX":
		return pcb.Registers.EBX
	case "ECX":
		return pcb.Registers.ECX
	case "EDX":
		return pcb.Registers.EDX
	default:
		return -1
	}
}

func setRegister(register string, value int, pcb *model.PCB) {
	switch register {
	case "AX":
		pcb.Registers.AX = value
	case "BX":
		pcb.Registers.BX = value
	case "CX":
		pcb.Registers.CX = value
	case "DX":
		pcb.Registers.DX = value
	}
}
