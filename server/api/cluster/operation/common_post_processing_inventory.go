package operation

import (
	"io/ioutil"
	"strings"

	"github.com/eip-work/kuboard-spray/api/cluster"
	"github.com/eip-work/kuboard-spray/api/cluster/state"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func contains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func PostProcessInventory(clusterName string, action string) (string, error) {
	inventoryPath := cluster.ClusterInventoryYamlPath(clusterName)
	inventoryNew, _ := common.ParseYamlFile(inventoryPath)

	kubectlResult, err := state.ExecuteShellOnControlPlane(clusterName, "kubectl get nodes -o jsonpath=\"{.items[*].metadata.name}\"")

	if err != nil {
		return "", err
	}

	nodeStr := strings.Trim(kubectlResult.Stdout, "\n")

	nodesInK8s := strings.Split(nodeStr, " ")
	logrus.Trace("Nodes in kubernetes: ", nodesInK8s)

	message := ""

	common.MapSet(inventoryNew, "all.hosts.localhost.kuboardspray_sync_nginx_config", false)

	if action == "add_node" || action == "install_cluster" {
		includeControlPlane := false
		for _, node := range nodesInK8s {
			if common.MapGet(inventoryNew, "all.hosts."+node+".kuboardspray_node_action") == "add_node" {
				controlPlanes := common.MapGet(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts").(map[string]interface{})
				if controlPlanes[node] != nil {
					includeControlPlane = true
				}
				common.MapDelete(inventoryNew, "all.hosts."+node+".kuboardspray_node_action")
				common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts."+node+".kuboardspray_node_action")
				common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_node.hosts."+node+".kuboardspray_node_action")
				common.MapDelete(inventoryNew, "all.children.target.children.etcd.hosts."+node+".kuboardspray_node_action")
			}
		}
		if action == "add_node" {
			if includeControlPlane {
				common.MapSet(inventoryNew, "all.hosts.localhost.kuboardspray_sync_nginx_config", true)
				message += "\033[31m\033[01m\033[05m[ " + "Apiserver list changed, it's required to update all nginx proxy in kube_nodes." + " ]\033[0m \n"
				message += "\033[31m\033[01m\033[05m[ " + "Apiserver 列表发生变化，请在集群页面更新所有工作节点的 nginx proxy 配置." + " ]\033[0m \n"
			}
		}
	}
	if action == "remove_node" {
		nodesInvent := common.MapGet(inventoryNew, "all.hosts").(map[string]interface{})
		includeControlPlane := false
		for key, nodeInfo := range nodesInvent {
			n := nodeInfo.(map[string]interface{})
			if n["kuboardspray_node_action"] == "remove_node" && !contains(nodesInK8s, key) {
				controlPlanes := common.MapGet(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts").(map[string]interface{})
				if controlPlanes[key] != nil {
					includeControlPlane = true
				}
				common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts."+key)
				common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_node.hosts."+key)
				common.MapDelete(inventoryNew, "all.children.target.children.etcd.hosts."+key)
				common.MapDelete(inventoryNew, "all.hosts."+key)
			}
		}
		if includeControlPlane {
			common.MapSet(inventoryNew, "all.hosts.localhost.kuboardspray_sync_nginx_config", true)
			message += "\n\033[31m\033[01m\033[05m[ " + "Apiserver list changed, it's required to update all nginx proxy in kube_nodes." + " ]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[ " + "Apiserver 列表发生变化，请在集群页面更新所有工作节点的 nginx proxy 配置." + " ]\033[0m \n"
		}
	}

	// 对比 etcd 成员
	etcdInvent := common.MapGet(inventoryNew, "all.children.target.children.etcd.hosts").(map[string]interface{})

	etcdOut, err := state.ExecuteShellOnETCD(clusterName, "etcdctl member list --write-out=json")
	if err == nil {
		count := 0
		for _, node := range etcdOut.StdoutObj["members"].([]interface{}) {
			nodeObj := node.(map[string]interface{})
			nodeName := nodeObj["name"].(string)
			logrus.Trace("etcdName: ", nodeName)
			for _, v := range etcdInvent {
				etcdvar := v.(map[string]interface{})
				if etcdvar["etcd_member_name"] == nodeName {
					count++
				}
			}
		}

		logrus.Trace(etcdOut.Stdout)
		logrus.Trace("current etcd_count: ", count)

		etcdMaxCount := common.MapGet(inventoryNew, "all.hosts.localhost.kuboardspray_etcd_max_count")
		maxCount := 0
		if etcdMaxCount != nil {
			maxCount = etcdMaxCount.(int)
			logrus.Trace("history etcd_max_count: ", maxCount)
		}
		if maxCount < count {
			maxCount = count
		}
		common.MapSet(inventoryNew, "all.hosts.localhost.kuboardspray_etcd_max_count", maxCount)
	}

	inventoryNewContent, err := yaml.Marshal(inventoryNew)
	if err != nil {
		return "", err
	}
	// logrus.Trace(string(inventoryNewContent))
	if err := ioutil.WriteFile(inventoryPath, inventoryNewContent, 0655); err != nil {
		logrus.Trace(err)
	}

	return message, nil
}
