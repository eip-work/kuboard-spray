package common

import "github.com/gin-gonic/gin"

func HandleError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
