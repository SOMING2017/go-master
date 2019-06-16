package message

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	channelWebSocket ".."
	"../../../chan"
	"../../../crypto/rsa"
	myNet "../../../net"
	"../../../user"
	pw "../../../user/password"
	channelMain "../main"
	main "./main"
)

var ErrorNoFoundMessageInfo = fmt.Errorf("消息信息无法获取")

type WebMessage struct {
	IsMe       bool      `json:"isMe"`
	Src        string    `josn:"src"`
	Name       string    `json:"name"`
	EnContent  string    `json:"enContent"`
	CreateTime time.Time `json:"createTime"`
}

//WebMessage to string
func (webMessage *WebMessage) String() string {
	webMessageJson, err := json.Marshal(webMessage)
	if err != nil {
		return ""
	}
	return string(webMessageJson)
}

//将Message转换WebMessage
func MessageToWebMessage(client *Client, message main.Message) (WebMessage, error) {
	webMessage := WebMessage{}
	webMessage.IsMe = (client.userID == message.CreateID)
	enContent, err := rsa.RsaStringEncryptToBase64String(client.publicKey, message.Content)
	if err != nil {
		return WebMessage{}, err
	}
	webMessage.EnContent = enContent
	webMessage.CreateTime = message.CreateTime
	userInfo, err := user.FindIDGetInfo(message.CreateID)
	if err != nil {
		return WebMessage{}, err
	}
	webMessage.Src = userInfo["faceSrc"].(string)
	webMessage.Name = userInfo["name"].(string)
	return webMessage, nil
}

//返回通道所有缓存消息
func GetAllMessage(client *Client) {
	MessageList, err := main.GetCacheMessage(client.channelID)
	if err != nil {
		client.WriteLoginError(err)
		return
	}
	messageStringList := make([]string, 0)
	for _, message := range MessageList {
		//message to webMessage and save new list
		webMessage, err := MessageToWebMessage(client, message)
		if err != nil {
			client.WriteLoginError(err)
			return
		}
		messageStringList = append(messageStringList, webMessage.String())
	}
	data := "[" + strings.Join(messageStringList, ",") + "]"
	json := myNet.ResponseAction{Status: http.StatusOK, Action: "All", Data: string(data)}
	goChan.SafeSend(client.send, []byte(json.String()))
}

//解析获取的消息
//必须带有action
//若action为Send时，须带有enContent
//解密enContent后获取content，则获取到正确消息
func ParseSendMessage(client *Client, message []byte) {
	var result map[string]interface{}
	if err := json.Unmarshal(message, &result); err != nil {
		client.WriteError(err)
		return
	}
	action, err := ParseAction(result)
	if err != nil {
		client.WriteError(err)
		return
	}
	if action == "Send" {
		enContent, err := ParseEnContent(result)
		if err != nil {
			client.WriteError(err)
			return
		}
		//content decrypt
		sid, err := user.FindSIDByID(client.userID)
		if err != nil {
			client.WriteError(err)
			return
		}
		content, err := pw.PasswordDecryptForRedis(sid, enContent)
		if err != nil {
			client.WriteError(err)
			return
		}
		go SendMessageToOnlineAndRedisCache(client, content)
		go channelWebSocket.TopChannel(client.channelID, content, client.userID)
		//return 200,send success message
		//client.WriteNormalMsg("消息发送成功")
	}
}

func ParseAction(result map[string]interface{}) (string, error) {
	action, ok := result["action"]
	if !ok {
		return "", ErrorMessageNoFormat
	}
	action, ok = action.(string)
	if !ok {
		return "", ErrorMessageNoFormat
	}
	return action.(string), nil
}

func ParseEnContent(result map[string]interface{}) (string, error) {
	enContent, ok := result["enContent"]
	if !ok {
		return "", ErrorMessageNoFormat
	}
	enContent, ok = enContent.(string)
	if !ok {
		return "", ErrorMessageNoFormat
	}
	return enContent.(string), nil
}

//消息的后续处理
//群组消息：发送至当前通道的所有在线成员，并存入redis
//个人消息：判断当前通道在线成员数是否满足最大成员，满足则直接发送
//不满足则增加存入redis步骤
//en:
//send to message panel
//if channel is group or all panel have offline,add to redis
func SendMessageToOnlineAndRedisCache(client *Client, content string) {
	isGroup, err := channelMain.IsGroupByID(client.channelID)
	if err != nil {
		client.WriteError(err)
		return
	}
	message := main.CreateMessageByIDContent(client.userID, content)
	go func() {
		var AddRedis = false
		if isGroup {
			AddRedis = true
		} else {
			if GetChannelIDUserSize(client.channelID) != channelMain.MaxUserIndividual {
				AddRedis = true
			}
		}
		if AddRedis {
			main.AddMessageToRedis(client.channelID, &message)
		}
	}()
	go SendMessageToThisChannelOnlineAllUser(client.channelID, message)
}

//获取某通道所有在线成员数
func GetChannelIDUserSize(channelID string) int {
	size := 0
	for _, client := range manager.clients {
		if client.channelID == channelID {
			size++
		}
	}
	return size
}

//发送消息至某通道所有成员
func SendMessageToThisChannelOnlineAllUser(channelID string, message main.Message) {
	for _, client := range manager.clients {
		if client.channelID == channelID {
			webMessage, err := MessageToWebMessage(client, message)
			webMessageJson, err := json.Marshal(webMessage)
			if err != nil {
				client.WriteError(err)
				return
			}
			json := myNet.ResponseAction{Status: http.StatusOK, Action: "Add", Data: string(webMessageJson)}
			client.SendContent(json.String())
		}
	}
}
