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
	cpuConfig := config.LoadConfiguration[global.CpuConfig]("./config/config.json", logger)

	processSlice, _ := requests.GetHTTP[[]ProcessState](cpuConfig.IPKernel, cpuConfig.PortKernel, "process", &logger)

	logger.Log(fmt.Sprintf("%+v", processSlice), log.INFO)

	logger.CloseLogger()
}
