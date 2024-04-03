package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

type ProcessState struct {
	PID string `json:"pid"`
	State string `json:"state"`
}

func ProcessHandler(w http.ResponseWriter, r *http.Request, logger log.Logger) {
	pid := r.PathValue("pid")

	// TODO: Buscar PID en slice de procesos
	state := "READY"

	processState := ProcessState{PID: pid, State: state}

	w.Header().Set("Content-Type", "application/json")
	
	if err := json.NewEncoder(w).Encode(&processState); err != nil {
		logger.Log(fmt.Sprintf("Failed to encode response: %v", err.Error()), log.ERROR)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Log(fmt.Sprintf("Process %s - State: %s", pid, state), log.DEBUG)
}