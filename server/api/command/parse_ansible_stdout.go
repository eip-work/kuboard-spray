package command

import (
	"encoding/json"
	"strings"

	"github.com/sirupsen/logrus"
)

type AnsibleStdout struct {
	ReturnCode string                 `json:"return_code"`
	Changed    string                 `json:"changed"`
	NodeName   string                 `json:"node_name"`
	Stdout     string                 `json:"stdout"`
	StdoutObj  map[string]interface{} `json:"stdout_obj"`
}

func ParseAnsibleStdout(stdout string) *AnsibleStdout {

	result := AnsibleStdout{}

	index := strings.Index(stdout, "\n")

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

	if content[0] == "{"[0] {
		temp := map[string]interface{}{}
		if err := json.Unmarshal([]byte(content), &temp); err != nil {
			logrus.Warn("ParseAnsibleStdout Error", err)
		} else {
			result.StdoutObj = temp
		}
	}

	return &result

}
