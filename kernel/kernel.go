package main

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/kernel/api"
	"github.com/sisoputnfrba/tp-golang/kernel/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

func main() {
	global.InitGlobal()
	defer global.Logger.CloseLogger()

	server := api.NewServer()
	err := server.Start()
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Failed to start kernel API server: %v", err), log.ERROR)
		return
	}
}
