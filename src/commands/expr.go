package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/StevenSermeus/goval/src/file"
	"github.com/StevenSermeus/goval/src/types"
)

func IsExpr(commandInfo types.CommandInfo) bool {
	return commandInfo.Command[:4] == "EXPR"
}

func ExecuteExprCommand(conn *types.Client, commandInfo types.CommandInfo) error {
	if !IsExpr(commandInfo) {
		return errors.New("invalid command")
	}
	splitCommand := strings.Split(commandInfo.Command, " ")
	if len(splitCommand) != 3 {
		fmt.Println("invalid command 1")
		return errors.New("invalid command")
	}
	key := splitCommand[1]
	in, err := strconv.Atoi(splitCommand[2])
	if err != nil {
		fmt.Println(err, "2")
		return err
	}
	return Expr(conn, commandInfo, key, int64(in))
}

func Expr(conn *types.Client, commandInfo types.CommandInfo, key string, in int64) error {
	fileContent, err := file.ReadFile(key, conn.ServerConfig)
	if err != nil {
		return err
	}
	expireAt := time.Now().UnixMilli() + in
	fmt.Println("Exprcall", fileContent)
	err = file.WriteFile(key, fileContent.ValueType, fileContent.Value, conn.ServerConfig, expireAt)
	fmt.Println("Exprcall", fileContent)
	if err != nil {
		fmt.Println(err, "err while update file")
		return err
	}
	conn.Send(types.ResponseInfo{ValueType: "string", Value: "OK"})
	conn.Cache.SetKey(key, fileContent.Value, fileContent.ValueType, expireAt)
	return nil
}
