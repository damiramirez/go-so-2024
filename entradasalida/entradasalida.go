package main

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/entradasalida/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
)

type ProcessState struct {
	PID   int    `json:"pid"`
	State string `json:"state"`
}

func main() {

	global.InitGlobal()

	processSlice, _ := requests.GetHTTP[ProcessState](global.IOConfig.IPKernel, global.IOConfig.PortKernel, "process/12")
	global.Logger.Log(fmt.Sprintf("%+v", processSlice), log.INFO)

	requests.DeleteHTTP[interface{}]("plani", global.IOConfig.PortKernel, nil, global.IOConfig.IPKernel)

	global.Logger.CloseLogger()
}
