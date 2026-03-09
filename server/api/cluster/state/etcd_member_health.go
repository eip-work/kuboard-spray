package state

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/api/command"
	"github.com/opencmit/pangee-cluster/common"
)

func CheckEtcdEndpoints(c *gin.Context) {
	var request GetNodesRequest
	c.ShouldBindUri(&request)

	shellReq1 := command.AnsibleCommandsRequest{
		Name:    "members",
		Command: `etcdctl member list --write-out=json`,
	}
	shellReq2 := command.AnsibleCommandsRequest{
		Name:    "health",
		Command: `etcdctl endpoint health --cluster=true --write-out=json`,
	}
	shellResult, err := command.ExecuteShellCommandsAbortOnFirstSuccess("cluster", request.ClusterName, "etcd[0]", []command.AnsibleCommandsRequest{shellReq1, shellReq2})

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed", err)
		return
	}

	// logrus.Trace(shellResult)

	c.JSON(http.StatusOK, common.PangeeClusterResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    shellResult,
	})

}
