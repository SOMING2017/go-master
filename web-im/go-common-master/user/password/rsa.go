package password

import (
	"../../crypto/rsa"
)

func GetPasswordKey() (string, string, error) {
	return rsa.GetRSAStringKey(rsa.RSAKeyBits)
}

func PasswordEncrypt(rsaPublicKey string, password string) (string, error) {
	return rsa.RsaStringEncryptToBase64String(rsaPublicKey, password)
}

func PasswordDecrypt(rsaPrivateKey string, enPassword string) (string, error) {
	return rsa.RsaStringDecryptForBase64String(rsaPrivateKey, enPassword)
}
