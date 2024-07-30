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
	logging.Info.Printf("Command :%s %s %sexecuted in %s (Response might have been faster because cache update is done asynchronously)", "\033[34m", commandInfo.Command, "\033[0m", time.Since(start))
}
