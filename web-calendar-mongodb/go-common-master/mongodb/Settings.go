// Settings
package mongoDB

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var user = ""
var password = ""
var ipv4 = "localhost"
var dbPort = "27017"
var dbName = "go_master"

var DriverName = "mongodb"

func GetDsn() string {
	loginMsg := ""
	if user != "" && password != "" {
		loginMsg = user + ":" + password + "@"
	}
	return DriverName + "://" + loginMsg + ipv4 + ":" + dbPort + "/?gssapiServiceName=" + DriverName
}

func GetClient() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(GetDsn()))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
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

func GetCollection(name string) (*mongo.Collection, error) {
	db, err := GetDatabase()
	if err != nil {
		return nil, err
	}
	return db.Collection(name), nil
}
