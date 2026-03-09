package cluster

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/api/cluster/cluster_common"
	"github.com/opencmit/pangee-cluster/api/command"
	"github.com/opencmit/pangee-cluster/common"
)

type GetClusterRequest struct {
	Cluster string `uri:"cluster"`
}

func GetCluster(c *gin.Context) {
	var req GetClusterRequest
	c.ShouldBindUri(&req)

	inventory, err := common.ParseYamlFile(cluster_common.ClusterInventoryYamlPath(req.Cluster))
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot parse file: "+cluster_common.ClusterInventoryYamlPath(req.Cluster), err)
		return
	}

	history, err := command.ReadTaskHistory("cluster", req.Cluster)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot read task history.", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"inventory": inventory,
			"history":   history,
			"name":      req.Cluster,
		},
	})
}
