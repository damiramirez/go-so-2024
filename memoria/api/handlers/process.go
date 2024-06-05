package handlers

import (
	"fmt"
	"net/http"
	
	"time"

	global "github.com/sisoputnfrba/tp-golang/memoria/global"
	internal "github.com/sisoputnfrba/tp-golang/memoria/internal"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

// recibe el codigo q manda kernel y lo guarda en slice de strings
func CodeReciever(w http.ResponseWriter, r *http.Request) {
	//recibo pid y path del archivo con instrucciones
	var pPath internal.ProcessPath
	err := serialization.DecodeHTTPBody(r, &pPath)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	//leo el archivo con las instrucciones y las guardo en un array de string
	ListInstructions, err := internal.ReadTxt(pPath.Path)
	if err != nil {
		global.Logger.Log("error al leer el archivo "+err.Error(), log.ERROR)
		http.Error(w, "Error al leer archivo", http.StatusBadRequest)
		return
	}
	//escribo en el map el pid y su lista de instrucciones y crea tabla de paginas
	internal.InstructionStorage(ListInstructions, pPath.Pid)
	global.Logger.Log(fmt.Sprintf("%+v\n", global.DictProcess), log.INFO)
	w.WriteHeader(http.StatusOK)
}



func SendInstruction(w http.ResponseWriter, r *http.Request) {
	var ProcessAssets internal.ProcessAssets
	//decodeo lo q me envia cpu (pid y pc)
	err := serialization.DecodeHTTPBody(r, &ProcessAssets)
	if err != nil {
		http.Error(w, "Error al decodear el PC", http.StatusBadRequest)
		return
	}
	//asocio el pc con la instruccion que tengo q enviar
	Instruction := ProcessAssets.Pc
	//me traigo toda la lista de instrucciones del pid correspondiente
	ListInstructions:=global.DictProcess[ProcessAssets.Pid].Instructions
	//de aca en adelante la logica es la misma
	if Instruction > len(ListInstructions) { //esto chequea si la intruccion esta dentro del rango
		global.Logger.Log("out of memory: ", log.ERROR)
		http.Error(w, "out of memory", http.StatusForbidden)
		return
	}
	if Instruction == len(ListInstructions) { //esto chequea que no lea memoria q no le corresponde
		w.WriteHeader(http.StatusNoContent)
		return
	}
	DelayResponse := time.Duration(global.MemoryConfig.DelayResponse)
	time.Sleep(DelayResponse * time.Millisecond) //genera el delay response de la respuesta
	err = serialization.EncodeHTTPResponse(w, ListInstructions[Instruction], http.StatusOK)
	if err != nil {
		global.Logger.Log("Error al encodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al encodear el body", http.StatusBadRequest)
		return
	}
}

func DeleteProcess(w http.ResponseWriter, r *http.Request){
	var ProcessDelete internal.ProcessDelete
	err := serialization.DecodeHTTPBody(r, &ProcessDelete)
	if err != nil {
		http.Error(w, "Error al decodear el PC", http.StatusBadRequest)
		return
	}
	// y elimino las intrucciones de memoria 
	global.Logger.Log(fmt.Sprintf("se elimino el proceso con el PID : %d ", ProcessDelete.Pid), log.DEBUG)
	global.DictProcess[ProcessDelete.Pid]=global.ListInstructions{}
	//elimino frames relacionados al proceso, marco los frames como libres sin eliminarlos 
}
