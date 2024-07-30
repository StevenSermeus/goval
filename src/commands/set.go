package commands

import (
	"errors"
	"strings"

	"github.com/StevenSermeus/goval/src/file"
	"github.com/StevenSermeus/goval/src/types"
)

func IsSet(commandInfo types.CommandInfo) bool {
	return commandInfo.Command[:3] == "SET"
}

func Set(conn *types.Client, commandInfo types.CommandInfo, key string, value any) error {
	file.WriteFile(key, commandInfo.ValueType, value, conn.ServerConfig, 0)
	conn.Send(types.ResponseInfo{ValueType: "string", Value: "OK"})
	conn.Cache.SetKey(key, value, commandInfo.ValueType)
	return nil
}

func ExecuteSetCommand(conn *types.Client, commandInfo types.CommandInfo) error {
	splitCommand := strings.Split(commandInfo.Command, " ")
	if len(splitCommand) < 3 {
		return errors.New("invalid command")
	}
	key := splitCommand[1]
	value := strings.Join(splitCommand[2:], " ")
	return Set(conn, commandInfo, key, value)
}
