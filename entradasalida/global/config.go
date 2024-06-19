package global

import (
	"bufio"
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

type IODevice struct {
	Name  string
	Type  string
	InUse bool
	Port  int
}

type Estructura_sleep struct {
	Name        string `json:"nombre"`
	Instruction string `json:"instruccion"`
	Time        int    `json:"tiempo"`
	Pid         int    `json:"pid"`
}
type ValoraMandar struct {
	Texto string `json:"texto"`
}
type MemStdIO struct {
	Pid       int    `json:"pid"`
	Content   string `json:"content"`
	Length    int    `json:"length"`
	NumFrames []int  `json:"numframe"`
	Offset    int    `json:"offset"`
}

type KernelIOStd struct {
	Pid         int    `json:"pid"`
	Instruction string `json:"instruccion"`
	Name        string `json:"name"`
	Length      int    `json:"length"`
	NumFrames   []int  `json:"numframe"`
	Offset      int    `json:"offset"`
}

type KernelIOFS_CD struct {
	Pid         int    `json:"pid"`
	Instruction string `json:"instruccion"`
	IOName      string `json:"nombre"`
	FileName    string `json:"filename"`
}

type KernelIOFS_Truncate struct {
	Pid         int    `json:"pid"`
	Instruction string `json:"instruccion"`
	IOName      string `json:"nombre"`
	FileName    string `json:"filename"`
	Tamanio     int    `json:"tamanio"`
}

type KernelIOFS_WR struct {
	Pid            int    `json:"pid"`
	Instruction    string `json:"instruccion"`
	IOName         string `json:"nombre"`
	FileName       string `json:"filename"`
	NumFrames      []int  `json:"numframe"`
	Offset         int    `json:"offset"`
	Tamanio        int    `json:"tamanio"`
	PunteroArchivo int    `json:"punteroArchivo"`
}

type File struct {
	Initial_block int `json:"initial_block"`
	Size          int `json:"size"`
	CurrentBlocks int
}

var Filestruct File

var Estructura_actualizada MemStdIO

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

	LevantarFS(IOConfig)

}

func InitIODevice(name string) *IODevice {

	dispositivo := IODevice{Name: name, Type: IOConfig.Type, Port: IOConfig.Port}

	Logger.Log(fmt.Sprintf("Nuevo IO inicializado: %+v", dispositivo), log.DEBUG)

	return &dispositivo

}

func AvisoKernelIOExistente() {

	_, err := requests.PutHTTPwithBody[IODevice, interface{}](IOConfig.IPKernel, IOConfig.PortKernel, "newio", *Dispositivo)
	if err != nil {
		Logger.Log(fmt.Sprintf("NO se pudo enviar al kernel el IODevice %s", err.Error()), log.ERROR)
		panic(1)
		// TODO: kernel falta que entienda el mensaje (hacer el endpoint) y nos envíe la respuesta que está todo ok
	}

}

func VerificacionTamanio(texto string, tamanio int) {

	BtT := []byte(Texto)

	Logger.Log(fmt.Sprintf("Slice de bytes: %+v", BtT), log.DEBUG)

	if len(BtT) == 0 {

		Logger.Log(fmt.Sprintf("No ingresó nada, ingrese un nuevo valor (tamaño máximo %d", tamanio)+"): ", log.INFO)

		reader := bufio.NewReader(os.Stdin)
		Texto, _ = reader.ReadString('\n')

		VerificacionTamanio(Texto, tamanio)
	}

	if len(BtT) <= tamanio+1 {
		Estructura_actualizada.Content = Texto[:len(BtT)-1]
		return
	}

	Logger.Log(fmt.Sprintf("Tamaño excedido, ingrese un nuevo valor (tamaño máximo %d", tamanio)+"): ", log.INFO)

	reader := bufio.NewReader(os.Stdin)
	Texto, _ = reader.ReadString('\n')

	VerificacionTamanio(Texto, tamanio)

}

func LevantarFS(config *Config) {

	if config.Type == "DIALFS" {

		// crear carpeta para los archivos del FS
		dir := config.DialFSPath + "/Filesystems" + "/" + Dispositivo.Name

		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			Logger.Log(fmt.Sprintf("Error al crear el directorio: %v", err), log.ERROR)
			return
		}

		// crear bloques.dat

		openBloquesDat(config)

		// crear bitmap.dat

		openBitmapDat(config)

	}

}

func openBloquesDat(config *Config) {

	filename := config.DialFSPath + "/Filesystems" + "/" + Dispositivo.Name + "/bloques.dat"
	size := config.DialFSBlockSize * config.DialFSBlockCount

	// crear el archivo
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		Logger.Log(fmt.Sprint("Error al crear el archivo:", err), log.ERROR)
		return
	}

	// cerrar el archivo
	defer file.Close()

	// ajustar el tamaño del archivo
	err = file.Truncate(int64(size))
	if err != nil {
		Logger.Log(fmt.Sprint("Error al ajustar el tamaño del archivo:", err), log.ERROR)
		return
	}

	data := make([]byte, IOConfig.DialFSBlockCount*IOConfig.DialFSBlockCount) // crea un slice de bytes de tamaño global.IOConfig.DialFSBlockCount*IOConfig.DialFSBlockCount, en el cual asigno los bytes que leo del archivo bloques.dat
	_, err = file.Read(data)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al leer el archivo: %s ", err.Error()), log.ERROR)
	}

	Logger.Log(fmt.Sprintf("Archivo %s abierto con éxito (tamaño de %d bytes): %+v", filename, size, data), log.DEBUG)
}

func openBitmapDat(config *Config) {

	filename := config.DialFSPath + "/Filesystems" + "/" + Dispositivo.Name + "/bitmap.dat"
	size := config.DialFSBlockCount

	// crear el archivo
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		Logger.Log(fmt.Sprint("Error al crear el archivo:", err), log.ERROR)
		return
	}

	// cerrar el archivo
	defer file.Close()

	// ajustar el tamaño del archivo
	err = file.Truncate(int64(size))
	if err != nil {
		Logger.Log(fmt.Sprint("Error al ajustar el tamaño del archivo:", err), log.ERROR)
		return
	}

	data := make([]byte, IOConfig.DialFSBlockCount) // crea un slice de bytes de tamaño global.IOConfig.DialFSBlockCount, en el cual asigno los bytes que leo del archivo bitmapfile
	_, err = file.Read(data)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al leer el archivo: %s ", err.Error()), log.ERROR)
	}

	Logger.Log(fmt.Sprintf("Archivo %s abierto con éxito (tamaño de %d bytes): %+v", filename, size, data), log.DEBUG)
}
