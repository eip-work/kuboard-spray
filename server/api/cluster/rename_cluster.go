package cluster

import (
	"net/http"
	"os"

	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
)

type RenameClusterRequest struct {
	GetClusterRequest
	NewClusterName string `form:"newName"`
}

func RenameCluster(c *gin.Context) {
	var req RenameClusterRequest
	c.ShouldBindUri(&req)
	c.ShouldBindQuery(&req)

	clusterDir := constants.GET_DATA_CLUSTER_DIR()
	os.Rename(clusterDir+"/"+req.Cluster, clusterDir+"/"+req.NewClusterName)

	inventoryPath := cluster_common.ClusterInventoryYamlPath(req.NewClusterName)
	inventory, err := common.ParseYamlFile(inventoryPath)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "读取 Cluster 信息失败", err)
		return
	}

	common.PopulateKuboardSprayVars(inventory, "cluster", req.NewClusterName)

	if err := cluster_common.SaveInventory(req.NewClusterName, inventory); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "保存 Cluster 信息失败", err)
		return
	}

	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
	})

}
