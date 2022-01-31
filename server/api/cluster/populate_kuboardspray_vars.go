package cluster

import (
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
)

func PopulateKuboardSprayVars(inventory map[string]interface{}, clusterName string) {
	common.MapSet(inventory, "all.vars.ansible_ssh_args", "{{ kuboardspray_ssh_args }} {{ '' if ansible_password is defined and 'bastion' in groups['all'] else kuboardspray_ssh_multiplexing }}")
	common.MapSet(inventory, "all.vars.ansible_ssh_pipelining", true)
	common.MapSet(inventory, "all.vars.kuboardspray_cluster_dir", constants.GET_DATA_CLUSTER_DIR()+"/"+clusterName)
	common.MapSet(inventory, "all.vars.kuboardspray_ssh_args", "-o ConnectionAttempts=100 -o UserKnownHostsFile=/dev/null -F /dev/null")
	common.MapSet(inventory, "all.vars.kuboardspray_ssh_multiplexing", "-o ControlMaster=auto -o ControlPersist=30m -o ControlPath={{ kuboardspray_cluster_dir }}/ansible-{{ansible_user}}@{{ansible_host}}:{{ansible_port}}")
}
