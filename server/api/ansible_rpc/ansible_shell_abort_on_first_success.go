package ansible_rpc

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

func ExecuteShellCommandsAbortOnFirstSuccess(owner_type, owner_name, target string, commands []AnsibleCommandsRequest) ([]AnsibleResultNode, error) {
	shellResult, err := ExecuteShellCommands(owner_type, owner_name, target, commands)
	if err != nil {
		return nil, err
	}

	result := []AnsibleResultNode{}
	for _, t := range shellResult.Plays[0].Tasks {
		for _, r := range t.Hosts {
			if r.StdOut != "" && (r.StdOut[0:1] == "{" || r.StdOut[0:1] == "[") {
				if err := json.Unmarshal([]byte(r.StdOut), &r.StdOutObj); err != nil {
					logrus.Warning("failed to parse", err)
				}
			}
			result = append(result, r)
			break
		}
	}
	return result, nil
}
