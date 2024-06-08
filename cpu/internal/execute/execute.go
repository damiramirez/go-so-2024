package execute

import (
	"fmt"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
)

// TODO: IO_GEN_SLEEP

const (
	CONTINUE       = 0
	RETURN_CONTEXT = 1
)

type Estructura_mov struct {
	DataValue      int `json:"data"`
	DirectionValue int `json:"direction"`
}

var estructura_mov Estructura_mov

type Estructura_resize struct {
	Pid  int `json:"pid"`
	Size int `json:"size"`
}

type Response struct {
	Respuesta string `json:"respuesta"`
}

var estructura_resize Estructura_resize
var result = 0

// Ejecuto -> sumo PC en dispatch?
func Execute(pcb *model.PCB, instruction *model.Instruction) int {

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
	case "WAIT":
		result = RETURN_CONTEXT
	case "SIGNAL":
		result = RETURN_CONTEXT
	case "MOV_IN":
		mov_in(pcb, instruction)
		result = CONTINUE
	case "MOV_OUT":
		mov_out(pcb, instruction)
		result = CONTINUE
	case "RESIZE":
		resize(pcb, instruction)
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

func mov_in(pcb *model.PCB, instruction *model.Instruction) {
	dataValue := instruction.Parameters[0]
	estructura_mov.DirectionValue = getRegister(instruction.Parameters[1], pcb)

	// put a memoria para que devuelva el valor solicitado

	resp, err := requests.PutHTTPwithBody[Estructura_mov, Estructura_mov](global.CPUConfig.IPMemory, global.CPUConfig.PortMemory, "mov_in", estructura_mov)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("NO se pudo enviar a memoria la estructura %s", err.Error()), log.INFO)
		panic(1)
		// TODO: falta que memoria vea si puede escribir o no (?)
	}
	global.Logger.Log(fmt.Sprintf("Resp %+v", resp), log.INFO)
	setRegister(dataValue, int(resp.DataValue), pcb)
}

func mov_out(pcb *model.PCB, instruction *model.Instruction) {
	dataValue := getRegister(instruction.Parameters[1], pcb)
	directionValue := getRegister(instruction.Parameters[0], pcb)

	estructura_mov.DataValue = dataValue
	estructura_mov.DataValue = directionValue // esta es la dirección que hay que traducir de Lógica a Física

	// put a memoria para que guarde

	_, err := requests.PutHTTPwithBody[Estructura_mov, interface{}](global.CPUConfig.IPMemory, global.CPUConfig.PortMemory, "mov_out", estructura_mov)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("NO se pudo enviar a memoria la estructura %s", err.Error()), log.INFO)
		panic(1)
		// TODO: falta que memoria vea si puede escribir o no (?)
	}

}

func resize(pcb *model.PCB, instruction *model.Instruction) {
	newSize, _ := strconv.Atoi(instruction.Parameters[0])
	estructura_resize.Pid = pcb.PID
	estructura_resize.Size = newSize
	// put a memoria para hacer el resize

	resp, err := requests.PutHTTPwithBody[Estructura_resize, Response](global.CPUConfig.IPMemory, global.CPUConfig.PortMemory, "resize", estructura_resize)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("NO se pudo enviar a memoria la estructura %s", err.Error()), log.INFO)
		panic(1)
		// TODO: falta que memoria vea si puede escribir o no (?)
	}

	if resp.Respuesta == "Out of Memory" {
		result = RETURN_CONTEXT
		return
	}
	result = CONTINUE

}

func copyString(pcb *model.PCB, instruction *model.Instruction) {

	tamanio, _ := strconv.Atoi(instruction.Parameters[0])

	// put a memoria para obtener tamanio bytes de lo que hay en el string apuntado por SI

	// put a memoria para guardar en DI lo que obtuve en el primer put

}
