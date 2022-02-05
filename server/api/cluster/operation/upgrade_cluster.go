package operation

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/eip-work/kuboard-spray/api/ansible_rpc"
	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type UpgradeClusterRequest struct {
	OperationCommonRequest
	Nodes         string `json:"nodes"`
	SkipDownloads bool   `json:"skip_downloads"`
}

func UpgradeCluster(c *gin.Context) {

	var req UpgradeClusterRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)
	req.Operation = "upgrade_cluster"

	doUpgrade(req, c)
}

func doUpgrade(req UpgradeClusterRequest, c *gin.Context) {
	inventory, resourcePackage, err := updateResourcePackageVarsToInventory(req.OperationCommonRequest)
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
		upgradedNodes = strings.Trim(upgradedNodes, ",")

		inventoryNewContent, _ := yaml.Marshal(clusterMetadata.Inventory)

		if err := ioutil.WriteFile(clusterMetadata.InventoryPath, inventoryNewContent, 0655); err != nil {
			logrus.Trace(err)
		}

		if count > 0 {
			message += "\n"
			message = "\033[32m[ " + "Succeeded in upgrading " + strconv.Itoa(count) + " nodes: " + upgradedNodes + " ]\033[0m \n"
			message += "\033[32m[ 已经成功升级了 " + strconv.Itoa(count) + " 个节点: " + upgradedNodes + " ]\033[0m \n"
		} else {
			message += "\n"
			message = "\033[31m\033[01m\033[05m[" + "Failed to upgrade. Please review the logs and fix the problem." + "]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[" + "集群升级失败，请回顾日志，找到错误信息，并解决问题后，再次尝试。" + "]\033[0m \n"
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
			result = appendCommonParams(result, req.OperationCommonRequest, true)
			result = append(result, "--limit", req.Nodes)
			if req.SkipDownloads {
				result = append(result, "-e", "skip_downloads=True")
			}
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

func getKubeletVersion(clusterName string) (map[string]string, error) {
	cmd := "kubectl get nodes -o go-template --template=\"{{ '{{' }}range.items{{'}}'}}{{ '{{' }}.metadata.name{{'}}'}} {{ '{{' }}.status.nodeInfo.kubeletVersion{{'}}'}},{{ '{{' }}end{{'}}'}}\""

	req := []ansible_rpc.AnsibleCommandsRequest{}
	req = append(req, ansible_rpc.AnsibleCommandsRequest{
		Name:    "getKubeletVersion",
		Command: cmd,
	})

	ansibleResult, err := ansible_rpc.ExecuteShellCommandsAbortOnFirstSuccess("cluster", clusterName, "kube_control_plane[0]", req)
	if err != nil {
		return nil, err
	}

	stdout := ansibleResult[0].StdOut
	logrus.Trace(stdout)
	lines := strings.Split(stdout, ",")

	result := map[string]string{}
	for _, line := range lines {
		if line != "" {
			s := strings.Split(line, " ")
			result[s[0]] = s[1]
		}
	}

	logrus.Trace(result)

	return result, nil
}
