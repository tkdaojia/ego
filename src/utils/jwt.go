package utils

import (
	"ego/src/boot/global"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTClaims struct { // token里面添加用户信息，验证token后可能会用到用户信息
	jwt.StandardClaims
	UserID   int      `json:"user_id"`
	Account  string   `json:"account"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

var (
	ExpireTime = 3600 * 24 // token有效期
)

func CreateJWT(id int, account string, username string, roles []string) string {
	claims := &JWTClaims{
		UserID:   id,
		Account:  account,
		Username: username,
		Roles:    roles,
	}
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(ExpireTime)).Unix()
	return getToken(claims)
}

func getToken(claims *JWTClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	Secret := global.C_CONFIG.System.Appid
	signedToken, err := token.SignedString([]byte(Secret))
	if err != nil {
		return ""
	}
	return signedToken
}

func DecodeJWT(strToken string) *JWTClaims {
	Secret := global.C_CONFIG.System.Appid
	token, err := jwt.ParseWithClaims(strToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		return nil
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil
	}
	if err := token.Claims.Valid(); err != nil {
		return nil
	}
	return claims
}
