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

func ProcessToIO(blockProcess *model.PCB) {
	type IOStruct struct {
		Name        string `json:"nombre"`
		Instruccion string `json:"instruccion"`
		Time        int    `json:"tiempo"`
		Pid         int    `json:"pid"`
	}

	// TODO: MUTEX

	//blockProcess := global.BlockedState.Front().Value.(*model.PCB)
	/*global.MutexBlockState.Lock()
	blockProcess:=global.BlockedState.Remove(global.BlockedState.Front()).(*model.PCB)
	global.MutexBlockState.Unlock()*/

	time, _ := strconv.Atoi(blockProcess.Instruction.Parameters[1])
	global.Logger.Log(fmt.Sprintf("Proceso bloqueado %+v", blockProcess), log.INFO)

	ioStruct := IOStruct{
		Name:        blockProcess.Instruction.Parameters[0],
		Instruccion: blockProcess.Instruction.Operation,
		Time:        time,
		Pid:         blockProcess.PID,
	}
	if !CheckIfExist(ioStruct.Name) || !CheckIfIsValid(ioStruct.Name, ioStruct.Instruccion) {
		global.MutexBlockState.Lock()
		global.BlockedState.Remove(global.BlockedState.Front())
		global.MutexBlockState.Unlock()

		blockProcess.State = "EXIT"

		global.MutexExitState.Lock()
		global.ExitState.PushBack(blockProcess)
		global.MutexExitState.Unlock()
		global.SemReadyList <- struct{}{}
		global.Logger.Log(fmt.Sprintf("Finaliza el proceso %d - Motivo: INVALID_RESOURCE", blockProcess.PID), log.INFO)
		global.Logger.Log(fmt.Sprintf("PID: <%d> - Estado Anterior: BLOCK - Estado Actual: %s ", blockProcess.PID, blockProcess.State), log.INFO)

		return
	}
	global.IoMap[ioStruct.Name].Sem <- 0
	_, err := requests.PutHTTPwithBody[IOStruct, interface{}](global.KernelConfig.IPIo, global.IoMap[ioStruct.Name].Port, ioStruct.Instruccion, ioStruct)
	if err != nil {
		global.Logger.Log("ERROR AL REQUEST IO:"+err.Error(), log.DEBUG)
	}

	// Saco de block cuando termino la IO
	global.MutexBlockState.Lock()
	global.BlockedState.Remove(global.BlockedState.Front())
	global.MutexBlockState.Unlock()
	<-global.IoMap[ioStruct.Name].Sem
	// TESTEO -  HACER EN OTRA FUNCION
	blockProcess.State = "READY"

	global.MutexReadyState.Lock()
	global.ReadyState.PushBack(blockProcess)

	global.MutexReadyState.Unlock()
	array := longterm.ConvertListToArray(global.ReadyState)
	global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: BLOCK - Estado Actual: %s", blockProcess.PID, blockProcess.State), log.INFO)
	global.Logger.Log(fmt.Sprintf("Cola Ready : %v", array), log.INFO)

	global.SemReadyList <- struct{}{}

}
