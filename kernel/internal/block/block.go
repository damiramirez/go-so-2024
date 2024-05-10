package block

import (
	"fmt"
	"strconv"

	
	"github.com/sisoputnfrba/tp-golang/kernel/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
)

// MAP de IOs "conectadas"
// nombre - puerto - tipo
//implemnetar checktype y chequear para distintas ios
func CheckIfExist(name string) bool{
	_,Ioexist:=global.IoMap[name]
	return Ioexist
}
func CheckType(name string,Type string)bool{
	if global.IoMap[name].Type==Type{

		return global.IoMap[name].Type==Type
	}
	return false 
}
func ProcessToIO()  {
	type IOStruct struct {
		Name        string `json:"nombre"`
		Instruccion string `json:"instruccion"`
		Time        int    `json:"tiempo"`
	}

	// TODO: MUTEX
	
	blockProcess := global.BlockedState.Front().Value.(*model.PCB)
	
	time, _ := strconv.Atoi(blockProcess.Instruction.Parameters[1])
	global.Logger.Log(fmt.Sprintf("Proceso bloqueado %+v", blockProcess), log.DEBUG)

	ioStruct := IOStruct{
		Name:        blockProcess.Instruction.Parameters[0],
		Instruccion: blockProcess.Instruction.Operation,
		Time:        time,
	}
	if !CheckIfExist(ioStruct.Name){
		global.MutexBlockState.Lock()
		global.BlockedState.Remove(global.BlockedState.Front())
		global.MutexBlockState.Unlock()

		blockProcess.State = "EXIT"

		global.MutexExitState.Lock()
		global.ExitState.PushBack(blockProcess)
		global.MutexExitState.Unlock()
		global.SemReadyList <- struct{}{}
		global.Logger.Log(fmt.Sprintf("Finaliza el proceso %d - Motivo: INVALID_RESOURCE",blockProcess.PID), log.INFO)
		global.Logger.Log(fmt.Sprintf("PID: <%d> - Estado Anterior: BLOCK - Estado Actual: %s ",blockProcess.PID,blockProcess.State), log.INFO)

		return
	}
	
	_, err := requests.PutHTTPwithBody[IOStruct, interface{}](global.KernelConfig.IPIo, global.IoMap[ioStruct.Name].Port, ioStruct.Instruccion, ioStruct)
	if err != nil {
		global.Logger.Log("ERROR AL REQUEST IO:"+err.Error(), log.DEBUG)
	}

	// Saco de block cuando termino la IO
	global.MutexBlockState.Lock()
	global.BlockedState.Remove(global.BlockedState.Front())
	global.MutexBlockState.Unlock()
	
	// TESTEO -  HACER EN OTRA FUNCION
	blockProcess.State = "READY"
	global.MutexReadyState.Lock()
	global.ReadyState.PushBack(blockProcess)
	global.PidList.PushBack(blockProcess.PID)
	global.MutexReadyState.Unlock()

	global.SemReadyList <- struct{}{}

	
}
