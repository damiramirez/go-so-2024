package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

func GetHTTP[T any](ip string, port int, endpoint string, logger *log.Logger) (*T, error) {
	url := fmt.Sprintf("http://%s:%d/%s", ip, port, endpoint)
	logger.Log("GET to: " + url, log.DEBUG)

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

func PutHTTPwithBody[T any, R any](ip string, port int, endpoint string, data T, logger *log.Logger) (R, error) {
	var RespData R

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
	

	err = json.NewDecoder(resp.Body).Decode(&RespData)
	if err != nil {
			logger.Log(fmt.Sprintf("Error al decodificar la respuesta: %s", err), log.ERROR)
			return RespData, err
	}

	return RespData, nil
}