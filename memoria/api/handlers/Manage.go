package handlers

import (
	"fmt"
	"net/http"

	global "github.com/sisoputnfrba/tp-golang/memoria/global"
	internal "github.com/sisoputnfrba/tp-golang/memoria/internal"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

//recibo tamaño en frames
func Resize(w http.ResponseWriter, r *http.Request) {
	var Process internal.Resize
	err := serialization.DecodeHTTPBody(r, &Process)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	//si la cantidad de frames que me envian(tamaño del proceso) es mas grande que el tamaño de mi tabla de paginas, amplio por la diferencia
	if len(global.DictProcess[Process.Pid].PageTable.Pages) < Process.Frames{
		//
		tamaño:=len(global.DictProcess[Process.Pid].PageTable.Pages)
		for i := 0; i < Process.Frames-tamaño; i++ {
			internal.AddPage(Process.Pid)
		}
		global.Logger.Log(fmt.Sprintf("se solicito ampliar la memoria del proceso %d con los siguiente frames %d", Process.Pid, Process.Frames), log.DEBUG)
		// buscar al proceso y ampliar el proceso si la memoria esta llena devolver out of memory
		global.Logger.Log(fmt.Sprintf("page table %d %+v", Process.Pid, global.DictProcess[Process.Pid].PageTable), log.DEBUG)
		global.Logger.Log(fmt.Sprintf("Bit Map  %+v", global.BitMap), log.DEBUG)

		
		w.WriteHeader(http.StatusNoContent)

	}else if len(global.DictProcess[Process.Pid].PageTable.Pages) >Process.Frames&&Process.Frames!=0{
		nuevoTam:=len(global.DictProcess[Process.Pid].PageTable.Pages) -Process.Frames
		for i := nuevoTam; i >0; i-- {
			global.BitMap[global.DictProcess[Process.Pid].PageTable.Pages[i]]=0
		}
		global.Logger.Log(fmt.Sprintf("Bit Map  %+v", global.BitMap), log.DEBUG)
		global.Logger.Log(fmt.Sprintf("se solicito reducir la memoria del proceso %d con los siguiente frames %d", Process.Pid, Process.Frames), log.DEBUG)
		
		global.DictProcess[Process.Pid].PageTable.Pages=global.DictProcess[Process.Pid].PageTable.Pages[:Process.Frames]
		global.Logger.Log(fmt.Sprintf("page table %d %+v", Process.Pid, global.DictProcess[Process.Pid].PageTable), log.DEBUG)
		
	}else if Process.Frames ==0{
		for i := len(global.DictProcess[Process.Pid].PageTable.Pages); i >0; i-- {
			global.BitMap[global.DictProcess[Process.Pid].PageTable.Pages[i]]=0
		}
		global.DictProcess[Process.Pid].PageTable.Pages=global.DictProcess[Process.Pid].PageTable.Pages[:0]

		global.Logger.Log("vaciando tabla de paginas", log.DEBUG)
		global.Logger.Log(fmt.Sprintf("Bit Map  %+v", global.BitMap), log.DEBUG)

		global.Logger.Log(fmt.Sprintf("page table %d %+v", Process.Pid, global.DictProcess[Process.Pid].PageTable), log.DEBUG)

	}
	

}

// dice q en marco esta asociado la pagina
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
