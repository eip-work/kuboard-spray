package operation

import (
	"net/http"
	"strings"

	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

type InstallAddonRequest struct {
	OperationCommonRequest
	AddonName string `json:"addon_name"`
}

func InstallAddon(c *gin.Context) {
	var req InstallAddonRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)
	req.Operation = "install_addon"

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
			message = "\033[32m[ " + "Addon " + req.AddonName + " has been installed successfully." + " ]\033[0m \n"
			message += "\033[32m[ " + "可选组件 " + req.AddonName + " 已成功安装。" + " ]\033[0m \n"
		} else {
			message += "\n"
			message = "\033[31m\033[01m\033[05m[" + "Failed to install addon " + req.AddonName + "." + "]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[" + "安装可选组件 " + req.AddonName + " 失败，请回顾错误提示，修正问题后重新尝试。" + "]\033[0m \n"
		}

		return "\n" + message, nil
	}

	playbook := common.MapGet(resourcePackage, "data.supported_playbooks.install_addon").(string)

	cmd := command.Execute{
		OwnerType: "cluster",
		OwnerName: req.Cluster,
		Cmd:       "ansible-playbook",
		Args: func(execute_dir string) []string {
			result := []string{"-i", execute_dir + "/inventory.yaml", playbook}

			lifecycle := addon["lifecycle"].(map[string]interface{})
			installTags := lifecycle["install_addon_tags"].([]interface{})
			tags := ""
			for _, t := range installTags {
				tags += t.(string) + ","
			}
			tags = strings.Trim(tags, ",")

			result = append(result, "--tags", tags)

			downloads := lifecycle["downloads"].([]interface{})
			download_addon := "["
			for _, d := range downloads {
				download_addon += "\"" + d.(string) + "\","
			}
			download_addon = strings.Trim(download_addon, ",")
			download_addon += "]"
			result = append(result, "--extra-vars", "{\"kuboardspray_download_addon\":"+download_addon+"}")
			result = appendCommonParams(result, req.OperationCommonRequest, false)
			return result
		},
		Dir:      cluster_common.ResourcePackageDirForInventory(inventory),
		Type:     req.Operation,
		PreExec:  func(execute_dir string) error { return common.SaveYamlFile(execute_dir+"/inventory.yaml", inventory) },
		PostExec: postExec,
	}

	if err := cmd.Exec(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Faild to InstallAddon. ", err)
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
