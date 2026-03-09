package cluster

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/api/cluster/cluster_common"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
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
	err := os.Rename(clusterDir+"/"+req.Cluster, clusterDir+"/"+req.NewClusterName)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "重命名 Cluster 失败", err)
		return
	}

	inventoryPath := cluster_common.ClusterInventoryYamlPath(req.NewClusterName)
	inventory, err := common.ParseYamlFile(inventoryPath)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "读取 Cluster 信息失败", err)
		return
	}

	common.PopulatePangeeClusterVars(inventory, "cluster", req.NewClusterName)

	if err := cluster_common.SaveInventory(req.NewClusterName, inventory); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "保存 Cluster 信息失败", err)
		return
	}

	c.JSON(http.StatusOK, common.PangeeClusterResponse{
		Code:    http.StatusOK,
		Message: "success",
	})

}
