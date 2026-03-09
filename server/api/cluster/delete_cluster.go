package cluster

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
)

func DeleteCluster(c *gin.Context) {

	var req GetClusterRequest
	c.ShouldBindUri(&req)

	if err := os.RemoveAll(constants.GET_DATA_CLUSTER_DIR() + "/" + req.Cluster); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to remove cluster", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
	})

}
