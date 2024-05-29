package resource

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

func Wait(Pcb *model.PCB){
	
	Resource:= global.ResourceMap[Pcb.Instruction.Parameters[0]] 
	Resource.Count-=1
	Resource.PidList=append(Resource.PidList,Pcb.PID)
	global.Logger.Log(fmt.Sprintf("%s %d",Resource.Name,Resource.Count),log.DEBUG)
	if Resource.Count<0 {
		Resource.MutexList.Lock()
		Resource.BlockedList.PushBack(Pcb)
		Resource.MutexList.Unlock()
		//poner en listar procesos
		global.Logger.Log(fmt.Sprintf("se bloqueo el proceso con PID %d",Pcb.PID),log.DEBUG)
	}else {
		
		global.MutexReadyState.Lock()
		global.ReadyState.PushFront(Pcb)
		global.MutexReadyState.Unlock()
		global.Logger.Log(fmt.Sprintf("se mando el proceso con PID %d primero a ready",Pcb.PID),log.DEBUG)
	}

}

func Signal(PcbExec *model.PCB){
	
	Resource:= global.ResourceMap[PcbExec.Instruction.Parameters[0]] 
	Resource.Count+=1
	global.Logger.Log(fmt.Sprintf("%s %d",Resource.Name,Resource.Count),log.DEBUG)
	if Resource.BlockedList.Len()>0 {
		Resource.MutexList.Lock()
		PCBBlock := Resource.BlockedList.Front().Value.(*model.PCB)
		Resource.BlockedList.Remove(Resource.BlockedList.Front())
		Resource.MutexList.Unlock()
		global.MutexReadyState.Lock()
		global.ReadyState.PushBack(PCBBlock)
		global.MutexReadyState.Unlock()
		global.Logger.Log(fmt.Sprintf("lo mandamos al fondo de ready %d",PCBBlock.PID),log.DEBUG)
		global.SemReadyList <- struct{}{}
		
	}
	Value:=CheckInArray(Resource.PidList,PcbExec.PID)
	if Value!=-1 {
		Resource.PidList=removeAt(Resource.PidList,Value)
	}

	global.Logger.Log(fmt.Sprintf("Pid list %+v",Resource.PidList),log.DEBUG)
	global.MutexReadyState.Lock()
	global.ReadyState.PushFront(PcbExec)
	global.MutexReadyState.Unlock()
	global.SemReadyList <- struct{}{}
	global.Logger.Log(fmt.Sprintf("se mando el proceso con PID %d primero a ready",PcbExec.PID),log.DEBUG)
}

func CheckInArray(Array []int,Pid int)int{
	for i,Value:=range Array{
		if Value==Pid {
			return i
		}
	}
	return -1
}

func removeAt(slice []int, index int) []int {
    return append(slice[:index], slice[index+1:]...)
}