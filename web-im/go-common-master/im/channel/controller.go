package channel

import (
	"net/http"
)

var ApiPath = "/im/channel/controller"

func Start(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	reqMethod := req.Method
	if reqMethod == "GET" {
	}
	if reqMethod == "POST" {
	}
	//websocket
	WebsocketMyChannel(rw, req)
}
