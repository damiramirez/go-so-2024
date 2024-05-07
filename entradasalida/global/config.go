package global

import (
	"fmt"
	"os"

	config "github.com/sisoputnfrba/tp-golang/utils/config"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
)

const IOLOG = "./entradasalida.log"

type Config struct {
	Port             int    `json:"port"`
	Type             string `json:"type"`
	UnitWorkTime     int    `json:"unit_work_time"`
	IPKernel         string `json:"ip_kernel"`
	PortKernel       int    `json:"port_kernel"`
	IPMemory         string `json:"ip_memory"`
	PortMemory       int    `json:"port_memory"`
	DialFSPath       string `json:"dialfs_path"`
	DialFSBlockSize  int    `json:"dialfs_block_size"`
	DialFSBlockCount int    `json:"dialfs_block_count"`
}

type GenericIODevice struct {
	Name  string
	Type  string
	InUse bool
	Port  int
}

var Dispositivo *GenericIODevice

var IOConfig *Config

var Logger *log.LoggerStruct

func InitGlobal() {
	args := os.Args[1:]
	if len(args) != 3 {
		fmt.Println("Uso: programa <go run `modulo`.go dev|prod>")
		os.Exit(1)
	}
	env := args[0]
	name := args[1]
	configuracion := args[2]

	Logger = log.ConfigureLogger(IOLOG, env)
	IOConfig = config.LoadConfiguration[Config](configuracion)

	Dispositivo = InitGenericIODevice(name)

	AvisoKernelIOExistente()

}

func InitGenericIODevice(name string) *GenericIODevice {

	dispositivo := GenericIODevice{Name: name, Type: IOConfig.Type, Port: IOConfig.Port}

	Logger.Log(fmt.Sprintf("Nuevo IO genérico inicializado: %+v", dispositivo), log.INFO)

	return &dispositivo

}

func AvisoKernelIOExistente() {

	_, err := requests.PutHTTPwithBody[GenericIODevice, interface{}](IOConfig.IPKernel, IOConfig.PortKernel, "newio", *Dispositivo)
	if err != nil {
		Logger.Log(fmt.Sprintf("NO se pudo enviar al kernel el IODevice %s", err.Error()), log.INFO)
		panic(1)
		// TODO: kernel falta que entienda el mensaje (hacer el endpoint) y nos envíe la respuesta que está todo ok
	}

}
