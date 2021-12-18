package private_key

import (
	"io/ioutil"
	"net/http"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

type GetPrivateKeyRequest struct {
	Cluster string `uri:"cluster"`
	Name    string `uri:"name"`
}

func GetPrivateKey(c *gin.Context) {
	var req GetPrivateKeyRequest
	c.ShouldBindUri(&req)

	_, err := ioutil.ReadFile(PrivateKeyPath(req.Cluster, req.Name))
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot open file: "+PrivateKeyPath(req.Cluster, req.Name), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    gin.H{},
	})
}

func PrivateKeyPath(cluster string, name string) string {
	return ClusterPrivateKeyPath(cluster) + "/" + name
}
