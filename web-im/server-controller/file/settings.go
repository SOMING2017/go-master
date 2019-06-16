package file

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//读入文件内容，以map形式返回
func ReadFileMap(filePath string) (map[string]interface{}, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var mapContent map[string]interface{}
	err = json.Unmarshal(fileContent, &mapContent)
	if err != nil {
		return nil, err
	}
	return mapContent, nil
}

//写入文件Map内容
func WriteFileMap(filePath string, content map[string]interface{}) error {
	contentJson, err := json.Marshal(content)
	if err != nil {
		return err
	}
	return WriteFileArrayByte(filePath, contentJson)
}

//写入文件[]byte内容
func WriteFileArrayByte(filePath string, content []byte) error {
	return ioutil.WriteFile(filePath, content, os.ModePerm)
}

//写入文件string内容
func WriteFileString(filePath string, content string) error {
	return WriteFileArrayByte(filePath, []byte(content))
}
