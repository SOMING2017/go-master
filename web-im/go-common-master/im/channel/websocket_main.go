package channel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	myNet "../../net"
	"../../user"

	main "./main"
	messageMain "./message/main"
	"go.mongodb.org/mongo-driver/bson"
)

type WebChannel struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	NewMessage string    `json:"newMessage"`
	NewTime    time.Time `json:"newTime"`
}

//WebChannel to string
func (webChannel *WebChannel) String() string {
	getChannelJson, err := json.Marshal(webChannel)
	if err != nil {
		return ""
	}
	return string(getChannelJson)
}

//Channel to WebChannel
func ChannelToWebChannel(channel bson.M, myUserID string) (WebChannel, error) {
	lastErr := fmt.Errorf("")
	lastErr = nil
	webChannel := WebChannel{}
	id, err := main.GetChannelID(channel)
	if err != nil {
		lastErr = err
	}
	webChannel.ID = id
	name, err := GetChannelName(channel, myUserID)
	if err != nil {
		lastErr = err
	}
	webChannel.Name = name
	message, err := messageMain.GetHeadMessage(id)
	if err != nil {
		lastErr = err
	}
	webChannel.NewMessage = message.Content
	webChannel.NewTime = message.CreateTime
	return webChannel, lastErr
}

//返回用户所有通道信息
//通道信息
func GetAllChannel(client *Client) {
	channelInfo, err := main.GetChannelInfoByUserID(client.userID)
	if err != nil {
		client.WriteLoginError(ErrorNoFoundChannelInfo)
		return
	}
	channelStringInfo := make([]string, 0)
	for _, channel := range channelInfo {
		webChannel, lastErr := ChannelToWebChannel(channel, client.userID)
		if lastErr != nil {
			client.WriteLoginError(lastErr)
			return
		}
		channelStringInfo = append(channelStringInfo, webChannel.String())
	}
	data := "[" + strings.Join(channelStringInfo, ",") + "]"
	json := myNet.ResponseAction{Status: http.StatusOK, Action: "All", Data: string(data)}
	client.SendContent(json.String())
}

//获取通道名字
//若通道为群组通道，返回群组名
//若通道为私人通道，返回另一个用户的名称
func GetChannelName(channel bson.M, myUserID string) (string, error) {
	area, err := main.GetChannelArea(channel)
	if err != nil {
		return "", err
	}
	isGroup := main.IsGroupByArea(area)
	if isGroup {
		return main.GetChannelAreaName(channel)
	} else {
		otherUserID, err := main.GetIndividualChannelOtherUserID(channel, myUserID)
		if err != nil {
			return "", err
		}
		name, err := user.FindIDGetName(otherUserID)
		if err != nil {
			return "", err
		}
		return name, nil
	}
	return "", nil
}

//置顶某通道
//向该通道所有在线的用户发送行为
//置顶该通道所有用户的redis
//en:
//send top action to this channel all user
//and update redis top
func TopChannel(channelID string, newMessage string, myUserID string) {
	getChannel, err := CreateWebChannelByIDMsg(channelID, newMessage, myUserID)
	if err != nil {
		return
	}
	isGroup, err := main.IsGroupByID(channelID)
	if err != nil {
		return
	}
	userIDs, err := main.GetUserIDsByChannelID(channelID)
	if err != nil {
		return
	}
	go func() {
		for _, userID := range userIDs {
			if client, ok := manager.clients[userID]; ok {
				json := myNet.ResponseAction{Status: http.StatusOK, Action: "Top", Data: getChannel.String()}
				client.SendContent(json.String())
			}
		}
	}()
	go main.TopChannelByUserIDs(userIDs, channelID, isGroup)
}

//传入channelID和newMessage，新建WebChannel
func CreateWebChannelByIDMsg(channelID string, newMessage string, myUserID string) (WebChannel, error) {
	webChannel := WebChannel{}
	channel, err := main.FindIDGetOneInfo(channelID)
	if err != nil {
		return webChannel, err
	}
	channelName, err := GetChannelName(channel, myUserID)
	if err != nil {
		return webChannel, err
	}
	webChannel.ID = channelID
	webChannel.Name = channelName
	webChannel.NewMessage = newMessage
	webChannel.NewTime = time.Now()
	return webChannel, nil
}

//新增某通道
//todo: send add action to this channel all user
//and binding user and channel for redis
func AddChannel() {

}

//删除某通道
//todo: send delete action to this channel all user
//and unbinding user and channel for redis
//last delete this channel for redis
func DeleteChannel() {

}
