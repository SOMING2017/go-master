package user

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"../crypto/md5"
	"../net"
	"../session"
	pw "./password"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrorTokenNoFound = fmt.Errorf("no found token")
	ErrorLoginInfo    = fmt.Errorf("登录信息错误")
	ErrorLoginOverNum = fmt.Errorf("错误登录次数达到限制(" + strconv.Itoa(loginLimitMaxNum) + "次)，今日不可使用密码登录。")

	loginLimitMaxNum = 10
)

//根据请求获取登录信息
func GetLoginInfoRQ(rw http.ResponseWriter, req *http.Request, publicKey string, user string) (string, error) {
	sid := session.GetSessionID(rw, req)
	ipv4 := net.IP(req)
	return GetLoginInfo(sid, ipv4, publicKey, user)
}

//获取登录信息
//若有旧的登录信息，则消除
//存入新的登录信息，具体：
//存token,ip进redis
//存sid进mongodb
func GetLoginInfo(sid string, ipv4 string, publicKey string, user string) (string, error) {
	//消除旧的登录记录
	oldSid, err := FindUserGetSID(user)
	if err == nil {
		ClearSissionID(sid)
		DelLoginInfo(oldSid)
	}
	//存入新的登录记录
	nowTimeString := time.Now().String()
	token := md5.Md5(sid + ipv4 + nowTimeString)
	enToken, err := pw.PasswordEncrypt(publicKey, token)
	if err != nil {
		return "", err
	}
	err = SetLoginToken(sid, token)
	if err != nil {
		return "", err
	}
	err = SetLoginIP(sid, ipv4)
	if err != nil {
		return "", err
	}
	err = SetSissionID(user, sid)
	if err != nil {
		return "", err
	}
	return enToken, nil
}

//登录验证
//通过验证请求的token和ip是否与登录存入的一致验证
func IsLogin(rw http.ResponseWriter, req *http.Request) (string, error) {
	token := req.FormValue("token")
	if token == "" {
		return "", ErrorTokenNoFound
	}
	sid := session.GetSessionID(rw, req)
	err := isToken(sid, token)
	if err != nil {
		return "", ErrorLoginInfo
	}
	ipv4 := net.IP(req)
	err = isIP(sid, ipv4)
	if err != nil {
		return "", ErrorLoginInfo
	}
	ExtendLoginKeyLife(sid)
	userID, err := FindUserIDBySID(sid)
	return userID, err
}

//登录token验证
func isToken(sid string, token string) error {
	dbToken, err := GetLoginToken(sid)
	if err != nil {
		return err
	}
	if token != dbToken {
		return ErrorLoginInfo
	}
	return nil
}

//登录ip验证
func isIP(sid string, ip string) error {
	dbIP, err := GetLoginIP(sid)
	if err != nil {
		return err
	}
	if ip != dbIP {
		return ErrorLoginInfo
	}
	return nil
}

//密码登录限制
//满足限制条件返回err
func IsPasswordLoginLimit(user string) error {
	nowTime := time.Now()
	nextTime := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day()+1,
		0, 0, 0, 0,
		nowTime.Location())
	loginLimitSeconds := int(nextTime.Unix() - nowTime.Unix())
	resultBool := false
	sidSum, err := GetSumPasswordLoginLimit(user, loginLimitSeconds)
	if err != nil {
		return err
	}
	resultBool = resultBool || (sidSum >= loginLimitMaxNum)
	if resultBool {
		return ErrorLoginOverNum
	}
	return nil
}

//用户密码注册操作
func RegisterPassword(user string, password string, name string) (string, error) {
	msg, err := RegisterForMongoDB(user, password, name)
	if err == nil {
		//加入世界组
		//todo: 采用消息队列进行失败缓存
		go func() {
			exist, userInfo, err := FindUserGetInfo(user)
			if err == nil && exist {
				bindingUserWorldChannel(userInfo["_id"].(primitive.ObjectID).Hex())
			}
		}()
	}
	return msg, err
}

//用户密码登录操作
func LoginForPassword(user string, password string) (string, error) {
	//判断用户是否已经使用尽登录次数
	err := IsPasswordLoginLimit(user)
	if err != nil {
		return "", err
	}
	msg, err := LoginForMongoDB(user, password)
	if err == ErrorPassword {
		//记录错误次数，达到10次，隔日才能进行再登录
		go IncrPasswordLoginLimit(user)
	}
	if err == nil {
		//登录成功，清除错误次数记录
		go DelPasswordLoginLimit(user)
	}
	return msg, err
}
