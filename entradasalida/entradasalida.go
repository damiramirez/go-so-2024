package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/entradasalida/global"
	config "github.com/sisoputnfrba/tp-golang/utils/config"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

const IOLOG = "./entradaysalida.log"

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Uso: programa <go run `modulo`.go dev|prod>")
		return
	}
	env := args[0]

	logger := log.ConfigureLogger(IOLOG, env)
	ioConfig := config.LoadConfiguration[global.Config]("./config/config.json", logger)

	logger.Log("Port: "+strconv.Itoa(ioConfig.Port), log.INFO)
	logger.Log("Type: "+ioConfig.Type, log.INFO)
	logger.Log("UnitWorkTime: "+strconv.Itoa(ioConfig.UnitWorkTime), log.INFO)
	logger.Log("IPKernel: "+ioConfig.IPKernel, log.INFO)
	logger.Log("PortKernel: "+strconv.Itoa(ioConfig.PortKernel), log.INFO)
	logger.Log("IPMemory: "+ioConfig.IPMemory, log.INFO)
	logger.Log("PortMemory: "+strconv.Itoa(ioConfig.PortMemory), log.INFO)
	logger.Log("DialFSPath: "+ioConfig.DialFSPath, log.INFO)
	logger.Log("DialFSBlockSize: "+strconv.Itoa(ioConfig.DialFSBlockSize), log.INFO)
	logger.Log("DialFSBlockCount: "+strconv.Itoa(ioConfig.DialFSBlockCount), log.INFO)

	logger.CloseLogger()
}
