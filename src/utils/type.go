package utils

import (
	"errors"
)

func Type(str string) (string, error) {
	if len(str) == 0 {
		return "", errors.New("empty string")
	}
	switch str[0] {
	case '+':
		return "string", nil
	case '-':
		return "error", nil
	case ':':
		return "int", nil
	case '!':
		return "command", nil
	default:
		return "", errors.New("invalid message format")
	}
}

func GetCodeFromType(typeValue string) (string, error) {
	switch typeValue {
	case "string":
		return "+", nil
	case "error":
		return "-", nil
	case "int":
		return ":", nil
	case "command":
		return "!", nil
	default:
		return "", errors.New("invalid type")
	}
}
