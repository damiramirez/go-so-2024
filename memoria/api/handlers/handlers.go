package handlers

import (
	"fmt"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/memoria/global"
	internal "github.com/sisoputnfrba/tp-golang/memoria/internal"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

// LEE DE MEMORIA
func Stdout_write(w http.ResponseWriter, r *http.Request) {
	var MemoryAccessIO internal.MemStdIO
	err := serialization.DecodeHTTPBody(r, &MemoryAccessIO)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
	}
	global.Logger.Log(fmt.Sprintf("Me llegó ésta mensaje: %+v", MemoryAccessIO), log.INFO)

	var Content []byte
	var ContentByte byte

	j := 1
	accu := 0
	for i := 0; i < MemoryAccessIO.Length; i++ {
		if MemoryAccessIO.Offset+i < global.MemoryConfig.PageSize {
			MemFrame := MemoryAccessIO.NumFrames[0]*global.MemoryConfig.PageSize + MemoryAccessIO.Offset + i
			ContentByte = global.Memory.Spaces[MemFrame]
			Content = append(Content, ContentByte)
		} else {
			if accu >= global.MemoryConfig.PageSize {
				accu = 0
				j++
			}
			//newFrame := global.DictProcess[Pid].PageTable.Pages[NumPage+1]
			MemFrame := MemoryAccessIO.NumFrames[j]*global.MemoryConfig.PageSize + accu
			ContentByte = global.Memory.Spaces[MemFrame]
			Content = append(Content, ContentByte)
			accu++
		}
	}
	global.Logger.Log(fmt.Sprintf("page table %d %+v", MemoryAccessIO.Pid, global.DictProcess[MemoryAccessIO.Pid].PageTable), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("Bit Map  %+v", global.BitMap), log.DEBUG)
	//global.Logger.Log(fmt.Sprintf("Memoria  %+v", global.Memory), log.DEBUG)

	str := string(Content)
	serialization.EncodeHTTPResponse(w, str, 200)

}

// ESCRIBE EN MEMORIA
func Stdin_read(w http.ResponseWriter, r *http.Request) {
	//var estructura estructura_write
	var MemoryAccessIO internal.MemStdIO
	err := serialization.DecodeHTTPBody(r, &MemoryAccessIO)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
	}
	global.Logger.Log(fmt.Sprintf("Me llegó este mensaje : %+v", MemoryAccessIO), log.DEBUG)

	byteArray := []byte(MemoryAccessIO.Content)
	global.Logger.Log(fmt.Sprintf("largo %+v", len(byteArray)), log.DEBUG)
	accu := 0
	//cantDeFrames:=len(MemoryAccessIO.NumFrames)
	j := 1
	for i := 0; i < MemoryAccessIO.Length; i++ {

		if i+MemoryAccessIO.Offset < global.MemoryConfig.PageSize {
			MemFrame := MemoryAccessIO.NumFrames[0]*global.MemoryConfig.PageSize + MemoryAccessIO.Offset + i
			global.Memory.Spaces[MemFrame] = byteArray[i]

		} else {
			if accu >= global.MemoryConfig.PageSize {
				accu = 0
				j++
			}
			//newFrame := AddPage(Pid)
			MemFrame := MemoryAccessIO.NumFrames[j]*global.MemoryConfig.PageSize + accu
			global.Memory.Spaces[MemFrame] = byteArray[i]
			accu++
		}
	}

	str := "lo pude escribir"
	global.Logger.Log(fmt.Sprintf("page table %d %+v", MemoryAccessIO.Pid, global.DictProcess[MemoryAccessIO.Pid].PageTable), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("Bit Map  %+v", global.BitMap), log.DEBUG)
	internal.PrintMemoryTable(global.Memory.Spaces, global.MemoryConfig.PageSize)
	//global.Logger.Log(fmt.Sprintf("Memoria  %+v", global.Memory), log.DEBUG)

	serialization.EncodeHTTPResponse(w, str, 200)
	if err != nil {
		http.Error(w, "Error encodeando respuesta", http.StatusInternalServerError)
		return
	}

}
