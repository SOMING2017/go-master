package main

import (
	"fmt"

	"./file"
)

func GetFileServerRunState() (bool, error) {
	mapSettings, err := getServerSettings()
	if err != nil {
		fmt.Println("GetFileServerRunState->err is ", err)
		return false, err
	}
	return getServerRunState(mapSettings)
}

func SetFileServerRunState(value bool) error {
	mapSettings, err := getServerSettings()
	if err != nil {
		fmt.Println("SetFileServerRunState->err is ", err)
		return err
	}
	mapSettings["ServerRunState"] = value
	err = setServerSettings(mapSettings)
	if err != nil {
		fmt.Println("SetFileServerRunState->err is ", err)
		return err
	}
	fmt.Println("SetFileServerRunState->ServerRunState is ", value)
	return nil
}

//获取ServerRunState
func getServerRunState(mapSettings map[string]interface{}) (bool, error) {
	ServerRunState := mapSettings["ServerRunState"]
	if ServerRunState == nil {
		return false, fmt.Errorf("no value")
	}
	result, ok := ServerRunState.(bool)
	if ok {
		return result, nil
	} else {
		return false, fmt.Errorf("no bool value")
	}
}

//获取配置文件
func getServerSettings() (map[string]interface{}, error) {
	mapSettings, err := file.ReadFileMap("./settings.cnf")
	if err != nil {
		return nil, err
	}
	return mapSettings, nil
}

func setServerSettings(content map[string]interface{}) error {
	return file.WriteFileMap("./settings.cnf", content)
}
