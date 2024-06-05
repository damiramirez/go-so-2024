package internal

import (
	"fmt"
	
	"os"
	"strings"

	"github.com/sisoputnfrba/tp-golang/memoria/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

var NumPages int

// Se inicializa cada página de la memoria con datos vacíos

func InstructionStorage(data []string, pid int) {
	//creo tabla de paginas (struct con array de paginas) y las inicio en -1
	pagetable:=global.NewPageTable()
	//le asigno al map la lista de instrucciones y la tabla de paginas del proceso q pase por id
	global.DictProcess[pid] = global.ListInstructions{Instructions: data, PageTable: pagetable}
	
	global.Logger.Log(fmt.Sprintf("contenido pagetable %+v",pagetable),log.DEBUG)
}

//

func ReadTxt(Path string) ([]string, error) {
	Data, err := os.ReadFile(Path)
	if err != nil {
		global.Logger.Log("error al leer el archivo "+err.Error(), log.ERROR)
		return nil, err
	}
	ListInstructions := strings.Split(string(Data), "\n")

	return ListInstructions, nil
}

//escribir en memoria
func MemOut(PageNumber int,Offset int,content int,Pid int){
	Byte:=byte(content)
	if Offset>=16 {
		global.Logger.Log("memoria inaccesible",log.ERROR)
		return
	}
	global.Logger.Log("El offset esta bien", log.DEBUG)
	//verifico si esta creada la pagina o si esta llena y necesito crear otra
	PageCheck(PageNumber,Pid,Offset)
	
	FrameBase:=global.DictProcess[Pid].PageTable.Pages[PageNumber-1]
	if  FrameBase==-1{
		return
	}
	MemFrame:=FrameBase+Offset
	global.Memory.Spaces[MemFrame]=Byte
}

func MemIn(PageNumber int,Offset int,Pid int)  byte{
	var Content byte
	if IsWritten(Pid,PageNumber,Offset){
		FrameBase:=global.DictProcess[Pid].PageTable.Pages[PageNumber-1]
		MemFrame:=FrameBase+Offset
		Content=global.Memory.Spaces[MemFrame]
		return Content
	}else {
		return 0 
	}
	
}

func PageCheck(PageNumber int,Pid int,Offset int)bool{
	//si la pagina es valida y que el ultimo valor es igual a 0, osea hay tengo espacio
	if CheckIfValid(PageNumber,Pid)&& global.Memory.Spaces[global.DictProcess[Pid].PageTable.Pages[PageNumber-1]+16]==0{
		return true
	}
	global.Logger.Log("La pagina esta bien", log.DEBUG)
	
	if  len(global.DictProcess[Pid].PageTable.Pages)==0 {
		//agrego pagina
		global.Logger.Log("entre al if",log.DEBUG)
		if PageNumber ==1{
			AddPage(PageNumber,Pid)
			global.Logger.Log("estoy dentro de la addpage",log.DEBUG)
			return true
		}
	}else if global.Memory.Spaces[global.DictProcess[Pid].PageTable.Pages[PageNumber-1]+16]!=0{
		global.Logger.Log("estoy dentro de la addpage del else",log.DEBUG)
		AddPage(PageNumber,Pid)
		return true
	}
	return false
}

func CheckIfValid(PageNumber int, Pid int) bool {
	if process, ok := global.DictProcess[Pid]; ok && process.PageTable != nil {
		if len(process.PageTable.Pages) > 0 {
			for pageNum := range process.PageTable.Pages {
				if PageNumber == pageNum {
					return true
				}
			}
		}
	}
	return false
}

func AddPage(PageNumber int,Pid int){
	for i := 0; i < len(global.BitMap); i++ {
		if global.BitMap[i]==0 {

			if i==0{
				//asigno a la a la tabla de paginas el valor de i y pongo el bit map ocupado en la pos i
				//para i = 0
				global.DictProcess[Pid].PageTable.Pages=append(global.DictProcess[Pid].PageTable.Pages, i)
				global.BitMap[i]=1
				break
				
			} else{
				//asigno a la a la tabla de paginas el valor de i y pongo el bit map ocupado en la pos i
				global.DictProcess[Pid].PageTable.Pages=append(global.DictProcess[Pid].PageTable.Pages, i*16)
				global.BitMap[i]=1
				break
			}
			
		}

	}

}

func IsWritten(Pid int,PageNumber int,Offset int)	bool {
	if  len(global.DictProcess[Pid].PageTable.Pages)==0|| Offset > 16 {
		return false
	}else {
		return true
	}
}
/*func stringsToBytes(strings []string) []byte {
    var bytesSlice []byte
    for _, str := range strings {
        // Convertir cada string a un slice de bytes y concatenarlo al slice resultante
        bytesSlice = append(bytesSlice, []byte(str)...)
    }
    return bytesSlice
}*/
