package utils

import (
	"encoding/json"
	"os"

	utils "github.com/sisoputnfrba/tp-golang/utils/logger"
)

func LoadConfiguration[T any](filePath string, logger utils.Logger) T {
	var config T

	configFile, err := os.Open(filePath)
	if err != nil {
		logger.Log(err.Error(), utils.ERROR)
		os.Exit(1)
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	logger.Log("Config created", utils.INFO)

	return config
}
