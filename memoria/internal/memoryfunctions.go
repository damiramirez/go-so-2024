package internal

import (
	"encoding/binary"
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
	pagetable := global.NewPageTable()
	//le asigno al map la lista de instrucciones y la tabla de paginas del proceso q pase por id
	global.DictProcess[pid] = global.ListInstructions{Instructions: data, PageTable: pagetable}

	global.Logger.Log(fmt.Sprintf("contenido pagetable %+v", pagetable), log.DEBUG)
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

// se le envia un contenido y una direccion para escribir en memoria
func MemOut(NumFrames []int, Offset int, content int, Pid int, Largo int) bool {
	var Slicebytes []byte
	accu:=0
	global.Logger.Log(fmt.Sprintf("largo %d",Largo), log.DEBUG)
	
	//MemFrame := NumFrames[0]*global.MemoryConfig.PageSize + Offset
	if Offset >= global.MemoryConfig.PageSize {
		global.Logger.Log("memoria inaccesible", log.ERROR)
		return false
	}
	global.Logger.Log("El offset esta bien", log.DEBUG)
	if Largo == 4 {
		Slicebytes = EncodeContent(uint32(content))
		global.Logger.Log(fmt.Sprintf("largo %+v",Slicebytes), log.DEBUG)
		for i := 0; i < Largo; i++ {
			if i+Offset < global.MemoryConfig.PageSize {
				MemFrame := NumFrames[0]*global.MemoryConfig.PageSize + Offset + i
				global.Memory.Spaces[MemFrame] = Slicebytes[i]
				
			} else {
				//newFrame := AddPage(Pid)
				MemFrame := NumFrames[1]*global.MemoryConfig.PageSize + accu
				global.Memory.Spaces[MemFrame] = Slicebytes[i]
				accu++
			}
		}
	} else if Largo == 1 {
		global.Memory.Spaces[NumFrames[0]*global.MemoryConfig.PageSize+Offset] = byte(content)
	}
	//accu :=0

	return true

}

// le paso un valor y me devuelve un slice de bytes en hexa
func EncodeContent(value uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, value)
	return bytes
}

func DecodeContent(slice []byte) uint32 {
	return binary.LittleEndian.Uint32(slice)
}

func MemIn(NumFrame []int, Offset int, Pid int, Largo int) int {
	var Content []byte
	var ContentByte byte

	if Largo == 4 {
		accu := 0
		for i := 0; i < 4; i++ {
			if Offset+i < global.MemoryConfig.PageSize {
				MemFrame := NumFrame[0]*global.MemoryConfig.PageSize + Offset + i
				ContentByte = global.Memory.Spaces[MemFrame]
				Content = append(Content, ContentByte)
			} else {
				//newFrame := global.DictProcess[Pid].PageTable.Pages[NumPage+1]
				MemFrame := NumFrame[1]*global.MemoryConfig.PageSize + accu
				ContentByte = global.Memory.Spaces[MemFrame]
				Content = append(Content, ContentByte)
				accu++
			}
		}
		return int(DecodeContent(Content))
	} else {
		MemFrame := NumFrame[0]*global.MemoryConfig.PageSize + Offset
		ContentByte = global.Memory.Spaces[MemFrame]
		return int(ContentByte)
	}

}

func PageCheck(PageNumber int, Pid int, Offset int) bool {

	global.Logger.Log("La pagina esta bien", log.DEBUG)
	global.Logger.Log(fmt.Sprintf(" %+v", global.DictProcess[Pid]), log.DEBUG)
	//si el largo de la pagina es 0

	if checkCompletedPage(PageNumber-1, Pid) {
		global.Logger.Log("estoy dentro de la addpage del else", log.DEBUG)
		AddPage(Pid)
		return true
	}
	return false
}

func checkCompletedPage(PageNumber int, Pid int) bool {

	for i := 0; i < 16; i++ {
		if global.Memory.Spaces[global.DictProcess[Pid].PageTable.Pages[PageNumber]+i] == 0 {
			return false
		}
	}
	return true
}

func GetFrame(PageNumber int, Pid int) int {

	//si es valida esta en la tabla de paginas, devuelvo el frame de la pagina pedida
	if CheckIfValid(PageNumber, Pid) {
		return global.DictProcess[Pid].PageTable.Pages[PageNumber]
	}
	//si es invalida
	return -1
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

// devuelve fram asociado a la pagina q se le mando y devuelve el frame
func AddPage(Pid int) int {
	for i := 0; i < len(global.BitMap); i++ {

		//compruebo que el frame este vacio, si lo esta agrego una pagina
		if global.BitMap[i] == 0 {
			//asigno a la a la tabla de paginas el valor de i y pongo el bit map ocupado en la pos i
			//tamTable:=len(global.DictProcess[Pid].PageTable.Pages)
			global.DictProcess[Pid].PageTable.Pages = append(global.DictProcess[Pid].PageTable.Pages, i)
			global.BitMap[i] = 1

			return i
		}

	}
	return -1
}
/*
func IsWritten(Pid int, Offset int) bool {
	if len(global.DictProcess[Pid].PageTable.Pages) == 0 || Offset > global.MemoryConfig.PageSize {
		return false
	} else {
		return true
	}
}*/

/*func stringsToBytes(strings []string) []byte {
    var bytesSlice []byte
    for _, str := range strings {
        // Convertir cada string a un slice de bytes y concatenarlo al slice resultante
        bytesSlice = append(bytesSlice, []byte(str)...)
    }
    return bytesSlice
}*/

func PrintMemoryTable(memory []byte, cols int) {
    // Imprimir encabezados de columna
    fmt.Print("Addr\t")
    for i := 0; i < cols; i++ {
        fmt.Printf("%02X ", i)
    }
    fmt.Println()

    // Imprimir separador de encabezado
    fmt.Print("----\t")
    for i := 0; i < cols; i++ {
        fmt.Print("---")
    }
    fmt.Println()

    // Imprimir contenido de la memoria
    for i := 0; i < len(memory); i += cols {
        // Imprimir dirección (índice de fila)
        fmt.Printf("%04X\t", i)

        // Imprimir los bytes en la fila
        for j := 0; j < cols; j++ {
            if i+j < len(memory) {
                fmt.Printf("%02X ", memory[i+j])
            } else {
                fmt.Print("   ")
            }
        }
        fmt.Println()
    }
}
