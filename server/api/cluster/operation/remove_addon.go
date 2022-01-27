package operation

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func RemoveAddon(c *gin.Context) {
	var req InstallAddonRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)
	req.Operation = "remove_addon"

	inventory, resourcePackage, err := updateResourcePackageVarsToInventory(req.OperationCommonRequest)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to process inventory", err)
		return
	}

	addonsInPackage := common.MapGet(resourcePackage, "data.addon").([]interface{})
	var addon map[string]interface{}
	for _, a := range addonsInPackage {
		temp := a.(map[string]interface{})
		if temp["name"] == req.AddonName {
			addon = temp
		}
	}

	postExec := func(status command.ExecuteExitStatus) (string, error) {

		success := status.Success
		var message string
		if success {
			message += "\n"
			message = "\033[32m[ " + "Addon " + req.AddonName + " has been removed successfully." + " ]\033[0m \n"
			message += "\033[32m[ " + "已成功卸载可选组件 " + req.AddonName + "。" + " ]\033[0m \n"

			inventoryPath := cluster_common.ClusterInventoryYamlPath(req.Cluster)
			inventoryNew, _ := common.ParseYamlFile(inventoryPath)

			addons := common.MapGet(resourcePackage, "data.addon").([]interface{})
			target := ""
			for _, a := range addons {
				addon := a.(map[string]interface{})
				if addon["name"] == req.AddonName {
					target = addon["target"].(string)
					break
				}
			}
			if target != "" {
				common.MapSet(inventoryNew, "all.children.target.children.k8s_cluster.vars."+target, false)

				inventoryNewContent, _ := yaml.Marshal(inventoryNew)
				if err := ioutil.WriteFile(inventoryPath, inventoryNewContent, 0655); err != nil {
					logrus.Trace(err)
				}
			}
		} else {
			message += "\n"
			message = "\033[31m\033[01m\033[05m[" + "Failed to remove addon " + req.AddonName + "." + "]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[" + "卸载可选组件 " + req.AddonName + " 失败，请回顾错误提示，修正问题后重新尝试。" + "]\033[0m \n"
		}

		return "\n" + message, nil
	}

	playbook := common.MapGet(resourcePackage, "data.supported_playbooks.remove_addon").(string)

	cmd := command.Execute{
		OwnerType: "cluster",
		OwnerName: req.Cluster,
		Cmd:       "ansible-playbook",
		Args: func(execute_dir string) []string {
			result := []string{"-i", execute_dir + "/inventory.yaml", playbook}

			lifecycle := addon["lifecycle"].(map[string]interface{})
			removeTags := lifecycle["remove_addon_tags"].([]interface{})
			tags := ""
			for _, t := range removeTags {
				tags += t.(string) + ","
			}
			tags = strings.Trim(tags, ",")

			result = append(result, "--tags", tags)

			result = appendCommonParams(result, req.OperationCommonRequest, false)
			return result
		},
		Dir:      cluster_common.ResourcePackagePathForInventory(inventory),
		Type:     req.Operation,
		PreExec:  func(execute_dir string) error { return common.SaveYamlFile(execute_dir+"/inventory.yaml", inventory) },
		PostExec: postExec,
	}

	if err := cmd.Exec(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Faild to RemoveAddon. ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"pid": cmd.R_Pid,
		},
	})
}
