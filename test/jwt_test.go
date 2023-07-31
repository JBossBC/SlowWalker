package test

import (
	"replite_web/internal/app/utils"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

// pass
func TestJWT(t *testing.T) {
	Claim := utils.JwtClaims{
		Username: "xiyang",
		Role:     "admin",
	}
	result, err := utils.CreateJWT("oParplZS7iTFisR6VLXGG1_4fPDpQo2qjQiH4By7wehSMhgSUM8OYFMuZ4kWi9ETVpA5K6BhWGoJdqq2uT8uTQ", Claim, time.Now().Add(3*time.Second))
	if err != nil {
		t.Fatalf("创建jwt出错:%s", err.Error())
	}
	var cliams = new(utils.JwtClaims)
	token, err := jwt.ParseWithClaims(result, cliams, func(token *jwt.Token) (interface{}, error) {
		return []byte("oParplZS7iTFisR6VLXGG1_4fPDpQo2qjQiH4By7wehSMhgSUM8OYFMuZ4kWi9ETVpA5K6BhWGoJdqq2uT8uTQ"), nil
	})
	if err != nil {
		t.Fatalf("测试utils.JWT失败:%s", err.Error())
	}

	// t.Logf(cliams.Role)
	//检查JWTToken是否已过期
	claims, ok := token.Claims.(*utils.JwtClaims)
	if !ok || !token.Valid {
		t.Fatalf("Token无效")
	}
	if claims.ExpiresAt < time.Now().Unix() {
		t.Fatalf("Token过期")
	}
}
