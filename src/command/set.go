package command

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/StevenSermeus/goval/src/cache"
	"github.com/StevenSermeus/goval/src/config"
)

func Set(key string, value string, cache *cache.Cache, serverConfig *config.Config) error {
	filePath := filepath.Join(serverConfig.DataDir, key)
	err := os.WriteFile(filePath, []byte(value), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}
	cache.SetKey(key, value)
	return nil
}

func IsSetCommand(command string) bool {
	return command[:3] == "SET"
}

func ExecuteSetCommand(command string, cache *cache.Cache, serverConfig *config.Config) error {
	commandParts := strings.Split(command, " ")
	if len(commandParts) < 3 {
		return errors.New("invalid SET command format")
	}
	key := commandParts[1]
	value := strings.Join(commandParts[2:], " ")
	return Set(key, value, cache, serverConfig)
}