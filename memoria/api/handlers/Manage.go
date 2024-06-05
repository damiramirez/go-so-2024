package handlers

import (
	"fmt"
	"net/http"

	global "github.com/sisoputnfrba/tp-golang/memoria/global"
	internal "github.com/sisoputnfrba/tp-golang/memoria/internal"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

func Resize(w http.ResponseWriter, r *http.Request) {
	var Process internal.Resize
	err := serialization.DecodeHTTPBody(r, &Process)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	if Process.Tipo == "enlarge" {
		global.Logger.Log(fmt.Sprintf("se solicito ampliar la memoria del proceso %d con los siguiente frames %d", Process.Pid, Process.Frames), log.DEBUG)
		// buscar al proceso y ampliar el proceso si la memoria esta llena devolver out of memory
	}
	if Process.Tipo == "reduce" {
		global.Logger.Log(fmt.Sprintf("se solicito reducir la memoria del proceso %d con los siguiente frames %d", Process.Pid, Process.Frames), log.DEBUG)
	}

}

// dice q en marco de esta asociado la pagina
func PageTableAccess(w http.ResponseWriter, r *http.Request) {
	var PageNumber internal.Page
	err := serialization.DecodeHTTPBody(r, &PageNumber)

	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	frame := internal.GetFrame(PageNumber.PageNumber, PageNumber.Pid)
	serialization.EncodeHTTPResponse(w, frame, http.StatusOK)

}

func MemoryAccessIn(w http.ResponseWriter, r *http.Request) {
	//content si es out tiene contenido, lo q quiero guardar, si es in no
	var MemoryAccess internal.MemAccess
	err := serialization.DecodeHTTPBody(r, &MemoryAccess)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("Me enviaron: %+v", MemoryAccess), log.DEBUG)

	MemoryAccess.Content = internal.MemIn(MemoryAccess.NumFrame,MemoryAccess.NumPage , MemoryAccess.Offset, MemoryAccess.Pid)
	serialization.EncodeHTTPResponse(w, MemoryAccess.Content, http.StatusOK)
	
}

func MemoryAccessOut(w http.ResponseWriter, r *http.Request) {
	
	global.Logger.Log("ENtrando a memoryAcess", log.DEBUG)
	var MemoryAccess internal.MemAccess
	err := serialization.DecodeHTTPBody(r, &MemoryAccess)

	global.Logger.Log(fmt.Sprintf("Me enviaron: %+v", MemoryAccess), log.DEBUG)

	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	if internal.MemOut(MemoryAccess.NumFrame,MemoryAccess.Offset,MemoryAccess.Content,MemoryAccess.Pid){
	
		global.Logger.Log(fmt.Sprintf("page table %d %+v", MemoryAccess.Pid, global.DictProcess[MemoryAccess.Pid].PageTable), log.DEBUG)
		global.Logger.Log(fmt.Sprintf("Bit Map  %+v", global.BitMap), log.DEBUG)
		global.Logger.Log(fmt.Sprintf("Memoria  %+v", global.Memory), log.DEBUG)
		w.WriteHeader(http.StatusNoContent)
	}else{
		w.WriteHeader(http.StatusBadRequest)
	}
	
}
