package limit

import (
	myRedis "../redis"
)

var apiLimitPrefix = "api:limit:"

//暂存api次数
func GetSumApiLimit(apiPath string, onlyString string, life int) (int, error) {
	myRedis.SetNXExpireKey(apiLimitPrefix+apiPath+":"+onlyString, 0, life)
	myRedis.INCRKey(apiLimitPrefix + apiPath + ":" + onlyString)
	return myRedis.GetInt(apiLimitPrefix + apiPath + ":" + onlyString)
}

//延长暂存时间
func ExpireApiLimit(apiPath string, onlyString string, life int) {
	myRedis.ExpireKey(apiLimitPrefix+apiPath+":"+onlyString, life)
}
