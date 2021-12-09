package login

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/common"
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
		common.HandleError(c, http.StatusBadRequest, "Wrong parameters.")
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

	common.HandleError(c, http.StatusUnauthorized, "Failed to authorize.")

}
