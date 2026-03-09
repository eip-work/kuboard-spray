package state

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/api/cluster/cluster_common"
	"github.com/opencmit/pangee-cluster/api/command"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
	"github.com/sirupsen/logrus"
)

type ClusterVersion map[string](map[string]command.AnsibleResultNode)

func CheckClusterVersion(c *gin.Context) {

	var request GetNodesRequest
	c.ShouldBindUri(&request)

	startTime := time.Now()

	cluster, err := cluster_common.ClusterMetadataByName(request.ClusterName)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to get cluster metadata", err)
		return
	}

	playbook := common.MapGetString(cluster.ResourcePackage, "data.supported_playbooks.cluster_version")
	if playbook == "" {
		common.HandleError(c, http.StatusInternalServerError, "current resource package doesnot support the operation to view cluster version", nil)
		return
	}

	cmd := command.Run{
		Cmd: "ansible-playbook",
		Args: []string{
			cluster.ResourcePackageDir + "/" + playbook,
			"-i", cluster.InventoryPath,
			"--forks", "200",
			"--timeout", "3",
			"-e", "pangeecluster_ssh_args='-o ConnectionAttempts=1 -o UserKnownHostsFile=/dev/null -F /dev/null'",
		},
		Env:     []string{"ANSIBLE_CONFIG=./ansible.cfg", "ANSIBLE_CACHE_PLUGIN_CONNECTION=" + constants.GET_DATA_CLUSTER_DIR() + "/" + request.ClusterName + "/fact"},
		Timeout: 60,
		Dir:     "./ansible-rpc",
	}

	stdout, stderr, err := cmd.Run()
	duration := time.Now().UnixNano() - startTime.UnixNano()
	logrus.Trace("duration: ", duration/1000000)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to run", err)
		return
	}

	// logrus.Trace("stdout: ", string(stdout), "\nstderr: ", string(stderr))

	result := &command.AnsibleResult{}
	if err := json.Unmarshal(stdout, result); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to Unmarshal result: ["+string(stdout)+"]", err)
		logrus.Trace("stdout: ", string(stdout), "\nstderr: ", string(stderr))
		logrus.Trace("duration: ", duration/1000000)
		return
	}

	clusterVersion := ClusterVersion{}

	for _, task := range result.Plays[0].Tasks {
		for nodeName, node := range task.Hosts {
			if clusterVersion[nodeName] == nil {
				clusterVersion[nodeName] = make(map[string]command.AnsibleResultNode)
			}
			clusterVersion[nodeName][task.Task.Name] = node
		}
	}

	c.JSON(http.StatusOK, common.PangeeClusterResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    clusterVersion,
	})
}
