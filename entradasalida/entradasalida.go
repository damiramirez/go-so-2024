package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/entradasalida/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
)

type IODevice struct {
	Name  int    `json:"name"`
	State string `json:"state"`
}

func main() {
	
	// Me crea el loger y la configuracion
	global.InitGlobal()

	// Funcion para iniciar una interfaz IO

	func init_IO (string name, config) {





	}



	// Conecto el modulo IO a Kernel como cliente

	for (state == waiting) {

		// Estoy constantemente esperando a que Kernel me mande un mensaje (operación a realizar)
		// Me llega el mensaje, paso a operating 
		state = operating
		
		//Realizo la operacion

		state = waiting
	}
	//processSlice, _ := requests.GetHTTP[ProcessState](global.IOConfig.IPKernel, global.IOConfig.PortKernel, "process/12")
	//global.Logger.Log(fmt.Sprintf("%+v", processSlice), log.DEBUG)

	// requests.DeleteHTTP[interface{}]("plani", global.IOConfig.PortKernel, nil, global.IOConfig.IPKernel)

	global.Logger.CloseLogger()
}
