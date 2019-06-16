package user

import (
	"net/http"

	myLog "../log"

	"../api-limit"
	"../crypto/rsa"
	"../net"
	"../session"
	pw "./password"
)

func Register(rw http.ResponseWriter, req *http.Request) {
	action := req.PostFormValue("action")
	if action != "Register" {
		return
	}
	sid := session.GetSessionID(rw, req)
	ipv4 := net.IP(req)
	if err := limit.IsNormalApiLimit(ApiPath+"/Register", sid, ipv4); err != nil {
		net.ResponseMsg(rw, "", err)
		return
	}
	user, password, name, err := getRegisterForm(rw, req)
	if err != nil {
		net.ResponseMsg(rw, "", err)
		return
	}
	myLog.LogApiInfo(ApiPath, "sid:"+sid, "ipv4:"+ipv4,
		"user:"+user, "password:"+password, "name:"+name)
	msg, errMsg := RegisterPassword(user, password, name)
	errMsgStr := ""
	if errMsg != nil {
		errMsgStr = errMsg.Error()
	}
	myLog.LogApiInfo(ApiPath, "sid:"+sid, "ipv4:"+ipv4,
		"msg:"+msg, "errMsg:"+errMsgStr)
	net.ResponseMsg(rw, msg, errMsg)
}

func getRegisterForm(rw http.ResponseWriter, req *http.Request) (string, string, string, error) {
	user, err := getUser(rw, req)
	if err != nil {
		return "", "", "", err
	}
	password, err := getPassword(rw, req)
	if err != nil {
		return "", "", "", err
	}
	name, err := getName(rw, req)
	if err != nil {
		return "", "", "", err
	}
	return user, password, name, nil
}

func Login(rw http.ResponseWriter, req *http.Request) {
	action := req.PostFormValue("action")
	if action != "Login" {
		return
	}
	sid := session.GetSessionID(rw, req)
	ipv4 := net.IP(req)
	if err := limit.IsNormalApiLimit(ApiPath+"/Login", sid, ipv4); err != nil {
		net.ResponseMsg(rw, "", err)
		return
	}
	publicKey, err := rsa.GetHttpPublicKey(req)
	if err != nil {
		net.ResponseMsg(rw, "", err)
		return
	}
	user, password, err := getLoginForm(rw, req)
	if err != nil {
		net.ResponseMsg(rw, "", err)
		return
	}
	myLog.LogApiInfo(ApiPath, "sid:"+sid, "ipv4:"+ipv4,
		"user:"+user, "password:"+password)
	msg, errMsg := LoginForPassword(user, password)
	if errMsg != nil {
		net.ResponseMsg(rw, msg, errMsg)
		return
	}
	enToken, errMsg := GetLoginInfoRQ(rw, req, publicKey, user)
	errMsgStr := ""
	if errMsg != nil {
		errMsgStr = errMsg.Error()
	}
	myLog.LogApiInfo(ApiPath, "sid:"+sid, "ipv4:"+ipv4,
		"enToken:"+enToken, "msg:"+msg, "errMsg:"+errMsgStr)
	net.ResponseMsg(rw, enToken, errMsg)
}

func getLoginForm(rw http.ResponseWriter, req *http.Request) (string, string, error) {
	user, err := getUser(rw, req)
	if err != nil {
		return "", "", err
	}
	password, err := getPassword(rw, req)
	if err != nil {
		return "", "", err
	}
	return user, password, nil
}

func getUser(rw http.ResponseWriter, req *http.Request) (string, error) {
	enUser := req.PostFormValue("enUser")
	sid := session.GetSessionID(rw, req)
	user, err := pw.PasswordDecryptForRedis(sid, enUser)
	if err != nil {
		return "", err
	}
	err = VerifyUser(user)
	if err != nil {
		return "", err
	}
	return user, nil
}

func getPassword(rw http.ResponseWriter, req *http.Request) (string, error) {
	enPassword := req.PostFormValue("enPassword")
	sid := session.GetSessionID(rw, req)
	password, err := pw.PasswordDecryptForRedis(sid, enPassword)
	if err != nil {
		return "", err
	}
	err = VerifyPassword(password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func getName(rw http.ResponseWriter, req *http.Request) (string, error) {
	enName := req.PostFormValue("enName")
	sid := session.GetSessionID(rw, req)
	name, err := pw.PasswordDecryptForRedis(sid, enName)
	if err != nil {
		return "", err
	}
	err = VerifyName(name)
	if err != nil {
		return "", err
	}
	return name, nil
}
