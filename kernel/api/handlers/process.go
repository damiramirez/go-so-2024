package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	global "github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/pcb"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

// Donde crear structs que solo me sirven para enviar body o recibir json???

type Process struct {
	PID   int    `json:"pid"`
	State string `json:"state"`
}

// Handler para devolver a memoria el estado del proceso
func ProcessByIdHandler(w http.ResponseWriter, r *http.Request) {
	pid, _ := strconv.Atoi(r.PathValue("pid"))
	global.Logger.Log(fmt.Sprintf("Buscando PID: %d", pid), log.DEBUG)

	// TODO: Buscar PID en slice de procesos
	pcb := utils.FindProcessInList(pid)

	if pcb == nil {
		global.Logger.Log(fmt.Sprintf("No existe el PID %d", pid), log.DEBUG)
		http.Error(w, fmt.Sprintf("No existe el PID %d", pid), http.StatusNotFound)
		return
	}

	processState := struct {
		PID   int    `json:"pid"`
		State string `json:"state"`
	}{
		PID:   pcb.PID,
		State: pcb.State,
	}

	err := serialization.EncodeHTTPResponse(w, processState, http.StatusOK)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	global.Logger.Log(fmt.Sprintf("Process %d - State: %s", pcb.PID, pcb.State), log.DEBUG)
}

func InitProcessHandler(w http.ResponseWriter, r *http.Request) {

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
	global.Logger.Log("Init process - Path: "+pPath.Path, log.DEBUG)

	pcb := pcb.CreateNewProcess()

	type ProcessMemory struct {
		Path string `json:"path"`
		PID  int    `json:"pid"`
	}
	processMemory:=ProcessMemory{
		Path:pPath.Path, 
		PID: pcb.PID,
	}

	global.Logger.Log(fmt.Sprintf("PCB: %+v", pcb), log.DEBUG)

	// TODO: Agregar a una cola de NEW
	global.NewState.PushBack(pcb)
	global.Logger.Log(fmt.Sprintf("Longitud cola de new: %d", global.NewState.Len()), log.DEBUG)

	// TODO: Request a memoria - enviar instrucciones
	requests.PutHTTPwithBody[ProcessMemory, interface{}](global.KernelConfig.IPMemory, global.KernelConfig.PortMemory, "process", processMemory)
	// _, err = requests.PutHTTPwithBody[ProcessMemory, interface{}](global.KernelConfig.IPMemory, global.KernelConfig.PortMemory, "process", processMemory)
	// if err != nil {
	// 	global.Logger.Log("Error al enviar instruccion "+err.Error(), log.ERROR)
	// 	http.Error(w, "Error al enviar instruccion", http.StatusBadRequest)
	// 	return
	// }

	global.Logger.Log(fmt.Sprintf("Se crea el proceso %d en NEW", pcb.PID), log.INFO)
	processPID := utils.ProcessPID{PID: pcb.PID}

	err = serialization.EncodeHTTPResponse(w, processPID, http.StatusCreated)
	if err != nil {
		http.Error(w, "Error encodeando respuesta", http.StatusInternalServerError)
		return
	}
}


func EndProcessHandler(w http.ResponseWriter, r *http.Request) {
	pid := r.PathValue("pid")
	global.Logger.Log("End process - PID: "+pid, log.DEBUG)

	// TODO: Delete process

	// global.Logger.Log(fmt.Sprintf("Finaliza el proceso %d - Motivo: %s", pcb.PID, pcb.FinalState), log.INFO)
	w.WriteHeader(http.StatusNoContent)
}

func ListProcessHandler(w http.ResponseWriter, r *http.Request) {
	global.Logger.Log("List process", log.DEBUG)

	// TODO: Buscar procesos y listarlos

	processes := []struct {
		PID   int    `json:"pid"`
		State string `json:"state"`
	}{
		{PID: 1, State: "READY"},
		{PID: 2, State: "EXEC"},
	}

	global.Logger.Log(fmt.Sprintf("Longitud cola de new: %d", global.NewState.Len()), log.DEBUG)

	err := serialization.EncodeHTTPResponse(w, processes, http.StatusOK)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
