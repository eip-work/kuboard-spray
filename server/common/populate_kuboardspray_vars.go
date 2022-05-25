package common

import (
	"github.com/eip-work/kuboard-spray/constants"
)

func PopulateKuboardSprayVars(inventory map[string]interface{}, ownerType, ownerName string) {
	MapSet(inventory, "all.hosts.localhost.ansible_python_interpreter", "/usr/bin/python3")
	MapSet(inventory, "all.vars.ansible_ssh_args", "{{ kuboardspray_ssh_args }} {{ '' if ansible_password is defined and 'bastion' in groups['all'] else kuboardspray_ssh_multiplexing }}")
	MapSet(inventory, "all.vars.ansible_ssh_pipelining", true)
	MapSet(inventory, "all.vars.kuboardspray_cluster_dir", constants.GET_DATA_DIR()+"/"+ownerType+"/"+ownerName)
	MapSet(inventory, "all.vars.kuboardspray_ssh_args", "-o ConnectionAttempts=3 -o UserKnownHostsFile=/dev/null -F /dev/null")
	MapSet(inventory, "all.vars.kuboardspray_ssh_multiplexing", "-o ControlMaster=auto -o ControlPersist=60m -o ControlPath={{ kuboardspray_cluster_dir }}/{{ansible_user}}@{{ansible_host}}-{{ansible_port}}")
}
