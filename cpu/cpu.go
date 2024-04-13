package main

import (
	"fmt"
	"os"

	"github.com/sisoputnfrba/tp-golang/cpu/global"
	config "github.com/sisoputnfrba/tp-golang/utils/config"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	requests "github.com/sisoputnfrba/tp-golang/utils/requests"
)

const CPULOG = "./cpu.log"

type ProcessState struct {
	PID   int    `json:"pid"`
	State string `json:"state"`
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Uso: programa <go run `modulo`.go dev|prod>")
		return
	}
	env := args[0]

	logger := log.ConfigureLogger(CPULOG, env)
	global.CPUConfig = config.LoadConfiguration[global.Config]("./config/config.json", logger)

	processSlice, _ := requests.GetHTTP[[]ProcessState](global.CPUConfig.IPKernel, global.CPUConfig.Port, "process", &logger)

	logger.Log(fmt.Sprintf("%+v", processSlice), log.INFO)

	logger.CloseLogger()
}
