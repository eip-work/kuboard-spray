package common

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HandleError(c *gin.Context, code int, message string, err error) {
	if err != nil {
		logrus.Warn(message, err)
		c.JSON(code, gin.H{
			"code":    code,
			"message": message + " : " + err.Error(),
		})
		return
	}
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
