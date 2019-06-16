package main

import (
	"fmt"

	".."
)

func main() {
	TestRSA()
}

func TestRSA() {
	origData := "abcd12345"
	fmt.Println("origData is ", origData)
	rsaPrivateKey, rsaPublicKey, err := rsa.GetRSAStringKey(rsa.RSAKeyBits)
	fmt.Println("rsaPrivateKey is ", rsaPrivateKey)
	fmt.Println("rsaPublicKey is ", rsaPublicKey)
	if err != nil {
		return
	}
	base64EnData, err := rsa.RsaStringEncryptToBase64String(rsaPublicKey, origData)
	if err != nil {
		return
	}
	fmt.Println("base64EnData is ", base64EnData)
	deData, err := rsa.RsaStringDecryptForBase64String(rsaPrivateKey, base64EnData)
	fmt.Println("deData is ", deData)
}
