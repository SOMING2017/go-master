package message

import (
	"net/http"

	"../../../crypto/rsa"
	myWebSocket "../../../net/websocket"
	"../../../user"
	channelMain "../main"
)

func WebsocketChannelMessage(rw http.ResponseWriter, req *http.Request) {
	action := req.FormValue("action")
	if action != "WebsocketChannelMessage" {
		return
	}
	conn, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		return
	}
	channelID := req.FormValue("channelID")
	err = VerifyChannelID(channelID)
	if err != nil {
		myWebSocket.WriteLoginError(conn, err)
		return
	}
	publicKey, err := rsa.GetHttpPublicKey(req)
	if err != nil {
		myWebSocket.WriteLoginError(conn, err)
		return
	}
	userID, err := user.IsLogin(rw, req)
	if err != nil {
		myWebSocket.WriteLoginError(conn, err)
		return
	}
	err = channelMain.VerifyUserIDChannelID(userID, channelID)
	if err != nil {
		myWebSocket.WriteLoginError(conn, err)
		return
	}
	if ok := HaveClient(userID, channelID); ok {
		myWebSocket.WriteLoginError(conn, ErrorUserOnline)
		return
	}
	client := &Client{userID: userID, channelID: channelID, publicKey: publicKey, conn: conn, send: make(chan []byte)}
	manager.register <- client
	go client.channelMessageRead()
	go client.channelMessageWrite()
	GetAllMessage(client)
}
