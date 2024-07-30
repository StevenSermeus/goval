package file

import (
	"os"
)

func FileExists(filename_path string) bool {
	if _, err := os.Stat(filename_path); err == nil || os.IsExist(err) {
		return true
	}
	return false
}
