package user

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	_ "go.mongodb.org/mongo-driver/mongo"

	"../mongodb"
	"./face"
)

var (
	ErrorPassword        = fmt.Errorf("密码错误")
	ErrorStateBanned     = fmt.Errorf("操作失败，账号处于封禁状态。")
	ErrorNoFoundUserInfo = fmt.Errorf("用户信息不存在。")

	StateDefault = 0
	StateBanned  = 1

	collectionName = "im_user"
)

//注册
func RegisterForMongoDB(user string, password string, name string) (string, error) {
	errorMsg := "注册失败"
	successMsg := "注册成功"
	exist, _, err := FindUserGetInfo(user)
	if err != nil {
		return "", mongoDB.PrintError(err, errorMsg)
	}
	if exist {
		return "", mongoDB.PrintErrorMsg(errorMsg + "：用户已存在")
	}
	faceSrc, err := face.GetRandomDefaultFaceSrc()
	if err != nil {
		return "", mongoDB.PrintError(err, errorMsg)
	}
	collection, err := mongoDB.GetCollection(collectionName)
	if err != nil {
		return "", mongoDB.PrintError(err, errorMsg)
	}
	//todo: permission group,ex:normal
	document := bson.M{
		"permission": "normal",
		"state":      StateDefault,
		"user":       user,
		"password":   password,
		"name":       name,
		"faceSrc":    faceSrc,
		"sessionid":  ""}
	_, err = collection.InsertOne(context.Background(), document)
	if err != nil {
		return "", mongoDB.PrintError(err, errorMsg)
	}
	return successMsg, nil
}

//登录，验证密码
func LoginForMongoDB(user string, password string) (string, error) {
	errorMsg := "登录失败"
	successMsg := "登录成功"
	exist, userInfo, err := FindUserGetInfo(user)
	if err != nil {
		return "", mongoDB.PrintError(err, errorMsg)
	}
	if !exist {
		return "", mongoDB.PrintErrorMsg(errorMsg + "：用户不存在")
	}
	if password != userInfo["password"] {
		return errorMsg, ErrorPassword
	}
	return successMsg, nil
}

//根据ID返回用户名
func FindIDGetName(id string) (string, error) {
	userInfo, err := FindIDGetInfo(id)
	if err != nil {
		return "", err
	}
	return userInfo["name"].(string), nil
}

//根据ID返回用户信息
func FindIDGetInfo(id string) (bson.M, error) {
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
		return nil, result.Err()
	}
	userInfo := bson.M{}
	err = result.Decode(&userInfo)
	if err != nil {
		return nil, ErrorNoFoundUserInfo
	}
	if int(userInfo["state"].(int32)) != StateDefault {
		return userInfo, ErrorStateBanned
	}
	return userInfo, nil
}

//根据用户名返回用户信息
func FindUserGetInfo(user string) (bool, bson.M, error) {
	errorMsg := "获取失败"
	collection, err := mongoDB.GetCollection(collectionName)
	if err != nil {
		return false, nil, mongoDB.PrintError(err, errorMsg)
	}
	filter := bson.M{
		"user": user}
	result := collection.FindOne(context.Background(), filter)
	if result.Err() != nil {
		return false, nil, mongoDB.PrintError(result.Err(), errorMsg)
	}
	userInfo := bson.M{}
	err = result.Decode(&userInfo)
	if err != nil {
		return false, nil, nil
	}
	if int(userInfo["state"].(int32)) != StateDefault {
		return true, userInfo, ErrorStateBanned
	}
	return true, userInfo, nil
}

//根据用户名返回sid
func FindUserGetSID(user string) (string, error) {
	exist, userInfo, err := FindUserGetInfo(user)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", fmt.Errorf("用户不存在")
	}
	sid, ok := userInfo["sessionid"]
	if !ok || sid.(string) == "" {
		return "", fmt.Errorf("sid no value")
	}
	return sid.(string), nil
}

//根据sid返回用户ID
func FindUserIDBySID(sid string) (string, error) {
	errorMsg := "获取失败"
	collection, err := mongoDB.GetCollection(collectionName)
	if err != nil {
		return "", mongoDB.PrintError(err, errorMsg)
	}
	filter := bson.M{
		"sessionid": sid}
	result := collection.FindOne(context.Background(), filter)
	if result.Err() != nil {
		return "", mongoDB.PrintError(result.Err(), errorMsg)
	}
	userInfo := bson.M{}
	err = result.Decode(&userInfo)
	if err != nil {
		return "", err
	}
	userID := userInfo["_id"].(primitive.ObjectID).Hex()
	if int(userInfo["state"].(int32)) != StateDefault {
		return userID, ErrorStateBanned
	}
	return userID, nil
}

//根据ID返回sid
func FindSIDByID(id string) (string, error) {
	userInfo, err := FindIDGetInfo(id)
	if err != nil {
		return "", err
	}
	return userInfo["sessionid"].(string), nil
}

//设置某用户名的sid
func SetSissionID(user string, sid string) error {
	errorMsg := "登录信息入库失败"
	//successMsg := "登录成功"
	collection, err := mongoDB.GetCollection(collectionName)
	if err != nil {
		return mongoDB.PrintError(err, errorMsg)
	}
	if err != nil {
		return mongoDB.PrintError(err, errorMsg)
	}
	filter := bson.M{
		"user": user}
	update := bson.M{
		"$set": bson.M{
			"sessionid": sid}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return mongoDB.PrintError(err, errorMsg)
	}
	return nil
}

//设置某用户名的sid
func ClearSissionID(sid string) error {
	errorMsg := "清除sid失败"
	collection, err := mongoDB.GetCollection(collectionName)
	if err != nil {
		return mongoDB.PrintError(err, errorMsg)
	}
	if err != nil {
		return mongoDB.PrintError(err, errorMsg)
	}
	filter := bson.M{
		"sessionid": sid}
	update := bson.M{
		"$set": bson.M{
			"sessionid": ""}}
	_, err = collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return mongoDB.PrintError(err, errorMsg)
	}
	return nil
}
