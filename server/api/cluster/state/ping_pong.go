package state

import (
	"net/http"
	"strings"

	"github.com/eip-work/kuboard-spray/api/cluster"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PingPongRequest struct {
	GetNodesRequest
	NodeNames string `json:"nodes"`
}

func Ping(c *gin.Context) {

	var req PingPongRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)

	pingpongResult, err := PingPong(req.ClusterName, req.NodeNames)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot test pingpong.", err)
		return
	}

	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    pingpongResult,
	})
}

type PingPongResult struct {
	Items map[string]PingPongItem `json:"items"`
}

type PingPongItem struct {
	NodeName string `json:"node_name"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

func PingPong(clusterName string, nodes string) (*PingPongResult, error) {
	inventoryYamlPath := cluster.ClusterInventoryYamlPath(clusterName)

	cmd := command.Run{
		Cmd: "ansible",
		Args: []string{
			nodes,
			"-m", "ping",
			"-i", inventoryYamlPath,
			"--fork", "100",
			"--timeout", "25",
			"-e", "kuboardspray_cluster_dir=" + constants.GET_DATA_DIR() + "/cluster/" + clusterName,
		},
		Env: []string{"ANSIBLE_CONFIG=" + constants.GET_ADHOC_CFG_PATH()},
		// Timeout: 20,
		// Dir: resourcePackageDir,
	}
	stdout, stderr, err := cmd.Run()
	if err != nil {
		return nil, err
	}
	logrus.Trace(string(stdout), string(stderr))

	result := PingPongResult{
		Items: make(map[string]PingPongItem),
	}

	splited := strings.Split(string(stdout), "\n}\n")
	for _, str := range splited {
		temp := str + "\n}"

		if !strings.Contains(temp, " => ") {
			break
		}

		tempSplited := strings.Split(temp, " => ")

		var item PingPongItem
		firstLineSplited := strings.Split(tempSplited[0], " | ")
		item.NodeName = strings.Trim(firstLineSplited[0], " ")
		item.Status = strings.Trim(firstLineSplited[1], " ")
		item.Message = tempSplited[1]

		result.Items[item.NodeName] = item
	}

	return &result, nil
}
