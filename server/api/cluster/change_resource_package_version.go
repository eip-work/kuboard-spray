package cluster

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/api/cluster/cluster_common"
	"github.com/opencmit/pangee-cluster/api/command"
	"github.com/opencmit/pangee-cluster/common"
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
			if host["pangeecluster_node_action"] != nil {
				common.HandleError(c, http.StatusInternalServerError, hostName+" is pending to "+host["pangeecluster_node_action"].(string), nil)
				return
			}
			if hostName != "localhost" && hostName != "bastion" {
				common.MapSet(metadata.Inventory, "all.hosts."+hostName+".pangeecluster_node_action", "upgrade_node")
				common.MapSet(metadata.Inventory, "all.hosts."+hostName+".pangeecluster_require_download", true)
			}
		}
	}

	common.MapSet(metadata.Inventory, "all.hosts.localhost.pangeecluster_resource_package_previous", common.MapGetString(metadata.Inventory, "all.hosts.localhost.pangeecluster_resource_package"))
	common.MapSet(metadata.Inventory, "all.hosts.localhost.pangeecluster_resource_package", req.TargetVersion)

	if common.MapGet(metadata.Inventory, "all.children.target.children.k8s_cluster.vars.calico_ippool_vxlan") == nil {
		shellReq1 := command.AnsibleCommandsRequest{
			Name:    "calico",
			Command: `kubectl get ippools.crd.projectcalico.org default-ipv4-ippool -o json`,
		}
		shellResult, err := command.ExecuteShellCommandsAbortOnFirstSuccess("cluster", req.Cluster, "kube_control_plane[0]", []command.AnsibleCommandsRequest{shellReq1})
		if err != nil {
			common.HandleError(c, http.StatusInternalServerError, "failed to get calico status", err)
			return
		}
		stdout := shellResult[0].StdOut
		calicoIpPool := map[string]interface{}{}
		if err := json.Unmarshal([]byte(stdout), &calicoIpPool); err == nil {
			vxlanMode := common.MapGetString(calicoIpPool, "spec.vxlanMode")
			common.MapSet(metadata.Inventory, "all.children.target.children.k8s_cluster.vars.calico_ippool_vxlan", vxlanMode)
		} else {
			common.HandleError(c, http.StatusInternalServerError, "failed to parse calicoIpPool", err)
			return
		}
	}

	if err := cluster_common.SaveInventory(req.Cluster, metadata.Inventory); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to save inventory", err)
		return
	}

	c.JSON(http.StatusOK, common.PangeeClusterResponse{
		Code:    http.StatusOK,
		Message: "success",
	})
}
