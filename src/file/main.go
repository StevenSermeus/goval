package file

import (
	"encoding/binary"
	"fmt"
	"os"
	"path"

	"github.com/StevenSermeus/goval/src/config"
	"github.com/StevenSermeus/goval/src/utils"
)

func FileExists(filename_path string) bool {
	if _, err := os.Stat(filename_path); err == nil || os.IsExist(err) {
		return true
	}
	return false
}

func WriteFile(key string, valueType string, value any, serverConfig *config.Config, expirateAtTimestamp int64) error {
	typeCode, err := utils.GetCodeFromType(valueType)
	if err != nil {
		return err
	}
	fmt.Println("typeCode", typeCode)
	fmt.Println("value", value)
	fmt.Println("expirateAtTimestamp", expirateAtTimestamp)
	fmt.Println("expireAtByte", byte(expirateAtTimestamp))
	byteExpire := make([]byte, 8)

	// Use the binary package to put the number into the byte slice in little-endian order
	binary.LittleEndian.PutUint64(byteExpire, uint64(expirateAtTimestamp))
	towrite := append([]byte{}, []byte(typeCode)[0])
	towrite = append(towrite, byteExpire...)
	towrite = append(towrite, []byte(value.(string))...)
	fmt.Println("towrite", towrite)
	err = os.WriteFile(path.Join(serverConfig.DataDir, key), towrite, 0644)
	if err != nil {
		return err
	}
	return nil
}

type FileContent struct {
	ValueType string
	Value     string
	Exp       int64
}

func ReadFile(key string, serverConfig *config.Config) (FileContent, error) {
	fileContent, err := os.ReadFile(path.Join(serverConfig.DataDir, key))

	if err != nil {
		return FileContent{}, err
	}
	typeCode := string(fileContent[0])

	fmt.Println("typeCode", typeCode)
	valueType, err := utils.Type(typeCode)
	if err != nil {
		return FileContent{}, err
	}
	expireAt := binary.LittleEndian.Uint64(fileContent[1:9])
	value := string(fileContent[9:])
	return FileContent{ValueType: valueType, Value: value, Exp: int64(expireAt)}, nil
}
