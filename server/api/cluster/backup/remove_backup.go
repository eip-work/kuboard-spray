package backup

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
	"github.com/sirupsen/logrus"
)

type RemoveBackupRequest struct {
	Cluster string   `uri:"cluster"`
	Backups []string `json:"backups_to_remove"`
}

func RemoveBackup(c *gin.Context) {

	var req RemoveBackupRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)

	for _, backup := range req.Backups {
		path := constants.GET_DATA_CLUSTER_DIR() + "/" + req.Cluster + "/backup/" + backup
		if err := os.Remove(path); err != nil {
			logrus.Warn("cannot remove backup.", err)
		}
	}

	c.JSON(http.StatusOK, common.PangeeClusterResponse{
		Code:    http.StatusOK,
		Message: "success",
	})
}
