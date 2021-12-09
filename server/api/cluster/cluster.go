package cluster

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListClusters(c *gin.Context) {
	username := c.MustGet("username").(string)
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    gin.H{"username": username},
	})
}
