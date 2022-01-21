package operation

import (
	"io/ioutil"
	"net/http"

	"github.com/eip-work/kuboard-spray/api/cluster"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type InstallClusterRequest struct {
	OperationCommonRequest
	SkipDownloads bool `json:"skip_downloads"`
}

func InstallCluster(c *gin.Context) {

	var req InstallClusterRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)

	inventory, resourcePackage, err := updateResourcePackageVarsToInventory(req.OperationCommonRequest)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to process inventory", err)
		return
	}

	postExec := func(status command.ExecuteExitStatus) (string, error) {

		success := status.Success
		var message string
		if success {
			message = "\033[32m[ " + "Kubernetes Cluster has been installed successfully, please go back to the cluster page for information about how to access the cluster." + " ]\033[0m \n"
			message += "\033[32m[ " + "Kubernetes 集群已成功安装，请回到集群详情页面查看如何访问该集群。" + " ]\033[0m \n"
		} else {
			message = "\033[31m\033[01m\033[05m[" + "Failed to install Kubernetes Cluster. Please review the logs and fix the problem." + "]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[" + "集群安装失败，请回顾日志，找到错误信息，并解决问题后，再次尝试。" + "]\033[0m \n"
		}

		inventoryPath := cluster.ClusterInventoryYamlPath(req.Cluster)
		inventoryNew, _ := common.ParseYamlFile(inventoryPath)
		nodesInK8s, _ := getNodesInK8s(req.Cluster)

		for _, node := range nodesInK8s {
			if common.MapGet(inventoryNew, "all.hosts."+node+".kuboardspray_node_action") == "add_node" {
				common.MapDelete(inventoryNew, "all.hosts."+node+".kuboardspray_node_action")
				common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts."+node+".kuboardspray_node_action")
				common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_node.hosts."+node+".kuboardspray_node_action")
				common.MapDelete(inventoryNew, "all.children.target.children.etcd.hosts."+node+".kuboardspray_node_action")
			}
		}
		inventoryNewContent, _ := yaml.Marshal(inventoryNew)

		if err := ioutil.WriteFile(inventoryPath, inventoryNewContent, 0655); err != nil {
			logrus.Trace(err)
		}

		return "\n" + message, nil
	}

	playbook := common.MapGet(resourcePackage, "data.supported_playbooks.install_cluster").(string)

	cmd := command.Execute{
		OwnerType: "cluster",
		OwnerName: req.Cluster,
		Cmd:       "ansible-playbook",
		Args: func(execute_dir string) []string {
			result := []string{"-i", execute_dir + "/inventory.yaml", playbook}
			if req.SkipDownloads {
				result = append(result, "-e", "skip_downloads=true")
			}
			result = appendCommonParams(result, req.OperationCommonRequest, false)
			return result
		},
		Dir:      cluster.ResourcePackagePathForInventory(inventory),
		Type:     "install_cluster",
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
