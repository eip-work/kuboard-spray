package login

import (
	"crypto/md5"
	"fmt"
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
		common.HandleError(c, http.StatusBadRequest, "Wrong parameters.", err)
		return
	}

	userInfo, err := readUserRepository(user.Username)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "", err)
		return
	}
	m := md5.Sum([]byte(user.Password))
	md5str1 := fmt.Sprintf("%x", m)
	if userInfo.Password == md5str1 {
		tokenString, _ := GenerateToken(user.Username)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "success",
			"data":    gin.H{"token": tokenString},
		})
		return
	}

	common.HandleError(c, http.StatusUnauthorized, "Failed to authorize.", err)

}
