package commands

import (
	"errors"
	"os"
	"path"
	"strings"
	"time"

	"github.com/StevenSermeus/goval/src/file"
	"github.com/StevenSermeus/goval/src/types"
)

func Get(conn *types.Client, commandInfo types.CommandInfo, key string) error {
	cacheEntry, err := conn.Cache.ReadKey(key)
	if err == nil {
		response := types.ResponseInfo{
			ValueType: cacheEntry.ValueType,
			Value:     cacheEntry.Value,
		}
		conn.Send(response)
		return nil
	} else {
		if err.Error() == "key expired" {
			go os.Remove(path.Join(conn.ServerConfig.DataDir, key))
			return errors.New("key not found")
		}
	}
	fileExists := file.FileExists(path.Join(conn.ServerConfig.DataDir, key))
	if fileExists {
		content, err := file.ReadFile(key, conn.ServerConfig)
		if err != nil {
			return err
		}
		if content.Exp > 0 && content.Exp < time.Now().UnixMilli() {
			go os.Remove(path.Join(conn.ServerConfig.DataDir, key))
			return errors.New("key not found")
		}
		conn.Send(types.ResponseInfo{ValueType: content.ValueType, Value: content.Value})
		conn.Cache.SetKey(key, content.Value, content.ValueType)
		return nil
	}
	return errors.New("key not found")
}

func IsGet(commandInfo types.CommandInfo) bool {
	return commandInfo.Command[:3] == "GET"
}

func ExecuteGetCommand(conn *types.Client, commandInfo types.CommandInfo) error {
	if !IsGet(commandInfo) {
		return errors.New("invalid command")
	}
	splitCommand := strings.Split(commandInfo.Command, " ")
	if len(splitCommand) != 2 {
		return errors.New("invalid command")
	}
	key := splitCommand[1]
	return Get(conn, commandInfo, key)
}
