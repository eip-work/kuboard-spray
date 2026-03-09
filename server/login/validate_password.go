package login

import (
	"net/http"

	"encoding/hex"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/common"
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
	result := "failed"

	hashedPassword, err := hex.DecodeString(userInfo.Password)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Cannot decode password in user repo.", err)
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(user.Password))
	if err == nil {
		result = "ok"
	} else {
		common.HandleError(c, http.StatusBadRequest, "Wrong password.", err)
		return
	}
	result = "ok"
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    result,
	})

}
