package state

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/api/cluster"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GetNodesRequest struct {
	ClusterName string `uri:"cluster"`
}

func GetNodes(c *gin.Context) {

	var request GetNodesRequest
	c.ShouldBindUri(&request)

	inventoryYamlPath := cluster.ClusterInventoryYamlPath(request.ClusterName)

	cmd := command.Run{
		Cmd: "ansible",
		Args: []string{
			"kube_control_plane[0]",
			"-m", "shell",
			"-a", "kubectl get nodes -o json",
			"-i", inventoryYamlPath,
		},
		Env: []string{"ANSIBLE_CONFIG=" + constants.GET_ADHOC_CFG_PATH()},
		// Dir: resourcePackageDir,
	}
	stdout, stderr, err := cmd.Run()
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot fetch kubeconfig", err)
		return
	}
	logrus.Trace(string(stdout), string(stderr))

	result := command.ParseAnsibleStdout(string(stdout))

	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    result,
	})

}
