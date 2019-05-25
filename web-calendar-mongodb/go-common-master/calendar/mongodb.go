// mongodb
package calendar

import (
	"context"
	"fmt"
	"strconv"

	"encoding/json"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"../mongodb"
)

var collectionName = "calendar_tag"

func GetCalendarNoticeInfoForMongoDB(startDatetime time.Time, endIndex *int) (string, error) {
	errorMsg := "获取失败"
	errorIndex := 0
	collection, err := mongoDB.GetCollection(collectionName)
	errorIndex++
	if err != nil {
		return errorMsg, printError(err, errorIndex)
	}
	endDatetime := time.Date(startDatetime.Year(), startDatetime.Month(), startDatetime.Day()+*endIndex-1,
		startDatetime.Hour(), startDatetime.Minute(), startDatetime.Second(), startDatetime.Nanosecond(),
		startDatetime.Location())
	match := bson.M{
		"isDiscard":        0,
		"createFormatDate": bson.M{"$gte": startDatetime, "$lte": endDatetime}}
	group := bson.M{
		"_id": "$createFormatDate"}
	sort := bson.M{
		"createFormatDate": 1}

	pipeline := []bson.M{
		bson.M{"$match": match},
		bson.M{"$group": group},
		bson.M{"$sort": sort}}
	cur, err := collection.Aggregate(context.Background(), pipeline)
	errorIndex++
	if err != nil {
		return errorMsg, printError(err, errorIndex)
	}
	defer cur.Close(context.Background())
	noticeInfos := make([]string, *endIndex)
	for key, _ := range noticeInfos {
		noticeInfos[key] = strconv.FormatBool(false)
	}
	noticeDateTimes := make([]time.Time, 0)
	errorIndex++
	for cur.Next(context.Background()) {
		tagInfo := bson.M{}
		err := cur.Decode(&tagInfo)
		if err != nil {
			return errorMsg, printError(err, errorIndex)
		}
		var createFormatDate time.Time
		_id := tagInfo["_id"]
		switch _id.(type) {
		case primitive.DateTime:
			createFormatDate = _id.(primitive.DateTime).Time()
			createFormatDate = createFormatDate.In(startDatetime.Location())
		default:
			errorIndex++
			return errorMsg, printError(err, errorIndex)
		}
		noticeDateTimes = append(noticeDateTimes, createFormatDate)
	}
	err = cur.Err()
	errorIndex++
	if err != nil {
		return errorMsg, printError(err, errorIndex)
	}
	for _, value := range noticeDateTimes {
		offsetDay := value.YearDay() - startDatetime.YearDay()
		noticeInfos[offsetDay] = strconv.FormatBool(true)
	}
	result := "[" + strings.Join(noticeInfos, ",") + "]"
	return result, nil
}

func GetTagInfoForMongoDB(selectDatetime time.Time) (string, error) {
	errorMsg := "获取失败"
	errorIndex := 0
	collection, err := mongoDB.GetCollection(collectionName)
	errorIndex++
	if err != nil {
		return errorMsg, printError(err, errorIndex)
	}
	filter := bson.M{
		"isDiscard":        0,
		"createFormatDate": selectDatetime}
	opts := &options.FindOptions{}
	opts.Projection = bson.M{
		"tag":        1,
		"content":    1,
		"createDate": 1}
	opts.Sort = bson.M{
		"createDate": -1}
	cur, err := collection.Find(context.Background(), filter, opts)
	errorIndex++
	if err != nil {
		return errorMsg, printError(err, errorIndex)
	}
	defer cur.Close(context.Background())
	tagInfos := make([]string, 0)
	errorIndex++
	for cur.Next(context.Background()) {
		tagInfo := bson.M{}
		err := cur.Decode(&tagInfo)
		if err != nil {
			return errorMsg, printError(err, errorIndex)
		}
		tagInfoJson, err := json.Marshal(tagInfo)
		if err != nil {
			return errorMsg, printError(err, errorIndex)
		}
		tagInfos = append(tagInfos, string(tagInfoJson))
	}
	err = cur.Err()
	errorIndex++
	if err != nil {
		return errorMsg, printError(err, errorIndex)
	}
	result := "[" + strings.Join(tagInfos, ",") + "]"
	return result, nil
}

func AddNewTagForMongoDB(tag *string, content *string, selectDatetime time.Time) (string, error) {
	errorMsg := "新增失败"
	successMsg := "新增成功"
	errorIndex := 0
	collection, err := mongoDB.GetCollection(collectionName)
	errorIndex++
	if err != nil {
		return errorMsg, printError(err, errorIndex)
	}
	createDate := time.Now()
	createFormatDate := selectDatetime
	document := bson.M{
		"isDiscard":        0,
		"tag":              tag,
		"content":          content,
		"createDate":       createDate,
		"createFormatDate": createFormatDate}
	_, err = collection.InsertOne(context.Background(), document)
	errorIndex++
	if err != nil {
		return errorMsg, printError(err, errorIndex)
	}
	return successMsg, nil
}

func DiscardOldTagForMongoDB(cid *string) (string, error) {
	errorMsg := "废除便签失败"
	successMsg := "废除便签成功"
	errorIndex := 0
	collection, err := mongoDB.GetCollection(collectionName)
	errorIndex++
	if err != nil {
		return errorMsg, printError(err, errorIndex)
	}
	_id, err := primitive.ObjectIDFromHex(*cid)
	errorIndex++
	if err != nil {
		return errorMsg, printError(err, errorIndex)
	}
	filter := bson.M{
		"_id": _id}
	update := bson.M{
		"$set": bson.M{
			"isDiscard": 1}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	errorIndex++
	if err != nil {
		return errorMsg, printError(err, errorIndex)
	}
	return successMsg, nil
}

func printError(err error, index int) error {
	fmt.Println("断点" + strconv.Itoa(index))
	fmt.Println(err)
	return fmt.Errorf("error" + strconv.Itoa(index))
}
