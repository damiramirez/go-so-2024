package global

import (
	"fmt"
	"os"

	config "github.com/sisoputnfrba/tp-golang/utils/config"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

const CPULOG = "./cpu.log"

type Config struct {
	Port             int    `json:"port"`
	IPKernel         string `json:"ip_kernel"`
	PortKernel       int    `json:"port_kernel"`
	IPMemory         string `json:"ip_memory"`
	PortMemory       int    `json:"port_memory"`
	NumberFellingTLB int    `json:"number_felling_tlb"`
	AlgorithmTLB     string `json:"algorithm_tlb"`
	Page_size        int    `json:"page_size"`
}

var CPUConfig *Config

var Logger *log.LoggerStruct

func InitGlobal() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Uso: programa <go run `modulo`.go dev|prod>")
		os.Exit(1)
	}
	env := args[0]

	Logger = log.ConfigureLogger(CPULOG, env)
	CPUConfig = config.LoadConfiguration[Config]("./config/config.json")
}
