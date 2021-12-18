package cluster

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
)

func InstallCluster(c *gin.Context) {

	var req GetClusterRequest
	c.ShouldBindUri(&req)

	inventoryPath := ClusterInventoryPath(req.Cluster)
	inventory, err := common.ParseYamlFile(inventoryPath)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to read inventory.", err)
		return
	}
	resourcePackagePath := constants.GET_DATA_RESOURCE_DIR() + "/" + common.MapGet(inventory, "all.hosts.localhost.kuboardspray_resource_package").(string)

	resourcePackage, err := common.ParseYamlFile(resourcePackagePath + "/package.yaml")
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to read package.", err)
	}

	// start merge resourcePackage info into inventory

	common.MapSet(inventory, "all.children.k8s_cluster.vars.kube_version", common.MapGet(resourcePackage, "kubernetes.kube_version"))
	common.MapSet(inventory, "all.children.k8s_cluster.vars.image_arch", common.MapGet(resourcePackage, "kubernetes.image_arch"))
	common.MapSet(inventory, "all.children.k8s_cluster.vars.gcr_image_repo", common.MapGet(resourcePackage, "kubernetes.gcr_image_repo"))
	common.MapSet(inventory, "all.children.k8s_cluster.vars.kube_image_repo", common.MapGet(resourcePackage, "kubernetes.kube_image_repo"))
	common.MapSet(inventory, "all.children.k8s_cluster.vars.container_manager", common.MapGet(resourcePackage, "docker_engine.container_manager"))
	dockerVersion := common.MapGet(resourcePackage, "docker_engine.target").(string)
	common.MapSet(inventory, "all.children.k8s_cluster.vars."+dockerVersion, common.MapGet(resourcePackage, "docker_engine.version"))
	common.MapSet(inventory, "all.children.k8s_cluster.vars.container_manager", common.MapGet(resourcePackage, "docker_engine.container_manager"))
	common.MapSet(inventory, "all.children.etcd.vars.etcd_version", common.MapGet(resourcePackage, "etcd.etcd_version"))
	common.MapSet(inventory, "all.children.etcd.vars.etcd_deployment_type", common.MapGet(resourcePackage, "etcd.etcd_deployment_type"))

	dependencies := resourcePackage["dependency"].([]interface{})
	for _, d := range dependencies {
		dependency := d.(map[string]interface{})
		field := dependency["target"].(string)
		version := dependency["version"]
		common.MapSet(inventory, "all.children.k8s_cluster.vars."+field, version)
	}

	cni := common.MapGet(inventory, "all.children.k8s_cluster.vars.kube_network_plugin").(string)
	var cni_dependency map[string]interface{}
	for _, c := range resourcePackage["cni"].([]interface{}) {
		cni_option := c.(map[string]interface{})
		if cni_option["name"] == cni {
			cni_dependency = cni_option
		}
	}
	common.MapSet(inventory, "all.children.k8s_cluster.vars."+cni_dependency["target"].(string), cni_dependency["version"])

	addons := resourcePackage["addons"].([]interface{})
	for _, a := range addons {
		addon := a.(map[string]interface{})
		enabled := common.MapGet(inventory, "all.children.k8s_cluster.vars."+addon["target"].(string))
		if enabled == true {
			for key, value := range addon["params"].(map[string]interface{}) {
				common.MapSet(inventory, "all.children.k8s_cluster.vars."+key, value)
			}
		}
	}

	// end merge resourcePackage info into inventory

	// common.MapSet(inventory, "all.vars.download_keep_remote_cache", true)
	common.MapSet(inventory, "all.vars.download_run_once", true)
	common.MapSet(inventory, "all.vars.download_localhost", true)
	common.MapSet(inventory, "all.vars.download_always_pull", false)
	common.MapSet(inventory, "all.vars.download_cache_dir", resourcePackagePath+"/kubespray_cache")
	common.MapSet(inventory, "all.vars.ansible_ssh_common_args", "-o StrictHostKeyChecking=no")
	common.MapSet(inventory, "all.vars.kuboardspray_cluster_dir", constants.GET_DATA_INVENTORY_DIR()+"/"+req.Cluster)

	command := command.Execute{
		Cluster: req.Cluster,
		Cmd:     "ansible-playbook",
		Args:    func(run_dir string) []string { return []string{"-i", run_dir + "/inventory.yaml", "cluster.yml"} },
		Type:    "install",
		Prepare: func(run_dir string) error { return common.SaveYamlFile(run_dir+"/inventory.yaml", inventory) },
		Dir:     resourcePackagePath + "/kubespray",
	}

	if err := command.Exec(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Faild to InstallCluster. ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
	})

}
