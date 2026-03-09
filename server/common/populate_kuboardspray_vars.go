package common

import (
	"github.com/opencmit/pangee-cluster/constants"
)

func PopulatePangeeClusterVars(inventory map[string]interface{}, ownerType, ownerName string) {
	MapSet(inventory, "all.hosts.localhost.ansible_python_interpreter", "/usr/bin/python3")
	if MapGet(inventory, "all.hosts.all.bastion.bastion_type") == "socks5" {
		MapSet(inventory, "all.vars.ansible_ssh_args", "{{ pangeecluster_ssh_args }}")
	} else {
		// MapSet(inventory, "all.vars.ansible_ssh_args", "{{ pangeecluster_ssh_args }} {{ '' if ansible_password is defined and 'bastion' in groups['all'] else pangeecluster_ssh_multiplexing }}")
		MapSet(inventory, "all.vars.ansible_ssh_args", "{{ pangeecluster_ssh_args }}")
	}
	MapSet(inventory, "all.vars.ansible_ssh_pipelining", true)
	MapSet(inventory, "all.vars.pangeecluster_cluster_dir", constants.GET_DATA_DIR()+"/"+ownerType+"/"+ownerName)
	MapSet(inventory, "all.vars.pangeecluster_ssh_args", "-o ConnectionAttempts=3 -o UserKnownHostsFile=/dev/null -F /dev/null")
	MapSet(inventory, "all.vars.pangeecluster_ssh_multiplexing", "-o ControlMaster=no -o ControlPersist=60m -o ControlPath={{ pangeecluster_cluster_dir }}/{{ansible_user}}@{{ansible_host}}-{{ansible_port}}")
}
