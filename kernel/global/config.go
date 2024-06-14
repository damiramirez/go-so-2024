package global

import (
	"container/list"
	"fmt"
	"os"
	"sync"

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
	IPIo              string   `json:"ip_io"`
	PlanningAlgorithm string   `json:"planning_algorithm"`
	Quantum           int      `json:"quantum"`
	Resources         []string `json:"resources"`
	ResourceInstances []int    `json:"resource_instances"`
	Multiprogramming  int      `json:"multiprogramming"`
}
type IoDevice struct {
	Port int    `json:"port"`
	Name string `json:"name"`
	Type string `json:"type"`
	Sem  chan int
}
type Resource struct {
	Name        string
	Count       int
	BlockedList *list.List
	MutexList   sync.Mutex
	PidList     []int
}

var KernelConfig *Config
var Logger *log.LoggerStruct
var nextPID int = 1

// States
var ReadyState *list.List
var NewState *list.List
var BlockedState *list.List
var ExecuteState *list.List
var ExitState *list.List
var ReadyPlus *list.List

var WorkingPlani bool

// Mutex
var MutexReadyState sync.Mutex
var MutexNewState sync.Mutex
var MutexExitState sync.Mutex
var MutexBlockState sync.Mutex
var MutexExecuteState sync.Mutex
var MutexPlani sync.Mutex
var MutexReadyPlus sync.Mutex

// Semaforos
var SemMulti chan int
var SemExecute chan int
var SemInterrupt chan int
var SemReadyList chan struct{}
var SemNewList chan struct{}
var SemStopPlani chan struct{}
var SemReadyPlus chan struct{}

// MAPs
var IoMap map[string]IoDevice
var ResourceMap map[string]*Resource
var PIDResourceMap map[int][]string

func InitGlobal() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("ARGS: ENV=dev|prod CONFIG=config_path")
		os.Exit(1)
	}
	env := args[0]
	configFile := args[1]

	Logger = log.ConfigureLogger(KERNELLOG, env)
	KernelConfig = config.LoadConfiguration[Config](configFile)

	NewState = list.New()
	ReadyState = list.New()
	BlockedState = list.New()
	ExecuteState = list.New()
	ExitState = list.New()
	ReadyPlus = list.New()

	SemStopPlani = make(chan struct{})
	SemMulti = make(chan int, KernelConfig.Multiprogramming)
	SemExecute = make(chan int, 1)
	SemInterrupt = make(chan int)
	SemReadyList = make(chan struct{}, KernelConfig.Multiprogramming)

	// Revisar el size
	SemNewList = make(chan struct{}, 20)
	ResourceMap = CreateResourceMap()
	IoMap = map[string]IoDevice{}
	PIDResourceMap = map[int][]string{}

	WorkingPlani = false
}

func GetNextPID() int {
	actualPID := nextPID
	nextPID++
	return actualPID
}

func CreateResourceMap() map[string]*Resource {
	ResourceMap = map[string]*Resource{}
	for i := 0; i < len(KernelConfig.Resources); i++ {
		ResourceMap[KernelConfig.Resources[i]] = &Resource{
			Name:        KernelConfig.Resources[i],
			Count:       KernelConfig.ResourceInstances[i],
			PidList:     make([]int, 0),
			BlockedList: list.New()}
	}

	return ResourceMap
}
