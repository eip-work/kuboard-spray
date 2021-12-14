package cluster

import (
	"net/http"
	"os"
	"os/exec"
	"syscall"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ApplyPlan(c *gin.Context) {

	var req GetClusterRequest
	c.ShouldBindUri(&req)

	// create lock.

	lockFilePath := constants.GET_DATA_INVENTORY_DIR() + "/" + req.Cluster + "/inventory.lock"
	logrus.Trace("lockFilePath: ", lockFilePath)
	lockFile, err := os.OpenFile(lockFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to open inventory file: "+lockFilePath, err)
		return
	}

	err = syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to lock inventory file: "+lockFilePath, err)
		return
	}

	defer lockFile.Close()
	defer syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN)
	// defer os.Remove(lockFilePath)

	// got lock.

	cmd := exec.Command("/bin/bash", "-c", "echo error >&2 && sleep 20 && echo hello world")

	cmd.Stdout = lockFile
	cmd.Stderr = lockFile

	if err := cmd.Run(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to execute command", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
	})

}
