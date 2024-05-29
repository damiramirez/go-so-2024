package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/memoria/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

type estructura_read struct {
	Texto     string
	Direccion string
	Tamanio   string
}

func Stdin_read(w http.ResponseWriter, r *http.Request) {
	var estructura estructura_read
	err := serialization.DecodeHTTPBody[*estructura_read](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
	}
	global.Logger.Log(fmt.Sprintf("Me llegó ésta instrucción: %+v", estructura), log.INFO)

	// escribe en memoria

	w.WriteHeader(http.StatusNoContent)
}

func Stdout_write(w http.ResponseWriter, r *http.Request) {
	var estructura estructura_read
	err := serialization.DecodeHTTPBody[*estructura_read](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
	}
	global.Logger.Log(fmt.Sprintf("Me llegó ésta instrucción: %+v", estructura), log.INFO)

	// busca en memoria y devuelve un valor

	valor := "10A"

	respuesta, err := json.Marshal(valor)
	if err != nil {
		global.Logger.Log("Error al convertir la respuesta a JSON: "+err.Error(), log.ERROR)
		return
	}

	w.Write(respuesta)

}
