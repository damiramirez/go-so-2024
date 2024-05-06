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

var IOConfig *Config

var Logger *log.LoggerStruct

var MapIOGenericsActivos map[string]GenericIODevice

func InitGlobal() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Uso: programa <go run `modulo`.go dev|prod>")
		os.Exit(1)
	}
	env := args[0]

	Logger = log.ConfigureLogger(IOLOG, env)
	IOConfig = config.LoadConfiguration[Config]("./config/config.json")

	/*
		Esperador := GenericIODevice{Name: "esperador", Type: "GENERIC"}
		Teclado := GenericIODevice{Name: "teclado", Type: "STDIN"}
		Pantalla := GenericIODevice{Name: "pantalla", Type: "STDOUT"}
	*/
	MapIOGenericsActivos = map[string]GenericIODevice{}

	InitGenericIODevice("Int1")
	InitGenericIODevice("teclado")
	/*
		MapIOGenericsActivos[Esperador.Name] = Esperador
		MapIOGenericsActivos[Teclado.Name] = Teclado
		MapIOGenericsActivos[Pantalla.Name] = Pantalla
	*/
	// AvisoKernelIOExistentes()

}

type GenericIODevice struct {
	Name      string
	Type      string
	EstaEnUso bool
}

func InitGenericIODevice(name string) {

	IODevice := GenericIODevice{Name: name, Type: IOConfig.Type}
	MapIOGenericsActivos[IODevice.Name] = IODevice
	Logger.Log(fmt.Sprintf("Nuevo IO genérico inicializado: %+v", IODevice), log.INFO)
	Logger.Log(fmt.Sprintf("Lista de IOs inicializados: %+v", MapIOGenericsActivos), log.INFO)
}

func AvisoKernelIOExistentes() {

	listaIODeviceRegistrados := []GenericIODevice{}

	for _, value := range MapIOGenericsActivos {

		listaIODeviceRegistrados = append(listaIODeviceRegistrados, value)

	}

	Logger.Log(fmt.Sprintf("%+v", listaIODeviceRegistrados), log.INFO)

	_, err := requests.PutHTTPwithBody[[]GenericIODevice, interface{}](IOConfig.IPKernel, IOConfig.PortKernel, "listIO", listaIODeviceRegistrados)
	if err != nil {
		Logger.Log(fmt.Sprintf("NO se pudo enviar al kernel los IODevices %s", err.Error()), log.INFO)
		panic(1)
		// TODO: kernel falta que entienda el mensaje y nos envíe la respuesta que está todo ok
	}

}
