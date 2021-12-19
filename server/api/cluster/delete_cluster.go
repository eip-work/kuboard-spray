package cluster

import (
	"net/http"
	"os"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
)

func DeleteCluster(c *gin.Context) {

	var req GetClusterRequest
	c.ShouldBindUri(&req)

	if err := os.RemoveAll(constants.GET_DATA_INVENTORY_DIR() + "/" + req.Cluster); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to remove cluster", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
	})

}
