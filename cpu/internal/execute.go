package internal

import (
	"strconv"

	"github.com/sisoputnfrba/tp-golang/utils/model"
)

// JNZ IO_GEN_SLEEP

// func Execute(instruction ) error {

// }

func set(pcb *model.PCB, instruction model.Instruction) {
	value, _ := strconv.Atoi(instruction.Parameters[1])
	setRegister(instruction.Parameters[0], value, pcb)
}

func sum(pcb *model.PCB, instruction model.Instruction) {

	destinationValue := getRegister(instruction.Parameters[0], pcb)
	sourceValue := getRegister(instruction.Parameters[1], pcb)
	destinationValue = destinationValue + sourceValue
	setRegister(instruction.Parameters[0], destinationValue, pcb)
}

func sub(pcb *model.PCB, instruction model.Instruction) {

	destinationValue := getRegister(instruction.Parameters[0], pcb)
	sourceValue := getRegister(instruction.Parameters[1], pcb)
	destinationValue = destinationValue - sourceValue
	setRegister(instruction.Parameters[0], destinationValue, pcb)
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
