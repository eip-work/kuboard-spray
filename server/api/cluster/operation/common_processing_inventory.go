package operation

import (
	"strings"

	"github.com/eip-work/kuboard-spray/api/cluster"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/sirupsen/logrus"
)

func resourcePackagePathForInventory(inventory map[string]interface{}) string {
	return constants.GET_DATA_RESOURCE_DIR() + "/" + common.MapGet(inventory, "all.hosts.localhost.kuboardspray_resource_package").(string) + "/content"
}

func updateResourcePackageVarsToInventory(req InstallClusterRequest) (map[string]interface{}, map[string]interface{}, error) {
	inventoryPath := cluster.ClusterInventoryYamlPath(req.Cluster)
	inventory, err := common.ParseYamlFile(inventoryPath)
	if err != nil {
		return nil, nil, err
	}
	resourcePackagePath := resourcePackagePathForInventory(inventory)

	resourcePackage, err := common.ParseYamlFile(resourcePackagePath + "/package.yaml")
	if err != nil {
		return nil, nil, err
	}

	common.MapSet(inventory, "all.vars.kuboardspray_no_log", !req.Verbose)

	// 设置 discovered_interpreter_python
	for k, v := range req.DiscoveredInterpreterPython {
		common.MapSet(inventory, "all.hosts."+k+".ansible_python_interpreter", v)
	}

	// >>>> 设置资源包相关参数

	// 设置 kubernetes 版本信息
	common.MapSet(inventory, "all.children.target.children.k8s_cluster.vars.kube_version", common.MapGet(resourcePackage, "data.kubernetes.kube_version"))
	common.MapSet(inventory, "all.children.target.children.k8s_cluster.vars.image_arch", common.MapGet(resourcePackage, "data.kubernetes.image_arch"))
	common.MapSet(inventory, "all.children.target.children.k8s_cluster.vars.gcr_image_repo", common.MapGet(resourcePackage, "data.kubernetes.gcr_image_repo"))
	common.MapSet(inventory, "all.children.target.children.k8s_cluster.vars.kube_image_repo", common.MapGet(resourcePackage, "data.kubernetes.kube_image_repo"))

	// 设置容器引擎相关参数
	container_manager := common.MapGet(inventory, "all.children.target.vars.container_manager")
	containerMangerObjArray := common.MapGet(resourcePackage, "data.container_engine").([]interface{})
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
	common.MapSet(inventory, "all.children.target.children.etcd.vars.etcd_version", common.MapGet(resourcePackage, "data.etcd.etcd_version"))
	params := common.MapGet(resourcePackage, "data.etcd.etcd_params")
	if params != nil {
		etcdParams := params.(map[string]interface{})
		for k, v := range etcdParams {
			common.MapSet(inventory, "all.children.target.children.etcd.vars."+k, v)
			logrus.Trace("set inventory: ", "all.children.target.children.etcd.vars."+k, v)
		}
	}

	common.MapSet(inventory, "all.children.target.children.etcd.vars.etcd_config_dir", "/etc/ssl/etcd")
	common.MapSet(inventory, "all.children.target.children.etcd.vars.etcd_cert_dir", "{{ etcd_config_dir }}/ssl")

	// 设置依赖组件版本信息
	dependencies := common.MapGet(resourcePackage, "data.dependency").([]interface{}) // resourcePackage["dependency"].([]interface{})
	for _, d := range dependencies {
		dependency := d.(map[string]interface{})
		field := dependency["target"].(string)
		version := dependency["version"]
		common.MapSet(inventory, "all.children.target.children.k8s_cluster.vars."+field, version)
	}

	// 设置网络插件信息  FIXME 调整 package.yaml 的格式
	np := common.MapGet(inventory, "all.children.target.children.k8s_cluster.vars.kube_network_plugin").(string)
	var np_dependency map[string]interface{}
	network_plugins := common.MapGet(resourcePackage, "data.network_plugin").([]interface{})
	for _, c := range network_plugins {
		np_option := c.(map[string]interface{})
		if np_option["name"] == np {
			np_dependency = np_option
		}
	}
	for key, value := range np_dependency["params"].(map[string]interface{}) {
		common.MapSet(inventory, "all.children.target.children.k8s_cluster.vars."+key, value)
	}

	// 设置可选组件参数
	addons := common.MapGet(resourcePackage, "data.addon").([]interface{})
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
					return nil, nil, err
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

	common.MapSet(inventory, "all.vars.download_keep_remote_cache", false)
	common.MapSet(inventory, "all.vars.download_run_once", true)
	common.MapSet(inventory, "all.vars.download_localhost", true)
	common.MapSet(inventory, "all.vars.download_force_cache", true)
	common.MapSet(inventory, "all.vars.download_always_pull", false)
	common.MapSet(inventory, "all.vars.download_cache_dir", resourcePackagePath+"/kubespray_cache")
	// common.MapSet(inventory, "all.vars.ansible_ssh_common_args", "-o StrictHostKeyChecking=no")
	common.MapSet(inventory, "all.vars.kuboardspray_cluster_dir", constants.GET_DATA_CLUSTER_DIR()+"/"+req.Cluster)
	common.MapSet(inventory, "all.children.target.vars.disable_service_firewall", true)
	common.MapSet(inventory, "all.children.target.vars.ansible_python_interpreter", "auto")

	return inventory, resourcePackage, nil
}
