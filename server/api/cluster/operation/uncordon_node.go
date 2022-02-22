package operation

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

type UncordonNodeRequest struct {
	OperationCommonRequest
	Nodes string `json:"nodes"`
}

func UncordonNode(c *gin.Context) {

	var req UncordonNodeRequest
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
			message = "\033[32m[ " + "Finished uncordoning node: " + req.Nodes + ". ]\033[0m \n"
			message += "\033[32m[ 成功恢复节点调度: " + req.Nodes + " ]\033[0m \n"
		} else {
			message += "\n"
			message = "\033[31m\033[01m\033[05m[ Failed to uncordon node: " + req.Nodes + ". ]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[ 恢复节点调度失败: " + req.Nodes + ". ]\033[0m \n"
		}

		return "\n" + message, nil
	}

	playbook := common.MapGet(resourcePackage, "data.supported_playbooks.uncordon_node").(string)

	cmd := command.Execute{
		OwnerType: "cluster",
		OwnerName: req.Cluster,
		Cmd:       "ansible-playbook",
		Args: func(execute_dir string) []string {
			result := []string{
				playbook,
				"-i", execute_dir + "/inventory.yaml",
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
