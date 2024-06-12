package internal

import (
	"fmt"
	"math"

	"github.com/sisoputnfrba/tp-golang/cpu/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
)

// handlers "github.com/sisoputnfrba/tp-golang/cpu/api/handlers"

type MemStruct struct{
	Pid int `json:"pid"`
	Content any `json:"content"`
	Length int `json:"length"`
	NumFrames []int `json:"numframe"`
	Offset  int `json:"offset"`
}

type GetFrame struct{
	Pid int`json:"pid"`
	Page_Number int	`json:"page_number"`
}

func AdressConverter(LogAdress int) (int, int) {
	Page_Size:=global.CPUConfig.Page_size

	Page_Number := LogAdress / Page_Size
	offset := LogAdress - (Page_Number * Page_Size)
	
	return Page_Number, offset
}

func CreateAdress(Register string,LogAdress int,Pid int,Content any)MemStruct{
	Page_Number,Offset:=AdressConverter(LogAdress)

	global.Logger.Log(fmt.Sprintf("Numero de pagina %d - Offset: %d", Page_Number, Offset), log.DEBUG)

	Length:=getLength(Register)

	global.Logger.Log(fmt.Sprintf("Registro %s - Longitud %d", Register, Length), log.DEBUG)
	
	Adresses :=MemStruct{Pid: Pid,Content: Content,Length: Length,Offset: Offset}

	Page:=GetFrame{Pid: Pid,Page_Number: Page_Number}

	NumPages:=math.Ceil(float64(Offset+Length)/float64(global.CPUConfig.Page_size))

	global.Logger.Log(fmt.Sprintf("Paginas necesarias %d", int(NumPages)), log.DEBUG)

	for i := 0; i < int(NumPages); i++ {
		global.Logger.Log(fmt.Sprintf("Busco a memoria %+v", Page), log.DEBUG)
		frame,_:=requests.PutHTTPwithBody[GetFrame,int](global.CPUConfig.IPMemory,global.CPUConfig.PortMemory,"framenumber",Page)
		global.Logger.Log(fmt.Sprintf("PID: %d - OBTENER MARCO - Página: %d - Marco: %d",Pid,Page.Page_Number,*frame),log.INFO)
		Page.Page_Number=+1
		Adresses.NumFrames = append(Adresses.NumFrames, *frame)
	}

	return Adresses
	//Adresses.Adress[0].Length=(Offset+Lenght)-global.CPUConfig.Page_size
}
func getLength(Register string)int{
	switch Register{
	case "AX", "BX", "CX", "DX":
		return 1
	case "EAX", "EBX", "ECX", "EDX":
		return 4	
	}
	return -1
}

func TLBcreator() {


}
