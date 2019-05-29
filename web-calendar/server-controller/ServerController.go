// ServerController
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func GetFileCloseServer() (bool, error) {
	mapSettings, err1 := getServerSettings()
	if err1 != nil {
		fmt.Println("获取设置失败:", err1)
	}
	fileCloseServer := mapSettings["closeServer"]
	if fileCloseServer == nil {
		return false, fmt.Errorf("no value")
	}
	result, ok := fileCloseServer.(bool)
	if ok {
		return result, nil
	} else {
		return false, fmt.Errorf("no bool value")
	}
}

func SetFileCloseServer(setBool bool) {
	mapSettings, err1 := getServerSettings()
	if err1 != nil {
		fmt.Println("获取设置失败:", err1)
	}
	err2 := setCloseServer(mapSettings, setBool)
	if err2 != nil {
		fmt.Println("设置参数失败:", err2)
	}
	mapSettingsJson, err3 := json.Marshal(mapSettings)
	if err3 != nil {
		fmt.Println("转换设置失败:", err3)
	}
	setServerSettings(mapSettingsJson)
	fmt.Println("SetFileCloseServer运行结束...setBool:", setBool)
}

func setCloseServer(mapSettings map[string]interface{}, setBool bool) error {
	fileCloseServer := mapSettings["closeServer"]
	if fileCloseServer == nil {
		return fmt.Errorf("no value")
	}
	_, ok := fileCloseServer.(bool)
	if ok {
		mapSettings["closeServer"] = setBool
		return nil
	} else {
		return fmt.Errorf("no bool value")
	}

}

func getServerSettings() (map[string]interface{}, error) {
	settings, err1 := ioutil.ReadFile("./ServerSettings.cnf")
	if err1 != nil {
		fmt.Println("err1:", err1)
		return nil, err1
	}
	var mapSettings map[string]interface{}
	err2 := json.Unmarshal(settings, &mapSettings)
	if err2 != nil {
		fmt.Println("err2:", err2)
		return nil, err2
	}
	return mapSettings, nil
}

func setServerSettings(content []byte) error {
	return ioutil.WriteFile("./ServerSettings.cnf", content, os.ModePerm)
}
