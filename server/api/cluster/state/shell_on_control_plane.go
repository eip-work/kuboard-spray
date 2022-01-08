package state

import (
	"errors"
	"strconv"

	"github.com/eip-work/kuboard-spray/api/cluster"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/sirupsen/logrus"
)

func ExecuteShellOnControlPlane(clusterName string, shellCommand string) (*command.AnsibleStdout, error) {
	inventoryYamlPath := cluster.ClusterInventoryYamlPath(clusterName)

	inventory, err := common.ParseYamlFile(inventoryYamlPath)
	if err != nil {
		return nil, err
	}

	control_planes := common.MapGet(inventory, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts").(map[string]interface{})

	index := 0
	errMsg := ""
	for key := range control_planes {
		logrus.Trace("try on ", key, " : ", shellCommand)
		cmd := command.Run{
			Cmd: "ansible",
			Args: []string{
				"kube_control_plane[" + strconv.Itoa(index) + "]",
				"-m", "shell",
				"-a", shellCommand,
				"-i", inventoryYamlPath,
			},
			Env: []string{"ANSIBLE_CONFIG=" + constants.GET_ADHOC_CFG_PATH()},
			// Dir: resourcePackageDir,
		}
		stdout, _, err := cmd.Run()
		if err != nil {
			return nil, err
		}
		// logrus.Trace(string(stdout), string(stderr))

		result := command.ParseAnsibleStdout(string(stdout))

		if result.Changed == "CHANGED" {
			return result, nil
		} else {
			errMsg += "\nfailed on kube_control_plane[" + strconv.Itoa(index) + "] " + result.Stdout
			index++
		}
	}
	return nil, errors.New("cannot reach any of the kube_control_plane \n" + errMsg)
}
