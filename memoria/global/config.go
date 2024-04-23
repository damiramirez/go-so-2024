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
var ConfigPag Config 
var numPages = ConfigPag.MemorySize/ConfigPag.PageSize

var ByteArray = make([]byte,ConfigPag.PageSize)
type Page struct {
    data []byte
}
type PageTable []*Page
type Memory struct {
    pages PageTable
}
// Se inicializa cada página de la memoria con datos vacíos
func NewMemory() *Memory {
    mem := Memory{}
    for i := 0; i < numPages; i++ {
        mem.pages[i] = &Page{}
    }
    return &mem
}
//escribe en memoria 
// falta desarollar la funcion 
func WriteinMemory(data []byte,mem *Memory) {
    /*pageIndex := address / ConfigPag.PageSize
    offset := address % ConfigPag.PageSize
    if pageIndex < numPages && offset+len(data) <= ConfigPag.PageSize {
        copy(mem.pages[pageIndex].data[offset:], data)
    } else {
        fmt.Println("Error: Dirección de memoria fuera de rango")
    }*/
    Logger.Log("Se escribio en memoria ", log.DEBUG)
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
}