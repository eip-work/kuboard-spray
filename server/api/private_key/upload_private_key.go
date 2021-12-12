package private_key

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func UploadPrivateKeyFile(c *gin.Context) {

	var req GetPrivateKeyRequest
	c.ShouldBindUri(&req)

	file, _ := c.FormFile("file")
	logrus.Info(file.Filename)
	c.SaveUploadedFile(file, PrivateKeyPatch(req.Cluster, req.Name))
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
	})
}
