package handlers

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

func NewIO(w http.ResponseWriter, r *http.Request){

	var Device utils.NewDevice
	err := serialization.DecodeHTTPBody(r, &Device)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	global.Logger.Log("se conecto un nuevo  i/o a kernel  ", log.DEBUG)
	w.WriteHeader(http.StatusNoContent)
}
