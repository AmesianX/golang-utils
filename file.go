package util

import (
	"os"
)

// IsExistFile : 파일 존재 유무 확인
func IsExistFile(name string) bool {
	if _, e := os.Stat(name); os.IsNotExist(e) {
		return false
	}
	return true
}
