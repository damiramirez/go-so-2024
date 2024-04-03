package main

import (
	"fmt"
	"os"

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

	logger.Log(fmt.Sprintf("Port: %d", kernelConfig.Port), log.INFO)
	logger.Log(fmt.Sprintf("IPMemory: %s", kernelConfig.IPMemory), log.INFO)
	logger.Log(fmt.Sprintf("PortMemory: %d", kernelConfig.PortMemory), log.INFO)
	logger.Log(fmt.Sprintf("IPCPU: %s", kernelConfig.IPCPU), log.INFO)
	logger.Log(fmt.Sprintf("PortCPU: %d", kernelConfig.PortCPU), log.INFO)
	logger.Log(fmt.Sprintf("PlanningAlgorithm: %s", kernelConfig.PlanningAlgorithm), log.INFO)
	logger.Log(fmt.Sprintf("Quantum: %d", kernelConfig.Quantum), log.INFO)

	resourcesStr := fmt.Sprintf("Resources: %v", kernelConfig.Resources)
	logger.Log(resourcesStr, log.INFO)

	resourceInstancesStr := fmt.Sprintf("ResourceInstances: %v", kernelConfig.ResourceInstances)
	logger.Log(resourceInstancesStr, log.INFO)

	logger.Log(fmt.Sprintf("Multiprogramming: %d", kernelConfig.Multiprogramming), log.INFO)

	logger.CloseLogger()
}
