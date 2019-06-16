package message

import (
	myRedis "../../../../redis"
)

var (
	channelPrefix = "msg:list:channelID."
)

//新增头部消息
func LpushNewMessage(channelID string, msgStringInfo string) error {
	return myRedis.LpushKey(channelPrefix+channelID, msgStringInfo)
}

//删除某通道尾部消息
func DelLastMessage(channelID string) error {
	return myRedis.RpopKey(channelPrefix + channelID)
}

//获取某通道顶部消息
func GetHeadMessage(channelID string) (Message, error) {
	messageString, err := myRedis.LheadKey(channelPrefix + channelID)
	if err != nil {
		return Message{}, err
	}
	if messageString == "" {
		return Message{}, nil
	}
	return StringToMessage(messageString)
}

//获取某通道缓存消息，以字符串数组返回
func GetCacheMessageStrings(channelID string) ([]string, error) {
	messageStrings := make([]string, 0)
	reply, err := myRedis.LrangeKey(channelPrefix + channelID)
	if err != nil {
		return nil, err
	}
	for _, dbMessage := range reply.([]interface{}) {
		messageString := string(dbMessage.([]uint8))
		messageStrings = append(messageStrings, messageString)
	}
	revMessageStrings := make([]string, 0)
	for index := len(messageStrings) - 1; index >= 0; index-- {
		revMessageStrings = append(revMessageStrings, messageStrings[index])
	}
	return revMessageStrings, nil
}

//删除某通道
func DelCache(channelID string) error {
	return myRedis.DelKey(channelPrefix + channelID)
}

//获取某通道缓存消息数量
func GetCacheMessageLen(channelID string) (int, error) {
	return myRedis.LlenKey(channelPrefix + channelID)
}
