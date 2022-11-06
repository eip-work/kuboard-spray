package operation

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

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

		clusterMetadata, _ := cluster_common.ClusterMetadataByName(req.Cluster)

		count := 0
		failedCount := 0
		nodes := ""
		failedNodes := ""
		for _, nodeStatus := range status.NodeStatus {
			logrus.Trace(nodeStatus)
			if nodeStatus.OK != "0" && nodeStatus.Unreachable == "0" && nodeStatus.Failed == "0" {
				common.MapDelete(clusterMetadata.Inventory, "all.hosts."+nodeStatus.NodeName+".kuboardspray_require_download")
				count++
				nodes += nodeStatus.NodeName + ","
			} else {
				failedCount++
				failedNodes += nodeStatus.NodeName + ","
			}
		}
		nodes = strings.Trim(nodes, ",")
		failedNodes = strings.Trim(failedNodes, ",")

		inventoryNewContent, _ := yaml.Marshal(clusterMetadata.Inventory)

		if err := ioutil.WriteFile(clusterMetadata.InventoryPath, inventoryNewContent, 0655); err != nil {
			logrus.Trace(err)
			message += "写入 inventory 失败: "
			message += err.Error()
		}
		if count > 0 {
			message += "\n"
			message = "\033[32m[ " + "Finished in distributing installation binaries/images for " + strconv.Itoa(count) + " nodes: " + nodes + ". ]\033[0m \n"
			message += "\033[32m[ 成功加载安装包/容器镜像到 " + strconv.Itoa(count) + " 个节点: " + nodes + " ]\033[0m \n"
		}
		if failedCount > 0 {
			message += "\n"
			message = "\033[31m\033[01m\033[05m[ Failed to distribute installation binaries/images for " + strconv.Itoa(failedCount) + " nodes: " + failedNodes + ". ]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[ " + strconv.Itoa(failedCount) + " 个节点加载安装包/容器镜像失败: " + failedNodes + ". ]\033[0m \n"
		}

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
			result = append(result, "--skip-tags", "upgrade")
			return result
		},
		Dir:      cluster_common.ResourcePackageDirForInventory(inventory),
		Type:     req.Operation,
		PreExec:  func(execute_dir string) error { return common.SaveYamlFile(execute_dir+"/inventory.yaml", inventory) },
		PostExec: postExec,
	}

	if err := cmd.Exec(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Faild to download binaries. ", err)
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
