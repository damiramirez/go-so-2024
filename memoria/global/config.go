package global

import (
	"fmt"
	"os"

	config "github.com/sisoputnfrba/tp-golang/utils/config"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

const MEMORYLOG = "./memoria.log"

type Config struct {
	Port             int    `json:"port"`
	IPKernel         string `json:"ip_kernel"`
	PortKernel       int    `json:"port_kernel"`
	MemorySize       int    `json:"memory_size"`
	PageSize         int    `json:"page_size"`
	InstructionsPath string `json:"instructions_path"`
	DelayResponse    int    `json:"delay_response"`
}
type ValoraMandar struct {
	Texto string `json:"texto"`
}

var ValoraM ValoraMandar

var MemoryConfig *Config
var Logger *log.LoggerStruct

type ListInstructions struct {
	Instructions []string
	PageTable    *PageTable
}

type Estructura_mov struct {
	DataValue      int `json:"data"`
	DirectionValue int `json:"direction"`
}

type Estructura_resize struct {
	Pid  int `json:"pid"`
	Size int `json:"size"`
}

var DictProcess map[int]ListInstructions

type MemoryST struct {
	Spaces []byte
}
type PageTable struct {
	Pages       []int
}

var Memory *MemoryST

func NewMemory() *MemoryST {

	ByteArray := make([]byte, MemoryConfig.MemorySize)
	mem := MemoryST{Spaces: ByteArray}
	return &mem
}


var PTable *PageTable

func NewPageTable() *PageTable {
	//inicializo las 16 paginas en -1
	Slice:=make([]int,0)
	
	//le asigno al "struct" pagetable el array con las paginas
	pagetable := PageTable{Pages: Slice}

	return &pagetable
}

func NewBitMap()[]int{
	NumPages:=MemoryConfig.MemorySize/MemoryConfig.PageSize
	Array := make([]int, NumPages)
	
	return Array
}

var BitMap []int
func InitGlobal() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("ARGS: ENV=dev|prod CONFIG=config_path")
		os.Exit(1)
	}
	env := args[0]
	configFile := args[1]

	Logger = log.ConfigureLogger(MEMORYLOG, env)
	MemoryConfig = config.LoadConfiguration[Config](configFile)
	DictProcess = map[int]ListInstructions{}
	Memory = NewMemory()
	BitMap = NewBitMap()
}
