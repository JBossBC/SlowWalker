package middleware

import (
	"errors"
	"fmt"
	"log"
	"replite_web/internal/app/config"
	"replite_web/internal/app/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// const SECRET_KEY = ""

func Auth(ctx *gin.Context) {
	rawJWT := ctx.Request.Header.Get("Authorization")

	fmt.Println("中间件2开始")

	/*
		map[Accept:[application/json, text/plain, *] Accept-Encoding:[gzip, deflate, br] Accept-Language:[zh-CN,zh;q=0.9] Connection:[keep-alive] Origin:[http://localhost:3000] Referer:[http://localhost:3000/main
		] Sec-Ch-Ua:["Not-A.Brand";v="24", "Chromium";v="14"] Sec-Ch-Ua-Mobile:[?0] Sec-Ch-Ua-Platform:["Windows"] Sec-Fetch-Dest:[empty] Sec-Fetch-Mode:[cors] Sec-Fetch-Site:[same-site] User-Agent:[Mozilla/5.0 (Wi
		ndows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.5359.125 Safari/537.36]]
	*/
	if rawJWT == "" {
		ctx.AbortWithStatus(utils.AuthFailedState)
		fmt.Println("1")
		return
	}
	//bearer 方案
	_, rawJWT, ok := strings.Cut(rawJWT, "\u0020")
	if !ok {
		ctx.AbortWithStatus(utils.AuthFailedState)
		fmt.Println("2")
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
		log.Printf("%v", err)
		ctx.AbortWithStatus(utils.AuthFailedState)
		return
	}
	// add the jwt claims to context

	fmt.Println("auth=", auth)
	fmt.Println(auth.Username)
	fmt.Println(auth.Role)

	ctx.Set("resource", auth.Username) //原来写的是ctx.Set("username", auth.Username)
	ctx.Set("role", auth.Role)
	fmt.Println(auth.Role)

	fmt.Println("中间件2没有出差错")

}
