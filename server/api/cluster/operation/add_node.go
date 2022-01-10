package operation

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/eip-work/kuboard-spray/api/cluster/state"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func AddNode(c *gin.Context) {
	var req InstallClusterRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)

	inventory, resourcePackage, err := updateResourcePackageVarsToInventory(req.Cluster)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to process inventory", err)
		return
	}
	common.MapSet(inventory, "all.vars.kuboardspray_no_log", !req.Verbose)

	// 获取在线的节点
	result, err := state.ExecuteShellOnControlPlane(req.Cluster, "kubectl get nodes -o json")

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed", err)
		return
	}
	onlineNodes := common.MapGet(result.StdoutObj, "items").([]interface{})

	// 找出 inventory 中待添加的节点
	hosts := common.MapGet(inventory, "all.hosts").(map[string]interface{})
	nodes := ""
	includesControlPlane := false
	controlPlaneHosts := common.MapGet(inventory, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts").(map[string]interface{})
	for key := range hosts {
		flag := false
		for _, node := range onlineNodes {
			onlineNode := node.(map[string]interface{})
			name := common.MapGet(onlineNode, "metadata.name")
			if key == name {
				flag = true
				break
			}
		}
		if !flag {
			nodes += key + ","
			if controlPlaneHosts[key] != nil {
				includesControlPlane = true
			}
		}
	}
	nodes = strings.Trim(nodes, ",")

	logrus.Trace("add_nodes: ", nodes)

	postExec := func(status command.ExecuteExitStatus) (string, error) {

		success := status.Success
		var message string
		if success {
			message = "\033[32m[ " + "Nodes are already added to cluster, please release the machine." + " ]\033[0m \n"
			message += "\033[32m[ " + "节点已添加到 Kubernetes 集群，请返回集群详情页查看" + " ]\033[0m \n"
		} else {
			message = "\033[31m\033[01m\033[05m[" + "Failed to add node. Please review the logs and fix the problem." + "]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[" + "添加节点失败，请回顾日志，找到错误信息，并解决问题后，再次尝试。" + "]\033[0m \n"
		}
		return "\n" + message, nil
	}

	playbook := common.MapGet(resourcePackage, "data.supported_playbooks.add_node").(string)
	clusterPlaybook := common.MapGet(resourcePackage, "data.supported_playbooks.install_cluster").(string)

	cmd := command.Execute{
		OwnerType: "cluster",
		OwnerName: req.Cluster,
		Cmd:       "ansible-playbook",
		Args: func(execute_dir string) []string {
			result := []string{"-i", execute_dir + "/inventory.yaml", playbook, "--fork", strconv.Itoa(req.Fork), "--limit", nodes}
			if includesControlPlane {
				result = []string{"-i", execute_dir + "/inventory.yaml", clusterPlaybook, "--fork", strconv.Itoa(req.Fork)}
			}
			if req.VVV {
				result = append(result, "-vvv")
			}
			return result
		},
		Dir:      resourcePackagePathForInventory(inventory),
		Type:     "add_node",
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
