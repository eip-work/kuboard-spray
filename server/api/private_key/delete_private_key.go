package private_key

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

func DeletePrivateKey(c *gin.Context) {
	var req GetPrivateKeyRequest
	c.ShouldBindUri(&req)

	_, err := ioutil.ReadFile(PrivateKeyPath(req.Cluster, req.Name))
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot open file: "+PrivateKeyPath(req.Cluster, req.Name), err)
		return
	}

	e := os.Remove(PrivateKeyPath(req.Cluster, req.Name))
	if e != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot delete file: "+PrivateKeyPath(req.Cluster, req.Name), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
	})
}
