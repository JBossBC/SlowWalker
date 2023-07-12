package middleware

import (
	"errors"
	"replite_web/internal/app/config"
	"replite_web/internal/app/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// const SECRET_KEY = ""

func Auth(ctx *gin.Context) {
	rawJWT := ctx.Request.Header.Get("Authorization")
	if rawJWT == "" {
		ctx.AbortWithStatus(utils.AuthFailedState)
		return
	}
	//bearer 方案
	_, rawJWT, ok := strings.Cut(rawJWT, " ")
	if !ok {
		ctx.AbortWithStatus(utils.AuthFailedState)
		return
	}
	var auth = new(utils.JwtClaims)
	token, err := jwt.ParseWithClaims(rawJWT, auth, func(t *jwt.Token) (interface{}, error) {
		if t.Method != utils.DEFUALT_JWT_METHOD {
			return nil, errors.New("验证失败")
		}
		return []byte(config.ServerConf.Secret), nil
	})
	if err != nil || !token.Valid {
		ctx.AbortWithStatus(utils.AuthFailedState)
		return
	}
	// add the jwt claims to context
	ctx.Set("username", auth.Username)
	ctx.Set("role", auth.Role)
}
