package global

import (
	"fmt"
	"os"

	config "github.com/sisoputnfrba/tp-golang/utils/config"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

const KERNELLOG = "./kernel.log"

type Config struct {
	Port              int      `json:"port"`
	IPMemory          string   `json:"ip_memory"`
	PortMemory        int      `json:"port_memory"`
	IPCPU             string   `json:"ip_cpu"`
	PortCPU           int      `json:"port_cpu"`
	PlanningAlgorithm string   `json:"planning_algorithm"`
	Quantum           int      `json:"quantum"`
	Resources         []string `json:"resources"`
	ResourceInstances []int    `json:"resource_instances"`
	Multiprogramming  int      `json:"multiprogramming"`
}

var KernelConfig *Config

var Logger *log.LoggerStruct

var nextPID int = 1

func InitGlobal() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Uso: programa <go run `modulo`.go dev|prod>")
		os.Exit(1)
	}
	env := args[0]

	Logger = log.ConfigureLogger(KERNELLOG, env)
	KernelConfig = config.LoadConfiguration[Config]("./config/config.json")
}

func GetNextPID() int {
	actualPID := nextPID
	nextPID++
	return actualPID
}
