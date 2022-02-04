package cluster_access

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/api/ansible_rpc"
	"github.com/eip-work/kuboard-spray/api/cluster"
	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

func GetKubeConfig(c *gin.Context) {
	var req cluster.GetClusterRequest
	c.ShouldBindUri(&req)

	inventoryYamlPath := cluster_common.ClusterInventoryYamlPath(req.Cluster)

	args := []string{"kube_control_plane[0]", "-m", "shell", "-a", "cat /root/.kube/config"}
	result, err := ansible_rpc.ExecuteAdhocCommandWithInventory(inventoryYamlPath, args)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot fetch kubeconfig", err)
		return
	}

	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    result[0],
	})
}
