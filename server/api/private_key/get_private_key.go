package private_key

import (
	"io/ioutil"
	"net/http"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

type GetPrivateKeyRequest struct {
	OwnerType string `uri:"owner_type"`
	OwnerName string `uri:"owner_name"`
	Name      string `uri:"name"`
}

func GetPrivateKey(c *gin.Context) {
	var req GetPrivateKeyRequest
	c.ShouldBindUri(&req)

	_, err := ioutil.ReadFile(PrivateKeyPath(req))
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot open file: "+PrivateKeyPath(req), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    gin.H{},
	})
}

func PrivateKeyPath(req GetPrivateKeyRequest) string {
	return PrivateKeyDir(req.OwnerType, req.OwnerName) + "/" + req.Name
}
