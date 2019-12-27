package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// startTime小于某时间
// {"$lt": timestamp}
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

// {"timePoint.startTime": {"$lt": timestamp}}
type DeleteCond struct {
	beforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

func main() {
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
		delCond    *DeleteCond
		delResult  *mongo.DeleteResult
	)
	// 1. 建立连接
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://172.17.10.210:27017")); err != nil {
		fmt.Println(err)
		return
	}

	// 2. 选择数据库
	database = client.Database("cron")

	// 3. 选择表
	collection = database.Collection("log")

	// 4. 要删除开始时间早于当前时间的所有日志($lt是less than)
	// delete({"timePoint.startTime": {"$lt": 当前时间}})
	delCond = &DeleteCond{
		beforeCond: TimeBeforeCond{
			Before: time.Now().Unix(),
		},
	}

	// 执行删除
	if delResult, err = collection.DeleteMany(context.TODO(), delCond); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("删除的行数:", delResult.DeletedCount)
}
