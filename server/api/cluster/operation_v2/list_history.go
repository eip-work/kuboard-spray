package operation_v2

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
	"github.com/sirupsen/logrus"
)

type ListExecuteHistoryRequest struct {
	Cluster   string `uri:"cluster" binding:"required"`
	Operation string `uri:"operation" binding:"required"`
	Step      string `uri:"step" binding:"required"`
}

type HistoryItem struct {
	Time     string `json:"time" binding:"required"`
	Status   string `json:"status" binding:"required"`
	Duration int    `json:"duration"`
}

func ListOperationStepHistory(c *gin.Context) {
	var req ListExecuteHistoryRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)

	historyDir := constants.GET_DATA_CLUSTER_DIR() + "/" + req.Cluster + "/history/" + req.Operation + "/" + req.Step
	logrus.Info(historyDir)

	result := []HistoryItem{}
	if common.PathExists(historyDir) {
		entries, err := os.ReadDir(historyDir)
		if err != nil {
			common.HandleError(c, http.StatusInternalServerError, "cannot read dir: "+historyDir, err)
			return
		}
		for _, entry := range entries {
			if entry.IsDir() {
				item := HistoryItem{
					Time: entry.Name(),
				}
				status, err := common.ParseYamlFile(historyDir + "/" + entry.Name() + "/status.yaml")
				if err == nil {
					item.Status = status["status"].(string)
					if status["duration"] != nil {
						item.Duration = status["duration"].(int)
					}
				}
				result = append(result, item)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"history": result,
		},
	})

}
