package net

import (
	"net"
	"net/http"
)

func IP(req *http.Request) string {
	if ip, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		return ip
	}
	return req.RemoteAddr
}
