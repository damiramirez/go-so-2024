package main

import (
	"fmt"
	"os"

	"github.com/sisoputnfrba/tp-golang/cpu/api"
	"github.com/sisoputnfrba/tp-golang/cpu/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

type ProcessState struct {
	PID   int    `json:"pid"`
	State string `json:"state"`
}

func main() {
	global.InitGlobal()

	s := api.CreateServer()

	// Levanto server
	global.Logger.Log(fmt.Sprintf("Starting cpu server on port: %d", global.CPUConfig.Port), log.INFO)
	if err := s.Start(); err != nil {
		global.Logger.Log(fmt.Sprintf("Failed to start cpu server: %v", err), log.ERROR)
		os.Exit(1)
	}

	global.Logger.CloseLogger()
}
