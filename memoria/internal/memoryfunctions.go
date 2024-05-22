package internal

import (
	"os"
	"strings"
    log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/memoria/global"
)



var NumPages int

var Memory *MemoryST

// Se inicializa cada página de la memoria con datos vacíos

func NewMemory() *MemoryST {

	ByteArray := make([]byte, global.MemoryConfig.MemorySize)
	mem := MemoryST{spaces: ByteArray}

	return &mem
}
func NewPageTable() *PageTable {
	ByteArray := make([]byte, global.MemoryConfig.PageSize)
	pagetable := PageTable{pages: ByteArray}

	return &pagetable
}
func InstructionStorage(data []string, pid int) {
	global.DictProcess[pid] = global.ListInstructions{Instructions: data}
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
