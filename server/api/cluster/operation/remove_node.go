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

type RemoveNodeRequest struct {
	InstallClusterRequest
	NodesToRemove          string `json:"nodes_to_remove"`
	ResetNodes             bool   `json:"reset_nodes"`
	AllowUngracefulRemoval bool   `json:"allow_ungraceful_removal"`
	DrainTimeout           string `json:"drain_timeout"`
	DrainOutRetries        string `json:"drain_retries"`
	DrainRetryDelaySeconds string `json:"drain_retry_delay_seconds"`
	DrainGracePeriod       string `json:"drain_grace_period"`
}

func RemoveNode(c *gin.Context) {
	var req RemoveNodeRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)

	inventory, resourcePackage, err := updateResourcePackageVarsToInventory(req.Cluster)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to process inventory", err)
		return
	}
	common.MapSet(inventory, "all.vars.kuboardspray_no_log", !req.Verbose)

	// 找出各类待删除的节点
	all_nodes_to_remove := strings.Split(req.NodesToRemove, ",")
	control_plane_to_remove := []string{}
	etcd_member_to_remove := []string{}

	for _, node := range all_nodes_to_remove {
		if common.MapGet(inventory, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts") != nil {
			control_plane_to_remove = append(control_plane_to_remove, node)
		}
		if common.MapGet(inventory, "all.children.target.children.etcd.hosts."+node) != nil {
			etcd_member := common.MapGet(inventory, "all.children.target.children.etcd.hosts."+node+".etcd_member_name").(string)
			etcd_member_to_remove = append(etcd_member_to_remove, etcd_member)
		}
	}

	postExec := func(status command.ExecuteExitStatus) (string, error) {

		inventoryPath := cluster.ClusterInventoryYamlPath(req.Cluster)
		inventoryNew, _ := common.ParseYamlFile(inventoryPath)
		nodesInK8s, _ := getNodesInK8s(req.Cluster)
		membersInEtcd, _ := getMembersInEtcd(req.Cluster)

		countRemoved := len(arraySubtract(all_nodes_to_remove, nodesInK8s))

		message := "\n"
		if countRemoved == len(all_nodes_to_remove) {
			// 所有节点添加成功
			message += "\033[32m[ " + strconv.Itoa(len(all_nodes_to_remove)) + " nodes are already removed, please release the machine." + " ]\033[0m \n"
			message += "\033[32m[ " + strconv.Itoa(len(all_nodes_to_remove)) + " 节点已从 Kubernetes 集群中删除，请释放该资源" + " ]\033[0m \n"

			message = newFunction(etcd_member_to_remove, membersInEtcd, message)
		} else if countRemoved > 0 {
			// 部分节点添加成功
			message += "\033[33m[ Intended to remove " + strconv.Itoa(len(all_nodes_to_remove)) + " nodes, and " + strconv.Itoa(countRemoved) + " of them are removed successfully." + " ]\033[0m \n"
			message += "\033[33m[ 计划删除 " + strconv.Itoa(len(all_nodes_to_remove)) + " 个节点，其中 " + strconv.Itoa(countRemoved) + " 个节点删除成功。" + " ]\033[0m \n"

			message = newFunction(etcd_member_to_remove, membersInEtcd, message)
		} else {
			// 添加节点失败
			message += "\033[31m\033[01m\033[05m[ " + "Failed to remove node. Please review the logs and fix the problem." + " ]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[ " + "删除节点失败，请回顾日志，找到错误信息，并解决问题后，再次尝试。" + " ]\033[0m \n"
		}

		if len(arraySubtract(control_plane_to_remove, nodesInK8s)) < len(control_plane_to_remove) {
			// 有控制节点被删除
			message += "\n\033[31m\033[01m\033[05m[ " + "Apiserver list changed, it's required to \"Update apiserver list in loadbalancer\"." + " ]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[ " + "Apiserver 列表发生变化，请在集群页面执行操作 \"更新负载均衡中 apiserver 列表\"." + " ]\033[0m \n"

			common.MapSet(inventoryNew, "all.hosts.localhost.kuboardspray_sync_nginx_config", true)
		}
		removedEtcdMembers := arraySubtract(etcd_member_to_remove, membersInEtcd)
		if len(removedEtcdMembers) > 0 {
			// 有 ETCD 节点被删除
			message += "\n\033[31m\033[01m\033[05m[ " + strconv.Itoa(len(removedEtcdMembers)) + " etcd members are removed, it's important for you to update --etcd-servers param in /etc/kubernetes/manifests/kube-apiserver.yaml" + " ]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[ " + strconv.Itoa(len(removedEtcdMembers)) + " etcd 成员被删除，请返回集群管理界面，执行操作 \"更新 apiserver 中 etcd 连接参数\"" + " ]\033[0m \n"
			common.MapSet(inventoryNew, "all.hosts.localhost.kuboardspray_sync_etcd_address", true)
		}

		for _, key := range arraySubtract(all_nodes_to_remove, nodesInK8s) {
			common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts."+key)
			common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_node.hosts."+key)
			common.MapDelete(inventoryNew, "all.children.target.children.etcd.hosts."+key)
			common.MapDelete(inventoryNew, "all.hosts."+key)
		}

		inventoryNewContent, _ := yaml.Marshal(inventoryNew)

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
			result := []string{
				"-i", execute_dir + "/inventory.yaml", playbook,
				"--fork", strconv.Itoa(req.Fork),
				"-e", "node=" + req.NodesToRemove,
				"-e", "reset_nodes=" + strconv.FormatBool(req.ResetNodes),
				"-e", "allow_ungraceful_removal=" + strconv.FormatBool(req.AllowUngracefulRemoval),
				"-e", "drain_grace_period=" + req.DrainGracePeriod,
				"-e", "drain_timeout=" + req.DrainTimeout,
				"-e", "drain_retries=" + req.DrainOutRetries,
				"-e", "drain_retry_delay_seconds=" + req.DrainRetryDelaySeconds,
			}
			if req.VVV {
				result = append(result, "-vvv")
			}
			return result
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

func newFunction(etcd_member_to_remove []string, membersInEtcd []string, message string) string {
	removedEtcdMembers := arraySubtract(etcd_member_to_remove, membersInEtcd)
	if len(removedEtcdMembers) != len(etcd_member_to_remove) {
		message += "\n\033[31m\033[01m\033[05m[ " + strconv.Itoa(len(etcd_member_to_remove)) + " etcd member should be removed, but only " + strconv.Itoa(len(removedEtcdMembers)) + " are removed, please verify and fix, or the cluster may misfunction." + " ]\033[0m \n"
		message += "\n\033[31m\033[01m\033[05m[ " + strconv.Itoa(len(etcd_member_to_remove)) + " 个 etcd 成员应该被删除，但是只有 " + strconv.Itoa(len(removedEtcdMembers)) + " 个删除成功，请检查 etcd 的状态并修复，否则集群可能出现不可预知的故障。" + " ]\033[0m \n"
	}
	return message
}
