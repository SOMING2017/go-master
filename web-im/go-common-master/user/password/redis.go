package password

import (
	"fmt"

	myRedis "../../redis"
)

var (
	passwordPublicKeyPrefix  = "password:public:key:sid."
	passwordPrivateKeyPrefix = "password:private:key:sid."

	passwordLife = 60
)

//获取redis暂存的公钥
func GetPublicKeyForRedis(sid string) (string, error) {
	publicKey, err := myRedis.GetString(passwordPublicKeyPrefix + sid)
	if err == nil {
		go ExtendKeyLife(sid)
	}
	if err != nil && err == myRedis.ErrorNoFoundValue {
		err = fmt.Errorf("没找到密钥")
	}
	return publicKey, err
}

//获取redis暂存的私钥
func GetPrivateKeyForRedis(sid string) (string, error) {
	privateKey, err := myRedis.GetString(passwordPrivateKeyPrefix + sid)
	if err != nil && err == myRedis.ErrorNoFoundValue {
		err = fmt.Errorf("没找到密钥")
	}
	return privateKey, err
}

//延长暂存的公钥私钥寿命
func ExtendKeyLife(sid string) {
	myRedis.ExpireKey(passwordPublicKeyPrefix+sid, passwordLife)
	myRedis.ExpireKey(passwordPrivateKeyPrefix+sid, passwordLife)
}

//公钥存入redis
func setPublicKeyToRedis(sid string, publicKey string) {
	myRedis.SetExpireKey(passwordPublicKeyPrefix+sid, publicKey, passwordLife)
	//todo: if err->log,class:redis

}

//私钥存入redis
func setPrivateKeyToRedis(sid string, privateKey string) {
	myRedis.SetExpireKey(passwordPrivateKeyPrefix+sid, privateKey, passwordLife)
	//todo: if err->log,class:redis
}
