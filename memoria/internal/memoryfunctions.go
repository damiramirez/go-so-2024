package internal

import (
	"os"
	"strings"
    log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/memoria/global"
)



var NumPages int

// Se inicializa cada página de la memoria con datos vacíos

func InstructionStorage(data []string, pid int) {
    global.DictProcess[pid]=global.ListInstructions{Instructions: data}
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
/*func stringsToBytes(strings []string) []byte {
    var bytesSlice []byte
    for _, str := range strings {
        // Convertir cada string a un slice de bytes y concatenarlo al slice resultante
        bytesSlice = append(bytesSlice, []byte(str)...)
    }
    return bytesSlice
}*/