package net

import (
	"encoding/json"
	"net/http"
)

type ResponseJson struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   string `json:"data"`
}

type ResponseAction struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Action string `json:"action"`
	Data   string `json:"data"`
}

func (rj *ResponseJson) String() string {
	rjJson, err := json.Marshal(rj)
	if err == nil {
		return string(rjJson)
	}
	return ""
}

func (rj *ResponseJson) SetNotFound() {
	rj.SetStatusMsg(http.StatusNotFound, "一个服务端问题，请联系开发人员")
}

func (rj *ResponseJson) SetError(err error) {
	rj.SetStatusMsg(http.StatusNotFound, err.Error())
}

func (rj *ResponseJson) SetStatusError(status int, err error) {
	rj.SetStatusMsg(status, err.Error())
}

func (rj *ResponseJson) SetStatusMsg(status int, msg string) {
	rj.Status = status
	rj.Msg = msg
	rj.Data = ""
}

func (ra *ResponseAction) String() string {
	raJson, err := json.Marshal(ra)
	if err == nil {
		return string(raJson)
	}
	return ""
}
