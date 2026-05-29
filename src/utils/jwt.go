package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTClaims struct {
	UserID   int      `json:"user_id"`
	Account  string   `json:"account"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

var (
	Secret     = "Hello EGO" // JWT 加盐
	ExpireTime = 3600 * 24   // token 有效期（秒）
)

func CreateJWT(id int, account string, username string, roles []string) string {
	now := time.Now()
	claims := &JWTClaims{
		UserID:   id,
		Account:  account,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(ExpireTime))),
		},
	}
	return getToken(claims)
}

func getToken(claims *JWTClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(Secret))
	if err != nil {
		return ""
	}
	return signedToken
}

func DecodeJWT(strToken string) *JWTClaims {
	token, err := jwt.ParseWithClaims(strToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 🔴 安全防线：v5 依然需要强校验签名算法，防止 none 攻击
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(Secret), nil
	})
	if err != nil || !token.Valid {
		return nil
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil
	}
	return claims
}
