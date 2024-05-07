package internal

import (
	"fmt"
	"strings"
	"time"
	"github.com/sisoputnfrba/tp-golang/cpu/global"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/execute"
	internal "github.com/sisoputnfrba/tp-golang/cpu/internal/fetch"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

func Dispatch(pcb *model.PCB) (*model.PCB, error) {
	global.Logger.Log(fmt.Sprintf("Recibi PCB %+v", pcb), log.DEBUG)

	global.ExecuteMutex.Lock()
	global.Execute=true
	global.ExecuteMutex.Unlock()
	//esta logica se deberia hacer en una funcion 
	inicioTime := time.Now()
	for global.Execute {
		
		instruction, err := internal.Fetch(pcb)
		if err != nil {
			return nil, err
		}

		exec_result := execute.Execute(pcb, instruction)
		if exec_result == execute.RETURN_CONTEXT{
			global.Execute = false
		}
		QuantumTime:=time.Duration(pcb.Quantum)
		InitTime:=time.Since(inicioTime)
		pcb.CPUTime=int(InitTime*time.Millisecond)
		global.Logger.Log(fmt.Sprintf("valor de cpu time %d",pcb.CPUTime),log.DEBUG)
		//al igual q esta 
		if InitTime*time.Millisecond>QuantumTime*time.Millisecond{
			global.ExecuteMutex.Lock()
			global.Execute=false
			global.ExecuteMutex.Unlock()
		}
		DisplaceReason(pcb)
	}
	global.Logger.Log(fmt.Sprintf("PCB Actualizada %+v", pcb), log.DEBUG)
	return pcb, nil
}

func DisplaceReason(pcb *model.PCB){
	if strings.Contains(pcb.Instruction.Operation, "IO"){
		pcb.DisplaceReason="BLOCKED"
	}else if pcb.Instruction.Operation == "EXIT"{
		pcb.DisplaceReason="EXIT"
	}else{
		pcb.DisplaceReason="QUANTUM"
	}
		
}