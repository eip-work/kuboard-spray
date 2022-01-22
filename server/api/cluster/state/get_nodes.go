package state

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

type GetNodesRequest struct {
	ClusterName string `uri:"cluster"`
}

func GetNodes(c *gin.Context) {

	var request GetNodesRequest
	c.ShouldBindUri(&request)

	result, err := cluster_common.ExecuteShellOnControlPlane(request.ClusterName, "kubectl get nodes -o json")

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
