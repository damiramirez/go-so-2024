package handlers

import (
	"fmt"
	"net/http"

	global "github.com/sisoputnfrba/tp-golang/memoria/global"
	internal "github.com/sisoputnfrba/tp-golang/memoria/internal"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)
 
func Resize(w http.ResponseWriter, r *http.Request){
	var Process internal.Resize
	err := serialization.DecodeHTTPBody(r, &Process)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	if Process.Tipo=="enlarge"{
		global.Logger.Log(fmt.Sprintf("se solicito ampliar la memoria del proceso %d con los siguiente frames %d",Process.Pid,Process.Frames),log.DEBUG)
		// buscar al proceso y ampliar el proceso si la memoria esta llena devolver out of memory
	}
	if Process.Tipo=="reduce"{
		global.Logger.Log(fmt.Sprintf("se solicito reducir la memoria del proceso %d con los siguiente frames %d",Process.Pid,Process.Frames),log.DEBUG)
	}
	
}
func PageTableAccess(w http.ResponseWriter, r *http.Request){
	var PageNumber internal.Page
	err := serialization.DecodeHTTPBody(r, &PageNumber)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	//recibe una pagina y envia numero de frame correspondiente
	global.Logger.Log(fmt.Sprintf("buscando frame asociado a pagina %d",PageNumber.PageNumber),log.DEBUG)
}

func MemoryAccess(w http.ResponseWriter, r *http.Request){
	var MemoryAcess internal.MemAccess
	err := serialization.DecodeHTTPBody(r, &MemoryAcess)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	if MemoryAcess.Tipo=="Read"{
		//devuelve el valor solicitado en ese espacio de memoria 
		global.Logger.Log(fmt.Sprintf("se quiere leer  memoria con esta direccion %d",MemoryAcess.Adress), log.DEBUG)
	}
	if MemoryAcess.Tipo=="Write" {
		//escribe en el espacio de memoria solicitado 
		global.Logger.Log(fmt.Sprintf("se quiere escribir memoria con esta direccion %d",MemoryAcess.Adress), log.DEBUG)
	}
	
}