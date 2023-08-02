package dao

import (
	"context"
	"fmt"
	"replite_web/internal/app/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var mapfuncCache map[string]FuncMap

const funcmapTable = "funcmap"

const default_funcmap_times = 3 * time.Second

func getFuncMapCollection() *mongo.Collection {
	return getMongoConn().Collection(config.CollectionConfig.Get(funcmapTable).(string))
}

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cur, err := getFuncMapCollection().Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	funcmapSlice := make([]FuncMap, 10)
	err = cur.All(ctx, funcmapSlice)
	if err != nil {
		panic(fmt.Sprintf("analysis the funcmap collection error:%s", err.Error()))
	}
	mapfuncCache = make(map[string]FuncMap)
	// build the index about the funcmap
	for i := 0; i < len(funcmapSlice); i++ {
		var funcmap = funcmapSlice[i]
		mapfuncCache[funcmap.Function] = funcmap
	}
}

// the use function  connect with the exec file relative location
type FuncMap struct {
	Function string `json:"function" bson:"function"`
	// the function should execution command,isolated operate system
	Command string `json:"execfile" bson:"execfile"`
	Type    Core   `json:"type" bson:"type"`
	OSType  OSType `json:"osType"`
	// the additional field represent the extending environment
	Additional string `json:"additional"`
}

func CreateFuncMap(funcmap FuncMap) error {
	ctx, cancel := context.WithTimeout(context.Background(), default_funcmap_times)
	defer cancel()
	_, err := getFuncMapCollection().InsertOne(ctx, funcmap)
	if err != nil {
		return err
	}
	// add success ,update the cache
	mapfuncCache[funcmap.Function] = funcmap
	return nil
}

func DeleteFuncMap(funcmap FuncMap) error {
	ctx, cancel := context.WithTimeout(context.Background(), default_funcmap_times)
	defer cancel()
	_, err := getFuncMapCollection().DeleteOne(ctx, funcmap)
	if err != nil {
		return err
	}
	// update cache
	delete(mapfuncCache, funcmap.Function)
	return nil
}

func GetFuncMap(function string) FuncMap {
	return mapfuncCache[function]
}
