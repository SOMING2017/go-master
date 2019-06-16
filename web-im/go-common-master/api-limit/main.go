package limit

import (
	"fmt"
)

var ErrorLimit = fmt.Errorf("请求繁忙，请稍后再试。")

//普通接口限制
//5秒内4次调用限制三秒
func IsNormalApiLimit(apiPath string, sid string, ipv4 string) error {
	return IsApiLimit(apiPath, sid, ipv4, 5, 4, 3)
}

//消息接口限制
//1秒内至多发送一次
func IsMessageApiLimit(apiPath string, sid string, ipv4 string) error {
	return IsApiLimit(apiPath, sid, ipv4, 1, 2, 1)
}

//含义：{apiLimitSeconds}秒内满足{apiLimitMaxNum}次接口调用，则需等待{apiLimitDelaySeconds}秒后方可正常调用
func IsApiLimit(apiPath string, sid string, ipv4 string,
	apiLimitSeconds int, apiLimitMaxNum int, apiLimitDelaySeconds int) error {
	resultBool := false
	sidSum, err := GetSumApiLimit(apiPath, sid, apiLimitSeconds)
	if err != nil {
		return err
	}
	resultBool = resultBool || (sidSum >= apiLimitMaxNum)
	ipv4Sum, err := GetSumApiLimit(apiPath, ipv4, apiLimitSeconds)
	if err != nil {
		return err
	}
	resultBool = resultBool || (ipv4Sum >= apiLimitMaxNum)
	if resultBool {
		go ExpireApiLimit(apiPath, sid, apiLimitDelaySeconds)
		go ExpireApiLimit(apiPath, ipv4, apiLimitDelaySeconds)
		return ErrorLimit
	}
	return nil
}
