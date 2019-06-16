package mongoDB

import (
	"context"
	"time"

	myLog "../log"

	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DriverName = "mongodb"

	user           = ""
	password       = ""
	ipv4           = "localhost"
	dbPort         = "27017"
	testPingDBName = ""
	dbName         = "go_master"
)

func GetDsn() string {
	loginMsg := ""
	if user != "" && password != "" {
		loginMsg = user + ":" + password + "@"
	}
	portMsg := ""
	if dbPort != "" {
		portMsg = ":" + dbPort
	}
	extraMsg := ""
	extraMsg += "gssapiServiceName=" + DriverName
	return DriverName + "://" + loginMsg + ipv4 + portMsg + "/" + testPingDBName + "?" + extraMsg
}

func GetClient() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(GetDsn()))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetDatabase() (*mongo.Database, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}
	return client.Database(dbName), nil
}

//返回集合
func GetCollection(name string) (*mongo.Collection, error) {
	db, err := GetDatabase()
	if err != nil {
		return nil, err
	}
	return db.Collection(name), nil
}

//mogodb相关错误处理
func PrintErrorMsg(errString string) error {
	return PrintError(fmt.Errorf(errString), errString)
}

func PrintError(err error, errString string) error {
	fmt.Println(err)
	myLog.LogErrorInfo("mongodb", "err:"+err.Error(), "errString:"+errString)
	return fmt.Errorf(errString)
}
