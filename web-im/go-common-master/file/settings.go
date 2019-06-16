package file

import (
	"encoding/json"
	"io"
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

//读入文件内容，以[]byte形式返回
func ReadFileArrayByte(filePath string) ([]byte, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}

//读入文件内容，以[]byte形式返回
func ReadFileString(filePath string) (string, error) {
	fileContent, err := ReadFileArrayByte(filePath)
	if err != nil {
		return "", err
	}
	return string(fileContent), nil
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

//写入文件[]byte内容，添加模式
func AppendFileArrayByte(filePath string, content []byte) error {
	return AppendFile(filePath, content)
}

//写入文件string内容，添加模式
func AppendFileString(filePath string, content string) error {
	return AppendFileArrayByte(filePath, []byte(content))
}

//不存在则创建目录
func MkdirAllNX(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		err = os.Chmod(path, os.ModePerm)
	}
	return nil
}

func AppendFile(filename string, data []byte) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
