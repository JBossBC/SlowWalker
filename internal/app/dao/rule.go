package dao

import (
	"context"
	"fmt"
	"log"
	"replite_web/internal/app/config"
	"sync"
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

const DEFAULT_RENEW_RULES_MAP_TIME = 24 * time.Hour

// begin the scheduled task to renew the systemSource
func init() {
	systemSource = make(map[string]map[string]any)
	go func() {
		// init the timer
		timer := time.NewTimer(DEFAULT_RENEW_RULES_MAP_TIME)
		// init the systemSource
		getRulesToMap()
		for {
			select {
			case <-timer.C:
				timer.Reset(DEFAULT_RENEW_RULES_MAP_TIME)
				//renew the systemSource
				getRulesToMap()
			default:
				time.Sleep(10 * time.Minute)
			}
		}
	}()
}

var (
	systemSource map[string]map[string]any
	rw           sync.RWMutex
)

func GetRule(owner string, authority string) (value any, ok bool) {
	rw.RLock()
	defer rw.RUnlock()
	owners, hasOwner := systemSource[owner]
	if !hasOwner {
		return nil, false
	}
	value, ok = owners[authority]
	return
}

// func getRulesCache() *map[string]map[string]any {
// 	return &systemSource
// }

func getRulesToMap() {
	rw.Lock()
	defer rw.Unlock()
	rules, err := QueryAllRules()
	if err != nil {
		panic(fmt.Sprintf("the rules(%v) init failed,please inspect the error:%s", rules, err.Error()))
	}
	//clear all key value for map
	for k, v := range systemSource {
		for secondK := range v {
			delete(v, secondK)
		}
		delete(systemSource, k)
	}
	for i := 0; i < len(rules); i++ {
		rule := rules[i]
		if systemSource[rule.Name] == nil {
			systemSource[rule.Name] = make(map[string]any)
		}
		log.Printf("正在添加 rule:%s authority:%s", rule.Name, rule.Authority)
		systemSource[rule.Name][rule.Authority] = nil
	}
	log.Println("renew the rules successfully")
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
