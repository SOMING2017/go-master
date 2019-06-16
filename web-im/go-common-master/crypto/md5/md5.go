package md5

import (
	"crypto/md5"
	"encoding/hex"
)

//https://www.jianshu.com/p/58dcbf490ef3
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
