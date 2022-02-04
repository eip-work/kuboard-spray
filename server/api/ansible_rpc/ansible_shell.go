package ansible_rpc

import (
	"encoding/json"
	"os"
	"time"

	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func ExecuteShellCommands(owner_type, owner_name, target string, commands []AnsibleCommandsRequest) (*AnsibleResult, error) {

	startTime := time.Now()

	inventoryPath := constants.GetInventoryPath(owner_type, owner_name)

	tempPlayBookFilePath := "./temp_play_" + owner_type + "_" + owner_name + common.RandomString(10) + ".yaml"
	playbookFile, err := os.OpenFile(tempPlayBookFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0655)
	if err != nil {
		return nil, err
	}

	defer playbookFile.Close()
	defer os.Remove(tempPlayBookFilePath)

	tasks := []gin.H{}
	for _, cmd := range commands {
		task := gin.H{
			"name":  cmd.Name,
			"shell": cmd.Command,
		}
		tasks = append(tasks, task)
	}

	playbook := gin.H{
		"name":         "temp_play",
		"gather_facts": false,
		"strategy":     "free",
		"hosts":        target,
		"tasks":        tasks,
		"environment": gin.H{
			"ETCDCTL_API":    "3",
			"ETCDCTL_CERT":   "/etc/ssl/etcd/ssl/admin-{{inventory_hostname}}.pem",
			"ETCDCTL_KEY":    "/etc/ssl/etcd/ssl/admin-{{inventory_hostname}}-key.pem",
			"ETCDCTL_CACERT": "/etc/ssl/etcd/ssl/ca.pem",
		},
	}

	playbookStr, err := yaml.Marshal([]gin.H{playbook})
	if err != nil {
		return nil, err
	}

	logrus.Trace(string(playbookStr))
	playbookFile.Write(playbookStr)

	cmd := command.Run{
		Cmd: "ansible-playbook",
		Args: []string{
			tempPlayBookFilePath,
			"-i", inventoryPath,
			"-f", "200",
			"-e", "kuboardspray_ssh_args='-o ConnectionAttempts=1 -o UserKnownHostsFile=/dev/null -F /dev/null'",
		},
		Env:     []string{"ANSIBLE_CONFIG=./ansible-rpc/ansible.cfg"},
		Timeout: 20,
	}

	stdout, stderr, err := cmd.Run()
	if err != nil {
		duration := time.Now().UnixNano() - startTime.UnixNano()
		logrus.Trace("duration: ", duration/1000000)
		return nil, err
	}

	result := &AnsibleResult{}
	if err := json.Unmarshal(stdout, result); err != nil {
		logrus.Trace("Failed to pase stdout, ERROR:", err)
		logrus.Trace("stdout:", string(stdout))
		logrus.Trace("stderr:", string(stderr))
		return nil, err
	}

	return result, nil
}
