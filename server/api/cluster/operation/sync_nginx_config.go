package operation

import (
	"io/ioutil"
	"net/http"

	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func SyncNginxConfigActions(c *gin.Context) {
	var req OperationCommonRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)
	req.Operation = "sync_nginx_config"

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
			message = "\033[32m[ " + "Completed task \"Update apiserver list in loadbalancer\"" + " ]\033[0m \n"
			message += "\033[32m[ " + "已完成 \"更新负载均衡中 apiserver 列表\" 的任务，您可以继续其他操作。" + " ]\033[0m \n"

			inventoryPath := cluster_common.ClusterInventoryYamlPath(req.Cluster)
			inventoryNew, _ := common.ParseYamlFile(inventoryPath)
			common.MapSet(inventoryNew, "all.hosts.localhost.kuboardspray_sync_nginx_config", false)
			inventoryNewContent, _ := yaml.Marshal(inventoryNew)
			if err := ioutil.WriteFile(inventoryPath, inventoryNewContent, 0655); err != nil {
				logrus.Trace(err)
			}
		} else {
			message += "\n"
			message = "\033[31m\033[01m\033[05m[ " + "Failed to do \"Update apiserver list in loadbalancer\". Please review the logs and fix the problem." + " ]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[ " + "\"更新负载均衡中 apiserver 列表\" 失败，请回顾日志，找到错误信息，并解决问题后，再次尝试。" + " ]\033[0m \n"
		}

		return "\n" + message, nil
	}

	cmd := command.Execute{
		OwnerType: "cluster",
		OwnerName: req.Cluster,
		Cmd:       "ansible-playbook",
		Args: func(execute_dir string) []string {

			playbook := common.MapGet(resourcePackage, "data.supported_playbooks.sync_nginx_config").(string)
			result := []string{"-i", execute_dir + "/inventory.yaml", playbook, "--tags", "nginx,bastion"}
			result = appendCommonParams(result, req, false)
			return result
		},
		Dir:      cluster_common.ResourcePackageDirForInventory(inventory),
		Type:     req.Operation,
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
