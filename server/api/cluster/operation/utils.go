package operation

import (
	"strconv"
	"strings"

	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/sirupsen/logrus"
)

func contains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func getNodesInK8s(clusterName string) ([]string, error) {
	kubectlResult, err := cluster_common.ExecuteShellOnControlPlane(clusterName, "kubectl get nodes -o jsonpath=\"{.items[*].metadata.name}\"")

	if err != nil {
		return nil, err
	}

	nodeStr := strings.Trim(kubectlResult.Stdout, "\n")

	nodesInK8s := strings.Split(nodeStr, " ")
	logrus.Trace("Nodes in kubernetes: ", nodesInK8s)
	return nodesInK8s, nil
}

func arraySubtract(array1, array2 []string) []string {
	result := []string{}
	for _, v := range array1 {
		if !contains(array2, v) {
			result = append(result, v)
		}
	}
	return result
}

func getMembersInEtcd(clusterName string) ([]string, error) {
	result := []string{}
	temp, err := cluster_common.ExecuteShellOnETCD(clusterName, "etcdctl member list --write-out=json")
	if err != nil {
		return nil, err
	}
	for _, member := range temp.StdoutObj["members"].([]interface{}) {
		memberObj := member.(map[string]interface{})
		logrus.Trace("etcd member:", memberObj)
		if memberObj["name"] == nil {
			logrus.Warn("etcd member doesnot have attribute name : ", memberObj)
			break
		}
		nodeName := memberObj["name"].(string)
		result = append(result, nodeName)
	}

	return result, nil
}

type OperationCommonRequest struct {
	Cluster                     string            `uri:"cluster" binding:"required"`
	Fork                        int               `json:"fork"`
	Verbose                     string            `json:"verbose"`
	ExcludeNodes                string            `json:"nodes_to_exclude"`
	DiscoveredInterpreterPython map[string]string `json:"discovered_interpreter_python"`
	DownloadsOption             string            `json:"downloads_option"`
	Operation                   string
}

func appendCommonParams(result []string, req OperationCommonRequest, skipLimitParam bool) []string {
	result = append(result, "--fork", strconv.Itoa(req.Fork))
	result = append(result, "-e", "kuboardspray_operation="+req.Operation)
	if req.ExcludeNodes != "" && !skipLimitParam {
		result = append(result, "--limit", req.ExcludeNodes)
	}
	if req.Verbose == "vvv" {
		result = append(result, "-vvv")
	}
	if req.Verbose == "v" {
		result = append(result, "-v")
	}
	return result
}
