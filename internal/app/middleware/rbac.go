package middleware

import (
	"fmt"
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

/*
the middleware must open after the jwt verify
*/
func RBACMiddleware(context *gin.Context) {
	role, bol := context.Get("role")
	fmt.Println("role=", role)
	if !bol {
		fmt.Println("错误1")
		context.AbortWithStatus(utils.AuthFailedState)
		return
	}
	resource, bol := context.Get("resource")
	fmt.Println("resource=", resource)
	if !bol {
		fmt.Println("错误2")
		context.AbortWithStatus(utils.AuthFailedState)
		return
	}
	// currency resource skip the rbac
	// if _, ok := currencyResource[resource.(string)]; ok {
	// 	return
	// }

	if !hasAuthority(role.(string), resource.(string)) {
		context.AbortWithStatus(utils.AuthFailedState)
		fmt.Println("错误3")
		return
	}

	fmt.Println("中间件3没有出差错")

}

// key: the Authentication level  value: what the level want to access
func hasAuthority(key string, value string) bool {
	// rw.RLock()
	// defer rw.RUnlock()
	// _, ok := systemSource[key][value]
	_, ok := dao.GetRule(key, value)
	fmt.Println(ok, 99999)
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
