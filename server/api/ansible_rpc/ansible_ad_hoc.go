package ansible_rpc

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/sirupsen/logrus"
)

type AdhocCommandRequestWithIP struct {
	NodeOwnerType     string `uri:"node_owner_type"`
	NodeOwner         string `uri:"node_owner"`
	Ip                string `json:"ansible_host" binding:"required"`
	Port              string `json:"ansible_port" binding:"required"`
	User              string `json:"ansible_user" binding:"required"`
	Password          string `json:"ansible_password"`
	PrivateKeyFile    string `json:"ansible_ssh_private_key_file"`
	Become            bool   `json:"ansible_become"`
	BecomeUser        string `json:"ansible_become_user"`
	BecomePassword    string `json:"ansible_become_password"`
	SshCommonArgs     string `json:"ansible_ssh_common_args"`
	PythonInterpreter string `json:"ansible_python_interpreter"`
}

func ExecuteAdhocCommandWithIp(req AdhocCommandRequestWithIP, args []string) ([]AnsibleResultNode, error) {

	inventory := map[string]interface{}{
		"all": map[string]interface{}{
			"hosts": map[string]interface{}{
				"temp": map[string]interface{}{
					"ansible_host":                 req.Ip,
					"ansible_port":                 req.Port,
					"ansible_user":                 req.User,
					"ansible_password":             req.Password,
					"ansible_ssh_private_key_file": req.PrivateKeyFile,
					"ansible_become":               req.Become,
					"ansible_become_user":          req.BecomeUser,
					"ansible_become_password":      req.BecomePassword,
					"ansible_ssh_common_args":      "-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectionAttempts=1 -F /dev/null" + req.SshCommonArgs,
					"ansible_ssh_pipelining":       true,
					"kuboardspray_cluster_dir":     constants.GET_DATA_DIR() + "/" + req.NodeOwnerType + "/" + req.NodeOwner,
					"ansible_python_interpreter":   "auto",
				},
			},
		},
	}
	if req.PythonInterpreter != "" {
		common.MapSet(inventory, "all.hosts.temp.ansible_python_interpreter", req.PythonInterpreter)
	}

	inventoryBytes, err := json.Marshal(inventory)
	if err != nil {
		return nil, errors.New("failed to marshal inventory file: " + err.Error())
	}

	inventoryPath := "adhoc_" + time.Now().Format("2006-01-02_15-04-05.999") + ".json"

	err = ioutil.WriteFile("./ansible-rpc/"+inventoryPath, inventoryBytes, 0666)
	if err != nil {
		return nil, errors.New("failed to create inventory file " + inventoryPath + err.Error())
	}
	logrus.Trace(string(inventoryBytes))

	defer os.Remove("./ansible-rpc/" + inventoryPath)

	arguments := []string{"temp", "-i", inventoryPath, "--timeout", "3", "--forks", "200"}
	arguments = append(arguments, args...)

	run := command.Run{
		Cmd:  "ansible",
		Args: arguments,
		Env:  []string{"ANSIBLE_CONFIG=./ansible.cfg"},
		Dir:  "./ansible-rpc",
	}

	stdout, stderr, err := run.Run()

	if err != nil {
		logrus.Trace(err)
		return nil, errors.New("failed to execute adhoc command " + err.Error())
	}

	var result AnsibleResult
	if err := json.Unmarshal(stdout, &result); err != nil {
		logrus.Trace("stdout: ", string(stdout))
		logrus.Trace("stderr: ", string(stderr))
		return nil, errors.New("Failed to Unmarshal result " + err.Error())
	}
	var nodeResults []AnsibleResultNode
	for nodeName, node := range result.Plays[0].Tasks[0].Hosts {
		temp := node
		temp.InventoryHostName = nodeName
		nodeResults = append(nodeResults, temp)
	}

	return nodeResults, nil
}

func ExecuteAdhocCommandWithInventory(inventoryPath string, args []string) ([]AnsibleResultNode, error) {

	arguments := []string{"-i", inventoryPath, "--timeout", "3", "--forks", "200", "-e", "kuboardspray_ssh_args='-o ConnectionAttempts=1 -o UserKnownHostsFile=/dev/null -F /dev/null'"}
	arguments = append(arguments, args...)

	run := command.Run{
		Cmd:  "ansible",
		Args: arguments,
		Env:  []string{"ANSIBLE_CONFIG=./ansible.cfg"},
		Dir:  "./ansible-rpc",
	}

	stdout, stderr, err := run.Run()

	if err != nil {
		logrus.Trace(err)
		return nil, errors.New("failed to execute adhoc command " + err.Error())
	}

	// logrus.Trace(string(stdout))

	var result AnsibleResult
	if err := json.Unmarshal(stdout, &result); err != nil {
		logrus.Trace("stdout: ", string(stdout))
		logrus.Trace("stderr: ", string(stderr))
		return nil, errors.New("Failed to Unmarshal result " + err.Error())
	}

	var nodeResults []AnsibleResultNode
	for nodeName, node := range result.Plays[0].Tasks[0].Hosts {
		temp := node
		temp.InventoryHostName = nodeName
		nodeResults = append(nodeResults, temp)
	}

	return nodeResults, nil
}
