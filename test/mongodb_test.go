package test

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
)

func TestMongoLog(t *testing.T) {

	fmt.Println(111111111)

	// 设置 MongoDB 连接选项
	clientOptions := options.Client().ApplyURI("mongodb://124.220.208.204:27017/")

	fmt.Println(11111)

	// 建立连接
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(err)
}
