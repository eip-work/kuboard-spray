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
		logrus.Warning("failed to read request", err)
		common.HandleError(c, http.StatusInternalServerError, "failed to read request: "+err.Error())
		return
	}

	var inventory gin.H
	err = json.Unmarshal(buf, &inventory)

	if err != nil {
		logrus.Warning("failed to parse request", err)
		common.HandleError(c, http.StatusInternalServerError, "failed to parse request: "+err.Error())
		return
	}

	inventoryYamleBytes, err := yaml.Marshal(inventory)

	if err != nil {
		logrus.Warning("failed to convert request into yaml: ", err)
		common.HandleError(c, http.StatusBadRequest, "failed to convert request into yaml: "+err.Error())
		return
	}

	inventoryPath := constants.GET_DATA_INVENTORY_DIR() + "/" + req.Cluster

	logrus.Info(inventoryPath)
	inventoryFilePath := inventoryPath + "/inventory.yaml"
	inventoryFile, err := os.OpenFile(inventoryFilePath, os.O_RDWR|os.O_TRUNC, 0777)

	if err != nil {
		logrus.Warning("failed to open inventory file: "+inventoryFilePath, err)
		common.HandleError(c, http.StatusInternalServerError, "failed to open inventory file: "+inventoryFilePath)
		return
	}

	err = syscall.Flock(int(inventoryFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)

	if err != nil {
		logrus.Warning("failed to lock inventory file: "+inventoryFilePath, err)
		common.HandleError(c, http.StatusInternalServerError, "failed to lock inventory file: "+inventoryFilePath)
		return
	}

	defer inventoryFile.Close()
	defer syscall.Flock(int(inventoryFile.Fd()), syscall.LOCK_UN)

	_, err = inventoryFile.Write(inventoryYamleBytes)
	if err != nil {
		logrus.Warning("failed to save inventory file: "+inventoryFilePath, err)
		common.HandleError(c, http.StatusInternalServerError, "failed to save inventory file: "+inventoryFilePath)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    inventory,
	})

}
