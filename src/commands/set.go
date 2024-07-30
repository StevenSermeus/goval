package commands

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/StevenSermeus/goval/src/types"
	"github.com/StevenSermeus/goval/src/utils"
)

func IsSet(commandInfo types.CommandInfo) bool {
	return commandInfo.Command[:3] == "SET"
}

func Set(conn *types.Client, commandInfo types.CommandInfo, key string, value any) error {
	typeCode, err := utils.GetCodeFromType(commandInfo.ValueType)
	if err != nil {
		return err
	}
	toWrite := []byte(typeCode + value.(string))
	err = os.WriteFile(path.Join(conn.ServerConfig.DataDir, key), toWrite, 0644)
	if err != nil {
		return err
	}
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
