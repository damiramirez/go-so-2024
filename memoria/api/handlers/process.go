package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	global "github.com/sisoputnfrba/tp-golang/memoria/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

type ProcessPath struct {
	Path string `json:"path"`
}
type Index struct {
	Index int  `json:"pc"`
}

// recibe el codigo q manda kernel y lo guarda en memoria
func CodeReciever(w http.ResponseWriter, r *http.Request) {

	var pPath ProcessPath
	err := serialization.DecodeHTTPBody(r, &pPath)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	//escribe en memoria
	_, err = ReadTxt(pPath.Path)
	if err != nil {
		global.Logger.Log("error al leer el archivo "+err.Error(), log.ERROR)
		http.Error(w, "Error al leer archivo", http.StatusBadRequest)
		return
	}
	//internal.WriteinMemory(pPath.Path)
	w.WriteHeader(http.StatusOK)
}

func ReadTxt(Path string) ([]string, error) {
	Data, err := os.ReadFile(Path)
	if err != nil {
		global.Logger.Log("error al leer el archivo "+err.Error(), log.ERROR)
		return nil, err
	}
	ListInstructions := strings.Split(string(Data), "\n")

	return ListInstructions, nil
}

func SendInstruction(w http.ResponseWriter, r *http.Request) {
	var PC Index
	err := serialization.DecodeHTTPBody(r, &(PC.Index))
	Instruction:= PC.Index
	if err != nil {
		http.Error(w, "Error al decodear el PC", http.StatusBadRequest)
		return
	}
	//en este paso deberia tomar la info de memoria y no del archivo 
	FilePath := fmt.Sprintf("/home/utnso/tp-2024-1c-sudoers/proceso%s.txt", r.PathValue("pid"))
	ListInstructions, err := ReadTxt(FilePath)
	if err != nil {
		http.Error(w, "Error al leer el archivo", http.StatusBadRequest)
		return
	}
	//de aca en adelante la logica es la misma 
	if Instruction>len(ListInstructions){//esto chequea si la intruccion esta dentro del rango 
		global.Logger.Log("out of memory: ", log.ERROR)
		http.Error(w, "out of memory", http.StatusForbidden)
		return
	}
	if Instruction==len(ListInstructions){//esto chequea que no lea memoria q no le corresponde
		mensaje := "out of memory"
		json.NewEncoder(w).Encode(mensaje)
		return 
	}
	DelayResponse := time.Duration(global.MemoryConfig.DelayResponse)
	time.Sleep(DelayResponse*time.Second)//genera el delay response de la respuesta
	err = serialization.EncodeHTTPResponse(w, ListInstructions[Instruction], http.StatusOK)
	if err != nil {
		global.Logger.Log("Error al encodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al encodear el body", http.StatusBadRequest)
		return
	}
}
