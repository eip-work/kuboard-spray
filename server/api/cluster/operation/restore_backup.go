package operation

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

type RestoreBackupRequest struct {
	OperationCommonRequest
	BackupFileName string `json:"backup_file_name"`
	BackupFilePath string `json:"backup_file_path"`
}

func RestoreBackup(c *gin.Context) {

	var req RestoreBackupRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)
	req.Operation = "restore_etcd"

	inventory, resourcePackage, err := updateResourcePackageVarsToInventory(req.OperationCommonRequest)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to process inventory", err)
		return
	}

	postExec := func(status command.ExecuteExitStatus) (string, error) {

		success := status.Success
		var message string
		if success {
			message += "\n"
			message = "\033[32m[ " + "Restore etcd data successfully." + " ]\033[0m \n"
			message += "\033[32m[ " + "已成功恢复 ETCD 数据。" + " ]\033[0m \n"
		} else {
			message += "\n"
			message = "\033[31m\033[01m\033[05m[" + "Failed to restore etcd." + "]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[" + "恢复 ETCD 数据失败，请回顾错误提示，修正问题后重新尝试。" + "]\033[0m \n"
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
			result = appendCommonParams(result, req.OperationCommonRequest, false)
			result = append(result, "-e", "backup_file_path="+req.BackupFilePath)
			result = append(result, "-e", "backup_file_name="+req.BackupFileName)
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
