package middleware

import (
	"fmt"
	"log"
	"replite_web/internal/app/dao"
	"replite_web/internal/app/utils"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// const DEFAULT_RBAC_CONFIG = "../../configs/rbac.json"

//TODO match the rule to allow user to use

// all user can access the resource list
// var currencyResource = map[string]any{
// 	"login":    nil,
// 	"register": nil,
// }

/*
the middleware must open after the jwt verify
*/
func RBACMiddleware(context *gin.Context) {
	role, bol := context.Get("role")
	if !bol {
		context.AbortWithStatus(utils.AuthFailedState)
		return
	}
	resource, bol := context.Get("resource")
	if !bol {
		context.AbortWithStatus(utils.AuthFailedState)
		return
	}
	// currency resource skip the rbac
	// if _, ok := currencyResource[resource.(string)]; ok {
	// 	return
	// }
	if !hasAuthority(role.(string), resource.(string)) {
		context.AbortWithStatus(utils.AuthFailedState)
		return
	}
}

// key: the Authentication level  value: what the level want to access
func hasAuthority(key string, value string) bool {
	rw.RLock()
	defer rw.RUnlock()
	_, ok := systemSource[key][value]
	return ok
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

func getRulesToMap() {
	rw.Lock()
	defer rw.Unlock()
	rules, err := dao.QueryAllRules()
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

// func Authorization(context *gin.Context) {

// rbac, err := grbac.New(grbac.WithJSON(DEFAULT_RBAC_CONFIG, 10*time.Minute))
// if err != nil {
// 	panic(fmt.Sprint("rbac config error:", err.Error()))
// }

// roles := queryRolesForUser(context)
// state, err := rbac.IsRequestGranted(context.Request, roles)
// if err != nil || !state.IsGranted() {
// 	context.Writer.Write(utils.Not_Granted_Error(context.Request))
// 	context.Abort()
// 	return
// }
// return
// }
