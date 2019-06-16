package password

import (
	"net/http"

	"../../net"
	"../../session"
)

func GetPublickey(rw http.ResponseWriter, req *http.Request) {
	action := req.FormValue("action")
	if action != "GetPublickey" {
		return
	}
	sid := session.GetSessionID(rw, req)
	msg, errMsg := GetPublicKey(sid)
	net.ResponseMsg(rw, msg, errMsg)
}
