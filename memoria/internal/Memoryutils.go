package internal

import (
	
	"os"
	"strings"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/cpu/global"
)

type ProcessPath struct {
	Path string `json:"path"`
	Pid  int    `json:"pid"`
}
type PCB struct {
	Pc  int `json:"pc"`
	Pid int `json:"pid"`
}

func ReadTxt(Path string) ([]string, error) {
	Data, err := os.ReadFile(Path)
	if err != nil {
		global.Logger.Log("error al leer el archivo "+err.Error(), log.ERROR)
		return nil, err
	}
	ListInstructions := strings.Split(string(Data), "\n")

	return ListInstructions, nil
}
