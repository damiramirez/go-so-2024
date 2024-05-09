package handlers

import (
	"fmt"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

func NewIO(w http.ResponseWriter, r *http.Request){

	var Device global.NewDevice
	err := serialization.DecodeHTTPBody(r, &Device)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	//global.Logger.Log(fmt.Sprintf("se conecto un nuevo  i/o a kernel %s ",Device.Name), log.DEBUG)
	global.IoMap[Device.Name]=global.NewDevice{Name: Device.Name,Port: Device.Port,Usage: Device.Usage,Type: Device.Type}
	global.Logger.Log(fmt.Sprintf("se conecto un nuevo  i/o a kernel %s ",global.IoMap[Device.Name].Name), log.DEBUG)

	w.WriteHeader(http.StatusNoContent)
}
