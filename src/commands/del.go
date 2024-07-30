package commands

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/StevenSermeus/goval/src/types"
)

func Del(conn *types.Client, commandInfo types.CommandInfo, key string) error {
	err := os.Remove(path.Join(conn.ServerConfig.DataDir, key))
	if err != nil {
		return err
	}
	conn.Cache.DeleteKey(key)
	conn.Send(types.ResponseInfo{ValueType: "string", Value: "OK"})
	return nil
}

func IsDel(commandInfo types.CommandInfo) bool {
	return commandInfo.Command[:3] == "DEL"
}

func ExecuteDelCommand(conn *types.Client, commandInfo types.CommandInfo) error {
	if !IsDel(commandInfo) {
		return errors.New("invalid command")
	}
	splitCommand := strings.Split(commandInfo.Command, " ")
	if len(splitCommand) != 2 {
		return errors.New("invalid command")
	}
	key := splitCommand[1]
	return Del(conn, commandInfo, key)
}
