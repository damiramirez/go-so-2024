package main

import (
	"fmt"
	"os"

	api "github.com/sisoputnfrba/tp-golang/memoria/api"
	global "github.com/sisoputnfrba/tp-golang/memoria/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
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

	s := api.CreateServer()

	global.Logger.Log(fmt.Sprintf("Starting kernel server on port: %d", global.MemoryConfig.Port), log.INFO)
	if err := s.Start(); err != nil {
		global.Logger.Log(fmt.Sprintf("Failed to start kernel server: %v", err), log.ERROR)
		os.Exit(1)
	}

	global.Logger.CloseLogger()
}
