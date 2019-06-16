package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

var RSAKeyBits = 1024

//获取RSA公钥私钥
func GetRSAStringKey(bits int) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}
	//get rsa private key
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	priBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	rsaPrivateKey := pem.EncodeToMemory(priBlock)
	//get rsa public key
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}
	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	rsaPublicKey := pem.EncodeToMemory(pubBlock)
	return string(rsaPrivateKey), string(rsaPublicKey), nil
}

//判断传入公钥是否正确
func IsPublicKey(publicKey string) bool {
	_, err := RsaStringEncryptToBase64String(publicKey, "11110000")
	return err == nil
}

//RSA加密，并且base64转码
func RsaStringEncryptToBase64String(publicKey string, origData string) (string, error) {
	enData, err := RsaStringEncrypt(publicKey, origData)
	if err != nil {
		return "", err
	}
	base64EnData := base64.StdEncoding.EncodeToString([]byte(enData))
	return base64EnData, nil
}

//RSA加密，返回字符串
func RsaStringEncrypt(publicKey string, origData string) (string, error) {
	enData, err := RsaEncrypt([]byte(publicKey), []byte(origData))
	return string(enData), err
}

//RSA加密
func RsaEncrypt(publicKey []byte, origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, fmt.Errorf("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

//RSA解密，并且base64转码
func RsaStringDecryptForBase64String(privateKey string, base64EnData string) (string, error) {
	enData, err := base64.StdEncoding.DecodeString(base64EnData)
	if err != nil {
		return "", err
	}
	return RsaStringDecrypt(privateKey, string(enData))
}

//RSA解密，返回字符串
func RsaStringDecrypt(privateKey string, ciphertext string) (string, error) {
	deData, err := RsaDecrypt([]byte(privateKey), []byte(ciphertext))
	return string(deData), err
}

//RSA解密
func RsaDecrypt(privateKey []byte, ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, fmt.Errorf("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
