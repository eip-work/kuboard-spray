package login

import (
	"net/http"

	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/common"
	"golang.org/x/crypto/bcrypt"
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
	// m := md5.Sum([]byte(user.Password))
	// m, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// md5str1 := fmt.Sprintf("%x", m)
	// if userInfo.Password == md5str1 {
	hashedPassword, err := hex.DecodeString(userInfo.Password)
	if err != nil {
		common.HandleError(c, http.StatusUnauthorized, "Failed to authorize.", err)
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(user.Password))
	if err == nil {
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
