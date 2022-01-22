package cluster_common

import (
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
)

func ResourcePackagePathForInventory(inventory map[string]interface{}) string {
	return constants.GET_DATA_RESOURCE_DIR() + "/" + common.MapGet(inventory, "all.hosts.localhost.kuboardspray_resource_package").(string) + "/content"
}

func ClusterInventoryYamlPath(cluster string) string {
	return constants.GET_DATA_CLUSTER_DIR() + "/" + cluster + "/inventory.yaml"
}
