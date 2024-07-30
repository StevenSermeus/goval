package commands

import (
	"errors"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/StevenSermeus/goval/src/types"
)

func IsIncr(commandInfo types.CommandInfo) bool {
	return commandInfo.Command[:4] == "INCR"
}

func Incr(conn *types.Client, commandInfo types.CommandInfo, key string, value any) error {
	//get current value of key from file
	fileContent, err := os.ReadFile(path.Join(conn.ServerConfig.DataDir, key))
	if err != nil {
		return err
	}
	fileContentString := string(fileContent)
	typeIndicator := fileContentString[0]
	valueFromFile := fileContentString[1:]
	if typeIndicator != ':' {
		return errors.New("invalid type for key")
	}
	//parse value from file
	valueFromFileInt, err := strconv.Atoi(valueFromFile)
	if err != nil {
		return err
	}
	//parse value from command
	valueInt, err := strconv.Atoi(value.(string))
	if err != nil {
		return err
	}
	//add values
	newValue := valueFromFileInt + valueInt
	//write new value to file
	toWrite := []byte(":" + strconv.Itoa(newValue))
	err = os.WriteFile(path.Join(conn.ServerConfig.DataDir, key), toWrite, 0644)
	if err != nil {
		return err
	}
	conn.Send(types.ResponseInfo{ValueType: "int", Value: strconv.Itoa(newValue)})
	conn.Cache.SetKey(key, strconv.Itoa(newValue), "int")
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
