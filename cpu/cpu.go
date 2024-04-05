package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/global"
	config "github.com/sisoputnfrba/tp-golang/utils/config"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

const CPULOG = "./cpu.log"

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Uso: programa <go run `modulo`.go dev|prod>")
		return
	}
	env := args[0]

	logger := log.ConfigureLogger(CPULOG, env)
	cpuConfig := config.LoadConfiguration[global.CpuConfig]("./config/config.json", logger)

	logger.Log("Port: "+strconv.Itoa(cpuConfig.Port), log.INFO)
	logger.Log("IpMemory: "+cpuConfig.IPMemory, log.INFO)
	logger.Log("PortMemory: "+strconv.Itoa(cpuConfig.PortMemory), log.INFO)
	logger.Log("NumberFellingTlb: "+strconv.Itoa(cpuConfig.NumberFellingTLB), log.INFO)
	logger.Log("AlgorithmTlb: "+cpuConfig.AlgorithmTLB, log.INFO)

	logger.CloseLogger()
}
