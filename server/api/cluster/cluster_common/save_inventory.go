package cluster_common

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

func SaveInventory(clusterName string, inventory map[string]interface{}) error {
	content, _ := yaml.Marshal(inventory)

	if err := ioutil.WriteFile(ClusterInventoryYamlPath(clusterName), content, 0655); err != nil {
		return err
	}
	return nil
}
