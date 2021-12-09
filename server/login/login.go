package login

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AuthHandler(c *gin.Context) {
	var user UserInfo
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Wrong parameters.",
		})
		return
	}

	if user.Username == "admin" && user.Password == "123456" {
		tokenString, _ := GenerateToken(user.Username)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "success",
			"data":    gin.H{"token": tokenString},
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "Failed to authorize."})
}
