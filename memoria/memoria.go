package main

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/memoria/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
)

const MEMORYLOG = "./memoria.log"

type ProcessPath struct {
	Path string `json:"path"`
}

type ProcessPID struct {
	PID int `json:"pid"`
}

func main() {
	global.InitGlobal()

	processPath := ProcessPath{
		Path: "sisop/tp-go/path",
	}

	processPID, err := requests.PutHTTPwithBody[ProcessPath, ProcessPID](global.MemoryConfig.IPKernel, global.MemoryConfig.PortKernel, "process", processPath)
	if err != nil {
		global.Logger.Log("Error con el put: "+err.Error(), log.ERROR)
	}
	global.Logger.Log(fmt.Sprintf("Struct: %+v", processPID), log.INFO)

	requests.PutHTTPwithBody[interface{}, interface{}](global.MemoryConfig.IPKernel, global.MemoryConfig.PortKernel, "plani", nil)

	global.Logger.CloseLogger()
}
