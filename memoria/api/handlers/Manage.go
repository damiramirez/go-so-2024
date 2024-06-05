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

//hacer 2 funciones distintas para in y para out
func MemoryAccess(w http.ResponseWriter, r *http.Request){
	//content si es out tiene contenido, lo q quiero guardar, si es in no
	var MemoryAccess internal.MemAccess
	err := serialization.DecodeHTTPBody(r, &MemoryAccess)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	if MemoryAccess.Tipo=="In"{
		//devuelve el valor solicitado en ese espacio de memoria 
		MemoryAccess.Content=int(internal.MemIn(MemoryAccess.NumPage,MemoryAccess.Offset,MemoryAccess.Pid))
		serialization.EncodeHTTPResponse(w,MemoryAccess.Content,r.Response.StatusCode)
	}else if MemoryAccess.Tipo=="Out" {
		//escribe en el espacio de memoria solicitado 
		internal.MemOut(MemoryAccess.NumPage,MemoryAccess.Offset,MemoryAccess.Content,MemoryAccess.Pid)
		global.Logger.Log(fmt.Sprintf("page table %d %+v",MemoryAccess.Pid,global.DictProcess[MemoryAccess.Pid].PageTable), log.DEBUG)
		global.Logger.Log(fmt.Sprintf("Bit Map  %+v",global.BitMap), log.DEBUG)
		global.Logger.Log(fmt.Sprintf("Memoria  %+v",global.Memory), log.DEBUG)
		w.WriteHeader(http.StatusNoContent)
	}else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}