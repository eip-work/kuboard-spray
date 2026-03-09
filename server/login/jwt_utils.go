package login

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
	"github.com/sirupsen/logrus"
)

type PangeeClusterClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 24 * 365

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

		password := os.Getenv("PANGEE_CLUSTER_DEFAULT_ADMIN_PASSWORD")

		defaultRepository := map[string]*UserInfo{
			"admin": {
				Username: "admin",
				Password: password,
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

	c := PangeeClusterClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "pangee-cluster",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(JwtSecret)
}

func ParseToken(tokenString string) (*PangeeClusterClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &PangeeClusterClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*PangeeClusterClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
