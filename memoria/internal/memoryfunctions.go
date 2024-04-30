package internal

import ("github.com/sisoputnfrba/tp-golang/memoria/global")

//log "github.com/sisoputnfrba/tp-golang/utils/logger"

var NumPages int

// Se inicializa cada página de la memoria con datos vacíos

func InstructionStorage(data []string, pid int) {
    global.DictProcess[pid]=global.ListInstructions{Instructions: data}
}

/*func stringsToBytes(strings []string) []byte {
    var bytesSlice []byte
    for _, str := range strings {
        // Convertir cada string a un slice de bytes y concatenarlo al slice resultante
        bytesSlice = append(bytesSlice, []byte(str)...)
    }
    return bytesSlice
}*/