package commands

import (
	"errors"

	"github.com/StevenSermeus/goval/src/types"
)

func IsVersion(commandInfo types.CommandInfo) bool {
	return commandInfo.Command[:7] == "VERSION"
}

func Version(conn *types.Client, commandInfo types.CommandInfo) error {
	conn.Send(types.ResponseInfo{ValueType: "string", Value: conn.ServerConfig.Version})
	return nil
}

func ExecuteVersionCommand(conn *types.Client, commandInfo types.CommandInfo) error {
	if !IsVersion(commandInfo) {
		return errors.New("invalid command")
	}
	return Version(conn, commandInfo)
}
