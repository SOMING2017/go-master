// ResponseController
package net

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

var JsonType = "application/json"

//todo post调用无效
func ReturnResponse(req *http.Request, rw http.ResponseWriter, rj ResponseJson) {
	reqContentType := req.Header.Get("Content-Type")
	rw.Header().Set("Content-Type", reqContentType)
	if reqContentType == JsonType {
		returnJsonResponse(rw, rj)
	}
}

func returnJsonResponse(rw http.ResponseWriter, rj ResponseJson) {
	rwJson, err := json.Marshal(rj)
	if err != nil {
		rw.WriteHeader(HttpNomalErrorStatus)
		rj.Status = strconv.Itoa(HttpNomalErrorStatus)
		rj.Msg = "一个服务端问题，请联系开发人员"
		rj.Data = ""
		rwJson, _ = json.Marshal(rj)
	}
	// rw.Write(rwJson)
	rwJsonStr := string(rwJson)
	fmt.Println(rwJsonStr)
	fmt.Fprintf(rw, "%#v", rwJsonStr)
}
