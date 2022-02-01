package state

import (
	"net/http"
	"strings"
	"time"

	"github.com/eip-work/kuboard-spray/api/ansible_rpc"
	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AddonStatus struct {
	Name            string      `json:"addon_name"`
	IntendToInstall bool        `json:"intend_to_install"`
	IsInstalled     bool        `json:"is_installed"`
	Code            int         `json:"code"`
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
	resourcePackagePath := cluster_common.ResourcePackageDirForInventory(inventory)

	resourcePackage, err := common.ParseYamlFile(resourcePackagePath + "/package.yaml")
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to parse resourcePackage", err)
		return
	}

	k8sVars := common.MapGet(inventory, "all.children.target.children.k8s_cluster.vars").(map[string]interface{})
	addons := common.MapGet(resourcePackage, "data.addon").([]interface{})
	startTime := time.Now()
	commands := []ansible_rpc.AnsibleCommandsRequest{}
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
		shell := checker["shell"].(string) + " || true"
		commands = append(commands, ansible_rpc.AnsibleCommandsRequest{
			Command: shell,
			Name:    addonName,
		})
	}

	out, err := ansible_rpc.ExecuteShellCommands("cluster", request.ClusterName, "kube_control_plane", commands)

	// logrus.Trace(out)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "error", err)
		return
	}

	if len(out.Plays) > 0 {
		for _, task := range out.Plays[0].Tasks {
			keyword := "KEYWORD_NOT_FOUND_YET"
			for _, v := range addons {
				addon := v.(map[string]interface{})
				if addon["name"] == task.Task.Name {
					keyword = common.MapGet(addon, "lifecycle.check.keyword").(string)
				}
			}
			for _, nodeResult := range task.Hosts {
				if nodeResult.Changed {
					result[task.Task.Name] = AddonStatus{
						Name:            task.Task.Name,
						IntendToInstall: true,
						IsInstalled:     strings.Contains(nodeResult.StdOut, keyword),
						StdOut:          nodeResult.StdOut,
					}
					break
				}
			}
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
