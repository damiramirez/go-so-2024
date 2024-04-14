package main

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/cpu/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	requests "github.com/sisoputnfrba/tp-golang/utils/requests"
)

type ProcessState struct {
	PID   int    `json:"pid"`
	State string `json:"state"`
}

func main() {
	global.InitGlobal()

	processSlice, err := requests.GetHTTP[[]ProcessState](global.CPUConfig.IPKernel, global.CPUConfig.PortKernel, "process")
	if err != nil {
		global.Logger.Log(err.Error(), log.ERROR)
		return
	}

	global.Logger.Log(fmt.Sprintf("%+v", processSlice), log.INFO)

	global.Logger.CloseLogger()
}
