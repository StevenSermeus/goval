package commands

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/StevenSermeus/goval/src/file"
	"github.com/StevenSermeus/goval/src/types"
)

func IsIncr(commandInfo types.CommandInfo) bool {
	return commandInfo.Command[:4] == "INCR"
}

func Incr(conn *types.Client, commandInfo types.CommandInfo, key string, value any) error {
	//get current value of key from file
	fileContent, err := file.ReadFile(key, conn.ServerConfig)
	fmt.Println("fileContent", fileContent)
	if err != nil {
		return err
	}
	if fileContent.Exp > 0 && fileContent.Exp < time.Now().UnixMilli() {
		conn.Cache.DeleteKey(key)
		err := os.Remove(path.Join(conn.ServerConfig.DataDir, key))
		if err != nil {
			return err
		}
		return errors.New("key not found")
	}
	if fileContent.ValueType != "int" {
		return errors.New("invalid value type")
	}
	//convert value to int
	intValue, err := strconv.Atoi(fileContent.Value)
	if err != nil {
		fmt.Println(err)
		return err
	}
	incrementValue, err := strconv.Atoi(value.(string))
	if err != nil {
		fmt.Println(err, "2")
		return err
	}
	//increment value
	intValue += incrementValue
	//write new value to file
	err = file.WriteFile(key, fileContent.ValueType, strconv.Itoa(intValue), conn.ServerConfig, fileContent.Exp)
	if err != nil {
		return err
	}
	conn.Send(types.ResponseInfo{ValueType: "int", Value: strconv.Itoa(intValue)})
	conn.Cache.SetKey(key, strconv.Itoa(intValue), "int", fileContent.Exp)
	return nil
}

func ExecuteIncrCommand(conn *types.Client, commandInfo types.CommandInfo) error {
	splitCommand := strings.Split(commandInfo.Command, " ")
	if commandInfo.ValueType != "int" {
		return errors.New("invalid command incr take only int")
	}
	if len(splitCommand) < 3 {
		return errors.New("invalid command")
	}
	key := splitCommand[1]
	value := strings.Join(splitCommand[2:], " ")
	return Incr(conn, commandInfo, key, value)
}
