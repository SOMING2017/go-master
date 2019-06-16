package message

import (
	"encoding/json"
	"fmt"

	"time"
)

var (
	ErrorNoFormatMessage = fmt.Errorf("缓存消息格式错误，缓存已清除。")

	cacheMessageMaxNum = 100
)

type Message struct {
	Content    string    `json:"content"`
	CreateID   string    `json:"createID"`
	CreateTime time.Time `json:"createTime"`
}

//消息存入通道
//若缓存消息达到最大值，则删除尾部消息
func AddMessageToRedis(channelID string, message *Message) error {
	cacheMessageLen, err := GetCacheMessageLen(channelID)
	if err == nil {
		go func() {
			if cacheMessageLen >= cacheMessageMaxNum {
				DelLastMessage(channelID)
			}
		}()
	}
	messageString, err := MessageToString(*message)
	if err != nil {
		return err
	}
	return LpushNewMessage(channelID, messageString)
}

//根据content生成完整消息
func CreateMessageByIDContent(createID string, content string) Message {
	return Message{Content: content, CreateID: createID, CreateTime: time.Now()}
}

//获取某通道缓存消息
func GetCacheMessage(channelID string) ([]Message, error) {
	messageInfo := make([]Message, 0)
	messageStrings, err := GetCacheMessageStrings(channelID)
	if err != nil {
		return nil, err
	}
	for _, messageString := range messageStrings {
		message, err := StringToMessage(messageString)
		if err != nil {
			go DelCache(channelID)
			return nil, ErrorNoFormatMessage
		}
		messageInfo = append(messageInfo, message)
	}
	return messageInfo, nil
}

//StringToMessage
func StringToMessage(messageString string) (Message, error) {
	message := Message{}
	err := json.Unmarshal([]byte(messageString), &message)
	if err != nil {
		return Message{}, err
	}
	return message, nil
}

//MessageToString
func MessageToString(message Message) (string, error) {
	messageJson, err := json.Marshal(message)
	if err != nil {
		return "", err
	}
	return string(messageJson), nil
}
