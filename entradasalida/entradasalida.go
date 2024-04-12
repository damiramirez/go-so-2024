package main

import (
	"fmt"
	"os"

	"github.com/sisoputnfrba/tp-golang/entradasalida/global"
	config "github.com/sisoputnfrba/tp-golang/utils/config"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
)

const IOLOG = "./entradaysalida.log"

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

	logger := log.ConfigureLogger(IOLOG, env)
	ioConfig := config.LoadConfiguration[global.IOConfig]("./config/config.json", logger)

	processSlice, _ := requests.GetHTTP[ProcessState](ioConfig.IPKernel, ioConfig.PortKernel, "process/12", &logger)

	logger.Log(fmt.Sprintf("%+v", processSlice), log.INFO)

	logger.CloseLogger()
}
