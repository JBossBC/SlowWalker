package dao

import (
	"context"
	"replite_web/internal/app/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DEFAULT_RULE_COLLECTION = "rule"

const DEFAULT_RULE_NUMBER = 3

type Rule struct {
	Name      string `json:"name" bson:"name"`
	Authority string `json:"authority" bson:"authority"`
}

func getRuleCollection() *mongo.Collection {
	return getMongoConn().Collection(config.CollectionConfig.Get(DEFAULT_RULE_COLLECTION).(string))
}

/* the method is be used by the init stage,the backend load the authority rules*/
func QueryRules(page int, pageNumber int) (rules []*Rule, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cur, err := getRuleCollection().Find(ctx, bson.M{}, options.Find().SetLimit(int64(pageNumber)), options.Find().SetSkip(int64(page-1)))
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	rules = make([]*Rule, page)
	err = cur.All(context.Background(), rules)
	return rules, err
}

func QueryAllRules() (rules []*Rule, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cur, err := getRuleCollection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	rules = make([]*Rule, 0, DEFAULT_RULE_NUMBER)
	err = cur.All(context.Background(), &rules)
	return rules, err
}
