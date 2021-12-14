package cluster

import (
	"io/ioutil"
	"net/http"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	yaml "gopkg.in/yaml.v3"
)

type GetClusterRequest struct {
	Cluster string `uri:"cluster"`
}

func GetCluster(c *gin.Context) {
	var req GetClusterRequest
	c.ShouldBindUri(&req)

	inventoryContent, err := ioutil.ReadFile(ClusterInventoryPath(req.Cluster))
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot open file: "+ClusterInventoryPath(req.Cluster), err)
		return
	}

	inventory := gin.H{}
	err = yaml.Unmarshal(inventoryContent, inventory)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot parse file: "+ClusterInventoryPath(req.Cluster), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"inventory": inventory,
			"name":      req.Cluster,
		},
	})
}

func ClusterInventoryPath(cluster string) string {
	return constants.GET_DATA_INVENTORY_DIR() + "/" + cluster + "/inventory.yaml"
}
