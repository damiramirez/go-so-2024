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

var MemoryConfig *Config
var Logger *log.LoggerStruct

var NumPages int

type MemoryST struct {
    pages []byte
}
var Memory *MemoryST
// Se inicializa cada página de la memoria con datos vacíos
func NewMemory() *MemoryST {
	NumPages = MemoryConfig.MemorySize/MemoryConfig.PageSize
	var ByteArray = make([]byte,MemoryConfig.PageSize)
    mem := MemoryST{pages: ByteArray}
    return &mem
}
	

func InitGlobal() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Uso: programa <go run `modulo`.go dev|prod>")
		os.Exit(1)
	}
	env := args[0]

	Logger = log.ConfigureLogger(MEMORYLOG, env)
	MemoryConfig = config.LoadConfiguration[Config]("./config/config.json")
	Memory=NewMemory()
}