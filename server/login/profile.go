package login

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	username, exists := c.Get("username")
	if exists {
		c.JSON(http.StatusOK, common.KuboardSprayResponse{
			Code:    http.StatusOK,
			Message: "success",
			Data: UserInfo{
				Username: username.(string),
			},
		})
	} else {
		common.HandleError(c, http.StatusInternalServerError, "you are not login", nil)
		return
	}
}

type ChangePasswordRequest struct {
	Password  string `json:"password"`
	Password2 string `json:"password2"`
}

func ChangePassword(c *gin.Context) {

	var req ChangePasswordRequest

	c.ShouldBindJSON(&req)

	username, exists := c.Get("username")
	if exists {

		userRepository, err := common.ParseYamlFile(userFilePath)
		if err != nil {
			common.HandleError(c, http.StatusInternalServerError, "", err)
			return
		}

		m := md5.Sum([]byte(req.Password))
		md5str1 := fmt.Sprintf("%x", m)
		common.MapSet(userRepository, username.(string)+".password", md5str1)

		if err := common.SaveYamlFile(userFilePath, userRepository); err != nil {
			common.HandleError(c, http.StatusInternalServerError, "cannot save password", err)
			return
		}

		c.JSON(http.StatusOK, common.KuboardSprayResponse{
			Code:    http.StatusOK,
			Message: "success",
			Data: UserInfo{
				Username: username.(string),
			},
		})
	} else {
		common.HandleError(c, http.StatusInternalServerError, "you are not login", nil)
		return
	}
}
