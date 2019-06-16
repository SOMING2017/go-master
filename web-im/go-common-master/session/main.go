package session

import (
	"net/http"

	"github.com/astaxie/beego/session"
)

var GlobalSessions *session.Manager

func init() {
	sessionConfig := &session.ManagerConfig{
		CookieName:      "gosessionid",
		EnableSetCookie: true,
		Gclifetime:      1296000,
		Maxlifetime:     1296000,
		Secure:          false,
		CookieLifeTime:  1296000,
		ProviderConfig:  "./tmp",
	}
	GlobalSessions, _ = session.NewManager("memory", sessionConfig)
	go GlobalSessions.GC()
}

//通过请求头获取sessionID
func GetSessionID(rw http.ResponseWriter, req *http.Request) string {
	sess, _ := GlobalSessions.SessionStart(rw, req)
	defer sess.SessionRelease(rw)
	return sess.SessionID()
}
