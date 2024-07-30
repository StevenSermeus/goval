package commands

import (
	"errors"
	"time"

	"github.com/StevenSermeus/goval/src/logging"
	"github.com/StevenSermeus/goval/src/types"
)

func ExecuteCommad(client *types.Client, commandInfo types.CommandInfo) {
	start := time.Now()
	var err error = nil
	switch {
	case IsGet(commandInfo):
		err = ExecuteGetCommand(client, commandInfo)
	case IsSet(commandInfo):
		err = ExecuteSetCommand(client, commandInfo)
	case IsDel(commandInfo):
		err = ExecuteDelCommand(client, commandInfo)
	case IsIncr(commandInfo):
		err = ExecuteIncrCommand(client, commandInfo)
	default:
		err = errors.New("invalid command")
	}
	if err != nil {
		client.Send(types.ResponseInfo{ValueType: "error", Value: err.Error()})
	}
	logging.Info.Println("Command executed in ", time.Since(start))
}
