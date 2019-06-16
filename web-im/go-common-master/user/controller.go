package user

import (
	"net/http"
)

var ApiPath = "/user/controller"

func Start(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	reqMethod := req.Method
	if reqMethod == "GET" {
	}
	if reqMethod == "POST" {
		Register(rw, req)
		Login(rw, req)
	}
}
