package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
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
		panic(1)
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
	}
	global.Logger.Log(fmt.Sprintf("Estructura: %+v", estructura), log.DEBUG)

	global.Logger.Log(fmt.Sprintf("Dispositivo: %+v", dispositivo), log.DEBUG)

	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s", estructura.Pid, estructura.Instruction+">"), log.INFO)

	// implementación

	// abrir el archivo bitmap.dat, acceder a la posición dada por initial_block del .txt
	// y cambiar solo ese bit, luego ejecutar FS_TRUNCATE y ocupar la cantidad real de bloques
	// dada por size (del .txt) / global.IOConfig.dialfs_block_size

	// abro el archivo bitmap.dat

	bitmappath := global.IOConfig.DialFSPath + "/Filesystems/" + global.Dispositivo.Name + "/bitmap.dat"

	bitmapfile, err := os.OpenFile(bitmappath, os.O_RDWR, 0644)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al abrir el archivo: %s ", err.Error()), log.ERROR)
	}

	defer bitmapfile.Close() // esta línea de código garantiza que el archivo en el que estoy trabajando se cierre cuando la función actual termina de ejecutarse

	// leo el archivo y logeo su contenido

	data := make([]byte, global.IOConfig.DialFSBlockCount)
	_, err = bitmapfile.Read(data)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al leer el archivo: %s ", err.Error()), log.ERROR)
	}
	global.Logger.Log(fmt.Sprintf("Bitmap del FS %s antes del cambio: %+v", global.Dispositivo.Name, data), log.DEBUG)

	// obtengo la posición (el bit) a cambiar, del archivo metadata (lo abro y decodeo su contenido)

	filepath := "/home/utnso/tp-2024-1c-sudoers/notas.txt"

	file, err := os.Open(filepath)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al abrir el archivo %s: %s ", filepath, err.Error()), log.DEBUG)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&global.Filestruct)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al decodear el archivo: %s: %s ", filepath, err.Error()), log.ERROR)
	}
	global.Logger.Log(fmt.Sprintf("Datos del archivo %s: %+v ", filepath, global.Filestruct), log.DEBUG)

	position := int64(global.Filestruct.Initial_block)

	// muevo el cursor a la posición deseada

	_, err = bitmapfile.Seek(position, 0)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al mover el cursor: %s ", err.Error()), log.ERROR)
		return
	}

	// cambio el bit de 0 a 1 (ver qué pasa si esa posición ya está ocupada, fragmentación externa, compactación)
	_, err = bitmapfile.Write([]byte{1})
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al escribir el byte: %s ", err.Error()), log.ERROR)
		return
	}

	global.Filestruct.CurrentBlocks = 1

	global.Logger.Log(fmt.Sprintf("Current blocks actualizado: %+v ", global.Filestruct), log.DEBUG)

	// muevo el cursor nuevamente al principio del archivo bitmap.dat
	_, err = bitmapfile.Seek(0, 0)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al mover el cursor: %s ", err.Error()), log.ERROR)
		return
	}

	// leo el archivo (desde la posición inicial) y logeo su contenido actualizado

	data = make([]byte, global.IOConfig.DialFSBlockCount)
	_, err = bitmapfile.Read(data)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al leer el archivo: %s ", err.Error()), log.ERROR)
	}

	global.Logger.Log(fmt.Sprintf("Bitmap del FS %s luego del cambio: %+v", global.Dispositivo.Name, data), log.DEBUG)

	dispositivo.InUse = false
	w.WriteHeader(http.StatusNoContent)
}

func Fs_delete(w http.ResponseWriter, r *http.Request) {
	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	var estructura global.KernelIOFS_CD

	err := serialization.DecodeHTTPBody[*global.KernelIOFS_CD](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
	}
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s", estructura.Pid, estructura.Instruction+">"), log.INFO)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.DEBUG)

	// implementación

	// abrir el archivo bloques.dat, acceder a la posición dada por initial_block (del .txt) * config.dialfs_block_size
	// y borrar/setear en 0(? el contenido que hay desde esa posición hasta la posición dada por size (del .txt) * config.dialfs_block_size

	// actualizar el bitmap!

	dispositivo.InUse = false
	w.WriteHeader(http.StatusNoContent)
}

func Fs_truncate(w http.ResponseWriter, r *http.Request) {
	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	var estructura global.KernelIOFS_Truncate

	err := serialization.DecodeHTTPBody[*global.KernelIOFS_Truncate](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
	}
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s", estructura.Pid, estructura.Instruction+">"), log.INFO)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.DEBUG)

	// implementación

	// chequear si hay lugar
	// chequear si hay que compactar (al compactar, actualizo los initial_block de los .txt?)

	// abro el archivo metadata y decodeo su contenido

	filepath := "/home/utnso/tp-2024-1c-sudoers/notas.txt"

	file, err := os.Open(filepath)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al abrir el archivo %s: %s ", filepath, err.Error()), log.DEBUG)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&global.Filestruct)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al decodear el archivo: %s: %s ", filepath, err.Error()), log.ERROR)
	}
	global.Logger.Log(fmt.Sprintf("Datos del archivo antes de truncar %s: %+v ", filepath, global.Filestruct), log.DEBUG)

	global.Logger.Log(fmt.Sprintf("Current blocks luego de decodear: %+v ", global.Filestruct), log.DEBUG)

	// modifico el bitmap usando los datos recién decodeados

	// abro el archivo bitmap.dat

	bitmappath := global.IOConfig.DialFSPath + "/Filesystems/" + global.Dispositivo.Name + "/bitmap.dat"

	bitmapfile, err := os.OpenFile(bitmappath, os.O_RDWR, 0644)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al abrir el archivo: %s ", err.Error()), log.ERROR)
	}

	defer bitmapfile.Close() // esta línea de código garantiza que el archivo en el que estoy trabajando se cierre cuando la función actual termina de ejecutarse

	// leo el archivo y logeo su contenido

	data := make([]byte, global.IOConfig.DialFSBlockCount)
	_, err = bitmapfile.Read(data)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al leer el archivo: %s ", err.Error()), log.ERROR)
	}
	global.Logger.Log(fmt.Sprintf("Bitmap del FS %s antes del cambio: %+v", global.Dispositivo.Name, data), log.DEBUG)

	// cambio todos los bits que corresponden a los bloques que va a ocupar el archivo metadata

	currentBlocks := global.Filestruct.CurrentBlocks

	neededBlocks := int(math.Ceil(float64(estructura.Tamanio) / float64(global.IOConfig.DialFSBlockSize)))

	global.Logger.Log(fmt.Sprintf("Bloques actuales: %d", currentBlocks), log.DEBUG)

	global.Logger.Log(fmt.Sprintf("Bloques necesarios: %d", neededBlocks), log.DEBUG)

	if neededBlocks >= currentBlocks { // tengo que ocupar más bloques

		global.Logger.Log(fmt.Sprintln("Entré a truncar para aumentar"), log.DEBUG)

		// muevo el cursor a la primera posición (global.Filestruct.Initial_block)

		for i := 0; i < neededBlocks; i++ {

			_, err = bitmapfile.Seek(int64(global.Filestruct.Initial_block+i), 0)
			if err != nil {
				global.Logger.Log(fmt.Sprintf("Error al mover el cursor: %s ", err.Error()), log.ERROR)
				return
			}

			// cambio el bit de 0 a 1 (ver qué pasa si esa posición ya está ocupada, fragmentación externa, compactación)
			_, err = bitmapfile.Write([]byte{1})
			if err != nil {
				global.Logger.Log(fmt.Sprintf("Error al escribir el byte: %s ", err.Error()), log.ERROR)
				return
			}
		}

		global.Filestruct.CurrentBlocks = neededBlocks

	} else { // neededBlocks < currentBlocks (tengo que liberar bloques)

		global.Logger.Log(fmt.Sprintln("Entré a truncar para disminuir"), log.DEBUG)

		for i := 0; i < currentBlocks-neededBlocks; i++ {

			_, err = bitmapfile.Seek(int64(neededBlocks+i), 0)
			if err != nil {
				global.Logger.Log(fmt.Sprintf("Error al mover el cursor: %s ", err.Error()), log.ERROR)
				return
			}

			// cambio el bit de 1 a 0
			_, err = bitmapfile.Write([]byte{0})
			if err != nil {
				global.Logger.Log(fmt.Sprintf("Error al escribir el byte: %s ", err.Error()), log.ERROR)
				return
			}
		}

		global.Filestruct.CurrentBlocks = neededBlocks

	}

	// modificar el size en el txt

	file, err = os.OpenFile(filepath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al abrir el archivo %s: %s ", filepath, err.Error()), log.ERROR)
		return
	}

	newSize := map[string]interface{}{
		"initial_block": global.Filestruct.Initial_block,
		"size":          estructura.Tamanio,
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(newSize)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al encodear el nuevo size en el archivo %s: %s ", filepath, err.Error()), log.ERROR)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al mover el cursor: %s ", err.Error()), log.ERROR)
		return
	}

	decoder = json.NewDecoder(file)
	err = decoder.Decode(&global.Filestruct)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al decodear el archivo: %s: %s ", filepath, err.Error()), log.ERROR)
	}

	global.Logger.Log(fmt.Sprintf("Datos del archivo luego de truncar %s: %+v ", filepath, global.Filestruct), log.DEBUG)

	// muevo el cursor nuevamente al principio del archivo bitmap.dat
	_, err = bitmapfile.Seek(0, 0)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al mover el cursor: %s ", err.Error()), log.ERROR)
		return
	}

	// leo el archivo (desde la posición inicial) y logeo su contenido actualizado

	data = make([]byte, global.IOConfig.DialFSBlockCount)
	_, err = bitmapfile.Read(data)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al leer el archivo: %s ", err.Error()), log.ERROR)
	}

	global.Logger.Log(fmt.Sprintf("Bitmap del FS %s luego del cambio: %+v", global.Dispositivo.Name, data), log.DEBUG)

	dispositivo.InUse = false
	w.WriteHeader(http.StatusNoContent)
}

func Fs_write(w http.ResponseWriter, r *http.Request) {
	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	var estructura global.KernelIOFS_WR

	err := serialization.DecodeHTTPBody[*global.KernelIOFS_WR](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
	}
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s", estructura.Pid, estructura.Instruction+">"), log.INFO)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.DEBUG)

	// implementación

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
	}
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s", estructura.Pid, estructura.Instruction+">"), log.INFO)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.DEBUG)

	// implementación

	dispositivo.InUse = false
	w.WriteHeader(http.StatusNoContent)
}
