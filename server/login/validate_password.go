package login

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

func ValidatePassword(c *gin.Context) {
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
	result := "failed"
	if userInfo.Password == md5str1 {
		result = "ok"
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    result,
	})

}
