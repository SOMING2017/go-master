package rsa

import (
	"fmt"
	"net/http"
)

var ErrorPublicKey = fmt.Errorf("传递公钥错误")

//返回请求体的公钥
func GetHttpPublicKey(req *http.Request) (string, error) {
	publicKey := req.FormValue("publicKey")
	if IsPublicKey(publicKey) {
		return publicKey, nil
	}
	return "", ErrorPublicKey
}
