package health_check

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/api/ansible_rpc"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

type CheckConnectivityDetailsRequest struct {
	CheckConnectivityRequest
	Namespace string `form:"namespace"`
	Pod       string `form:"pod"`
	Type      string `form:"type"`
}

func CheckConnectivityDetails(c *gin.Context) {
	var request CheckConnectivityDetailsRequest
	c.ShouldBindUri(&request)
	c.ShouldBindQuery(&request)

	var shell string
	if request.Type == "events" {
		shell = "kubectl get events -n " + request.Namespace + " --field-selector involvedObject.name=" + request.Pod
	} else if request.Type == "logs" {
		shell = "kubectl logs -n " + request.Namespace + " " + request.Pod + " --tail 100"
	} else if request.Type == "podList" {
		shell = "kubectl get pods -n " + request.Namespace + " --selector 'app in (netchecker-agent,netchecker-server,netchecker-agent-hostnet)' -o wide"
	} else { // podDetails
		shell = "kubectl get pods -n " + request.Namespace + " " + request.Pod + " -o yaml"
	}

	shellReq := ansible_rpc.AnsibleCommandsRequest{
		Name:    "shell",
		Command: shell,
	}
	shellResult, err := ansible_rpc.ExecuteShellCommandsAbortOnFirstSuccess("cluster", request.ClusterName, "kube_control_plane[0]", []ansible_rpc.AnsibleCommandsRequest{shellReq})
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed", err)
		return
	}

	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    shellResult[0],
	})
}
