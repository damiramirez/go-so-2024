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
type ProcessInstructions struct {
	Instructions []string
}
type ListInstructions ProcessInstructions
var DictProcess map[int]ListInstructions
type MemoryST struct {
	Spaces []byte
}
type PageTable struct {
	Page []byte
}
var Memory *MemoryST
func NewMemory() *MemoryST {
	
	ByteArray := make([]byte,MemoryConfig.MemorySize)
    mem := MemoryST{Spaces: ByteArray}
    return &mem
}
var PTable *PageTable
func NewPageTable()*PageTable{
    
	ByteArray := make([]byte,MemoryConfig.PageSize)
	pagetable:=PageTable{Page: ByteArray}

	return &pagetable
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
	DictProcess=map[int]ListInstructions{}
	Memory=NewMemory()
	PTable=NewPageTable()
}
