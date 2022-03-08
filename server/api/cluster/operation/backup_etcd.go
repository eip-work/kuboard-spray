package operation

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

func BackupEtcd(c *gin.Context) {

	var req OperationCommonRequest

	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)

	req.Operation = "backup_etcd"

	inventory, resourcePackage, err := updateResourcePackageVarsToInventory(req)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to process inventory", err)
		return
	}

	postExec := func(status command.ExecuteExitStatus) (string, error) {

		success := status.Success
		var message string
		if success {
			message += "\n"
			message = "\033[32m[ " + "Backup etcd successfully." + " ]\033[0m \n"
			message += "\033[32m[ " + "已成功备份 ETCD。" + " ]\033[0m \n"
		} else {
			message += "\n"
			message = "\033[31m\033[01m\033[05m[" + "Failed to backup etcd." + "]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[" + "备份 ETCD 失败，请回顾错误提示，修正问题后重新尝试。" + "]\033[0m \n"
		}

		return "\n" + message, nil
	}

	playbook := common.MapGet(resourcePackage, "data.supported_playbooks."+req.Operation).(string)

	cmd := command.Execute{
		OwnerType: "cluster",
		OwnerName: req.Cluster,
		Cmd:       "ansible-playbook",
		Args: func(execute_dir string) []string {
			result := []string{"-i", execute_dir + "/inventory.yaml", playbook}

			result = appendCommonParams(result, req, false)
			return result
		},
		Dir:      cluster_common.ResourcePackageDirForInventory(inventory),
		Type:     req.Operation,
		PreExec:  func(execute_dir string) error { return common.SaveYamlFile(execute_dir+"/inventory.yaml", inventory) },
		PostExec: postExec,
	}

	if err := cmd.Exec(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Faild to backup etcd. ", err)
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
