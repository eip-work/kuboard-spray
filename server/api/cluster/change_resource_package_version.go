package cluster

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

type ChangeVersionRequest struct {
	GetClusterRequest
	Type          string `json:"type"`
	TargetVersion string `json:"target_version"`
}

func ChangeResourcePackageVersion(c *gin.Context) {
	var req ChangeVersionRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)
	metadata, err := cluster_common.ClusterMetadataByName(req.Cluster)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to load cluster metadata", err)
		return
	}

	hosts := common.MapGet(metadata.Inventory, "all.hosts").(map[string]interface{})

	if req.Type == "upgrade" {
		for hostName, v := range hosts {
			host := v.(map[string]interface{})
			if host["kuboardspray_node_action"] != nil {
				common.HandleError(c, http.StatusInternalServerError, hostName+" is pending to "+host["kuboardspray_node_action"].(string), nil)
				return
			}
			if hostName != "localhost" && hostName != "bastion" {
				common.MapSet(metadata.Inventory, "all.hosts."+hostName+".kuboardspray_node_action", "upgrade_node")
				common.MapSet(metadata.Inventory, "all.hosts."+hostName+".kuboardspray_require_download", true)
			}
		}
	}

	common.MapSet(metadata.Inventory, "all.hosts.localhost.kuboardspray_resource_package_previous", common.MapGetString(metadata.Inventory, "all.hosts.localhost.kuboardspray_resource_package"))
	common.MapSet(metadata.Inventory, "all.hosts.localhost.kuboardspray_resource_package", req.TargetVersion)

	if err := cluster_common.SaveInventory(req.Cluster, metadata.Inventory); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to save inventory", err)
	}

	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
	})
}
