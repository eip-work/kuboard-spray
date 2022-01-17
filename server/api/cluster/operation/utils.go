package operation

import (
	"strings"

	"github.com/eip-work/kuboard-spray/api/cluster/state"
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
	kubectlResult, err := state.ExecuteShellOnControlPlane(clusterName, "kubectl get nodes -o jsonpath=\"{.items[*].metadata.name}\"")

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

// func getMembersInEtcd(clusterName string) ([]string, error) {
// 	result, err := state.ExecuteShellOnETCD(clusterName, "etcdctl member list --write-out=json")
// 	if err != nil {
// 		return nil, err
// 	}

// }
