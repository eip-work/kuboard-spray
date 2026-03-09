package cluster_access

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/api/cluster"
	"github.com/opencmit/pangee-cluster/api/cluster/cluster_common"
	"github.com/opencmit/pangee-cluster/api/command"
	"github.com/opencmit/pangee-cluster/common"
)

func GetKubeConfig(c *gin.Context) {
	var req cluster.GetClusterRequest
	c.ShouldBindUri(&req)

	inventoryYamlPath := cluster_common.ClusterInventoryYamlPath(req.Cluster)

	args := []string{"kube_control_plane[0]", "-m", "shell", "-a", "cat /root/.kube/config"}
	result, err := command.ExecuteAdhocCommandWithInventory(inventoryYamlPath, args)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot fetch kubeconfig", err)
		return
	}

	c.JSON(http.StatusOK, common.PangeeClusterResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    result[0],
	})
}
