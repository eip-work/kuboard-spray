package state

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/api/command"
	"github.com/opencmit/pangee-cluster/common"
)

type CheckCertExpirationRequest struct {
	ClusterName string `uri:"cluster"`
}

func CheckCertExpiration(c *gin.Context) {
	var req CheckCertExpirationRequest
	c.ShouldBindUri(&req)

	cmds := []command.AnsibleCommandsRequest{
		{
			Name:    "check_cert",
			Command: "/usr/local/bin/kubeadm certs check-expiration",
		},
	}
	results, err := command.ExecuteShellCommandsWithStrategy("cluster", req.ClusterName, "kube_control_plane", cmds, "linear")
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to check certificate expiration", err)
	}
	c.JSON(http.StatusOK, common.PangeeClusterResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    results,
	})
}
