package ansible_rpc

import "time"

type AnsibleResult struct {
	Plays []AnsibleResultPlay `json:"plays"`
}

type AnsibleResultPlay struct {
	Tasks []AnsibleResultTask `json:"tasks"`
}

type AnsibleResultTask struct {
	Hosts map[string]AnsibleResultNode `json:"hosts"`
	Task  AnsibleResultTaskMeta        `json:"task"`
}

type AnsibleResultNode struct {
	Action     string `json:"action"`
	Changed    bool   `json:"changed"`
	Cmd        string `json:"cmd"`
	StdErr     string `json:"stderr"`
	StdOut     string `json:"stdout"`
	ReturnCode int    `json:"rc"`
	Skipped    bool   `json:"skipped"`
	SkipReason string `json:"skip_reason"`
	Delta      string `json:"delta"`
}

type AnsibleResultTaskMeta struct {
	Id       string                `json:"id"`
	Name     string                `json:"name"`
	Duration AnsibleResultDuration `json:"duration"`
}

type AnsibleResultDuration struct {
	End   time.Time `json:"end"`
	Start time.Time `json:"start"`
}

type AnsibleCommandsRequest struct {
	Name    string
	Command string
}
