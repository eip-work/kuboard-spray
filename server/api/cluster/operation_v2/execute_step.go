package operation_v2

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/api/cluster/cluster_common"
	"github.com/opencmit/pangee-cluster/api/command"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/sirupsen/logrus"
)

type ExecuteStepRequest struct {
	Cluster      string `uri:"cluster" binding:"required"`
	Operation    string `uri:"operation" binding:"required"`
	Step         string `uri:"step" binding:"required"`
	Fork         int    `json:"fork"`
	Verbose      string `json:"verbose"`
	ExcludeNodes string `json:"nodes_to_exclude"`
}

// 对于 inventory 中以 path 指定的 hosts 内所有属性 pangeecluster_resource_package 为 originalAction 的 host，
// 重置其 pangeecluster_resource_package 字段值为 none
func resetNodeAction(inventory map[string]interface{}, path string, originalAction string) {
	hosts := common.MapGet(inventory, path).(map[string]interface{})
	if originalAction == "delete_node" {
		for key := range hosts {
			hostName := "all.hosts." + key
			if common.MapGetString(inventory, hostName+".pangeecluster_node_action") == "delete_node" {
				common.MapDelete(hosts, key)
				common.MapDelete(inventory, "all.hosts."+key)
			}
		}
	} else {
		for key := range hosts {
			hostActionPath := "all.hosts." + key + ".pangeecluster_node_action"
			if common.MapGet(inventory, hostActionPath) == originalAction {
				common.MapSet(inventory, hostActionPath, "none")
			}
		}
	}
}

func ExecuteStep(c *gin.Context) {
	var req ExecuteStepRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)

	inventory, _, err := updateResourcePackageVarsToInventory(req)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to process inventory", err)
		return
	}

	postExec := func(status command.ExecuteExitStatus) (string, error) {
		if status.Success {
			// FIXME 待验证

			tempClusterMetadata, err := cluster_common.ClusterMetadataByName(req.Cluster)
			if err != nil {
				return "解析数据出错", err
			}

			// logrus.Info("将要获取 data.operations")
			operations := common.MapGet(tempClusterMetadata.ResourcePackage, "data.operations").([]interface{})
			for index, operationItf := range operations { // 遍历所有的 operation
				logrus.Info("遍历 operation", index)
				operation := operationItf.(map[string]interface{})
				if common.MapGetString(operation, "name") == req.Operation { // 匹配到当前 operation
					logrus.Info("匹配到 operation: ", req.Operation)
					steps := common.MapGet(operation, "steps").([]interface{})
					for _, stepItf := range steps { // 遍历所有的 steps
						step := stepItf.(map[string]interface{})
						if common.MapGetString(step, "name") == req.Step { // 匹配到当前 step
							logrus.Info("匹配到 step: ", req.Step)
							if common.MapGet(step, "mark_pangeecluster_node_action_done") != nil {
								markDone := common.MapGet(step, "mark_pangeecluster_node_action_done").(map[string]interface{})
								ansibleGroups := common.MapGet(markDone, "ansible_group").([]interface{})
								action := common.MapGetString(markDone, "action")
								if ansibleGroups != nil && action != "" {
									for _, ansibleGroup := range ansibleGroups {
										resetNodeAction(inventory, ansibleGroup.(string)+".hosts", action)
										logrus.Infof("重置了 ansible group %s 下 hosts 的 pangeecluster_node_action 字段值为 none", ansibleGroup)
									}
								}
							}

							if common.MapGet(step, "mark_pangeecluster_status") != nil {
								statusName := common.MapGetString(step, "mark_pangeecluster_status.name")
								statusValue := common.MapGetString(step, "mark_pangeecluster_status.value")
								if statusName != "" && statusValue != "" {
									common.MapSet(inventory, statusName, statusValue)
									logrus.Infof("设置了 inventory 中变量 %s 的值为 %s", statusName, statusValue)
								}
							}

							// 保存 inventory 文件到磁盘
							common.SaveYamlFile(tempClusterMetadata.InventoryPath, inventory)

							break // 已经找到对应的 step 并处理完毕，跳出循环
						}
					}
					break // 已经找到对应的 operation 并处理完毕，跳出循环
				}
			}

			return "\n执行成功: " + req.Cluster + "/" + req.Operation + "/" + req.Step, nil
		}
		return "\n执行失败: " + req.Cluster + "/" + req.Operation + "/" + req.Step, nil
	}

	playbook := "operations/" + req.Operation + "/" + req.Step + "/playbook.yaml"

	pid := req.Operation + "/" + req.Step + "/" + time.Now().Format("2006-01-02_15-04-05.999")

	cmd := command.Execute{
		OwnerType: "cluster",
		OwnerName: req.Cluster,
		Cmd:       "ansible-playbook",
		Args: func(execute_dir string) []string {
			result := []string{"-i", execute_dir + "/inventory.yaml", playbook}
			result = appendCommonParams(result, req, false)
			return result
		},
		Env:  []string{"ANSIBLE_CONFIG=./ansible.cfg"},
		Dir:  cluster_common.ResourcePackageDirForInventory(inventory),
		Type: req.Operation,
		PreExec: func(execute_dir string) error {
			return common.SaveYamlFile(execute_dir+"/inventory.yaml", inventory)
		},
		PostExec: postExec,
		Pid:      pid,
	}

	if err := cmd.Exec(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Faild to InstallCluster. ", err)
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
