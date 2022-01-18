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

type AddNodeRequest struct {
	InstallClusterRequest
	NodesToAdd string `json:"nodes_to_add"`
}

func AddNode(c *gin.Context) {
	var req AddNodeRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)

	inventory, resourcePackage, err := updateResourcePackageVarsToInventory(req.Cluster, req.DiscoveredInterpreterPython)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to process inventory", err)
		return
	}
	common.MapSet(inventory, "all.vars.kuboardspray_no_log", !req.Verbose)

	// 判断待添加节点是否有控制节点或者 etcd 节点
	nodes_to_add := strings.Split(req.NodesToAdd, ",")
	nodesToAdd := ""
	includesControlPlane := false
	includesEtcd := false
	controlPlaneHosts := common.MapGet(inventory, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts").(map[string]interface{})
	etcdHosts := common.MapGet(inventory, "all.children.target.children.etcd.hosts").(map[string]interface{})
	for _, key := range nodes_to_add {
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

		var message string

		nodesInK8s, _ := getNodesInK8s(req.Cluster)
		countAdded := len(nodes_to_add) - len(arraySubtract(nodes_to_add, nodesInK8s))

		if len(arraySubtract(nodes_to_add, nodesInK8s)) == 0 {
			message += "\033[32m[ " + strconv.Itoa(countAdded) + " nodes are already added to cluster, please release the machine." + " ]\033[0m \n"
			message += "\033[32m[ " + strconv.Itoa(countAdded) + " 个节点已添加到 Kubernetes 集群，请返回集群详情页查看。" + " ]\033[0m \n"
		} else if countAdded > 0 {
			message += "\033[33m[ Intended to add " + strconv.Itoa(len(nodes_to_add)) + " nodes, and " + strconv.Itoa(countAdded) + " of them are added successfully." + " ]\033[0m \n"
			message += "\033[33m[ 计划添加 " + strconv.Itoa(len(nodes_to_add)) + " 个节点，其中 " + strconv.Itoa(countAdded) + " 个节点添加成功。" + " ]\033[0m \n"
		} else {
			message += "\033[31m\033[01m\033[05m[" + "Failed to add node. Please review the logs and fix the problem." + "]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[" + "添加节点失败，请回顾日志，找到错误信息，并解决问题后，再次尝试。" + "]\033[0m \n"
		}

		inventoryPath := cluster.ClusterInventoryYamlPath(req.Cluster)
		inventoryNew, _ := common.ParseYamlFile(inventoryPath)

		addedControlePlane := false
		for _, node := range nodesInK8s {
			if common.MapGet(inventoryNew, "all.hosts."+node+".kuboardspray_node_action") == "add_node" {
				if common.MapGet(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts."+node) != nil {
					addedControlePlane = true
				}
				common.MapDelete(inventoryNew, "all.hosts."+node+".kuboardspray_node_action")
				common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts."+node+".kuboardspray_node_action")
				common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_node.hosts."+node+".kuboardspray_node_action")
				common.MapDelete(inventoryNew, "all.children.target.children.etcd.hosts."+node+".kuboardspray_node_action")
			}
		}
		if addedControlePlane {
			common.MapSet(inventoryNew, "all.hosts.localhost.kuboardspray_sync_nginx_config", true)
			common.MapSet(inventoryNew, "all.hosts.localhost.kuboardspray_sync_etcd_address", false)
			message += "\033[31m\033[01m\033[05m[ " + "Apiserver list changed, it's required to \"Update apiserver list in loadbalancer\"." + " ]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[ " + "Apiserver 列表发生变化，请在集群页面执行操作 \"更新负载均衡中 apiserver 列表\"." + " ]\033[0m \n"
		}
		inventoryNewContent, _ := yaml.Marshal(inventoryNew)

		if err := ioutil.WriteFile(inventoryPath, inventoryNewContent, 0655); err != nil {
			logrus.Trace(err)
		}

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
					nodesToAdd = "etcd," + nodesToAdd
				}
				nodesToAdd = "kube_control_plane," + nodesToAdd
				if req.ExcludeNodes != "" {
					nodesToAdd += "," + req.ExcludeNodes
				}
				result = []string{"-i", execute_dir + "/inventory.yaml", playbook, "--fork", strconv.Itoa(req.Fork), "--limit", nodesToAdd, "-e", "etcd_retries=15"}
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
