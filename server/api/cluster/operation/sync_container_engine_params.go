package operation

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

func SyncContainerEngineParams(c *gin.Context) {

	var req OperationCommonRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)

	req.Operation = "sync_container_engine_params"

	inventory, resourcePackage, err := updateResourcePackageVarsToInventory(req)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to process inventory", err)
		return
	}

	postExec := func(status command.ExecuteExitStatus) (string, error) {

		var message string
		if status.Success {
			message += "\n"
			message = "\033[32m[ " + "Finished synchronizing container engine params. ]\033[0m \n"
			message += "\033[32m[ 成功同步容器引擎参数 ]\033[0m \n"
		} else {
			message += "\n"
			message = "\033[31m\033[01m\033[05m[ Failed to synchronize container engine params. ]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[ 同步容器引擎参数失败. ]\033[0m \n"
		}

		return "\n" + message, nil
	}

	playbook := common.MapGet(resourcePackage, "data.supported_playbooks."+req.Operation).(string)

	cmd := command.Execute{
		OwnerType: "cluster",
		OwnerName: req.Cluster,
		Cmd:       "ansible-playbook",
		Args: func(execute_dir string) []string {
			result := []string{
				playbook,
				"-i", execute_dir + "/inventory.yaml",
			}
			result = appendCommonParams(result, req, false)
			return result
		},
		Dir:      cluster_common.ResourcePackageDirForInventory(inventory),
		Type:     req.Operation,
		PreExec:  func(execute_dir string) error { return common.SaveYamlFile(execute_dir+"/inventory.yaml", inventory) },
		PostExec: postExec,
	}

	if err := cmd.Exec(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Faild to drain node. ", err)
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
