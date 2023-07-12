package utils

import (
	"github.com/golang-jwt/jwt"
)

type JwtClaims struct {
	jwt.StandardClaims
	Role     string `json:"role"`
	Username string `json:"username"`
}

var DEFUALT_JWT_METHOD = jwt.SigningMethodHS512

func CreateJWT(secret string, claims jwt.Claims) (string, error) {
	// if _, ok := claims.(jwt.Claims); !ok {
	// 	return "", errors.New("the input params cant correspond the jwt.StandardClaims")
	// }
	token := jwt.NewWithClaims(DEFUALT_JWT_METHOD, claims)
	//return the jwt string
	return token.SignedString([]byte(secret))
}
