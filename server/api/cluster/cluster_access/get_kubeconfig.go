package cluster_access

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/api/cluster"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetKubeConfig(c *gin.Context) {
	var req cluster.GetClusterRequest
	c.ShouldBindUri(&req)

	inventoryYamlPath := cluster.ClusterInventoryYamlPath(req.Cluster)

	cmd := command.Run{
		Cmd: "ansible",
		Args: []string{
			"kube_control_plane[0]",
			"-m", "shell",
			"-a", "cat /root/.kube/config",
			"-i", inventoryYamlPath,
			"-e", "kuboardspray_cluster_dir=" + constants.GET_DATA_DIR() + "/cluster/" + req.Cluster,
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
