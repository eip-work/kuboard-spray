package state

import (
	"errors"

	"github.com/eip-work/kuboard-spray/api/cluster"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/sirupsen/logrus"
)

func ExecuteShellOnETCD(clusterName string, shellCommand string) (*command.AnsibleStdout, error) {
	inventoryYamlPath := cluster.ClusterInventoryYamlPath(clusterName)

	inventory, err := common.ParseYamlFile(inventoryYamlPath)
	if err != nil {
		return nil, err
	}

	etcds := common.MapGet(inventory, "all.children.target.children.etcd.hosts").(map[string]interface{})

	errMsg := ""

	for key := range etcds {
		logrus.Trace("try on ", key, " : ", shellCommand)
		str := "export ETCDCTL_API=3 && export ETCDCTL_CERT=/etc/ssl/etcd/ssl/admin-" + key + ".pem && export ETCDCTL_KEY=/etc/ssl/etcd/ssl/admin-" + key + "-key.pem"
		str += " && export ETCDCTL_CACERT=/etc/ssl/etcd/ssl/ca.pem"
		str += " && " + shellCommand
		cmd := command.Run{
			Cmd: "ansible",
			Args: []string{
				key,
				"-m", "shell",
				"-a", str,
				"-i", inventoryYamlPath,
				"-e", "kuboardspray_cluster_dir=" + constants.GET_DATA_DIR() + "/cluster/" + clusterName,
			},
			Env: []string{
				"ANSIBLE_CONFIG=" + constants.GET_ADHOC_CFG_PATH(),
			},
			Timeout: 20,
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
			logrus.Warn("failed on " + key + ". " + result.Stdout)
			errMsg += "\nfailed on " + key + ". " + result.Stdout
		}
	}
	return nil, errors.New("cannot execute on any of the etcd \n" + errMsg)
}
