package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sisoputnfrba/tp-golang/entradasalida/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

type estructura_sleep struct {
	Nombre      string `json:"nombre"`
	Instruccion string `json:"instruccion"`
	Tiempo      int    `json:"tiempo"`
	Pid         int    `json:"pid"`
}

type estructura_STDIN_read struct {
	Nombre      string `json:"nombre"`
	Instruccion string `json:"instruccion"`
	Direccion   string `json:"direccion"`
	Tamanio     string `json:"tamanio"`
}

type estructura_read struct {
	Texto     string
	Direccion string
	Tamanio   string
}

type estructura_STDOUT_write struct {
	Nombre      string `json:"nombre"`
	Instruccion string `json:"instruccion"`
	Direccion   string `json:"direccion"`
	Tamanio     string `json:"tamanio"`
}

type estructura_write struct {
	Direccion string
	Tamanio   string
}

func Sleep(w http.ResponseWriter, r *http.Request) {
	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	var estructura estructura_sleep
	err := serialization.DecodeHTTPBody[*estructura_sleep](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
	}
	global.Logger.Log(fmt.Sprintf("Me llegó ésta instrucción: %+v", estructura), log.DEBUG)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.INFO)

	global.Logger.Log(fmt.Sprintf("a punto de dormir: %+v", dispositivo), log.DEBUG)

	dispositivo.InUse = true
	global.Logger.Log(fmt.Sprintf("PID: %d - Operacion: %s", estructura.Pid, estructura.Instruccion), log.INFO)
	global.Logger.Log(fmt.Sprintf("durmiendo: %+v", dispositivo), log.INFO)

	time.Sleep(time.Duration(estructura.Tiempo*global.IOConfig.UnitWorkTime) * time.Millisecond)

	global.Logger.Log(fmt.Sprintf("terminé de dormir: %+v", dispositivo), log.INFO)

	w.WriteHeader(http.StatusNoContent)
	dispositivo.InUse = false
}

func Stdin_read(w http.ResponseWriter, r *http.Request) {
	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	var estructura estructura_STDIN_read
	var estructura_actualizada estructura_read
	err := serialization.DecodeHTTPBody[*estructura_STDIN_read](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
	}
	global.Logger.Log(fmt.Sprintf("Me llegó ésta instrucción: %+v", estructura), log.INFO)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.INFO)

	global.Logger.Log(fmt.Sprintf("Ingrese un valor de tamaño (%s", estructura.Tamanio)+"): ", log.INFO)

	fmt.Scanf("%s", &global.Texto)

	estructura_actualizada.Direccion = estructura.Direccion
	estructura_actualizada.Tamanio = estructura.Tamanio
	estructura_actualizada.Texto = global.Texto

	global.Logger.Log(fmt.Sprintf("Estructura actualizada para mandar a memoria: %+v", estructura_actualizada), log.INFO)

	// PUT a memoria de la estructura
	_, err = requests.PutHTTPwithBody[estructura_read, interface{}](global.IOConfig.IPMemory, global.IOConfig.PortMemory, "stdin_read", estructura_actualizada)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("NO se pudo enviar a memoria la estructura %s", err.Error()), log.INFO)
		panic(1)
		// TODO: falta que memoria vea si puede escribir o no (?)
	}

	dispositivo.InUse = false
}

func Stdout_write(w http.ResponseWriter, r *http.Request) {
	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	var estructura estructura_STDOUT_write
	var estructura_actualizada estructura_write
	err := serialization.DecodeHTTPBody[*estructura_STDOUT_write](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
	}
	global.Logger.Log(fmt.Sprintf("Me llegó ésta instrucción: %+v", estructura), log.INFO)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.INFO)

	estructura_actualizada.Direccion = estructura.Direccion
	estructura_actualizada.Tamanio = estructura.Tamanio

	// PUT a memoria (le paso un registro y me devuelve el valor)

	valor, err := requests.PutHTTPwithBody[estructura_write, interface{}](global.IOConfig.IPMemory, global.IOConfig.PortMemory, "stdout_write", estructura_actualizada)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("NO se pudo enviar a memoria el valor a escribir %s", err.Error()), log.INFO)
		panic(1)
		// TODO: memoria falta que entienda el mensaje (hacer el endpoint) y me devuelva el valor del registro
	}
	global.Logger.Log(fmt.Sprintf("Memoria devolvió este valor: %d", valor), log.INFO)

	dispositivo.InUse = false

}
