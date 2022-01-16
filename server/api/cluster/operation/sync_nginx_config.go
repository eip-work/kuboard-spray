package operation

import (
	"net/http"
	"strconv"

	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

func SyncNginxConfigActions(c *gin.Context) {
	var req InstallClusterRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)

	inventory, resourcePackage, err := updateResourcePackageVarsToInventory(req.Cluster)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to process inventory", err)
		return
	}

	common.MapSet(inventory, "all.vars.kuboardspray_no_log", !req.Verbose)

	postExec := func(status command.ExecuteExitStatus) (string, error) {
		success := status.Success
		var message string
		if success {
			message = "\033[32m[ " + "Completed sync config tasks." + " ]\033[0m \n"
			message += "\033[32m[ " + "已完成同步集群配置的任务，您可以继续其他操作。" + " ]\033[0m \n"
			PostProcessInventory(req.Cluster, "sync_nginx_config")
		} else {
			message = "\033[31m\033[01m\033[05m[ " + "Failed to synchronize config to nodes. Please review the logs and fix the problem." + " ]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[ " + "同步集群配置失败，请回顾日志，找到错误信息，并解决问题后，再次尝试。" + " ]\033[0m \n"
		}

		return "\n" + message, nil
	}

	cmd := command.Execute{
		OwnerType: "cluster",
		OwnerName: req.Cluster,
		Cmd:       "ansible-playbook",
		Args: func(execute_dir string) []string {

			playbook := common.MapGet(resourcePackage, "data.supported_playbooks.sync_nginx_config").(string)
			result := []string{"-i", execute_dir + "/inventory.yaml", playbook, "--fork", strconv.Itoa(req.Fork), "--tags", "nginx"}
			if req.ExcludeNodes != "" {
				result = append(result, "--limit ", req.ExcludeNodes)
			}
			if req.VVV {
				result = append(result, "-vvv")
			}
			return result
		},
		Dir:      resourcePackagePathForInventory(inventory),
		Type:     "sync_nginx_config",
		PreExec:  func(execute_dir string) error { return common.SaveYamlFile(execute_dir+"/inventory.yaml", inventory) },
		PostExec: postExec,
	}

	if err := cmd.Exec(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Faild to synchronize config to nodes. ", err)
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
