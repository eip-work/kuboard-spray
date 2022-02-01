package cluster_common

import (
	"github.com/eip-work/kuboard-spray/common"
)

type ClusterMetadata struct {
	ResourcePackageDir string
	ResourcePackage    map[string]interface{}
	InventoryPath      string
	Inventory          map[string]interface{}
	ClusterName        string
}

func ClusterMetadataByName(name string) (*ClusterMetadata, error) {
	inventoryPath := ClusterInventoryYamlPath(name)
	inventory, err := common.ParseYamlFile(inventoryPath)
	if err != nil {
		return nil, err
	}
	resourcePackageDir := ResourcePackageDirForInventory(inventory)
	resourcePackage, err := common.ParseYamlFile(resourcePackageDir + "/package.yaml")
	if err != nil {
		return nil, err
	}

	result := ClusterMetadata{
		ResourcePackageDir: resourcePackageDir,
		ResourcePackage:    resourcePackage,
		InventoryPath:      inventoryPath,
		Inventory:          inventory,
		ClusterName:        name,
	}
	return &result, nil
}
