package internal

import (
	
	global "github.com/sisoputnfrba/tp-golang/memoria/global"
	log  "github.com/sisoputnfrba/tp-golang/utils/logger"

)

//escribe en memoria
// falta desarollar la funcion
func WriteinMemory(data string) {
    /*pageIndex := address / ConfigPag.PageSize
    offset := address % ConfigPag.PageSize
    if pageIndex < numPages && offset+len(data) <= ConfigPag.PageSize {
        copy(mem.pages[pageIndex].data[offset:], data)
    } else {
        fmt.Println("Error: Dirección de memoria fuera de rango")
    }*/
    global.Logger.Log("Se escribio en memoria ", log.DEBUG)
}