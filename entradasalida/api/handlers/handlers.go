package handlers

import (
	"fmt"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/entradasalida/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

type estructura_sleep struct {
	Nombre      string `json:"name"`
	Instruccion string `json:"instruccion"`
	Tiempo      int    `json:"tiempo"`
}

func Ping(w http.ResponseWriter, r *http.Request) {
	global.Logger.Log("me hicieron un request de ping", log.INFO)
	message := "Hello, World!\n"
	w.Write([]byte(message))
}

func Sleep(w http.ResponseWriter, r *http.Request) {
	var estructura estructura_sleep
	global.Logger.Log(fmt.Sprintf("%+v", estructura), log.INFO)
	err := serialization.DecodeHTTPBody[*estructura_sleep](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
	}
	global.Logger.Log(fmt.Sprintf("%+v", estructura), log.INFO)

}
