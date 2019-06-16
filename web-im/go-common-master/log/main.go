package log

import (
	"fmt"
	"time"

	"../file"
)

var (
	//"../../"
	rootPath     = "../"
	logName      = "go-master-log"
	apiLogPath   = "/api"
	errorLogPath = "/error"
)

func LogApiInfo(apiPath string, args ...string) {
	go func() {
		filePath := rootPath + logName + apiLogPath + apiPath
		file.MkdirAllNX(filePath)
		now := time.Now()
		fileName := fmt.Sprint(now.Year(), "-", int(now.Month()), "-", now.Day(), ".log")
		file.AppendFileString(filePath+"/"+fileName, formatApiArgs(args))
	}()
}

func LogErrorInfo(classPath string, args ...string) {
	go func() {
		filePath := rootPath + logName + errorLogPath + classPath
		file.MkdirAllNX(filePath)
		now := time.Now()
		fileName := fmt.Sprint(now.Year(), "-", int(now.Month()), "-", now.Day(), ".log")
		file.AppendFileString(filePath+"/"+fileName, formatApiArgs(args))
	}()
}

func formatApiArgs(args []string) string {
	fix := "\r\n//************************************************\r\n"
	log := fix
	log += "logTime:" + time.Now().String() + "\r\n"
	log += "start input args:\r\n"
	for _, arg := range args {
		log += "/**/" + arg + "/**/" + "\r\n"
	}
	log += "end input args."
	log += fix
	return log
}
