package internal

import (
	"github.com/sisoputnfrba/tp-golang/memoria/global"
	//log "github.com/sisoputnfrba/tp-golang/utils/logger"
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

/*func stringsToBytes(strings []string) []byte {
    var bytesSlice []byte
    for _, str := range strings {
        // Convertir cada string a un slice de bytes y concatenarlo al slice resultante
        bytesSlice = append(bytesSlice, []byte(str)...)
    }
    return bytesSlice
}*/
