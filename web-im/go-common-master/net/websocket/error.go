package websocket

import (
	"encoding/json"
	"net/http"

	myLog "../../log"

	myNet "../"
	"github.com/gorilla/websocket"
)

func WriteNormalMsg(conn *websocket.Conn, msg string) {
	WritStatusMsg(conn, http.StatusOK, msg)
}

func WritStatusMsg(conn *websocket.Conn, status int, msg string) {
	responseJson := myNet.ResponseJson{}
	responseJson.SetStatusMsg(status, msg)
	result, _ := json.Marshal(responseJson)
	myLog.LogErrorInfo("websocket", "result:"+string(result))
	conn.WriteMessage(websocket.TextMessage, result)
}

func WriteError(conn *websocket.Conn, err error) {
	responseJson := myNet.ResponseJson{}
	responseJson.SetError(err)
	result, _ := json.Marshal(responseJson)
	myLog.LogErrorInfo("websocket", "result:"+string(result))
	conn.WriteMessage(websocket.TextMessage, result)
}

func WriteLoginError(conn *websocket.Conn, err error) {
	WritStatusError(conn, http.StatusUnauthorized, err)
}

func WritStatusError(conn *websocket.Conn, status int, err error) {
	responseJson := myNet.ResponseJson{}
	responseJson.SetStatusError(status, err)
	result, _ := json.Marshal(responseJson)
	myLog.LogErrorInfo("websocket", "result:"+string(result))
	conn.WriteMessage(websocket.TextMessage, result)
}
