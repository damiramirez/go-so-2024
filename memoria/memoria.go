package main

import (
	"fmt"
	"os"

	"github.com/sisoputnfrba/tp-golang/memoria/global"
	config "github.com/sisoputnfrba/tp-golang/utils/config"
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
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Uso: programa <go run `modulo`.go dev|prod>")
		return
	}
	env := args[0]

	logger := log.ConfigureLogger(MEMORYLOG, env)
	global.MemoryConfig = config.LoadConfiguration[global.Config]("./config/config.json", logger)

	processPath := ProcessPath{
		Path: "sisop/tp-go/path",
	}

	processPID, err := requests.PutHTTPwithBody[ProcessPath, ProcessPID](global.MemoryConfig.IPKernel, global.MemoryConfig.PortKernel, "process", processPath, &logger)

	if err != nil {
		logger.Log("Error con el put: "+err.Error(), log.ERROR)
	}
	logger.Log(fmt.Sprintf("Struct: %+v", processPID), log.INFO)

	requests.PutHTTPwithBody[interface{}, interface{}](global.MemoryConfig.IPKernel, global.MemoryConfig.PortKernel, "plani", nil, &logger)


	logger.CloseLogger()
}
