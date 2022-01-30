package state

import (
	"net/http"
	"strings"
	"time"

	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AddonStatus struct {
	Name            string      `json:"addon_name"`
	IntendToInstall bool        `json:"intend_to_install"`
	IsInstalled     bool        `json:"is_installed"`
	Code            string      `json:"code"`
	StdOut          string      `json:"stdout"`
	StdOutObj       interface{} `json:"stdout_obj"`
	Index           int
}

func CheckAddonStatus(c *gin.Context) {
	var request GetNodesRequest
	c.ShouldBindUri(&request)

	result := map[string]AddonStatus{}

	inventoryPath := cluster_common.ClusterInventoryYamlPath(request.ClusterName)
	inventory, err := common.ParseYamlFile(inventoryPath)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to parse inventory", err)
		return
	}
	resourcePackagePath := cluster_common.ResourcePackagePathForInventory(inventory)

	resourcePackage, err := common.ParseYamlFile(resourcePackagePath + "/package.yaml")
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to parse resourcePackage", err)
		return
	}

	k8sVars := common.MapGet(inventory, "all.children.target.children.k8s_cluster.vars").(map[string]interface{})
	addons := common.MapGet(resourcePackage, "data.addon").([]interface{})
	startTime := time.Now()
	commands := []string{}
	count := 0
	for _, v := range addons {
		addon := v.(map[string]interface{})
		addonName := addon["name"].(string)
		addonTarget := addon["target"].(string)
		if k8sVars[addonTarget] != true {
			continue
		}
		var lifecycle map[string]interface{}
		if addon["lifecycle"] != nil {
			lifecycle = addon["lifecycle"].(map[string]interface{})
		}
		if lifecycle["check"] == nil {
			continue
		}
		checker := lifecycle["check"].(map[string]interface{})
		shell := checker["shell"].(string)
		commands = append(commands, shell)
		result[addonName] = AddonStatus{Index: count}
		count++
	}

	out, err := cluster_common.ExecuteShellCommandsOnControlPlane(request.ClusterName, commands)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "error", err)
		return
	}

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

		keyword := checker["keyword"].(string)
		result[addonName] = AddonStatus{
			Name:            addonName,
			IntendToInstall: true,
			IsInstalled:     strings.Contains(out.Stdouts[result[addonName].Index], keyword),
			Code:            out.ReturnCode,
			StdOut:          out.Stdouts[result[addonName].Index],
			StdOutObj:       out.StdoutObjs[result[addonName].Index],
		}
	}

	duration := time.Now().UnixNano() - startTime.UnixNano()
	logrus.Trace("duration: ", duration/1000000)

	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    result,
	})
}
