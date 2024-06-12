package handlers

import (
	"fmt"
	"net/http"

	global "github.com/sisoputnfrba/tp-golang/memoria/global"
	internal "github.com/sisoputnfrba/tp-golang/memoria/internal"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

// recibo tamaño en frames
func Resize(w http.ResponseWriter, r *http.Request) {

	var Process internal.Resize

	err := serialization.DecodeHTTPBody(r, &Process)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	
	//bitMap := global.BitMap
	tablaPag := global.DictProcess[Process.Pid].PageTable.Pages
	//Aumento tamaño
	if len(tablaPag) < Process.Frames {
		global.Logger.Log(fmt.Sprintf("frames a aumentar %d", Process.Frames-len(tablaPag)), log.DEBUG)
		for i := 0; i < Process.Frames-len(tablaPag); i++ {

			if internal.AddPage(Process.Pid) == -1 {
				global.Logger.Log("Error memoria llena", log.DEBUG)
				http.Error(w, "Out of memory", http.StatusForbidden)
				//w.WriteHeader(http.StatusForbidden)
				return
			}
		}	
		global.Logger.Log(fmt.Sprintf("PID: %d - Tamaño Actual: %d- Tamaño a Ampliar: %d",Process.Pid,len(tablaPag),Process.Frames ), log.INFO)

		//global.Logger.Log(fmt.Sprintf("se solicito ampliar la memoria del proceso %d con los siguiente frames %d", Process.Pid, Process.Frames), log.DEBUG)
		// buscar al proceso y ampliar el proceso si la memoria esta llena devolver out of memory
		global.Logger.Log(fmt.Sprintf("page table %d %+v", Process.Pid, global.DictProcess[Process.Pid].PageTable), log.DEBUG)
		global.Logger.Log(fmt.Sprintf("Bit Map  %+v", global.BitMap), log.DEBUG)
		w.WriteHeader(http.StatusNoContent)
		return
		//Reduzco tamaño
	} else if len(global.DictProcess[Process.Pid].PageTable.Pages) > Process.Frames && Process.Frames != 0 {
		difTam := len(global.DictProcess[Process.Pid].PageTable.Pages) - Process.Frames
		//global.Logger.Log(fmt.Sprintf("se solicito reducir la memoria del proceso %d con los siguiente frames %d", Process.Pid, Process.Frames), log.DEBUG)
		global.Logger.Log(fmt.Sprintf("PID: %d - Tamaño Actual: %d- Tamaño a Reducir: %d",Process.Pid,len(global.DictProcess[Process.Pid].PageTable.Pages),Process.Frames ), log.INFO)
		for i := 0; i < difTam; i++ {
			global.BitMap[global.DictProcess[Process.Pid].PageTable.Pages[len(global.DictProcess[Process.Pid].PageTable.Pages)-1-i]] = 0
		}
		//tablaPag=tablaPag[:Process.Frames]
		global.DictProcess[Process.Pid].PageTable.Pages = global.DictProcess[Process.Pid].PageTable.Pages[:Process.Frames]
		//global.Logger.Log(fmt.Sprintf("page table %d %+v", Process.Pid, tablaPag), log.DEBUG)

		global.Logger.Log(fmt.Sprintf("page table %d %+v", Process.Pid, global.DictProcess[Process.Pid].PageTable), log.DEBUG)
		global.Logger.Log(fmt.Sprintf("Bit Map  %+v", global.BitMap), log.DEBUG)

		//Limpio bitmap y tabla de paginas
	} else if Process.Frames == 0 {

		for i := 0; i < len(global.DictProcess[Process.Pid].PageTable.Pages); i++ {
			global.BitMap[global.DictProcess[Process.Pid].PageTable.Pages[len(global.DictProcess[Process.Pid].PageTable.Pages)-1-i]] = 0
		}
		global.DictProcess[Process.Pid].PageTable.Pages = global.DictProcess[Process.Pid].PageTable.Pages[:0]
		//tablaPag=tablaPag[:0]
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
	if frame == -1 {
		global.Logger.Log("Error no existe la pagina", log.DEBUG)
		http.Error(w, "Invalid page", http.StatusForbidden)
	} else {
		//global.Logger.Log(fmt.Sprintf("Se consulto por la pagina %d",frame), log.DEBUG)
		global.Logger.Log(fmt.Sprintf("PID: %d- Pagina: %d - Marco: %d", PageNumber.Pid, PageNumber.PageNumber,frame), log.INFO)

		serialization.EncodeHTTPResponse(w, frame, http.StatusOK)
	}

}

func MemoryAccessIn(w http.ResponseWriter, r *http.Request) {

	var MemoryAccess internal.MemStruct
	err := serialization.DecodeHTTPBody(r, &MemoryAccess)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("MOVIN: Me enviaron: %+v", MemoryAccess), log.DEBUG)

	MemoryAccess.Content = internal.MemIn(MemoryAccess.NumFrames, MemoryAccess.Offset, MemoryAccess.Pid, MemoryAccess.Length)
	serialization.EncodeHTTPResponse(w, MemoryAccess.Content, http.StatusOK)

}

func MemoryAccessOut(w http.ResponseWriter, r *http.Request) {

	global.Logger.Log("Entrando a memoryAcess", log.DEBUG)
	var MemoryAccess internal.MemStruct
	err := serialization.DecodeHTTPBody(r, &MemoryAccess)

	global.Logger.Log(fmt.Sprintf("MOVOUT: Me enviaron: %+v", MemoryAccess), log.DEBUG)

	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}

	if internal.MemOut(MemoryAccess.NumFrames, MemoryAccess.Offset, MemoryAccess.Content, MemoryAccess.Pid, MemoryAccess.Length) {

		global.Logger.Log(fmt.Sprintf("page table %d %+v", MemoryAccess.Pid, global.DictProcess[MemoryAccess.Pid].PageTable), log.DEBUG)
		global.Logger.Log(fmt.Sprintf("Bit Map  %+v", global.BitMap), log.DEBUG)
		global.Logger.Log(fmt.Sprintf("Memoria  %+v", global.Memory), log.DEBUG)
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}
