package ssh

import (
	"errors"
	"strconv"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type ptyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

type TerminalSpec struct {
	Columns uint32 `json:"cols"`
	Rows    uint32 `json:"rows"`
}

type NodeInfo struct {
	Host           string `json:"ansible_host"`
	User           string `json:"ansible_user"`
	Password       string `json:"ansible_password"`
	PrivateKeyPath string `json:"ansible_ssh_private_key_file"`
	Become         bool   `json:"ansible_become"`
	BecomeUser     string `json:"ansible_become_user"`
	BecomePassword string `json:"ansible_become_password"`
	Port           int    `json:"ansible_port"`
}

type SSHClient struct {
	NodeInfo
	Bastion *NodeInfo `json:"bastion"`
	Session *ssh.Session
	Client  *ssh.Client
	channel ssh.Channel
}

func NewSSHClient(shellRequest ShellRequest) (*SSHClient, error) {
	client := SSHClient{}

	inventoryPath := constants.GET_DATA_DIR() + "/" + shellRequest.NodeOwnerType + "/" + shellRequest.NodeOwner + "/inventory.yaml"
	inventory, err := common.ParseYamlFile(inventoryPath)
	if err != nil {
		return nil, err
	}
	nodeRef := common.MapGet(inventory, "all.hosts."+shellRequest.Node)
	if nodeRef == nil {
		return nil, errors.New("cannot find node " + shellRequest.Node + " in inventory " + inventoryPath)
	}
	client.NodeInfo = *populateNodeInfo(inventory, shellRequest.Node)

	bastionRef := common.MapGet(inventory, "all.hosts.bastion")
	if bastionRef != nil {
		client.Bastion = populateNodeInfo(inventory, "bastion")
	}

	return &client, nil
}

func populateNodeInfo(inventory map[string]interface{}, node string) *NodeInfo {
	result := &NodeInfo{}
	result.Host = getNodeInfo(inventory, node, "ansible_host")
	result.User = getNodeInfo(inventory, node, "ansible_user")
	result.Password = getNodeInfo(inventory, node, "ansible_password")
	port := getNodeInfo(inventory, node, "ansible_port")
	var err error
	result.Port, err = strconv.Atoi(port)
	if err != nil {
		logrus.Warn("ansible_port cannot be parsed into int")
		result.Port = 22
	}
	result.PrivateKeyPath = getNodeInfo(inventory, node, "ansible_ssh_private_key_file")
	result.Become = getNodeInfoBool(inventory, node, "ansible_become")
	result.BecomeUser = getNodeInfo(inventory, node, "ansible_become_user")
	result.BecomePassword = getNodeInfo(inventory, node, "ansible_become_password")
	return result
}

func getNodeInfo(inventory map[string]interface{}, node, prop string) string {
	v := common.MapGetString(inventory, "all.hosts."+node+"."+prop)
	if v == "" {
		v = common.MapGetString(inventory, "all.children.target.vars."+prop)
	}
	return v
}
func getNodeInfoBool(inventory map[string]interface{}, node, prop string) bool {
	v := common.MapGet(inventory, "all.hosts."+node+"."+prop)
	if v == nil {
		v = common.MapGet(inventory, "all.children.target.vars."+prop)
	}
	if v != nil {
		return v.(bool)
	}
	return false
}
