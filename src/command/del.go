package command

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/StevenSermeus/goval/src/cache"
	"github.com/StevenSermeus/goval/src/config"
)

func Delete(key string, cache *cache.Cache, serverConfig *config.Config) error {
	filePath := filepath.Join(serverConfig.DataDir, key)
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	cache.DeleteKey(key)
	return nil
}

func IsDelCommand(command string) bool {
	return command[:3] == "DEL"
}

func ExecuteDelCommand(command string, cache *cache.Cache, serverConfig *config.Config) error {
	commandParts := strings.Split(command, " ")
	if len(commandParts) != 2 {
		return errors.New("invalid DEL command")
	}
	key := commandParts[1]
	return Delete(key, cache, serverConfig)
}
