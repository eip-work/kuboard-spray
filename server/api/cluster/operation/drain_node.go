package operation

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

type DrainNodeRequest struct {
	OperationCommonRequest
	Nodes                  string `json:"nodes"`
	DrainTimeout           string `json:"drain_timeout"`
	DrainOutRetries        string `json:"drain_retries"`
	DrainRetryDelaySeconds string `json:"drain_retry_delay_seconds"`
	DrainGracePeriod       string `json:"drain_grace_period"`
}

func DrainNode(c *gin.Context) {

	var req DrainNodeRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)
	req.Operation = "drain_node"

	inventory, resourcePackage, err := updateResourcePackageVarsToInventory(req.OperationCommonRequest)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to process inventory", err)
		return
	}

	postExec := func(status command.ExecuteExitStatus) (string, error) {

		var message string
		if status.Success {
			message += "\n"
			message = "\033[32m[ " + "Finished draining node: " + req.Nodes + ". ]\033[0m \n"
			message += "\033[32m[ 成功排空节点: " + req.Nodes + " ]\033[0m \n"
		} else {
			message += "\n"
			message = "\033[31m\033[01m\033[05m[ Failed to drain node: " + req.Nodes + ". ]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[ 排空节点失败: " + req.Nodes + ". ]\033[0m \n"
		}

		return "\n" + message, nil
	}

	playbook := common.MapGet(resourcePackage, "data.supported_playbooks.drain_node").(string)

	cmd := command.Execute{
		OwnerType: "cluster",
		OwnerName: req.Cluster,
		Cmd:       "ansible-playbook",
		Args: func(execute_dir string) []string {
			result := []string{
				playbook,
				"-i", execute_dir + "/inventory.yaml",
				"-e", "upgrade_node_always_cordon=true",
				"-e", "drain_grace_period=" + req.DrainGracePeriod,
				"-e", "drain_timeout=" + req.DrainTimeout,
				"-e", "drain_retries=" + req.DrainOutRetries,
				"-e", "drain_retry_delay_seconds=" + req.DrainRetryDelaySeconds,
			}
			result = appendCommonParams(result, req.OperationCommonRequest, false)
			result = append(result, "--limit", req.Nodes)
			return result
		},
		Dir:      cluster_common.ResourcePackageDirForInventory(inventory),
		Type:     req.Operation,
		PreExec:  func(execute_dir string) error { return common.SaveYamlFile(execute_dir+"/inventory.yaml", inventory) },
		PostExec: postExec,
	}

	if err := cmd.Exec(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Faild to drain node. ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"pid": cmd.R_Pid,
		},
	})
}
