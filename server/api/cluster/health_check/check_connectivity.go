package health_check

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/api/cluster/cluster_common"
	"github.com/opencmit/pangee-cluster/api/command"
	"github.com/opencmit/pangee-cluster/common"
)

type CheckConnectivityRequest struct {
	ClusterName string `uri:"cluster"`
}

func CheckConnectivity(c *gin.Context) {
	var request CheckConnectivityRequest
	c.ShouldBindUri(&request)

	inventoryPath := cluster_common.ClusterInventoryYamlPath(request.ClusterName)
	inventory, err := common.ParseYamlFile(inventoryPath)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to parse inventory", err)
		return
	}

	netcheckerPort := 31081
	np := common.MapGet(inventory, "all.children.target.children.k8s_cluster.vars.netchecker_port")
	if np != nil {
		netcheckerPort = np.(int)
	}

	shellReq := command.AnsibleCommandsRequest{
		Name:    "shell",
		Command: "curl -s http://$(hostname):" + strconv.Itoa(netcheckerPort) + "/api/v1/connectivity_check",
	}
	shellResult, err := command.ExecuteShellCommandsAbortOnFirstSuccess("cluster", request.ClusterName, "kube_control_plane[0]", []command.AnsibleCommandsRequest{shellReq})
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed", err)
		return
	}

	c.JSON(http.StatusOK, common.PangeeClusterResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    shellResult[0],
	})
}
