package os_mirror

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"syscall"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func ModifyOsMirro(c *gin.Context) {
	var req GetMirrorRequest
	c.ShouldBindUri(&req)

	buf, err := c.GetRawData()

	logrus.Trace(string(buf))

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to read request", err)
		return
	}

	var os_mirror gin.H
	err = json.Unmarshal(buf, &os_mirror)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to parse request", err)
		return
	}

	inventoryYamleBytes, err := yaml.Marshal(os_mirror["inventory"])

	if err != nil {
		common.HandleError(c, http.StatusBadRequest, "failed to convert request into yaml", err)
		return
	}

	inventoryFilePath := constants.GET_DATA_MIRROR_DIR() + "/" + req.Name + "/inventory.yaml"
	inventoryFile, err := os.OpenFile(inventoryFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0655)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to open os_mirror file.", err)
		return
	}

	err = syscall.Flock(int(inventoryFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to lock os_mirror file.", err)
		return
	}

	defer inventoryFile.Close()
	defer syscall.Flock(int(inventoryFile.Fd()), syscall.LOCK_UN)

	statusFilePath := constants.GET_DATA_MIRROR_DIR() + "/" + req.Name + "/status.yaml"
	statusContent, err := ioutil.ReadFile(statusFilePath)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to read status file.", err)
		return
	}
	status := gin.H{}
	if err := yaml.Unmarshal(statusContent, status); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to parse status.", err)
		return
	}

	newStatus := os_mirror["status"].(map[string]interface{})
	newStatus["status"] = status["status"]

	os_mirror["status"] = newStatus

	newStatusContent, err := yaml.Marshal(newStatus)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to marshal status.", err)
		return
	}

	if err := ioutil.WriteFile(statusFilePath, newStatusContent, 0655); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to save status.", err)
		return
	}

	if err := ioutil.WriteFile(inventoryFilePath, inventoryYamleBytes, 0655); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to save os_mirror file", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    os_mirror,
	})
}
