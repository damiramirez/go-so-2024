package global

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"

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

var FilesMap map[string]File

var Bloques []byte

var Bitmap []byte

var Estructura_truncate KernelIOFS_Truncate

// var Filestruct File

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

	FilesMap = map[string]File{}

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

		// crear directorio específico de la IO
		dir := config.DialFSPath + "/" + Dispositivo.Name

		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			Logger.Log(fmt.Sprintf("Error al crear el directorio: %s", err.Error()), log.ERROR)
			return
		}

		// crear-abrir bloques.dat

		openBloquesDat(config)

		// crear-abrir bitmap.dat

		openBitmapDat(config)

		// reconstruir mi map de files

		rebuildFilesMap(config)

	}

}

func openBloquesDat(config *Config) {

	filename := IOConfig.DialFSPath + "/" + Dispositivo.Name + "/bloques.dat"
	size := config.DialFSBlockSize * config.DialFSBlockCount
	Bloques = make([]byte, IOConfig.DialFSBlockCount*IOConfig.DialFSBlockSize)

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

	_, err = file.Read(Bloques)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al leer el archivo: %s ", err.Error()), log.ERROR)
	}

	Logger.Log(fmt.Sprintf("Archivo %s abierto con éxito (tamaño de %d bytes): %+v", filename, size, Bloques), log.DEBUG)
}

func openBitmapDat(config *Config) {

	filename := IOConfig.DialFSPath + "/" + Dispositivo.Name + "/bitmap.dat"
	size := config.DialFSBlockCount
	Bitmap = make([]byte, IOConfig.DialFSBlockCount)

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

	_, err = file.Read(Bitmap)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al leer el archivo: %s ", err.Error()), log.ERROR)
	}

	Logger.Log(fmt.Sprintf("Archivo %s abierto con éxito (tamaño de %d bits): %+v", filename, size, Bitmap), log.DEBUG)
}

func GetCurrentBlocks(file string, w http.ResponseWriter) int {

	filestruct := FilesMap[file]
	Logger.Log(fmt.Sprintf("Filestruct %s: %+v", file, filestruct), log.DEBUG)
	if filestruct.Size > 0 {
		filestruct.CurrentBlocks = int(math.Ceil(float64(filestruct.Size) / float64(IOConfig.DialFSBlockSize)))
	} else if filestruct.Size == 0 {
		filestruct.CurrentBlocks = 1
	}
	Logger.Log(fmt.Sprintf("Current blocks: %d", filestruct.CurrentBlocks), log.DEBUG)
	return filestruct.CurrentBlocks

}

func GetFreeContiguousBlocks(file string, w http.ResponseWriter) int {

	currentBlocks := GetCurrentBlocks(file, w)

	freeContiguousBlocks := 0

	bitmappath := IOConfig.DialFSPath + "/" + Dispositivo.Name + "/bitmap.dat"

	bitmapfile, err := os.OpenFile(bitmappath, os.O_RDWR, 0644)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al abrir el archivo: %s ", err.Error()), log.ERROR)
		http.Error(w, "Error al abrir el archivo", http.StatusBadRequest)
		return -1
	}

	defer bitmapfile.Close()

	filestruct := FilesMap[file]

	_, err = bitmapfile.Seek(int64(filestruct.Initial_block+currentBlocks), 0)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al mover el cursor: %s ", err.Error()), log.ERROR)
		http.Error(w, "Error al mover el cursor", http.StatusBadRequest)
		return -1
	}
	value := make([]byte, 1)

	bitmapfile.Read(value)

	for value[0] != 1 && filestruct.Initial_block+currentBlocks+freeContiguousBlocks <= IOConfig.DialFSBlockCount-1 {

		freeContiguousBlocks++
		_, err = bitmapfile.Seek(int64(filestruct.Initial_block+currentBlocks+freeContiguousBlocks), 0)
		if err != nil {
			Logger.Log(fmt.Sprintf("Error al mover el cursor: %s ", err.Error()), log.ERROR)
			http.Error(w, "Error al mover el cursor", http.StatusBadRequest)
			return -1
		}

		bitmapfile.Read(value)
	}
	Logger.Log(fmt.Sprintf("Free contiguous blocks: %d ", freeContiguousBlocks), log.DEBUG)
	return freeContiguousBlocks
}

func GetNeededBlocks(w http.ResponseWriter, estructura KernelIOFS_Truncate) int {

	var neededBlocks int

	if estructura.Tamanio == 0 {
		neededBlocks = 1
	} else {
		neededBlocks = int(math.Ceil((float64(estructura.Tamanio) / float64(IOConfig.DialFSBlockSize))))
	}
	Logger.Log(fmt.Sprintf("Needed blocks: %d ", neededBlocks), log.DEBUG)
	return neededBlocks
}

func GetTotalFreeBlocks(w http.ResponseWriter) int { // no usar archivo, usar el slice de bytes

	bitmappath := IOConfig.DialFSPath + "/" + Dispositivo.Name + "/bitmap.dat"

	bitmapfile, err := os.OpenFile(bitmappath, os.O_RDWR, 0644)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al abrir el archivo: %s ", err.Error()), log.ERROR)
		http.Error(w, "Error al abrir el archivo", http.StatusBadRequest)
		return -1
	}

	defer bitmapfile.Close()

	_, err = bitmapfile.Seek(0, 0)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al mover el cursor: %s ", err.Error()), log.ERROR)
		http.Error(w, "Error al mover el cursor", http.StatusBadRequest)
		return -1
	}

	value := make([]byte, 1)
	var totalFreeBlocks int = 0
	var i int = -1
	for i < IOConfig.DialFSBlockCount-2 {

		if value[0] == 0 {
			totalFreeBlocks++
		}
		i++
		_, err = bitmapfile.Seek(int64(i+1), 0)
		if err != nil {
			Logger.Log(fmt.Sprintf("Error al mover el cursor: %s ", err.Error()), log.ERROR)
			http.Error(w, "Error al mover el cursor", http.StatusBadRequest)
			return -1
		}

		bitmapfile.Read(value)
	}
	Logger.Log(fmt.Sprintf("Total free blocks: %d ", totalFreeBlocks), log.DEBUG)
	return totalFreeBlocks
}

func PrintBitmap(w http.ResponseWriter) {

	// leo el archivo y logeo su contenido

	bitmappath := IOConfig.DialFSPath + "/" + Dispositivo.Name + "/bitmap.dat"

	bitmapfile, err := os.OpenFile(bitmappath, os.O_RDWR, 0644)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al abrir el archivo: %s ", err.Error()), log.ERROR)
		http.Error(w, "Error al abrir el archivo", http.StatusBadRequest)
		return
	}

	defer bitmapfile.Close()

	_, err = bitmapfile.Read(Bitmap)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al leer el archivo: %s ", err.Error()), log.ERROR)
		http.Error(w, "Error al leer el archivo", http.StatusBadRequest)
		return
	}
	Logger.Log(fmt.Sprintf("Bitmap del FS: %+v", Bitmap), log.DEBUG)

}

func UpdateSize(file string, newSize int, w http.ResponseWriter) { // modificar el size en el metadata

	filepath := IOConfig.DialFSPath + "/" + Dispositivo.Name + "/" + file

	metadatafile, err := os.OpenFile(filepath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al abrir el archivo %s: %s ", filepath, err.Error()), log.ERROR)
		http.Error(w, "Error al abrir el archivo", http.StatusBadRequest)
		return
	}

	Logger.Log(fmt.Sprintf("NewSize: %d", newSize), log.DEBUG)

	defer metadatafile.Close()

	filestruct := FilesMap[file]

	newSizemap := map[string]interface{}{
		"initial_block": filestruct.Initial_block,
		"size":          newSize,
	}

	encoder := json.NewEncoder(metadatafile)
	err = encoder.Encode(newSizemap)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al encodear el nuevo size en el archivo %s: %s ", filepath, err.Error()), log.ERROR)
		http.Error(w, "Error al encodear el nuevo size en el archivo", http.StatusBadRequest)
		return
	}

	filestruct.Size = newSize
	FilesMap[file] = filestruct
}

func UpdateInitialBlock(file string, newInitialBlock int, w http.ResponseWriter) { // modificar el initial block en el metadata

	// declarar variable filestruct
	// accedo a map con la key file
	// actualizar map con la key file

	filepath := IOConfig.DialFSPath + "/" + Dispositivo.Name + "/" + file

	metadatafile, err := os.OpenFile(filepath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al abrir el archivo %s: %s ", filepath, err.Error()), log.ERROR)
		http.Error(w, "Error al abrir el archivo", http.StatusBadRequest)
		return
	}

	defer metadatafile.Close()

	filestruct := FilesMap[file]

	newInitialBlockmap := map[string]interface{}{
		"initial_block": newInitialBlock,
		"size":          filestruct.Size,
	}

	encoder := json.NewEncoder(metadatafile)
	err = encoder.Encode(newInitialBlockmap)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al encodear el nuevo initial block en el archivo %s: %s ", filepath, err.Error()), log.ERROR)
		http.Error(w, "Error al encodear el nuevo initial block en el archivo", http.StatusBadRequest)
		return
	}

	filestruct.Initial_block = newInitialBlock
	FilesMap[file] = filestruct
}

func UpdateBitmap(writeValue int, initialBit int, bitAmount int, w http.ResponseWriter) {

	// actualizo el slice de bytes

	for i := 0; i < bitAmount; i++ {
		Bitmap[initialBit+i] = byte(writeValue)
	}

	// actualizo el archivo bitmap.dat

	bitmappath := IOConfig.DialFSPath + "/" + Dispositivo.Name + "/bitmap.dat"

	bitmapfile, err := os.OpenFile(bitmappath, os.O_RDWR, 0644)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al abrir el archivo: %s ", err.Error()), log.ERROR)
		http.Error(w, "Error al abrir el archivo", http.StatusBadRequest)
		return
	}

	defer bitmapfile.Close()

	_, err = bitmapfile.Write(Bitmap)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al actualizar el bitmap: %s ", err.Error()), log.ERROR)
		http.Error(w, "Error al actualizar el bitmap", http.StatusBadRequest)
		return
	}

}

func UpdateBlocksFile(newInfo []byte, file string, punteroArchivo int, w http.ResponseWriter) {

	// convierto el archivo bloques.dat en un slice de bytes
	/*
		data, err := os.ReadFile(filepath)
		if err != nil {
			Logger.Log(fmt.Sprintf("Error al leer el archivo: %s ", err.Error()), log.ERROR)
			return
		}
	*/

	filestruct := FilesMap[file]

	bloquesdatpath := IOConfig.DialFSPath + "/" + Dispositivo.Name + "/bloques.dat"

	Logger.Log(fmt.Sprintf("Contenido original del slice bloques: %+v", Bloques), log.DEBUG)

	//Actualizo slice Bloques

	ubicacionDeseada := IOConfig.DialFSBlockSize*filestruct.Initial_block + punteroArchivo

	for i := 0; i < len(newInfo); i++ {
		Bloques[ubicacionDeseada+i] = newInfo[i]

	}

	Logger.Log(fmt.Sprintf("Contenido del slice modificado: %+v", Bloques), log.DEBUG)

	//Actualizo el archivo bloques.dat con el slice Bloques

	err := os.WriteFile(bloquesdatpath, Bloques, 0644)
	if err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
		return
	}

	Logger.Log(fmt.Sprintf("El archivo se actualizó con éxito: %+v", Bloques), log.INFO)

}

func rebuildFilesMap(config *Config) {

	// recorrer el directorio en busca de metadatas existentes

	dirPath := config.DialFSPath + "/" + Dispositivo.Name

	// Leer los contenidos del directorio
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		Logger.Log(fmt.Sprintf("No se pudo leer el directorio que contiene los metadata %s", err.Error()), log.ERROR)
	}
	// Iterar sobre los archivos y agregarlos al map
	for _, entry := range entries {
		if !entry.IsDir() && strings.Contains(entry.Name(), "txt") {
			addToFilesMapDecoding(entry.Name())
		}
	}

	Logger.Log(fmt.Sprintf("FilesMap %+v", FilesMap), log.DEBUG)
}

func addToFilesMapDecoding(file string) {

	var filestruct File

	// obtengo los datos del archivo metadata

	metadatapath := IOConfig.DialFSPath + "/" + Dispositivo.Name + "/" + file

	metadatafile, err := os.Open(metadatapath)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al abrir el archivo %s: %s ", metadatapath, err.Error()), log.DEBUG)
		return
	}

	defer metadatafile.Close()

	decoder := json.NewDecoder(metadatafile)
	err = decoder.Decode(&filestruct)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al decodear el archivo %s: %s ", metadatapath, err.Error()), log.ERROR)
		return
	}

	FilesMap[file] = filestruct

	// actualizar el valor de CurrentBlocks de filestruct (usar GetCurrentBlocks?)

	Logger.Log(fmt.Sprintf("Filestruct recién decodeado: %+v", filestruct), log.DEBUG)

}

/*
func TruncateLess(file string, w http.ResponseWriter) {

		filepath := IOConfig.DialFSPath + "/" + Dispositivo.Name  "/" + file

		bitmappath := IOConfig.DialFSPath + "/" + Dispositivo.Name  "/bitmap.dat"

		bitmapfile, err := os.OpenFile(bitmappath, os.O_RDWR, 0644)
		if err != nil {
			Logger.Log(fmt.Sprintf("Error al abrir el archivo: %s ", err.Error()), log.ERROR)
			http.Error(w, "Error al abrir el archivo", http.StatusBadRequest)
			return
		}

		defer bitmapfile.Close() // esta línea de código garantiza que el archivo en el que estoy trabajando se cierre cuando la función actual termina de ejecutarse

		// leo el archivo y logeo su contenido

		data := make([]byte, IOConfig.DialFSBlockCount) // crea un slice de bytes de tamaño global.IOConfig.DialFSBlockCount, en el cual asigno los bytes que leo del archivo bitmapfile
		_, err = bitmapfile.Read(data)
		if err != nil {
			Logger.Log(fmt.Sprintf("Error al leer el archivo: %s ", err.Error()), log.ERROR)
			http.Error(w, "Error al leer el archivo", http.StatusBadRequest)
			return
		}
		Logger.Log(fmt.Sprintf("Bitmap del FS %s antes de truncar: %+v", Dispositivo.Name, data), log.DEBUG)

		Logger.Log(fmt.Sprintf("Datos del archivo %s antes de truncar: %+v ", filepath, Filestruct), log.DEBUG)

		neededBlocks := GetNeededBlocks(w, Estructura_truncate)

		currentBlocks := GetCurrentBlocks(file, w)

		Logger.Log(fmt.Sprintf("Current Blocks: %d", currentBlocks), log.DEBUG)

		Logger.Log(fmt.Sprintf("Needed Blocks: %d", neededBlocks), log.DEBUG)

		for i := 0; i < currentBlocks-neededBlocks; i++ {

			_, err = bitmapfile.Seek(int64(Filestruct.Initial_block+neededBlocks+i), 0)
			if err != nil {
				Logger.Log(fmt.Sprintf("Error al mover el cursor: %s ", err.Error()), log.ERROR)
				http.Error(w, "Error al mover el cursor", http.StatusBadRequest)
				return
			}

			// cambio el bit de 1 a 0
			_, err = bitmapfile.Write([]byte{0})
			if err != nil {
				Logger.Log(fmt.Sprintf("Error al escribir el byte: %s ", err.Error()), log.ERROR)
				http.Error(w, "Error al escribir el byte", http.StatusBadRequest)
				return
			}
		}

		Filestruct.Size = Estructura_truncate.Tamanio
		Filestruct.CurrentBlocks = neededBlocks

		Logger.Log(fmt.Sprintf("Datos del archivo %s luego de truncar: %+v ", filepath, Filestruct), log.DEBUG)

		// muevo el cursor nuevamente al principio del archivo bitmap
		_, err = bitmapfile.Seek(0, 0)
		if err != nil {
			Logger.Log(fmt.Sprintf("Error al mover el cursor: %s ", err.Error()), log.ERROR)
			return
		}

		// leo el archivo (desde la posición inicial) y logeo su contenido actualizado

		_, err = bitmapfile.Read(data) // asigno los bytes que leo del archivo bitmapfile (actualizado) a mi slice de bytes data, creado anteriormente
		if err != nil {
			Logger.Log(fmt.Sprintf("Error al leer el archivo: %s ", err.Error()), log.ERROR)
			http.Error(w, "Error al leer el archivo", http.StatusBadRequest)
			return
		}

		Logger.Log(fmt.Sprintf("Bitmap del FS %s luego de truncar: %+v", Dispositivo.Name, data), log.DEBUG)
	}
*/

/*
func TruncateMore(file string, w http.ResponseWriter) {

	filepath := IOConfig.DialFSPath + "/" + Dispositivo.Name  "/" + file

	bitmappath := IOConfig.DialFSPath + "/" + Dispositivo.Name  "/bitmap.dat"

	bitmapfile, err := os.OpenFile(bitmappath, os.O_RDWR, 0644)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al abrir el archivo: %s ", err.Error()), log.ERROR)
		http.Error(w, "Error al abrir el archivo", http.StatusBadRequest)
		return
	}

	defer bitmapfile.Close() // esta línea de código garantiza que el archivo en el que estoy trabajando se cierre cuando la función actual termina de ejecutarse

	// leo el archivo y logeo su contenido
	_, err = bitmapfile.Read(Bitmap)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al leer el archivo: %s ", err.Error()), log.ERROR)
		http.Error(w, "Error al leer el archivo", http.StatusBadRequest)
		return
	}
	Logger.Log(fmt.Sprintf("Bitmap del FS %s antes de truncar: %+v", Dispositivo.Name, Bitmap), log.DEBUG)

	Logger.Log(fmt.Sprintf("Datos del archivo %s antes de truncar: %+v ", filepath, Filestruct), log.DEBUG)

	neededBlocks := GetNeededBlocks(w, Estructura_truncate)

	currentBlocks := GetCurrentBlocks(file, w)

	Logger.Log(fmt.Sprintf("Current Blocks: %d", currentBlocks), log.DEBUG)

	Logger.Log(fmt.Sprintf("Needed Blocks: %d", neededBlocks), log.DEBUG)

	for i := 0; i < neededBlocks; i++ {

		_, err = bitmapfile.Seek(int64(Filestruct.Initial_block+i), 0)
		if err != nil {
			Logger.Log(fmt.Sprintf("Error al mover el cursor: %s ", err.Error()), log.ERROR)
			http.Error(w, "Error al mover el cursor", http.StatusBadRequest)
			return
		}

		// cambio el bit de 0 a 1 (ver qué pasa si esa posición ya está ocupada, fragmentación externa, compactación)
		_, err = bitmapfile.Write([]byte{1})
		if err != nil {
			Logger.Log(fmt.Sprintf("Error al escribir el byte: %s ", err.Error()), log.ERROR)
			http.Error(w, "Error al escribir el byte", http.StatusBadRequest)
			return
		}
	}

	Filestruct.Size = Estructura_truncate.Tamanio
	Filestruct.CurrentBlocks = neededBlocks

	Logger.Log(fmt.Sprintf("Datos del archivo %s luego de truncar: %+v ", filepath, Filestruct), log.DEBUG)

	// muevo el cursor nuevamente al principio del archivo bitmap
	_, err = bitmapfile.Seek(0, 0)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al mover el cursor: %s ", err.Error()), log.ERROR)
		return
	}

	// leo el archivo (desde la posición inicial) y logeo su contenido actualizado

	_, err = bitmapfile.Read(Bitmap) // asigno los bytes que leo del archivo bitmapfile (actualizado) a mi slice de bytes data, creado anteriormente
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al leer el archivo: %s ", err.Error()), log.ERROR)
		http.Error(w, "Error al leer el archivo", http.StatusBadRequest)
		return
	}

	Logger.Log(fmt.Sprintf("Bitmap del FS %s luego de truncar: %+v", Dispositivo.Name, Bitmap), log.DEBUG)

}
*/

/*
func openTruncatedFilesDirectory(config *Config) {

	// crear carpeta para los archivos del FS que fueron truncados
	dir := config.DialFSPath + "/" + Dispositivo.Name + "/" + "truncated-files"

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		Logger.Log(fmt.Sprintf("Error al crear el directorio: %s", err.Error()), log.ERROR)
		return
	}

	Logger.Log(fmt.Sprintf("Archivo %s abierto con éxito", dir), log.DEBUG)
}

func AddToTruncatedFiles(file string) {

	// crear carpeta para los archivos del FS que fueron truncados
	dir := IOConfig.DialFSPath + "/" + Dispositivo.Name + "/" + "truncated-files"

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		Logger.Log(fmt.Sprintf("Error al crear el directorio: %s", err.Error()), log.ERROR)
		return
	}

	truncatedpath := IOConfig.DialFSPath + "/" + Dispositivo.Name + "/truncated-files/truncated-" + file

	truncatedfile, err := os.OpenFile(truncatedpath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al crear el archivo: %s", err.Error()), log.ERROR)
	}

	defer truncatedfile.Close()
}

/*
func hasBeenTruncated(file string) int { // 1 si fue truncado, 0 si no lo fue

	dirPath := IOConfig.DialFSPath + "/" + Dispositivo.Name + "/" + "truncated-files"

	dir, err := os.Open(dirPath)
	if err != nil {
		fmt.Printf("Error al abrir el directorio %s: %s\n", dirPath, err.Error())
		return 0
	}
	defer dir.Close()

	fileNames, err := dir.Readdirnames(0)
	if err != nil {
		fmt.Printf("Error al leer los nombres de los archivos en el directorio %s: %s", dirPath, err.Error())
		return 0
	}

	// Comprobar si el archivo específico existe
	for _, fName := range fileNames {
		if fName == "truncated-"+file {
			Logger.Log(fmt.Sprintf("El archivo %s ha sido truncado anteriormente", file), log.DEBUG)
			return 1
		}
	}

	Logger.Log(fmt.Sprintf("El archivo %s no ha sido truncado anteriormente", file), log.DEBUG)
	return 0

func AddToActiveFiles(file string) {

	// crear carpeta para los archivos del FS que están activos
	dir := IOConfig.DialFSPath + "/" + Dispositivo.Name + "/" + "active-files"

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		Logger.Log(fmt.Sprintf("Error al crear el directorio: %s", err.Error()), log.ERROR)
		return
	}

	activepath := IOConfig.DialFSPath + "/" + Dispositivo.Name + "/active-files/" + "active-" + file

	activefile, err := os.OpenFile(activepath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		Logger.Log(fmt.Sprintf("Error al crear el archivo: %s", err.Error()), log.ERROR)
	}

	defer activefile.Close()
}


func openActiveFilesDirectory(config *Config) {

	// crear carpeta para los archivos del FS que están activos
	dir := config.DialFSPath + "/" + Dispositivo.Name + "/" + "active-files"

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		Logger.Log(fmt.Sprintf("Error al crear el directorio: %s", err.Error()), log.ERROR)
		return
	}

	Logger.Log(fmt.Sprintf("Archivo %s abierto con éxito", dir), log.DEBUG)
}
*/
