package main

import (
	"fmt"
	"os"

	"github.com/sisoputnfrba/tp-golang/kernel/api"
	"github.com/sisoputnfrba/tp-golang/kernel/global"
	config "github.com/sisoputnfrba/tp-golang/utils/config"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

const KERNELLOG = "./kernel.log"

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Uso: programa <go run `modulo`.go dev|prod>")
		return
	}
	env := args[0]

	logger := log.ConfigureLogger(KERNELLOG, env)
	kernelConfig := config.LoadConfiguration[global.Config]("./config/config.json", logger)

	server := api.NewServer(logger)
	server.Start(kernelConfig.Port)

	logger.CloseLogger()
}
