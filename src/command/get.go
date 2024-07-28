package command

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/StevenSermeus/goval/src/cache"
	"github.com/StevenSermeus/goval/src/config"
	"github.com/StevenSermeus/goval/src/file"
)

func Get(key string, cache *cache.Cache, serverConfig *config.Config) (string, error) {
	value, err := cache.ReadKey(key)
	if err == nil {
		return value, nil
	}
	filePath := filepath.Join(serverConfig.DataDir, key)
	fileExists := file.FileExists(filePath)
	if !fileExists {
		return "", errors.New("key not found")
	}
	file_bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	cache.SetKey(key, string(file_bytes))
	return string(file_bytes), nil
}

func IsGetCommand(command string) bool {
	return command[:3] == "GET"
}

func ExecuteGetCommand(command string, cache *cache.Cache, serverConfig *config.Config) (string, error) {
	commandParts := strings.Split(command, " ")
	if len(commandParts) != 2 {
		return "", errors.New("invalid GET command")
	}
	key := commandParts[1]
	return Get(key, cache, serverConfig)
}
