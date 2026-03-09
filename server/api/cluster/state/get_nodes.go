package state

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/api/command"
	"github.com/opencmit/pangee-cluster/common"
)

type GetNodesRequest struct {
	ClusterName string `uri:"cluster"`
}

func GetNodes(c *gin.Context) {

	var request GetNodesRequest
	c.ShouldBindUri(&request)

	shellReq := command.AnsibleCommandsRequest{
		Name:    "nodes",
		Command: `kubectl get nodes -o json`,
	}
	shellResult, err := command.ExecuteShellCommandsAbortOnFirstSuccess("cluster", request.ClusterName, "kube_control_plane[0]", []command.AnsibleCommandsRequest{shellReq})

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed", err)
		return
	}

	stdOut, stdOutObj := IPLink(shellResult[0].StdOut, request.ClusterName)
	shellResult[0].StdOut = stdOut
	shellResult[0].StdOutObj = stdOutObj

	c.JSON(http.StatusOK, common.PangeeClusterResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    shellResult[0],
	})

}
