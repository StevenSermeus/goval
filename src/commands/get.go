package commands

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/StevenSermeus/goval/src/file"
	"github.com/StevenSermeus/goval/src/types"
	"github.com/StevenSermeus/goval/src/utils"
)

func Get(conn *types.Client, commandInfo types.CommandInfo, key string) error {
	cacheEntry, err := conn.Cache.ReadKey(key)
	if err == nil {
		fmt.Println("Cache hit", cacheEntry)
		response := types.ResponseInfo{
			ValueType: cacheEntry.ValueType,
			Value:     cacheEntry.Value,
		}
		conn.Send(response)
		return nil
	}
	fileExists := file.FileExists(path.Join(conn.ServerConfig.DataDir, key))
	if fileExists {
		fileContent, err := os.ReadFile(path.Join(conn.ServerConfig.DataDir, key))
		if err != nil {
			return err
		}
		fileContentString := string(fileContent)
		responseType, err := utils.Type(fileContentString)
		if err != nil {
			return err
		}
		response := types.ResponseInfo{
			ValueType: responseType,
			Value:     fileContentString[0:],
		}
		conn.Send(response)
		conn.Cache.SetKey(key, fileContentString[0:], responseType)
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
