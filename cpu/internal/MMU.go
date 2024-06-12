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
	Pid int
	Content any
	Length int
	NumFrames []int
	Offset  int
}

type GetFrame struct{
	Pid int
	Page_Number int
}

func AdressConverter(LogAdress int) (int, int) {
	Page_Size:=global.CPUConfig.Page_size

	Page_Number := LogAdress / Page_Size
	offset := LogAdress - (Page_Number * Page_Size)
	
	return Page_Number, offset
}

func CreateAdress(Register string,LogAdress int,Pid int,Content any)MemStruct{
	Page_Number,Offset:=AdressConverter(LogAdress)

	Length:=getLength(Register)
	
	Adresses :=MemStruct{Pid: Pid,Content: Content,Length: Length,Offset: Offset}

	Page:=GetFrame{Pid: Pid,Page_Number: Page_Number}

	NumPages:=math.Ceil(float64(Offset+Length)/float64(global.CPUConfig.Page_size))


	for i := 0; i < int(NumPages); i++ {
		Page.Page_Number=+i
		frame,_:=requests.PutHTTPwithBody[GetFrame,int](global.CPUConfig.IPMemory,global.CPUConfig.PortMemory,"framenumber",Page)
		global.Logger.Log(fmt.Sprintf("PID: %d - OBTENER MARCO - Página: %d - Marco: %d",Pid,Page_Number,*frame),log.INFO)
		Adresses.NumFrames = append(Adresses.NumFrames, *frame)

	}

	return Adresses
	//Adresses.Adress[0].Length=(Offset+Lenght)-global.CPUConfig.Page_size
}
func getLength(Register string)int{
	switch Register{
	case "AX":
	case "BX":
	case "CX":
	case "DX":
		return 1
	case "EAX":
	case "EBX":
	case "ECX":
	case "EDX":
		return 4	
	}

	return -1
}

func TLBcreator() {


}
