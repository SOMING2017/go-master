package channel

import (
	"time"

	myRedis "../../../redis"
)

var (
	channelPrefix        = "userSet:channelID."
	userPrefix           = "channelSet:userID."
	userIndividualPrefix = "individual:" + userPrefix
	userGroupPrefix      = "group:" + userPrefix
)

//用户与通道相互绑定
func BindingUserChannelForRedis(userID string, channelID string, isGroup bool) error {
	err1 := AddUserToChannel(channelID, userID)
	err2 := AddChannelToUser(userID, channelID, isGroup)
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}

//向某用户新增通道
func AddChannelToUser(userID string, channelID string, isGroup bool) error {
	nowUnix := time.Now().Unix()
	if isGroup {
		return myRedis.ZaddKey(userGroupPrefix+userID, int(nowUnix), channelID)
	} else {
		return myRedis.ZaddKey(userIndividualPrefix+userID, int(nowUnix), channelID)
	}
}

//向某通道增加用户
func AddUserToChannel(channelID string, userID string) error {
	nowUnix := time.Now().Unix()
	return myRedis.ZaddKey(channelPrefix+channelID, int(nowUnix), userID)
}

//获取某用户所有通道ID
func GetChannelIDsByUserID(userID string) ([]string, error) {
	channelIDs := make([]string, 0)
	err := myRedis.ZunionStoreTwo(userPrefix+userID,
		userIndividualPrefix+userID, userGroupPrefix+userID)
	if err != nil {
		return channelIDs, err
	}
	reply, err := myRedis.ZrevRangeKey(userPrefix + userID)
	if err != nil {
		return channelIDs, err
	}
	for _, channelID := range reply.([]interface{}) {
		channelStringID := string(channelID.([]uint8))
		channelIDs = append(channelIDs, channelStringID)
	}
	return channelIDs, nil
}

//获取某通道所有用户ID
func GetUserIDsByChannelID(channelID string) ([]string, error) {
	userIDs := make([]string, 0)
	reply, err := myRedis.ZrevRangeKey(channelPrefix + channelID)
	if err != nil {
		return userIDs, err
	}
	for _, userID := range reply.([]interface{}) {
		userStringID := string(userID.([]uint8))
		userIDs = append(userIDs, userStringID)
	}
	return userIDs, nil
}
