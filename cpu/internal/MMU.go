package internal

import "github.com/sisoputnfrba/tp-golang/cpu/global"

// handlers "github.com/sisoputnfrba/tp-golang/cpu/api/handlers"

func AdressConverter(LogAdress int) (int, int) {
	Page_Size:=global.CPUConfig.Page_size
	Page_Number := LogAdress / Page_Size
	offset := LogAdress - (Page_Number * Page_Size)
	
	return Page_Number, offset
}

func TLBcreator() {


}
