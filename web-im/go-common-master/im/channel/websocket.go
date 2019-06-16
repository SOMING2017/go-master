package channel

import (
	"net/http"

	"../../user"

	myWebSocket "../../net/websocket"
)

func WebsocketMyChannel(rw http.ResponseWriter, req *http.Request) {
	action := req.FormValue("action")
	if action != "WebsocketMyChannel" {
		return
	}
	conn, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		return
	}
	userID, err := user.IsLogin(rw, req)
	if err != nil {
		myWebSocket.WriteLoginError(conn, err)
		return
	}
	if _, ok := manager.clients[userID]; ok {
		myWebSocket.WriteLoginError(conn, ErrorUserOnline)
		return
	}
	client := &Client{userID: userID, conn: conn, send: make(chan []byte)}
	manager.register <- client
	go client.myChannelRead()
	go client.myChannelWrite()
	GetAllChannel(client)
}
