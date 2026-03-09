package cluster

import (
	"encoding/json"
	"net/http"
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func ModifyCluster(c *gin.Context) {
	var req GetClusterRequest
	c.ShouldBindUri(&req)

	buf, err := c.GetRawData()

	logrus.Trace(string(buf))

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to read request", err)
		return
	}

	var inventory gin.H
	err = json.Unmarshal(buf, &inventory)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to parse request", err)
		return
	}

	common.PopulatePangeeClusterVars(inventory, "cluster", req.Cluster)

	inventoryYamleBytes, err := yaml.Marshal(inventory)

	if err != nil {
		common.HandleError(c, http.StatusBadRequest, "failed to convert request into yaml", err)
		return
	}

	inventoryFilePath := constants.GET_DATA_CLUSTER_DIR() + "/" + req.Cluster + "/inventory.yaml"
	inventoryFile, err := os.Create(inventoryFilePath)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to open inventory file.", err)
		return
	}

	err = syscall.Flock(int(inventoryFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to lock inventory file.", err)
		return
	}

	defer inventoryFile.Close()
	defer syscall.Flock(int(inventoryFile.Fd()), syscall.LOCK_UN)

	err = inventoryFile.Truncate(0)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to Truncate inventory file.", err)
		return
	}
	_, err = inventoryFile.Write(inventoryYamleBytes)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to save inventory file.", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    inventory,
	})

}
