package command

import (
	"encoding/json"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func ExecuteShellCommands(owner_type, owner_name, target string, commands []AnsibleCommandsRequest) (*AnsibleResult, error) {
	return ExecuteShellCommandsWithStrategy(owner_type, owner_name, target, commands, "free")
}

func ExecuteShellCommandsWithStrategy(owner_type, owner_name, target string, commands []AnsibleCommandsRequest, strategy string) (*AnsibleResult, error) {

	startTime := time.Now()

	inventoryPath := constants.GetInventoryPath(owner_type, owner_name)

	tempPlayBookFilePath := "./temp_play_" + owner_type + "_" + owner_name + common.RandomString(10) + ".yaml"
	playbookFile, err := os.Create(tempPlayBookFilePath)
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
		"strategy":     strategy,
		"hosts":        target,
		"tasks":        tasks,
		"environment": gin.H{
			"ETCDCTL_API":       "3",
			"ETCDCTL_CERT":      "{{ kube_ssl_dir }}/etcd_server.crt",
			"ETCDCTL_KEY":       "{{ kube_ssl_dir }}/etcd_server.key",
			"ETCDCTL_CACERT":    "{{ kube_ssl_dir }}/ca.crt",
			"ETCDCTL_ENDPOINTS": "https://{{ hostvars[inventory_hostname]['ip'] }}:{{etcd_client_port}}",
		},
	}

	playbookStr, err := yaml.Marshal([]gin.H{playbook})
	if err != nil {
		return nil, err
	}

	logrus.Trace(string(playbookStr))
	_, err = playbookFile.Write(playbookStr)
	if err != nil {
		return nil, err
	}

	cmd := Run{
		Cmd: "ansible-playbook",
		Args: []string{
			tempPlayBookFilePath,
			"-i", inventoryPath,
			"-f", "200",
			"-e", "pangeecluster_ssh_args='-o ConnectionAttempts=1 -o UserKnownHostsFile=/dev/null -F /dev/null'",
		},
		Env:     []string{"ANSIBLE_CONFIG=./ansible-rpc/ansible.cfg", "ANSIBLE_CACHE_PLUGIN_CONNECTION=" + constants.GET_DATA_DIR() + "/" + owner_type + "/" + owner_name + "/fact"},
		Timeout: 20,
	}

	stdout, stderr, err := cmd.Run()
	if err != nil {
		duration := time.Now().UnixNano() - startTime.UnixNano()
		logrus.Trace("stdout:", string(stdout))
		logrus.Trace("stderr:", string(stderr))
		logrus.Trace("err: ", err)
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
