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
	Name string `uri:"name"`
}

func GetCluster(c *gin.Context) {
	var req GetClusterRequest
	c.ShouldBindUri(&req)

	inventoryContent, err := ioutil.ReadFile(ClusterInventoryPath(req.Name))
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot open file: "+ClusterInventoryPath(req.Name))
		return
	}

	inventory := gin.H{}
	yaml.Unmarshal(inventoryContent, inventory)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"inventory": inventory,
		},
	})
}

func ClusterInventoryPath(name string) string {
	return constants.GET_DATA_INVENTORY_DIR() + "/" + name + "/inventory.yaml"
}
