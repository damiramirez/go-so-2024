package handler

import (
	"fmt"
	"net/http"

	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

// Donde crear structs que solo me sirven para enviar body o recibir json???

type Process struct {
	PID   int    `json:"pid"`
	State string `json:"state"`
}

// Handler para devolver a memoria el estado del proceso
func ProcessStateHandler(w http.ResponseWriter, r *http.Request, logger log.Logger) {
	pid := r.PathValue("pid")
	logger.Log("State - PID: "+pid, log.DEBUG)

	// TODO: Buscar PID en slice de procesos
	// Ver que pasa si no existe

	state := "READY"
	processState := struct {
		PID   string `json:"pid"`
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

	logger.Log(fmt.Sprintf("Process %s - State: %s", pid, state), log.DEBUG)
}

// Se encargará de ejecutar un nuevo proceso en base a un archivo
// dentro del file system de Linux. Dicho mensaje se encargará de
// la creación del proceso (PCB) y dejará el mismo en el estado NEW.
func InitProcessHandler(w http.ResponseWriter, r *http.Request, logger log.Logger) {
	type ProcessPath struct {
		Path string `json:"path"`
	}

	var pPath ProcessPath

	err := serialization.DecodeHTTPBody[ProcessPath](r, pPath)

	if err != nil {
		logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}

	// TODO: Crear proceso - PCB y dejarlo en NEW
	logger.Log("Init process - Path: "+pPath.Path, log.DEBUG)

	processPID := struct {
		PID int `json:"pid"`
	}{
		PID: 1,
	}

	err = serialization.EncodeHTTPResponse(w, processPID, http.StatusCreated)
	if err != nil {
		logger.Log("Error al encodear la respuesta: "+err.Error(), log.ERROR)
		http.Error(w, "Error encodeando respuesta", http.StatusInternalServerError)
		return
	}
}

// Se encargará de finalizar un proceso que se encuentre dentro del sistema.
// Este mensaje se encargará de realizar las mismas operaciones como si el proceso
// llegara a EXIT por sus caminos habituales (deberá liberar recursos, archivos y memoria).
func EndProcessHandler(w http.ResponseWriter, r *http.Request, logger log.Logger) {
	pid := r.PathValue("pid")
	logger.Log("End process - PID: "+pid, log.DEBUG)

	// TODO: Delete process

	w.WriteHeader(http.StatusNoContent)
}

func ListProcessHandler(w http.ResponseWriter, r *http.Request, logger log.Logger) {
	logger.Log("List process", log.DEBUG)

	// TODO: Buscar procesos y listarlos

	processes := []struct {
		PID   int    `json:"pid"`
		State string `json:"state"`
	}{
		{PID: 1, State: "READY"},
		{PID: 2, State: "EXEC"},
	}

	err := serialization.EncodeHTTPResponse(w, processes, http.StatusOK)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
