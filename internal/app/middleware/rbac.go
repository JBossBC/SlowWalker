package middleware

import (
	"replite_web/internal/app/dao"
	"replite_web/internal/app/utils"

	"github.com/gin-gonic/gin"
)

// const DEFAULT_RBAC_CONFIG = "../../configs/rbac.json"

//TODO match the rule to allow user to use

// all user can access the resource list
// var currencyResource = map[string]any{
// 	"login":    nil,
// 	"register": nil,
// }

var validRule map[string]any = map[string]any{"admin": nil, "member": nil, "audit": nil}

/*
the middleware must open after the jwt verify
*/
func RBACMiddleware(context *gin.Context) {
	role, bol := context.Get("role")
	if _, ok := validRule[role.(string)]; !bol || !ok {
		context.AbortWithStatus(utils.AuthFailedState)
		return
	}
	// resource, bol := context.Get("resource")
	// if !bol {
	// 	context.AbortWithStatus(utils.AuthFailedState)
	// 	return
	// }
	// currency resource skip the rbac
	// if _, ok := currencyResource[resource.(string)]; ok {
	// 	return
	// }
	resource := context.Request.URL.Path
	if !hasAuthority(role.(string), resource) {
		context.AbortWithStatus(utils.AuthFailedState)
		return
	}
}

// key: the Authentication level  value: what the level want to access
func hasAuthority(key string, value string) bool {
	// rw.RLock()
	// defer rw.RUnlock()
	// _, ok := systemSource[key][value]
	_, ok := dao.GetRuleDao().GetRule(key, value)
	return ok
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
