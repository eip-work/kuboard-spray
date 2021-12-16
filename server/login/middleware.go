package login

import (
	"net/http"
	"strings"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {

		// Token 放在 Header 的 Authorization 中，并使用 Bearer 开头
		authHeader := c.Request.Header.Get("Authorization")
		var token string
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if !(len(parts) == 2 && parts[0] == "Bearer") {
				common.HandleError(c, http.StatusBadRequest, "Authorization header is in bad format.", nil)
				return
			}
			token = parts[1]
		} else { // Token 也可以放在 Cookie 中
			var err error
			token, err = c.Cookie("KuboardToken")
			if err != nil {
				common.HandleError(c, http.StatusBadRequest, "Authorization header or Cookie KuboardToken can't be empty", err)
				return
			}
		}

		logrus.Trace(token)
		mc, err := ParseToken(token)
		if err != nil {
			common.HandleError(c, http.StatusUnauthorized, "Invalid authorization token.", err)
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.Username)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
