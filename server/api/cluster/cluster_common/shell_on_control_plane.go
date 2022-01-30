package cluster_common

import (
	"errors"
	"time"

	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/sirupsen/logrus"
)

func ExecuteShellOnControlPlane(clusterName string, shellCommand string) (*command.AnsibleStdout, error) {
	startTime := time.Now()

	inventoryYamlPath := ClusterInventoryYamlPath(clusterName)

	inventory, err := common.ParseYamlFile(inventoryYamlPath)
	if err != nil {
		return nil, err
	}

	control_planes := common.MapGet(inventory, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts").(map[string]interface{})

	errMsg := ""
	for key := range control_planes {
		logrus.Trace("try on ", key, " : ", shellCommand)
		cmd := command.Run{
			Cmd: "ansible",
			Args: []string{
				key,
				"-m", "shell",
				"-a", shellCommand,
				"-i", inventoryYamlPath,
				"-e", "kuboardspray_cluster_dir=" + constants.GET_DATA_DIR() + "/cluster/" + clusterName,
			},
			Env:     []string{"ANSIBLE_CONFIG=" + constants.GET_ADHOC_CFG_PATH()},
			Timeout: 20,
			// Dir: resourcePackageDir,
		}
		stdout, _, err := cmd.Run()
		if err != nil {
			duration := time.Now().UnixNano() - startTime.UnixNano()
			logrus.Trace("duration: ", duration/1000000)
			return nil, err
		}
		// logrus.Trace(string(stdout), string(stderr))

		result := command.ParseAnsibleStdout(string(stdout))

		if result.Changed == "CHANGED" {
			duration := time.Now().UnixNano() - startTime.UnixNano()
			logrus.Trace("duration: ", duration/1000000)
			return result, nil
		} else {
			logrus.Warn("failed on " + key + ". " + result.Stdout)
			errMsg += "\nfailed on " + key + ". " + result.Stdout
		}
	}
	duration := time.Now().UnixNano() - startTime.UnixNano()
	logrus.Trace("duration: ", duration/1000000)
	return nil, errors.New("cannot execute on any of the kube_control_plane \n" + errMsg)
}
