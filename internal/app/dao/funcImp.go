package dao

import (
	"context"
	"fmt"
	"log"
	"replite_web/internal/app/cache"
	"replite_web/internal/app/config"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FunctionDao struct {
}

var (
	functionDao  *FunctionDao
	functionOnce sync.Once
)

func getFunctionDao() *FunctionDao {
	fmt.Println("func")
	functionOnce.Do(func() {
		functionDao = new(FunctionDao)
	})
	return functionDao
}

const funcmapTable = "funcmap"

const default_funcmap_times = 10 * time.Second

var (
	funcMapCollection     *mongo.Collection
	funcMapCollectionOnce sync.Once
)

func getFuncMapCollection() *mongo.Collection {
	funcMapCollectionOnce.Do(func() {
		funcMapCollection = getMongoConn().Collection(config.GetCollectionConfig().Get(funcmapTable).(string))
	})
	return funcMapCollection
}

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), default_funcmap_times)
	defer cancel()
	cur, err := getFuncMapCollection().Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	funcmapSlice := make([]*FuncMap, 10)
	err = cur.All(ctx, &funcmapSlice)
	if err != nil {
		panic(fmt.Sprintf("analysis the funcmap collection error:%s", err.Error()))
	}
	// build the index about the funcmap
	for i := 0; i < len(funcmapSlice); i++ {
		var funcmap = funcmapSlice[i]
		cache.GetCachePool().Store(funcmap.Function, funcmap, 10*time.Minute)
	}
}

// the use function  connect with the exec file relative location
type FuncMap struct {
	Function string `json:"function" bson:"function"`
	// the function should execution command,isolated operate system
	Command string `json:"execfile" bson:"execfile"`
	// // the params template params
	// Params []string `json:"params" bson:"params"`
	Type   Core   `json:"type" bson:"type"`
	OSType OSType `json:"osType"`
	// the additional field represent the extending environment
	Additional string `json:"additional"`
}

func (functionDao *FunctionDao) CreateFuncMap(funcmap FuncMap) error {
	ctx, cancel := context.WithTimeout(context.Background(), default_funcmap_times)
	defer cancel()
	_, err := getFuncMapCollection().InsertOne(ctx, funcmap)
	if err != nil {
		return err
	}
	// add success ,update the cache
	cache.GetCachePool().Store(funcmap.Function, &funcmap, 10*time.Minute)
	return nil
}

func (functionDao *FunctionDao) DeleteFuncMap(funcmap FuncMap) error {
	ctx, cancel := context.WithTimeout(context.Background(), default_funcmap_times)
	defer cancel()
	_, err := getFuncMapCollection().DeleteOne(ctx, funcmap)
	if err != nil {
		return err
	}
	// update cache
	cache.GetCachePool().Delete(funcmap.Function)
	return nil
}

func (functionDao *FunctionDao) GetFuncMap(function string) (fm *FuncMap) {
	if fm, ok := cache.GetCachePool().TryGet(function); ok {
		return fm.(*FuncMap)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), default_funcmap_times)
	defer cancel()
	single := getFuncMapCollection().FindOne(ctx, bson.M{"function": function})
	if single.Err() != nil {
		log.Printf("query funcmap(%s) error:%s", function, single.Err().Error())
		return nil
	}
	err := single.Decode(&fm)
	if err != nil {
		return nil
	}
	cache.GetCachePool().Store(function, fm, 10*time.Minute)
	return fm
}
