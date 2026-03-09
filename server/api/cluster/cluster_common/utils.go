package cluster_common

import (
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
)

func ResourcePackageDirForInventory(inventory map[string]interface{}) string {
	return constants.GET_DATA_RESOURCE_DIR() + "/" + common.MapGet(inventory, "all.hosts.localhost.pangeecluster_resource_package").(string) + "/content"
}

func ClusterInventoryYamlPath(cluster string) string {
	return constants.GET_DATA_CLUSTER_DIR() + "/" + cluster + "/inventory.yaml"
}
