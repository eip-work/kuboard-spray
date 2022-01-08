package operation

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/eip-work/kuboard-spray/api/cluster"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func RemoveNode(c *gin.Context) {
	var req InstallClusterRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)

	inventory, resourcePackage, err := updateResourcePackageVarsToInventory(req.Cluster)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to process inventory", err)
		return
	}
	common.MapSet(inventory, "all.vars.kuboardspray_no_log", !req.Verbose)

	hosts := common.MapGet(inventory, "all.hosts").(map[string]interface{})
	nodes := ""
	for key, value := range hosts {
		vars := value.(map[string]interface{})
		if vars["kuboard_spray_remove_node"] == true {
			nodes += key + ","
		}
	}
	nodes = "node=" + nodes
	nodes = strings.Trim(nodes, ",")

	logrus.Trace("remove_nodes: ", nodes)

	postExec := func(status command.ExecuteExitStatus) (string, error) {

		success := status.Success
		var message string
		if success {
			message = "\033[32m[ " + "Nodes are already removed, please release the machine." + " ]\033[0m \n"
			message += "\033[32m[ " + "节点已从 Kubernetes 集群中删除，请释放该资源" + " ]\033[0m \n"
		} else {
			message = "\033[31m\033[01m\033[05m[" + "Failed to remove node. Please review the logs and fix the problem." + "]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[" + "删除节点失败，请回顾日志，找到错误信息，并解决问题后，再次尝试。" + "]\033[0m \n"
		}

		inventoryPath := cluster.ClusterInventoryYamlPath(req.Cluster)
		inventoryNew, _ := common.ParseYamlFile(inventoryPath)
		for _, nodeStatus := range status.NodeStatus {
			logrus.Trace(nodeStatus.NodeName, nodeStatus.Failed, "-", nodeStatus.Changed)
			if nodeStatus.NodeName != "localhost" && nodeStatus.NodeName != "bastion" && nodeStatus.Failed == "0" && nodeStatus.Changed != "0" {
				logrus.Trace("deleteNode [", "all.hosts."+nodeStatus.NodeName, "]")
				common.MapDelete(inventoryNew, "all.hosts."+nodeStatus.NodeName)
				common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts."+nodeStatus.NodeName)
				common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_node.hosts."+nodeStatus.NodeName)
				common.MapDelete(inventoryNew, "all.children.target.children.etcd.hosts."+nodeStatus.NodeName)
			}
		}

		inventoryNewContent, err := yaml.Marshal(inventoryNew)
		if err != nil {
			return "", err
		}
		// logrus.Trace(string(inventoryNewContent))
		if err := ioutil.WriteFile(inventoryPath, inventoryNewContent, 0655); err != nil {
			logrus.Trace(err)
		}

		return "\n" + message, nil
	}

	playbook := common.MapGet(resourcePackage, "data.supported_playbooks.remove_node").(string)

	cmd := command.Execute{
		OwnerType: "cluster",
		OwnerName: req.Cluster,
		Cmd:       "ansible-playbook",
		Args: func(execute_dir string) []string {
			if req.VVV {
				return []string{"-i", execute_dir + "/inventory.yaml", playbook, "-vvv", "--fork", strconv.Itoa(req.Fork), "-e", nodes, "-e", "allow_ungraceful_removal=true", "-e", "reset_nodes=true"}
			}
			return []string{"-i", execute_dir + "/inventory.yaml", playbook, "--fork", strconv.Itoa(req.Fork), "-e", nodes, "-e", "allow_ungraceful_removal=true", "-e", "reset_nodes=true"}
		},
		Dir:      resourcePackagePathForInventory(inventory),
		Type:     "remove_node",
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
