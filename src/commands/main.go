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
	case IsExpr(commandInfo):
		err = ExecuteExprCommand(client, commandInfo)
	case IsVersion(commandInfo):
		err = ExecuteVersionCommand(client, commandInfo)
	default:
		err = errors.New("invalid command")
	}
	if err != nil {
		client.Send(types.ResponseInfo{ValueType: "error", Value: err.Error()})
	}

	if time.Since(start) > time.Second {
		logging.Warning.Printf("%sCommand took more than 1 second to execute%s", "\033[34m", "\033[0m")
	} else {
		logging.Info.Printf("%s Command :executed in %s (Response might have been faster because cache update is done asynchronously)%s", "\033[34m", time.Since(start), "\033[0m")
	}
}
