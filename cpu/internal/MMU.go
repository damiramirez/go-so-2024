package internal

import (
//handlers "github.com/sisoputnfrba/tp-golang/cpu/api/handlers"
)

func AdressConverter(Page_size int, LogAdress int) (int, int) {
	Page_Number := LogAdress / Page_size
	offset := LogAdress - (Page_Number * Page_size)

	return Page_Number, offset
}
func TLBcreator() {

}
