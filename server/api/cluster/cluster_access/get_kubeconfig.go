package cluster_access

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/api/cluster"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetKubeConfig(c *gin.Context) {
	var req cluster.GetClusterRequest
	c.ShouldBindUri(&req)

	inventoryYamlPath := cluster.ClusterInventoryYamlPath(req.Cluster)
	// resourcePackageDir := cluster.ClusterResourcePackageDir(req.Cluster)
	// clusterDir := constants.GET_DATA_CLUSTER_DIR() + "/" + req.Cluster

	cmd := command.Run{
		Cmd: "ansible",
		Args: []string{
			"kube_control_plane[0]",
			// "-m", "fetch",
			// "-a", "src=/root/.kube/config dest=" + clusterDir + "/kubeconfig.yaml force=yes",
			"-m", "shell",
			"-a", "cat /root/.kube/config",
			"-i", inventoryYamlPath,
		},
		// Dir: resourcePackageDir,
	}
	stdout, stderr, err := cmd.Run()
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot fetch kubeconfig", err)
		return
	}
	logrus.Trace(string(stdout), string(stderr))

	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    string(stdout),
	})
}
