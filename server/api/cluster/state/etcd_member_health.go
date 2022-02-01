package state

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CheckEtcdEndpoints(c *gin.Context) {
	var request GetNodesRequest
	c.ShouldBindUri(&request)

	cmd := "sh -c 'members=$(etcdctl member list --write-out=json || true) && health=$(etcdctl endpoint health --cluster=true --write-out=json || true) && echo ''{\\\"members\\\":${members},\\\"health\\\":${health}}'''"

	result, err := cluster_common.ExecuteShellOnETCD(request.ClusterName, cmd)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed", err)
		return
	}

	index := strings.Index(result.Stdout, "]}{\"level\":\"")
	if index > 0 {
		if err := json.Unmarshal([]byte(result.Stdout[0:index+2]), &result.StdoutObj); err != nil {
			logrus.Warn("failed to parse: ", result.Stdout[0:index+2])
			logrus.Warn(err)
		}
	}

	logrus.Trace(result.Stdout)

	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    result,
	})

}
