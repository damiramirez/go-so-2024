package main

import (
	"fmt"
	"os"

	"github.com/sisoputnfrba/tp-golang/memoria/global"
	config "github.com/sisoputnfrba/tp-golang/utils/config"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

const MEMORYLOG = "./memoria.log"

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Uso: programa <go run `modulo`.go dev|prod>")
		return
	}
	env := args[0]

	logger := log.ConfigureLogger(MEMORYLOG, env)
	memoryConfig := config.LoadConfiguration[global.MemoryConfig]("./config/config.json", logger)

	logger.Log(fmt.Sprintf("Port: %d", memoryConfig.Port), log.INFO)
	logger.Log(fmt.Sprintf("MemorySize: %d", memoryConfig.MemorySize), log.INFO)
	logger.Log(fmt.Sprintf("PageSize: %d", memoryConfig.PageSize), log.INFO)
	logger.Log(fmt.Sprintf("InstructionsPath: %s", memoryConfig.InstructionsPath), log.INFO)
	logger.Log(fmt.Sprintf("DelayResponse: %d", memoryConfig.DelayResponse), log.INFO)

	logger.CloseLogger()
}
