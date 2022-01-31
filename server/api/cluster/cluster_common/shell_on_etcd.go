package cluster_common

import (
	"errors"
	"time"

	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/sirupsen/logrus"
)

func ExecuteShellOnETCD(clusterName string, shellCommand string) (*command.AnsibleStdout, error) {
	startTime := time.Now()

	inventoryYamlPath := ClusterInventoryYamlPath(clusterName)

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
			},
			Env: []string{
				"ANSIBLE_CONFIG=" + constants.GET_ADHOC_CFG_PATH(),
			},
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
	return nil, errors.New("cannot execute on any of the etcd \n" + errMsg)
}
