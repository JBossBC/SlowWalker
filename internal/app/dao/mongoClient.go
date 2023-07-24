package dao

import (
	"context"
	"fmt"
	"log"
	"replite_web/internal/app/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mClient *mongo.Client
)

func init() {
	var mongoConfig = config.DBConfig.MongoConfig
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tempClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConfig.URL))
	if err != nil {
		panic(fmt.Sprintf("connecting the mongoDB %s is error: %v", mongoConfig.URL, err))
	}
	err = tempClient.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to MongoDB")
	mClient = tempClient
	// to init the mongo schema file
	if mongoConfig.Init == "false" {
		InitMogoSchema()
	}
}

func newMongoConn(dbName string) *mongo.Database {
	return mClient.Database(dbName)
}
func getMongoConn() *mongo.Database {
	return newMongoConn(config.DBConfig.MongoConfig.Database)
}
func getMongoClient() *mongo.Client {
	return mClient
}
