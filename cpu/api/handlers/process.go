package handlers

import (
	"fmt"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/cpu/global"
	internal "github.com/sisoputnfrba/tp-golang/cpu/internal/dispatch"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

type PCB struct {
	Pid int `json:"pid"`
	Pc  int `json:"pc"`
}

var pcb PCB

func PCBreciever(w http.ResponseWriter, r *http.Request) {
	err := serialization.DecodeHTTPBody(r, &pcb)
	if err != nil {
		http.Error(w, "Error al decodear el PCB", http.StatusBadRequest)
		return
	}
	for {

		instruction, err := requests.PutHTTPwithBody[PCB, string]("127.0.0.1", 8002, "process/1", pcb)
		if err != nil {
			global.Logger.Log(fmt.Sprintf("Failed to send PC: %v", err), log.ERROR)
			return
		}
		if *instruction != "out of memory" {
			global.Logger.Log(fmt.Sprintf("intruccion nro %d: %s", pcb.Pc+1, *instruction), log.DEBUG)
			pcb.Pc++
		} else {
			return
		}
	}
}

func Dispatch(w http.ResponseWriter, r *http.Request) {
	pcb := &model.PCB{}
	err := serialization.DecodeHTTPBody(r, pcb)
	if err != nil {
		http.Error(w, "Error al decodear PCB", http.StatusBadRequest)
		global.Logger.Log(fmt.Sprintf("Error al decodear PCB: %v", err), log.ERROR)
		return
	}

	pcb, _ = internal.Dispatch(pcb)

	serialization.EncodeHTTPResponse(w, pcb, http.StatusAccepted)
}
