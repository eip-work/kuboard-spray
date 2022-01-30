package cluster_common

import (
	"encoding/json"
	"strings"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/sirupsen/logrus"
)

func ResourcePackagePathForInventory(inventory map[string]interface{}) string {
	return constants.GET_DATA_RESOURCE_DIR() + "/" + common.MapGet(inventory, "all.hosts.localhost.kuboardspray_resource_package").(string) + "/content"
}

func ClusterInventoryYamlPath(cluster string) string {
	return constants.GET_DATA_CLUSTER_DIR() + "/" + cluster + "/inventory.yaml"
}

const multi_command_spliter = "----------kuboardspray_multi_command_spliter--------\n"

type AnsibleMultiCommandsStdout struct {
	ReturnCode string                   `json:"return_code"`
	Changed    string                   `json:"changed"`
	NodeName   string                   `json:"node_name"`
	Stdout     string                   `json:"stdout"`
	Stdouts    []string                 `json:"stdouts"`
	StdoutObjs []map[string]interface{} `json:"stdout_objs"`
}

func ParseAnsibleMultiCommandsStdout(stdout string) *AnsibleMultiCommandsStdout {

	result := AnsibleMultiCommandsStdout{}

	index := strings.Index(stdout, "\n")

	if index <= 0 {
		result.ReturnCode = ""
		result.NodeName = ""
		result.Stdouts = []string{}
		result.StdoutObjs = []map[string]interface{}{}
		return &result
	}

	line0 := stdout[:index]
	content := stdout[index+1:]

	if line0[len(line0)-1] == "{"[0] {
		content = "{\n" + content
	}

	line0Splited := strings.Split(line0, "|")
	result.NodeName = strings.Trim(line0Splited[0], " ")
	result.Changed = strings.Trim(line0Splited[1], " ")
	if len(line0Splited) > 2 {
		result.ReturnCode = strings.Trim(line0Splited[2], " ")
	}

	result.Stdout = content
	temp := strings.Split(content, multi_command_spliter)

	for _, stdout := range temp {
		if len(stdout) == 0 {
			continue
		}
		result.Stdouts = append(result.Stdouts, stdout)
		if stdout[0] == "{"[0] {
			temp := map[string]interface{}{}
			if err := json.Unmarshal([]byte(stdout), &temp); err != nil {
				logrus.Warn("ParseAnsibleStdout Error", err)
			} else {
				result.StdoutObjs = append(result.StdoutObjs, temp)
			}
		}
	}

	return &result

}
