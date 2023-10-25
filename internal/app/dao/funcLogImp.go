package dao

import (
	"replite_web/internal/app/config"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type FuncLogDao struct {
}

var (
	functionLogDao  *FuncLogDao
	functionLogOnce sync.Once
)

func getFuncLogDao() *FuncLogDao {
	functionLogOnce.Do(func() {
		functionLogDao = new(FuncLogDao)
	})
	return functionLogDao
}

// if wanting to scan the details, the operation should display the task object
type FunctionLog struct {
	FunctionName    string    `json:"functionName" bson:"functionName"`
	FunctionCreator string    `json:"functionCreator" bson:"functionCreator"`
	User            string    `json:"user" bson:"user"`
	UseTime         time.Time `json:"useTime" bson:"useTime"`
	Ip              string    `json:"ip" bson:"ip"`
	// according the task status,maybe has situation as follow: push to kafka,wait to execute,has executed.the field is approximately introduce about the stage for task
	Message string `json:"message" bson:"message"`
}

const DEFAULT_FUNCTIONLOG_COLLECTION = "funclog"

var (
	funcLogCollection     *mongo.Collection
	funcLogCollectionOnce sync.Once
)

func getFuncLogCollection() *mongo.Collection {
	funcLogCollectionOnce.Do(func() {
		funcLogCollection = getMongoConn().Collection(config.GetCollectionConfig().Get(funcViewTable).(string))
	})
	return funcLogCollection
}
