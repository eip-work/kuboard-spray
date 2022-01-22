package state

import (
	"net/http"
	"strings"

	"github.com/eip-work/kuboard-spray/api/cluster"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

type AddonStatus struct {
	Name            string      `json:"addon_name"`
	IntendToInstall bool        `json:"intend_to_install"`
	IsInstalled     bool        `json:"is_installed"`
	Code            string      `json:"code"`
	StdOut          string      `json:"stdout"`
	StdOutObj       interface{} `json:"stdout_obj"`
}

func CheckAddonStatus(c *gin.Context) {
	var request GetNodesRequest
	c.ShouldBindUri(&request)

	result := map[string]AddonStatus{}

	inventoryPath := cluster.ClusterInventoryYamlPath(request.ClusterName)
	inventory, err := common.ParseYamlFile(inventoryPath)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to parse inventory", err)
		return
	}
	resourcePackagePath := cluster.ResourcePackagePathForInventory(inventory)

	resourcePackage, err := common.ParseYamlFile(resourcePackagePath + "/package.yaml")
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to parse resourcePackage", err)
		return
	}

	k8sVars := common.MapGet(inventory, "all.children.target.children.k8s_cluster.vars").(map[string]interface{})
	addons := common.MapGet(resourcePackage, "data.addon").([]interface{})
	for _, v := range addons {
		addon := v.(map[string]interface{})
		addonName := addon["name"].(string)
		addonTarget := addon["target"].(string)
		if k8sVars[addonTarget] != true {
			result[addonName] = AddonStatus{
				Name:            addonName,
				IntendToInstall: false,
				IsInstalled:     false,
			}
			continue
		}
		var lifecycle map[string]interface{}
		if addon["lifecycle"] != nil {
			lifecycle = addon["lifecycle"].(map[string]interface{})
		}
		if lifecycle["check"] == nil {
			result[addonName] = AddonStatus{
				Name:            addonName,
				IntendToInstall: true,
				IsInstalled:     false,
				StdOut:          "No data.addon.*.check field in package.yaml",
			}
			continue
		}
		checker := lifecycle["check"].(map[string]interface{})
		shell := checker["shell"].(string)
		keyword := checker["keyword"].(string)
		out, err := ExecuteShellOnControlPlane(request.ClusterName, shell+" || true")
		if err != nil {
			result[addonName] = AddonStatus{
				Name:            addonName,
				IntendToInstall: true,
				IsInstalled:     false,
				Code:            "500",
				StdOut:          err.Error(),
			}
			continue
		}
		result[addonName] = AddonStatus{
			Name:            addonName,
			IntendToInstall: true,
			IsInstalled:     strings.Contains(out.Stdout, keyword),
			Code:            out.ReturnCode,
			StdOut:          out.Stdout,
			StdOutObj:       out.StdoutObj,
		}
	}

	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    result,
	})
}
