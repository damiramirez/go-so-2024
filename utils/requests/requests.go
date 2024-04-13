package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

func GetHTTP[T any](ip string, port int, endpoint string, logger *log.Logger) (*T, error) {
	url := fmt.Sprintf("http://%s:%d/%s", ip, port, endpoint)
	logger.Log("GET to: "+url, log.DEBUG)

	resp, err := http.Get(url)
	if err != nil {
		logger.Log("Error al hacer el request: "+err.Error(), log.ERROR)
		return nil, err
	}
	defer resp.Body.Close()

	var data T
	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		logger.Log("Error al decodear el body de la respuesta: "+err.Error(), log.ERROR)
		return nil, err
	}
	return &data, nil
}

func PutHTTPwithBody[T any, R any](ip string, port int, endpoint string, data T, logger *log.Logger) (*R, error) {
	var RespData *R
	url := fmt.Sprintf("http://%s:%d/%s", ip, port, endpoint)

	logger.Log("POST to: "+url, log.DEBUG)

	body, err := json.Marshal(data)
	if err != nil {
		logger.Log("Error al encodear la estructura: "+err.Error(), log.ERROR)
		return RespData, err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		logger.Log("error al hacer el request: "+err.Error(), log.ERROR)
		return RespData, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log("Error al enviar la solicitud: "+err.Error(), log.ERROR)
		return RespData, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent{
		return nil, nil
	}

	err = json.NewDecoder(resp.Body).Decode(&RespData)
	if err != nil {
		logger.Log(fmt.Sprintf("Error al decodificar la respuesta: %s", err), log.ERROR)
		return RespData, err
	}

	return RespData, nil
}

//desarolle una funcion q permite hacer los deletes con plani y process con PID aplicable a la funcion de arriba
//para q funcione deberiamos pasarles como parametro en donde dice endpointwithpid en el caso de un proccess/pid
// en el caso de plani/  
func DeleteHTTP[T any](endpointwithPID string,logger log.Logger,port int,data T,ip string) (*T,error){
	var RespData T
	Delimitador := "/"
	SplitString := strings.Split(endpointwithPID, Delimitador)
	endpoint:=SplitString[0]
	var url string
	Pid:=SplitString[1]
	if len(Pid)==0 {
		url = fmt.Sprintf("http://%s:%d/%s", ip, port, endpoint)
	} 
	if len(Pid)!=0{
		url = fmt.Sprintf("http://%s:%d/%s/%s", ip, port, endpoint,Pid)
	}
	body, err := json.Marshal(data)

	if err != nil {
		logger.Log("Error al encodear la estructura: "+err.Error(), log.ERROR)
		return &RespData, err
	}
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(body))

	if err != nil {
		logger.Log("error al hacer el request: "+err.Error(), log.ERROR)
		return &RespData, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log("Error al enviar la solicitud: "+err.Error(), log.ERROR)
		return &RespData, err
	}
	defer resp.Body.Close()
	_ = json.NewDecoder(resp.Body).Decode(&RespData)
	//no devuelve nada por ahora por lo que despues se deberia modificar
	/*if err != nil {
		logger.Log(fmt.Sprintf("Error al decodificar la respuesta: %s", err), log.ERROR)
		return &RespData, err
	}*/ //comento para despues implementarlo en el caso de que no se pueda decodificar la respuesta 

	return &RespData,nil
}


