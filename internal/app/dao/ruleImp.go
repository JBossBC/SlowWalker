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

type RuleDao struct {
}

var (
	ruleDao  *RuleDao
	ruleOnce sync.Once
)

func getRuleDao() *RuleDao {
	ruleOnce.Do(func() {
		ruleDao = new(RuleDao)
	})
	return ruleDao
}

const DEFAULT_RULE_COLLECTION = "rule"

const DEFAULT_RULE_NUMBER = 3

type RuleInfo struct {
	//the project relative url
	Name      string `json:"name" bson:"name"`
	Authority string `json:"authority" bson:"authority"`
	// Type      string `json:"type" bson:"type"`
	//the owner of rule can operate the function level
	//1. Query Level: the model can be found by this authority
	//2. Scan: the model should eval the details for this authority
	//3. Update,Create,Delete: the authority can process the operation for this model(the field has difference for different model)
	// Operation string `json:"operation" bson:"operation"`
}

const DEFAULT_RENEW_RULES_MAP_TIME = 24 * time.Hour

// begin the scheduled task to renew the systemSource
func init() {
	systemSource = make(map[string]map[string]any)
	go func() {
		// init the timer
		// timer := time.NewTimer(DEFAULT_RENEW_RULES_MAP_TIME)
		// init the systemSource
		getRulesToMap()
		// for {
		// 	select {
		// 	case <-timer.C:
		// 		timer.Reset(DEFAULT_RENEW_RULES_MAP_TIME)
		// 		//renew the systemSource
		// 		getRulesToMap()
		// 	default:
		// 		time.Sleep(1 * time.Hour)
		// 	}
		// }
	}()
}

var (
	systemSource map[string]map[string]any
	rw           sync.RWMutex
)

// ok is true represent the rule has authority
func (ruleDao *RuleDao) GetRule(owner string, authority string) (value any, ok bool) {
	rw.RLock()
	defer rw.RUnlock()
	owners, hasOwner := systemSource[owner]
	if !hasOwner {
		return nil, false
	}
	value, ok = owners[authority]
	return
}

func (ruleDao *RuleDao) GetAuthority(owner string) []any {
	var result = make([]any, 0, 3)
	for _, value := range systemSource[owner] {
		result = append(result, value)
	}
	return result
}

// func getRulesCache() *map[string]map[string]any {
// 	return &systemSource
// }

func getRulesToMap() {
	rw.Lock()
	defer rw.Unlock()
	rules, err := getRuleDao().QueryAllRules()
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
		systemSource[rule.Name][rule.Authority] = rule
	}
	log.Println("renew the rules successfully")
}

var (
	ruleCollection     *mongo.Collection
	ruleCollectionOnce sync.Once
)

func getRuleCollection() *mongo.Collection {
	ruleCollectionOnce.Do(func() {
		ruleCollection = getMongoConn().Collection(config.GetCollectionConfig().Get(DEFAULT_RULE_COLLECTION).(string))
	})
	return ruleCollection
}

/* the method is be used by the init stage,the backend load the authority rules*/
func (ruleDao *RuleDao) QueryRules(page int, pageNumber int) (rules []*RuleInfo, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cur, err := getRuleCollection().Find(ctx, bson.M{}, options.Find().SetLimit(int64(pageNumber)), options.Find().SetSkip(int64(page-1)))
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	rules = make([]*RuleInfo, page)
	err = cur.All(context.Background(), rules)
	return rules, err
}

func (ruleDao *RuleDao) QueryAllRules() (rules []*RuleInfo, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cur, err := getRuleCollection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	rules = make([]*RuleInfo, 0, DEFAULT_RULE_NUMBER)
	err = cur.All(context.Background(), &rules)
	return rules, err
}
