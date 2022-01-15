package state

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

func GetEtcdNodes(c *gin.Context) {
	var request GetNodesRequest
	c.ShouldBindUri(&request)

	result, err := ExecuteShellOnETCD(request.ClusterName, "etcdctl member list --write-out=json")

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed", err)
		return
	}

	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    result,
	})
}
