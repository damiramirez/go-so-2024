package handlers

import (
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
	"net/http"

	global "github.com/sisoputnfrba/tp-golang/memoria/global"
)

// recibe el codigo q manda kernel y lo guarda en memoria
func CodeReciever(w http.ResponseWriter, r *http.Request, Memory *global.Memory) {
	type Code struct {
		Code []string
	}
	var ProcessCode Code
	var bytesSlice []byte
	err := serialization.DecodeHTTPBody(r, ProcessCode)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	for _, str := range ProcessCode.Code {
		// Convertir cada string a un slice de bytes y concatenarlo al slice resultante
		bytesSlice = append(bytesSlice, []byte(str)...)
	}
	//escribe en memoria
	global.WriteinMemory(bytesSlice, Memory)
    w.WriteHeader(http.StatusOK)
}
