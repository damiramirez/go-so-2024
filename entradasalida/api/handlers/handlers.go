package handlers

import (
	"bufio"
	"fmt"
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
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s", estructura.Pid, estructura.Instruction+">"), log.INFO)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.DEBUG)

	// implementación

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
