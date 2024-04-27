package internal

import (
	"fmt"

	global "github.com/sisoputnfrba/tp-golang/memoria/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)
func stringsToBytes(strings []string) []byte {
    var bytesSlice []byte
    for _, str := range strings {
        // Convertir cada string a un slice de bytes y concatenarlo al slice resultante
        bytesSlice = append(bytesSlice, []byte(str)...)
    }
    return bytesSlice
}
//escribe en memoria
// falta desarollar la funcion
func WriteinMemory(data []string) {
    Info:=stringsToBytes(data)
    fmt.Printf("Datos en slice : %+v, largo de slice: %d\n",Info,len(Info))
    global.Logger.Log("Se escribio en memoria ", log.DEBUG)
}
