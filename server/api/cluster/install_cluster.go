package cluster

import (
	"net/http"
	"strconv"

	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

var count = 0

func InstallCluster(c *gin.Context) {

	var req GetClusterRequest
	c.ShouldBindUri(&req)

	count++

	command := command.Execute{
		Cluster: req.Cluster,
		Cmd:     "/root/test.sh",
		Args:    []string{"param1", "param2", strconv.Itoa(count)},
		Type:    "install",
	}

	if err := command.Exec(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Faild to InstallCluster. ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
	})

}
