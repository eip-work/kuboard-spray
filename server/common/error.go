package common

import (
	"os"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var printstack = os.Getenv("GIN_MODE") != "release"

func HandleError(c *gin.Context, code int, message string, err error) {
	if err != nil {
		logrus.Warning(message, err)
		if printstack {
			debug.PrintStack()
		}
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
