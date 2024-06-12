package block

import (
	"fmt"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/longterm"

	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
)

// var acceptedInstructions map[string] []string
var acceptedInstructions = map[string][]string{
	"GEN":    {"IO_GEN_SLEEP"},
	"STDIN":  {"IO_STDIN_READ"},
	"STDOUT": {"IO_STDOUT_WRITE"},
	"DIALFS": {"IO_FS_CREATE", "IO_FS_DELETE", "IO_FS_TRUNCATE", "IO_FS_WRITE", "IO_FS_READ"},
}

// MAP de IOs "conectadas"
// nombre - puerto - tipo
// implemnetar checktype y chequear para distintas ios
func CheckIfExist(name string) bool {
	_, Ioexist := global.IoMap[name]
	return Ioexist
}
func CheckIfIsValid(name, instruccion string) bool {
	validInstructions := acceptedInstructions[global.IoMap[name].Type]
	for _, ins := range validInstructions {
		if instruccion == ins {
			return true
		}
	}
	return false
}

func ProcessToIO(pcb *model.PCB) {
	type IOStruct struct {
		Name        string `json:"nombre"`
		Instruccion string `json:"instruccion"`
		Time        int    `json:"tiempo"`
		Pid         int    `json:"pid"`
	}

	time, _ := strconv.Atoi(pcb.Instruction.Parameters[1])
	global.Logger.Log(fmt.Sprintf("Proceso bloqueado %+v", pcb), log.DEBUG)

	ioStruct := IOStruct{
		Name:        pcb.Instruction.Parameters[0],
		Instruccion: pcb.Instruction.Operation,
		Time:        time,
		Pid:         pcb.PID,
	}
	if !CheckIfExist(ioStruct.Name) || !CheckIfIsValid(ioStruct.Name, ioStruct.Instruccion) {
		moveToExit(pcb)
		return 
	}
	global.IoMap[ioStruct.Name].Sem <- 0
	_, err := requests.PutHTTPwithBody[IOStruct, interface{}](global.KernelConfig.IPIo, global.IoMap[ioStruct.Name].Port, ioStruct.Instruccion, ioStruct)
	if err != nil {
		global.Logger.Log("Se desconecto IO:"+err.Error(), log.DEBUG)
		delete(global.IoMap, ioStruct.Name)
		moveToExit(pcb)
		return
	}
	<-global.IoMap[ioStruct.Name].Sem

	BlockToReady(pcb)

	arrayReady := longterm.ConvertListToArray(global.ReadyState)
	arrayPlus := longterm.ConvertListToArray(global.ReadyPlus)

	global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: BLOCK - Estado Actual: %s", pcb.PID, pcb.State), log.INFO)
	global.Logger.Log(fmt.Sprintf("Cola Ready : %v, Cola Ready+ : %v", arrayReady, arrayPlus), log.INFO)
	global.SemReadyList <- struct{}{}

}

func moveToExit(pcb *model.PCB) {
	global.MutexBlockState.Lock()
	global.BlockedState.Remove(global.BlockedState.Front())
	global.MutexBlockState.Unlock()

	pcb.State = "EXIT"

	global.MutexExitState.Lock()
	global.ExitState.PushBack(pcb)
	global.MutexExitState.Unlock()

	global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: BLOCK - Estado Actual: %s ", pcb.PID, pcb.State), log.INFO)
	global.Logger.Log(fmt.Sprintf("Finaliza el proceso %d - Motivo: INVALID_RESOURCE", pcb.PID), log.INFO)
}

func BlockToReady(pcb *model.PCB) {
	// Saco de block cuando termino la IO
	global.MutexBlockState.Lock()
	global.BlockedState.Remove(global.BlockedState.Front())
	global.MutexBlockState.Unlock()

	pcb.State = "READY"

	if global.KernelConfig.PlanningAlgorithm == "VRR" && pcb.RemainingQuantum > 0 {
		global.MutexReadyPlus.Lock()
		global.ReadyPlus.PushBack(pcb)
		global.MutexReadyPlus.Unlock()

		global.Logger.Log(fmt.Sprintf("PID: %d - Bloqueado a Ready Plus", pcb.PID), log.DEBUG)

	} else {
		global.MutexReadyState.Lock()
		global.ReadyState.PushBack(pcb)
		global.MutexReadyState.Unlock()
		global.Logger.Log(fmt.Sprintf("PID: %d - Bloqueado normal", pcb.PID), log.DEBUG)

	}

}
