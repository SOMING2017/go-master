package face

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

var (
	//"../../../html"
	htmlPath        = "../html"
	facePath        = "/user/assets/face/"
	defaultFacePath = facePath + "default/"
)

//返回默认头像src
//找到默认头像目录，获取目录文件数
//随机某个图片文件作为默认头像
func GetRandomDefaultFaceSrc() (string, error) {
	defaultFaceFiles := make([]string, 0)
	filepath.Walk(htmlPath+defaultFacePath, func(path string, info os.FileInfo, err error) error {
		fileinfo, err := os.Stat(path)
		if err != nil || fileinfo.IsDir() {
			return err
		}
		defaultFaceFiles = append(defaultFaceFiles, fileinfo.Name())
		return nil
	})
	faceFileSize := len(defaultFaceFiles)
	if faceFileSize == 0 {
		return "", fmt.Errorf("not found default face image")
	}
	var randomNum = 0
	if faceFileSize > 1 {
		rand := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomNum = rand.Intn(faceFileSize - 1)
	}
	return defaultFacePath + defaultFaceFiles[randomNum], nil
}
