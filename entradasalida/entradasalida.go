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
	endpoint := "plani/"
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Uso: programa <go run `modulo`.go dev|prod>")
		return
	}
	env := args[0]
	logger := log.ConfigureLogger(IOLOG, env)
	global.IOConfig = config.LoadConfiguration[global.Config]("./config/config.json", logger)

	processSlice, _ := requests.GetHTTP[ProcessState](global.IOConfig.IPKernel, global.IOConfig.PortKernel, "process/12", &logger)

	requests.DeleteHTTP[interface{}](endpoint, logger, global.IOConfig.PortKernel, nil, global.IOConfig.IPKernel)

	logger.Log(fmt.Sprintf("%+v", processSlice), log.INFO)

	logger.CloseLogger()
}
