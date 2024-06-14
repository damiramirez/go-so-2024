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

type IO interface {
	GetName() string
	GetInstruction() string
}

type IOGen struct {
	Name        string `json:"nombre"`
	Instruction string `json:"instruccion"`
	Time        int    `json:"tiempo"`
	Pid         int    `json:"pid"`
}

func (io IOGen) GetName() string {
	return io.Name
}

func (io IOGen) GetInstruction() string {
	return io.Instruction
}

type IOStd struct {
	Pid       int    `json:"pid"`
	Instruction string `json:"instruccion"`
	Name      string `json:"name"`
	Length    int    `json:"length"`
	NumFrames []int  `json:"numframe"`
	Offset    int    `json:"offset"`
}

func (io IOStd) GetName() string {
	return io.Name
}

func (io IOStd) GetInstruction() string {
	return io.Instruction
}

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
	// time, _ := strconv.Atoi(pcb.Instruction.Parameters[1])
	global.Logger.Log(fmt.Sprintf("Proceso bloqueado %+v", pcb), log.DEBUG)

	io := factoryIO(pcb)

	if !CheckIfExist(io.GetName()) || !CheckIfIsValid(io.GetName(), io.GetInstruction()) {
		moveToExit(pcb)
		return
	}
	global.IoMap[io.GetName()].Sem <- 0
	_, err := requests.PutHTTPwithBody[IO, interface{}](global.KernelConfig.IPIo, global.IoMap[io.GetName()].Port, io.GetInstruction(), io)
	if err != nil {
		global.Logger.Log("Se desconecto IO:"+err.Error(), log.DEBUG)
		delete(global.IoMap, io.GetName())
		moveToExit(pcb)
		return
	}
	<-global.IoMap[io.GetName()].Sem

	BlockToReady(pcb)

	arrayReady := longterm.ConvertListToArray(global.ReadyState)
	arrayPlus := longterm.ConvertListToArray(global.ReadyPlus)

	global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: BLOCK - Estado Actual: %s", pcb.PID, pcb.State), log.INFO)

	if global.KernelConfig.PlanningAlgorithm == "VRR" {
		global.Logger.Log(fmt.Sprintf("Cola Ready : %v, Cola Ready+ : %v", arrayReady, arrayPlus), log.INFO)
	} else {
		global.Logger.Log(fmt.Sprintf("Cola Ready : %v", arrayReady), log.INFO)
	}

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

	if pcb.DisplaceReason == "QUANTUM" {
		pcb.RemainingQuantum = global.KernelConfig.Quantum
	}
}

func factoryIO(pcb *model.PCB) IO {
	switch pcb.Instruction.Operation {
	case "IO_GEN_SLEEP":
		time, _ := strconv.Atoi(pcb.Instruction.Parameters[1])
		return IOGen{
			Name:        pcb.Instruction.Parameters[0],
			Instruction: pcb.Instruction.Operation,
			Time:        time,
			Pid:         pcb.PID,
		}

	case "IO_STDIN_READ", "IO_STDOUT_WRITE":
		return IOStd{
			Name: pcb.Instruction.Parameters[0],
			Pid: pcb.PID,
			Length: pcb.Instruction.Size,
			NumFrames: pcb.Instruction.NumFrames,
			Offset:pcb.Instruction.Offset,
		}
	}

	return nil
}