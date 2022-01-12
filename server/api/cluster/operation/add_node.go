package operation

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AddNodeRequest struct {
	InstallClusterRequest
	Nodes string `json:"nodes_to_add"`
}

func AddNode(c *gin.Context) {
	var req AddNodeRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)

	inventory, resourcePackage, err := updateResourcePackageVarsToInventory(req.Cluster)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to process inventory", err)
		return
	}
	common.MapSet(inventory, "all.vars.kuboardspray_no_log", !req.Verbose)

	// 判断待添加节点是否有控制节点或者 etcd 节点
	nodes := strings.Split(req.Nodes, ",")
	nodesToAdd := ""
	includesControlPlane := false
	includesEtcd := false
	controlPlaneHosts := common.MapGet(inventory, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts").(map[string]interface{})
	etcdHosts := common.MapGet(inventory, "all.children.target.children.etcd.hosts").(map[string]interface{})
	for _, key := range nodes {
		nodesToAdd += key + ","
		if controlPlaneHosts[key] != nil {
			includesControlPlane = true
		}
		if etcdHosts[key] != nil {
			includesEtcd = true
		}
	}
	nodesToAdd = strings.Trim(nodesToAdd, ",")

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

		PostProcessInventory(req.Cluster, "add_node")

		return "\n" + message, nil
	}

	cmd := command.Execute{
		OwnerType: "cluster",
		OwnerName: req.Cluster,
		Cmd:       "ansible-playbook",
		Args: func(execute_dir string) []string {
			playbook := common.MapGet(resourcePackage, "data.supported_playbooks.add_node").(string)
			logrus.Trace("add_nodes: ", nodesToAdd)
			result := []string{"-i", execute_dir + "/inventory.yaml", playbook, "--fork", strconv.Itoa(req.Fork), "-e", "node=" + nodesToAdd}
			if req.ExcludeNodes != "" {
				result = append(result, "--limit", req.ExcludeNodes)
			}
			if includesControlPlane || includesEtcd {
				playbook = common.MapGet(resourcePackage, "data.supported_playbooks.install_cluster").(string)
				if includesEtcd {
					nodesToAdd = "kube_control_plane,etcd," + nodesToAdd
				} else {
					nodesToAdd = "kube_control_plane," + nodesToAdd
				}
				if req.ExcludeNodes != "" {
					nodesToAdd += "," + req.ExcludeNodes
				}
				result = []string{"-i", execute_dir + "/inventory.yaml", playbook, "--fork", strconv.Itoa(req.Fork), "--limit", nodesToAdd}
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
