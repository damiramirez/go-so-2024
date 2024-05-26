package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sisoputnfrba/tp-golang/entradasalida/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
)

type estructura_sleep struct {
	Nombre      string `json:"nombre"`
	Instruccion string `json:"instruccion"`
	Tiempo      int    `json:"tiempo"`
}

type estructura_STDIN_read struct {
	Nombre      string `json:"nombre"`
	Instruccion string `json:"instruccion"`
	Direccion 	string `json:"direccion"`
	Tamanio 	string `json:"tamanio"`
}

type estructura_STDOUT_write struct {
	Nombre      string `json:"nombre"`
	Instruccion string `json:"instruccion"`
	Direccion 	string `json:"direccion"`
	Tamanio 	string `json:"tamanio"`
}

func Sleep(w http.ResponseWriter, r *http.Request) {
	var estructura estructura_sleep
	err := serialization.DecodeHTTPBody[*estructura_sleep](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
	}
	global.Logger.Log(fmt.Sprintf("Me llegó ésta instrucción: %+v", estructura), log.INFO)

	dispositivo := global.Dispositivo

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.INFO)

	global.Logger.Log(fmt.Sprintf("a punto de dormir: %+v", dispositivo), log.INFO)

	dispositivo.InUse = true

	global.Logger.Log(fmt.Sprintf("durmiendo: %+v", dispositivo), log.INFO)

	time.Sleep(time.Duration(estructura.Tiempo*global.IOConfig.UnitWorkTime) * time.Millisecond)

	dispositivo.InUse = false

	global.Logger.Log(fmt.Sprintf("terminé de dormir: %+v", dispositivo), log.INFO)

}

func Stdin_read(w http.ResponseWriter, r *http.Request) {
	var estructura estructura_STDIN_read
	err := serialization.DecodeHTTPBody[*estructura_STDIN_read](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
	}
	global.Logger.Log(fmt.Sprintf("Me llegó ésta instrucción: %+v", estructura), log.INFO)

	dispositivo := global.Dispositivo

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.INFO)

	global.Logger.Log(fmt.Sprintf("Ingrese un valor: "), log.INFO)

	fmt.Scanf("%s", &global.Texto)

	fmt.Println(global.Texto)

	// PUT a memoria de "texto"
	//stdin_read()

}

func stdin_read() {

	_, err := requests.PutHTTPwithBody[string, interface{}](global.IOConfig.IPMemory, global.IOConfig.PortMemory, "stdin_read", global.Texto)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("NO se pudo enviar a memoria el valor a escribir %s", err.Error()), log.INFO)
		panic(1)
		// TODO: memoria falta que entienda el mensaje (hacer el endpoint) y vea si puede escribir o no (?)
	}
}

func Stdout_write(w http.ResponseWriter, r *http.Request) {
	var estructura estructura_STDOUT_write
	err := serialization.DecodeHTTPBody[*estructura_STDOUT_write](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
	}
	global.Logger.Log(fmt.Sprintf("Me llegó ésta instrucción: %+v", estructura), log.INFO)

	dispositivo := global.Dispositivo

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.INFO)

	// PUT a memoria (le paso un registro y me devuelve el valor)
	stdout_write()

}

func stdout_write() {

	valor, err := requests.PutHTTPwithBody[string, interface{}](global.IOConfig.IPMemory, global.IOConfig.PortMemory, "stdout_write", global.Texto)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("NO se pudo enviar a memoria el valor a escribir %s", err.Error()), log.INFO)
		panic(1)
		// TODO: memoria falta que entienda el mensaje (hacer el endpoint) y me devuelva el valor del registro
	}

	fmt.Sprintf("%d", valor)
}