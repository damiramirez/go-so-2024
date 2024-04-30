package global

import (
	"container/list"
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


//States
var ReadyState *list.List
var NewState *list.List
var BlockedState *list.List
var RunningState *list.List
var Exit *list.List



func InitGlobal() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Uso: programa <go run `modulo`.go dev|prod>")
		os.Exit(1)
	}
	env := args[0]

	Logger = log.ConfigureLogger(KERNELLOG, env)
	KernelConfig = config.LoadConfiguration[Config]("./config/config.json")
	NewState = list.New()
	ReadyState= list.New()
	BlockedState= list.New()
	RunningState= list.New()
	Exit= list.New()



}

func GetNextPID() int {
	actualPID := nextPID
	nextPID++
	return actualPID
}
