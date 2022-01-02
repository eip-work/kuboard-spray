package cluster

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type InstallClusterRequest struct {
	Cluster string `uri:"cluster" binding:"required"`
	Fork    int    `json:"fork"`
	Verbose bool   `json:"verbose"`
	VVV     bool   `json:"vvv"`
}

func InstallCluster(c *gin.Context) {

	var req InstallClusterRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)

	inventoryPath := ClusterInventoryYamlPath(req.Cluster)
	inventory, err := common.ParseYamlFile(inventoryPath)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to read inventory.", err)
		return
	}
	resourcePackagePath := constants.GET_DATA_RESOURCE_DIR() + "/" + common.MapGet(inventory, "all.hosts.localhost.kuboardspray_resource_package").(string) + "/content"

	resourcePackage, err := common.ParseYamlFile(resourcePackagePath + "/package.yaml")
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to read package.", err)
	}

	// >>>> 设置资源包相关参数

	// 设置 kubernetes 版本信息
	common.MapSet(inventory, "all.children.target.children.k8s_cluster.vars.kube_version", common.MapGet(resourcePackage, "kubernetes.kube_version"))
	common.MapSet(inventory, "all.children.target.children.k8s_cluster.vars.image_arch", common.MapGet(resourcePackage, "kubernetes.image_arch"))
	common.MapSet(inventory, "all.children.target.children.k8s_cluster.vars.gcr_image_repo", common.MapGet(resourcePackage, "kubernetes.gcr_image_repo"))
	common.MapSet(inventory, "all.children.target.children.k8s_cluster.vars.kube_image_repo", common.MapGet(resourcePackage, "kubernetes.kube_image_repo"))

	// 设置容器引擎相关参数
	container_manager := common.MapGet(inventory, "all.children.target.vars.container_manager")
	containerMangerObjArray := common.MapGet(resourcePackage, "container_engine").([]interface{})
	var containerManagerObj map[string]interface{}
	for _, cmo := range containerMangerObjArray {
		cmObj := cmo.(map[string]interface{})
		if cmObj["container_manager"] == container_manager {
			containerManagerObj = cmObj
		}
	}
	for key, value := range containerManagerObj["params"].(map[string]interface{}) {
		common.MapSet(inventory, "all.children.target.vars."+key, value)
	}

	// 设置 etcd 版本信息
	common.MapSet(inventory, "all.children.target.children.etcd.vars.etcd_version", common.MapGet(resourcePackage, "etcd.etcd_version"))

	// 设置依赖组件版本信息
	dependencies := resourcePackage["dependency"].([]interface{})
	for _, d := range dependencies {
		dependency := d.(map[string]interface{})
		field := dependency["target"].(string)
		version := dependency["version"]
		common.MapSet(inventory, "all.children.target.children.k8s_cluster.vars."+field, version)
	}

	// 设置网络插件信息  FIXME 调整 package.yaml 的格式
	cni := common.MapGet(inventory, "all.children.target.children.k8s_cluster.vars.kube_network_plugin").(string)
	var cni_dependency map[string]interface{}
	for _, c := range resourcePackage["cni"].([]interface{}) {
		cni_option := c.(map[string]interface{})
		if cni_option["name"] == cni {
			cni_dependency = cni_option
		}
	}
	common.MapSet(inventory, "all.children.target.children.k8s_cluster.vars."+cni_dependency["target"].(string), cni_dependency["version"])

	// 设置可选组件参数
	addons := resourcePackage["addons"].([]interface{})
	for _, a := range addons {
		addon := a.(map[string]interface{})
		enabled := common.MapGet(inventory, "all.children.target.children.k8s_cluster.vars."+addon["target"].(string))
		if enabled == true {
			for key, value := range addon["params"].(map[string]interface{}) {
				common.MapSet(inventory, "all.children.target.children.k8s_cluster.vars."+key, value)
			}
		}
	}

	// 设置软件源相关参数
	targetVars := common.MapGet(inventory, "all.children.target.vars").(map[string]interface{})
	for key, value := range targetVars {
		if strings.Index(key, "kuboardspray_repo_") == 0 { // 忽略 kuboardspray_repo_ 和 kuboardspray_repo_docker_ 的差异
			v := value.(string)
			if v == "AS_IS" {
				logrus.Trace("使用操作系统已经配置的软件源 -> ", key, " : ", value)
			} else {
				repo, err := common.ParseYamlFile(constants.GET_DATA_DIR() + "/mirror/" + v + "/status.yaml")
				if err != nil {
					common.HandleError(c, http.StatusInternalServerError, "cannot read repo. ", err)
					return
				}
				params := repo["params"].(map[string]interface{})
				logrus.Trace("配置软件源参数  -> ", key, " : ", value)
				for k, v := range params {
					targetVars[k] = v
					logrus.Trace("    ", k, " : ", v)
				}
			}
		}
	}

	// <<<< 设置资源包相关参数

	common.MapSet(inventory, "all.vars.kuboardspray_no_log", !req.Verbose)
	common.MapSet(inventory, "all.vars.download_keep_remote_cache", false)
	common.MapSet(inventory, "all.vars.download_run_once", true)
	common.MapSet(inventory, "all.vars.download_localhost", true)
	common.MapSet(inventory, "all.vars.download_force_cache", true)
	common.MapSet(inventory, "all.vars.download_always_pull", false)
	common.MapSet(inventory, "all.vars.download_cache_dir", resourcePackagePath+"/kubespray_cache")
	common.MapSet(inventory, "all.vars.ansible_ssh_common_args", "-o StrictHostKeyChecking=no")
	common.MapSet(inventory, "all.vars.kuboardspray_cluster_dir", constants.GET_DATA_CLUSTER_DIR()+"/"+req.Cluster)
	common.MapSet(inventory, "all.children.target.vars.disable_service_firewall", true)
	common.MapSet(inventory, "all.children.target.vars.ansible_python_interpreter", "auto")

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
		return "\n" + message, nil
	}

	playbook := common.MapGet(resourcePackage, "metadata.supported_playbooks.install_cluster").(string)

	cmd := command.Execute{
		OwnerType: "cluster",
		OwnerName: req.Cluster,
		Cmd:       "ansible-playbook",
		Args: func(execute_dir string) []string {
			if req.VVV {
				return []string{"-i", execute_dir + "/inventory.yaml", playbook, "-vvv", "--fork", strconv.Itoa(req.Fork)}
			}
			return []string{"-i", execute_dir + "/inventory.yaml", playbook, "--fork", strconv.Itoa(req.Fork)}
		},
		Dir:      resourcePackagePath,
		Type:     "install",
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
