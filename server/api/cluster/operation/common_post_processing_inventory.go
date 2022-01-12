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

func PostProcessInventory(clusterName string, action string) error {
	inventoryPath := cluster.ClusterInventoryYamlPath(clusterName)
	inventoryNew, _ := common.ParseYamlFile(inventoryPath)

	kubectlResult, err := state.ExecuteShellOnControlPlane(clusterName, "kubectl get nodes -o jsonpath=\"{.items[*].metadata.name}\"")

	if err != nil {
		return err
	}

	nodeStr := strings.Trim(kubectlResult.Stdout, "\n")

	nodesInK8s := strings.Split(nodeStr, " ")
	logrus.Trace("Nodes in kubernetes: ", nodesInK8s)

	if action == "add_node" || action == "install_cluster" {
		for _, node := range nodesInK8s {
			logrus.Trace(node, "--", common.MapGet(inventoryNew, "all.hosts."+node+".kuboardspray_node_action"))
			if common.MapGet(inventoryNew, "all.hosts."+node+".kuboardspray_node_action") == "add_node" {
				common.MapDelete(inventoryNew, "all.hosts."+node+".kuboardspray_node_action")
				common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts."+node+".kuboardspray_node_action")
				common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_node.hosts."+node+".kuboardspray_node_action")
				common.MapDelete(inventoryNew, "all.children.target.children.etcd.hosts."+node+".kuboardspray_node_action")
			}
		}
	}
	if action == "remove_node" {
		nodesInvent := common.MapGet(inventoryNew, "all.hosts").(map[string]interface{})
		for key, nodeInfo := range nodesInvent {
			n := nodeInfo.(map[string]interface{})
			if n["kuboardspray_node_action"] == "remove_node" && !contains(nodesInK8s, key) {
				common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts."+key)
				common.MapDelete(inventoryNew, "all.children.target.children.k8s_cluster.children.kube_node.hosts."+key)
				common.MapDelete(inventoryNew, "all.children.target.children.etcd.hosts."+key)
				common.MapDelete(inventoryNew, "all.hosts."+key)
			}
		}
	}

	inventoryNewContent, err := yaml.Marshal(inventoryNew)
	if err != nil {
		return err
	}
	// logrus.Trace(string(inventoryNewContent))
	if err := ioutil.WriteFile(inventoryPath, inventoryNewContent, 0655); err != nil {
		logrus.Trace(err)
	}

	return nil
}
