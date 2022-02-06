package login

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/sirupsen/logrus"
)

type KuboardSprayClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 24

var JwtSecret = []byte("mytoken")

var userDir = constants.GET_DATA_DIR() + "/user"
var userFilePath = userDir + "/user.yaml"

func prepareJwtSecret() {
	secretFilePath := userDir + "/jwt_secret"
	JwtSecret, err := os.ReadFile(secretFilePath)
	if err != nil {
		secret := common.RandomString(32)
		JwtSecret = []byte(secret)
		if err := common.CreateDirIfNotExists(userDir); err != nil {
			logrus.Fatal("cannot create user dir", err)
			os.Exit(1)
		}
		if err := os.WriteFile(secretFilePath, JwtSecret, 0655); err != nil {
			logrus.Fatal("cannot generate jwt secret", err)
			os.Exit(1)
		}
	}
	logrus.Trace("JwtSecret:", string(JwtSecret))
}

func readUserRepository(username string) (*UserInfo, error) {
	userRepository, err := common.ParseYamlFile(userFilePath)
	if err != nil {
		if err := common.CreateDirIfNotExists(userDir); err != nil {
			logrus.Fatal("cannot create user dir", err)
			os.Exit(1)
		}
		defaultRepository := map[string]*UserInfo{
			"admin": {
				Username: "admin",
				Password: "4df59565aa2497c16ac4f49a073ee318",
			},
		}
		if err := common.SaveYamlFile(userFilePath, defaultRepository); err != nil {
			return nil, err
		}
		if username == "admin" {
			return defaultRepository["admin"], nil
		}
	}
	if userRepository[username] != nil {
		user := userRepository[username].(map[string]interface{})
		userInfo := &UserInfo{
			Username: username,
			Password: user["password"].(string),
		}
		return userInfo, nil
	}
	return nil, errors.New("username / password doesnot match")
}

func GenerateToken(username string) (string, error) {

	prepareJwtSecret()

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
