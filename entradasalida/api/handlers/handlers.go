package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sisoputnfrba/tp-golang/entradasalida/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

type estructura_sleep struct {
	Nombre      string `json:"name"`
	Instruccion string `json:"instruccion"`
	Tiempo      int    `json:"tiempo"`
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
