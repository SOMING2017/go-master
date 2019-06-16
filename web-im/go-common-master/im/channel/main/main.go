package channel

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrorNoUserChannel       = fmt.Errorf("该群组不属于用户")
	ErrorNoIndividualChannel = fmt.Errorf("不是私人通道")
)

//用户绑定世界组
func BindingUserWorldChannel(userID string) error {
	worldChannelID, err := GetWorldChannelID(userID)
	if err != nil {
		return err
	}
	return BindingUserChannelForRedis(userID, worldChannelID, true)
}

//获取用户通道信息
//根据用户ID获取此redis存入的所有通道ID
//根据通道ID获取存入mongodb的通道信息
func GetChannelInfoByUserID(userID string) ([]bson.M, error) {
	channelInfo := make([]bson.M, 0)
	channelIDs, err := GetChannelIDsByUserID(userID)
	if err != nil {
		return channelInfo, err
	}
	for _, channelStringID := range channelIDs {
		channel, err := FindIDGetOneInfo(channelStringID)
		if err == nil {
			channelInfo = append(channelInfo, channel)
		}
	}
	return channelInfo, nil
}

//置顶某用户数组的某通道
func TopChannelByUserIDs(userIDs []string, channelID string, isGroup bool) error {
	lastErr := fmt.Errorf("")
	lastErr = nil
	for _, userID := range userIDs {
		err := AddChannelToUser(userID, channelID, isGroup)
		if err != nil {
			lastErr = err
		}
	}
	return lastErr
}

//验证某通道是否属于某用户
func VerifyUserIDChannelID(userID string, channelID string) error {
	channelIDs, err := GetChannelIDsByUserID(userID)
	if err != nil {
		return err
	}
	for _, channelStringID := range channelIDs {
		if channelStringID == channelID {
			return nil
		}
	}
	return ErrorNoUserChannel
}

//获取通道另一个用户ID
func GetIndividualChannelOtherUserID(channel bson.M, myUserID string) (string, error) {
	area := channel["area"].(string)
	isGroup := IsGroupByArea(area)
	if !isGroup {
		return "", ErrorNoIndividualChannel
	} else {
		id, err := GetChannelID(channel)
		if err != nil {
			return "", err
		}
		userIDs, err := GetUserIDsByChannelID(id)
		if err != nil {
			return "", err
		}
		for _, userID := range userIDs {
			if userID != myUserID {
				return userID, nil
			}
		}
	}
	return "", nil
}
