package handlers

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"

	global "github.com/sisoputnfrba/tp-golang/memoria/global"
	//"github.com/sisoputnfrba/tp-golang/memoria/internal"
)

// recibe el codigo q manda kernel y lo guarda en memoria
func CodeReciever(w http.ResponseWriter, r *http.Request) {
	type ProcessPath struct {
		Path string `json:"path"`
	}
	var pPath ProcessPath
	err := serialization.DecodeHTTPBody(r, &pPath)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	//escribe en memoria
	ReadTxt(pPath.Path)
	//internal.WriteinMemory(pPath.Path)
    w.WriteHeader(http.StatusOK)
}

func ReadTxt(Path string)([]byte,error){
	Data ,err := os.ReadFile(Path)
	if err!=nil{
		global.Logger.Log("error al leer el archivo "+err.Error(),log.ERROR)
		return nil,err
	}
	global.Logger.Log(fmt.Sprintf("datos del arhcivo %s",string(Data)),log.INFO)

	return Data,nil
}
