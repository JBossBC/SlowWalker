package middleware

// var default_get_params = map[string]any{
// 	"resource": nil,
// }

//TODO abandoned,because this operation will consume the huge performance
// func BeforeHandler(ctx *gin.Context) {
// 	switch ctx.Request.Method {
// 	case http.MethodPost:
// 		for k := range default_get_params {
// 			value := ctx.PostForm(k)
// 			if value != "" {
// 				ctx.Set(k, value)
// 			}
// 		}
// 	case http.MethodGet:
// 		for k := range default_get_params {
// 			value := ctx.Query(k)
// 			if value != "" {
// 				ctx.Set(k, value)
// 			}
// 		}
// 	}
// }

// RBAC  /context
