package net

import (
	"fmt"
	"net/http"
)

var JsonType = "application/json"

func ResponseMsg(rw http.ResponseWriter, msg string, err error) {
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprint(rw, err)
		return
	}
	fmt.Fprint(rw, msg)
}
