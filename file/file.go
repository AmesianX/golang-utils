package file

import (
	"io/ioutil"
	"os"
)

// IsExistFile 파일 존재 유무 확인
func IsExistFile(name string) bool {
	if _, e := os.Stat(name); os.IsNotExist(e) {
		return false
	}
	return true
}

// GetFileNameList 해당 경로의 파일 목록 반환
func GetFileNameList(path string) ([]string, error) {
	var fileNameList []string
	files, e := ioutil.ReadDir(path)
	if e != nil {
		return nil, e
	}

	for _, file := range files {
		if file.IsDir() {
			retFileNameList, e := GetFileNameList(path + "/" + file.Name())
			if e != nil {
				return nil, e
			}
			for _, fileName := range retFileNameList {
				fileNameList = append(fileNameList, fileName)
			}
		} else {
			fileNameList = append(fileNameList, path+"/"+file.Name())
		}
	}

	return fileNameList, nil
}
