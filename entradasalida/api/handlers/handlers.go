package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sisoputnfrba/tp-golang/entradasalida/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

func Sleep(w http.ResponseWriter, r *http.Request) {
	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	var estructura global.Estructura_sleep
	err := serialization.DecodeHTTPBody[*global.Estructura_sleep](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s>", estructura.Pid, estructura.Instruction), log.INFO)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.DEBUG)

	global.Logger.Log(fmt.Sprintf("a punto de dormir: %+v", dispositivo), log.DEBUG)

	global.Logger.Log(fmt.Sprintf("durmiendo: %+v", dispositivo), log.DEBUG)

	time.Sleep(time.Duration(estructura.Time*global.IOConfig.UnitWorkTime) * time.Millisecond)

	global.Logger.Log(fmt.Sprintf("terminé de dormir: %+v", dispositivo), log.DEBUG)

	w.WriteHeader(http.StatusNoContent)
	dispositivo.InUse = false
}

func Stdin_read(w http.ResponseWriter, r *http.Request) {
	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	var estructura global.KernelIOStd

	err := serialization.DecodeHTTPBody[*global.KernelIOStd](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("PID: %d - Operacion: <%s>", estructura.Pid, estructura.Instruction), log.INFO)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.DEBUG)

	global.Logger.Log(fmt.Sprintf("Ingrese un valor (tamaño máximo %d): ", estructura.Length), log.INFO)

	global.Estructura_actualizada.Pid = estructura.Pid
	global.Estructura_actualizada.NumFrames = estructura.NumFrames
	global.Estructura_actualizada.Offset = estructura.Offset

	// fmt.Scanf("%s", global.Texto)
	reader := bufio.NewReader(os.Stdin)
	global.Texto, _ = reader.ReadString('\n')

	global.Logger.Log("De consola escribi: "+global.Texto, log.DEBUG)

	global.VerificacionTamanio(global.Texto, estructura.Length)

	global.Estructura_actualizada.Length = len(global.Texto) - 1

	global.Logger.Log(fmt.Sprintf("Estructura actualizada para mandar a memoria: %+v", global.Estructura_actualizada), log.DEBUG)

	// PUT a memoria de la estructura
	_, err = requests.PutHTTPwithBody[global.MemStdIO, interface{}](global.IOConfig.IPMemory, global.IOConfig.PortMemory, "stdin_read", global.Estructura_actualizada)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("NO se pudo enviar a memoria la estructura %s", err.Error()), log.ERROR)
		w.WriteHeader(http.StatusBadRequest)
		return
		// TODO: falta que memoria vea si puede escribir o no (?)
	}

	dispositivo.InUse = false

	w.WriteHeader(http.StatusNoContent)
}

func Stdout_write(w http.ResponseWriter, r *http.Request) {
	dispositivo := global.Dispositivo
	dispositivo.InUse = true
	var estructura global.KernelIOStd
	var estructura_actualizada global.MemStdIO
	err := serialization.DecodeHTTPBody[*global.KernelIOStd](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s", estructura.Pid, estructura.Instruction+">"), log.INFO)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.DEBUG)

	estructura_actualizada.Pid = estructura.Pid
	estructura_actualizada.NumFrames = estructura.NumFrames
	estructura_actualizada.Offset = estructura.Offset
	estructura_actualizada.Length = estructura.Length

	global.Logger.Log(fmt.Sprintf("Intentando leer con %s", estructura.Name), log.DEBUG)

	time.Sleep(time.Duration(global.IOConfig.UnitWorkTime) * time.Millisecond)

	// PUT a memoria (le paso un registro y me devuelve el valor)

	resp, err := requests.PutHTTPwithBody[global.MemStdIO, string](global.IOConfig.IPMemory, global.IOConfig.PortMemory, "stdout_write", estructura_actualizada)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("NO se pudo enviar a memoria el valor a escribir %s", err.Error()), log.ERROR)
		http.Error(w, "Error al enviar a memoria el valor a escribir", http.StatusBadRequest)
		return
		// TODO: memoria falta que entienda el mensaje (hacer el endpoint) y me devuelva el valor del registro
	}
	global.Logger.Log(fmt.Sprintf("Memoria devolvió este valor: %s", *resp), log.DEBUG)

	dispositivo.InUse = false

	w.WriteHeader(http.StatusNoContent)
}

func Fs_create(w http.ResponseWriter, r *http.Request) {
	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	var estructura global.KernelIOFS_CD

	err := serialization.DecodeHTTPBody[*global.KernelIOFS_CD](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s", estructura.Pid, estructura.Instruction+">"), log.INFO)
	global.Logger.Log(fmt.Sprintf("Estructura: %+v", estructura), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("Dispositivo: %+v", dispositivo), log.DEBUG)

	// implementación

	// 1) busco en mi bitmap el primer bloque libre, uso ese dato para asignarlo como initial_block del archivo metadata estructura.Filename

	firstFreeBlock := getFirstFreeBlock()

	// 2) creo el archivo metadata, de nombre estructura.Filename, con size = 0 e initial_block = al valor hallado en 2)
	filename := global.IOConfig.DialFSPath + "/" + global.Dispositivo.Name + "/" + estructura.FileName
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		global.Logger.Log(fmt.Sprint("Error al crear el archivo:", err), log.ERROR)
		return
	}

	defer file.Close()

	global.UpdateSize(estructura.FileName, 0, w)
	global.UpdateInitialBlock(estructura.FileName, firstFreeBlock, w)

	// 3) actualizo el bitmap, tanto el slice bytes como el archivo (podría hacerlo en el paso 1))

	global.Logger.Log(fmt.Sprintf("Bitmap del FS %s antes de crear el nuevo archivo: %+v", global.Dispositivo.Name, global.Bitmap), log.DEBUG)
	global.UpdateBitmap(1, firstFreeBlock, 1, w)
	global.Logger.Log(fmt.Sprintf("Bitmap del FS %s luego de crear el nuevo archivo: %+v", global.Dispositivo.Name, global.Bitmap), log.DEBUG)

	var filestruct global.File

	filestruct.CurrentBlocks = 0
	filestruct.Initial_block = -1
	filestruct.Size = -1
	global.Logger.Log(fmt.Sprintf("Datos del archivo antes de ser creado (%s): %+v ", filename, filestruct), log.DEBUG)
	filestruct.CurrentBlocks = 1
	filestruct.Initial_block = firstFreeBlock
	filestruct.Size = 0
	global.Logger.Log(fmt.Sprintf("Datos del archivo luego de ser creado (%s): %+v ", filename, filestruct), log.DEBUG)
	global.FilesMap[estructura.FileName] = filestruct

	dispositivo.InUse = false
	w.WriteHeader(http.StatusNoContent)
}

func Fs_delete(w http.ResponseWriter, r *http.Request) { // actualizar bitmap, eliminar archivo del directorio y eliminar elemento del map
	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	var estructura global.KernelIOFS_CD

	err := serialization.DecodeHTTPBody[*global.KernelIOFS_CD](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s", estructura.Pid, estructura.Instruction+">"), log.INFO)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.DEBUG)

	// implementación

	filestruct := global.FilesMap[estructura.FileName]

	// actualizar el bitmap!

	global.Logger.Log(fmt.Sprintf("Bitmap antes de eliminar archivo: %+v", global.Bitmap), log.DEBUG)
	global.UpdateBitmap(0, filestruct.Initial_block, filestruct.CurrentBlocks, w)
	global.Logger.Log(fmt.Sprintf("Bitmap luego de eliminar archivo: %+v", global.Bitmap), log.DEBUG)

	//actualizar la cerpeta de archivos
	metadatapath := global.IOConfig.DialFSPath + "/" + global.Dispositivo.Name + "/" + estructura.FileName

	// Eliminar el archivo metadata
	err = os.Remove(metadatapath)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al eliminar el archivo: "+err.Error()), log.ERROR)
		http.Error(w, "Error al eliminar el archivo", http.StatusInternalServerError)
	}
	dispositivo.InUse = false
	w.WriteHeader(http.StatusNoContent)
}

func Fs_truncate(w http.ResponseWriter, r *http.Request) {

	// decodeo el json que me acaba de llegar para los logs obligatorios

	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	err := serialization.DecodeHTTPBody[*global.KernelIOFS_Truncate](r, &global.Estructura_truncate)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s", global.Estructura_truncate.Pid, global.Estructura_truncate.Instruction+">"), log.INFO)
	global.Logger.Log(fmt.Sprintf("Dispositivo: %+v", dispositivo), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("Instrucción: %+v", global.Estructura_truncate), log.DEBUG)

	// obtengo los datos del archivo metadata

	filestruct := global.FilesMap[global.Estructura_truncate.FileName]

	metadatapath := global.IOConfig.DialFSPath + "/" + global.Dispositivo.Name + "/" + global.Estructura_truncate.FileName

	metadatafile, err := os.Open(metadatapath)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al abrir el archivo %s: %s ", metadatapath, err.Error()), log.DEBUG)
		http.Error(w, "Error al abrir el archivo", http.StatusBadRequest)
		return
	}

	defer metadatafile.Close()

	global.Logger.Log(fmt.Sprintf("Filestruct recién decodeado: %+v", filestruct), log.DEBUG)

	currentBlocks := global.GetCurrentBlocks(global.Estructura_truncate.FileName, w)
	freeContiguousBlocks := global.GetFreeContiguousBlocks(global.Estructura_truncate.FileName, w)
	neededBlocks := global.GetNeededBlocks(w, global.Estructura_truncate)
	totalFreeBlocks := global.GetTotalFreeBlocks(w)

	if currentBlocks == neededBlocks {
		global.UpdateSize(global.Estructura_truncate.FileName, global.Estructura_truncate.Tamanio, w)
		global.Logger.Log(fmt.Sprintf("No es necesario truncar pero actualicé el size: %+v", global.Estructura_truncate), log.DEBUG)
		w.WriteHeader(http.StatusNoContent)
		dispositivo.InUse = false
		return
	} else if !(totalFreeBlocks >= neededBlocks-currentBlocks) {
		global.Logger.Log(fmt.Sprintf("No es posible agrandar el archivo: %+v", global.Estructura_truncate), log.ERROR)
		w.WriteHeader(http.StatusNoContent)
		dispositivo.InUse = false
		return
	} else if currentBlocks > neededBlocks {
		global.Logger.Log(fmt.Sprintf("Trunco a menos %+v", global.Estructura_truncate), log.DEBUG)
		//global.TruncateLess(global.Estructura_truncate.FileName, w)
		//global.AddToTruncatedFiles(global.Estructura_truncate.FileName)
		global.UpdateSize(global.Estructura_truncate.FileName, global.Estructura_truncate.Tamanio, w)
		global.PrintBitmap(w)
		global.UpdateBitmap(0, filestruct.Initial_block+neededBlocks, currentBlocks-neededBlocks, w)
		global.PrintBitmap(w)
		w.WriteHeader(http.StatusNoContent)
		dispositivo.InUse = false
		return
	} else if neededBlocks-currentBlocks <= freeContiguousBlocks {
		global.Logger.Log(fmt.Sprintf("Trunco a más %+v", global.Estructura_truncate), log.DEBUG)
		//global.TruncateMore(global.Estructura_truncate.FileName, w)
		//global.AddToTruncatedFiles(global.Estructura_truncate.FileName)
		global.UpdateSize(global.Estructura_truncate.FileName, global.Estructura_truncate.Tamanio, w)
		global.PrintBitmap(w)
		global.UpdateBitmap(1, filestruct.Initial_block+currentBlocks, neededBlocks-currentBlocks, w)
		global.PrintBitmap(w)
		w.WriteHeader(http.StatusNoContent)
		dispositivo.InUse = false
		return
	} else {
		global.Logger.Log(fmt.Sprintf("Es necesario compactar: %+v", global.Estructura_truncate), log.DEBUG)

		// compactar huecos libres entre bloques ocupados (1 a la izq)

		compactar(global.Estructura_truncate.FileName, totalFreeBlocks, w)

		w.WriteHeader(http.StatusNoContent)
		dispositivo.InUse = false
		return
	}

}

func Fs_write(w http.ResponseWriter, r *http.Request) {
	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	var estructura global.KernelIOFS_WR
	var estructura_actualizada global.MemStdIO

	err := serialization.DecodeHTTPBody[*global.KernelIOFS_WR](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s", estructura.Pid, estructura.Instruction+">"), log.INFO)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.DEBUG)

	// implementación

	estructura_actualizada.Pid = estructura.Pid
	estructura_actualizada.NumFrames = estructura.NumFrames
	estructura_actualizada.Offset = estructura.Offset

	// hago una request a memoria para obtener un valor

	resp, err := requests.PutHTTPwithBody[global.MemStdIO, string](global.IOConfig.IPMemory, global.IOConfig.PortMemory, "stdout_write", estructura_actualizada)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("NO se pudo enviar a memoria el valor a escribir %s", err.Error()), log.ERROR)
		http.Error(w, "Error al enviar a memoria el valor a escribir", http.StatusBadRequest)
		return
		// TODO: memoria falta que entienda el mensaje (hacer el endpoint) y me devuelva el valor del registro
	}
	global.Logger.Log(fmt.Sprintf("Memoria devolvió este valor: %s", *resp), log.DEBUG)

	// convierto la response en un slice de bytes

	valor := []byte(*resp)

	global.Logger.Log(fmt.Sprintf("Conversión de la respuesta de memoria en un slice de bytes: %v", valor), log.INFO)

	// TODO: chequear que donde escribo pertenece al archivo

	//modifico el archivo de bloques

	global.UpdateBlocksFile(valor, estructura.FileName, estructura.PunteroArchivo, w)

	// abro el archivo bloques
	/*
		bloquesdatfile, err := os.OpenFile(bloquesdatpath, os.O_RDWR, 0644)
		if err != nil {
			global.Logger.Log(fmt.Sprintf("Error al abrir el archivo %s: %s ", bloquesdatpath, err.Error()), log.ERROR)
			http.Error(w, "Error al abrir el archivo", http.StatusBadRequest)
			return
		}

		// esta línea de código garantiza que el archivo en el que estoy trabajando se cierre cuando la función actual termina de ejecutarse
		defer bloquesdatfile.Close()



		// TODO: chequear que esté bien colocada la ubicación deseada
		// ubico el puntero en la ubicación deseada

			ubicacionDeseada := global.IOConfig.DialFSBlockSize*global.Filestruct.Initial_block + estructura.PunteroArchivo

			for i := 0; i < len(valor); i++ {

				// Mueve el cursor a medida que vas escribiendo(lenght de valor)
				_, err = bloquesdatfile.Seek(int64(ubicacionDeseada+i), 0)
				if err != nil {
					global.Logger.Log(fmt.Sprintf("Error al mover el cursor: %s ", err.Error()), log.ERROR)
					return
				}

				// escribo el contenido que me llegó de memoria en el archivo de bloques

				_, err = bloquesdatfile.Write(valor[:i])
				if err != nil {
					global.Logger.Log(fmt.Sprintf("Error al escribir en el archivo %s: %s ", bloquesdatpath, err.Error()), log.ERROR)
					http.Error(w, "Error al escribir en el archivo", http.StatusInternalServerError)
					return
				}

			}
	*/
	global.Logger.Log("Datos escritos exitosamente en el archivo bloques.dat", log.INFO)

	dispositivo.InUse = false
	w.WriteHeader(http.StatusNoContent)
}

func Fs_read(w http.ResponseWriter, r *http.Request) {
	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	var estructura global.KernelIOFS_WR

	err := serialization.DecodeHTTPBody[*global.KernelIOFS_WR](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s", estructura.Pid, estructura.Instruction+">"), log.INFO)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.DEBUG)

	// implementación

	filestruct := global.FilesMap[estructura.FileName]

	// abro el archivo metadata y decodeo su contenido

	filepath := global.IOConfig.DialFSPath + "/" + estructura.FileName

	file, err := os.Open(filepath)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al abrir el archivo %s: %s ", filepath, err.Error()), log.DEBUG)
		http.Error(w, "Error al abrir el archivo", http.StatusBadRequest)
		return
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&filestruct)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al decodear el archivo %s: %s ", filepath, err.Error()), log.ERROR)
		http.Error(w, "Error al decodear el archivo", http.StatusBadRequest)
		return
	}

	// abro el archivo bloques

	bloquesdatpath := global.IOConfig.DialFSPath + "/bloques.dat"

	bloquesdatfile, err := os.OpenFile(bloquesdatpath, os.O_RDONLY, 0644)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al abrir el archivo %s: %s ", bloquesdatpath, err.Error()), log.ERROR)
		http.Error(w, "Error al abrir el archivo", http.StatusBadRequest)
		return
	}

	// ubico el puntero en la ubicación deseada

	ubicacionDeseada := global.IOConfig.DialFSBlockSize*filestruct.Initial_block + estructura.PunteroArchivo

	_, err = bloquesdatfile.Seek(int64(ubicacionDeseada), 0)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al mover el cursor: %s ", err.Error()), log.ERROR)
		return
	}

	// crea un slice de estructura.Tamanio bytes

	data := make([]byte, estructura.Tamanio)

	// leo estructura.Tamanio bytes desde el archivo y los asigno a mi data
	_, err = bloquesdatfile.Read(data)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al leer el archivo: %s ", err.Error()), log.ERROR)
		http.Error(w, "Error al leer el archivo", http.StatusBadRequest)
		return
	}

	// la lógica de leer ya está implementada pero por ahora lo hardcodeo
	// porque no tengo forma de escribir en bloques (salvo con hex editor) y está todo en 0

	data[0] = 72
	data[1] = 79
	data[2] = 76
	data[3] = 65
	data[4] = 33

	global.Logger.Log(fmt.Sprintf("Del archivo leí: %+v ", data), log.DEBUG)

	// armo la estructura para mandar a memoria

	global.Estructura_actualizada.Pid = estructura.Pid
	global.Estructura_actualizada.Content = string(data)
	global.Estructura_actualizada.NumFrames = estructura.NumFrames
	global.Estructura_actualizada.Offset = estructura.Offset
	global.Estructura_actualizada.Length = len(data)

	global.Logger.Log(fmt.Sprintf("String a mandar a memoria: %+v", string(data)), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("Estructura a mandar a memoria: %+v", global.Estructura_actualizada), log.DEBUG)

	// Put a memoria de la estructura
	_, err = requests.PutHTTPwithBody[global.MemStdIO, interface{}](global.IOConfig.IPMemory, global.IOConfig.PortMemory, "stdin_read", global.Estructura_actualizada)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("NO se pudo enviar a memoria la estructura %s", err.Error()), log.ERROR)
		w.WriteHeader(http.StatusBadRequest)
		return
		// TODO: falta que memoria vea si puede escribir o no (?)
	}

	dispositivo.InUse = false
	w.WriteHeader(http.StatusNoContent)
}

func compactar(file string, totalfreeblocks int, w http.ResponseWriter) {

	//sacar el truncado

	filestruct := global.FilesMap[file]

	global.UpdateBitmap(0, filestruct.Initial_block, filestruct.CurrentBlocks, w)

	//actualizar bitmap (mover todos los 1 a la izquierda)
	totalUsedBlocks := global.IOConfig.DialFSBlockCount - totalfreeblocks
	global.UpdateBitmap(1, 0, totalUsedBlocks, w)
	global.UpdateBitmap(0, totalUsedBlocks, totalfreeblocks, w)

	// set initial block del truncado
	global.UpdateInitialBlock(global.Estructura_truncate.FileName, getFirstFreeBlock(), w)

	//actualizar los initial block de los archivos de metadata
	updateMetadataFiles(w)

	// actualizar el size del archivo truncado
	global.UpdateSize(global.Estructura_truncate.FileName, global.Estructura_truncate.Tamanio, w)

	global.PrintBitmap(w)

}

func updateMetadataFiles(w http.ResponseWriter) {

	var fileNames []string

	dirPath := global.IOConfig.DialFSPath + "/" + global.Dispositivo.Name

	// Leer los contenidos del directorio
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("No se pudo leer el directorio que contiene los metadata %s", err.Error()), log.ERROR)
	}
	// Iterar sobre los archivos y agregar sus nombres al slice
	for _, entry := range entries {
		if !entry.IsDir() && strings.Contains(entry.Name(), "txt") {
			fileNames = append(fileNames, entry.Name())
		}
	}

	// Imprimir los nombres de los archivos
	global.Logger.Log(fmt.Sprintf("Nombre de los archivos: %+v", fileNames), log.DEBUG)

	currentInitialBlock := 0

	for i := 0; i < len(fileNames); i++ {
		global.Logger.Log(fmt.Sprintf("Archivo a actualizar: %s - %d", fileNames[i], currentInitialBlock), log.DEBUG)
		global.UpdateInitialBlock(fileNames[i], currentInitialBlock, w)
		currentInitialBlock = currentInitialBlock + global.GetCurrentBlocks(fileNames[i], w)
	}

}

func getFirstFreeBlock() int {

	var firstFreeBlock int

	found := false
	for i, v := range global.Bitmap {
		if v == byte(0) {
			global.Logger.Log(fmt.Sprintf("FirstFreeBlock: %d", i), log.DEBUG)
			firstFreeBlock = i
			found = true
			break
		}
	}
	if !found {
		global.Logger.Log("No hay bloques libres", log.DEBUG)
	}
	return firstFreeBlock
}
