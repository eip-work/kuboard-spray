package operation_v2

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/opencmit/pangee-cluster/api/cluster/cluster_common"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
)

func updateResourcePackageVarsToInventory(req ExecuteStepRequest) (map[string]interface{}, map[string]interface{}, error) {

	clusterMetadata, err := cluster_common.ClusterMetadataByName(req.Cluster)
	if err != nil {
		return nil, nil, err
	}
	inventory := clusterMetadata.Inventory
	resourcePackage := clusterMetadata.ResourcePackage

	common.MapSet(inventory, "all.vars.pangeecluster_no_log", !(req.Verbose == "v" || req.Verbose == "vvv"))

	versionJsonPath := constants.GetPangeeClusterWebDir() + "/version.json"
	versionJson, err := os.ReadFile(versionJsonPath)
	if err != nil {
		return nil, nil, err
	}
	version := map[string]string{}
	err = json.Unmarshal(versionJson, &version)
	if err != nil {
		return nil, nil, err
	}

	v := version["version"]
	v = strings.TrimSuffix(v, "-amd64")
	v = strings.TrimSuffix(v, "-arm64")

	common.MapSet(inventory, "all.vars.pangeecluster_version", v)

	common.PopulatePangeeClusterVars(inventory, "cluster", req.Cluster)

	common.MapSet(inventory, "all.vars.ansible_ssh_pipelining", false)

	return inventory, resourcePackage, nil
}
