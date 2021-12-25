package cluster

import (
	"runtime/debug"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/sirupsen/logrus"
)

func ClusterResourcePackageDir(cluster string) string {
	inventoryYamlPath := ClusterInventoryYamlPath(cluster)
	inventory, err := common.ParseYamlFile(inventoryYamlPath)
	if err != nil {
		logrus.Warn("failed to parse inventory yaml", err)
		debug.PrintStack()
		return ""
	}
	resourcePackagePath := constants.GET_DATA_RESOURCE_DIR() + "/" + common.MapGet(inventory, "all.hosts.localhost.kuboardspray_resource_package").(string)
	return resourcePackagePath
}
