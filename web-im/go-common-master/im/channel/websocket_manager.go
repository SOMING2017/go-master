package channel

import (
	"fmt"

	"../../chan"
	myWebSocket "../../net/websocket"
	"github.com/gorilla/websocket"
)

const (
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	ErrorTest               = fmt.Errorf("错误测试")
	ErrorUserOnline         = fmt.Errorf("用户已在线")
	ErrorNoFoundChannelInfo = fmt.Errorf("通道信息无法获取")

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
	userID string
	conn   *websocket.Conn
	send   chan []byte
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
			manager.clients[client.userID] = client
		case client := <-manager.unregister:
			if _, ok := manager.clients[client.userID]; ok {
				client.conn.Close()
				goChan.SafeClose(client.send)
				delete(manager.clients, client.userID)
			}
		}
	}
}

//某通道持续读取
//用于监听websocket是否关闭
func (client *Client) myChannelRead() {
	client.conn.SetReadLimit(maxMessageSize)
	for {
		_, _, err := client.conn.ReadMessage()
		if err != nil {
			manager.unregister <- client
			return
		}
	}
}

//某通道发送至客户端
//{status:200,action:"All",data:[{},{}]}
//{status:200,action:"Add",data:{}}
//{status:200,action:"Top",data:{}}
func (client *Client) myChannelWrite() {
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
			client.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

//发送内容至某通道
func (client *Client) SendContent(content string) {
	select {
	case client.send <- []byte(content):
	default:
		manager.unregister <- client
	}
}

//返回登录相关错误
func (client *Client) WriteLoginError(err error) {
	myWebSocket.WriteLoginError(client.conn, err)
	manager.unregister <- client
}
