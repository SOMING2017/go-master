package channel

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/mongo"

	"../../../mongodb"
)

var (
	ErrorNoFindInfo        = fmt.Errorf("通道不存在")
	ErrorInsertNoRule      = fmt.Errorf("置入值不满足规则")
	ErrorStateBanned       = fmt.Errorf("操作失败，账号处于封禁状态。")
	ErrorChannelNoID       = fmt.Errorf("通道不存在ID")
	ErrorChannelNoArea     = fmt.Errorf("通道不存在Area")
	ErrorChannelNoAreaName = fmt.Errorf("通道不存在AreaName")

	StateDefault      = 0
	StateBanned       = 1
	AreaIndividual    = "individual"
	AreaGroup         = "group"
	AreaWorld         = "world"
	MaxUserIndividual = 2
	MaxUserRank1      = 200
	MaxUserRank2      = 500

	collectionName = "im_channel"
)

//获取世界通道ID
//没有则新增世界通道，且当前用户为创建者
func GetWorldChannelID(masterID string) (string, error) {
	channelInfo, err := FindAreaGetOneInfo(AreaWorld)
	if err == nil || err == ErrorStateBanned {
		return channelInfo["_id"].(primitive.ObjectID).Hex(), nil
	}
	err = InsertChannel(StateDefault, AreaWorld, "", 0, masterID)
	if err != nil {
		return "", err
	}
	channelInfo, err = FindAreaGetOneInfo(AreaWorld)
	if err == nil || err == ErrorStateBanned {
		return channelInfo["_id"].(primitive.ObjectID).Hex(), nil
	}
	return "", err
}

//获取一条某区域的通道信息
//主要用途：获取世界区域的信息
func FindAreaGetOneInfo(area string) (bson.M, error) {
	errorMsg := "获取失败"
	collection, err := mongoDB.GetCollection(collectionName)
	if err != nil {
		return nil, mongoDB.PrintError(err, errorMsg)
	}
	filter := bson.M{
		"area": area}
	result := collection.FindOne(context.Background(), filter)
	if result.Err() != nil {
		return nil, mongoDB.PrintError(result.Err(), errorMsg)
	}
	userInfo := bson.M{}
	err = result.Decode(&userInfo)
	if err != nil {
		return nil, ErrorNoFindInfo
	}
	if userInfo["state"] != StateDefault {
		return userInfo, ErrorStateBanned
	}
	return userInfo, nil
}

//获取某ID的通道信息
func FindIDGetOneInfo(id string) (bson.M, error) {
	errorMsg := "获取失败"
	collection, err := mongoDB.GetCollection(collectionName)
	if err != nil {
		return nil, err
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id": objectID}
	result := collection.FindOne(context.Background(), filter)
	if result.Err() != nil {
		return nil, mongoDB.PrintError(result.Err(), errorMsg)
	}
	userInfo := bson.M{}
	err = result.Decode(&userInfo)
	if err != nil {
		return nil, ErrorNoFindInfo
	}
	if int(userInfo["state"].(int32)) != StateDefault {
		return userInfo, ErrorStateBanned
	}
	return userInfo, nil
}

//获取某ID的区域
func GetAreaByID(id string) (string, error) {
	channelInfo, err := FindIDGetOneInfo(id)
	if err != nil {
		return "", err
	}
	return channelInfo["area"].(string), nil
}

//新增通道，过滤函数
//存入值说明：
//state:默认0,封禁1
//area:通道域[individual]个人，[group]组，[world]世界组
//areaName:个人域无域名，组域[groupName]
//maxUser:个人域为2人，组域有200或500，世界组-1
//masterID:群主ID
func InsertChannel(state int, area string, areaName string, maxUser int, masterID string) error {
	if state != StateDefault && state != StateBanned {
		return ErrorInsertNoRule
	}
	if area != "individual" && area != "group" && area != "world" {
		return ErrorInsertNoRule
	}
	if area == "group" {
		if maxUser != MaxUserRank1 && maxUser != MaxUserRank2 {
			return ErrorInsertNoRule
		}
	}
	if area == "individual" {
		areaName = ""
		maxUser = MaxUserIndividual
	}
	if area == "world" {
		areaName = "世界消息"
		maxUser = -1
	}
	return insertChannel(state, area, areaName, maxUser, masterID)
}

//实际新增通道函数
func insertChannel(state int, area string, areaName string, maxUser int, masterID string) error {
	collection, err := mongoDB.GetCollection(collectionName)
	if err != nil {
		return err
	}
	document := bson.M{
		"state":    state,
		"area":     area,
		"areaName": areaName,
		"maxUser":  maxUser,
		"masterID": masterID}
	_, err = collection.InsertOne(context.Background(), document)
	if err != nil {
		return err
	}
	return nil
}

//判断通道是否群组
func IsGroupByID(id string) (bool, error) {
	area, err := GetAreaByID(id)
	if err != nil {
		return false, err
	}
	return IsGroupByArea(area), nil
}

//判断通道域是否群组
func IsGroupByArea(area string) bool {
	return area != AreaIndividual
}

//获取通道ID
func GetChannelID(channel bson.M) (string, error) {
	cid, ok := channel["_id"]
	if !ok {
		return "", ErrorChannelNoID
	}
	id, ok := cid.(primitive.ObjectID)
	if !ok {
		return "", ErrorChannelNoID
	}
	return id.Hex(), nil
}

//获取通道Area
func GetChannelArea(channel bson.M) (string, error) {
	carea, ok := channel["area"]
	if !ok {
		return "", ErrorChannelNoArea
	}
	area, ok := carea.(string)
	if !ok {
		return "", ErrorChannelNoArea
	}
	return area, nil
}

//获取通道AreaName
func GetChannelAreaName(channel bson.M) (string, error) {
	careaName, ok := channel["areaName"]
	if !ok {
		return "", ErrorChannelNoAreaName
	}
	areaName, ok := careaName.(string)
	if !ok {
		return "", ErrorChannelNoAreaName
	}
	return areaName, nil
}
