package health_check

import (
	"net/http"
	"strconv"

	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

type CheckConnectivityRequest struct {
	ClusterName string `uri:"cluster"`
}

func CheckConnectivity(c *gin.Context) {
	var request CheckConnectivityRequest
	c.ShouldBindUri(&request)

	inventoryPath := cluster_common.ClusterInventoryYamlPath(request.ClusterName)
	inventory, err := common.ParseYamlFile(inventoryPath)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to parse inventory", err)
		return
	}

	netcheckerPort := 31081
	np := common.MapGet(inventory, "all.children.target.children.k8s_cluster.vars.netchecker_port")
	if np != nil {
		netcheckerPort = np.(int)
	}

	result, err := cluster_common.ExecuteShellOnControlPlane(request.ClusterName, "curl -s http://$(hostname):"+strconv.Itoa(netcheckerPort)+"/api/v1/connectivity_check")

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
