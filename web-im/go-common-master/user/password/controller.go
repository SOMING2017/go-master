package password

import (
	"net/http"
)

var ApiPath = "/user/password/controller"

func Start(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	reqMethod := req.Method
	if reqMethod == "GET" {
		GetPublickey(rw, req)
	}
	if reqMethod == "POST" {

	}
}
