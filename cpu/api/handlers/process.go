package handlers

import (
	"fmt"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/cpu/global"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

type PCB struct {
	Pid string `json:"pid"`
	Pc  int    `json:"pc"`
}
var pcb PCB
func PCBreciever(w http.ResponseWriter, r *http.Request) {
	err := serialization.DecodeHTTPBody(r, &pcb)
	if err != nil {
		http.Error(w, "Error al decodear el PCB", http.StatusBadRequest)
		return
	}
	for {
	    instruction, err:= requests.PutHTTPwithBody[int,string]("127.0.0.1",8002,"process/1",pcb.Pc)
		if err!=nil{
	    	global.Logger.Log(fmt.Sprintf("Failed to send PC: %v", err), log.ERROR)
			return
	    }
		if *instruction!="out of memory"{
	    global.Logger.Log(fmt.Sprintf("intruccion nro %d: %s",pcb.Pc+1,*instruction),log.DEBUG)
		pcb.Pc++
		}	else {return}
	}
}
