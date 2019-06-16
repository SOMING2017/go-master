package user

import (
	myRedis "../redis"
)

var (
	loginTokenPrefix = "login:token:sid."
	loginIPPrefix    = "login:ip:sid."

	loginLife = 1296000

	loginLimitPrefix = "login:password:limit:user."
)

//设置登录token
func SetLoginToken(sid string, token string) error {
	return myRedis.SetExpireKey(loginTokenPrefix+sid, token, loginLife)
}

//设置登录ip
//todo: 也可以做成zset存放，处理其他可能性情况
func SetLoginIP(sid string, ip string) error {
	return myRedis.SetExpireKey(loginIPPrefix+sid, ip, loginLife)
}

//获取登录token
func GetLoginToken(sid string) (string, error) {
	return myRedis.GetString(loginTokenPrefix + sid)
}

//获取登录ip
func GetLoginIP(sid string) (string, error) {
	return myRedis.GetString(loginIPPrefix + sid)
}

//删除登录信息(含token,ip)
func DelLoginInfo(sid string) {
	myRedis.DelKey(loginTokenPrefix + sid)
	myRedis.DelKey(loginIPPrefix + sid)
}

//延长登录时长
func ExtendLoginKeyLife(sid string) {
	myRedis.ExpireKey(loginTokenPrefix+sid, loginLife)
	myRedis.ExpireKey(loginIPPrefix+sid, loginLife)
}

//登录失败次数+1
func IncrPasswordLoginLimit(user string) {
	myRedis.INCRKey(loginLimitPrefix + user)
}

//获取登录失败次数
func GetSumPasswordLoginLimit(user string, life int) (int, error) {
	myRedis.SetNXExpireKey(loginLimitPrefix+user, 0, life)
	return myRedis.GetInt(loginLimitPrefix + user)
}

//删除登录失败记录
func DelPasswordLoginLimit(user string) {
	myRedis.DelKey(loginLimitPrefix + user)
}
