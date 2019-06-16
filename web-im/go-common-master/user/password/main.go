package password

//获取公钥
//若没有则新建
func GetPublicKey(sid string) (string, error) {
	publicKey, err := GetPublicKeyForRedis(sid)
	if err == nil {
		return publicKey, nil
	}
	publicKey, err = NewPasswordKey(sid)
	if err != nil {
		return "", err
	}
	return publicKey, nil
}

//新建公钥私钥，并存于redis
//返回公钥
func NewPasswordKey(sid string) (string, error) {
	rsaPrivateKey, rsaPublicKey, err := GetPasswordKey()
	if err != nil {
		return "", err
	}
	go setPublicKeyToRedis(sid, rsaPublicKey)
	go setPrivateKeyToRedis(sid, rsaPrivateKey)
	return rsaPublicKey, nil
}

//获取redis私钥并解密
func PasswordDecryptForRedis(sid string, enPassword string) (string, error) {
	privateKey, err := GetPrivateKeyForRedis(sid)
	if err != nil {
		return "", err
	}
	password, err := PasswordDecrypt(privateKey, enPassword)
	if err != nil {
		return "", err
	}
	return password, nil
}
