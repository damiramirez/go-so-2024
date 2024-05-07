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

func ProcessToIO() (*model.PCB, error) {
	type IOStruct struct {
		Name        string `json:"name"`
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

	_, err := requests.PutHTTPwithBody[IOStruct, interface{}](global.KernelConfig.IPIo, 8005, "Sleep", ioStruct)
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
	global.MutexReadyState.Unlock()

	global.SemReadyList <- struct{}{}

	return blockProcess, nil
}
