package command

import (
	"errors"

	"github.com/StevenSermeus/goval/src/cache"
	"github.com/StevenSermeus/goval/src/config"
	"github.com/StevenSermeus/goval/src/logging"
)

func Exec(command string, cache *cache.Cache, serverConfig *config.Config) (string, error) {
	logging.Info.Println("Executing command :", command)
	switch {
	case IsGetCommand(command):
		return ExecuteGetCommand(command, cache, serverConfig)
	case IsSetCommand(command):
		return "", ExecuteSetCommand(command, cache, serverConfig)
	case IsDelCommand(command):
		return "", ExecuteDelCommand(command, cache)
	case IsHealthCheckCommand(command):
		return ExecuteHealthCheckCommand(command)
	default:
		return "", errors.New("invalid command")
	}
}
