package requests

import (
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
		logger.Log("Error with the request: "+err.Error(), log.ERROR)
		return nil, err
	}
	defer resp.Body.Close()
	
	var data T
	err = json.NewDecoder(resp.Body).Decode(&data)
	
	if err != nil {
		logger.Log("Error with the decode: "+err.Error(), log.ERROR)
		return nil, err
	}
	return &data, nil
}