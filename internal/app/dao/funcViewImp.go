package dao

import (
	"context"
	"log"
	"replite_web/internal/app/config"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FuncViewDao struct {
}

var (
	funcviewDao  *FuncViewDao
	funcviewOnce sync.Once
)

func getFuncViewDao() *FuncViewDao {
	funcviewOnce.Do(func() {
		funcviewDao = new(FuncViewDao)
	})
	return funcviewDao
}

const funcViewTable = "funcmap"

const default_funcview_times = 3 * time.Second

type FuncViewInfo struct {
	Function    string `json:"function" bson:"function"`
	View        string `json:"view" bson:"view"`
	Params      string `json:"params" bson:"params"`
	Sign        bool   `json:"sign" bson:"sign"`
	EmptyPrefix bool   `json:"emptyPrefix" bson:"emptyPrefix"`
	IsMedium    bool   `json:"isMedium" bson:"isMedium"`
}

func getFuncViewCollection() *mongo.Collection {
	return getMongoConn().Collection(config.CollectionConfig.Get(funcViewTable).(string))
}

func (funcviewDao *FuncViewDao) CreateFuncViews(funcs ...FuncViewInfo) error {
	ctx, cancel := context.WithTimeout(context.TODO(), default_funcview_times)
	defer cancel()
	//TODO has more effictive function
	var replicas = make([]interface{}, len(funcs))
	for i := 0; i < len(funcs); i++ {
		replicas[i] = funcs[i]
	}
	_, err := getFuncViewCollection().InsertMany(ctx, replicas)
	return err
}

func (funcviewDao *FuncViewDao) GetFuncViews(function string) ([]*FuncViewInfo, error) {
	rs := make([]*FuncViewInfo, 0, 3)
	ctx, cancel := context.WithTimeout(context.TODO(), default_funcview_times)
	defer cancel()
	many, err := getFuncViewCollection().Find(ctx, bson.M{"function": function})
	if err != nil {
		log.Printf("查询 function(%s) 出错:%s", function, err.Error())
		return nil, err
	}
	err = many.All(ctx, rs)
	if err != nil {
		log.Printf("解析funcview collection(%v) 出错:%s", many, err.Error())
	}
	return rs, nil
}
