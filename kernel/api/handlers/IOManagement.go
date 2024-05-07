package handlers

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

type NewDevice struct {
	Port  int    `json:"port"`
	Usage bool   `json:"usage"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

func NewIO(w http.ResponseWriter, r *http.Request) {

	var Device NewDevice
	err := serialization.DecodeHTTPBody(r, &Device)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body ", http.StatusBadRequest)
		return
	}
	global.Logger.Log("se conecto un nuevo  i/o a kernel  ", log.DEBUG)
	w.WriteHeader(http.StatusNoContent)
}
