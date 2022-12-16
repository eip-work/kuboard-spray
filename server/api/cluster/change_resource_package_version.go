package cluster

import (
	"encoding/json"
	"net/http"

	"github.com/eip-work/kuboard-spray/api/ansible_rpc"
	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

	netmanager := common.MapGetString(metadata.Inventory, "all.children.target.children.k8s_cluster.vars.kube_network_plugin")

	if netmanager == "calico" && (common.MapGet(metadata.Inventory, "all.children.target.children.k8s_cluster.vars.calico_ipip_mode") == nil ||
		common.MapGet(metadata.Inventory, "all.children.target.children.k8s_cluster.vars.calico_vxlan_mode") == nil ||
		common.MapGet(metadata.Inventory, "all.children.target.children.k8s_cluster.vars.calico_network_backend") == nil) {
		shellReq1 := ansible_rpc.AnsibleCommandsRequest{
			Name:    "calico",
			Command: `calicoctl.sh get ipPool default-pool -o json`,
		}
		shellReq2 := ansible_rpc.AnsibleCommandsRequest{
			Name:    "calico",
			Command: `kubectl get cm -n kube-system calico-config -o json`,
		}
		shellResult, err := ansible_rpc.ExecuteShellCommandsAbortOnFirstSuccess("cluster", req.Cluster, "kube_control_plane[0]", []ansible_rpc.AnsibleCommandsRequest{shellReq1, shellReq2})
		if err != nil {
			common.HandleError(c, http.StatusInternalServerError, "failed to get calico status", err)
			return
		}
		stdout := shellResult[0].StdOut
		calicoIpPool := map[string]interface{}{}
		if err := json.Unmarshal([]byte(stdout), &calicoIpPool); err == nil {
			ipipMode := common.MapGetString(calicoIpPool, "spec.ipipMode")
			vxlanMode := common.MapGetString(calicoIpPool, "spec.vxlanMode")
			common.MapSet(metadata.Inventory, "all.children.target.children.k8s_cluster.vars.calico_ipip_mode", ipipMode)
			common.MapSet(metadata.Inventory, "all.children.target.children.k8s_cluster.vars.calico_vxlan_mode", vxlanMode)
		} else {
			common.HandleError(c, http.StatusInternalServerError, "failed to parse calicoIpPool", err)
			return
		}

		stdout2 := shellResult[1].StdOut
		calicoBackend := map[string]interface{}{}
		if err := json.Unmarshal([]byte(stdout2), &calicoBackend); err == nil {
			calico_network_backend := common.MapGetString(calicoBackend, "data.calico_backend")
			common.MapSet(metadata.Inventory, "all.children.target.children.k8s_cluster.vars.calico_network_backend", calico_network_backend)
		} else {
			common.HandleError(c, http.StatusInternalServerError, "failed to parse calico_network_backend", err)
			return
		}

	}

	addons := common.MapGet(metadata.ResourcePackage, "data.addon").([]interface{})
	for _, addon := range addons {
		addonMap := addon.(map[string]interface{})
		targetName := common.MapGetString(addonMap, "target")

		if common.MapGet(metadata.Inventory, "all.children.target.children.k8s_cluster.vars."+targetName) == nil {
			logrus.Trace("设置变量： all.children.target.children.k8s_cluster.vars." + targetName + " 为 false")
			common.MapSet(metadata.Inventory, "all.children.target.children.k8s_cluster.vars."+targetName, false)
		} else {
			logrus.Trace("已经设置变量： all.children.target.children.k8s_cluster.vars." + targetName)
		}
	}

	if err := cluster_common.SaveInventory(req.Cluster, metadata.Inventory); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to save inventory", err)
		return
	}

	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
	})
}
