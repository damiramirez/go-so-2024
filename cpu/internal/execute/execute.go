package execute

import (
	"fmt"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/global"
	"github.com/sisoputnfrba/tp-golang/cpu/internal"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
	"github.com/sisoputnfrba/tp-golang/utils/requests"

)

// TODO: IO_GEN_SLEEP

const (
	CONTINUE       = 0
	RETURN_CONTEXT = 1
)


type Estructura_resize struct {
	Pid  int `json:"pid"`
	NumFrames int `json:"size"`
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
	case "EAX":
		pcb.Registers.EAX = value
	case "EBX":
		pcb.Registers.EBX = value
	case "ECX":
		pcb.Registers.ECX = value
	case "EDX":
		pcb.Registers.EDX = value
	case "PC":
		pcb.PC = value
	}
}

func mov_in(pcb *model.PCB, instruction *model.Instruction) {
	dataValue := instruction.Parameters[0]
	LogAdress:=getRegister(instruction.Parameters[1],pcb)

	SendStruct:=internal.CreateAdress(dataValue,LogAdress,pcb.PID,getRegister(dataValue,pcb))
	
	// put a memoria para que devuelva el valor solicitado

	resp, err := requests.PutHTTPwithBody[internal.MemStruct, int](global.CPUConfig.IPMemory, global.CPUConfig.PortMemory, "mov_in", SendStruct)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("NO se pudo enviar a memoria la estructura %s", err.Error()), log.INFO)
		panic(1)
		// TODO: falta que memoria vea si puede escribir o no (?)
	}
	global.Logger.Log(fmt.Sprintf("Resp %+v", resp), log.DEBUG)
	setRegister(dataValue, *resp, pcb)
	global.Logger.Log(fmt.Sprintf("PID: %d - Acción: LEER - Dirección Física: %d %d - Valor: %d",pcb.PID,SendStruct.NumFrames[0],SendStruct.Offset,*resp),log.INFO)
}

func mov_out(pcb *model.PCB, instruction *model.Instruction) {
	dataValue := instruction.Parameters[1]
	LogAdress:=getRegister(instruction.Parameters[0],pcb)
	SendStruct:=internal.CreateAdress(dataValue,LogAdress,pcb.PID,getRegister(dataValue,pcb))

	// put a memoria para que guarde

	_, err := requests.PutHTTPwithBody[internal.MemStruct, interface{}](global.CPUConfig.IPMemory, global.CPUConfig.PortMemory, "mov_out",SendStruct)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("NO se pudo enviar a memoria la estructura %s", err.Error()), log.INFO)
		panic(1)
		// TODO: falta que memoria vea si puede escribir o no (?)
	}
	global.Logger.Log(fmt.Sprintf("PID: %d - Acción: ESCRIBIR - Dirección Física: %d %d - Valor: %d",pcb.PID,SendStruct.NumFrames[0],SendStruct.Offset,SendStruct.Content),log.INFO)
}

func resize(pcb *model.PCB, instruction *model.Instruction) {
	newSize, _ := strconv.Atoi(instruction.Parameters[0])
	estructura_resize.Pid = pcb.PID
	estructura_resize.NumFrames = newSize/global.CPUConfig.Page_size
	// put a memoria para hacer el resize

	_, err := requests.PutHTTPwithBody[Estructura_resize, Response](global.CPUConfig.IPMemory, global.CPUConfig.PortMemory, "resize", estructura_resize)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("OUT OF MEMORY %s", err.Error()), log.INFO)
		result = RETURN_CONTEXT
		return
		// TODO: falta que memoria vea si puede escribir o no (?)
	}

	result = CONTINUE

}

/*func copyString(pcb *model.PCB, instruction *model.Instruction) {

	tamanio, _ := strconv.Atoi(instruction.Parameters[0])

	// put a memoria para obtener tamanio bytes de lo que hay en el string apuntado por SI

	// put a memoria para guardar en DI lo que obtuve en el primer put

}*/
