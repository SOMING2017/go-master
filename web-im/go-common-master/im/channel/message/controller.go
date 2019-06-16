package message

import (
	"net/http"
)

var ApiPath = "/im/channel/message/controller"

func Start(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	reqMethod := req.Method
	if reqMethod == "GET" {
	}
	if reqMethod == "POST" {
	}
	//websocket
	WebsocketChannelMessage(rw, req)
}
