package state

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/api/ansible_rpc"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

type CheckCertExpirationRequest struct {
	ClusterName string `uri:"cluster"`
}

func CheckCertExpiration(c *gin.Context) {
	var req CheckCertExpirationRequest
	c.ShouldBindUri(&req)

	cmds := []ansible_rpc.AnsibleCommandsRequest{
		{
			Name:    "check_cert",
			Command: "/usr/local/bin/kubeadm certs check-expiration",
		},
	}
	results, err := ansible_rpc.ExecuteShellCommandsWithStrategy("cluster", req.ClusterName, "kube_control_plane", cmds, "linear")
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to check certificate expiration", err)
	}
	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    results,
	})
}
