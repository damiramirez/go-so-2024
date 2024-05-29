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
	pagetable:=global.NewPageTable()
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


func MemOut(PageNumber int,Offset int,content int,Pid int){
	Byte:=byte(content)
	if Offset>=16 {
		global.Logger.Log("memoria inaccesible",log.ERROR)
		return
	}
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

func PageCheck(PageNumber int,Pid int,Offset int){
	if  global.DictProcess[Pid].PageTable.Pages[PageNumber-1]==-1 {
		AddPage(PageNumber,Pid)
		global.Logger.Log("estoy dentro de la addpage",log.DEBUG)
	}
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
