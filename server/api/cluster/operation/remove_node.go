package operation

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/eip-work/kuboard-spray/api/cluster"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type RemoveNodeRequest struct {
	InstallClusterRequest
	Nodes                  string `json:"nodes_to_remove"`
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

	nodes := "node=" + req.Nodes

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
			result := []string{
				"-i", execute_dir + "/inventory.yaml", playbook,
				"--fork", strconv.Itoa(req.Fork),
				"-e", nodes,
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
