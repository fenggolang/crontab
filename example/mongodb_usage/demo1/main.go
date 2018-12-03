package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
)

func main() {
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
	)

	// 1. 建立连接
	if client, err = mongo.Connect(context.TODO(), "mongodb://172.17.10.210:27017", clientopt.ConnectTimeout(5*time.Second)); err != nil {
		fmt.Println(err)
		return
	}

	// 2. 选择数据库my_db
	database = client.Database("my_db")

	// 3. 选择表my_collection
	collection = database.Collection("my_collection")

	collection = collection
}
