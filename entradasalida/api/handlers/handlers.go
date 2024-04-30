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

func Ping(w http.ResponseWriter, r *http.Request) {
	global.Logger.Log("me hicieron un request de ping", log.INFO)
	message := "Tu ping es infinito 777\n"
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
	global.Logger.Log(fmt.Sprintf("Me llegó ésta instrucción: %+v", estructura), log.INFO)
	
	global.Logger.Log(fmt.Sprintf("%+v", global.MapIOGenericsActivos[estructura.Nombre]), log.INFO)

	dispositivo := global.MapIOGenericsActivos[estructura.Nombre]

	// if es un IO existente y Generic -> ejecuto el sleep, sino -> "IO intexistente / este tipo de IO no duerme" [este chequeo lo hace Kernel o Entradasalida?]

	global.Logger.Log(fmt.Sprintf("a punto de dormir: %+v", global.MapIOGenericsActivos[estructura.Nombre]), log.INFO)

	dispositivo.EstaEnUso = true // fix to do: no se ve reflejado en el log

	global.Logger.Log(fmt.Sprintf("durmiendo: %+v", global.MapIOGenericsActivos[estructura.Nombre]), log.INFO)

	time.Sleep(time.Duration(estructura.Tiempo * global.IOConfig.UnitWorkTime)*time.Millisecond)

	global.Logger.Log(fmt.Sprintf("terminé de dormir: %+v", global.MapIOGenericsActivos[estructura.Nombre]), log.INFO)

}
