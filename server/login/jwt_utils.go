package login

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type KuboardSprayClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 24

var JwtSecret = []byte("mytoken")

func GenerateToken(username string) (string, error) {
	c := KuboardSprayClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "kuboard-spray",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(JwtSecret)
}

func ParseToken(tokenString string) (*KuboardSprayClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &KuboardSprayClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*KuboardSprayClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
