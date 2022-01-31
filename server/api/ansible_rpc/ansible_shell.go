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
	playbookFile, err := os.OpenFile(tempPlayBookFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
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
		"hosts":        target,
		"tasks":        tasks,
	}

	playbookStr, err := yaml.Marshal([]gin.H{playbook})
	if err != nil {
		return nil, err
	}

	playbookFile.Write(playbookStr)

	cmd := command.Run{
		Cmd: "ansible-playbook",
		Args: []string{
			tempPlayBookFilePath,
			"-i", inventoryPath,
			"-e", "kuboardspray_cluster_dir=" + constants.GET_DATA_DIR() + "/" + owner_type + "/" + owner_name,
		},
		Env:     []string{"ANSIBLE_CONFIG=./ansible-rpc/ansible.cfg"},
		Timeout: 20,
		// Dir: resourcePackageDir,
	}

	stdout, _, err := cmd.Run()
	if err != nil {
		duration := time.Now().UnixNano() - startTime.UnixNano()
		logrus.Trace("duration: ", duration/1000000)
		return nil, err
	}

	result := &AnsibleResult{}
	if err := json.Unmarshal(stdout, result); err != nil {
		logrus.Trace("ERROR: ", err)
	}

	return result, nil
}
