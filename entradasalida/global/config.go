package global

import (
	"fmt"
	"os"

	"github.com/sisoputnfrba/tp-golang/entradasalida/global"
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

type IODevice struct {
	Name  string
	Type  string
	InUse bool
	Port  int
}

var Dispositivo *IODevice

var Texto string

var IOConfig *Config

var Logger *log.LoggerStruct

func InitGlobal() {
	args := os.Args[1:]
	if len(args) != 3 {
		fmt.Println("Uso: programa <go run `modulo`.go dev|prod N=name P=path>")
		os.Exit(1)
	}
	env := args[0]
	name := args[1]
	configuracion := args[2]

	Logger = log.ConfigureLogger(IOLOG, env)
	IOConfig = config.LoadConfiguration[Config](configuracion)

	Dispositivo = InitIODevice(name)

	AvisoKernelIOExistente()

}

func InitIODevice(name string) *IODevice {

	dispositivo := IODevice{Name: name, Type: IOConfig.Type, Port: IOConfig.Port}

	Logger.Log(fmt.Sprintf("Nuevo IO inicializado: %+v", dispositivo), log.INFO)

	return &dispositivo

}

func AvisoKernelIOExistente() {

	_, err := requests.PutHTTPwithBody[IODevice, interface{}](IOConfig.IPKernel, IOConfig.PortKernel, "newio", *Dispositivo)
	if err != nil {
		Logger.Log(fmt.Sprintf("NO se pudo enviar al kernel el IODevice %s", err.Error()), log.INFO)
		panic(1)
		// TODO: kernel falta que entienda el mensaje (hacer el endpoint) y nos envíe la respuesta que está todo ok
	}

}

func VerificacionTamanio(texto string, tamanio string) {

	if estructura_actualizada.Tamanio == global.Texto { // TO DO: implementar la comparacion, hacer estructura global?
		estructura_actualizada.Texto = global.Texto
		return
	}

	global.Logger.Log(fmt.Sprintf("Tamaño incorrecto, ingrese un nuevo valor de tamaño (%s", estructura.Tamanio)+"): ", log.INFO)

	fmt.Scanf("%s", &global.Texto)

	VerificacionTamanio(global.Texto, estructura_actualizada.Tamanio)

}
