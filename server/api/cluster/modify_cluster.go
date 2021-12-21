package cluster

import (
	"encoding/json"
	"net/http"
	"os"
	"syscall"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
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

	inventoryYamleBytes, err := yaml.Marshal(inventory)

	if err != nil {
		common.HandleError(c, http.StatusBadRequest, "failed to convert request into yaml", err)
		return
	}

	inventoryFilePath := constants.GET_DATA_CLUSTER_DIR() + "/" + req.Cluster + "/inventory.yaml"
	inventoryFile, err := os.OpenFile(inventoryFilePath, os.O_RDWR|os.O_TRUNC, 0666)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to open inventory file: "+inventoryFilePath, err)
		return
	}

	err = syscall.Flock(int(inventoryFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to lock inventory file: "+inventoryFilePath, err)
		return
	}

	defer inventoryFile.Close()
	defer syscall.Flock(int(inventoryFile.Fd()), syscall.LOCK_UN)

	_, err = inventoryFile.Write(inventoryYamleBytes)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to save inventory file: "+inventoryFilePath, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    inventory,
	})

}
