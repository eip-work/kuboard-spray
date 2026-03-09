package state

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/api/cluster/cluster_common"
	"github.com/opencmit/pangee-cluster/api/command"
	"github.com/opencmit/pangee-cluster/common"
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

	c.JSON(http.StatusOK, common.PangeeClusterResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    pingpongResult,
	})
}

type PingPongResult struct {
	Items map[string]PingPongItem `json:"items"`
}

type PingPongItem struct {
	NodeName    string `json:"node_name"`
	Ping        string `json:"ping"`
	UnReachable bool   `json:"unreachable"`
	Message     string `json:"message"`
	StdOut      string `json:"stdout"`
	StdErr      string `json:"stderr"`
}

func PingPong(clusterName string, nodes string) (*PingPongResult, error) {
	inventoryYamlPath := cluster_common.ClusterInventoryYamlPath(clusterName)

	args := []string{nodes, "-m", "ping"}
	results, err := command.ExecuteAdhocCommandWithInventory(inventoryYamlPath, args)
	if err != nil {
		return nil, err
	}

	result := PingPongResult{
		Items: make(map[string]PingPongItem),
	}

	for _, temp := range results {
		item := PingPongItem{
			NodeName:    temp.InventoryHostName,
			Ping:        temp.Ping,
			UnReachable: temp.UnReachable,
			Message:     temp.Msg,
			StdOut:      temp.ModuleStdOut,
			StdErr:      temp.ModuleStdErr,
		}
		result.Items[item.NodeName] = item
	}
	return &result, nil
}
