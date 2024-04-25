package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	global "github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/pcb"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
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
	global.Logger.Log(fmt.Sprintf("State - PID: %d", pid), log.DEBUG)

	// TODO: Buscar PID en slice de procesos
	// Ver que pasa si no existe

	state := "READY"
	processState := struct {
		PID   int    `json:"pid"`
		State string `json:"state"`
	}{
		PID:   pid,
		State: state,
	}

	err := serialization.EncodeHTTPResponse(w, processState, http.StatusOK)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	global.Logger.Log(fmt.Sprintf("Process %d - State: %s", pid, state), log.DEBUG)
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

	// TODO: Request a memoria - enviar instrucciones

	pcb := pcb.CreateNewProcess()
	global.Logger.Log(fmt.Sprintf("PCB: %+v", pcb), log.DEBUG)

	// TODO: Agregar a una cola de NEW
	global.NewState.PushBack(pcb)
	global.Logger.Log(fmt.Sprintf("Longitud cola de new: %d", global.NewState.Len()), log.DEBUG)

	
	processPID := utils.ProcessPID{PID: pcb.PID}

	err = serialization.EncodeHTTPResponse(w, processPID, http.StatusCreated)
	if err != nil {
		http.Error(w, "Error encodeando respuesta", http.StatusInternalServerError)
		return
	}
}

// Se encargará de finalizar un proceso que se encuentre dentro del sistema.
// Este mensaje se encargará de realizar las mismas operaciones como si el proceso
// llegara a EXIT por sus caminos habituales (deberá liberar recursos, archivos y memoria).
func EndProcessHandler(w http.ResponseWriter, r *http.Request) {
	pid := r.PathValue("pid")
	global.Logger.Log("End process - PID: "+pid, log.DEBUG)

	// TODO: Delete process

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
