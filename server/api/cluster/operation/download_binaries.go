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

func DownloadBinaries(c *gin.Context) {

	var req OperationCommonRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)
	req.Operation = "download_binaries"

	inventory, resourcePackage, err := updateResourcePackageVarsToInventory(req)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to process inventory", err)
		return
	}

	postExec := func(status command.ExecuteExitStatus) (string, error) {

		var message string

		kubeletVersions, err := getKubeletVersion(req.Cluster)
		if err != nil {
			message += "\n"
			message += "\033[31m\033[01m\033[05m[ 获取更新后的版本时出现错误：" + err.Error() + "]\033[0m \n"
			return "\n" + message, nil
		}

		clusterMetadata, _ := cluster_common.ClusterMetadataByName(req.Cluster)

		newVersion := common.MapGetString(clusterMetadata.ResourcePackage, "data.kubernetes.kube_version")
		logrus.Trace("newVersion:", newVersion)
		upgradedNodes := ""
		count := 0
		for node, version := range kubeletVersions {
			logrus.Trace("node:", node, " version:", version)
			if newVersion == version {
				common.MapDelete(clusterMetadata.Inventory, "all.hosts."+node+".kuboardspray_node_action")
				upgradedNodes += node + ","
				count++
			}
		}
		// upgradedNodes = strings.Trim(upgradedNodes, ",")

		inventoryNewContent, _ := yaml.Marshal(clusterMetadata.Inventory)

		if err := ioutil.WriteFile(clusterMetadata.InventoryPath, inventoryNewContent, 0655); err != nil {
			logrus.Trace(err)
		}

		// if count > 0 {
		// 	message += "\n"
		// 	message = "\033[32m[ " + "Succeeded in upgrading " + strconv.Itoa(count) + " nodes: " + upgradedNodes + ". ]\033[0m \n"
		// 	message += "\033[32m[ 成功升级了 " + strconv.Itoa(count) + " 个节点: " + upgradedNodes + " ]\033[0m \n"
		// } else {
		// 	message += "\n"
		// 	message = "\033[31m\033[01m\033[05m[" + "Failed to upgrade. Please review the logs and fix the problem." + "]\033[0m \n"
		// 	message += "\033[31m\033[01m\033[05m[" + "集群升级失败，请回顾日志，找到错误信息，并解决问题后，再次尝试。" + "]\033[0m \n"
		// }

		return "\n" + message, nil
	}

	playbook := common.MapGet(resourcePackage, "data.supported_playbooks.upgrade_cluster").(string)

	cmd := command.Execute{
		OwnerType: "cluster",
		OwnerName: req.Cluster,
		Cmd:       "ansible-playbook",
		Args: func(execute_dir string) []string {
			result := []string{"-i", execute_dir + "/inventory.yaml", playbook}
			result = appendCommonParams(result, req, false)
			result = append(result, "--tags", "download")
			result = append(result, "-skip-tags", "upgrade")
			return result
		},
		Dir:      cluster_common.ResourcePackageDirForInventory(inventory),
		Type:     req.Operation,
		PreExec:  func(execute_dir string) error { return common.SaveYamlFile(execute_dir+"/inventory.yaml", inventory) },
		PostExec: postExec,
	}

	if err := cmd.Exec(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Faild to InstallCluster. ", err)
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
