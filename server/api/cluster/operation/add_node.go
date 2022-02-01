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

type AddNodeRequest struct {
	OperationCommonRequest
	NodesToAdd string `json:"nodes_to_add"`
}

func AddNode(c *gin.Context) {
	var req AddNodeRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)
	req.Operation = "add_node"

	inventory, resourcePackage, err := updateResourcePackageVarsToInventory(req.OperationCommonRequest)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to process inventory", err)
		return
	}

	// 判断待添加节点是否有控制节点或者 etcd 节点
	nodes_to_add := strings.Split(req.NodesToAdd, ",")
	nodesToAdd := ""

	controlPlaneHosts := common.MapGet(inventory, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts").(map[string]interface{})
	kubeNodeHosts := common.MapGet(inventory, "all.children.target.children.k8s_cluster.children.kube_node.hosts").(map[string]interface{})
	etcdHosts := common.MapGet(inventory, "all.children.target.children.etcd.hosts").(map[string]interface{})

	kube_nodes_to_add := []string{}
	kube_control_planes_to_add := []string{}
	etcd_members_to_add := []string{}
	for _, key := range nodes_to_add {
		nodesToAdd += key + ","
		if controlPlaneHosts[key] != nil {
			kube_control_planes_to_add = append(kube_control_planes_to_add, key)
		}
		if kubeNodeHosts[key] != nil {
			kube_nodes_to_add = append(kube_nodes_to_add, key)
		}
		if etcdHosts[key] != nil {
			etcd := etcdHosts[key].(map[string]interface{})
			if etcd["etcd_member_name"] == nil {
				common.HandleError(c, http.StatusInternalServerError, "node "+key+"'s etcd_member_name field cannot be empty.", nil)
				return
			}
			etcd_members_to_add = append(etcd_members_to_add, etcd["etcd_member_name"].(string))
		}
	}
	nodesToAdd = strings.Trim(nodesToAdd, ",")

	postExec := func(status command.ExecuteExitStatus) (string, error) {

		var message string

		nodesInK8s, _ := getNodesInK8s(req.Cluster)
		membersInEtcd, _ := getMembersInEtcd(req.Cluster)
		countAddedKubeNode := len(kube_nodes_to_add) - len(arraySubtract(kube_nodes_to_add, nodesInK8s))
		countAddedKubeControlPlane := len(kube_control_planes_to_add) - len(arraySubtract(kube_control_planes_to_add, nodesInK8s))
		countAddedEtcdMembers := len(etcd_members_to_add) - len(arraySubtract(etcd_members_to_add, membersInEtcd))

		inventoryPath := cluster_common.ClusterInventoryYamlPath(req.Cluster)
		inventoryNew, _ := common.ParseYamlFile(inventoryPath)

		if countAddedKubeControlPlane > 0 {
			if countAddedKubeControlPlane == len(kube_control_planes_to_add) {
				message += "\n"
				message += "\033[32m[ " + strconv.Itoa(countAddedKubeControlPlane) + " kube_control_planes are already added to the cluster." + " ]\033[0m \n"
				message += "\033[32m[ " + strconv.Itoa(countAddedKubeControlPlane) + " 个控制节点已添加到 Kubernetes 集群，请返回集群详情页查看。" + " ]\033[0m \n"
			} else {
				message += "\n"
				message += "\033[33m[ Intended to add " + strconv.Itoa(len(kube_control_planes_to_add)) + " kube_control_planes, and " + strconv.Itoa(countAddedKubeControlPlane) + " of them are added successfully." + " ]\033[0m \n"
				message += "\033[33m[ 计划添加 " + strconv.Itoa(len(kube_control_planes_to_add)) + " 个控制节点，其中 " + strconv.Itoa(countAddedKubeControlPlane) + " 个节点添加成功。" + " ]\033[0m \n"
			}
			common.MapSet(inventoryNew, "all.hosts.localhost.kuboardspray_sync_nginx_config", true)
			common.MapSet(inventoryNew, "all.hosts.localhost.kuboardspray_sync_etcd_address", false)
			message += "\n"
			message += "\033[31m\033[01m\033[05m[ " + "Apiserver list changed, it's required to \"Update apiserver list in loadbalancer\"." + " ]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[ " + "Apiserver 列表发生变化，请在集群页面执行操作 \"更新负载均衡中 apiserver 列表\"." + " ]\033[0m \n"
		}

		if countAddedKubeNode > 0 {
			if countAddedKubeNode == len(kube_nodes_to_add) {
				message += "\n"
				message += "\033[32m[ " + strconv.Itoa(countAddedKubeNode) + " kube_nodes are already added to the cluster." + " ]\033[0m \n"
				message += "\033[32m[ " + strconv.Itoa(countAddedKubeNode) + " 个工作节点已添加到 Kubernetes 集群，请返回集群详情页查看。" + " ]\033[0m \n"
			} else {
				message += "\n"
				message += "\033[33m[ Intended to add " + strconv.Itoa(len(kube_nodes_to_add)) + " kube_nodes, and " + strconv.Itoa(countAddedKubeNode) + " of them are added successfully." + " ]\033[0m \n"
				message += "\033[33m[ 计划添加 " + strconv.Itoa(len(kube_nodes_to_add)) + " 个工作节点，其中 " + strconv.Itoa(countAddedKubeNode) + " 个节点添加成功。" + " ]\033[0m \n"
			}
		}

		if countAddedEtcdMembers > 0 {
			if countAddedEtcdMembers == len(etcd_members_to_add) {
				message += "\n"
				message += "\033[32m[ " + strconv.Itoa(countAddedEtcdMembers) + " etcd members are already added to the cluster." + " ]\033[0m \n"
				message += "\033[32m[ " + strconv.Itoa(countAddedEtcdMembers) + " 个 etcd 成员已添加到 etcd 集群，请返回集群详情页查看。" + " ]\033[0m \n"
			} else {
				message += "\n"
				message += "\033[33m[ Intended to add " + strconv.Itoa(len(etcd_members_to_add)) + " etcd members, and " + strconv.Itoa(countAddedEtcdMembers) + " of them are added successfully." + " ]\033[0m \n"
				message += "\033[33m[ 计划添加 " + strconv.Itoa(len(etcd_members_to_add)) + " 个 etcd 成员，其中 " + strconv.Itoa(countAddedEtcdMembers) + " 个成员添加成功。" + " ]\033[0m \n"
			}
			message += "\n"
			message += "\033[32m[ Etcd member list changed, all --etcd-servers in /etc/kubernetes/manifests/kube-apiserver.yaml on control_planes are already refreshed." + " ]\033[0m \n"
			message += "\033[32m[ Etcd 成员列表发生变化，所有控制节点上 /etc/kubernetes/manifests/kube-apiserver.yaml 文件中的参数 --etcd-servers 已更新。 ]\033[0m \n"
		}

		if countAddedKubeNode == 0 && countAddedKubeControlPlane == 0 && countAddedEtcdMembers == 0 {
			message += "\n"
			message += "\033[31m\033[01m\033[05m[" + "Failed to add node. Please review the logs and fix the problem." + "]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[" + "添加节点失败，请回顾日志，找到错误信息，并解决问题后，再次尝试。" + "]\033[0m \n"
		}

		for _, node := range nodesInK8s {
			if common.MapGet(inventoryNew, "all.hosts."+node+".kuboardspray_node_action") == "add_node" {
				common.MapDelete(inventoryNew, "all.hosts."+node+".kuboardspray_node_action")
				common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts."+node+".kuboardspray_node_action")
				common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_node.hosts."+node+".kuboardspray_node_action")
				common.MapDelete(inventoryNew, "all.children.target.children.etcd.hosts."+node+".kuboardspray_node_action")
			}
		}

		for _, member := range membersInEtcd {
			for etcdNodeName, n := range etcdHosts {
				etcdHost := n.(map[string]interface{})
				if etcdHost["etcd_member_name"] == member {
					if common.MapGet(inventoryNew, "all.hosts."+etcdNodeName+".kuboardspray_node_action") == "add_node" {
						common.MapDelete(inventoryNew, "all.hosts."+etcdNodeName+".kuboardspray_node_action")
						common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts."+etcdNodeName+".kuboardspray_node_action")
						common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_node.hosts."+etcdNodeName+".kuboardspray_node_action")
						common.MapDelete(inventoryNew, "all.children.target.children.etcd.hosts."+etcdNodeName+".kuboardspray_node_action")
					}
				}
			}
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
			result := []string{"-i", execute_dir + "/inventory.yaml", playbook, "-e", "node=" + nodesToAdd}
			if len(kube_control_planes_to_add) > 0 || len(etcd_members_to_add) > 0 {
				playbook = common.MapGet(resourcePackage, "data.supported_playbooks.install_cluster").(string)
				if len(etcd_members_to_add) > 0 {
					nodesToAdd = "etcd," + nodesToAdd
				}
				nodesToAdd = "kube_control_plane," + nodesToAdd
				if req.ExcludeNodes != "" {
					nodesToAdd += "," + req.ExcludeNodes
				}
				result = []string{"-i", execute_dir + "/inventory.yaml", playbook, "--limit", nodesToAdd, "-e", "etcd_retries=15"}
			}
			result = appendCommonParams(result, req.OperationCommonRequest, true)
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
