package cluster

import (
	"io/ioutil"
	"net/http"

	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
)

type GetClusterRequest struct {
	Cluster string `uri:"cluster"`
}

func GetCluster(c *gin.Context) {
	var req GetClusterRequest
	c.ShouldBindUri(&req)

	inventory, err := common.ParseYamlFile(ClusterInventoryPath(req.Cluster))
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot parse file: "+ClusterInventoryPath(req.Cluster), err)
		return
	}

	successTasks, err := command.ReadSuccessTasks(req.Cluster)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot read cluster status", err)
		return
	}

	processing := false
	lockFile, err := command.LockOwner("cluster", req.Cluster)
	if err != nil {
		processing = true
	}
	lockFilePath := constants.GET_DATA_CLUSTER_DIR() + "/" + req.Cluster + "/inventory.lastrun"
	pid, err := ioutil.ReadFile(lockFilePath)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot read current pid", err)
		return
	}
	command.UnlockOwner(lockFile)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"inventory":     inventory,
			"name":          req.Cluster,
			"success_tasks": successTasks,
			"processing":    processing,
			"current_pid":   string(pid),
		},
	})
}

func ClusterInventoryPath(cluster string) string {
	return constants.GET_DATA_CLUSTER_DIR() + "/" + cluster + "/inventory.yaml"
}
