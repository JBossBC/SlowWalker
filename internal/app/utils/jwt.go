package utils

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

type JwtClaims struct {
	jwt.StandardClaims
	Role     string `json:"role"`
	Username string `json:"username"`
}

var DEFUALT_JWT_METHOD = jwt.SigningMethodHS512

func CreateJWT(secret string, claims jwt.Claims, expirationTime time.Time) (string, error) {
	// if _, ok := claims.(jwt.Claims); !ok {
	// 	return "", errors.New("the input params cant correspond the jwt.StandardClaims")
	// }
	if c, ok := claims.(*JwtClaims); ok {
		c.ExpiresAt = expirationTime.Unix()
	} else {
		return "", errors.New("the input params cant correspond to JwtClaims")
	}
	token := jwt.NewWithClaims(DEFUALT_JWT_METHOD, claims)
	//return the jwt string
	return token.SignedString([]byte(secret))
}
