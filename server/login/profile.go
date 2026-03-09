package login

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/common"
)

func GetProfile(c *gin.Context) {
	username, exists := c.Get("username")
	if exists {
		c.JSON(http.StatusOK, common.PangeeClusterResponse{
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

		// m := md5.Sum([]byte(req.Password))
		m, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		md5str1 := fmt.Sprintf("%x", m)
		common.MapSet(userRepository, username.(string)+".password", md5str1)

		if err := common.SaveYamlFile(userFilePath, userRepository); err != nil {
			common.HandleError(c, http.StatusInternalServerError, "cannot save password", err)
			return
		}

		c.JSON(http.StatusOK, common.PangeeClusterResponse{
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
