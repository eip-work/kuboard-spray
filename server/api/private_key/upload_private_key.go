package private_key

import (
	"net/http"
	"os"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func UploadPrivateKeyFile(c *gin.Context) {

	var req GetPrivateKeyRequest
	c.ShouldBindUri(&req)

	file, _ := c.FormFile("file")
	logrus.Info(file.Filename)
	c.SaveUploadedFile(file, PrivateKeyPath(req.Cluster, req.Name))
	err := os.Chmod(PrivateKeyPath(req.Cluster, req.Name), 0600)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "chmod faild: ", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
	})
}
