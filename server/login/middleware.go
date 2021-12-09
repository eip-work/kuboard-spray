package login

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Token 放在 Header 的 Authorization 中，并使用 Bearer 开头
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Authorization header can't be empty",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "Authorization header is in bad format.",
			})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Invalid authorization token.",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.Username)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
