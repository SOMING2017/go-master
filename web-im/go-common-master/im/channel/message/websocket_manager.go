package message

import (
	"fmt"

	"../../../chan"
	myWebSocket "../../../net/websocket"
	"github.com/gorilla/websocket"
)

const (
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	ErrorTest            = fmt.Errorf("错误测试")
	ErrorUserOnline      = fmt.Errorf("消息面板已打开")
	ErrorMessageNoFormat = fmt.Errorf("传入信息错误")

	newline  = []byte{'\n'}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	manager = &ClientManager{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
)

type Client struct {
	userID    string
	channelID string
	publicKey string
	conn      *websocket.Conn
	send      chan []byte
}

type ClientManager struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
}

func init() {
	go manager.run()
}

//开启客户端管理组
func (manager *ClientManager) run() {
	for {
		select {
		case client := <-manager.register:
			manager.clients[client.userID+client.channelID] = client
		case client := <-manager.unregister:
			if ok := HaveClient(client.userID, client.channelID); ok {
				client.conn.Close()
				goChan.SafeClose(client.send)
				delete(manager.clients, client.userID+client.channelID)
			}
		}
	}
}

//客户端是否已在管理组
func HaveClient(userID string, channelID string) bool {
	_, ok := manager.clients[userID+channelID]
	return ok
}

//某通道消息持续读取
func (client *Client) channelMessageRead() {
	client.conn.SetReadLimit(maxMessageSize)
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			manager.unregister <- client
			return
		}
		go ParseSendMessage(client, message)
	}
}

//某通道消息发送至客户端
//{status:200,action:"All",data:[{},{}]}
//{status:200,action:"Add",data:{}}
func (client *Client) channelMessageWrite() {
	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				// The hub closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			// Add queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

//发送内容至某通道消息
func (client *Client) SendContent(content string) {
	select {
	case client.send <- []byte(content):
	default:
		manager.unregister <- client
	}
}

//返回正确提示
func (client *Client) WriteNormalMsg(msg string) {
	myWebSocket.WriteNormalMsg(client.conn, msg)
}

//返回错误
func (client *Client) WriteError(err error) {
	myWebSocket.WriteError(client.conn, err)
}

//返回登录相关错误
func (client *Client) WriteLoginError(err error) {
	myWebSocket.WriteLoginError(client.conn, err)
	manager.unregister <- client
}
