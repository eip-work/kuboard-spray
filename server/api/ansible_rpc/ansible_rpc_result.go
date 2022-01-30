package ansible_rpc

type AnsibleResult struct {
	Plays []AnsibleResultPlay `json:"plays"`
}

type AnsibleResultPlay struct {
	Tasks []AnsibleResultTask `json:"tasks"`
}

type AnsibleResultTask struct {
	Hosts map[string]AnsibleResultNode `json:"hosts"`
}

type AnsibleResultNode struct {
	Action     string `json:"action"`
	Changed    bool   `json:"changed"`
	Cmd        string `json:"cmd"`
	StdErr     string `json:"stderr"`
	StdOut     string `json:"stdout"`
	ReturnCode int    `json:"rc"`
}
